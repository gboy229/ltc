// Package withdrawal
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package withdrawal

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	chainconsts "hotgo/addons/chain/consts"
	"hotgo/addons/chain/logic/adapters"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	"hotgo/addons/chain/utility"
	"hotgo/internal/consts"
	"hotgo/internal/model/input/assetin"
	basesvc "hotgo/internal/service"
)

const (
	// withdrawSendBatchSize 单次最多处理的待打币记录数
	withdrawSendBatchSize = 10
	// withdrawConfirmBatchSize 单次最多处理的待确认记录数
	withdrawConfirmBatchSize = 10
	// withdrawSendingTimeout 打币中状态超时时间，超时后仅告警、不自动重试
	withdrawSendingTimeout = 10 * time.Minute
)

var processMu sync.Mutex

// Process 处理提现：已审核的记录抢占后打币，已打币的记录确认上链结果并结算
func (s *sWithdrawal) Process(ctx context.Context) (err error) {
	if !processMu.TryLock() {
		g.Log().Debug(ctx, "chain withdraw process is already running, skip")
		return
	}
	defer processMu.Unlock()

	if e := s.sendWithdraw(ctx); e != nil {
		g.Log().Warningf(ctx, "chain withdraw send failed: %+v", e)
	}
	if e := s.confirmWithdraw(ctx); e != nil {
		g.Log().Warningf(ctx, "chain withdraw confirm failed: %+v", e)
	}
	if e := s.checkStuckSending(ctx); e != nil {
		g.Log().Warningf(ctx, "chain withdraw sending check failed: %+v", e)
	}
	return
}

// sendWithdraw 处理已审核的提现记录，发起链上打币
func (s *sWithdrawal) sendWithdraw(ctx context.Context) (err error) {
	cols := dao.ChainWithdraw.Columns()
	var list []*entity.ChainWithdraw
	if err = dao.ChainWithdraw.Ctx(ctx).
		Where(cols.Status, chainconsts.ChainWithdrawStatusAudited).
		OrderAsc(cols.Id).
		Limit(0, withdrawSendBatchSize).
		Scan(&list); err != nil {
		return gerror.Wrap(err, "查询待打币提现记录失败")
	}

	for _, w := range list {
		if e := s.sendOne(ctx, w); e != nil {
			g.Log().Warningf(ctx, "chain withdraw send error, id:%d err:%+v", w.Id, e)
		}
	}
	return
}

// sendOne 抢占并打币。先原子抢占 audited->sending，只有抢占成功者才发起链上转账，
// 从根上避免多实例/多轮次重复打币；即使后续更新失败，状态也停在 sending，不会被再次打币。
func (s *sWithdrawal) sendOne(ctx context.Context, w *entity.ChainWithdraw) (err error) {
	cols := dao.ChainWithdraw.Columns()

	// 1. 原子抢占（条件更新由数据库保证原子性）
	result, err := dao.ChainWithdraw.Ctx(ctx).
		WherePri(w.Id).
		Where(cols.Status, chainconsts.ChainWithdrawStatusAudited).
		Where(cols.Version, w.Version).
		Data(g.Map{
			cols.Status:  chainconsts.ChainWithdrawStatusSending,
			cols.Version: w.Version + 1,
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "抢占打币任务失败")
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		g.Log().Debugf(ctx, "chain withdraw send already handled, id:%d", w.Id)
		return nil
	}

	// 2. 抢占成功后再发起链上转账
	hash, err := s.sendTransaction(ctx, w)
	if err != nil {
		// 转账未发出：标记失败并解冻
		if f := s.markFailed(ctx, w, err.Error()); f != nil {
			g.Log().Warningf(ctx, "chain withdraw mark failed error, id:%d err:%+v", w.Id, f)
		}
		return err
	}

	// 3. 交易已广播，更新为已打币。即使该更新失败，状态也停在 sending，不会被再次打币
	_, err = dao.ChainWithdraw.Ctx(ctx).
		WherePri(w.Id).
		Where(cols.Status, chainconsts.ChainWithdrawStatusSending).
		Data(g.Map{
			cols.Status: chainconsts.ChainWithdrawStatusTransferred,
			cols.Hash:   hash,
		}).Update()
	if err != nil {
		g.Log().Errorf(ctx, "chain withdraw sent but update status failed, id:%d hash:%s err:%+v", w.Id, hash, err)
		return gerror.Wrap(err, "更新提现状态失败")
	}
	g.Log().Infof(ctx, "chain withdraw sent, id:%d hash:%s", w.Id, hash)
	return
}

// sendTransaction 组装并发送链上转账
func (s *sWithdrawal) sendTransaction(ctx context.Context, w *entity.ChainWithdraw) (hash string, err error) {
	cfg, err := s.getChainConfig(ctx, w.Chain)
	if err != nil {
		return "", err
	}
	if cfg.TransferAddress == "" || cfg.TransferPrivate == "" {
		return "", gerror.New("链配置未设置提币地址或私钥")
	}

	privateKey, err := utility.DecryptPrivateKey(ctx, cfg.TransferPrivate)
	if err != nil {
		return "", gerror.Wrap(err, "解密提币私钥失败")
	}

	contract, err := s.getContractByAddress(ctx, w.Chain, w.ContractAddress)
	if err != nil {
		return "", err
	}
	if contract == nil {
		return "", gerror.New("提币合约不存在")
	}

	adapter, err := adapters.NewChainAdapter(cfg.Driver)
	if err != nil {
		return "", gerror.Wrap(err, "不支持的链驱动")
	}
	adapter.Chain(cfg.Chain)
	if err = adapter.Client(cfg.RpcUrl); err != nil {
		return "", gerror.Wrap(err, "连接 rpc 失败")
	}

	hash, _, err = adapter.SendTransaction(adapters.TransactionOptions{
		PrivateKey:   privateKey,
		FromAddress:  cfg.TransferAddress,
		ToAddress:    w.ToAddress,
		TokenAddress: w.ContractAddress,
		Amount:       w.ReceiveAmount,
		Decimal:      int8(contract.Decimal),
	})
	if err != nil {
		return "", gerror.Wrap(err, "发起链上转账失败")
	}
	return hash, nil
}

// confirmWithdraw 处理已打币的提现记录，确认上链结果并结算
func (s *sWithdrawal) confirmWithdraw(ctx context.Context) (err error) {
	cols := dao.ChainWithdraw.Columns()
	var list []*entity.ChainWithdraw
	if err = dao.ChainWithdraw.Ctx(ctx).
		Where(cols.Status, chainconsts.ChainWithdrawStatusTransferred).
		WhereNot(cols.Hash, "").
		OrderAsc(cols.Id).
		Limit(0, withdrawConfirmBatchSize).
		Scan(&list); err != nil {
		return gerror.Wrap(err, "查询待确认提现记录失败")
	}

	for _, w := range list {
		done, e := s.checkReceipt(ctx, w)
		if e != nil {
			// 交易仍在 pending 或节点查询失败，下次再试
			g.Log().Debugf(ctx, "chain withdraw receipt pending, id:%d err:%+v", w.Id, e)
			continue
		}
		if !done {
			// 链上交易执行失败，回滚冻结
			if f := s.markFailed(ctx, w, "链上交易执行失败"); f != nil {
				g.Log().Warningf(ctx, "chain withdraw mark failed error, id:%d err:%+v", w.Id, f)
			}
			continue
		}
		if f := s.finishWithdraw(ctx, w); f != nil {
			g.Log().Warningf(ctx, "chain withdraw finish failed, id:%d err:%+v", w.Id, f)
		}
	}
	return
}

// checkReceipt 查询交易回执，true=成功，false=链上失败，err=pending/查询异常
func (s *sWithdrawal) checkReceipt(ctx context.Context, w *entity.ChainWithdraw) (done bool, err error) {
	cfg, err := s.getChainConfig(ctx, w.Chain)
	if err != nil {
		return false, err
	}
	adapter, err := adapters.NewChainAdapter(cfg.Driver)
	if err != nil {
		return false, err
	}
	adapter.Chain(cfg.Chain)
	if err = adapter.Client(cfg.RpcUrl); err != nil {
		return false, err
	}
	return adapter.Receipt(w.Hash)
}

// finishWithdraw 完成提现：先结算冻结资产（幂等），再更新状态为完成
func (s *sWithdrawal) finishWithdraw(ctx context.Context, w *entity.ChainWithdraw) (err error) {
	if err = s.settleAsset(ctx, w); err != nil {
		return err
	}

	cols := dao.ChainWithdraw.Columns()
	result, err := dao.ChainWithdraw.Ctx(ctx).
		WherePri(w.Id).
		Where(cols.Status, chainconsts.ChainWithdrawStatusTransferred).
		Where(cols.Version, w.Version).
		Data(g.Map{
			cols.Status:  chainconsts.ChainWithdrawStatusFinished,
			cols.Version: w.Version + 1,
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "更新提现状态失败")
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return nil
	}
	g.Log().Infof(ctx, "chain withdraw finished, id:%d hash:%s", w.Id, w.Hash)
	return
}

// markFailed 标记提现失败：先解冻资产（幂等），再更新状态为失败
func (s *sWithdrawal) markFailed(ctx context.Context, w *entity.ChainWithdraw, reason string) (err error) {
	if err = s.unfreezeAsset(ctx, w); err != nil {
		return err
	}

	cols := dao.ChainWithdraw.Columns()
	result, err := dao.ChainWithdraw.Ctx(ctx).
		WherePri(w.Id).
		WhereIn(cols.Status, []int{
			chainconsts.ChainWithdrawStatusAudited,
			chainconsts.ChainWithdrawStatusSending,
			chainconsts.ChainWithdrawStatusTransferred,
		}).
		Data(g.Map{
			cols.Status: chainconsts.ChainWithdrawStatusFailed,
			cols.Remark: reason,
		}).Update()
	if err != nil {
		return gerror.Wrap(err, "更新提现失败状态失败")
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return nil
	}
	g.Log().Warningf(ctx, "chain withdraw failed, id:%d reason:%s", w.Id, reason)
	return
}

// checkStuckSending 检查长时间停留在打币中的记录：只告警，不自动重试（可能已转出，需人工核对）
func (s *sWithdrawal) checkStuckSending(ctx context.Context) (err error) {
	cols := dao.ChainWithdraw.Columns()
	var list []*entity.ChainWithdraw
	if err = dao.ChainWithdraw.Ctx(ctx).
		Where(cols.Status, chainconsts.ChainWithdrawStatusSending).
		WhereLT(cols.UpdatedAt, gtime.Now().Add(-withdrawSendingTimeout)).
		Scan(&list); err != nil {
		return gerror.Wrap(err, "查询打币中超时记录失败")
	}
	for _, w := range list {
		g.Log().Errorf(ctx,
			"chain withdraw stuck in sending，请人工核对链上是否已转出, id:%d uid:%d chain:%s amount:%v",
			w.Id, w.Uid, w.Chain, w.PayAmount,
		)
	}
	return
}

// settleAsset 提现成功结算：扣减冻结金额
func (s *sWithdrawal) settleAsset(ctx context.Context, w *entity.ChainWithdraw) error {
	return basesvc.Asset().Change(ctx, &assetin.ChangeInp{
		UserId:     uint64(w.Uid),
		Account:    consts.AssetAccountUSDT,
		AddonName:  consts.AssetAddonChain,
		ChangeType: consts.AssetChangeTypeSettle,
		BizType:    consts.AssetBizTypeChainWithdraw,
		BizId:      w.Id,
		BizSn:      strconv.FormatUint(w.Id, 10),
		Num:        w.PayAmount,
		Remark:     "提现成功结算",
	})
}

// unfreezeAsset 提现失败解冻：退回冻结金额
func (s *sWithdrawal) unfreezeAsset(ctx context.Context, w *entity.ChainWithdraw) error {
	return basesvc.Asset().Change(ctx, &assetin.ChangeInp{
		UserId:     uint64(w.Uid),
		Account:    consts.AssetAccountUSDT,
		AddonName:  consts.AssetAddonChain,
		ChangeType: consts.AssetChangeTypeUnfreeze,
		BizType:    consts.AssetBizTypeChainWithdraw,
		BizId:      w.Id,
		BizSn:      strconv.FormatUint(w.Id, 10),
		Num:        w.PayAmount,
		Remark:     "提现失败解冻",
	})
}

// getChainConfig 获取链配置
func (s *sWithdrawal) getChainConfig(ctx context.Context, chain string) (*entity.ChainConfig, error) {
	cols := dao.ChainConfig.Columns()
	var cfg *entity.ChainConfig
	if err := dao.ChainConfig.Ctx(ctx).Where(cols.Chain, chain).Scan(&cfg); err != nil {
		return nil, gerror.Wrap(err, "查询链配置失败")
	}
	if cfg == nil {
		return nil, gerror.Newf("链配置不存在：%s", chain)
	}
	if cfg.RpcUrl == "" {
		return nil, gerror.New("链配置 rpc 地址为空")
	}
	return cfg, nil
}

// getContractByAddress 按链和合约地址获取合约配置
func (s *sWithdrawal) getContractByAddress(ctx context.Context, chain, address string) (*entity.ChainContract, error) {
	cols := dao.ChainContract.Columns()
	var contract *entity.ChainContract
	if err := dao.ChainContract.Ctx(ctx).
		Where(cols.Chain, chain).
		Where("LOWER("+cols.Address+") = ?", strings.ToLower(address)).
		Scan(&contract); err != nil {
		return nil, gerror.Wrap(err, "查询提币合约失败")
	}
	return contract, nil
}

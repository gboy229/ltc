// Package queues
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package queues

import (
	"context"
	"encoding/json"
	"math/big"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shopspring/decimal"

	chainconsts "hotgo/addons/chain/consts"
	chainmodel "hotgo/addons/chain/model"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	iconsts "hotgo/internal/consts"
	"hotgo/internal/library/queue"
	mainmodel "hotgo/internal/model"
	isc "hotgo/internal/service"
)

func init() {
	queue.RegisterConsumer(ScanTransfer)
}

// ScanTransfer 扫块 Transfer 事件消费者
// 扫描端只负责发现链上转账并推送；归属、幂等、落库在这里处理。
var ScanTransfer = &qScanTransfer{}

type qScanTransfer struct{}

func (q *qScanTransfer) GetTopic() string {
	return chainconsts.QueueChainScanTransfer
}

func (q *qScanTransfer) Handle(ctx context.Context, mqMsg queue.MqMsg) (err error) {
	var event chainmodel.ScanTransferEvent
	if err = json.Unmarshal(mqMsg.Body, &event); err != nil {
		return
	}

	// 1. 校验合约是否开启充币
	contractCols := dao.ChainContract.Columns()
	var contract *entity.ChainContract
	if err = dao.ChainContract.Ctx(ctx).
		Where(contractCols.Chain, event.Chain).
		Where(contractCols.Address, event.Contract).
		Scan(&contract); err != nil {
		return gerror.Wrap(err, "查询支付合约失败")
	}
	if contract == nil || contract.RechargeSwitch != chainconsts.ChainSwitchOn {
		return nil
	}

	// 2. 只处理平台收款地址的转账
	configCols := dao.ChainConfig.Columns()
	var cfg *entity.ChainConfig
	if err = dao.ChainConfig.Ctx(ctx).Where(configCols.Chain, event.Chain).Scan(&cfg); err != nil {
		return gerror.Wrap(err, "查询链配置失败")
	}
	if cfg == nil {
		return nil
	}

	// 平台收款地址：pledge 配置的 depositAddress 优先，兼容 chain_config.transfer_address
	matched := cfg.TransferAddress != "" && strings.EqualFold(cfg.TransferAddress, event.To)
	if !matched {
		if depositAddress, cacheErr := isc.SysAddonsConfig().GetConfigByKey(ctx, "pledge", "depositAddress"); cacheErr == nil {
			matched = depositAddress != "" && strings.EqualFold(depositAddress, event.To)
		}
	}
	if !matched {
		return nil
	}

	// 3. raw 金额按代币精度转十进制
	rawAmount, ok := new(big.Int).SetString(event.Amount, 10)
	if !ok {
		return gerror.Newf("非法链上金额: %s", event.Amount)
	}
	amount := decimal.NewFromBigInt(rawAmount, -int32(event.Decimals))

	// 4. 付款地址反查用户ID，查不到也落库，但业务侧不会自动匹配
	uid, cacheErr := getWalletUserId(ctx, event.From)
	if cacheErr != nil {
		g.Log().Warningf(ctx, "查询钱包用户缓存失败, from:%s err:%+v", event.From, cacheErr)
	}

	// 5. 充值记录幂等：hash + log_index
	rechargeCols := dao.ChainRecharge.Columns()
	var recharge *entity.ChainRecharge
	if err = dao.ChainRecharge.Ctx(ctx).
		Where(rechargeCols.Hash, event.TxHash).
		Where(rechargeCols.LogIndex, event.LogIndex).
		Scan(&recharge); err != nil {
		return gerror.Wrap(err, "查询充值记录失败")
	}
	if recharge == nil {
		id, insertErr := dao.ChainRecharge.Ctx(ctx).Data(g.Map{
			rechargeCols.Type:            chainconsts.ChainRechargeTypeOrder,
			rechargeCols.Uid:             int(uid),
			rechargeCols.Chain:           event.Chain,
			rechargeCols.Symbol:          strings.ToUpper(event.Symbol),
			rechargeCols.FromAddress:     event.From,
			rechargeCols.ToAddress:       event.To,
			rechargeCols.ContractAddress: event.Contract,
			rechargeCols.Amount:          amount.InexactFloat64(),
			rechargeCols.BlockHeight:     int64(event.BlockHeight),
			rechargeCols.Hash:            event.TxHash,
			rechargeCols.LogIndex:        event.LogIndex,
		}).InsertAndGetId()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "写入充值记录失败")
		}
		recharge = &entity.ChainRecharge{Id: uint64(id)}
	}

	// 6. 推送支付事件，由业务插件幂等匹配
	return queue.Push(iconsts.QueueChainPayment, &mainmodel.ChainPaymentEvent{
		RechargeId:  recharge.Id,
		Uid:         uid,
		Chain:       event.Chain,
		Symbol:      strings.ToUpper(event.Symbol),
		Contract:    event.Contract,
		From:        event.From,
		To:          event.To,
		Amount:      amount.String(),
		TxHash:      event.TxHash,
		LogIndex:    event.LogIndex,
		BlockHeight: int64(event.BlockHeight),
	})
}

// getWalletUserId 从 users 插件写入的 Redis Hash 读取钱包地址对应的用户ID。
// 这里只共享 Redis key，不 import users 插件，保持插件隔离。
func getWalletUserId(ctx context.Context, address string) (int64, error) {
	address = strings.ToLower(strings.TrimSpace(address))
	if address == "" {
		return 0, nil
	}
	value, err := g.Redis().Do(ctx, "HGET", chainconsts.WalletUidCacheKey, address)
	if err != nil {
		return 0, gerror.Wrap(err, "读取钱包用户缓存失败")
	}
	if value == nil || value.IsNil() || value.IsEmpty() {
		return 0, nil
	}
	return value.Int64(), nil
}

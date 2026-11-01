// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/chain/consts"
	"hotgo/addons/chain/logic/adapters"
	"hotgo/addons/chain/model"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	"hotgo/addons/chain/service"
	"hotgo/internal/library/queue"
)

const (
	scanBatchSize            = uint64(1000)
	scanMaxBatchesPerRun     = 5
	scanReorgRewind          = uint64(100)
	scanDefaultConfirmations = 15
)

var scanMu sync.Mutex

type sSysScan struct{}

func NewSysScan() *sSysScan {
	return &sSysScan{}
}

func init() {
	service.RegisterSysScan(NewSysScan())
}

// Scan 扫描所有开启充币的合约，把标准 Transfer 事件推送到队列
func (s *sSysScan) Scan(ctx context.Context) (err error) {
	if !scanMu.TryLock() {
		g.Log().Debug(ctx, "chain scan is already running, skip")
		return
	}
	defer scanMu.Unlock()

	configCols := dao.ChainConfig.Columns()
	var configs []*entity.ChainConfig
	if err = dao.ChainConfig.Ctx(ctx).
		Fields(
			configCols.Id,
			configCols.Driver,
			configCols.Chain,
			configCols.RpcUrl,
			configCols.Confirmations,
		).
		Scan(&configs); err != nil {
		return gerror.Wrap(err, "获取链配置失败")
	}

	for _, cfg := range configs {
		if e := s.scanChain(ctx, cfg); e != nil {
			g.Log().Warningf(ctx, "chain scan failed, chain:%s err:%+v", cfg.Chain, e)
		}
	}
	return
}

// scanChain 扫描单条链下所有开启充币的合约
func (s *sSysScan) scanChain(ctx context.Context, cfg *entity.ChainConfig) (err error) {
	if cfg.RpcUrl == "" {
		return gerror.New("rpc地址为空")
	}
	contractCols := dao.ChainContract.Columns()
	var contracts []*entity.ChainContract
	if err = dao.ChainContract.Ctx(ctx).
		Where(contractCols.Chain, cfg.Chain).
		Where(contractCols.RechargeSwitch, consts.ChainSwitchOn).
		Scan(&contracts); err != nil {
		return gerror.Wrap(err, "获取合约配置失败")
	}
	if len(contracts) == 0 {
		return
	}

	adapter, err := adapters.NewChainAdapter(cfg.Driver)
	if err != nil {
		return gerror.Wrap(err, "不支持的链驱动")
	}
	adapter.Chain(cfg.Chain)
	if err = adapter.Client(cfg.RpcUrl); err != nil {
		return gerror.Wrap(err, "连接rpc失败")
	}

	latest, err := adapter.GetBlockNumber(ctx)
	if err != nil {
		return gerror.Wrap(err, "获取最新块高失败")
	}

	confirmations := cfg.Confirmations
	if confirmations < 0 {
		confirmations = scanDefaultConfirmations
	}
	if uint64(confirmations) >= latest {
		return
	}
	target := latest - uint64(confirmations)

	for _, contract := range contracts {
		if e := s.scanContract(ctx, adapter, cfg, contract, target, confirmations); e != nil {
			g.Log().Warningf(ctx, "chain scan contract failed, chain:%s contract:%s err:%+v", cfg.Chain, contract.Address, e)
		}
	}
	return
}

// scanContract 扫描单个合约，按批次推进游标
func (s *sSysScan) scanContract(ctx context.Context, adapter adapters.ChainAdapter, cfg *entity.ChainConfig, contract *entity.ChainContract, target uint64, confirmations int) (err error) {
	if contract.Address == "" {
		return gerror.New("合约地址为空")
	}

	scanBlock := contract.ScanBlock
	if scanBlock > target {
		// 游标比目标块还高：通常是节点落后、链重置或换链，先告警不自动回退
		g.Log().Warningf(ctx, "chain scan cursor ahead of target, chain:%s contract:%s cursor:%d target:%d", cfg.Chain, contract.Address, scanBlock, target)
		return nil
	}

	// reorg 检测：校验上次扫到的块 hash 是否还在主链上
	if scanBlock > 0 && contract.ScanBlockHash != "" {
		onChainHash, e := adapter.GetBlockHash(ctx, scanBlock)
		if e != nil {
			return gerror.Wrapf(e, "校验块哈希失败, block:%d", scanBlock)
		}
		if !strings.EqualFold(onChainHash, contract.ScanBlockHash) {
			if scanBlock <= scanReorgRewind {
				return gerror.Newf("检测到 reorg，但游标过低(%d)，请人工确认", scanBlock)
			}
			scanBlock -= scanReorgRewind
			g.Log().Warningf(ctx, "chain reorg detected, chain:%s contract:%s rewind to:%d", cfg.Chain, contract.Address, scanBlock)
		}
	}

	cols := dao.ChainContract.Columns()

	for i := 0; i < scanMaxBatchesPerRun && scanBlock < target; i++ {
		var from uint64
		if scanBlock == 0 {
			// 首次扫描从当前确认块开始，避免从创世块扫
			from = target
		} else {
			from = scanBlock + 1
		}
		if from > target {
			break
		}

		to := from + scanBatchSize - 1
		if to > target {
			to = target
		}

		// 先取批次结束块 hash，推送成功后再和游标一起落库
		toHash, e := adapter.GetBlockHash(ctx, to)
		if e != nil {
			return gerror.Wrapf(e, "获取块哈希失败, block:%d", to)
		}

		transfers, e := adapter.ScanTransferLogs(ctx, from, to, []string{contract.Address})
		if e != nil {
			return gerror.Wrapf(e, "扫描 Transfer 日志失败, blocks:%d-%d", from, to)
		}

		for _, transfer := range transfers {
			event := &model.ScanTransferEvent{
				Chain:         cfg.Chain,
				Driver:        cfg.Driver,
				Contract:      contract.Address,
				Symbol:        contract.Symbol,
				Decimals:      contract.Decimal,
				From:          transfer.From,
				To:            transfer.To,
				Amount:        transfer.Amount,
				TxHash:        transfer.TxHash,
				LogIndex:      transfer.LogIndex,
				BlockHeight:   transfer.BlockHeight,
				BlockHash:     transfer.BlockHash,
				Confirmations: confirmations,
			}

			body, e := json.Marshal(event)
			if e != nil {
				return gerror.Wrap(e, "序列化扫描事件失败")
			}
			if e = queue.Push(consts.QueueChainScanTransfer, body); e != nil {
				return gerror.Wrapf(e, "推送扫描事件失败, tx:%s", transfer.TxHash)
			}
		}

		if _, e = dao.ChainContract.Ctx(ctx).WherePri(contract.Id).Data(g.Map{
			cols.ScanBlock:     to,
			cols.ScanBlockHash: toHash,
		}).Update(); e != nil {
			return gerror.Wrap(e, "更新扫描游标失败")
		}

		scanBlock = to

		time.Sleep(100 * time.Millisecond)
	}
	return
}

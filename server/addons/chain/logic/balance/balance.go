// Package balance
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package balance

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/errors/gerror"

	apiBalance "hotgo/addons/chain/api/api/balance"
	chainconsts "hotgo/addons/chain/consts"
	"hotgo/addons/chain/logic/adapters"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	"hotgo/addons/chain/service"
)

type sBalance struct{}

func NewBalance() *sBalance {
	return &sBalance{}
}

func init() {
	service.RegisterBalance(NewBalance())
}

// Get 查询钱包在支付链上的代币余额
func (s *sBalance) Get(ctx context.Context, in *apiBalance.GetReq) (res *apiBalance.GetRes, err error) {
	address := strings.TrimSpace(in.Address)
	if !common.IsHexAddress(address) {
		return nil, gerror.New("钱包地址不合法")
	}

	cfg, err := s.chainConfig(ctx, strings.TrimSpace(in.Chain))
	if err != nil {
		return nil, err
	}

	contract, err := s.contractConfig(ctx, cfg.Chain, in.Symbol)
	if err != nil {
		return nil, err
	}

	adapter, err := adapters.NewChainAdapter(cfg.Driver)
	if err != nil {
		return nil, gerror.Wrap(err, "暂不支持该链")
	}
	if err = adapter.Client(cfg.RpcUrl); err != nil {
		return nil, gerror.Wrap(err, "连接链节点失败")
	}

	amount, err := adapter.GetTokenBalance(address, contract.Address, contract.Decimal)
	if err != nil {
		return nil, gerror.Wrap(err, "查询链上余额失败")
	}

	res = &apiBalance.GetRes{
		Chain:    cfg.Chain,
		Symbol:   strings.ToUpper(contract.Symbol),
		Address:  address,
		Contract: contract.Address,
		Decimals: contract.Decimal,
		Amount:   amount,
	}
	return
}

func (s *sBalance) chainConfig(ctx context.Context, chain string) (*entity.ChainConfig, error) {
	cols := dao.ChainConfig.Columns()
	model := dao.ChainConfig.Ctx(ctx)
	if chain != "" {
		model = model.Where(cols.Chain, chain)
	} else {
		model = model.WhereNot(cols.TransferAddress, "").OrderAsc(cols.Id)
	}

	var cfg *entity.ChainConfig
	if err := model.Scan(&cfg); err != nil {
		return nil, gerror.Wrap(err, "查询链配置失败")
	}
	if cfg == nil {
		return nil, gerror.New("暂未配置支付链")
	}
	if cfg.RpcUrl == "" {
		return nil, gerror.New("链节点 RPC 地址未配置")
	}
	return cfg, nil
}

func (s *sBalance) contractConfig(ctx context.Context, chain, symbol string) (*entity.ChainContract, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		symbol = "USDT"
	}

	cols := dao.ChainContract.Columns()
	var contract *entity.ChainContract
	if err := dao.ChainContract.Ctx(ctx).
		Where(cols.Chain, chain).
		Where(cols.Symbol, symbol).
		Where(cols.RechargeSwitch, chainconsts.ChainSwitchOn).
		OrderAsc(cols.Id).
		Scan(&contract); err != nil {
		return nil, gerror.Wrap(err, "查询支付合约失败")
	}

	// 指定币种未配置时退回该链第一个开启充币的合约
	if contract == nil {
		if err := dao.ChainContract.Ctx(ctx).
			Where(cols.Chain, chain).
			Where(cols.RechargeSwitch, chainconsts.ChainSwitchOn).
			OrderAsc(cols.Id).
			Scan(&contract); err != nil {
			return nil, gerror.Wrap(err, "查询支付合约失败")
		}
	}
	if contract == nil {
		return nil, gerror.New("暂未配置支付合约")
	}
	return contract, nil
}

// Package config
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package config

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	apiConfig "hotgo/addons/chain/api/api/config"
	"hotgo/addons/chain/consts"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	"hotgo/addons/chain/service"
)

type sConfig struct{}

func NewConfig() *sConfig {
	return &sConfig{}
}

func init() {
	service.RegisterConfig(NewConfig())
}

// List 获取前台支付链配置
func (s *sConfig) List(ctx context.Context, _ *apiConfig.ListReq) (res *apiConfig.ListRes, err error) {
	cols := dao.ChainConfig.Columns()
	var configs []*entity.ChainConfig
	if err = dao.ChainConfig.Ctx(ctx).
		WhereNot(cols.TransferAddress, "").
		OrderAsc(cols.Id).
		Scan(&configs); err != nil {
		return nil, gerror.Wrap(err, "查询链配置失败")
	}

	res = new(apiConfig.ListRes)
	res.Chains = make([]*apiConfig.Chain, 0, len(configs))

	for _, cfg := range configs {
		contractCols := dao.ChainContract.Columns()
		var contracts []*entity.ChainContract
		if err = dao.ChainContract.Ctx(ctx).
			Where(contractCols.Chain, cfg.Chain).
			Where(contractCols.RechargeSwitch, consts.ChainSwitchOn).
			OrderAsc(contractCols.Id).
			Scan(&contracts); err != nil {
			return nil, gerror.Wrap(err, "查询支付合约失败")
		}
		if len(contracts) == 0 {
			continue
		}

		chainName, nativeName, nativeSymbol := chainMeta(cfg.Chain, cfg.ChainId)
		item := &apiConfig.Chain{
			Chain:     cfg.Chain,
			ChainId:   cfg.ChainId,
			ChainName: chainName,
			NativeCurrency: apiConfig.NativeCurrency{
				Name:     nativeName,
				Symbol:   nativeSymbol,
				Decimals: 18,
			},
			Contracts: make([]*apiConfig.Contract, 0, len(contracts)),
		}
		if cfg.RpcUrl != "" {
			item.RpcUrls = []string{cfg.RpcUrl}
		}
		if cfg.BrowserUrl != "" {
			item.BlockExplorerUrls = []string{cfg.BrowserUrl}
		}

		for _, contract := range contracts {
			item.Contracts = append(item.Contracts, &apiConfig.Contract{
				Symbol:   strings.ToUpper(contract.Symbol),
				Address:  contract.Address,
				Decimals: contract.Decimal,
			})
		}
		res.Chains = append(res.Chains, item)
	}
	return
}

func chainMeta(chain string, chainId int) (chainName, nativeName, nativeSymbol string) {
	switch strings.ToLower(chain) {
	case "bsc", "bnb":
		if chainId == 97 {
			return "BNB Smart Chain Testnet", "BNB", "BNB"
		}
		return "BNB Smart Chain", "BNB", "BNB"
	case "eth", "ethereum":
		return "Ethereum", "Ether", "ETH"
	case "polygon", "matic":
		return "Polygon", "MATIC", "MATIC"
	default:
		upper := strings.ToUpper(chain)
		return upper, upper, upper
	}
}

// Package config
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package config

import "github.com/gogf/gf/v2/frame/g"

// NativeCurrency 钱包添加链所需的原生币信息
type NativeCurrency struct {
	Name     string `json:"name"     dc:"原生币名称"`
	Symbol   string `json:"symbol"   dc:"原生币符号"`
	Decimals int    `json:"decimals" dc:"原生币精度"`
}

// Contract 支付合约
type Contract struct {
	Symbol   string `json:"symbol"   dc:"币种"`
	Address  string `json:"address"  dc:"合约地址"`
	Decimals int    `json:"decimals" dc:"精度"`
}

// Chain 前台钱包支付链配置
type Chain struct {
	Chain             string         `json:"chain"             dc:"链名"`
	ChainId           int            `json:"chainId"           dc:"链ID"`
	ChainName         string         `json:"chainName"         dc:"链名称"`
	NativeCurrency    NativeCurrency `json:"nativeCurrency"    dc:"原生币"`
	RpcUrls           []string       `json:"rpcUrls"           dc:"RPC地址"`
	BlockExplorerUrls []string       `json:"blockExplorerUrls" dc:"区块浏览器"`
	Contracts         []*Contract    `json:"contracts"         dc:"可用支付合约"`
}

// ListReq 获取前台支付链配置
type ListReq struct {
	g.Meta `path:"/config" method:"get" tags:"链配置" summary:"获取前台支付链配置"`
}

type ListRes struct {
	Chains []*Chain `json:"chains" dc:"支付链列表"`
}

// Package balance
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package balance

import "github.com/gogf/gf/v2/frame/g"

// GetReq 查询钱包链上代币余额
type GetReq struct {
	g.Meta  `path:"/balance" method:"get" tags:"链余额" summary:"查询钱包代币余额"`
	Chain   string `json:"chain"   dc:"链名，默认取第一条启用的链"`
	Address string `json:"address" v:"required#钱包地址不能为空" dc:"钱包地址"`
	Symbol  string `json:"symbol"  dc:"币种，默认 USDT"`
}

// GetRes 钱包链上代币余额
type GetRes struct {
	Chain    string  `json:"chain"    dc:"链名"`
	Symbol   string  `json:"symbol"   dc:"币种"`
	Address  string  `json:"address"  dc:"钱包地址"`
	Contract string  `json:"contract" dc:"合约地址"`
	Decimals int     `json:"decimals" dc:"精度"`
	Amount   float64 `json:"amount"   dc:"余额（已按精度换算）"`
}

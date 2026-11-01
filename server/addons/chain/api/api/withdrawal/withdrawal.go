// Package withdrawal
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package withdrawal

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/model/input/form"
)

// OptionsReq 获取可用提币网络
type OptionsReq struct {
	g.Meta `path:"/withdrawal/options" method:"get" tags:"链上提币" summary:"获取可用提币网络"`
}

type Option struct {
	Chain     string  `json:"chain"     dc:"链标识"`
	Network   string  `json:"network"   dc:"网络名称，如 BEP20"`
	Symbol    string  `json:"symbol"    dc:"币种"`
	Decimals  int     `json:"decimals"  dc:"精度"`
	FeeType   int     `json:"feeType"   dc:"手续费类型：0=百分比，1=固定数量"`
	Rate      float64 `json:"rate"      dc:"手续费率/固定手续费"`
	MinAmount float64 `json:"minAmount" dc:"最小提币数量"`
	MaxAmount float64 `json:"maxAmount" dc:"最大提币数量，0=不限"`
}

type OptionsRes struct {
	List []*Option `json:"list" dc:"可用提币网络"`
}

// ApplyReq 提交提币申请
type ApplyReq struct {
	g.Meta    `path:"/withdrawal/apply" method:"post" tags:"链上提币" summary:"提交提币申请"`
	Chain     string  `json:"chain"     v:"required#请选择提币网络" dc:"链标识"`
	Symbol    string  `json:"symbol"    v:"required#请选择提币币种" dc:"币种"`
	Amount    float64 `json:"amount"    v:"required" dc:"提币数量（含手续费）"`
	Nonce     string  `json:"nonce"     v:"required#签名挑战值不能为空" dc:"签名挑战值"`
	Signature string  `json:"signature" v:"required#钱包签名不能为空" dc:"钱包签名"`
}

type ApplyRes struct{}

// ListReq 用户提币记录
type ListReq struct {
	g.Meta `path:"/withdrawal/list" method:"get" tags:"链上提币" summary:"用户提币记录"`
	form.PageReq
	Status string `json:"status" dc:"状态筛选：processing/success，不传为全部"`
}

type Record struct {
	Chain         string  `json:"chain"         dc:"链标识"`
	Network       string  `json:"network"       dc:"网络名称"`
	Symbol        string  `json:"symbol"        dc:"币种"`
	Status        string  `json:"status"        dc:"状态：processing/success/failed"`
	Address       string  `json:"address"       dc:"收币地址"`
	Hash          string  `json:"hash"          dc:"交易哈希"`
	Date          string  `json:"date"          dc:"申请时间"`
	Amount        float64 `json:"amount"        dc:"提币数量"`
	Fee           float64 `json:"fee"           dc:"手续费"`
	ReceiveAmount float64 `json:"receiveAmount" dc:"到账数量"`
}

type ListRes struct {
	form.PageRes
	List []*Record `json:"list" dc:"提币记录"`
}

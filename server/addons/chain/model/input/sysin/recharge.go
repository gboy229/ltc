// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"github.com/gogf/gf/v2/os/gtime"
	"hotgo/addons/chain/model/entity"
	"hotgo/internal/model/input/form"
)

// RechargeViewInp 获取指定充值记录信息
type RechargeViewInp struct {
	Id int64 `json:"id" v:"required#充值记录ID不能为空" dc:"充值记录ID"`
}

type RechargeViewModel struct {
	entity.ChainRecharge
}

// RechargeListInp 获取充值记录列表
type RechargeListInp struct {
	form.PageReq
	Keyword         string        `json:"keyword" dc:"用户关键词：用户ID/钱包地址"`
	Type            string        `json:"type" dc:"类型"`
	Uid             int           `json:"uid" dc:"用户ID"`
	Chain           string        `json:"chain" dc:"链"`
	Symbol          string        `json:"symbol" dc:"币种"`
	FromAddress     string        `json:"fromAddress" dc:"用户钱包地址（转出）"`
	ToAddress       string        `json:"toAddress" dc:"充值地址（转入）"`
	ContractAddress string        `json:"contractAddress" dc:"合约地址"`
	Hash            string        `json:"hash" dc:"交易哈希"`
	AmountMin       float64       `json:"amountMin" dc:"最小金额"`
	AmountMax       float64       `json:"amountMax" dc:"最大金额"`
	BlockHeightMin  int64         `json:"blockHeightMin" dc:"起始块高"`
	BlockHeightMax  int64         `json:"blockHeightMax" dc:"结束块高"`
	CreatedAt       []*gtime.Time `json:"createdAt" dc:"充值时间"`
}

type RechargeListModel struct {
	entity.ChainRecharge
}

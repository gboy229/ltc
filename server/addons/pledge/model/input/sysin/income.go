// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"hotgo/addons/pledge/model/entity"
	"hotgo/internal/model/input/form"
)

// IncomeListInp 收益流水列表
type IncomeListInp struct {
	form.PageReq
	UserId        uint64 `json:"userId"        dc:"收益用户ID"`
	OrderId       uint64 `json:"orderId"       dc:"订单ID"`
	WalletAddress string `json:"walletAddress" dc:"钱包地址"`
	Type          string `json:"type"          dc:"收益类型"`
	Status        int    `json:"status"        dc:"状态"`
}

type IncomeListModel struct {
	entity.PledgeIncomeLog
	WalletAddress string `json:"walletAddress" dc:"钱包地址"`
}

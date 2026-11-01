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

// OrderListInp 订单列表
type OrderListInp struct {
	form.PageReq
	UserId        uint64 `json:"userId"        dc:"用户ID"`
	OrderNo       string `json:"orderNo"       dc:"订单号"`
	WalletAddress string `json:"walletAddress" dc:"钱包地址"`
	Status        int    `json:"status"        dc:"状态"`
	IsEarlyRefund *int   `json:"isEarlyRefund" dc:"是否提前退本"`
}

type OrderListModel struct {
	entity.PledgeOrder
	WalletAddress string `json:"walletAddress" dc:"钱包地址"`
}

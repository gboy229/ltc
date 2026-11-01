// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"hotgo/addons/pledge/model/input/sysin"
)

type (
	ISysOrder interface {
		// List 订单列表
		List(ctx context.Context, in *sysin.OrderListInp) (list []*sysin.OrderListModel, totalCount int, err error)
	}
)

var (
	localSysOrder ISysOrder
)

func SysOrder() ISysOrder {
	if localSysOrder == nil {
		panic("implement not found for interface ISysOrder, forgot register?")
	}
	return localSysOrder
}

func RegisterSysOrder(i ISysOrder) {
	localSysOrder = i
}

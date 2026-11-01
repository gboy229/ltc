// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"hotgo/addons/pledge/api/api/order"
	mainmodel "hotgo/internal/model"
)

type (
	IOrder interface {
		// Buy 提交报单
		Buy(ctx context.Context, userId int64, in *order.BuyReq) (res *order.BuyRes, err error)
		// List 报单列表
		List(ctx context.Context, userId int64, in *order.ListReq) (list []*order.ListItem, totalCount int, err error)
		// Statistics 报单统计
		Statistics(ctx context.Context, userId int64, in *order.StatisticsReq) (res *order.StatisticsRes, err error)
		// Refund 提前退本
		Refund(ctx context.Context, userId int64, in *order.RefundReq) (res *order.RefundRes, err error)
		// Status 查询报单状态
		Status(ctx context.Context, userId int64, in *order.StatusReq) (res *order.StatusRes, err error)
		// HandlePayment 处理链上支付事件
		HandlePayment(ctx context.Context, event *mainmodel.ChainPaymentEvent) (err error)
	}
)

var (
	localOrder IOrder
)

func Order() IOrder {
	if localOrder == nil {
		panic("implement not found for interface IOrder, forgot register?")
	}
	return localOrder
}

func RegisterOrder(i IOrder) {
	localOrder = i
}

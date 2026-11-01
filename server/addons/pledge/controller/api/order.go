// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"

	"hotgo/addons/pledge/api/api/order"
	"hotgo/addons/pledge/service"
	"hotgo/internal/library/contexts"
)

var (
	Order = cOrder{}
)

type cOrder struct{}

// Buy 提交报单
func (c *cOrder) Buy(ctx context.Context, req *order.BuyReq) (res *order.BuyRes, err error) {
	return service.Order().Buy(ctx, contexts.GetUserId(ctx), req)
}

// List 报单列表
func (c *cOrder) List(ctx context.Context, req *order.ListReq) (res *order.ListRes, err error) {
	list, totalCount, err := service.Order().List(ctx, contexts.GetUserId(ctx), req)
	if err != nil {
		return
	}

	res = new(order.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Statistics 报单统计
func (c *cOrder) Statistics(ctx context.Context, req *order.StatisticsReq) (res *order.StatisticsRes, err error) {
	return service.Order().Statistics(ctx, contexts.GetUserId(ctx), req)
}

// Refund 提前退本
func (c *cOrder) Refund(ctx context.Context, req *order.RefundReq) (res *order.RefundRes, err error) {
	return service.Order().Refund(ctx, contexts.GetUserId(ctx), req)
}

// Status 查询报单状态
func (c *cOrder) Status(ctx context.Context, req *order.StatusReq) (res *order.StatusRes, err error) {
	return service.Order().Status(ctx, contexts.GetUserId(ctx), req)
}

// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/pledge/api/admin/order"
	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/addons/pledge/service"
)

var (
	Order = cOrder{}
)

type cOrder struct{}

func (c *cOrder) List(ctx context.Context, req *order.ListReq) (res *order.ListRes, err error) {
	list, totalCount, err := service.SysOrder().List(ctx, &req.OrderListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.OrderListModel{}
	}
	res = new(order.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

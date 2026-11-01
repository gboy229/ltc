// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"

	"hotgo/addons/pledge/api/api/income"
	"hotgo/addons/pledge/service"
	"hotgo/internal/library/contexts"
)

var (
	Income = cIncome{}
)

type cIncome struct{}

// List 收益流水列表
func (c *cIncome) List(ctx context.Context, req *income.ListReq) (res *income.ListRes, err error) {
	list, totalCount, err := service.Income().List(ctx, uint64(contexts.GetUserId(ctx)), req)
	if err != nil {
		return
	}
	if list == nil {
		list = []*income.ListItem{}
	}
	res = new(income.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

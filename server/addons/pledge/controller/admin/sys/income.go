// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/pledge/api/admin/income"
	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/addons/pledge/service"
)

var (
	Income = cIncome{}
)

type cIncome struct{}

func (c *cIncome) List(ctx context.Context, req *income.ListReq) (res *income.ListRes, err error) {
	list, totalCount, err := service.SysIncome().List(ctx, &req.IncomeListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.IncomeListModel{}
	}
	res = new(income.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

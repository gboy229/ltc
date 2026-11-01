// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/chain/api/admin/recharge"
	"hotgo/addons/chain/service"
)

var (
	Recharge = cRecharge{}
)

type cRecharge struct{}

// View 获取指定充值记录信息
func (c *cRecharge) View(ctx context.Context, req *recharge.ViewReq) (res *recharge.ViewRes, err error) {
	data, err := service.SysRecharge().View(ctx, &req.RechargeViewInp)
	if err != nil {
		return
	}

	res = new(recharge.ViewRes)
	res.RechargeViewModel = data
	return
}

// List 查看充值记录列表
func (c *cRecharge) List(ctx context.Context, req *recharge.ListReq) (res *recharge.ListRes, err error) {
	list, totalCount, err := service.SysRecharge().List(ctx, &req.RechargeListInp)
	if err != nil {
		return
	}

	res = new(recharge.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

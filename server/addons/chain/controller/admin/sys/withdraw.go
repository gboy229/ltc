// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/chain/api/admin/withdraw"
	"hotgo/addons/chain/service"
)

var (
	Withdraw = cWithdraw{}
)

type cWithdraw struct{}

// Status 提现审核
func (c *cWithdraw) Status(ctx context.Context, req *withdraw.StatusReq) (res *withdraw.StatusRes, err error) {
	err = service.SysWithdraw().Status(ctx, &req.WithdrawStatusInp)
	return
}

// View 获取指定提现记录信息
func (c *cWithdraw) View(ctx context.Context, req *withdraw.ViewReq) (res *withdraw.ViewRes, err error) {
	data, err := service.SysWithdraw().View(ctx, &req.WithdrawViewInp)
	if err != nil {
		return
	}

	res = new(withdraw.ViewRes)
	res.WithdrawViewModel = data
	return
}

// List 查看提现记录列表
func (c *cWithdraw) List(ctx context.Context, req *withdraw.ListReq) (res *withdraw.ListRes, err error) {
	list, totalCount, err := service.SysWithdraw().List(ctx, &req.WithdrawListInp)
	if err != nil {
		return
	}

	res = new(withdraw.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

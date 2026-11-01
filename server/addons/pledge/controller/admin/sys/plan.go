// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/addons/pledge/api/admin/plan"
	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/addons/pledge/service"
)

var (
	Plan = cPlan{}
)

type cPlan struct{}

// Delete 删除
func (c *cPlan) Delete(ctx context.Context, req *plan.DeleteReq) (res *plan.DeleteRes, err error) {
	err = service.SysPlan().Delete(ctx, &req.PlanDeleteInp)
	return
}

// Edit 更新
func (c *cPlan) Edit(ctx context.Context, req *plan.EditReq) (res *plan.EditRes, err error) {
	err = service.SysPlan().Edit(ctx, &req.PlanEditInp)
	return
}

// View 获取指定信息
func (c *cPlan) View(ctx context.Context, req *plan.ViewReq) (res *plan.ViewRes, err error) {
	data, err := service.SysPlan().View(ctx, &req.PlanViewInp)
	if err != nil {
		return
	}

	res = new(plan.ViewRes)
	res.PlanViewModel = data
	return
}

// List 查看列表
func (c *cPlan) List(ctx context.Context, req *plan.ListReq) (res *plan.ListRes, err error) {
	list, totalCount, err := service.SysPlan().List(ctx, &req.PlanListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.PlanListModel{}
	}

	res = new(plan.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Status 更新状态
func (c *cPlan) Status(ctx context.Context, req *plan.StatusReq) (res *plan.StatusRes, err error) {
	err = service.SysPlan().Status(ctx, &req.PlanStatusInp)
	return
}

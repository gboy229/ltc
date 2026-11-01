// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/pledge/api/admin/periodconfig"
	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/addons/pledge/service"
)

var (
	PeriodConfig = cPeriodConfig{}
)

type cPeriodConfig struct{}

func (c *cPeriodConfig) Delete(ctx context.Context, req *periodconfig.DeleteReq) (res *periodconfig.DeleteRes, err error) {
	err = service.SysPeriodConfig().Delete(ctx, &req.PeriodConfigDeleteInp)
	return
}

func (c *cPeriodConfig) Edit(ctx context.Context, req *periodconfig.EditReq) (res *periodconfig.EditRes, err error) {
	err = service.SysPeriodConfig().Edit(ctx, &req.PeriodConfigEditInp)
	return
}

func (c *cPeriodConfig) Status(ctx context.Context, req *periodconfig.StatusReq) (res *periodconfig.StatusRes, err error) {
	err = service.SysPeriodConfig().Status(ctx, &req.PeriodConfigStatusInp)
	return
}

func (c *cPeriodConfig) View(ctx context.Context, req *periodconfig.ViewReq) (res *periodconfig.ViewRes, err error) {
	data, err := service.SysPeriodConfig().View(ctx, &req.PeriodConfigViewInp)
	if err != nil {
		return
	}
	res = new(periodconfig.ViewRes)
	res.PeriodConfigViewModel = data
	return
}

func (c *cPeriodConfig) List(ctx context.Context, req *periodconfig.ListReq) (res *periodconfig.ListRes, err error) {
	list, totalCount, err := service.SysPeriodConfig().List(ctx, &req.PeriodConfigListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.PeriodConfigListModel{}
	}
	res = new(periodconfig.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/addons/users/api/admin/level"
	"hotgo/addons/users/service"
)

var (
	Level = cLevel{}
)

type cLevel struct{}

// Delete 删除
func (c *cLevel) Delete(ctx context.Context, req *level.DeleteReq) (res *level.DeleteRes, err error) {
	err = service.SysLevel().Delete(ctx, &req.LevelDeleteInp)
	return
}

// Edit 更新
func (c *cLevel) Edit(ctx context.Context, req *level.EditReq) (res *level.EditRes, err error) {
	err = service.SysLevel().Edit(ctx, &req.LevelEditInp)
	return
}

// MaxSort 最大排序
func (c *cLevel) MaxSort(ctx context.Context, req *level.MaxSortReq) (res *level.MaxSortRes, err error) {
	res = new(level.MaxSortRes)
	res.LevelMaxSortModel, err = service.SysLevel().MaxSort(ctx, &req.LevelMaxSortInp)
	return
}

// View 获取指定信息
func (c *cLevel) View(ctx context.Context, req *level.ViewReq) (res *level.ViewRes, err error) {
	data, err := service.SysLevel().View(ctx, &req.LevelViewInp)
	if err != nil {
		return
	}

	res = new(level.ViewRes)
	res.LevelViewModel = data
	return
}

// List 查看列表
func (c *cLevel) List(ctx context.Context, req *level.ListReq) (res *level.ListRes, err error) {
	list, totalCount, err := service.SysLevel().List(ctx, &req.LevelListInp)
	if err != nil {
		return
	}

	res = new(level.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Status 更新状态
func (c *cLevel) Status(ctx context.Context, req *level.StatusReq) (res *level.StatusRes, err error) {
	err = service.SysLevel().Status(ctx, &req.LevelStatusInp)
	return
}

// Option 获取等级选项
func (c *cLevel) Option(ctx context.Context, req *level.OptionReq) (res *level.OptionRes, err error) {
	res = new(level.OptionRes)
	res.List, err = service.SysLevel().Option(ctx, &req.LevelOptionInp)
	return
}

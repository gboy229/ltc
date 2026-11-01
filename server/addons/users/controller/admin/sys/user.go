// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/addons/users/api/admin/user"
	"hotgo/addons/users/service"
)

var (
	User = cUser{}
)

type cUser struct{}

// Delete 删除
func (c *cUser) Delete(ctx context.Context, req *user.DeleteReq) (res *user.DeleteRes, err error) {
	err = service.SysUser().Delete(ctx, &req.UserDeleteInp)
	return
}

// Edit 更新
func (c *cUser) Edit(ctx context.Context, req *user.EditReq) (res *user.EditRes, err error) {
	data, err := service.SysUser().Edit(ctx, &req.UserEditInp)
	if err != nil {
		return
	}

	res = new(user.EditRes)
	res.UserEditModel = data
	return
}

// View 获取指定信息
func (c *cUser) View(ctx context.Context, req *user.ViewReq) (res *user.ViewRes, err error) {
	data, err := service.SysUser().View(ctx, &req.UserViewInp)
	if err != nil {
		return
	}

	res = new(user.ViewRes)
	res.UserViewModel = data
	return
}

// List 查看列表
func (c *cUser) List(ctx context.Context, req *user.ListReq) (res *user.ListRes, err error) {
	list, totalCount, err := service.SysUser().List(ctx, &req.UserListInp)
	if err != nil {
		return
	}

	res = new(user.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Status 更新状态
func (c *cUser) Status(ctx context.Context, req *user.StatusReq) (res *user.StatusRes, err error) {
	err = service.SysUser().Status(ctx, &req.UserStatusInp)
	return
}

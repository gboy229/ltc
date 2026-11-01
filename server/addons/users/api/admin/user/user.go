// Package user
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package user

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/users/model/input/sysin"
	"hotgo/internal/model/input/form"
)

// ListReq 查询用户列表
type ListReq struct {
	g.Meta `path:"/user/list" method:"get" tags:"用户管理" summary:"获取用户列表"`
	sysin.UserListInp
}

type ListRes struct {
	List []*sysin.UserListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// ViewReq 获取指定信息
type ViewReq struct {
	g.Meta `path:"/user/view" method:"get" tags:"用户管理" summary:"获取指定用户信息"`
	sysin.UserViewInp
}

type ViewRes struct {
	*sysin.UserViewModel
}

// EditReq 修改/新增
type EditReq struct {
	g.Meta `path:"/user/edit" method:"post" tags:"用户管理" summary:"修改/新增用户"`
	sysin.UserEditInp
}

type EditRes struct {
	*sysin.UserEditModel
}

// DeleteReq 删除
type DeleteReq struct {
	g.Meta `path:"/user/delete" method:"post" tags:"用户管理" summary:"删除用户"`
	sysin.UserDeleteInp
}

type DeleteRes struct{}

// StatusReq 更新状态
type StatusReq struct {
	g.Meta `path:"/user/status" method:"post" tags:"用户管理" summary:"更新用户状态"`
	sysin.UserStatusInp
}

type StatusRes struct{}

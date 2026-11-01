// Package middleware
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package middleware

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/net/ghttp"

	usersconsts "hotgo/addons/users/consts"
	"hotgo/addons/users/library/contexts"
	"hotgo/addons/users/model/dao"
	"hotgo/addons/users/model/entity"
	"hotgo/internal/consts"
	"hotgo/internal/library/response"
	"hotgo/internal/library/token"
)

// Auth 插件 API 鉴权中间件：解析 token 并加载用户信息到插件上下文
func Auth(r *ghttp.Request) {
	// 解析登录 token，获取用户ID
	identity, err := token.ParseLoginUser(r)
	if err != nil {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), err.Error())
		return
	}

	// 加载用户信息
	var user *entity.UsersUser
	if err = dao.UsersUser.Ctx(r.Context()).WherePri(identity.Id).Scan(&user); err != nil {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "获取用户信息失败")
		return
	}
	if user == nil {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "用户不存在")
		return
	}
	if user.Status != consts.StatusEnabled {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "账号已被禁用")
		return
	}

	// 加载当前用户钱包地址
	var account *entity.UsersUserAccount
	if err = dao.UsersUserAccount.Ctx(r.Context()).
		Where(dao.UsersUserAccount.Columns().UserId, user.Id).
		Where(dao.UsersUserAccount.Columns().Type, usersconsts.UserAuthTypeWallet).
		Scan(&account); err != nil {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "获取钱包地址失败")
		return
	}
	wallet := ""
	if account != nil {
		wallet = account.Account
	}

	// 填充插件用户上下文
	r.SetCtx(contexts.SetUsersContext(r.Context(), user, wallet))
	r.Middleware.Next()
}

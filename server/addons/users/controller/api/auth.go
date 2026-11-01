// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"
	"hotgo/addons/users/api/api/auth"
	"hotgo/addons/users/service"
)

var (
	Auth = cAuth{}
)

type cAuth struct{}

// Nonce 获取钱包签名挑战值
func (c *cAuth) Nonce(ctx context.Context, req *auth.NonceReq) (res *auth.NonceRes, err error) {
	return service.Auth().Nonce(ctx, req)
}

// Login 统一登录
func (c *cAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	return service.Auth().Login(ctx, req)
}

// Logout 退出登录
func (c *cAuth) Logout(ctx context.Context, req *auth.LogoutReq) (res *auth.LogoutRes, err error) {
	return service.Auth().Logout(ctx, req)
}

// Status 查询账号注册状态
func (c *cAuth) Status(ctx context.Context, req *auth.StatusReq) (res *auth.StatusRes, err error) {
	return service.Auth().Status(ctx, req)
}

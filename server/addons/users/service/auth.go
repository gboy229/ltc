// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	apiAuth "hotgo/addons/users/api/api/auth"
)

type (
	IAuth interface {
		// Nonce 获取钱包签名挑战值
		Nonce(ctx context.Context, in *apiAuth.NonceReq) (res *apiAuth.NonceRes, err error)
		// Status 查询账号注册状态
		Status(ctx context.Context, in *apiAuth.StatusReq) (res *apiAuth.StatusRes, err error)
		// Login 统一登录
		Login(ctx context.Context, in *apiAuth.LoginReq) (res *apiAuth.LoginRes, err error)
		// Logout 退出登录
		Logout(ctx context.Context, in *apiAuth.LogoutReq) (res *apiAuth.LogoutRes, err error)
	}
)

var (
	localAuth IAuth
)

func Auth() IAuth {
	if localAuth == nil {
		panic("implement not found for interface IAuth, forgot register?")
	}
	return localAuth
}

func RegisterAuth(i IAuth) {
	localAuth = i
}

// Package contexts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package contexts

import (
	"context"

	"hotgo/addons/users/model/entity"
)

// UsersContextKey 用户上下文键名
const UsersContextKey = "users_context"

// UsersContext 用户上下文结构
type UsersContext struct {
	User   *entity.UsersUser `json:"user"`   // 用户信息
	Wallet string            `json:"wallet"` // 当前登录钱包地址
}

// SetUsersContext 设置用户上下文
func SetUsersContext(ctx context.Context, user *entity.UsersUser, wallet string) context.Context {
	usersCtx := &UsersContext{
		User:   user,
		Wallet: wallet,
	}
	return context.WithValue(ctx, UsersContextKey, usersCtx)
}

// GetUsersContext 获取用户上下文
func GetUsersContext(ctx context.Context) *UsersContext {
	value := ctx.Value(UsersContextKey)
	if value == nil {
		return nil
	}
	if usersCtx, ok := value.(*UsersContext); ok {
		return usersCtx
	}
	return nil
}

// GetUsersUser 获取用户信息
func GetUsersUser(ctx context.Context) *entity.UsersUser {
	usersCtx := GetUsersContext(ctx)
	if usersCtx == nil {
		return nil
	}
	return usersCtx.User
}

// GetUsersUserId 获取用户ID
func GetUsersUserId(ctx context.Context) int64 {
	user := GetUsersUser(ctx)
	if user == nil {
		return 0
	}
	return user.Id
}

// GetUsersWallet 获取当前登录钱包地址
func GetUsersWallet(ctx context.Context) string {
	usersCtx := GetUsersContext(ctx)
	if usersCtx == nil {
		return ""
	}
	return usersCtx.Wallet
}

// IsUsersLoggedIn 检查用户是否已登录
func IsUsersLoggedIn(ctx context.Context) bool {
	return GetUsersUser(ctx) != nil
}

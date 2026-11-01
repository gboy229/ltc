// Package router
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"hotgo/addons/users/controller/api"
	"hotgo/addons/users/global"
	"hotgo/addons/users/logic/middleware"
	"hotgo/internal/consts"
	"hotgo/internal/library/addons"
)

// Api 前台路由
func Api(ctx context.Context, group *ghttp.RouterGroup) {
	prefix := addons.RouterPrefix(ctx, consts.AppApi, global.GetSkeleton().Name)
	group.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Bind(
			// 无需验证的路由
			api.Auth, // 认证入口：nonce/login/logout
		)
		group.Middleware(middleware.Auth)
		group.Bind(
			// 需要验证的路由
			api.Center, // 用户中心
			api.Asset,  // 用户资产
			api.Team,   // 我的团队
		)
	})
}

// Package pledge
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package pledge

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	_ "hotgo/addons/pledge/consts"
	_ "hotgo/addons/pledge/crons"
	"hotgo/addons/pledge/global"
	_ "hotgo/addons/pledge/logic"
	_ "hotgo/addons/pledge/queues"
	"hotgo/addons/pledge/router"
	"hotgo/internal/library/addons"
	"hotgo/internal/library/cron"
	"hotgo/internal/library/menu"
	"hotgo/internal/service"
	"sync"
)

type module struct {
	skeleton *addons.Skeleton
	ctx      context.Context
	sync.Mutex
}

func init() {
	newModule()
}

func newModule() {
	m := &module{
		skeleton: &addons.Skeleton{
			Label:       `质押管理`,
			Name:        `pledge`,
			Group:       1,
			Logo:        "",
			Brief:       `投资挡位与套餐配置`,
			Description: `管理质押投资套餐、周期收益、直推收益、每日交易手数与提前退本违约金配置`,
			Author:      ``,
			Version:     `v1.0.4`, // 当该版本号高于已安装的版本号时，会提示可以更新
		},
		ctx: gctx.New(),
	}
	addons.RegisterModule(m)
}

// Start 启动模块
func (m *module) Start(option *addons.Option) (err error) {
	// 初始化模块
	global.Init(m.ctx, m.skeleton)

	// 注册插件路由
	option.Server.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(service.Middleware().Addon)
		router.Admin(m.ctx, group)
		router.Api(m.ctx, group)
	})
	return
}

// Stop 停止模块
func (m *module) Stop() (err error) {
	return
}

// Ctx 上下文
func (m *module) Ctx() context.Context {
	return m.ctx
}

// GetSkeleton 获取模块
func (m *module) GetSkeleton() *addons.Skeleton {
	return m.skeleton
}

// Install 安装模块
func (m *module) Install(ctx context.Context) (err error) {
	if err = menu.SyncPlugin(ctx, m.skeleton.Name); err != nil {
		return
	}
	return cron.SyncPlugin(ctx, m.skeleton.Name)
}

// Upgrade 更新模块
func (m *module) Upgrade(ctx context.Context) (err error) {
	return m.Install(ctx)
}

// UnInstall 卸载模块
func (m *module) UnInstall(ctx context.Context) (err error) {
	if err = menu.DisablePlugin(ctx, m.skeleton.Name); err != nil {
		return
	}
	return cron.DeletePlugin(ctx, m.skeleton.Name)
}

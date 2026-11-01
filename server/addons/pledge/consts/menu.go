// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import "hotgo/internal/library/menu"

// init 注册插件菜单
func init() {
	menu.RegisterPlugin("pledge",
		&menu.Item{Name: "pledge", Title: "质押管理", Path: "/pledge", Icon: "AppstoreOutlined", Dir: true, Sort: 200,
			Children: []*menu.Item{
				{
					Name: "pledge_stat", Title: "数据统计", Path: "stat",
					Component: "/addons/pledge/stat/index", Api: "pledge/stat",
					Perm: "/pledge/stat",
				},
				{
					Name: "pledge_period_config", Title: "周期配置", Path: "periodconfig",
					Component: "/addons/pledge/periodconfig/index", Api: "pledge/periodConfig",
					Perm: "/pledge/periodConfig/list,/pledge/periodConfig/view",
					Ops:  "edit,delete,status",
				},
				{
					Name: "pledge_plan", Title: "套餐管理", Path: "plan",
					Component: "/addons/pledge/plan/index", Api: "pledge/plan",
					Perm: "/pledge/plan/list,/pledge/plan/view",
					Ops:  "edit,delete,status",
				},
				{
					Name: "pledge_order", Title: "订单管理", Path: "order",
					Component: "/addons/pledge/order/index", Api: "pledge/order",
					Perm: "/pledge/order/list",
				},
				{
					Name: "pledge_income", Title: "收益流水", Path: "income",
					Component: "/addons/pledge/income/index", Api: "pledge/income",
					Perm: "/pledge/income/list",
				},
				{
					Name: "pledge_config", Title: "配置管理", Path: "config",
					Component: "/addons/pledge/config/system",
					Perm:      "/pledge/config/get,/pledge/config/items",
					Buttons: []*menu.Button{
						{Name: "pledge_config_get", Title: "查看配置", Permissions: "/pledge/config/get"},
						{Name: "pledge_config_items", Title: "查看全部配置", Permissions: "/pledge/config/items"},
						{Name: "pledge_config_update", Title: "更新配置", Permissions: "/pledge/config/update"},
					},
				},
			},
		},
	)
}

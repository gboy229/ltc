// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import "hotgo/internal/library/menu"

// init 注册插件菜单（顶级目录）
func init() {
	menu.RegisterPlugin("users",
		&menu.Item{Name: "users", Title: "用户管理", Path: "/users", Icon: "AppstoreOutlined", Dir: true, Sort: 100,
			Children: []*menu.Item{

				{Name: "users_user", Title: "用户管理", Path: "user", Component: "/addons/users/user/index", Api: "users/user", Perm: "/users/user/list,/users/level/option", Ops: "edit,status"},
				{Name: "users_level", Title: "等级管理", Path: "level", Component: "/addons/users/level/index", Api: "users/level", Ops: "edit,delete,status"},
				{Name: "users_asset", Title: "资金管理", Path: "asset", Dir: true,
					Children: []*menu.Item{
						{Name: "users_asset_account", Title: "资产账户", Path: "account", Component: "/addons/users/asset/account/index", Perm: "/users/asset/accountList"},
						{Name: "users_asset_log", Title: "资金流水", Path: "log", Component: "/addons/users/asset/log/index", Perm: "/users/asset/logList"},
						{Name: "users_asset_stat", Title: "资产统计", Path: "stat", Component: "/addons/users/asset/stat/index", Perm: "/users/asset/stat"},
					}},
				{Name: "users_config", Title: "配置管理", Path: "config", Component: "/addons/users/config/system", Perm: "/users/config/get",
					Buttons: []*menu.Button{
						{Name: "users_config_get", Title: "查看配置", Permissions: "/users/config/get"},
						{Name: "users_config_update", Title: "更新配置", Permissions: "/users/config/update"},
					}},
			},
		},
	)
}

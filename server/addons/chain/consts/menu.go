// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import "hotgo/internal/library/menu"

// init 注册插件菜单（顶级目录）
func init() {
	menu.RegisterPlugin("chain",
		&menu.Item{Name: "chain", Title: "区块管理", Path: "/chain", Icon: "AppstoreOutlined", Dir: true, Sort: 900,
			Children: []*menu.Item{
				{Name: "chain_stat", Title: "数据统计", Path: "stat", Component: "/addons/chain/stat/index", Api: "chain/stat"},
				{Name: "chain_config", Title: "网络配置", Path: "config", Component: "/addons/chain/config/index", Api: "chain/config", Ops: "edit,delete",
					Buttons: []*menu.Button{
						{Name: "chain_config_generateAddress", Title: "生成提币地址", Permissions: "/chain/config/generateAddress"},
					}},
				{Name: "chain_contract", Title: "合约管理", Path: "contract", Component: "/addons/chain/contract/index", Api: "chain/contract", Ops: "edit,delete"},
				{Name: "chain_recharge", Title: "充值记录", Path: "recharge", Component: "/addons/chain/recharge/index", Api: "chain/recharge"},
				{Name: "chain_withdraw", Title: "提现记录", Path: "withdraw", Component: "/addons/chain/withdraw/index", Api: "chain/withdraw", Ops: "status"},
			},
		},
	)
}

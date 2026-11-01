// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import "hotgo/internal/library/menu"

// init 注册插件菜单
func init() {
	menu.RegisterPlugin("cms",
		&menu.Item{Name: "cms", Title: "内容管理", Path: "/cms", Icon: "AppstoreOutlined", Dir: true, Sort: 1000,
			Children: []*menu.Item{
				{Name: "cms_article", Title: "公告管理", Path: "article", Component: "/addons/cms/article/index", Api: "cms/article", Ops: "edit,delete,status"},
			},
		},
	)
}

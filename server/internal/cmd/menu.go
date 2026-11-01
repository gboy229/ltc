// Package cmd
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package cmd

import (
	"context"
	"github.com/gogf/gf/v2/os/gcmd"
	"hotgo/internal/library/menu"
)

var (
	Menu = &gcmd.Command{
		Name:        "menu",
		Brief:       "菜单同步",
		Description: `将代码中注册的菜单同步到数据库，自动识别新增、修改、删除（以 name 为唯一键）。`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return menu.Sync(ctx)
		},
	}
)

// Package cmd
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package cmd

import (
	"context"
	"github.com/gogf/gf/v2/os/gcmd"
	"hotgo/internal/library/migrate"
)

var (
	Migrate = &gcmd.Command{
		Name:        "migrate",
		Brief:       "数据库结构迁移",
		Description: `执行数据库结构迁移（建表/改表），按版本号只执行未应用的迁移。`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return migrate.Up(ctx, migrate.DefaultScope, migrate.DefaultDir)
		},
	}
)

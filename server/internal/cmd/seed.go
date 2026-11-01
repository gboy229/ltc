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
	Seed = &gcmd.Command{
		Name:        "seed",
		Brief:       "数据库数据填充",
		Description: `填充数据库种子数据，幂等、可重复执行。`,
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			return migrate.Seed(ctx, migrate.DefaultSeedDir)
		},
	}
)

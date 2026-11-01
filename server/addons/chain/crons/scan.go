// Package crons
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package crons

import (
	"context"

	"hotgo/addons/chain/service"
	"hotgo/internal/consts"
	"hotgo/internal/library/cron"
)

func init() {
	cron.Register(ScanTask)
	cron.RegisterPlugin("chain", &cron.CronItem{
		Name:    "chain_scan",
		Title:   "区块扫描",
		Pattern: "*/5 * * * * *", // 每 5 秒
		Policy:  consts.CronPolicySingle,
		GroupId: 1,
		Sort:    10,
		Remark:  "扫描 ERC20 充值 Transfer 事件并推送到 Redis 队列",
		Status:  1,
	})
}

// ScanTask 区块扫描任务
var ScanTask = &cScanTask{name: "chain_scan"}

type cScanTask struct{ name string }

func (c *cScanTask) GetName() string { return c.name }

func (c *cScanTask) Execute(ctx context.Context, parser *cron.Parser) (err error) {
	return service.SysScan().Scan(ctx)
}

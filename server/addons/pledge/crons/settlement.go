// Package crons
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package crons

import (
	"context"

	"hotgo/addons/pledge/service"
	"hotgo/internal/consts"
	"hotgo/internal/library/cron"
)

func init() {
	cron.Register(SettlementTask)
	cron.RegisterPlugin("pledge", &cron.CronItem{
		Name:    "pledge_settlement",
		Title:   "质押静态结算",
		Pattern: "0 * * * * *", // 每分钟第 0 秒，单次处理一批，靠频率消化积压
		Policy:  consts.CronPolicySingle,
		GroupId: 1,
		Sort:    10,
		Remark:  "结算进行中订单的静态收益，并处理到期本金返还",
		Status:  1,
	})
}

var SettlementTask = &cSettlementTask{name: "pledge_settlement"}

type cSettlementTask struct{ name string }

func (c *cSettlementTask) GetName() string { return c.name }

func (c *cSettlementTask) Execute(ctx context.Context, parser *cron.Parser) (err error) {
	return service.Settlement().Run(ctx)
}

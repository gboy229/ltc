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
	cron.Register(WithdrawTask)
	cron.RegisterPlugin("chain", &cron.CronItem{
		Name:    "chain_withdraw",
		Title:   "提币处理",
		Pattern: "*/30 * * * * *", // 每 30 秒
		Policy:  consts.CronPolicySingle,
		GroupId: 1,
		Sort:    20,
		Remark:  "处理已审核提现：链上打币、确认上链结果并结算资产",
		Status:  1,
	})
}

// WithdrawTask 提币处理任务
var WithdrawTask = &cWithdrawTask{name: "chain_withdraw"}

type cWithdrawTask struct{ name string }

func (c *cWithdrawTask) GetName() string { return c.name }

func (c *cWithdrawTask) Execute(ctx context.Context, parser *cron.Parser) (err error) {
	return service.Withdrawal().Process(ctx)
}

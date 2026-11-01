// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"

	"hotgo/addons/chain/api/api/withdrawal"
	"hotgo/addons/chain/service"
	"hotgo/internal/library/contexts"
	basesvc "hotgo/internal/service"
)

var (
	Withdrawal = cWithdrawal{}
)

type cWithdrawal struct{}

// Options 获取可用提币网络
func (c *cWithdrawal) Options(ctx context.Context, req *withdrawal.OptionsReq) (res *withdrawal.OptionsRes, err error) {
	return service.Withdrawal().Options(ctx, req)
}

// Apply 提交提币申请
func (c *cWithdrawal) Apply(ctx context.Context, req *withdrawal.ApplyReq) (res *withdrawal.ApplyRes, err error) {
	userId := contexts.GetUserId(ctx)
	wallet, err := basesvc.UserWallet().GetWallet(ctx, userId)
	if err != nil {
		return nil, err
	}
	return service.Withdrawal().Apply(ctx, userId, wallet, req)
}

// List 用户提币记录
func (c *cWithdrawal) List(ctx context.Context, req *withdrawal.ListReq) (res *withdrawal.ListRes, err error) {
	return service.Withdrawal().List(ctx, contexts.GetUserId(ctx), req)
}

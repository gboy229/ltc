// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"

	apiBalance "hotgo/addons/chain/api/api/balance"
	"hotgo/addons/chain/service"
)

var (
	Balance = cBalance{}
)

type cBalance struct{}

// Get 查询钱包链上代币余额
func (c *cBalance) Get(ctx context.Context, req *apiBalance.GetReq) (res *apiBalance.GetRes, err error) {
	return service.Balance().Get(ctx, req)
}

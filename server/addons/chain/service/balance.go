// Package service
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package service

import (
	"context"

	apiBalance "hotgo/addons/chain/api/api/balance"
)

type (
	IBalance interface {
		// Get 查询钱包链上代币余额
		Get(ctx context.Context, in *apiBalance.GetReq) (res *apiBalance.GetRes, err error)
	}
)

var (
	localBalance IBalance
)

func Balance() IBalance {
	if localBalance == nil {
		panic("implement not found for interface IBalance, forgot register?")
	}
	return localBalance
}

func RegisterBalance(i IBalance) {
	localBalance = i
}

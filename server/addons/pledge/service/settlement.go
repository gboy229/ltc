// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import "context"

type (
	ISettlement interface {
		// Run 执行静态收益结算
		Run(ctx context.Context) (err error)
		// SettleOrder 结算指定订单，用于测试或手动补结算
		SettleOrder(ctx context.Context, orderId uint64) (err error)
		// CreatePushReward 用真实订单生成直推奖励，用于测试/补发
		CreatePushReward(ctx context.Context, orderId uint64, staticAmount float64, periodNo int) (err error)
	}
)

var (
	localSettlement ISettlement
)

func Settlement() ISettlement {
	if localSettlement == nil {
		panic("implement not found for interface ISettlement, forgot register?")
	}
	return localSettlement
}

func RegisterSettlement(i ISettlement) {
	localSettlement = i
}

// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"hotgo/addons/chain/api/api/withdrawal"
)

type (
	IWithdrawal interface {
		// Options 获取可用提币网络
		Options(ctx context.Context, in *withdrawal.OptionsReq) (res *withdrawal.OptionsRes, err error)
		// Apply 提交提币申请
		Apply(ctx context.Context, userId int64, wallet string, in *withdrawal.ApplyReq) (res *withdrawal.ApplyRes, err error)
		// List 用户提币记录
		List(ctx context.Context, userId int64, in *withdrawal.ListReq) (res *withdrawal.ListRes, err error)
		// Process 处理提现打币与结算
		Process(ctx context.Context) (err error)
	}
)

var (
	localWithdrawal IWithdrawal
)

func Withdrawal() IWithdrawal {
	if localWithdrawal == nil {
		panic("implement not found for interface IWithdrawal, forgot register?")
	}
	return localWithdrawal
}

func RegisterWithdrawal(i IWithdrawal) {
	localWithdrawal = i
}

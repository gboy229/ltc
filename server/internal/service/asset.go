// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"hotgo/internal/model/input/assetin"
)

type (
	IAsset interface {
		// EnsureAccount 确保用户资产账户存在
		EnsureAccount(ctx context.Context, userId uint64, account string) (err error)
		// Change 资产变更（增减/冻结/解冻/结算），同一业务同一动作幂等，仅返回成功或失败
		Change(ctx context.Context, in *assetin.ChangeInp) (err error)
		// Account 获取用户资产账户
		Account(ctx context.Context, in *assetin.AccountInp) (res *assetin.AccountModel, err error)
		// AccountList 获取资产账户列表
		AccountList(ctx context.Context, in *assetin.AccountListInp) (list []*assetin.AccountListModel, totalCount int, err error)
		// Logs 获取资产流水列表
		Logs(ctx context.Context, in *assetin.LogListInp) (list []*assetin.LogListModel, totalCount int, err error)
		// Stat 资产统计
		Stat(ctx context.Context, in *assetin.StatInp) (res *assetin.StatModel, err error)
	}
)

var (
	localAsset IAsset
)

func Asset() IAsset {
	if localAsset == nil {
		panic("implement not found for interface IAsset, forgot register?")
	}
	return localAsset
}

func RegisterAsset(i IAsset) {
	localAsset = i
}

// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"hotgo/addons/pledge/api/api/income"
	"hotgo/addons/pledge/model/input/incomein"
	"hotgo/addons/pledge/model/input/sysin"
)

type (
	ISysIncome interface {
		// List 收益流水列表
		List(ctx context.Context, in *sysin.IncomeListInp) (list []*sysin.IncomeListModel, totalCount int, err error)
	}
)

var (
	localSysIncome ISysIncome
)

func SysIncome() ISysIncome {
	if localSysIncome == nil {
		panic("implement not found for interface ISysIncome, forgot register?")
	}
	return localSysIncome
}

func RegisterSysIncome(i ISysIncome) {
	localSysIncome = i
}

type (
	IIncome interface {
		// Create 写入收益流水并资产入账
		Create(ctx context.Context, in *incomein.CreateInp) (id uint64, inserted bool, err error)
		// Sum 收益金额汇总，date为空时不限制日期
		Sum(ctx context.Context, userId uint64, types []string, date string) (amount float64, err error)
		// List 当前用户收益流水
		List(ctx context.Context, userId uint64, in *income.ListReq) (list []*income.ListItem, totalCount int, err error)
	}
)

var (
	localIncome IIncome
)

func Income() IIncome {
	if localIncome == nil {
		panic("implement not found for interface IIncome, forgot register?")
	}
	return localIncome
}

func RegisterIncome(i IIncome) {
	localIncome = i
}

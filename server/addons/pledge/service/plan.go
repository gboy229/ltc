// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"hotgo/addons/pledge/api/api/plan"
)

type (
	IPlan interface {
		// List 获取启用的投资套餐
		List(ctx context.Context, in *plan.ListReq) (res *plan.ListRes, err error)
	}
)

var (
	localPlan IPlan
)

func Plan() IPlan {
	if localPlan == nil {
		panic("implement not found for interface IPlan, forgot register?")
	}
	return localPlan
}

func RegisterPlan(i IPlan) {
	localPlan = i
}

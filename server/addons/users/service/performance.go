// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import "context"

type (
	IPerformance interface {
		// RefreshByUser 刷新该用户及所有上级的业绩
		RefreshByUser(ctx context.Context, userId int64) (err error)
	}
)

var (
	localPerformance IPerformance
)

func Performance() IPerformance {
	if localPerformance == nil {
		panic("implement not found for interface IPerformance, forgot register?")
	}
	return localPerformance
}

func RegisterPerformance(i IPerformance) {
	localPerformance = i
}

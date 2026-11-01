// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"hotgo/addons/pledge/model/input/sysin"
)

type (
	ISysPeriodConfig interface {
		// Delete 删除周期配置
		Delete(ctx context.Context, in *sysin.PeriodConfigDeleteInp) (err error)
		// Edit 修改/新增周期配置
		Edit(ctx context.Context, in *sysin.PeriodConfigEditInp) (err error)
		// Status 更新周期配置状态
		Status(ctx context.Context, in *sysin.PeriodConfigStatusInp) (err error)
		// View 获取周期配置详情
		View(ctx context.Context, in *sysin.PeriodConfigViewInp) (res *sysin.PeriodConfigViewModel, err error)
		// List 获取周期配置列表
		List(ctx context.Context, in *sysin.PeriodConfigListInp) (list []*sysin.PeriodConfigListModel, totalCount int, err error)
	}
)

var (
	localSysPeriodConfig ISysPeriodConfig
)

func SysPeriodConfig() ISysPeriodConfig {
	if localSysPeriodConfig == nil {
		panic("implement not found for interface ISysPeriodConfig, forgot register?")
	}
	return localSysPeriodConfig
}

func RegisterSysPeriodConfig(i ISysPeriodConfig) {
	localSysPeriodConfig = i
}

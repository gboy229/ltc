// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	apiConfig "hotgo/addons/chain/api/api/config"
)

type (
	IConfig interface {
		// List 获取前台支付链配置
		List(ctx context.Context, in *apiConfig.ListReq) (res *apiConfig.ListRes, err error)
	}
)

var (
	localConfig IConfig
)

func Config() IConfig {
	if localConfig == nil {
		panic("implement not found for interface IConfig, forgot register?")
	}
	return localConfig
}

func RegisterConfig(i IConfig) {
	localConfig = i
}

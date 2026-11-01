// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"

	apiConfig "hotgo/addons/chain/api/api/config"
	"hotgo/addons/chain/service"
)

var (
	Config = cConfig{}
)

type cConfig struct{}

// List 获取前台支付链配置
func (c *cConfig) List(ctx context.Context, req *apiConfig.ListReq) (res *apiConfig.ListRes, err error) {
	return service.Config().List(ctx, req)
}

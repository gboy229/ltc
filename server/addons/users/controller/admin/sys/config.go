// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/addons/users/api/admin/config"
	"hotgo/addons/users/service"
)

var (
	Config = cConfig{}
)

type cConfig struct{}

// GetConfig 获取指定分组的配置
func (c *cConfig) GetConfig(ctx context.Context, req *config.GetReq) (res *config.GetRes, err error) {
	data, err := service.SysConfig().GetConfigByGroup(ctx, &req.GetConfigInp)

	res = new(config.GetRes)
	res.GetConfigModel = data
	return
}

// UpdateConfig 更新指定分组的配置
func (c *cConfig) UpdateConfig(ctx context.Context, req *config.UpdateReq) (res *config.UpdateRes, err error) {
	err = service.SysConfig().UpdateConfigByGroup(ctx, &req.UpdateConfigInp)
	return
}

// Items 获取插件全部配置
func (c *cConfig) Items(ctx context.Context, req *config.ItemsReq) (res *config.ItemsRes, err error) {
	data, err := service.SysConfig().GetConfigItems(ctx)
	if err != nil {
		return
	}
	res = new(config.ItemsRes)
	res.GetConfigItemsModel = data
	return
}

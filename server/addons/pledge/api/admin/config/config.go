// Package config
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/pledge/model/input/sysin"
)

// GetReq 获取指定分组的配置
type GetReq struct {
	g.Meta `path:"/config/get" method:"get" tags:"质押管理" summary:"获取指定分组的配置"`
	sysin.GetConfigInp
}

type GetRes struct {
	*sysin.GetConfigModel
}

// UpdateReq 获取指定分组的配置
type UpdateReq struct {
	g.Meta `path:"/config/update" method:"post" tags:"质押管理" summary:"获取指定分组的配置"`
	sysin.UpdateConfigInp
}

type UpdateRes struct {
}

// ItemsReq 获取插件全部配置
type ItemsReq struct {
	g.Meta `path:"/config/items" method:"get" tags:"质押管理" summary:"获取插件全部配置"`
}

type ItemsRes struct {
	*sysin.GetConfigItemsModel
}

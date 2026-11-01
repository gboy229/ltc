// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/internal/model/input/sysin"
)

// UpdateConfigInp 更新指定配置
type UpdateConfigInp struct {
	sysin.UpdateAddonsConfigInp
}

type GetConfigInp struct {
	sysin.GetAddonsConfigInp
}

type GetConfigModel struct {
	List g.Map `json:"list"`
}

// ConfigItemModel 配置项（含元数据）
type ConfigItemModel struct {
	Group string      `json:"group"`
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Tip   string      `json:"tip"`
	Sort  int         `json:"sort"`
}

// ConfigGroupModel 配置分组
type ConfigGroupModel struct {
	Group string `json:"group"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Sort  int    `json:"sort"`
}

// GetConfigItemsModel 插件全部配置（分组+配置项）
type GetConfigItemsModel struct {
	Groups []*ConfigGroupModel `json:"groups"`
	Items  []*ConfigItemModel  `json:"items"`
}

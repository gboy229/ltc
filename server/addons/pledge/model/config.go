// Package model
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package model

// BasicConfig 基础配置
type BasicConfig struct {
	DepositAddress string `json:"depositAddress"` // 入金地址
	OrderEnabled   bool   `json:"orderEnabled"`   // 是否开始下单
}

// ConfigGroup 插件配置分组定义
type ConfigGroup struct {
	Group string // 分组标识
	Name  string // 分组名称（tab标题）
	Desc  string // 分组描述
	Sort  int    // 排序
}

// ConfigGroups 插件配置分组（供前端动态渲染 tab 使用）
var ConfigGroups = []ConfigGroup{
	{Group: "basic", Name: "基本设置", Desc: "入金地址与下单开关", Sort: 1},
}

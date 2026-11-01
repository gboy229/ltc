// Package model
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package model

// BasicConfig 基础配置
type BasicConfig struct {
	Test string `json:"basicTest"`
}

// LoginConfig 登录/注册配置
type LoginConfig struct {
	RegisterSwitch bool `json:"registerSwitch"` // 注册开关
	ForceInvite    bool `json:"forceInvite"`    // 强制邀请
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
	{Group: "login", Name: "用户配置", Desc: "登录/注册设置", Sort: 1},
}

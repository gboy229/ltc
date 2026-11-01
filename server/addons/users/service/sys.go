// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"hotgo/addons/users/model"
	"hotgo/addons/users/model/input/sysin"
)

type (
	ISysConfig interface {
		// GetBasic 获取基础配置
		GetBasic(ctx context.Context) (conf *model.BasicConfig, err error)
		// GetConfigByGroup 获取指定分组配置
		GetConfigByGroup(ctx context.Context, in *sysin.GetConfigInp) (res *sysin.GetConfigModel, err error)
		// UpdateConfigByGroup 更新指定分组的配置
		UpdateConfigByGroup(ctx context.Context, in *sysin.UpdateConfigInp) error
		// GetConfigItems 获取插件全部配置（分组+配置项）
		GetConfigItems(ctx context.Context) (res *sysin.GetConfigItemsModel, err error)
	}
	ISysLevel interface {
		// Delete 删除用户等级
		Delete(ctx context.Context, in *sysin.LevelDeleteInp) (err error)
		// Edit 修改/新增用户等级
		Edit(ctx context.Context, in *sysin.LevelEditInp) (err error)
		// Status 更新用户等级状态
		Status(ctx context.Context, in *sysin.LevelStatusInp) (err error)
		// MaxSort 获取用户等级最大排序
		MaxSort(ctx context.Context, in *sysin.LevelMaxSortInp) (res *sysin.LevelMaxSortModel, err error)
		// View 获取指定用户等级信息
		View(ctx context.Context, in *sysin.LevelViewInp) (res *sysin.LevelViewModel, err error)
		// List 获取用户等级列表
		List(ctx context.Context, in *sysin.LevelListInp) (list []*sysin.LevelListModel, totalCount int, err error)
		// Option 获取用户等级选项
		Option(ctx context.Context, in *sysin.LevelOptionInp) (list []*sysin.LevelOptionModel, err error)
	}
	ISysUser interface {
		// Delete 删除用户
		Delete(ctx context.Context, in *sysin.UserDeleteInp) (err error)
		// Edit 修改/新增用户
		Edit(ctx context.Context, in *sysin.UserEditInp) (res *sysin.UserEditModel, err error)
		// Status 更新用户状态
		Status(ctx context.Context, in *sysin.UserStatusInp) (err error)
		// View 获取指定用户信息
		View(ctx context.Context, in *sysin.UserViewInp) (res *sysin.UserViewModel, err error)
		// List 获取用户列表
		List(ctx context.Context, in *sysin.UserListInp) (list []*sysin.UserListModel, totalCount int, err error)
	}
	ISysAsset interface {
		// AccountList 获取资产账户列表
		AccountList(ctx context.Context, in *sysin.AssetAccountListInp) (list []*sysin.AssetAccountListModel, totalCount int, err error)
		// LogList 获取资金流水列表
		LogList(ctx context.Context, in *sysin.AssetLogListInp) (list []*sysin.AssetLogListModel, totalCount int, err error)
		// Stat 资产统计
		Stat(ctx context.Context, in *sysin.AssetStatInp) (res *sysin.AssetStatModel, err error)
	}
)

var (
	localSysConfig ISysConfig
	localSysLevel  ISysLevel
	localSysUser   ISysUser
	localSysAsset  ISysAsset
)

func SysConfig() ISysConfig {
	if localSysConfig == nil {
		panic("implement not found for interface ISysConfig, forgot register?")
	}
	return localSysConfig
}

func RegisterSysConfig(i ISysConfig) {
	localSysConfig = i
}

func SysLevel() ISysLevel {
	if localSysLevel == nil {
		panic("implement not found for interface ISysLevel, forgot register?")
	}
	return localSysLevel
}

func RegisterSysLevel(i ISysLevel) {
	localSysLevel = i
}

func SysUser() ISysUser {
	if localSysUser == nil {
		panic("implement not found for interface ISysUser, forgot register?")
	}
	return localSysUser
}

func RegisterSysUser(i ISysUser) {
	localSysUser = i
}

func SysAsset() ISysAsset {
	if localSysAsset == nil {
		panic("implement not found for interface ISysAsset, forgot register?")
	}
	return localSysAsset
}

func RegisterSysAsset(i ISysAsset) {
	localSysAsset = i
}

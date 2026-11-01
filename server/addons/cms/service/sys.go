// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"hotgo/addons/cms/model"
	"hotgo/addons/cms/model/input/sysin"
)

type (
	ISysArticle interface {
		// Delete 删除公告
		Delete(ctx context.Context, in *sysin.ArticleDeleteInp) (err error)
		// Edit 修改/新增公告
		Edit(ctx context.Context, in *sysin.ArticleEditInp) (err error)
		// Status 更新公告状态
		Status(ctx context.Context, in *sysin.ArticleStatusInp) (err error)
		// MaxSort 获取公告最大排序
		MaxSort(ctx context.Context, in *sysin.ArticleMaxSortInp) (res *sysin.ArticleMaxSortModel, err error)
		// View 获取指定公告
		View(ctx context.Context, in *sysin.ArticleViewInp) (res *sysin.ArticleViewModel, err error)
		// List 获取公告列表
		List(ctx context.Context, in *sysin.ArticleListInp) (list []*sysin.ArticleListModel, totalCount int, err error)
	}
	ISysConfig interface {
		// GetBasic 获取基础配置
		GetBasic(ctx context.Context) (conf *model.BasicConfig, err error)
		// GetConfigByGroup 获取指定分组配置
		GetConfigByGroup(ctx context.Context, in *sysin.GetConfigInp) (res *sysin.GetConfigModel, err error)
		// UpdateConfigByGroup 更新指定分组的配置
		UpdateConfigByGroup(ctx context.Context, in *sysin.UpdateConfigInp) error
	}
)

var (
	localSysArticle ISysArticle
	localSysConfig  ISysConfig
)

func SysArticle() ISysArticle {
	if localSysArticle == nil {
		panic("implement not found for interface ISysArticle, forgot register?")
	}
	return localSysArticle
}

func RegisterSysArticle(i ISysArticle) {
	localSysArticle = i
}

func SysConfig() ISysConfig {
	if localSysConfig == nil {
		panic("implement not found for interface ISysConfig, forgot register?")
	}
	return localSysConfig
}

func RegisterSysConfig(i ISysConfig) {
	localSysConfig = i
}

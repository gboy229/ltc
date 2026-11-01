// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/cms/api/admin/article"
	"hotgo/addons/cms/model/input/sysin"
	"hotgo/addons/cms/service"
)

var (
	Article = cArticle{}
)

type cArticle struct{}

// Delete 删除公告
func (c *cArticle) Delete(ctx context.Context, req *article.DeleteReq) (res *article.DeleteRes, err error) {
	err = service.SysArticle().Delete(ctx, &req.ArticleDeleteInp)
	return
}

// Edit 修改/新增公告
func (c *cArticle) Edit(ctx context.Context, req *article.EditReq) (res *article.EditRes, err error) {
	err = service.SysArticle().Edit(ctx, &req.ArticleEditInp)
	return
}

// MaxSort 获取最大排序
func (c *cArticle) MaxSort(ctx context.Context, req *article.MaxSortReq) (res *article.MaxSortRes, err error) {
	res = new(article.MaxSortRes)
	res.ArticleMaxSortModel, err = service.SysArticle().MaxSort(ctx, &req.ArticleMaxSortInp)
	return
}

// View 获取指定公告
func (c *cArticle) View(ctx context.Context, req *article.ViewReq) (res *article.ViewRes, err error) {
	data, err := service.SysArticle().View(ctx, &req.ArticleViewInp)
	if err != nil {
		return
	}
	res = new(article.ViewRes)
	res.ArticleViewModel = data
	return
}

// List 获取公告列表
func (c *cArticle) List(ctx context.Context, req *article.ListReq) (res *article.ListRes, err error) {
	list, totalCount, err := service.SysArticle().List(ctx, &req.ArticleListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.ArticleListModel{}
	}
	res = new(article.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Status 更新公告状态
func (c *cArticle) Status(ctx context.Context, req *article.StatusReq) (res *article.StatusRes, err error) {
	err = service.SysArticle().Status(ctx, &req.ArticleStatusInp)
	return
}

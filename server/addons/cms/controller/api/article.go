// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"

	"hotgo/addons/cms/api/api/article"
	"hotgo/addons/cms/service"
)

var (
	Article = cArticle{}
)

type cArticle struct{}

// List 获取已发布公告列表
func (c *cArticle) List(ctx context.Context, req *article.ListReq) (res *article.ListRes, err error) {
	return service.Article().List(ctx, req)
}

// Recommend 获取首页推荐公告
func (c *cArticle) Recommend(ctx context.Context, req *article.RecommendReq) (res *article.RecommendRes, err error) {
	return service.Article().Recommend(ctx, req)
}

// View 获取已发布公告详情
func (c *cArticle) View(ctx context.Context, req *article.ViewReq) (res *article.ViewRes, err error) {
	return service.Article().View(ctx, req)
}

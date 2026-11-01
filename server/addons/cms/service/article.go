// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"hotgo/addons/cms/api/api/article"
)

type (
	IArticle interface {
		// List 获取已发布公告列表
		List(ctx context.Context, in *article.ListReq) (res *article.ListRes, err error)
		// Recommend 获取首页推荐公告
		Recommend(ctx context.Context, in *article.RecommendReq) (res *article.RecommendRes, err error)
		// View 获取已发布公告详情
		View(ctx context.Context, in *article.ViewReq) (res *article.ViewRes, err error)
	}
)

var (
	localArticle IArticle
)

func Article() IArticle {
	if localArticle == nil {
		panic("implement not found for interface IArticle, forgot register?")
	}
	return localArticle
}

func RegisterArticle(i IArticle) {
	localArticle = i
}

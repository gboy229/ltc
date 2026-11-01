// Package article
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package article

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/model/input/form"
)

// Item 前台公告列表项（仅返回展示所需字段）
type Item struct {
	Id        string `json:"id"        dc:"公告ID"`
	Title     string `json:"title"     dc:"标题"`
	CreatedAt string `json:"createdAt" dc:"发布日期，格式：Y-m-d"`
}

// ListReq 获取已发布公告列表（前台，无需登录）
type ListReq struct {
	g.Meta `path:"/article/list" method:"get" tags:"公告" summary:"获取公告列表"`
	form.PageReq
	Title string `json:"title" dc:"标题关键词"`
}

type ListRes struct {
	List []*Item `json:"list" dc:"数据列表"`
	form.PageRes
}

// RecommendReq 首页推荐公告（前台，无需登录，不分页）
type RecommendReq struct {
	g.Meta `path:"/article/recommend" method:"get" tags:"公告" summary:"获取首页推荐公告"`
}

type RecommendRes struct {
	List []*Item `json:"list" dc:"推荐公告列表"`
}

// ViewReq 获取已发布公告详情（前台，无需登录）
type ViewReq struct {
	g.Meta    `path:"/article/view/{article_no}" method:"get" tags:"公告" summary:"获取公告详情"`
	ArticleNo string `json:"article_no" in:"path" v:"required#公告ID不能为空" dc:"公告ID"`
}

type ViewRes struct {
	Id        string `json:"id"        dc:"公告ID"`
	Title     string `json:"title"     dc:"标题"`
	Content   string `json:"content"   dc:"内容"`
	CreatedAt string `json:"createdAt" dc:"发布日期，格式：Y-m-d"`
}

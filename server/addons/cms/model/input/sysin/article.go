// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"hotgo/addons/cms/consts"
	"hotgo/addons/cms/model/entity"
	"hotgo/internal/model/input/form"
)

// ArticleEditInp 修改/新增公告
type ArticleEditInp struct {
	entity.CmsArticle
}

func (in *ArticleEditInp) Filter(ctx context.Context) (err error) {
	if in.Title == "" {
		err = gerror.New("公告标题不能为空")
		return
	}
	if in.Content == "" {
		err = gerror.New("公告内容不能为空")
		return
	}
	if in.Type <= 0 {
		in.Type = consts.CmsArticleTypeNotice
	}
	if in.IsRecommend != consts.CmsArticleRecommendNo && in.IsRecommend != consts.CmsArticleRecommendYes {
		in.IsRecommend = consts.CmsArticleRecommendYes
	}
	return
}

type ArticleEditModel struct{}

// ArticleUpdateFields 修改公告字段过滤
type ArticleUpdateFields struct {
	Title       string      `json:"title"       description:"标题"`
	Cover       string      `json:"cover"       description:"封面图"`
	Content     string      `json:"content"     description:"内容"`
	Sort        int         `json:"sort"        description:"排序"`
	Status      int         `json:"status"      description:"状态"`
	IsRecommend int         `json:"isRecommend" description:"是否推荐"`
	UpdatedBy   int64       `json:"updatedBy"   description:"修改人"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"更新时间"`
}

// ArticleInsertFields 新增公告字段过滤
type ArticleInsertFields struct {
	ArticleNo   string      `json:"articleNo"   description:"前台公开ID"`
	Title       string      `json:"title"       description:"标题"`
	Type        int         `json:"type"        description:"类型"`
	Cover       string      `json:"cover"       description:"封面图"`
	Content     string      `json:"content"     description:"内容"`
	Sort        int         `json:"sort"        description:"排序"`
	Status      int         `json:"status"      description:"状态"`
	IsRecommend int         `json:"isRecommend" description:"是否推荐"`
	CreatedBy   int64       `json:"createdBy"   description:"创建人"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
}

// ArticleDeleteInp 删除公告
type ArticleDeleteInp struct {
	Id interface{} `json:"id" v:"required#公告ID不能为空" dc:"公告ID"`
}

type ArticleDeleteModel struct{}

// ArticleViewInp 获取指定公告
type ArticleViewInp struct {
	Id int64 `json:"id" v:"required#公告ID不能为空" dc:"公告ID"`
}

type ArticleViewModel struct {
	entity.CmsArticle
}

// ArticleListInp 获取公告列表
type ArticleListInp struct {
	form.PageReq
	form.StatusReq
	Title       string `json:"title"       dc:"标题"`
	IsRecommend *int   `json:"isRecommend" dc:"是否推荐"`
}

type ArticleListModel struct {
	entity.CmsArticle
}

// ArticleStatusInp 更新公告状态
type ArticleStatusInp struct {
	Id     int64 `json:"id"     v:"required#公告ID不能为空" dc:"公告ID"`
	Status int   `json:"status" dc:"状态"`
}

func (in *ArticleStatusInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		err = gerror.New("公告ID不能为空")
		return
	}
	if in.Status <= 0 {
		err = gerror.New("状态不能为空")
		return
	}
	return
}

type ArticleStatusModel struct{}

// ArticleMaxSortInp 获取公告最大排序
type ArticleMaxSortInp struct{}

type ArticleMaxSortModel struct {
	Sort int `json:"sort" dc:"排序"`
}

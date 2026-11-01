// Package article
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package article

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	apiArticle "hotgo/addons/cms/api/api/article"
	cmsconsts "hotgo/addons/cms/consts"
	"hotgo/addons/cms/model/dao"
	"hotgo/addons/cms/service"
	"hotgo/internal/consts"
)

type sArticle struct{}

func NewArticle() *sArticle {
	return &sArticle{}
}

func init() {
	service.RegisterArticle(NewArticle())
}

// articleRow 前台公告查询行
type articleRow struct {
	Id        string
	Title     string
	Content   string
	CreatedAt *gtime.Time
}

func toItem(row *articleRow) *apiArticle.Item {
	item := &apiArticle.Item{
		Id:    row.Id,
		Title: row.Title,
	}
	if row.CreatedAt != nil {
		item.CreatedAt = row.CreatedAt.Format("Y-m-d")
	}
	return item
}

// List 获取已发布公告列表
func (s *sArticle) List(ctx context.Context, in *apiArticle.ListReq) (res *apiArticle.ListRes, err error) {
	cols := dao.CmsArticle.Columns()
	mod := dao.CmsArticle.Ctx(ctx).
		Where(cols.Type, cmsconsts.CmsArticleTypeNotice).
		Where(cols.Status, consts.StatusEnabled)

	if in.Title != "" {
		mod = mod.WhereLike(cols.Title, "%"+in.Title+"%")
	}

	totalCount, err := mod.Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询公告总数失败")
	}

	res = new(apiArticle.ListRes)
	res.List = []*apiArticle.Item{}
	res.PageRes.Pack(in, totalCount)
	if totalCount == 0 {
		return
	}

	var rows []*articleRow
	if err = mod.Fields(cols.ArticleNo+" AS id", cols.Title, cols.CreatedAt).
		Page(in.Page, in.PerPage).
		OrderDesc(cols.Sort).
		OrderDesc(cols.Id).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "查询公告列表失败")
	}

	for _, row := range rows {
		res.List = append(res.List, toItem(row))
	}
	return
}

// Recommend 获取首页全部推荐公告（不分页）
func (s *sArticle) Recommend(ctx context.Context, _ *apiArticle.RecommendReq) (res *apiArticle.RecommendRes, err error) {
	cols := dao.CmsArticle.Columns()

	var rows []*articleRow
	if err = dao.CmsArticle.Ctx(ctx).
		Where(cols.Type, cmsconsts.CmsArticleTypeNotice).
		Where(cols.Status, consts.StatusEnabled).
		Where(cols.IsRecommend, cmsconsts.CmsArticleRecommendYes).
		Fields(cols.ArticleNo+" AS id", cols.Title, cols.CreatedAt).
		OrderDesc(cols.Sort).
		OrderDesc(cols.Id).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "查询推荐公告失败")
	}

	res = new(apiArticle.RecommendRes)
	res.List = make([]*apiArticle.Item, 0, len(rows))
	for _, row := range rows {
		res.List = append(res.List, toItem(row))
	}
	return
}

// View 获取已发布公告详情
func (s *sArticle) View(ctx context.Context, in *apiArticle.ViewReq) (res *apiArticle.ViewRes, err error) {
	articleNo := strings.TrimSpace(in.ArticleNo)
	if articleNo == "" {
		return nil, gerror.New("公告ID不能为空")
	}

	cols := dao.CmsArticle.Columns()
	var row *articleRow
	if err = dao.CmsArticle.Ctx(ctx).
		Where(cols.ArticleNo, articleNo).
		Where(cols.Type, cmsconsts.CmsArticleTypeNotice).
		Where(cols.Status, consts.StatusEnabled).
		Fields(cols.ArticleNo+" AS id", cols.Title, cols.Content, cols.CreatedAt).
		Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "查询公告详情失败")
	}
	if row == nil || row.Id == "" {
		return nil, gerror.New("公告不存在或已下架")
	}

	res = &apiArticle.ViewRes{
		Id:      row.Id,
		Title:   row.Title,
		Content: row.Content,
	}
	if row.CreatedAt != nil {
		res.CreatedAt = row.CreatedAt.Format("Y-m-d")
	}
	return
}

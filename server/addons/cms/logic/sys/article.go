// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"

	cmsconsts "hotgo/addons/cms/consts"
	"hotgo/addons/cms/model/dao"
	"hotgo/addons/cms/model/input/sysin"
	"hotgo/addons/cms/service"
	"hotgo/internal/library/contexts"
	"hotgo/internal/model/input/form"
	"hotgo/utility/validate"
)

type sSysArticle struct{}

func NewSysArticle() *sSysArticle {
	return &sSysArticle{}
}

func init() {
	service.RegisterSysArticle(NewSysArticle())
}

// Delete 删除公告
func (s *sSysArticle) Delete(ctx context.Context, in *sysin.ArticleDeleteInp) (err error) {
	_, err = dao.CmsArticle.Ctx(ctx).Where(dao.CmsArticle.Columns().Id, in.Id).Delete()
	return
}

// Edit 修改/新增公告
func (s *sSysArticle) Edit(ctx context.Context, in *sysin.ArticleEditInp) (err error) {
	if err = validate.PreFilter(ctx, in); err != nil {
		return
	}

	var memberId int64
	if c := contexts.Get(ctx); c != nil && c.User != nil {
		memberId = c.User.Id
	}

	// 修改
	if in.Id > 0 {
		if memberId > 0 {
			in.UpdatedBy = memberId
		}
		in.UpdatedAt = gtime.Now()
		_, err = dao.CmsArticle.Ctx(ctx).Fields(sysin.ArticleUpdateFields{}).WherePri(in.Id).Data(in).Update()
		return
	}

	// 新增
	if memberId > 0 {
		in.CreatedBy = memberId
	}
	if in.ArticleNo, err = s.genArticleNo(ctx); err != nil {
		return
	}
	in.CreatedAt = gtime.Now()
	_, err = dao.CmsArticle.Ctx(ctx).Fields(sysin.ArticleInsertFields{}).Data(in).Insert()
	return
}

// genArticleNo 生成前台公开ID，避免暴露数据库主键
func (s *sSysArticle) genArticleNo(ctx context.Context) (string, error) {
	cols := dao.CmsArticle.Columns()
	for i := 0; i < 5; i++ {
		articleNo := grand.Str("0123456789abcdefghijklmnopqrstuvwxyz", 20)
		exist, err := dao.CmsArticle.Ctx(ctx).Where(cols.ArticleNo, articleNo).Exist()
		if err != nil {
			return "", gerror.Wrap(err, "生成公告ID失败")
		}
		if !exist {
			return articleNo, nil
		}
	}
	return grand.Str("0123456789abcdefghijklmnopqrstuvwxyz", 24), nil
}

// Status 更新公告状态
func (s *sSysArticle) Status(ctx context.Context, in *sysin.ArticleStatusInp) (err error) {
	if err = validate.PreFilter(ctx, in); err != nil {
		return
	}
	_, err = dao.CmsArticle.Ctx(ctx).WherePri(in.Id).Data(dao.CmsArticle.Columns().Status, in.Status).Update()
	return
}

// MaxSort 获取公告最大排序
func (s *sSysArticle) MaxSort(ctx context.Context, in *sysin.ArticleMaxSortInp) (res *sysin.ArticleMaxSortModel, err error) {
	if err = dao.CmsArticle.Ctx(ctx).
		Where(dao.CmsArticle.Columns().Type, cmsconsts.CmsArticleTypeNotice).
		Fields(dao.CmsArticle.Columns().Sort).
		OrderDesc(dao.CmsArticle.Columns().Sort).
		Scan(&res); err != nil {
		return
	}
	if res == nil {
		res = new(sysin.ArticleMaxSortModel)
	}
	res.Sort = form.DefaultMaxSort(res.Sort)
	return
}

// View 获取指定公告
func (s *sSysArticle) View(ctx context.Context, in *sysin.ArticleViewInp) (res *sysin.ArticleViewModel, err error) {
	if err = dao.CmsArticle.Ctx(ctx).WherePri(in.Id).Scan(&res); err != nil {
		return
	}
	return
}

// List 获取公告列表
func (s *sSysArticle) List(ctx context.Context, in *sysin.ArticleListInp) (list []*sysin.ArticleListModel, totalCount int, err error) {
	cols := dao.CmsArticle.Columns()
	mod := dao.CmsArticle.Ctx(ctx).Where(cols.Type, cmsconsts.CmsArticleTypeNotice)

	if in.Title != "" {
		mod = mod.WhereLike(cols.Title, "%"+in.Title+"%")
	}

	if in.Status > 0 {
		mod = mod.Where(cols.Status, in.Status)
	}

	if in.IsRecommend != nil {
		mod = mod.Where(cols.IsRecommend, *in.IsRecommend)
	}

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	if totalCount == 0 {
		return
	}

	if err = mod.Page(in.Page, in.PerPage).OrderDesc(cols.Sort).OrderDesc(cols.Id).Scan(&list); err != nil {
		return
	}
	return
}

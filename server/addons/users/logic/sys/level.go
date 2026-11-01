// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/addons/users/model/dao"
	"hotgo/addons/users/model/entity"
	"hotgo/addons/users/model/input/sysin"
	"hotgo/addons/users/service"
	"hotgo/internal/consts"
	"hotgo/internal/model/input/form"
)

type sSysLevel struct{}

func NewSysLevel() *sSysLevel {
	return &sSysLevel{}
}

func init() {
	service.RegisterSysLevel(NewSysLevel())
}

// Delete 删除用户等级
func (s *sSysLevel) Delete(ctx context.Context, in *sysin.LevelDeleteInp) (err error) {
	var models *entity.UsersLevel
	if err = dao.UsersLevel.Ctx(ctx).WherePri(in.Id).Scan(&models); err != nil {
		return
	}
	if models == nil {
		err = gerror.New("等级不存在或已删除")
		return
	}
	_, err = dao.UsersLevel.Ctx(ctx).WherePri(in.Id).Delete()
	return
}

// Edit 修改/新增用户等级
func (s *sSysLevel) Edit(ctx context.Context, in *sysin.LevelEditInp) (err error) {
	// 修改
	if in.Id > 0 {
		_, err = dao.UsersLevel.Ctx(ctx).Fields(sysin.LevelUpdateFields{}).WherePri(in.Id).Data(in).Update()
		return
	}

	// 新增
	_, err = dao.UsersLevel.Ctx(ctx).Fields(sysin.LevelInsertFields{}).Data(in).OmitEmptyData().Insert()
	return
}

// Status 更新用户等级状态
func (s *sSysLevel) Status(ctx context.Context, in *sysin.LevelStatusInp) (err error) {
	_, err = dao.UsersLevel.Ctx(ctx).WherePri(in.Id).Data(dao.UsersLevel.Columns().Status, in.Status).Update()
	return
}

// MaxSort 获取用户等级最大排序
func (s *sSysLevel) MaxSort(ctx context.Context, in *sysin.LevelMaxSortInp) (res *sysin.LevelMaxSortModel, err error) {
	if err = dao.UsersLevel.Ctx(ctx).Fields(dao.UsersLevel.Columns().Sort).OrderDesc(dao.UsersLevel.Columns().Sort).Scan(&res); err != nil {
		return
	}
	if res == nil {
		res = new(sysin.LevelMaxSortModel)
	}
	res.Sort = form.DefaultMaxSort(res.Sort)
	return
}

// View 获取指定用户等级信息
func (s *sSysLevel) View(ctx context.Context, in *sysin.LevelViewInp) (res *sysin.LevelViewModel, err error) {
	if err = dao.UsersLevel.Ctx(ctx).WherePri(in.Id).Scan(&res); err != nil {
		return
	}
	return
}

// List 获取用户等级列表
func (s *sSysLevel) List(ctx context.Context, in *sysin.LevelListInp) (list []*sysin.LevelListModel, totalCount int, err error) {
	mod := dao.UsersLevel.Ctx(ctx)
	cols := dao.UsersLevel.Columns()

	if in.Name != "" {
		mod = mod.WhereLike(cols.Name, "%"+in.Name+"%")
	}

	if in.Status > 0 {
		mod = mod.Where(cols.Status, in.Status)
	}

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	if totalCount == 0 {
		return
	}

	if err = mod.Page(in.Page, in.PerPage).Order("id asc").Scan(&list); err != nil {
		return
	}
	return
}

// Option 获取用户等级选项
func (s *sSysLevel) Option(ctx context.Context, in *sysin.LevelOptionInp) (list []*sysin.LevelOptionModel, err error) {
	cols := dao.UsersLevel.Columns()
	err = dao.UsersLevel.Ctx(ctx).
		Fields(cols.Id, cols.Name).
		Where(cols.Status, consts.StatusEnabled).
		Order("sort desc, id asc").
		Scan(&list)
	return
}

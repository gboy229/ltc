// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"hotgo/addons/pledge/model/dao"
	"hotgo/addons/pledge/model/entity"
	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/addons/pledge/service"
)

type sSysPlan struct{}

func NewSysPlan() *sSysPlan {
	return &sSysPlan{}
}

func init() {
	service.RegisterSysPlan(NewSysPlan())
}

// Delete 删除投资套餐
func (s *sSysPlan) Delete(ctx context.Context, in *sysin.PlanDeleteInp) (err error) {
	cols := dao.PledgePlan.Columns()
	if len(in.Ids) == 0 {
		return gerror.New("请选择要删除的套餐")
	}
	if _, err = dao.PledgePlan.Ctx(ctx).WhereIn(cols.Id, in.Ids).Delete(); err != nil {
		return gerror.Wrap(err, "删除投资套餐失败")
	}
	return
}

// Edit 修改/新增投资套餐
func (s *sSysPlan) Edit(ctx context.Context, in *sysin.PlanEditInp) (err error) {
	cols := dao.PledgePlan.Columns()
	mod := dao.PledgePlan.Ctx(ctx).
		Where(cols.Name, in.Name).
		Where(cols.PeriodMonths, in.PeriodMonths)
	if in.Id > 0 {
		mod = mod.WhereNot(cols.Id, in.Id)
	}
	exist, err := mod.Exist()
	if err != nil {
		return gerror.Wrap(err, "查询投资套餐失败")
	}
	if exist {
		return gerror.New("该套餐周期已存在")
	}

	if in.Id > 0 {
		_, err = dao.PledgePlan.Ctx(ctx).
			Fields(sysin.PlanUpdateFields{}).
			WherePri(in.Id).
			Data(in).
			Update()
		if err != nil {
			return gerror.Wrap(err, "更新投资套餐失败")
		}
		return
	}

	_, err = dao.PledgePlan.Ctx(ctx).
		Fields(sysin.PlanInsertFields{}).
		Data(in).
		OmitEmptyData().
		Insert()
	if err != nil {
		return gerror.Wrap(err, "新增投资套餐失败")
	}
	return
}

// Status 更新投资套餐状态
func (s *sSysPlan) Status(ctx context.Context, in *sysin.PlanStatusInp) (err error) {
	_, err = dao.PledgePlan.Ctx(ctx).
		WherePri(in.Id).
		Data(dao.PledgePlan.Columns().Status, in.Status).
		Update()
	if err != nil {
		return gerror.Wrap(err, "更新投资套餐状态失败")
	}
	return
}

// View 获取指定投资套餐信息
func (s *sSysPlan) View(ctx context.Context, in *sysin.PlanViewInp) (res *sysin.PlanViewModel, err error) {
	var plan *entity.PledgePlan
	if err = dao.PledgePlan.Ctx(ctx).WherePri(in.Id).Scan(&plan); err != nil {
		return
	}
	if plan == nil {
		return
	}
	res = &sysin.PlanViewModel{PledgePlan: *plan}
	return
}

// List 获取投资套餐列表
func (s *sSysPlan) List(ctx context.Context, in *sysin.PlanListInp) (list []*sysin.PlanListModel, totalCount int, err error) {
	cols := dao.PledgePlan.Columns()
	mod := dao.PledgePlan.Ctx(ctx)
	if in.Name != "" {
		mod = mod.WhereLike(cols.Name, "%"+in.Name+"%")
	}
	if in.PeriodMonths > 0 {
		mod = mod.Where(cols.PeriodMonths, in.PeriodMonths)
	}
	if in.Status > 0 {
		mod = mod.Where(cols.Status, in.Status)
	}

	totalCount, err = mod.Clone().Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, "查询投资套餐总数失败")
	}
	if totalCount == 0 {
		return
	}
	if err = mod.Page(in.Page, in.PerPage).OrderAsc(cols.Name).OrderAsc(cols.PeriodMonths).Scan(&list); err != nil {
		return nil, 0, gerror.Wrap(err, "查询投资套餐失败")
	}
	return
}

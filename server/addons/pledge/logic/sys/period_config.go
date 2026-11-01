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

type sSysPeriodConfig struct{}

func NewSysPeriodConfig() *sSysPeriodConfig {
	return &sSysPeriodConfig{}
}

func init() {
	service.RegisterSysPeriodConfig(NewSysPeriodConfig())
}

// Delete 删除周期配置
func (s *sSysPeriodConfig) Delete(ctx context.Context, in *sysin.PeriodConfigDeleteInp) (err error) {
	cols := dao.PledgePeriodConfig.Columns()
	if len(in.Ids) == 0 {
		return gerror.New("请选择要删除的配置")
	}
	if _, err = dao.PledgePeriodConfig.Ctx(ctx).WhereIn(cols.Id, in.Ids).Delete(); err != nil {
		return gerror.Wrap(err, "删除周期配置失败")
	}
	return
}

// Edit 修改/新增周期配置
func (s *sSysPeriodConfig) Edit(ctx context.Context, in *sysin.PeriodConfigEditInp) (err error) {
	cols := dao.PledgePeriodConfig.Columns()
	mod := dao.PledgePeriodConfig.Ctx(ctx).Where(cols.PeriodMonths, in.PeriodMonths)
	if in.Id > 0 {
		mod = mod.WhereNot(cols.Id, in.Id)
	}
	exist, err := mod.Exist()
	if err != nil {
		return gerror.Wrap(err, "查询周期配置失败")
	}
	if exist {
		return gerror.New("该投资周期已存在配置")
	}

	if in.Id > 0 {
		_, err = dao.PledgePeriodConfig.Ctx(ctx).
			Fields(
				cols.PeriodMonths,
				cols.MonthlyRewardInfo,
				cols.MonthlyRewardRate,
				cols.DirectRewardRate,
				cols.Status,
			).
			WherePri(in.Id).
			Data(in).
			Update()
		if err != nil {
			return gerror.Wrap(err, "更新周期配置失败")
		}
		return
	}

	_, err = dao.PledgePeriodConfig.Ctx(ctx).Data(in).Insert()
	if err != nil {
		return gerror.Wrap(err, "新增周期配置失败")
	}
	return
}

// Status 更新周期配置状态
func (s *sSysPeriodConfig) Status(ctx context.Context, in *sysin.PeriodConfigStatusInp) (err error) {
	cols := dao.PledgePeriodConfig.Columns()
	_, err = dao.PledgePeriodConfig.Ctx(ctx).
		WherePri(in.Id).
		Data(cols.Status, in.Status).
		Update()
	if err != nil {
		return gerror.Wrap(err, "更新周期配置状态失败")
	}
	return
}

// View 获取周期配置详情
func (s *sSysPeriodConfig) View(ctx context.Context, in *sysin.PeriodConfigViewInp) (res *sysin.PeriodConfigViewModel, err error) {
	var row *entity.PledgePeriodConfig
	if err = dao.PledgePeriodConfig.Ctx(ctx).WherePri(in.Id).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "查询周期配置失败")
	}
	if row == nil {
		return nil, gerror.New("周期配置不存在")
	}
	res = &sysin.PeriodConfigViewModel{PledgePeriodConfig: *row}
	return
}

// List 获取周期配置列表
func (s *sSysPeriodConfig) List(ctx context.Context, in *sysin.PeriodConfigListInp) (list []*sysin.PeriodConfigListModel, totalCount int, err error) {
	cols := dao.PledgePeriodConfig.Columns()
	mod := dao.PledgePeriodConfig.Ctx(ctx)
	if in.PeriodMonths > 0 {
		mod = mod.Where(cols.PeriodMonths, in.PeriodMonths)
	}
	if in.Status > 0 {
		mod = mod.Where(cols.Status, in.Status)
	}

	totalCount, err = mod.Clone().Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, "查询周期配置总数失败")
	}
	if totalCount == 0 {
		return
	}

	if err = mod.Page(in.Page, in.PerPage).OrderAsc(cols.PeriodMonths).Scan(&list); err != nil {
		return nil, 0, gerror.Wrap(err, "查询周期配置列表失败")
	}
	return
}

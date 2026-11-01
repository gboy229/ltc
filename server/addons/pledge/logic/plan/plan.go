// Package plan
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package plan

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	apiPlan "hotgo/addons/pledge/api/api/plan"
	"hotgo/addons/pledge/model/dao"
	"hotgo/addons/pledge/model/entity"
	"hotgo/addons/pledge/service"
	"hotgo/internal/consts"
)

type sPlan struct{}

func NewPlan() *sPlan {
	return &sPlan{}
}

func init() {
	service.RegisterPlan(NewPlan())
}

// List 获取启用的投资套餐
func (s *sPlan) List(ctx context.Context, _ *apiPlan.ListReq) (res *apiPlan.ListRes, err error) {
	var plans []*entity.PledgePlan
	if err = dao.PledgePlan.Ctx(ctx).
		Where(dao.PledgePlan.Columns().Status, consts.StatusEnabled).
		OrderAsc(dao.PledgePlan.Columns().Name).
		OrderAsc(dao.PledgePlan.Columns().PeriodMonths).
		Scan(&plans); err != nil {
		return nil, gerror.Wrap(err, "查询投资套餐失败")
	}

	periodConfigs, err := s.periodConfigs(ctx)
	if err != nil {
		return nil, err
	}

	res = new(apiPlan.ListRes)
	res.List = make([]*apiPlan.Item, 0, len(plans))
	for _, plan := range plans {
		if item := calcPlanQuote(plan, periodConfigs[plan.PeriodMonths]); item != nil {
			res.List = append(res.List, item)
		}
	}
	return
}

func (s *sPlan) periodConfigs(ctx context.Context) (map[int]*entity.PledgePeriodConfig, error) {
	cols := dao.PledgePeriodConfig.Columns()
	var rows []*entity.PledgePeriodConfig
	if err := dao.PledgePeriodConfig.Ctx(ctx).
		Where(cols.Status, consts.StatusEnabled).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "查询周期配置失败")
	}
	result := make(map[int]*entity.PledgePeriodConfig, len(rows))
	for _, row := range rows {
		result[row.PeriodMonths] = row
	}
	return result, nil
}

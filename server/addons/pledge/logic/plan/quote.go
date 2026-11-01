// Package plan
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package plan

import (
	"github.com/shopspring/decimal"

	apiPlan "hotgo/addons/pledge/api/api/plan"
	"hotgo/addons/pledge/model/entity"
)

// calcPlanQuote 计算前台套餐收益试算结果
func calcPlanQuote(plan *entity.PledgePlan, cfg *entity.PledgePeriodConfig) *apiPlan.Item {
	if plan == nil || cfg == nil {
		return nil
	}

	var (
		hundred       = decimal.NewFromInt(100)
		amount        = decimal.NewFromFloat(plan.Amount)
		monthlyRate   = decimal.NewFromFloat(cfg.MonthlyRewardRate)
		periods       = decimal.NewFromInt(int64(plan.PeriodMonths))
		monthlyReward = amount.Mul(monthlyRate).Div(hundred).Round(2)
		totalReward   = monthlyReward.Mul(periods).Round(2)
		totalReturn   = amount.Add(totalReward).Round(2)
	)

	return &apiPlan.Item{
		Name:              plan.Name,
		Amount:            amount.InexactFloat64(),
		PeriodMonths:      plan.PeriodMonths,
		MonthlyRewardInfo: cfg.MonthlyRewardInfo,
		MonthlyRewardRate: cfg.MonthlyRewardRate,
		TotalReturn:       totalReturn.InexactFloat64(),
	}
}

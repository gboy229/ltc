// Package plan
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package plan

import (
	"testing"

	"github.com/shopspring/decimal"

	"hotgo/addons/pledge/model/entity"
)

func money(v float64) string {
	return decimal.NewFromFloat(v).Round(2).StringFixed(2)
}

func TestCalcPlanQuote(t *testing.T) {
	tests := []struct {
		name string
		plan *entity.PledgePlan
		want string
	}{
		{
			name: "V1-3个月",
			plan: &entity.PledgePlan{Id: 1, Name: "V1", Amount: 1000, PeriodMonths: 3, Status: 1},
			want: "1270.00",
		},
		{
			name: "V5-12个月",
			plan: &entity.PledgePlan{Id: 15, Name: "V5", Amount: 10000, PeriodMonths: 12, Status: 1},
			want: "26200.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcPlanQuote(tt.plan, &entity.PledgePeriodConfig{PeriodMonths: tt.plan.PeriodMonths, MonthlyRewardInfo: "8%～10%", MonthlyRewardRate: map[int]float64{3: 9, 12: 13.5}[tt.plan.PeriodMonths]})
			if got == nil {
				t.Fatal("calcPlanQuote returned nil")
			}
			if money(got.TotalReturn) != tt.want {
				t.Fatalf("totalReturn = %s, want %s", money(got.TotalReturn), tt.want)
			}
		})
	}

	if got := calcPlanQuote(nil, nil); got != nil {
		t.Fatal("calcPlanQuote(nil, nil) should return nil")
	}
}

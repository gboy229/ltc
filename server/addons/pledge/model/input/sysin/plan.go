// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/addons/pledge/model/entity"
	"hotgo/internal/consts"
	"hotgo/internal/model/input/form"
)

// PlanEditInp 修改/新增投资套餐
type PlanEditInp struct {
	entity.PledgePlan
}

func (in *PlanEditInp) Filter(ctx context.Context) (err error) {
	in.Name = gstr.Trim(in.Name)
	if in.Name == "" {
		return gerror.New("套餐名称不能为空")
	}
	if in.Amount <= 0 {
		return gerror.New("入金金额必须大于0")
	}
	if in.PeriodMonths < 1 || in.PeriodMonths > 120 {
		return gerror.New("投资周期必须在1~120个月之间")
	}
	if in.DailyLots < 0 || in.DailyLots > 100000 {
		return gerror.New("每日交易手数必须在0~100000之间")
	}
	if in.RewardPerLot < 0 {
		return gerror.New("每手金额不能小于0")
	}
	if in.EarlyExitPenaltyRate < 0 || in.EarlyExitPenaltyRate > 100 {
		return gerror.New("提前退本违约金比例必须在0~100之间")
	}
	if in.Status == 0 {
		in.Status = consts.StatusEnabled
	}
	if in.Status != consts.StatusEnabled && in.Status != consts.StatusDisable {
		return gerror.New("状态值不合法")
	}
	return
}

type PlanEditModel struct{}

// PlanUpdateFields 修改投资套餐字段过滤
type PlanUpdateFields struct {
	Name                 string  `json:"name"                 description:"套餐名称"`
	Amount               float64 `json:"amount"               description:"入金金额"`
	PeriodMonths         int     `json:"periodMonths"         description:"投资周期（月）"`
	DailyLots            int     `json:"dailyLots"            description:"每日交易手数"`
	RewardPerLot         float64 `json:"rewardPerLot"         description:"每手金额"`
	EarlyExitPenaltyRate float64 `json:"earlyExitPenaltyRate" description:"提前退本违约金比例（%）"`
	Status               int     `json:"status"               description:"状态"`
}

// PlanInsertFields 新增投资套餐字段过滤
type PlanInsertFields struct {
	Name                 string  `json:"name"                 description:"套餐名称"`
	Amount               float64 `json:"amount"               description:"入金金额"`
	PeriodMonths         int     `json:"periodMonths"         description:"投资周期（月）"`
	DailyLots            int     `json:"dailyLots"            description:"每日交易手数"`
	RewardPerLot         float64 `json:"rewardPerLot"         description:"每手金额"`
	EarlyExitPenaltyRate float64 `json:"earlyExitPenaltyRate" description:"提前退本违约金比例（%）"`
	Status               int     `json:"status"               description:"状态"`
}

// PlanDeleteInp 删除投资套餐
type PlanDeleteInp struct {
	Ids []int64 `json:"ids" v:"required#请选择要删除的套餐" dc:"套餐ID"`
}

func (in *PlanDeleteInp) Filter(ctx context.Context) (err error) {
	if len(in.Ids) == 0 {
		return gerror.New("请选择要删除的套餐")
	}
	return
}

type PlanDeleteModel struct{}

// PlanViewInp 获取指定投资套餐信息
type PlanViewInp struct {
	Id int64 `json:"id" v:"required#套餐ID不能为空" dc:"套餐ID"`
}

type PlanViewModel struct {
	entity.PledgePlan
}

// PlanListInp 获取投资套餐列表
type PlanListInp struct {
	form.PageReq
	form.StatusReq
	Name         string `json:"name"         dc:"套餐名称"`
	PeriodMonths int    `json:"periodMonths" dc:"投资周期（月）"`
}

type PlanListModel struct {
	entity.PledgePlan
}

// PlanStatusInp 更新投资套餐状态
type PlanStatusInp struct {
	Id     int64 `json:"id" v:"required#套餐ID不能为空" dc:"套餐ID"`
	Status int   `json:"status" dc:"状态"`
}

func (in *PlanStatusInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		return gerror.New("套餐ID不能为空")
	}
	if in.Status != consts.StatusEnabled && in.Status != consts.StatusDisable {
		return gerror.New("状态值不合法")
	}
	return
}

type PlanStatusModel struct{}

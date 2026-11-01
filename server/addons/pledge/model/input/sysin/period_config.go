// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/addons/pledge/model/entity"
	"hotgo/internal/consts"
	"hotgo/internal/model/input/form"
)

// PeriodConfigListInp 周期配置列表
type PeriodConfigListInp struct {
	form.PageReq
	PeriodMonths int `json:"periodMonths" dc:"投资周期（月）"`
	Status       int `json:"status"       dc:"状态"`
}

type PeriodConfigListModel struct {
	entity.PledgePeriodConfig
}

// PeriodConfigViewInp 周期配置详情
type PeriodConfigViewInp struct {
	Id int64 `json:"id" v:"required#ID不能为空" dc:"ID"`
}

type PeriodConfigViewModel struct {
	entity.PledgePeriodConfig
}

// PeriodConfigEditInp 修改/新增周期配置
type PeriodConfigEditInp struct {
	entity.PledgePeriodConfig
}

func (in *PeriodConfigEditInp) Filter(ctx context.Context) (err error) {
	in.MonthlyRewardInfo = gstr.Trim(in.MonthlyRewardInfo)
	if in.PeriodMonths < 1 || in.PeriodMonths > 120 {
		return gerror.New("投资周期必须在1~120个月之间")
	}
	if in.MonthlyRewardRate <= 0 || in.MonthlyRewardRate > 100 {
		return gerror.New("月化收益率必须大于0且不超过100")
	}
	if in.DirectRewardRate < 0 || in.DirectRewardRate > 100 {
		return gerror.New("直推收益率必须在0~100之间")
	}
	if gstr.LenRune(in.MonthlyRewardInfo) > 64 {
		return gerror.New("月化收益展示区间不能超过64个字符")
	}
	if in.MonthlyRewardInfo == "" {
		in.MonthlyRewardInfo = fmt.Sprintf("%s%%", strconv.FormatFloat(in.MonthlyRewardRate, 'f', -1, 64))
	}
	if in.Status == 0 {
		in.Status = consts.StatusEnabled
	}
	if in.Status != consts.StatusEnabled && in.Status != consts.StatusDisable {
		return gerror.New("状态值不合法")
	}
	return
}

// PeriodConfigDeleteInp 删除周期配置
type PeriodConfigDeleteInp struct {
	Ids []int64 `json:"ids" v:"required#请选择要删除的配置" dc:"ID列表"`
}

// PeriodConfigStatusInp 更新周期配置状态
type PeriodConfigStatusInp struct {
	Id     int64 `json:"id"     v:"required#ID不能为空" dc:"ID"`
	Status int   `json:"status" dc:"状态"`
}

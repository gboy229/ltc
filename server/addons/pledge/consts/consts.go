// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import (
	"hotgo/internal/library/dict"
	"hotgo/internal/model"
)

// 正式订单状态
const (
	PledgeOrderStatusActive   = 1 // 进行中
	PledgeOrderStatusFinished = 2 // 已完成

	PledgeOrderIsEarlyRefundNo  = 0 // 未提前退本
	PledgeOrderIsEarlyRefundYes = 1 // 已提前退本
)

// 投资周期（月）
const (
	PledgePeriod3  = 3
	PledgePeriod6  = 6
	PledgePeriod12 = 12
)

// PledgePeriodOptions 投资周期选项
var PledgePeriodOptions = []*model.Option{
	dict.GenPrimaryOption(PledgePeriod3, "3个月"),
	dict.GenPrimaryOption(PledgePeriod6, "6个月"),
	dict.GenPrimaryOption(PledgePeriod12, "12个月"),
}

// 收益类型
const (
	PledgeIncomeTypeStatic    = "static"    // 静态收益
	PledgeIncomeTypeDirect    = "direct"    // 直推收益
	PledgeIncomeTypeDiff      = "diff"      // 级差收益
	PledgeIncomeTypeLot       = "lot"       // 手数收益
	PledgeIncomeTypeIndirect  = "indirect"  // 间推收益
	PledgeIncomeTypeTeam      = "team"      // 团队收益
	PledgeIncomeTypePrincipal = "principal" // 本金返还
)

// PledgeIncomeTypeOptions 收益类型选项
var PledgeIncomeTypeOptions = []*model.Option{
	dict.GenSuccessOption(PledgeIncomeTypeStatic, "静态收益"),
	dict.GenSuccessOption(PledgeIncomeTypeDirect, "直推收益"),
	dict.GenSuccessOption(PledgeIncomeTypeLot, "手数收益"),
	dict.GenPrimaryOption(PledgeIncomeTypeTeam, "团队收益"),
	dict.GenInfoOption(PledgeIncomeTypePrincipal, "本金返还"),
}

// 收益状态
const (
	PledgeIncomeStatusCredited = 1 // 已入账
	PledgeIncomeStatusRevoked  = 2 // 已撤销
)

func init() {
	dict.RegisterEnums("PledgePeriodOptions", "投资周期选项", PledgePeriodOptions)
	dict.RegisterEnums("PledgeIncomeTypeOptions", "收益类型选项", PledgeIncomeTypeOptions)
}

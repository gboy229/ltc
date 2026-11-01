// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgeOrder is the golang structure for table pledge_order.
type PledgeOrder struct {
	Id                   uint64      `json:"id"                   orm:"id"                      description:"主键ID"`
	OrderNo              string      `json:"orderNo"              orm:"order_no"                description:"订单号（与报单号一致）"`
	UserId               uint64      `json:"userId"               orm:"user_id"                 description:"用户ID"`
	PlanId               uint64      `json:"planId"               orm:"plan_id"                 description:"套餐ID"`
	Amount               float64     `json:"amount"               orm:"amount"                  description:"订单金额"`
	PeriodMonths         int         `json:"periodMonths"         orm:"period_months"           description:""`
	TodayIncome          float64     `json:"todayIncome"          orm:"today_income"            description:"今天收益"`
	IncIncome            float64     `json:"incIncome"            orm:"inc_income"              description:"累计收益"`
	CycleNum             int         `json:"cycleNum"             orm:"cycle_num"               description:"总天数"`
	IncIncomeNum         int         `json:"incIncomeNum"         orm:"inc_income_num"          description:"累计收益次数"`
	EarlyExitPenaltyRate float64     `json:"earlyExitPenaltyRate" orm:"early_exit_penalty_rate" description:"提前退本违约金比例（%）"`
	RechargeId           uint64      `json:"rechargeId"           orm:"recharge_id"             description:"关联充值记录ID"`
	Status               int         `json:"status"               orm:"status"                  description:"状态：1=进行中"`
	IsEarlyRefund        int         `json:"isEarlyRefund"        orm:"is_early_refund"         description:"是否提前退本：0=否，1=是"`
	SettleAt             *gtime.Time `json:"settleAt"             orm:"settle_at"               description:"最后结算时间"`
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"              description:"创建时间"`
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"              description:"更新时间"`
}

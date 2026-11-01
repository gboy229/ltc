// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgePlan is the golang structure for table pledge_plan.
type PledgePlan struct {
	Id                   uint64      `json:"id"                   orm:"id"                      description:""`
	Name                 string      `json:"name"                 orm:"name"                    description:"å¥—é¤åç§°"`
	Amount               float64     `json:"amount"               orm:"amount"                  description:"å…¥é‡‘é‡‘é¢"`
	PeriodMonths         int         `json:"periodMonths"         orm:"period_months"           description:"投资周期（月）"`
	DailyLots            int         `json:"dailyLots"            orm:"daily_lots"              description:"每日交易手数"`
	RewardPerLot         float64     `json:"rewardPerLot"         orm:"reward_per_lot"          description:"每手金额（USDT）"`
	EarlyExitPenaltyRate float64     `json:"earlyExitPenaltyRate" orm:"early_exit_penalty_rate" description:"提前退本违约金比例（%）"`
	Status               int         `json:"status"               orm:"status"                  description:"状态：1=启用，2=禁用"`
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"              description:"创建时间"`
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"              description:"更新时间"`
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgePeriodConfig is the golang structure for table pledge_period_config.
type PledgePeriodConfig struct {
	Id                uint64      `json:"id"                orm:"id"                  description:""`
	PeriodMonths      int         `json:"periodMonths"      orm:"period_months"       description:""`
	MonthlyRewardInfo string      `json:"monthlyRewardInfo" orm:"monthly_reward_info" description:""`
	MonthlyRewardRate float64     `json:"monthlyRewardRate" orm:"monthly_reward_rate" description:""`
	DirectRewardRate  float64     `json:"directRewardRate"  orm:"direct_reward_rate"  description:""`
	Status            int         `json:"status"            orm:"status"              description:""`
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          description:""`
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          description:""`
}

// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgePeriodConfig is the golang structure of table hg_addon_pledge_period_config for DAO operations like Where/Data.
type PledgePeriodConfig struct {
	g.Meta            `orm:"table:hg_addon_pledge_period_config, do:true"`
	Id                interface{} //
	PeriodMonths      interface{} //
	MonthlyRewardInfo interface{} //
	MonthlyRewardRate interface{} //
	DirectRewardRate  interface{} //
	Status            interface{} //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
}

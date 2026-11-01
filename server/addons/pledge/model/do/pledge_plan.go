// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgePlan is the golang structure of table hg_addon_pledge_plan for DAO operations like Where/Data.
type PledgePlan struct {
	g.Meta               `orm:"table:hg_addon_pledge_plan, do:true"`
	Id                   interface{} //
	Name                 interface{} // å¥—é¤åç§°
	Amount               interface{} // å…¥é‡‘é‡‘é¢
	PeriodMonths         interface{} // 投资周期（月）
	DailyLots            interface{} // 每日交易手数
	RewardPerLot         interface{} // 每手金额（USDT）
	EarlyExitPenaltyRate interface{} // 提前退本违约金比例（%）
	Status               interface{} // 状态：1=启用，2=禁用
	CreatedAt            *gtime.Time // 创建时间
	UpdatedAt            *gtime.Time // 更新时间
}

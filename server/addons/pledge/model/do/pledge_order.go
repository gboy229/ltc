// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgeOrder is the golang structure of table hg_addon_pledge_order for DAO operations like Where/Data.
type PledgeOrder struct {
	g.Meta               `orm:"table:hg_addon_pledge_order, do:true"`
	Id                   interface{} // 主键ID
	OrderNo              interface{} // 订单号（与报单号一致）
	UserId               interface{} // 用户ID
	PlanId               interface{} // 套餐ID
	Amount               interface{} // 订单金额
	PeriodMonths         interface{} //
	TodayIncome          interface{} // 今天收益
	IncIncome            interface{} // 累计收益
	CycleNum             interface{} // 总天数
	IncIncomeNum         interface{} // 累计收益次数
	EarlyExitPenaltyRate interface{} // 提前退本违约金比例（%）
	RechargeId           interface{} // 关联充值记录ID
	Status               interface{} // 状态：1=进行中
	IsEarlyRefund        interface{} // 是否提前退本：0=否，1=是
	SettleAt             *gtime.Time // 最后结算时间
	CreatedAt            *gtime.Time // 创建时间
	UpdatedAt            *gtime.Time // 更新时间
}

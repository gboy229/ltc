// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgeIncomeLog is the golang structure of table hg_addon_pledge_income_log for DAO operations like Where/Data.
type PledgeIncomeLog struct {
	g.Meta     `orm:"table:hg_addon_pledge_income_log, do:true"`
	Id         interface{} //
	UserId     interface{} // 收益用户
	OrderId    interface{} // 来源报单
	FromUserId interface{} // 来源用户（直推/团队）
	Type       interface{} // static/direct/indirect/team/...
	Amount     interface{} // 收益金额
	BaseAmount interface{} // 计算基数
	Rate       interface{} // 收益率
	PeriodNo   interface{} // 第几期
	Status     interface{} // 1=已入账 2=已撤销
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
}

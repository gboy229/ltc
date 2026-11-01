// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgeAddOrder is the golang structure of table hg_addon_pledge_add_order for DAO operations like Where/Data.
type PledgeAddOrder struct {
	g.Meta    `orm:"table:hg_addon_pledge_add_order, do:true"`
	Id        interface{} // 主键ID
	OrderNo   interface{} // 公开订单号
	UserId    interface{} // 用户ID
	PlanId    interface{} // 套餐ID
	Amount    interface{} // 报单金额
	Status    interface{} // 状态：0=待处理，1=已处理
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}

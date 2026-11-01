// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgeAddOrder is the golang structure for table pledge_add_order.
type PledgeAddOrder struct {
	Id        uint64      `json:"id"        orm:"id"         description:"主键ID"`
	OrderNo   string      `json:"orderNo"   orm:"order_no"   description:"公开订单号"`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"用户ID"`
	PlanId    uint64      `json:"planId"    orm:"plan_id"    description:"套餐ID"`
	Amount    float64     `json:"amount"    orm:"amount"     description:"报单金额"`
	Status    int         `json:"status"    orm:"status"     description:"状态：0=待处理，1=已处理"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`
}

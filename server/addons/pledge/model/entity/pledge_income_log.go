// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PledgeIncomeLog is the golang structure for table pledge_income_log.
type PledgeIncomeLog struct {
	Id         uint64      `json:"id"         orm:"id"           description:""`
	UserId     uint64      `json:"userId"     orm:"user_id"      description:"收益用户"`
	OrderId    uint64      `json:"orderId"    orm:"order_id"     description:"来源报单"`
	FromUserId uint64      `json:"fromUserId" orm:"from_user_id" description:"来源用户（直推/团队）"`
	Type       string      `json:"type"       orm:"type"         description:"static/direct/indirect/team/..."`
	Amount     float64     `json:"amount"     orm:"amount"       description:"收益金额"`
	BaseAmount float64     `json:"baseAmount" orm:"base_amount"  description:"计算基数"`
	Rate       float64     `json:"rate"       orm:"rate"         description:"收益率"`
	PeriodNo   int         `json:"periodNo"   orm:"period_no"    description:"第几期"`
	Status     int         `json:"status"     orm:"status"       description:"1=已入账 2=已撤销"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"   description:""`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"   description:""`
}

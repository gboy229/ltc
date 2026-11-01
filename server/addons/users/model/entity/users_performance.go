// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersPerformance is the golang structure for table users_performance.
type UsersPerformance struct {
	UserId         uint64      `json:"userId"         orm:"user_id"         description:"用户ID"`
	PersonalAmount float64     `json:"personalAmount" orm:"personal_amount" description:"个人业绩"`
	DirectAmount   float64     `json:"directAmount"   orm:"direct_amount"   description:"直推业绩"`
	TeamAmount     float64     `json:"teamAmount"     orm:"team_amount"     description:"团队业绩"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:"更新时间"`
}

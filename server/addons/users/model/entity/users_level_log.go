// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersLevelLog is the golang structure for table users_level_log.
type UsersLevelLog struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`
	UserId         uint64      `json:"userId"         orm:"user_id"         description:"用户ID"`
	FromLevelId    int         `json:"fromLevelId"    orm:"from_level_id"   description:"变更前等级ID"`
	ToLevelId      int         `json:"toLevelId"      orm:"to_level_id"     description:"变更后等级ID"`
	PersonalAmount float64     `json:"personalAmount" orm:"personal_amount" description:"变更时个人业绩"`
	TeamAmount     float64     `json:"teamAmount"     orm:"team_amount"     description:"变更时团队业绩"`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`
}

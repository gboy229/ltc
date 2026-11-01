// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersLevel is the golang structure for table users_level.
type UsersLevel struct {
	Id             uint64      `json:"id"             orm:"id"              description:"主键ID"`
	Name           string      `json:"name"           orm:"name"            description:"等级名称，例如 LV1、新手"`
	PersonalAmount float64     `json:"personalAmount" orm:"personal_amount" description:"个人业绩"`
	TeamAmount     float64     `json:"teamAmount"     orm:"team_amount"     description:"团队业绩"`
	DiffReward     float64     `json:"diffReward"     orm:"diff_reward"     description:"级差奖励"`
	Remark         string      `json:"remark"         orm:"remark"          description:"备注说明"`
	Status         int         `json:"status"         orm:"status"          description:"状态：1=启用，2=禁用"`
	Sort           int         `json:"sort"           orm:"sort"            description:"排序值，越大越靠前"`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:"更新时间"`
}

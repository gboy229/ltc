// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersLevel is the golang structure of table hg_addon_users_level for DAO operations like Where/Data.
type UsersLevel struct {
	g.Meta         `orm:"table:hg_addon_users_level, do:true"`
	Id             interface{} // 主键ID
	Name           interface{} // 等级名称，例如 LV1、新手
	PersonalAmount interface{} // 个人业绩
	TeamAmount     interface{} // 团队业绩
	DiffReward     interface{} // 级差奖励
	Remark         interface{} // 备注说明
	Status         interface{} // 状态：1=启用，2=禁用
	Sort           interface{} // 排序值，越大越靠前
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
}

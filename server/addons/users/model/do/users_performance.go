// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersPerformance is the golang structure of table hg_addon_users_performance for DAO operations like Where/Data.
type UsersPerformance struct {
	g.Meta         `orm:"table:hg_addon_users_performance, do:true"`
	UserId         interface{} // 用户ID
	PersonalAmount interface{} // 个人业绩
	DirectAmount   interface{} // 直推业绩
	TeamAmount     interface{} // 团队业绩
	UpdatedAt      *gtime.Time // 更新时间
}

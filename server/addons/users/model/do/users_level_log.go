// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersLevelLog is the golang structure of table hg_addon_users_level_log for DAO operations like Where/Data.
type UsersLevelLog struct {
	g.Meta         `orm:"table:hg_addon_users_level_log, do:true"`
	Id             interface{} //
	UserId         interface{} // 用户ID
	FromLevelId    interface{} // 变更前等级ID
	ToLevelId      interface{} // 变更后等级ID
	PersonalAmount interface{} // 变更时个人业绩
	TeamAmount     interface{} // 变更时团队业绩
	CreatedAt      *gtime.Time // 创建时间
}

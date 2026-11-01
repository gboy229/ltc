// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UsersTree is the golang structure of table hg_addon_users_tree for DAO operations like Where/Data.
type UsersTree struct {
	g.Meta `orm:"table:hg_addon_users_tree, do:true"`
	Id     interface{} //
	Uid    interface{} // uid
	Pid    interface{} // pid
	Level  interface{} // level
}

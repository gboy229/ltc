// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainAddress is the golang structure of table hg_addon_chain_address for DAO operations like Where/Data.
type ChainAddress struct {
	g.Meta    `orm:"table:hg_addon_chain_address, do:true"`
	Id        interface{} //
	Chain     interface{} // 链名
	Uid       interface{} // 用户ID
	Address   interface{} // 地址
	Keys      interface{} // 私钥
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}

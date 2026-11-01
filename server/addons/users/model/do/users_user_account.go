// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersUserAccount is the golang structure of table hg_addon_users_user_account for DAO operations like Where/Data.
type UsersUserAccount struct {
	g.Meta     `orm:"table:hg_addon_users_user_account, do:true"`
	Id         interface{} // 凭证ID
	UserId     interface{} // 用户ID
	Type       interface{} // 凭证类型:username,mobile,email,wallet,wechat,google...
	Account    interface{} // 凭证值(手机号/邮箱/用户名/钱包地址/openid)
	Credential interface{} // 凭证附加信息(可选,如钱包签名参数等)
	Verified   interface{} // 是否已验证(手机/邮箱/钱包通常需验证后才可登录)
	Status     interface{} // 状态
	CreatedAt  *gtime.Time // 创建时间
	UpdatedAt  *gtime.Time // 更新时间
	DeletedAt  *gtime.Time // 删除时间
}

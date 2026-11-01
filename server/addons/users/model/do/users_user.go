// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersUser is the golang structure of table hg_addon_users_user for DAO operations like Where/Data.
type UsersUser struct {
	g.Meta               `orm:"table:hg_addon_users_user, do:true"`
	Id                   interface{} // 用户ID
	Uid                  interface{} // 对外唯一用户号
	Pid                  interface{} // 上级用户ID(邀请关系)
	LevelId              interface{} // 会员等级ID
	LockLevel            interface{} // 是否锁定等级
	Nickname             interface{} // 昵称
	Avatar               interface{} // 头像
	Password             interface{} // 登录密码(bcrypt哈希,自带盐)
	TradePassword        interface{} // 交易密码(bcrypt哈希,自带盐)
	AuthenticatorKey     interface{} // TOTP验证器密钥
	AuthenticatorEnabled interface{} // 是否开启验证器
	AuthStatus           interface{} // 实名认证状态
	Status               interface{} // 状态
	Extend               *gjson.Json // 扩展资料(零散资料兜底,如性别/生日/地区等)
	LastLoginIp          interface{} // 最后登录IP
	LastLoginTime        *gtime.Time // 最后登录时间
	CreatedAt            *gtime.Time // 创建时间
	UpdatedAt            *gtime.Time // 更新时间
	DeletedAt            *gtime.Time // 删除时间
}

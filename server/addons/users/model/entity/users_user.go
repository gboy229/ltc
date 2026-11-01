// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersUser is the golang structure for table users_user.
type UsersUser struct {
	Id                   int64       `json:"id"                   orm:"id"                    description:"用户ID"`
	Uid                  string      `json:"uid"                  orm:"uid"                   description:"对外唯一用户号"`
	Pid                  int64       `json:"pid"                  orm:"pid"                   description:"上级用户ID(邀请关系)"`
	LevelId              int         `json:"levelId"              orm:"level_id"              description:"会员等级ID"`
	LockLevel            int         `json:"lockLevel"            orm:"lock_level"            description:"是否锁定等级"`
	Nickname             string      `json:"nickname"             orm:"nickname"              description:"昵称"`
	Avatar               string      `json:"avatar"               orm:"avatar"                description:"头像"`
	Password             string      `json:"password"             orm:"password"              description:"登录密码(bcrypt哈希,自带盐)"`
	TradePassword        string      `json:"tradePassword"        orm:"trade_password"        description:"交易密码(bcrypt哈希,自带盐)"`
	AuthenticatorKey     string      `json:"authenticatorKey"     orm:"authenticator_key"     description:"TOTP验证器密钥"`
	AuthenticatorEnabled int         `json:"authenticatorEnabled" orm:"authenticator_enabled" description:"是否开启验证器"`
	AuthStatus           int         `json:"authStatus"           orm:"auth_status"           description:"实名认证状态"`
	Status               int         `json:"status"               orm:"status"                description:"状态"`
	Extend               *gjson.Json `json:"extend"               orm:"extend"                description:"扩展资料(零散资料兜底,如性别/生日/地区等)"`
	LastLoginIp          string      `json:"lastLoginIp"          orm:"last_login_ip"         description:"最后登录IP"`
	LastLoginTime        *gtime.Time `json:"lastLoginTime"        orm:"last_login_time"       description:"最后登录时间"`
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"            description:"创建时间"`
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"            description:"更新时间"`
	DeletedAt            *gtime.Time `json:"deletedAt"            orm:"deleted_at"            description:"删除时间"`
}

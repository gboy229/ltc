// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UsersUserAccount is the golang structure for table users_user_account.
type UsersUserAccount struct {
	Id         int64       `json:"id"         orm:"id"         description:"凭证ID"`
	UserId     int64       `json:"userId"     orm:"user_id"    description:"用户ID"`
	Type       string      `json:"type"       orm:"type"       description:"凭证类型:username,mobile,email,wallet,wechat,google..."`
	Account    string      `json:"account"    orm:"account"    description:"凭证值(手机号/邮箱/用户名/钱包地址/openid)"`
	Credential string      `json:"credential" orm:"credential" description:"凭证附加信息(可选,如钱包签名参数等)"`
	Verified   int         `json:"verified"   orm:"verified"   description:"是否已验证(手机/邮箱/钱包通常需验证后才可登录)"`
	Status     int         `json:"status"     orm:"status"     description:"状态"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at" description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at" description:"更新时间"`
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at" description:"删除时间"`
}

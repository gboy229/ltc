// Package auth
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package auth

import (
	"github.com/gogf/gf/v2/frame/g"
)

// NonceReq 获取钱包签名挑战值
type NonceReq struct {
	g.Meta  `path:"/auth/nonce" method:"get" tags:"用户模块" summary:"获取钱包签名挑战值"`
	Address string `json:"address" v:"required#钱包地址不能为空" dc:"钱包地址"`
	Purpose string `json:"purpose" dc:"用途：login/action，默认login"`
}

// NonceRes 获取钱包签名挑战值响应
type NonceRes struct {
	Nonce string `json:"nonce" dc:"签名挑战值"`
}

// LoginReq 统一登录
type LoginReq struct {
	g.Meta    `path:"/auth/login" method:"post" tags:"用户模块" summary:"统一登录"`
	Type      string `json:"type" v:"required#登录类型不能为空" dc:"登录类型:wallet/username/mobile/email"`
	Account   string `json:"account"   dc:"账号"`
	Signature string `json:"signature" dc:"钱包签名"`
	Nonce     string `json:"nonce"     dc:"签名挑战值"`
	Password  string `json:"password"  dc:"密码"`
	Code      string `json:"code"      dc:"验证码"`
	Uid       string `json:"uid"       dc:"推荐人用户号(可选)"`
}

// LoginRes 统一登录响应
type LoginRes struct {
	Token   string `json:"token"   dc:"登录token"`
	Expires int64  `json:"expires" dc:"登录有效期"`
}

// UserInfo 用户信息
type UserInfo struct {
	Uid       string `json:"uid"       dc:"用户号"`
	Nickname  string `json:"nickname"  dc:"昵称"`
	Avatar    string `json:"avatar"    dc:"头像"`
	Address   string `json:"address"   dc:"钱包地址"`
	LevelName string `json:"levelName" dc:"等级名称"`
}

// StatusReq 查询账号注册状态
type StatusReq struct {
	g.Meta  `path:"/auth/status" method:"post" tags:"用户模块" summary:"查询账号注册状态"`
	Type    string `json:"type"    v:"required#账号类型不能为空" dc:"账号类型:wallet/username/mobile/email"`
	Account string `json:"account" v:"required#账号不能为空" dc:"账号"`
}

// StatusRes 账号注册状态响应
type StatusRes struct {
	Registered  bool `json:"registered"  dc:"是否已注册"`
	IsFirstUser bool `json:"isFirstUser" dc:"是否系统第一个账号"`
	ForceInvite bool `json:"forceInvite" dc:"注册是否强制邀请"`
}

// LogoutReq 退出登录
type LogoutReq struct {
	g.Meta `path:"/auth/logout" method:"post" tags:"用户模块" summary:"退出登录"`
}

// LogoutRes 退出登录响应
type LogoutRes struct{}

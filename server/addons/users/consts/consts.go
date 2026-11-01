// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import (
	"hotgo/internal/library/dict"
	"hotgo/internal/model"
)

// 用户凭证类型（登录方式）
const (
	UserAuthTypeWallet   = "wallet"   // 钱包地址
	UserAuthTypeUsername = "username" // 用户名
	UserAuthTypeMobile   = "mobile"   // 手机号
	UserAuthTypeEmail    = "email"    // 邮箱
)

// 用户实名认证状态
const (
	UserAuthStatusNot = 0 // 未认证
	UserAuthStatusYes = 1 // 已认证
)

// 登录凭证验证状态
const (
	UserVerifiedNo  = 0 // 未验证
	UserVerifiedYes = 1 // 已验证
)

// 用户等级锁定状态
const (
	UserLockLevelNo  = 0 // 未锁定
	UserLockLevelYes = 1 // 已锁定
)

// 缓存 key
const (
	CacheWalletNonce = "users:wallet:nonce:" // 钱包签名挑战值（后接钱包地址）
)

// UserAuthTypeOptions 凭证类型选项
var UserAuthTypeOptions = []*model.Option{
	dict.GenHashOption(UserAuthTypeWallet, "钱包地址"),
	dict.GenHashOption(UserAuthTypeUsername, "用户名"),
	dict.GenHashOption(UserAuthTypeMobile, "手机号"),
	dict.GenHashOption(UserAuthTypeEmail, "邮箱"),
}

// UserAuthStatusOptions 用户实名认证状态选项
var UserAuthStatusOptions = []*model.Option{
	dict.GenDefaultOption(UserAuthStatusNot, "未认证"),
	dict.GenSuccessOption(UserAuthStatusYes, "已认证"),
}

// UserVerifiedOptions 登录凭证验证状态选项
var UserVerifiedOptions = []*model.Option{
	dict.GenDefaultOption(UserVerifiedNo, "未验证"),
	dict.GenSuccessOption(UserVerifiedYes, "已验证"),
}

// UserLockLevelOptions 用户等级锁定状态选项
var UserLockLevelOptions = []*model.Option{
	dict.GenDefaultOption(UserLockLevelNo, "未锁定"),
	dict.GenWarningOption(UserLockLevelYes, "已锁定"),
}

func init() {
	dict.RegisterEnums("UserAuthTypeOptions", "用户凭证类型选项", UserAuthTypeOptions)
	dict.RegisterEnums("UserAuthStatusOptions", "用户实名认证状态选项", UserAuthStatusOptions)
	dict.RegisterEnums("UserVerifiedOptions", "登录凭证验证状态选项", UserVerifiedOptions)
	dict.RegisterEnums("UserLockLevelOptions", "用户等级锁定状态选项", UserLockLevelOptions)
}

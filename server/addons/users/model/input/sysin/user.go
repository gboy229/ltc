// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"hotgo/addons/users/model/entity"
	"hotgo/internal/model/input/form"
)

// UserEditInp 修改/新增用户
type UserEditInp struct {
	entity.UsersUser
}

func (in *UserEditInp) Filter(ctx context.Context) (err error) {
	if in.Nickname == "" {
		err = gerror.New("昵称不能为空")
		return
	}
	return
}

type UserEditModel struct {
	Id int64 `json:"id" dc:"用户ID"`
}

// UserUpdateFields 修改用户字段过滤（仅允许修改昵称、等级、状态、是否锁定等级）
type UserUpdateFields struct {
	Nickname  string `json:"nickname"  description:"昵称"`
	LevelId   int    `json:"levelId"   description:"会员等级ID"`
	Status    int    `json:"status"    description:"状态"`
	LockLevel int    `json:"lockLevel" description:"是否锁定等级"`
}

// UserInsertFields 新增用户字段过滤
type UserInsertFields struct {
	Uid       string `json:"uid"       description:"对外唯一用户号"`
	Nickname  string `json:"nickname"  description:"昵称"`
	LevelId   int    `json:"levelId"   description:"会员等级ID"`
	Status    int    `json:"status"    description:"状态"`
	LockLevel int    `json:"lockLevel" description:"是否锁定等级"`
}

// UserDeleteInp 删除用户
type UserDeleteInp struct {
	Id int64 `json:"id" v:"required#用户ID不能为空" dc:"用户ID"`
}

type UserDeleteModel struct{}

// UserViewInp 获取指定用户信息
type UserViewInp struct {
	Id int64 `json:"id" v:"required#用户ID不能为空" dc:"用户ID"`
}

// UserListModel 用户列表返回
type UserListModel struct {
	Id             int64       `json:"id"             dc:"用户ID"`
	Uid            string      `json:"uid"            dc:"对外唯一用户号"`
	Pid            int64       `json:"pid"            dc:"上级用户ID"`
	ParentUid      string      `json:"parentUid"      dc:"推荐人用户号"`
	ParentNickname string      `json:"parentNickname" dc:"推荐人昵称"`
	LevelId        int         `json:"levelId"        dc:"会员等级ID"`
	LevelName      string      `json:"levelName"      dc:"等级名称"`
	LockLevel      int         `json:"lockLevel"      dc:"是否锁定等级"`
	WalletAddress  string      `json:"walletAddress"  dc:"钱包地址"`
	PersonalAmount float64     `json:"personalAmount" dc:"个人业绩"`
	DirectAmount   float64     `json:"directAmount"   dc:"直推业绩"`
	TeamAmount     float64     `json:"teamAmount"     dc:"团队业绩"`
	Nickname       string      `json:"nickname"       dc:"昵称"`
	Avatar         string      `json:"avatar"         dc:"头像"`
	AuthStatus     int         `json:"authStatus"     dc:"实名认证状态"`
	Status         int         `json:"status"         dc:"状态"`
	LastLoginIp    string      `json:"lastLoginIp"    dc:"最后登录IP"`
	LastLoginTime  *gtime.Time `json:"lastLoginTime"  dc:"最后登录时间"`
	CreatedAt      *gtime.Time `json:"createdAt"      dc:"创建时间"`
}

// UserViewModel 用户详情返回
type UserViewModel struct {
	UserListModel
}

// UserListInp 获取用户列表
type UserListInp struct {
	form.PageReq
	form.StatusReq
	Uid           string  `json:"uid"           dc:"用户号"`
	Nickname      string  `json:"nickname"      dc:"昵称"`
	LevelId       int     `json:"levelId"       dc:"会员等级ID"`
	WalletAddress string  `json:"walletAddress" dc:"钱包地址"`
	ParentWallet  string  `json:"parentWallet"  dc:"推荐人钱包地址"`
	TeamWallet    string  `json:"teamWallet"    dc:"团队成员钱包地址"`
	AuthStatus    int     `json:"authStatus"    dc:"实名认证状态"`
	CreatedAt     []int64 `json:"createdAt"     dc:"注册时间"`
}

// UserStatusInp 更新用户状态
type UserStatusInp struct {
	Id     int64 `json:"id" v:"required#用户ID不能为空" dc:"用户ID"`
	Status int   `json:"status" dc:"状态"`
}

func (in *UserStatusInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		err = gerror.New("用户ID不能为空")
		return
	}
	if in.Status <= 0 {
		err = gerror.New("状态不能为空")
		return
	}
	return
}

type UserStatusModel struct{}

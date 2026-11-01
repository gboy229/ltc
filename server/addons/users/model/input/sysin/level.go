// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/addons/users/model/entity"
	"hotgo/internal/model/input/form"
)

// LevelEditInp 修改/新增用户等级
type LevelEditInp struct {
	entity.UsersLevel
}

func (in *LevelEditInp) Filter(ctx context.Context) (err error) {
	if in.Name == "" {
		return gerror.New("等级名称不能为空")
	}
	if in.PersonalAmount < 0 {
		return gerror.New("个人业绩不能小于0")
	}
	if in.TeamAmount < 0 {
		return gerror.New("团队业绩不能小于0")
	}
	if in.DiffReward < 0 {
		return gerror.New("级差奖励不能小于0")
	}
	if in.Status != 1 && in.Status != 2 {
		return gerror.New("状态值不合法")
	}
	return
}

type LevelEditModel struct{}

// LevelUpdateFields 修改用户等级字段过滤
type LevelUpdateFields struct {
	Name           string  `json:"name"           description:"等级名称"`
	PersonalAmount float64 `json:"personalAmount" description:"个人业绩"`
	TeamAmount     float64 `json:"teamAmount"     description:"团队业绩"`
	DiffReward     float64 `json:"diffReward"     description:"级差奖励"`
	Remark         string  `json:"remark"         description:"备注说明"`
	Status         int     `json:"status"         description:"状态"`
	Sort           int     `json:"sort"           description:"排序值"`
}

// LevelInsertFields 新增用户等级字段过滤
type LevelInsertFields struct {
	Name           string  `json:"name"           description:"等级名称"`
	PersonalAmount float64 `json:"personalAmount" description:"个人业绩"`
	TeamAmount     float64 `json:"teamAmount"     description:"团队业绩"`
	DiffReward     float64 `json:"diffReward"     description:"级差奖励"`
	Remark         string  `json:"remark"         description:"备注说明"`
	Status         int     `json:"status"         description:"状态"`
	Sort           int     `json:"sort"           description:"排序值"`
}

// LevelDeleteInp 删除用户等级
type LevelDeleteInp struct {
	Id int64 `json:"id" v:"required#等级ID不能为空" dc:"等级ID"`
}

type LevelDeleteModel struct{}

// LevelViewInp 获取指定用户等级信息
type LevelViewInp struct {
	Id int64 `json:"id" v:"required#等级ID不能为空" dc:"等级ID"`
}

type LevelViewModel struct {
	entity.UsersLevel
}

// LevelListInp 获取用户等级列表
type LevelListInp struct {
	form.PageReq
	form.StatusReq
	Name string `json:"name" dc:"等级名称"`
}

type LevelListModel struct {
	entity.UsersLevel
}

// LevelStatusInp 更新用户等级状态
type LevelStatusInp struct {
	Id     int64 `json:"id" v:"required#等级ID不能为空" dc:"等级ID"`
	Status int   `json:"status" dc:"状态"`
}

func (in *LevelStatusInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		err = gerror.New("等级ID不能为空")
		return
	}
	if in.Status <= 0 {
		err = gerror.New("状态不能为空")
		return
	}
	return
}

type LevelStatusModel struct{}

// LevelMaxSortInp 获取用户等级最大排序
type LevelMaxSortInp struct{}

type LevelMaxSortModel struct {
	Sort int `json:"sort" dc:"排序"`
}

// LevelOptionInp 获取用户等级选项
type LevelOptionInp struct{}

// LevelOptionModel 用户等级选项
type LevelOptionModel struct {
	Id   uint64 `json:"id"   dc:"等级ID"`
	Name string `json:"name" dc:"等级名称"`
}

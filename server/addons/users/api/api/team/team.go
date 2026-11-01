// Package team
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package team

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/model/input/form"
)

// PushReq 我的直推
type PushReq struct {
	g.Meta `path:"/team/push" method:"get" tags:"用户团队" summary:"我的直推"`
	form.PageReq
}

type PushItem struct {
	Uid          string `json:"uid"          dc:"用户UID"`
	LevelName    string `json:"levelName"    dc:"等级名称"`
	RegisteredAt string `json:"registeredAt" dc:"注册时间"`
}

type PushRes struct {
	form.PageRes
	List []*PushItem `json:"list" dc:"直推成员"`
}

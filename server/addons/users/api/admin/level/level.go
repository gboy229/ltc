// Package level
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package level

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/users/model/input/sysin"
	"hotgo/internal/model/input/form"
)

// ListReq 查询等级列表
type ListReq struct {
	g.Meta `path:"/level/list" method:"get" tags:"用户等级" summary:"获取等级列表"`
	sysin.LevelListInp
}

type ListRes struct {
	List []*sysin.LevelListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// ViewReq 获取指定信息
type ViewReq struct {
	g.Meta `path:"/level/view" method:"get" tags:"用户等级" summary:"获取指定等级信息"`
	sysin.LevelViewInp
}

type ViewRes struct {
	*sysin.LevelViewModel
}

// EditReq 修改/新增
type EditReq struct {
	g.Meta `path:"/level/edit" method:"post" tags:"用户等级" summary:"修改/新增等级"`
	sysin.LevelEditInp
}

type EditRes struct{}

// DeleteReq 删除
type DeleteReq struct {
	g.Meta `path:"/level/delete" method:"post" tags:"用户等级" summary:"删除等级"`
	sysin.LevelDeleteInp
}

type DeleteRes struct{}

// MaxSortReq 最大排序
type MaxSortReq struct {
	g.Meta `path:"/level/maxSort" method:"get" tags:"用户等级" summary:"等级最大排序"`
	sysin.LevelMaxSortInp
}

type MaxSortRes struct {
	*sysin.LevelMaxSortModel
}

// StatusReq 更新状态
type StatusReq struct {
	g.Meta `path:"/level/status" method:"post" tags:"用户等级" summary:"更新等级状态"`
	sysin.LevelStatusInp
}

type StatusRes struct{}

// OptionReq 等级选项
type OptionReq struct {
	g.Meta `path:"/level/option" method:"get" tags:"用户等级" summary:"获取等级选项"`
	sysin.LevelOptionInp
}

type OptionRes struct {
	List []*sysin.LevelOptionModel `json:"list" dc:"数据列表"`
}

// Package plan
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package plan

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/internal/model/input/form"
)

// ListReq 查询投资套餐列表
type ListReq struct {
	g.Meta `path:"/plan/list" method:"get" tags:"投资套餐" summary:"获取投资套餐列表"`
	sysin.PlanListInp
}

type ListRes struct {
	List []*sysin.PlanListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// ViewReq 获取指定投资套餐信息
type ViewReq struct {
	g.Meta `path:"/plan/view" method:"get" tags:"投资套餐" summary:"获取指定投资套餐信息"`
	sysin.PlanViewInp
}

type ViewRes struct {
	*sysin.PlanViewModel
}

// EditReq 修改/新增投资套餐
type EditReq struct {
	g.Meta `path:"/plan/edit" method:"post" tags:"投资套餐" summary:"修改/新增投资套餐"`
	sysin.PlanEditInp
}

type EditRes struct{}

// DeleteReq 删除投资套餐
type DeleteReq struct {
	g.Meta `path:"/plan/delete" method:"post" tags:"投资套餐" summary:"删除投资套餐"`
	sysin.PlanDeleteInp
}

type DeleteRes struct{}

// StatusReq 更新投资套餐状态
type StatusReq struct {
	g.Meta `path:"/plan/status" method:"post" tags:"投资套餐" summary:"更新投资套餐状态"`
	sysin.PlanStatusInp
}

type StatusRes struct{}

// Package periodconfig
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package periodconfig

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/internal/model/input/form"
)

type ListReq struct {
	g.Meta `path:"/periodConfig/list" method:"get" tags:"周期配置" summary:"周期配置列表"`
	sysin.PeriodConfigListInp
}

type ListRes struct {
	List []*sysin.PeriodConfigListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

type ViewReq struct {
	g.Meta `path:"/periodConfig/view" method:"get" tags:"周期配置" summary:"周期配置详情"`
	sysin.PeriodConfigViewInp
}

type ViewRes struct {
	*sysin.PeriodConfigViewModel
}

type EditReq struct {
	g.Meta `path:"/periodConfig/edit" method:"post" tags:"周期配置" summary:"修改/新增周期配置"`
	sysin.PeriodConfigEditInp
}

type EditRes struct{}

type DeleteReq struct {
	g.Meta `path:"/periodConfig/delete" method:"post" tags:"周期配置" summary:"删除周期配置"`
	sysin.PeriodConfigDeleteInp
}

type DeleteRes struct{}

type StatusReq struct {
	g.Meta `path:"/periodConfig/status" method:"post" tags:"周期配置" summary:"更新周期配置状态"`
	sysin.PeriodConfigStatusInp
}

type StatusRes struct{}

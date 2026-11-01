// Package income
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package income

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/model/input/form"
)

// ListReq 收益流水列表
type ListReq struct {
	g.Meta `path:"/income/list" method:"get" tags:"收益流水" summary:"收益流水列表"`
	form.PageReq
	Type string `json:"type" dc:"收益类型：static/direct/diff/lot/indirect/team/principal，不传为全部"`
}

type ListItem struct {
	Type   string  `json:"type"   dc:"收益类型"`
	Label  string  `json:"label"  dc:"收益名称"`
	Date   string  `json:"date"   dc:"发生时间"`
	Amount float64 `json:"amount" dc:"收益金额"`
}

type ListRes struct {
	form.PageRes
	List []*ListItem `json:"list" dc:"收益流水"`
}

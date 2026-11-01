// Package income
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package income

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/internal/model/input/form"
)

type ListReq struct {
	g.Meta `path:"/income/list" method:"get" tags:"收益流水" summary:"收益流水列表"`
	sysin.IncomeListInp
}

type ListRes struct {
	List []*sysin.IncomeListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

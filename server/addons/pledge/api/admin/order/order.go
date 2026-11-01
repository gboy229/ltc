// Package order
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package order

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/internal/model/input/form"
)

type ListReq struct {
	g.Meta `path:"/order/list" method:"get" tags:"质押订单" summary:"订单列表"`
	sysin.OrderListInp
}

type ListRes struct {
	List []*sysin.OrderListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

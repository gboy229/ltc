// Package recharge
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package recharge

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/chain/model/input/sysin"
	"hotgo/internal/model/input/form"
)

// ListReq 查询充值记录列表
type ListReq struct {
	g.Meta `path:"/recharge/list" method:"get" tags:"充值记录" summary:"获取充值记录列表"`
	sysin.RechargeListInp
}

type ListRes struct {
	List []*sysin.RechargeListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// ViewReq 获取指定信息
type ViewReq struct {
	g.Meta `path:"/recharge/view" method:"get" tags:"充值记录" summary:"获取指定充值记录信息"`
	sysin.RechargeViewInp
}

type ViewRes struct {
	*sysin.RechargeViewModel
}

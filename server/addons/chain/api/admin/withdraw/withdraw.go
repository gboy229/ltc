// Package withdraw
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package withdraw

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/chain/model/input/sysin"
	"hotgo/internal/model/input/form"
)

// ListReq 查询提现记录列表
type ListReq struct {
	g.Meta `path:"/withdraw/list" method:"get" tags:"提现记录" summary:"获取提现记录列表"`
	sysin.WithdrawListInp
}

type ListRes struct {
	List []*sysin.WithdrawListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// ViewReq 获取指定信息
type ViewReq struct {
	g.Meta `path:"/withdraw/view" method:"get" tags:"提现记录" summary:"获取指定提现记录信息"`
	sysin.WithdrawViewInp
}

type ViewRes struct {
	*sysin.WithdrawViewModel
}

// StatusReq 提现审核
type StatusReq struct {
	g.Meta `path:"/withdraw/status" method:"post" tags:"提现记录" summary:"提现审核"`
	sysin.WithdrawStatusInp
}

type StatusRes struct{}

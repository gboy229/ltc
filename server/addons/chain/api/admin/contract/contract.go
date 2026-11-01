// Package contract
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package contract

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/chain/model/input/sysin"
	"hotgo/internal/model/input/form"
)

// ListReq 查询合约列表
type ListReq struct {
	g.Meta `path:"/contract/list" method:"get" tags:"合约管理" summary:"获取合约列表"`
	sysin.ContractListInp
}

type ListRes struct {
	List []*sysin.ContractListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// ViewReq 获取指定信息
type ViewReq struct {
	g.Meta `path:"/contract/view" method:"get" tags:"合约管理" summary:"获取指定合约信息"`
	sysin.ContractViewInp
}

type ViewRes struct {
	*sysin.ContractViewModel
}

// EditReq 修改/新增
type EditReq struct {
	g.Meta `path:"/contract/edit" method:"post" tags:"合约管理" summary:"修改/新增合约"`
	sysin.ContractEditInp
}

type EditRes struct{}

// DeleteReq 删除
type DeleteReq struct {
	g.Meta `path:"/contract/delete" method:"post" tags:"合约管理" summary:"删除合约"`
	sysin.ContractDeleteInp
}

type DeleteRes struct{}

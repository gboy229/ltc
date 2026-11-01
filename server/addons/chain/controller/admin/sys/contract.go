// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/chain/api/admin/contract"
	"hotgo/addons/chain/service"
)

var (
	Contract = cContract{}
)

type cContract struct{}

// Delete 删除合约
func (c *cContract) Delete(ctx context.Context, req *contract.DeleteReq) (res *contract.DeleteRes, err error) {
	err = service.SysContract().Delete(ctx, &req.ContractDeleteInp)
	return
}

// Edit 修改/新增合约
func (c *cContract) Edit(ctx context.Context, req *contract.EditReq) (res *contract.EditRes, err error) {
	err = service.SysContract().Edit(ctx, &req.ContractEditInp)
	return
}

// View 获取指定合约信息
func (c *cContract) View(ctx context.Context, req *contract.ViewReq) (res *contract.ViewRes, err error) {
	data, err := service.SysContract().View(ctx, &req.ContractViewInp)
	if err != nil {
		return
	}

	res = new(contract.ViewRes)
	res.ContractViewModel = data
	return
}

// List 查看合约列表
func (c *cContract) List(ctx context.Context, req *contract.ListReq) (res *contract.ListRes, err error) {
	list, totalCount, err := service.SysContract().List(ctx, &req.ContractListInp)
	if err != nil {
		return
	}

	res = new(contract.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

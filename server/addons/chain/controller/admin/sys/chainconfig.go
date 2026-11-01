// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/chain/api/admin/chainconfig"
	"hotgo/addons/chain/service"
)

var (
	ChainConfig = cChainConfig{}
)

type cChainConfig struct{}

// Delete 删除网络配置
func (c *cChainConfig) Delete(ctx context.Context, req *chainconfig.DeleteReq) (res *chainconfig.DeleteRes, err error) {
	err = service.SysChainConfig().Delete(ctx, &req.ChainConfigDeleteInp)
	return
}

// Edit 修改/新增网络配置
func (c *cChainConfig) Edit(ctx context.Context, req *chainconfig.EditReq) (res *chainconfig.EditRes, err error) {
	err = service.SysChainConfig().Edit(ctx, &req.ChainConfigEditInp)
	return
}

// View 获取指定网络配置信息
func (c *cChainConfig) View(ctx context.Context, req *chainconfig.ViewReq) (res *chainconfig.ViewRes, err error) {
	data, err := service.SysChainConfig().View(ctx, &req.ChainConfigViewInp)
	if err != nil {
		return
	}

	res = new(chainconfig.ViewRes)
	res.ChainConfigViewModel = data
	return
}

// List 查看网络配置列表
func (c *cChainConfig) List(ctx context.Context, req *chainconfig.ListReq) (res *chainconfig.ListRes, err error) {
	list, totalCount, err := service.SysChainConfig().List(ctx, &req.ChainConfigListInp)
	if err != nil {
		return
	}

	res = new(chainconfig.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// GenerateAddress 生成提币地址
func (c *cChainConfig) GenerateAddress(ctx context.Context, req *chainconfig.GenerateAddressReq) (res *chainconfig.GenerateAddressRes, err error) {
	data, err := service.SysChainConfig().GenerateAddress(ctx, &req.ChainConfigGenerateAddressInp)
	if err != nil {
		return
	}

	res = new(chainconfig.GenerateAddressRes)
	res.ChainConfigGenerateAddressModel = data
	return
}

// Options 获取网络配置选项
func (c *cChainConfig) Options(ctx context.Context, req *chainconfig.OptionsReq) (res *chainconfig.OptionsRes, err error) {
	list, err := service.SysChainConfig().Options(ctx)
	if err != nil {
		return
	}

	res = new(chainconfig.OptionsRes)
	res.List = list
	return
}

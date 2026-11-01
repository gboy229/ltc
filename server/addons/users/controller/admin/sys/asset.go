// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/users/api/admin/asset"
	"hotgo/addons/users/model/input/sysin"
	"hotgo/addons/users/service"
)

var (
	Asset = cAsset{}
)

type cAsset struct{}

// AccountList 资产账户列表
func (c *cAsset) AccountList(ctx context.Context, req *asset.AccountListReq) (res *asset.AccountListRes, err error) {
	list, totalCount, err := service.SysAsset().AccountList(ctx, &req.AssetAccountListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.AssetAccountListModel{}
	}

	res = new(asset.AccountListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// LogList 资金流水列表
func (c *cAsset) LogList(ctx context.Context, req *asset.LogListReq) (res *asset.LogListRes, err error) {
	list, totalCount, err := service.SysAsset().LogList(ctx, &req.AssetLogListInp)
	if err != nil {
		return
	}
	if list == nil {
		list = []*sysin.AssetLogListModel{}
	}

	res = new(asset.LogListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

// Stat 资产统计
func (c *cAsset) Stat(ctx context.Context, req *asset.StatReq) (res *asset.StatRes, err error) {
	res = new(asset.StatRes)
	res.AssetStatModel, err = service.SysAsset().Stat(ctx, &req.AssetStatInp)
	return
}

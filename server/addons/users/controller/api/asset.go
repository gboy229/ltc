// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"
	"math"

	"github.com/gogf/gf/v2/errors/gerror"

	"hotgo/addons/users/api/api/asset"
	"hotgo/addons/users/library/contexts"
	"hotgo/internal/consts"
	"hotgo/internal/model/input/assetin"
	basesvc "hotgo/internal/service"
)

var (
	Asset = cAsset{}
)

type cAsset struct{}

// Account 当前用户资产账户
func (c *cAsset) Account(ctx context.Context, req *asset.AccountReq) (res *asset.AccountRes, err error) {
	userId := contexts.GetUsersUserId(ctx)
	if userId <= 0 {
		return nil, gerror.New("用户身份异常，请重新登录")
	}

	account, err := basesvc.Asset().Account(ctx, &assetin.AccountInp{UserId: uint64(userId), Account: consts.AssetAccountUSDT})
	if err != nil {
		return nil, err
	}

	res = new(asset.AccountRes)
	if account != nil {
		res.Available = account.Available
		res.Frozen = account.Frozen
		res.Total = account.Available + account.Frozen
	}
	return
}

// LogList 当前用户资产流水
func (c *cAsset) LogList(ctx context.Context, req *asset.LogListReq) (res *asset.LogListRes, err error) {
	userId := contexts.GetUsersUserId(ctx)
	if userId <= 0 {
		return nil, gerror.New("用户身份异常，请重新登录")
	}

	list, totalCount, err := basesvc.Asset().Logs(ctx, &assetin.LogListInp{
		PageReq:   req.PageReq,
		UserId:    uint64(userId),
		Direction: req.Type,
	})
	if err != nil {
		return nil, err
	}

	res = new(asset.LogListRes)
	res.List = make([]*asset.LogItem, 0, len(list))
	for _, item := range list {
		if item.ChangeAvailable == 0 {
			continue
		}
		date := ""
		if item.CreatedAt != nil {
			date = item.CreatedAt.Format("Y-m-d H:i")
		}
		res.List = append(res.List, &asset.LogItem{
			Type:   req.Type,
			Remark: item.Remark,
			Date:   date,
			Amount: math.Abs(item.ChangeAvailable),
		})
	}
	res.PageRes.Pack(req, totalCount)
	return
}

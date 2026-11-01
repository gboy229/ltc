// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/pledge/model/dao"
	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/addons/pledge/service"
	isc "hotgo/internal/service"
)

type sSysOrder struct{}

func NewSysOrder() *sSysOrder { return &sSysOrder{} }

func init() { service.RegisterSysOrder(NewSysOrder()) }

// List 订单列表
func (s *sSysOrder) List(ctx context.Context, in *sysin.OrderListInp) (list []*sysin.OrderListModel, totalCount int, err error) {
	cols := dao.PledgeOrder.Columns()
	mod := dao.PledgeOrder.Ctx(ctx)
	if in.UserId > 0 {
		mod = mod.Where(cols.UserId, in.UserId)
	}
	if in.WalletAddress != "" {
		ids, resolveErr := isc.UserWallet().GetUserIdsByWallet(ctx, in.WalletAddress)
		if resolveErr != nil {
			return nil, 0, resolveErr
		}
		if len(ids) == 0 {
			return nil, 0, nil
		}
		mod = mod.WhereIn(cols.UserId, ids)
	}
	if in.OrderNo != "" {
		mod = mod.WhereLike(cols.OrderNo, "%"+in.OrderNo+"%")
	}
	if in.Status > 0 {
		mod = mod.Where(cols.Status, in.Status)
	}
	if in.IsEarlyRefund != nil {
		mod = mod.Where(cols.IsEarlyRefund, *in.IsEarlyRefund)
	}

	totalCount, err = mod.Clone().Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, "查询订单总数失败")
	}
	if totalCount == 0 {
		return
	}
	if err = mod.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Scan(&list); err != nil {
		return nil, 0, gerror.Wrap(err, "查询订单失败")
	}

	// 批量补充钱包地址，便于运营在列表中直接识别用户
	if len(list) > 0 {
		userIds := make([]uint64, 0, len(list))
		for _, item := range list {
			userIds = append(userIds, item.UserId)
		}

		wallets, walletErr := isc.UserWallet().GetWalletsByUserIds(ctx, userIds)
		if walletErr != nil {
			return nil, 0, walletErr
		}
		if len(wallets) == 0 {
			g.Log().Warningf(ctx, "pledge order wallet fill empty, userIds:%v", userIds)
		}
		for _, item := range list {
			item.WalletAddress = wallets[item.UserId]
		}
	}
	return
}

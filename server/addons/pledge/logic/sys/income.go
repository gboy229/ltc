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

type sSysIncome struct{}

func NewSysIncome() *sSysIncome {
	return &sSysIncome{}
}

func init() {
	service.RegisterSysIncome(NewSysIncome())
}

// List 收益流水列表
func (s *sSysIncome) List(ctx context.Context, in *sysin.IncomeListInp) (list []*sysin.IncomeListModel, totalCount int, err error) {
	cols := dao.PledgeIncomeLog.Columns()
	mod := dao.PledgeIncomeLog.Ctx(ctx)
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
	if in.OrderId > 0 {
		mod = mod.Where(cols.OrderId, in.OrderId)
	}
	if in.Type != "" {
		mod = mod.Where(cols.Type, in.Type)
	}
	if in.Status > 0 {
		mod = mod.Where(cols.Status, in.Status)
	}

	totalCount, err = mod.Clone().Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, "查询收益流水总数失败")
	}
	if totalCount == 0 {
		return
	}
	if err = mod.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Scan(&list); err != nil {
		return nil, 0, gerror.Wrap(err, "查询收益流水失败")
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
			g.Log().Warningf(ctx, "pledge income wallet fill empty, userIds:%v", userIds)
		}
		for _, item := range list {
			item.WalletAddress = wallets[item.UserId]
		}
	}
	return
}

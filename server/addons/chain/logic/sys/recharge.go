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
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	"hotgo/addons/chain/model/input/sysin"
	"hotgo/addons/chain/service"
)

type sSysRecharge struct{}

func NewSysRecharge() *sSysRecharge {
	return &sSysRecharge{}
}

func init() {
	service.RegisterSysRecharge(NewSysRecharge())
}

// View 获取指定充值记录信息
func (s *sSysRecharge) View(ctx context.Context, in *sysin.RechargeViewInp) (res *sysin.RechargeViewModel, err error) {
	var row *entity.ChainRecharge
	if err = dao.ChainRecharge.Ctx(ctx).WherePri(in.Id).Scan(&row); err != nil {
		return
	}
	if row == nil {
		err = gerror.New("充值记录不存在")
		return
	}

	res = &sysin.RechargeViewModel{ChainRecharge: *row}
	return
}

// List 获取充值记录列表
func (s *sSysRecharge) List(ctx context.Context, in *sysin.RechargeListInp) (list []*sysin.RechargeListModel, totalCount int, err error) {
	mod := dao.ChainRecharge.Ctx(ctx)
	cols := dao.ChainRecharge.Columns()

	// 用户关键词：纯数字按充值 uid 精确匹配，其余按用户钱包/充值地址匹配
	if in.Keyword != "" {
		if gstr.IsNumeric(in.Keyword) {
			mod = mod.Where(cols.Uid, gconv.Int64(in.Keyword))
		} else {
			like := "%" + in.Keyword + "%"
			mod = mod.Where(g.Model().Builder().
				WhereLike(cols.FromAddress, like).
				WhereOrLike(cols.ToAddress, like),
			)
		}
	}

	if in.Type != "" {
		mod = mod.Where(cols.Type, in.Type)
	}

	if in.Uid > 0 {
		mod = mod.Where(cols.Uid, in.Uid)
	}

	if in.Chain != "" {
		mod = mod.WhereLike(cols.Chain, "%"+in.Chain+"%")
	}

	if in.Symbol != "" {
		mod = mod.WhereLike(cols.Symbol, "%"+in.Symbol+"%")
	}

	if in.FromAddress != "" {
		mod = mod.WhereLike(cols.FromAddress, "%"+in.FromAddress+"%")
	}

	if in.ToAddress != "" {
		mod = mod.WhereLike(cols.ToAddress, "%"+in.ToAddress+"%")
	}

	if in.ContractAddress != "" {
		mod = mod.WhereLike(cols.ContractAddress, "%"+in.ContractAddress+"%")
	}

	if in.Hash != "" {
		mod = mod.WhereLike(cols.Hash, "%"+in.Hash+"%")
	}

	if in.AmountMin > 0 {
		mod = mod.WhereGTE(cols.Amount, in.AmountMin)
	}

	if in.AmountMax > 0 {
		mod = mod.WhereLTE(cols.Amount, in.AmountMax)
	}

	if in.BlockHeightMin > 0 {
		mod = mod.WhereGTE(cols.BlockHeight, in.BlockHeightMin)
	}

	if in.BlockHeightMax > 0 {
		mod = mod.WhereLTE(cols.BlockHeight, in.BlockHeightMax)
	}

	if len(in.CreatedAt) == 2 {
		mod = mod.WhereBetween(cols.CreatedAt, in.CreatedAt[0], in.CreatedAt[1])
	}

	totalCount, err = mod.Count()
	if err != nil || totalCount == 0 {
		return
	}

	if err = mod.Page(in.Page, in.PerPage).Order("id desc").Scan(&list); err != nil {
		return
	}
	return
}

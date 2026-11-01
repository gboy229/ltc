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
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"hotgo/addons/chain/consts"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	"hotgo/addons/chain/model/input/sysin"
	"hotgo/addons/chain/service"
	"hotgo/internal/library/contexts"
)

type sSysWithdraw struct{}

func NewSysWithdraw() *sSysWithdraw {
	return &sSysWithdraw{}
}

func init() {
	service.RegisterSysWithdraw(NewSysWithdraw())
}

// Status 提现审核，仅允许审核中的记录变更为审核通过或拒绝
func (s *sSysWithdraw) Status(ctx context.Context, in *sysin.WithdrawStatusInp) (err error) {
	cols := dao.ChainWithdraw.Columns()

	var models *entity.ChainWithdraw
	if err = dao.ChainWithdraw.Ctx(ctx).WherePri(in.Id).Scan(&models); err != nil {
		return
	}
	if models == nil {
		err = gerror.New("提现记录不存在")
		return
	}
	if models.Status != consts.ChainWithdrawStatusPending {
		err = gerror.New("当前提现记录不是审核中状态，不能重复审核")
		return
	}

	// 乐观锁：version 未变化才允许更新，避免并发重复审核
	result, err := dao.ChainWithdraw.Ctx(ctx).
		WherePri(in.Id).
		Where(cols.Version, models.Version).
		Data(g.Map{
			cols.Status:  in.Status,
			cols.AuditBy: contexts.GetUserId(ctx),
			cols.AuditAt: gtime.Now(),
			cols.Remark:  in.Remark,
			cols.Version: models.Version + 1,
		}).Update()
	if err != nil {
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rows == 0 {
		err = gerror.New("提现记录已被其他操作修改，请刷新后重试")
		return
	}
	return
}

// View 获取指定提现记录信息
func (s *sSysWithdraw) View(ctx context.Context, in *sysin.WithdrawViewInp) (res *sysin.WithdrawViewModel, err error) {
	var row *entity.ChainWithdraw
	if err = dao.ChainWithdraw.Ctx(ctx).WherePri(in.Id).Scan(&row); err != nil {
		return
	}
	if row == nil {
		err = gerror.New("提现记录不存在")
		return
	}

	res = &sysin.WithdrawViewModel{ChainWithdraw: *row}
	return
}

// List 获取提现记录列表
func (s *sSysWithdraw) List(ctx context.Context, in *sysin.WithdrawListInp) (list []*sysin.WithdrawListModel, totalCount int, err error) {
	mod := dao.ChainWithdraw.Ctx(ctx)
	cols := dao.ChainWithdraw.Columns()

	// 用户关键词：纯数字按 uid 精确匹配，其余按用户钱包/出币地址匹配
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

	if in.Uid > 0 {
		mod = mod.Where(cols.Uid, in.Uid)
	}

	if in.Chain != "" {
		mod = mod.WhereLike(cols.Chain, "%"+in.Chain+"%")
	}

	if in.Symbol != "" {
		mod = mod.WhereLike(cols.Symbol, "%"+in.Symbol+"%")
	}

	if in.Hash != "" {
		mod = mod.WhereLike(cols.Hash, "%"+in.Hash+"%")
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

	if in.Status != nil {
		mod = mod.Where(cols.Status, *in.Status)
	}

	if in.AmountMin > 0 {
		mod = mod.WhereGTE(cols.PayAmount, in.AmountMin)
	}

	if in.AmountMax > 0 {
		mod = mod.WhereLTE(cols.PayAmount, in.AmountMax)
	}

	if len(in.CreatedAt) == 2 {
		mod = mod.WhereBetween(cols.CreatedAt, in.CreatedAt[0], in.CreatedAt[1])
	}

	if len(in.AuditAt) == 2 {
		mod = mod.WhereBetween(cols.AuditAt, in.AuditAt[0], in.AuditAt[1])
	}

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	if totalCount == 0 {
		return
	}

	if err = mod.Page(in.Page, in.PerPage).Order("id desc").Scan(&list); err != nil {
		return
	}
	return
}

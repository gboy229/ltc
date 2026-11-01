// Package income
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package income

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"hotgo/addons/pledge/api/api/income"
	"hotgo/addons/pledge/model/dao"
	"hotgo/addons/pledge/model/entity"
	"hotgo/addons/pledge/model/input/incomein"
	"hotgo/addons/pledge/service"
	"hotgo/internal/consts"
	"hotgo/internal/model/input/assetin"
	isc "hotgo/internal/service"
)

type sIncome struct{}

func NewIncome() *sIncome { return &sIncome{} }

func init() { service.RegisterIncome(NewIncome()) }

// Create 写入收益流水并资产入账
func (s *sIncome) Create(ctx context.Context, in *incomein.CreateInp) (id uint64, inserted bool, err error) {
	if in.UserId == 0 || in.Amount <= 0 {
		return 0, false, gerror.New("收益用户和金额不能为空")
	}

	err = dao.PledgeIncomeLog.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.PledgeIncomeLog.Columns()
		now := gtime.Now()
		result, insertErr := dao.PledgeIncomeLog.Ctx(ctx).Data(g.Map{
			cols.UserId:     in.UserId,
			cols.OrderId:    in.OrderId,
			cols.FromUserId: in.FromUserId,
			cols.Type:       in.Type,
			cols.Amount:     in.Amount,
			cols.BaseAmount: in.BaseAmount,
			cols.Rate:       in.Rate,
			cols.PeriodNo:   in.PeriodNo,
			cols.Status:     1,
			cols.CreatedAt:  now,
			cols.UpdatedAt:  now,
		}).InsertIgnore()
		if insertErr != nil {
			return gerror.Wrap(insertErr, "写入收益流水失败")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			inserted = false
			return nil
		}

		logId, idErr := result.LastInsertId()
		if idErr != nil {
			return gerror.Wrap(idErr, "获取收益流水ID失败")
		}

		bizType := in.AssetBizType
		if bizType == "" {
			bizType = assetBizType(in.Type)
		}
		if changeErr := isc.Asset().Change(ctx, &assetin.ChangeInp{
			UserId:     in.UserId,
			Account:    consts.AssetAccountUSDT,
			AddonName:  consts.AssetAddonPledge,
			ChangeType: consts.AssetChangeTypeIncr,
			BizType:    bizType,
			BizId:      uint64(logId),
			BizSn:      in.BizSn,
			Num:        in.Amount,
			Remark:     in.Remark,
		}); changeErr != nil {
			return changeErr
		}

		id = uint64(logId)
		inserted = true
		return nil
	})
	if err != nil {
		return 0, false, err
	}
	return id, inserted, nil
}

// Sum 收益金额汇总，date为空时不限制日期
func (s *sIncome) Sum(ctx context.Context, userId uint64, types []string, date string) (amount float64, err error) {
	cols := dao.PledgeIncomeLog.Columns()
	mod := dao.PledgeIncomeLog.Ctx(ctx).Where(cols.UserId, userId)
	if len(types) > 0 {
		mod = mod.WhereIn(cols.Type, types)
	}
	if date != "" {
		mod = mod.Where("DATE("+cols.CreatedAt+") = ?", date)
	}
	if amount, err = mod.Sum(cols.Amount); err != nil {
		return 0, gerror.Wrap(err, "统计收益失败")
	}
	return
}

var incomeLabels = map[string]string{
	"static":    "静态收益",
	"direct":    "直推收益",
	"diff":      "级差收益",
	"lot":       "手数收益",
	"indirect":  "间推收益",
	"team":      "团队收益",
	"principal": "本金返还",
}

// List 当前用户收益流水
func (s *sIncome) List(ctx context.Context, userId uint64, in *income.ListReq) (list []*income.ListItem, totalCount int, err error) {
	cols := dao.PledgeIncomeLog.Columns()
	if in.Type != "" {
		if _, ok := incomeLabels[in.Type]; !ok {
			return nil, 0, gerror.New("收益类型不合法")
		}
	}

	mod := dao.PledgeIncomeLog.Ctx(ctx).Where(cols.UserId, userId)
	if in.Type != "" {
		mod = mod.Where(cols.Type, in.Type)
	}

	if totalCount, err = mod.Clone().Count(); err != nil {
		return nil, 0, gerror.Wrap(err, "查询收益流水总数失败")
	}
	if totalCount == 0 {
		return
	}

	var rows []*entity.PledgeIncomeLog
	if err = mod.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Scan(&rows); err != nil {
		return nil, 0, gerror.Wrap(err, "查询收益流水失败")
	}

	list = make([]*income.ListItem, 0, len(rows))
	for _, row := range rows {
		date := ""
		if row.CreatedAt != nil {
			date = row.CreatedAt.Format("Y-m-d H:i")
		}
		list = append(list, &income.ListItem{
			Type:   row.Type,
			Label:  incomeLabels[row.Type],
			Date:   date,
			Amount: row.Amount,
		})
	}
	return
}

func assetBizType(incomeType string) string {
	switch incomeType {
	case "static":
		return "static_income"
	case "direct":
		return "direct_income"
	case "diff":
		return "diff_income"
	case "lot":
		return "lot_income"
	case "principal":
		return "principal_return"
	default:
		return incomeType
	}
}

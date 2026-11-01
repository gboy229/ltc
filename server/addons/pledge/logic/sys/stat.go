// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"fmt"
	"sort"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"hotgo/addons/pledge/consts"
	"hotgo/addons/pledge/model/dao"
	"hotgo/addons/pledge/model/input/sysin"
	"hotgo/addons/pledge/service"
)

type sSysStat struct{}

func NewSysStat() *sSysStat {
	return &sSysStat{}
}

func init() {
	service.RegisterSysStat(NewSysStat())
}

// Stat 质押订单与收益统计
func (s *sSysStat) Stat(ctx context.Context, in *sysin.StatInp) (res *sysin.StatModel, err error) {
	res = new(sysin.StatModel)
	// 未选择时间范围时统计全部
	start, end := in.StartAt, in.EndAt

	if res.Order, err = s.orderStat(ctx, start, end); err != nil {
		return
	}
	if res.Income, err = s.incomeStat(ctx, start, end); err != nil {
		return
	}
	if res.Trend, err = s.trend(ctx, start, end); err != nil {
		return
	}
	return
}

// orderStat 订单统计
func (s *sSysStat) orderStat(ctx context.Context, start, end *gtime.Time) (res *sysin.OrderStat, err error) {
	cols := dao.PledgeOrder.Columns()
	mod := dao.PledgeOrder.Ctx(ctx)
	if start != nil && end != nil {
		mod = mod.WhereBetween(cols.CreatedAt, start, end)
	}

	res = new(sysin.OrderStat)

	var total struct {
		Count  int     `orm:"cnt"`
		Amount float64 `orm:"amount"`
	}
	if err = mod.Fields(fmt.Sprintf(
		"COUNT(*) AS cnt, COALESCE(SUM(%s),0) AS amount",
		cols.Amount,
	)).Scan(&total); err != nil {
		return nil, err
	}
	res.Count = total.Count
	res.Amount = total.Amount

	if res.ActiveCount, err = mod.Where(cols.Status, consts.PledgeOrderStatusActive).Count(); err != nil {
		return nil, err
	}
	if res.FinishedCount, err = mod.Where(cols.Status, consts.PledgeOrderStatusFinished).Count(); err != nil {
		return nil, err
	}
	return
}

// incomeStat 收益统计，已撤销的收益不计入产出
func (s *sSysStat) incomeStat(ctx context.Context, start, end *gtime.Time) (res *sysin.IncomeStat, err error) {
	cols := dao.PledgeIncomeLog.Columns()
	mod := dao.PledgeIncomeLog.Ctx(ctx).Where(cols.Status, consts.PledgeIncomeStatusCredited)
	if start != nil && end != nil {
		mod = mod.WhereBetween(cols.CreatedAt, start, end)
	}

	res = new(sysin.IncomeStat)

	var total struct {
		Count  int     `orm:"cnt"`
		Amount float64 `orm:"amount"`
	}
	if err = mod.Fields(fmt.Sprintf(
		"COUNT(*) AS cnt, COALESCE(SUM(%s),0) AS amount",
		cols.Amount,
	)).Scan(&total); err != nil {
		return nil, err
	}
	res.Count = total.Count
	res.Amount = total.Amount

	// 按类型分组
	var rows []struct {
		Type   string  `orm:"type"`
		Count  int     `orm:"cnt"`
		Amount float64 `orm:"amount"`
	}
	if err = mod.Fields(fmt.Sprintf(
		"%s AS type, COUNT(*) AS cnt, COALESCE(SUM(%s),0) AS amount",
		cols.Type, cols.Amount,
	)).Group(cols.Type).Scan(&rows); err != nil {
		return nil, err
	}

	typeMap := make(map[string]*sysin.IncomeTypeItem, len(rows))
	for _, r := range rows {
		typeMap[r.Type] = &sysin.IncomeTypeItem{Type: r.Type, Count: r.Count, Amount: r.Amount}
	}

	// 按字典顺序补全所有类型，没有数据的补 0；间推收益、级差收益不参与统计
	res.ByType = make([]*sysin.IncomeTypeItem, 0, len(consts.PledgeIncomeTypeOptions))
	for _, opt := range consts.PledgeIncomeTypeOptions {
		key := gconv.String(opt.Value)
		if key == consts.PledgeIncomeTypeIndirect || key == consts.PledgeIncomeTypeDiff {
			continue
		}
		item := typeMap[key]
		if item == nil {
			item = &sysin.IncomeTypeItem{Type: key}
		}
		item.Label = opt.Label
		res.ByType = append(res.ByType, item)
	}
	return
}

// trend 按天趋势：订单金额与产出金额
func (s *sSysStat) trend(ctx context.Context, start, end *gtime.Time) (list []*sysin.TrendItem, err error) {
	orderCols := dao.PledgeOrder.Columns()
	var orderDays []struct {
		Date   string  `orm:"date"`
		Amount float64 `orm:"amount"`
	}
	orderMod := dao.PledgeOrder.Ctx(ctx)
	if start != nil && end != nil {
		orderMod = orderMod.WhereBetween(orderCols.CreatedAt, start, end)
	}
	if err = orderMod.Fields(fmt.Sprintf(
		"DATE(%s) AS date, COALESCE(SUM(%s),0) AS amount",
		orderCols.CreatedAt, orderCols.Amount,
	)).Group(fmt.Sprintf("DATE(%s)", orderCols.CreatedAt)).Order("date ASC").Scan(&orderDays); err != nil {
		return nil, err
	}

	incomeCols := dao.PledgeIncomeLog.Columns()
	var incomeDays []struct {
		Date   string  `orm:"date"`
		Amount float64 `orm:"amount"`
	}
	incomeMod := dao.PledgeIncomeLog.Ctx(ctx).Where(incomeCols.Status, consts.PledgeIncomeStatusCredited)
	if start != nil && end != nil {
		incomeMod = incomeMod.WhereBetween(incomeCols.CreatedAt, start, end)
	}
	if err = incomeMod.Fields(fmt.Sprintf(
		"DATE(%s) AS date, COALESCE(SUM(%s),0) AS amount",
		incomeCols.CreatedAt, incomeCols.Amount,
	)).Group(fmt.Sprintf("DATE(%s)", incomeCols.CreatedAt)).Order("date ASC").Scan(&incomeDays); err != nil {
		return nil, err
	}

	orderMap := make(map[string]float64, len(orderDays))
	incomeMap := make(map[string]float64, len(incomeDays))
	allDates := make(map[string]struct{}, len(orderDays)+len(incomeDays))
	for _, d := range orderDays {
		orderMap[d.Date] = d.Amount
		allDates[d.Date] = struct{}{}
	}
	for _, d := range incomeDays {
		incomeMap[d.Date] = d.Amount
		allDates[d.Date] = struct{}{}
	}

	dates := make([]string, 0, len(allDates))
	for d := range allDates {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	list = make([]*sysin.TrendItem, 0, len(dates))
	for _, d := range dates {
		list = append(list, &sysin.TrendItem{
			Date:         d,
			OrderAmount:  orderMap[d],
			IncomeAmount: incomeMap[d],
		})
	}
	return
}

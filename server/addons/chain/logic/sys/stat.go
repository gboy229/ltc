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

	"hotgo/addons/chain/consts"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/input/sysin"
	"hotgo/addons/chain/service"
)

type sSysStat struct{}

func NewSysStat() *sSysStat {
	return &sSysStat{}
}

func init() {
	service.RegisterSysStat(NewSysStat())
}

// Stat 充值提现统计
func (s *sSysStat) Stat(ctx context.Context, in *sysin.StatInp) (res *sysin.StatModel, err error) {
	res = new(sysin.StatModel)
	start, end := s.normalizeRange(in)

	if res.Recharge, err = s.rechargeStat(ctx, in, start, end); err != nil {
		return
	}
	if res.Withdraw, err = s.withdrawStat(ctx, in, start, end); err != nil {
		return
	}
	if res.Trend, err = s.trend(ctx, in, start, end); err != nil {
		return
	}
	return
}

// trend 按天趋势：充值金额与提现金额
func (s *sSysStat) trend(ctx context.Context, in *sysin.StatInp, start, end *gtime.Time) (list []*sysin.TrendItem, err error) {
	rechargeCols := dao.ChainRecharge.Columns()
	var rechargeDays []struct {
		Date   string  `orm:"date"`
		Amount float64 `orm:"amount"`
	}
	rechargeMod := dao.ChainRecharge.Ctx(ctx).Where(rechargeCols.Type, consts.ChainRechargeTypeOrder)
	if start != nil && end != nil {
		rechargeMod = rechargeMod.WhereBetween(rechargeCols.CreatedAt, start, end)
	}
	if in.Chain != "" {
		rechargeMod = rechargeMod.WhereLike(rechargeCols.Chain, "%"+in.Chain+"%")
	}
	if in.Symbol != "" {
		rechargeMod = rechargeMod.WhereLike(rechargeCols.Symbol, "%"+in.Symbol+"%")
	}
	if err = rechargeMod.Fields(fmt.Sprintf(
		"DATE(%s) AS date, COALESCE(SUM(%s),0) AS amount",
		rechargeCols.CreatedAt, rechargeCols.Amount,
	)).Group(fmt.Sprintf("DATE(%s)", rechargeCols.CreatedAt)).Order("date ASC").Scan(&rechargeDays); err != nil {
		return nil, err
	}

	withdrawCols := dao.ChainWithdraw.Columns()
	var withdrawDays []struct {
		Date   string  `orm:"date"`
		Amount float64 `orm:"amount"`
	}
	withdrawMod := dao.ChainWithdraw.Ctx(ctx)
	if start != nil && end != nil {
		withdrawMod = withdrawMod.WhereBetween(withdrawCols.CreatedAt, start, end)
	}
	if in.Chain != "" {
		withdrawMod = withdrawMod.WhereLike(withdrawCols.Chain, "%"+in.Chain+"%")
	}
	if in.Symbol != "" {
		withdrawMod = withdrawMod.WhereLike(withdrawCols.Symbol, "%"+in.Symbol+"%")
	}
	if err = withdrawMod.Fields(fmt.Sprintf(
		"DATE(%s) AS date, COALESCE(SUM(%s),0) AS amount",
		withdrawCols.CreatedAt, withdrawCols.PayAmount,
	)).Group(fmt.Sprintf("DATE(%s)", withdrawCols.CreatedAt)).Order("date ASC").Scan(&withdrawDays); err != nil {
		return nil, err
	}

	rechargeMap := make(map[string]float64, len(rechargeDays))
	withdrawMap := make(map[string]float64, len(withdrawDays))
	allDates := make(map[string]struct{}, len(rechargeDays)+len(withdrawDays))
	for _, d := range rechargeDays {
		rechargeMap[d.Date] = d.Amount
		allDates[d.Date] = struct{}{}
	}
	for _, d := range withdrawDays {
		withdrawMap[d.Date] = d.Amount
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
			Date:           d,
			RechargeAmount: rechargeMap[d],
			WithdrawAmount: withdrawMap[d],
		})
	}
	return
}

// normalizeRange 时间范围，为空时默认近7天
func (s *sSysStat) normalizeRange(in *sysin.StatInp) (start, end *gtime.Time) {
	if in.StartAt != nil && in.EndAt != nil {
		return in.StartAt, in.EndAt
	}
	now := gtime.Now()
	return now.AddDate(0, 0, -6).StartOfDay(), now
}

// rechargeStat 充值统计：只统计订单类型，排除 Gas
func (s *sSysStat) rechargeStat(ctx context.Context, in *sysin.StatInp, start, end *gtime.Time) (res *sysin.RechargeStat, err error) {
	cols := dao.ChainRecharge.Columns()
	mod := dao.ChainRecharge.Ctx(ctx).Where(cols.Type, consts.ChainRechargeTypeOrder)
	if start != nil && end != nil {
		mod = mod.WhereBetween(cols.CreatedAt, start, end)
	}
	if in.Chain != "" {
		mod = mod.WhereLike(cols.Chain, "%"+in.Chain+"%")
	}
	if in.Symbol != "" {
		mod = mod.WhereLike(cols.Symbol, "%"+in.Symbol+"%")
	}

	res = new(sysin.RechargeStat)

	// 汇总：笔数/总额
	var total struct {
		TotalCount  int     `orm:"total_count"`
		TotalAmount float64 `orm:"total_amount"`
	}
	if err = mod.Fields(fmt.Sprintf("COUNT(*) AS total_count, SUM(%s) AS total_amount", cols.Amount)).Scan(&total); err != nil {
		return
	}
	res.TotalCount = total.TotalCount
	res.TotalAmount = total.TotalAmount

	// 按天：笔数/金额
	var days []struct {
		Date   string  `orm:"date"`
		Count  int     `orm:"cnt"`
		Amount float64 `orm:"amount"`
	}
	if err = mod.Fields(fmt.Sprintf("DATE(%s) AS date, COUNT(*) AS cnt, SUM(%s) AS amount", cols.CreatedAt, cols.Amount)).
		Group(fmt.Sprintf("DATE(%s)", cols.CreatedAt)).
		Order("date ASC").
		Scan(&days); err != nil {
		return
	}
	res.List = make([]*sysin.StatDayItem, 0, len(days))
	for _, d := range days {
		res.List = append(res.List, &sysin.StatDayItem{
			Date:   d.Date,
			Count:  d.Count,
			Amount: d.Amount,
		})
	}
	return
}

// withdrawStat 提现统计：按创建时间统计，不分状态
func (s *sSysStat) withdrawStat(ctx context.Context, in *sysin.StatInp, start, end *gtime.Time) (res *sysin.WithdrawStat, err error) {
	cols := dao.ChainWithdraw.Columns()
	mod := dao.ChainWithdraw.Ctx(ctx)
	if start != nil && end != nil {
		mod = mod.WhereBetween(cols.CreatedAt, start, end)
	}
	if in.Chain != "" {
		mod = mod.WhereLike(cols.Chain, "%"+in.Chain+"%")
	}
	if in.Symbol != "" {
		mod = mod.WhereLike(cols.Symbol, "%"+in.Symbol+"%")
	}

	res = new(sysin.WithdrawStat)

	// 汇总：笔数/申请额/到账额/手续费
	var total struct {
		TotalCount    int     `orm:"total_count"`
		PayAmount     float64 `orm:"pay_amount"`
		ReceiveAmount float64 `orm:"receive_amount"`
		Fee           float64 `orm:"fee"`
	}
	if err = mod.Fields(fmt.Sprintf("COUNT(*) AS total_count, SUM(%s) AS pay_amount, SUM(%s) AS receive_amount, SUM(%s) AS fee", cols.PayAmount, cols.ReceiveAmount, cols.Fee)).Scan(&total); err != nil {
		return
	}
	res.TotalCount = total.TotalCount
	res.PayAmount = total.PayAmount
	res.ReceiveAmount = total.ReceiveAmount
	res.Fee = total.Fee

	// 按天：笔数/申请额/到账额/手续费
	var days []struct {
		Date          string  `orm:"date"`
		Count         int     `orm:"cnt"`
		PayAmount     float64 `orm:"pay_amount"`
		ReceiveAmount float64 `orm:"receive_amount"`
		Fee           float64 `orm:"fee"`
	}
	if err = mod.Fields(fmt.Sprintf("DATE(%s) AS date, COUNT(*) AS cnt, SUM(%s) AS pay_amount, SUM(%s) AS receive_amount, SUM(%s) AS fee", cols.CreatedAt, cols.PayAmount, cols.ReceiveAmount, cols.Fee)).
		Group(fmt.Sprintf("DATE(%s)", cols.CreatedAt)).
		Order("date ASC").
		Scan(&days); err != nil {
		return
	}
	res.List = make([]*sysin.WithdrawDayItem, 0, len(days))
	for _, d := range days {
		res.List = append(res.List, &sysin.WithdrawDayItem{
			Date:          d.Date,
			Count:         d.Count,
			PayAmount:     d.PayAmount,
			ReceiveAmount: d.ReceiveAmount,
			Fee:           d.Fee,
		})
	}
	return
}

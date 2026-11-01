// Package settlement
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package settlement

import (
	"context"
	"math"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	pledgeconsts "hotgo/addons/pledge/consts"
	perf "hotgo/addons/pledge/logic/performance"
	"hotgo/addons/pledge/model/dao"
	"hotgo/addons/pledge/model/entity"
	"hotgo/addons/pledge/model/input/incomein"
	"hotgo/addons/pledge/service"
	"hotgo/internal/consts"
)

const settlementBatchSize = 50

type sSettlement struct{}

func NewSettlement() *sSettlement {
	return &sSettlement{}
}

func init() {
	service.RegisterSettlement(NewSettlement())
}

// Run 执行静态收益结算
func (s *sSettlement) Run(ctx context.Context) (err error) {
	configs, err := s.periodConfigs(ctx)
	if err != nil {
		return err
	}

	todayStart := gtime.Now().StartOfDay()
	cols := dao.PledgeOrder.Columns()
	var orders []*entity.PledgeOrder
	if err = dao.PledgeOrder.Ctx(ctx).
		Where(cols.Status, pledgeconsts.PledgeOrderStatusActive).
		WhereLT(cols.SettleAt, todayStart).
		Limit(settlementBatchSize).
		Scan(&orders); err != nil {
		return gerror.Wrap(err, "查询待结算订单失败")
	}

	for _, order := range orders {
		finished, e := s.settleOrder(ctx, order.Id, configs)
		if e != nil {
			g.Log().Warningf(ctx, "pledge settlement order failed, order:%s err:%+v", order.OrderNo, e)
			continue
		}
		if finished {
			if pushErr := perf.Push(ctx, order.UserId); pushErr != nil {
				g.Log().Warningf(ctx, "push performance event failed, uid:%d err:%+v", order.UserId, pushErr)
			}
		}
	}
	return
}

// CreatePushReward 用真实订单生成直推奖励
func (s *sSettlement) CreatePushReward(ctx context.Context, orderId uint64, staticAmount float64, periodNo int) error {
	var order *entity.PledgeOrder
	if err := dao.PledgeOrder.Ctx(ctx).WherePri(orderId).Scan(&order); err != nil {
		return gerror.Wrap(err, "查询订单失败")
	}
	if order == nil {
		return gerror.New("订单不存在")
	}
	if periodNo <= 0 {
		periodNo = order.IncIncomeNum + 1
	}
	return s.createPushReward(ctx, order, staticAmount, periodNo)
}

// SettleOrder 结算指定订单
func (s *sSettlement) SettleOrder(ctx context.Context, orderId uint64) error {
	configs, err := s.periodConfigs(ctx)
	if err != nil {
		return err
	}
	finished, err := s.settleOrder(ctx, orderId, configs)
	if err != nil {
		return err
	}
	if finished {
		var order *entity.PledgeOrder
		if queryErr := dao.PledgeOrder.Ctx(ctx).WherePri(orderId).Scan(&order); queryErr == nil && order != nil {
			if pushErr := perf.Push(ctx, order.UserId); pushErr != nil {
				g.Log().Warningf(ctx, "push performance event failed, uid:%d err:%+v", order.UserId, pushErr)
			}
		}
	}
	return nil
}

func (s *sSettlement) periodConfigs(ctx context.Context) (map[int]*entity.PledgePeriodConfig, error) {
	cols := dao.PledgePeriodConfig.Columns()
	var rows []*entity.PledgePeriodConfig
	if err := dao.PledgePeriodConfig.Ctx(ctx).
		Where(cols.Status, consts.StatusEnabled).
		Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "查询周期配置失败")
	}

	configs := make(map[int]*entity.PledgePeriodConfig, len(rows))
	for _, row := range rows {
		configs[row.PeriodMonths] = row
	}
	return configs, nil
}

func (s *sSettlement) settleOrder(ctx context.Context, orderId uint64, configs map[int]*entity.PledgePeriodConfig) (finished bool, err error) {
	err = dao.PledgeOrder.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.PledgeOrder.Columns()
		var order *entity.PledgeOrder
		if err := dao.PledgeOrder.Ctx(ctx).
			WherePri(orderId).
			LockUpdate().
			Scan(&order); err != nil {
			return gerror.Wrap(err, "锁定订单失败")
		}
		if order == nil || order.Status != pledgeconsts.PledgeOrderStatusActive {
			return nil
		}

		remaining := order.CycleNum - order.IncIncomeNum
		if remaining <= 0 {
			return nil
		}

		lastDate := order.SettleAt
		if lastDate == nil {
			lastDate = order.CreatedAt
		}
		if lastDate == nil {
			return nil
		}
		dueDays := dateDiffDays(lastDate, gtime.Now())
		if dueDays <= 0 {
			return nil
		}
		if dueDays > remaining {
			dueDays = remaining
		}

		cfg, ok := configs[order.PeriodMonths]
		if !ok || cfg.MonthlyRewardRate <= 0 {
			return gerror.Newf("订单周期配置不存在，period_months:%d", order.PeriodMonths)
		}
		rate := cfg.MonthlyRewardRate

		daily := round4(order.Amount * rate / 100 / 30)
		amount := round4(daily * float64(dueDays))
		if amount <= 0 {
			return nil
		}

		now := gtime.Now()
		_, inserted, err := service.Income().Create(ctx, &incomein.CreateInp{
			UserId:     order.UserId,
			OrderId:    order.Id,
			Type:       "static",
			Amount:     amount,
			BaseAmount: order.Amount,
			Rate:       rate,
			PeriodNo:   order.IncIncomeNum + 1,
			BizSn:      order.OrderNo,
			Remark:     "静态收益",
		})
		if err != nil {
			return err
		}
		if !inserted {
			return nil
		}

		if err = s.createPushReward(ctx, order, amount, order.IncIncomeNum+1); err != nil {
			return err
		}
		if err = s.createDiffReward(ctx, order, amount, order.IncIncomeNum+1); err != nil {
			return err
		}

		lotAmount, lotToday, err := s.createLotReward(ctx, order, dueDays, order.IncIncomeNum+1)
		if err != nil {
			return err
		}

		newIncNum := order.IncIncomeNum + dueDays
		completed := newIncNum >= order.CycleNum

		data := g.Map{
			cols.TodayIncome:  round4(daily + lotToday),
			cols.IncIncome:    round4(order.IncIncome + amount + lotAmount),
			cols.IncIncomeNum: newIncNum,
			cols.SettleAt:     now,
		}
		if completed {
			data[cols.Status] = pledgeconsts.PledgeOrderStatusFinished
		}
		result, err := dao.PledgeOrder.Ctx(ctx).
			WherePri(order.Id).
			Where(cols.Status, pledgeconsts.PledgeOrderStatusActive).
			Data(data).
			Update()
		if err != nil {
			return gerror.Wrap(err, "更新订单结算状态失败")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.New("订单状态已变更，跳过本次结算")
		}

		if completed {
			finished = true
			_, _, err := service.Income().Create(ctx, &incomein.CreateInp{
				UserId:     order.UserId,
				OrderId:    order.Id,
				Type:       "principal",
				Amount:     order.Amount,
				BaseAmount: order.Amount,
				Rate:       0,
				PeriodNo:   newIncNum,
				BizSn:      order.OrderNo,
				Remark:     "本金返还",
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return finished, err
}

func dateDiffDays(from, to *gtime.Time) int {
	fromDate := time.Date(from.Year(), time.Month(from.Month()), from.Day(), 0, 0, 0, 0, time.Local)
	toDate := time.Date(to.Year(), time.Month(to.Month()), to.Day(), 0, 0, 0, 0, time.Local)
	return int(toDate.Sub(fromDate).Hours() / 24)
}

func round4(value float64) float64 {
	return math.Round(value*1e4) / 1e4
}

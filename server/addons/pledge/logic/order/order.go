// Package order
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package order

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"

	apiOrder "hotgo/addons/pledge/api/api/order"
	pledgeconsts "hotgo/addons/pledge/consts"
	perf "hotgo/addons/pledge/logic/performance"
	"hotgo/addons/pledge/model/dao"
	"hotgo/addons/pledge/model/entity"
	"hotgo/addons/pledge/model/input/incomein"
	"hotgo/addons/pledge/service"
	"hotgo/internal/consts"
	"hotgo/internal/library/ratelimit"
	mainmodel "hotgo/internal/model"
	isc "hotgo/internal/service"
)

type sOrder struct{}

func NewOrder() *sOrder {
	return &sOrder{}
}

func init() {
	service.RegisterOrder(NewOrder())
}

// Buy 提交报单
func (s *sOrder) Buy(ctx context.Context, userId int64, in *apiOrder.BuyReq) (res *apiOrder.BuyRes, err error) {
	planName := strings.TrimSpace(in.PlanName)
	var plan *entity.PledgePlan
	if err = dao.PledgePlan.Ctx(ctx).
		Where(dao.PledgePlan.Columns().Name, planName).
		Where(dao.PledgePlan.Columns().PeriodMonths, in.PeriodMonths).
		Where(dao.PledgePlan.Columns().Status, consts.StatusEnabled).
		Scan(&plan); err != nil {
		return nil, gerror.Wrap(err, "查询投资套餐失败")
	}
	if plan == nil {
		return nil, gerror.New("该投资方案不存在或已下架")
	}

	config, err := isc.SysAddonsConfig().GetCurrentAddonConfig(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, "获取插件配置失败")
	}
	toAddress := strings.TrimSpace(config["depositAddress"])
	if toAddress == "" {
		return nil, gerror.New("未配置入金地址，请联系客服")
	}

	if !gconv.Bool(config["orderEnabled"]) {
		return nil, gerror.New("暂未开启报单，请稍后再试")
	}

	res = new(apiOrder.BuyRes)
	res.ToAddress = toAddress
	err = dao.PledgeAddOrder.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		addCols := dao.PledgeAddOrder.Columns()

		// 一个用户同时只能有一条待处理报单：新报单提交时，旧待处理报单自动取消
		var pending *entity.PledgeAddOrder
		if err = dao.PledgeAddOrder.Ctx(ctx).
			Where(addCols.UserId, userId).
			Where(addCols.Status, 0).
			LockUpdate().
			Scan(&pending); err != nil {
			return gerror.Wrap(err, "查询待处理报单失败")
		}
		if pending != nil {
			if _, err = dao.PledgeAddOrder.Ctx(ctx).
				WherePri(pending.Id).
				Where(addCols.Status, 0).
				Data(addCols.Status, 2).
				Update(); err != nil {
				return gerror.Wrap(err, "取消旧报单失败")
			}
		}

		orderNo, err := genOrderNo(ctx)
		if err != nil {
			return err
		}

		if _, err = dao.PledgeAddOrder.Ctx(ctx).Data(g.Map{
			addCols.OrderNo: orderNo,
			addCols.UserId:  userId,
			addCols.PlanId:  plan.Id,
			addCols.Amount:  plan.Amount,
			addCols.Status:  0,
		}).Insert(); err != nil {
			return gerror.Wrap(err, "创建报单失败")
		}

		res.OrderNo = orderNo
		res.Amount = plan.Amount
		res.Status = 0
		return nil
	})
	return
}

// List 报单列表
func (s *sOrder) List(ctx context.Context, userId int64, in *apiOrder.ListReq) (list []*apiOrder.ListItem, totalCount int, err error) {
	cols := dao.PledgeOrder.Columns()
	mod := dao.PledgeOrder.Ctx(ctx).Where(cols.UserId, userId)

	switch strings.TrimSpace(in.Status) {
	case "active":
		mod = mod.
			Where(cols.Status, 1).
			Where("(" + cols.CycleNum + " = 0 OR " + cols.IncIncomeNum + " < " + cols.CycleNum + ")")
	case "completed":
		mod = mod.Where("(" + cols.Status + " != 1 OR (" + cols.CycleNum + " > 0 AND " + cols.IncIncomeNum + " >= " + cols.CycleNum + "))")
	}

	totalCount, err = mod.Clone().Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, "查询报单总数失败")
	}
	if totalCount == 0 {
		return
	}

	var orders []*entity.PledgeOrder
	if err = mod.
		Page(in.Page, in.PerPage).
		OrderDesc(cols.Id).
		Scan(&orders); err != nil {
		return nil, 0, gerror.Wrap(err, "查询报单列表失败")
	}

	list = make([]*apiOrder.ListItem, 0, len(orders))
	for _, order := range orders {
		investDate := ""
		if order.CreatedAt != nil {
			investDate = order.CreatedAt.Format("Y-m-d")
		}

		status := "active"
		if order.Status != 1 || (order.CycleNum > 0 && order.IncIncomeNum >= order.CycleNum) {
			status = "completed"
		}

		list = append(list, &apiOrder.ListItem{
			Id:                   order.OrderNo,
			Amount:               order.Amount,
			TodayRelease:         order.TodayIncome,
			TotalRelease:         order.IncIncome,
			Period:               order.CycleNum / 30,
			InvestDate:           investDate,
			Status:               status,
			EarlyExitPenaltyRate: order.EarlyExitPenaltyRate,
		})
	}
	return
}

// Statistics 报单统计
func (s *sOrder) Statistics(ctx context.Context, userId int64, in *apiOrder.StatisticsReq) (res *apiOrder.StatisticsRes, err error) {
	orderCols := dao.PledgeOrder.Columns()
	res = new(apiOrder.StatisticsRes)

	if res.Total, err = dao.PledgeOrder.Ctx(ctx).
		Where(orderCols.UserId, userId).
		Sum(orderCols.Amount); err != nil {
		return nil, gerror.Wrap(err, "统计入金总额失败")
	}

	uid := uint64(userId)
	today := gtime.Now().Format("Y-m-d")
	earningTypes := []string{"static", "direct"}
	if res.TodayRelease, err = service.Income().Sum(ctx, uid, earningTypes, today); err != nil {
		return nil, err
	}
	if res.TotalRelease, err = service.Income().Sum(ctx, uid, earningTypes, ""); err != nil {
		return nil, err
	}
	if res.StaticIncome, err = service.Income().Sum(ctx, uid, []string{"static"}, ""); err != nil {
		return nil, err
	}
	if res.ReferralIncome, err = service.Income().Sum(ctx, uid, []string{"direct"}, ""); err != nil {
		return nil, err
	}
	if res.EarningsTotal, err = service.Income().Sum(ctx, uid, []string{}, ""); err != nil {
		return nil, err
	}
	return
}

// Refund 提前退本
func (s *sOrder) Refund(ctx context.Context, userId int64, in *apiOrder.RefundReq) (res *apiOrder.RefundRes, err error) {
	if err = isc.UserWallet().VerifySignature(ctx, userId, in.Nonce, in.Signature); err != nil {
		return nil, err
	}

	allowed, err := ratelimit.Allow(ctx, fmt.Sprintf("rate:refund:apply:%d", userId), 10*time.Second)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, gerror.New("操作过于频繁，请稍后再试")
	}

	cols := dao.PledgeOrder.Columns()
	var order *entity.PledgeOrder
	if err = dao.PledgeOrder.Ctx(ctx).
		Where(cols.UserId, userId).
		Where(cols.OrderNo, strings.TrimSpace(in.OrderNo)).
		Scan(&order); err != nil {
		return nil, gerror.Wrap(err, "查询订单失败")
	}
	if order == nil {
		return nil, gerror.New("订单不存在")
	}
	if order.Status != pledgeconsts.PledgeOrderStatusActive || order.IsEarlyRefund == pledgeconsts.PledgeOrderIsEarlyRefundYes {
		return nil, gerror.New("当前订单不能提前退本")
	}

	refund := round4(order.Amount*(100-order.EarlyExitPenaltyRate)/100 - order.IncIncome)
	if refund <= 0 {
		return nil, gerror.New("当前订单可退金额为0")
	}

	err = dao.PledgeOrder.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		result, updateErr := dao.PledgeOrder.Ctx(ctx).
			WherePri(order.Id).
			Where(cols.Status, pledgeconsts.PledgeOrderStatusActive).
			Where(cols.IsEarlyRefund, pledgeconsts.PledgeOrderIsEarlyRefundNo).
			Data(g.Map{
				cols.Status:        pledgeconsts.PledgeOrderStatusFinished,
				cols.IsEarlyRefund: pledgeconsts.PledgeOrderIsEarlyRefundYes,
				cols.SettleAt:      gtime.Now(),
			}).
			Update()
		if updateErr != nil {
			return gerror.Wrap(updateErr, "更新订单状态失败")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.New("订单状态已变更，请刷新后重试")
		}

		_, inserted, createErr := service.Income().Create(ctx, &incomein.CreateInp{
			UserId:       uint64(userId),
			OrderId:      order.Id,
			Type:         "principal",
			Amount:       refund,
			BaseAmount:   order.Amount,
			Rate:         0,
			PeriodNo:     order.IncIncomeNum + 1,
			BizSn:        order.OrderNo,
			Remark:       "提前退本",
			AssetBizType: consts.AssetBizTypePledgeEarlyRefund,
		})
		if createErr != nil {
			return createErr
		}
		if !inserted {
			return gerror.New("退款已处理")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if pushErr := perf.Push(ctx, uint64(userId)); pushErr != nil {
		g.Log().Warningf(ctx, "push performance event failed, uid:%d err:%+v", userId, pushErr)
	}
	return &apiOrder.RefundRes{Amount: refund}, nil
}

// Status 查询报单状态
func (s *sOrder) Status(ctx context.Context, userId int64, in *apiOrder.StatusReq) (res *apiOrder.StatusRes, err error) {
	addCols := dao.PledgeAddOrder.Columns()
	var addOrder *entity.PledgeAddOrder
	if err = dao.PledgeAddOrder.Ctx(ctx).
		Where(addCols.UserId, userId).
		Where(addCols.OrderNo, strings.TrimSpace(in.OrderNo)).
		Scan(&addOrder); err != nil {
		return nil, gerror.Wrap(err, "查询报单失败")
	}
	if addOrder == nil {
		return nil, gerror.New("报单不存在")
	}

	res = &apiOrder.StatusRes{
		OrderNo: addOrder.OrderNo,
		Amount:  addOrder.Amount,
		Status:  addOrder.Status,
	}

	if addOrder.Status == 1 {
		var order *entity.PledgeOrder
		if err = dao.PledgeOrder.Ctx(ctx).
			Where(dao.PledgeOrder.Columns().OrderNo, addOrder.OrderNo).
			Scan(&order); err != nil {
			return nil, gerror.Wrap(err, "查询正式订单失败")
		}
		if order != nil {
			res.PledgeOrderStatus = order.Status
		}
	}
	return
}

// HandlePayment 处理链上支付事件
func (s *sOrder) HandlePayment(ctx context.Context, event *mainmodel.ChainPaymentEvent) (err error) {
	if event == nil || event.Uid <= 0 || event.RechargeId == 0 {
		return nil
	}

	orderCols := dao.PledgeOrder.Columns()
	exist, err := dao.PledgeOrder.Ctx(ctx).
		Where(orderCols.RechargeId, event.RechargeId).
		Exist()
	if err != nil {
		return gerror.Wrap(err, "查询订单失败")
	}
	if exist {
		return nil
	}

	var created bool
	var createdUserId uint64
	err = dao.PledgeAddOrder.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		addCols := dao.PledgeAddOrder.Columns()

		var addOrder *entity.PledgeAddOrder
		if err = dao.PledgeAddOrder.Ctx(ctx).
			Where(addCols.UserId, event.Uid).
			Where(addCols.Amount, event.Amount).
			Where(addCols.Status, 0).
			OrderAsc(addCols.Id).
			LockUpdate().
			Scan(&addOrder); err != nil {
			return gerror.Wrap(err, "查询待处理报单失败")
		}
		if addOrder == nil {
			return nil
		}

		var plan *entity.PledgePlan
		if err = dao.PledgePlan.Ctx(ctx).WherePri(addOrder.PlanId).Scan(&plan); err != nil {
			return gerror.Wrap(err, "查询套餐失败")
		}
		if plan == nil {
			return gerror.New("报单关联的套餐不存在")
		}

		result, err := dao.PledgeOrder.Ctx(ctx).Data(g.Map{
			orderCols.OrderNo:              addOrder.OrderNo,
			orderCols.UserId:               addOrder.UserId,
			orderCols.PlanId:               addOrder.PlanId,
			orderCols.Amount:               addOrder.Amount,
			orderCols.PeriodMonths:         plan.PeriodMonths,
			orderCols.CycleNum:             plan.PeriodMonths * 30,
			orderCols.EarlyExitPenaltyRate: plan.EarlyExitPenaltyRate,
			orderCols.IncIncomeNum:         0,
			orderCols.IncIncome:            0,
			orderCols.TodayIncome:          0,
			orderCols.RechargeId:           event.RechargeId,
			orderCols.Status:               1,
			orderCols.SettleAt:             gtime.Now(),
		}).InsertIgnore()
		if err != nil {
			return gerror.Wrap(err, "创建正式订单失败")
		}
		inserted, _ := result.RowsAffected()
		if inserted == 0 {
			return nil
		}

		_, _ = dao.PledgeAddOrder.Ctx(ctx).
			WherePri(addOrder.Id).
			Where(addCols.Status, 0).
			Data(addCols.Status, 1).
			Update()

		created = true
		createdUserId = addOrder.UserId
		return nil
	})
	if err != nil {
		return err
	}
	if created {
		if pushErr := perf.Push(ctx, createdUserId); pushErr != nil {
			g.Log().Warningf(ctx, "push performance event failed, uid:%d err:%+v", createdUserId, pushErr)
		}
	}
	return nil
}

// genOrderNo 生成公开订单号
func genOrderNo(ctx context.Context) (string, error) {
	cols := dao.PledgeAddOrder.Columns()
	for i := 0; i < 5; i++ {
		no := "P" + gtime.Now().Format("YmdHis") + strings.ToUpper(grand.S(6))
		exist, err := dao.PledgeAddOrder.Ctx(ctx).Where(cols.OrderNo, no).Exist()
		if err != nil {
			return "", gerror.Wrap(err, "生成订单号失败")
		}
		if !exist {
			return no, nil
		}
	}
	return "", gerror.New("生成订单号失败，请重试")
}

func round4(value float64) float64 {
	return math.Round(value*1e4) / 1e4
}

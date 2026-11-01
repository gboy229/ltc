// Package settlement
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package settlement

import (
	"context"
	pledgeconsts "hotgo/addons/pledge/consts"
	mainconsts "hotgo/internal/consts"
	isc "hotgo/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/pledge/model/dao"
	"hotgo/addons/pledge/model/entity"
	"hotgo/addons/pledge/model/input/incomein"
	"hotgo/addons/pledge/service"
)

// createPushReward 生成下级静态收益对应的一级直推收益
func (s *sSettlement) createPushReward(
	ctx context.Context,
	order *entity.PledgeOrder,
	staticAmount float64,
	periodNo int,
) error {
	uplineId, err := isc.UserRelation().GetDirectUpline(ctx, int64(order.UserId))
	if err != nil {
		return err
	}
	if uplineId <= 0 {
		return nil
	}

	orderTable := dao.PledgeOrder.Table()
	periodTable := dao.PledgePeriodConfig.Table()
	var row struct {
		DirectRate float64 `json:"direct_rate"`
	}
	if err = dao.PledgeOrder.Ctx(ctx).
		LeftJoin(periodTable+" c", "c.period_months="+orderTable+".period_months").
		Fields("COALESCE(MAX(c.direct_reward_rate), 0) AS direct_rate").
		Where(orderTable+".user_id", uplineId).
		Where(orderTable+".status", pledgeconsts.PledgeOrderStatusActive).
		Where(orderTable+".is_early_refund", pledgeconsts.PledgeOrderIsEarlyRefundNo).
		Where("c.status", mainconsts.StatusEnabled).
		Scan(&row); err != nil {
		return gerror.Wrap(err, "查询上级最高直推率失败")
	}
	directRate := row.DirectRate
	if directRate <= 0 {
		return nil
	}

	directAmount := round4(staticAmount * directRate / 100)
	if directAmount <= 0 {
		return nil
	}

	_, _, err = service.Income().Create(ctx, &incomein.CreateInp{
		UserId:       uint64(uplineId),
		OrderId:      order.Id,
		FromUserId:   order.UserId,
		Type:         "direct",
		Amount:       directAmount,
		BaseAmount:   staticAmount,
		Rate:         directRate,
		PeriodNo:     periodNo,
		BizSn:        order.OrderNo,
		Remark:       "直推收益",
		AssetBizType: "direct_income",
	})
	if err != nil {
		return err
	}
	return nil
}

// createLotReward 按套餐 reward_per_lot 发放手数收益。
// 返回本次发放金额和今日手数收益；重复发放时返回 0。
func (s *sSettlement) createLotReward(
	ctx context.Context,
	order *entity.PledgeOrder,
	dueDays int,
	periodNo int,
) (lotAmount, lotToday float64, err error) {
	var plan *entity.PledgePlan
	if err = dao.PledgePlan.Ctx(ctx).WherePri(order.PlanId).Scan(&plan); err != nil {
		return 0, 0, gerror.Wrap(err, "查询套餐手数收益失败")
	}
	if plan == nil || plan.RewardPerLot <= 0 {
		return 0, 0, nil
	}

	lotToday = plan.RewardPerLot
	lotAmount = round4(plan.RewardPerLot * float64(dueDays))
	if lotAmount <= 0 {
		return 0, 0, nil
	}

	_, inserted, err := service.Income().Create(ctx, &incomein.CreateInp{
		UserId:       order.UserId,
		OrderId:      order.Id,
		Type:         "lot",
		Amount:       lotAmount,
		BaseAmount:   plan.RewardPerLot,
		Rate:         0,
		PeriodNo:     periodNo,
		BizSn:        order.OrderNo,
		Remark:       "手数收益",
		AssetBizType: "lot_income",
	})
	if err != nil {
		return 0, 0, err
	}
	if !inserted {
		return 0, 0, nil
	}
	return lotAmount, lotToday, nil
}

// createDiffReward 按链式差额发放级差奖励
func (s *sSettlement) createDiffReward(ctx context.Context, order *entity.PledgeOrder, staticAmount float64, periodNo int) error {
	if staticAmount <= 0 {
		return nil
	}

	var rows []struct {
		Pid        int64   `json:"pid"`
		DiffReward float64 `json:"diff_reward"`
	}
	if err := g.DB().Model("hg_addon_users_tree t").
		LeftJoin("hg_addon_users_user u", "u.id=t.pid").
		LeftJoin("hg_addon_users_level l", "l.id=u.level_id").
		Fields("t.pid", "COALESCE(l.diff_reward,0) AS diff_reward").
		Where("t.uid", order.UserId).
		Where("t.level > ?", 0).
		OrderAsc("t.level").
		Scan(&rows); err != nil {
		return gerror.Wrap(err, "查询级差关系失败")
	}

	lastRate := 0.0
	for _, row := range rows {
		ownRate := row.DiffReward
		rate := round4(ownRate - lastRate)
		if rate <= 0 {
			continue
		}

		amount := round4(staticAmount * rate / 100)
		if amount <= 0 {
			continue
		}

		_, _, err := service.Income().Create(ctx, &incomein.CreateInp{
			UserId:       uint64(row.Pid),
			OrderId:      order.Id,
			FromUserId:   order.UserId,
			Type:         "team",
			Amount:       amount,
			BaseAmount:   staticAmount,
			Rate:         rate,
			PeriodNo:     periodNo,
			BizSn:        order.OrderNo,
			Remark:       "团队收益",
			AssetBizType: "team_income",
		})
		if err != nil {
			return err
		}

		lastRate = ownRate
	}
	return nil
}

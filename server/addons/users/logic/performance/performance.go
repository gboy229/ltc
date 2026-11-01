// Package performance
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package performance

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"hotgo/addons/users/model/dao"
	"hotgo/addons/users/model/entity"
	"hotgo/addons/users/service"
	mainconsts "hotgo/internal/consts"
	"hotgo/internal/library/ratelimit"
)

type sPerformance struct{}

func NewPerformance() *sPerformance { return &sPerformance{} }

func init() { service.RegisterPerformance(NewPerformance()) }

// RefreshByUser 刷新该用户及所有上级的业绩
func (s *sPerformance) RefreshByUser(ctx context.Context, userId int64) error {
	if userId <= 0 {
		return nil
	}

	treeCols := dao.UsersTree.Columns()
	var rows []*entity.UsersTree
	if err := dao.UsersTree.Ctx(ctx).
		Fields(treeCols.Pid).
		Where(treeCols.Uid, userId).
		OrderDesc(treeCols.Level).
		Scan(&rows); err != nil {
		return gerror.Wrap(err, "查询用户上级失败")
	}

	userIds := map[int64]struct{}{userId: {}}
	for _, row := range rows {
		userIds[int64(row.Pid)] = struct{}{}
	}

	for id := range userIds {
		allowed, err := ratelimit.Allow(ctx, fmt.Sprintf("users:performance:refresh:%d", id), 30*time.Second)
		if err != nil {
			return err
		}
		if !allowed {
			continue
		}

		if err = s.refreshUser(ctx, id); err != nil {
			g.Log().Warningf(ctx, "refresh performance failed, uid:%d err:%+v", id, err)
			continue
		}
		_, _ = g.Redis().Do(ctx, "DEL", fmt.Sprintf("users:performance:refresh:%d", id))
	}
	return nil
}

func (s *sPerformance) refreshUser(ctx context.Context, userId int64) error {
	var personal float64
	var err error
	if personal, err = g.DB().Model("hg_addon_pledge_order").
		Where("user_id", userId).
		Where("status", 1).
		Sum("amount"); err != nil {
		return gerror.Wrap(err, "查询个人业绩失败")
	}

	var downline struct {
		DirectAmount float64 `json:"direct_amount"`
		TeamAmount   float64 `json:"team_amount"`
	}
	if err = g.DB().Model("hg_addon_users_tree t").
		LeftJoin("hg_addon_pledge_order o", "o.user_id=t.uid AND o.status=1").
		Fields(
			"COALESCE(SUM(CASE WHEN t.level=1 THEN o.amount ELSE 0 END),0) AS direct_amount",
			"COALESCE(SUM(o.amount),0) AS team_amount",
		).
		Where("t.pid", userId).
		Where("t.level > ?", 0).
		Scan(&downline); err != nil {
		return gerror.Wrap(err, "查询直推/团队业绩失败")
	}

	cols := dao.UsersPerformance.Columns()
	_, err = dao.UsersPerformance.Ctx(ctx).Data(g.Map{
		cols.UserId:         userId,
		cols.PersonalAmount: round4(personal),
		cols.DirectAmount:   round4(downline.DirectAmount),
		cols.TeamAmount:     round4(downline.TeamAmount),
		cols.UpdatedAt:      gtime.Now(),
	}).Save()
	if err != nil {
		return gerror.Wrap(err, "更新业绩汇总失败")
	}

	return s.syncLevel(ctx, userId, personal, downline.TeamAmount)
}

func (s *sPerformance) syncLevel(ctx context.Context, userId int64, personal, team float64) error {
	userCols := dao.UsersUser.Columns()
	var user *entity.UsersUser
	if err := dao.UsersUser.Ctx(ctx).
		Fields(userCols.LevelId, userCols.LockLevel).
		WherePri(userId).
		Scan(&user); err != nil {
		return gerror.Wrap(err, "查询用户等级失败")
	}
	if user == nil {
		return nil
	}

	levelCols := dao.UsersLevel.Columns()
	var target *entity.UsersLevel
	if err := dao.UsersLevel.Ctx(ctx).
		Where(levelCols.Status, mainconsts.StatusEnabled).
		Where(levelCols.PersonalAmount+" <= ?", personal).
		Where(levelCols.TeamAmount+" <= ?", team).
		OrderDesc(levelCols.PersonalAmount).
		OrderDesc(levelCols.TeamAmount).
		OrderDesc(levelCols.Id).
		Scan(&target); err != nil {
		return gerror.Wrap(err, "查询目标等级失败")
	}
	if target == nil {
		return nil
	}

	targetId := int(target.Id)
	currentId := user.LevelId
	if targetId == currentId {
		return nil
	}
	if user.LockLevel == 1 && targetId < currentId {
		return nil
	}

	err := dao.UsersUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, updateErr := dao.UsersUser.Ctx(ctx).
			WherePri(userId).
			Data(userCols.LevelId, targetId).
			Update(); updateErr != nil {
			return gerror.Wrap(updateErr, "更新用户等级失败")
		}

		logCols := dao.UsersLevelLog.Columns()
		if _, insertErr := dao.UsersLevelLog.Ctx(ctx).Data(g.Map{
			logCols.UserId:         userId,
			logCols.FromLevelId:    currentId,
			logCols.ToLevelId:      targetId,
			logCols.PersonalAmount: round4(personal),
			logCols.TeamAmount:     round4(team),
			logCols.CreatedAt:      gtime.Now(),
		}).Insert(); insertErr != nil {
			return gerror.Wrap(insertErr, "写入等级变更记录失败")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func round4(value float64) float64 {
	return math.Round(value*1e4) / 1e4
}

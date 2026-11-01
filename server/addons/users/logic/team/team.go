// Package team
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package team

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"hotgo/addons/users/api/api/team"
	"hotgo/addons/users/model/dao"
)

// Push 我的直推
func (s *sTeam) Push(ctx context.Context, userId int64, in *team.PushReq) (list []*team.PushItem, totalCount int, err error) {

	userTable := dao.UsersUser.Table()
	treeTable := dao.UsersTree.Table()
	levelTable := dao.UsersLevel.Table()

	mod := dao.UsersUser.Ctx(ctx).
		LeftJoin(treeTable+" t", "t.uid="+userTable+".id").
		LeftJoin(levelTable+" l", "l.id="+userTable+".level_id").
		Where("t.pid", userId).
		Where("t.level", 1)

	if totalCount, err = mod.Count(); err != nil {
		return nil, 0, gerror.Wrap(err, "统计直推成员失败")
	}
	if totalCount == 0 {
		return
	}

	var rows []struct {
		Uid       string      `json:"uid"`
		LevelName string      `json:"level_name"`
		CreatedAt *gtime.Time `json:"created_at"`
	}
	if err = mod.Fields(userTable+".uid", "l.name AS level_name", userTable+".created_at").
		OrderDesc("t.id").
		Page(in.Page, in.PerPage).
		Scan(&rows); err != nil {
		return nil, 0, gerror.Wrap(err, "获取直推成员失败")
	}

	list = make([]*team.PushItem, 0, len(rows))
	for _, row := range rows {
		registeredAt := ""
		if row.CreatedAt != nil {
			registeredAt = row.CreatedAt.Format("Y-m-d H:i")
		}
		list = append(list, &team.PushItem{
			Uid:          row.Uid,
			LevelName:    row.LevelName,
			RegisteredAt: registeredAt,
		})
	}
	return
}

// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"

	"hotgo/addons/users/api/api/team"
	"hotgo/addons/users/library/contexts"
	"hotgo/addons/users/service"
)

var (
	Team = cTeam{}
)

type cTeam struct{}

// Push 我的直推
func (c *cTeam) Push(ctx context.Context, req *team.PushReq) (res *team.PushRes, err error) {
	list, totalCount, err := service.Team().Push(ctx, contexts.GetUsersUserId(ctx), req)
	if err != nil {
		return
	}
	if list == nil {
		list = []*team.PushItem{}
	}
	res = new(team.PushRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}

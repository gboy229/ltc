// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"
	"hotgo/addons/users/api/api/center"
	"hotgo/addons/users/library/contexts"
	"hotgo/addons/users/service"
)

var (
	Center = cCenter{}
)

type cCenter struct{}

// Info 当前用户信息
func (c *cCenter) Info(ctx context.Context, req *center.InfoReq) (res *center.InfoRes, err error) {
	return service.UserCenter().Info(ctx, contexts.GetUsersUserId(ctx), req)
}

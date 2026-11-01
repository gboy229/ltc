// Package api
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package api

import (
	"context"

	"hotgo/addons/pledge/api/api/plan"
	"hotgo/addons/pledge/service"
)

var (
	Plan = cPlan{}
)

type cPlan struct{}

// List 获取启用的投资套餐
func (c *cPlan) List(ctx context.Context, req *plan.ListReq) (res *plan.ListRes, err error) {
	return service.Plan().List(ctx, req)
}

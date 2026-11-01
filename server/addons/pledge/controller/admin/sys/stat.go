// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"hotgo/addons/pledge/api/admin/stat"
	"hotgo/addons/pledge/service"
)

var (
	Stat = cStat{}
)

type cStat struct{}

// Stat 数据统计
func (c *cStat) Stat(ctx context.Context, req *stat.StatReq) (res *stat.StatRes, err error) {
	data, err := service.SysStat().Stat(ctx, &req.StatInp)
	if err != nil {
		return
	}

	res = new(stat.StatRes)
	res.StatModel = data
	return
}

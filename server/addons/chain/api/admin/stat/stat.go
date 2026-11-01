// Package stat
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package stat

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/chain/model/input/sysin"
)

// StatReq 数据统计
type StatReq struct {
	g.Meta `path:"/stat" method:"get" tags:"数据统计" summary:"获取充值提现统计"`
	sysin.StatInp
}

type StatRes struct {
	*sysin.StatModel
}

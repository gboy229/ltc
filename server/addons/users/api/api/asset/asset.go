// Package asset
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package asset

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/model/input/form"
)

// AccountReq 当前用户资产账户
type AccountReq struct {
	g.Meta `path:"/asset/account" method:"get" tags:"用户资产" summary:"当前用户资产账户"`
}

type AccountRes struct {
	Total     float64 `json:"total"     dc:"总资产"`
	Available float64 `json:"available" dc:"可用余额"`
	Frozen    float64 `json:"frozen"    dc:"冻结金额"`
}

// LogListReq 当前用户资产流水
type LogListReq struct {
	g.Meta `path:"/asset/logs" method:"get" tags:"用户资产" summary:"当前用户资产流水"`
	form.PageReq
	Type string `json:"type" v:"required|in:income,expense#类型不能为空|类型不合法" dc:"类型：income/expense"`
}

type LogItem struct {
	Type   string  `json:"type"   dc:"类型：income/expense"`
	Remark string  `json:"remark" dc:"备注"`
	Date   string  `json:"date"   dc:"发生时间"`
	Amount float64 `json:"amount" dc:"金额（正数）"`
}

type LogListRes struct {
	form.PageRes
	List []*LogItem `json:"list" dc:"资产流水"`
}

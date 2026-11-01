// Package plan
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package plan

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Item 前台投资套餐项（仅返回前端使用字段，不含数据库主键）
type Item struct {
	Name              string  `json:"name"              dc:"套餐名称"`
	Amount            float64 `json:"amount"            dc:"投资金额"`
	PeriodMonths      int     `json:"periodMonths"      dc:"投资周期（月）"`
	MonthlyRewardInfo string  `json:"monthlyRewardInfo" dc:"月化收益展示区间"`
	MonthlyRewardRate float64 `json:"monthlyRewardRate" dc:"固定月化收益率（%）"`
	TotalReturn       float64 `json:"totalReturn"       dc:"预计到期总额"`
}

// ListReq 获取启用的投资套餐
type ListReq struct {
	g.Meta `path:"/plan/list" method:"get" tags:"质押投资" summary:"获取启用的投资套餐"`
}

type ListRes struct {
	List []*Item `json:"list" dc:"数据列表"`
}

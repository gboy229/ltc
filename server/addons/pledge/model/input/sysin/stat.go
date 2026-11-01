// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// StatInp 数据统计
type StatInp struct {
	StartAt *gtime.Time `json:"startAt" dc:"开始时间"`
	EndAt   *gtime.Time `json:"endAt"   dc:"结束时间"`
}

// StatModel 数据统计
type StatModel struct {
	Order  *OrderStat   `json:"order"  dc:"订单统计"`
	Income *IncomeStat  `json:"income" dc:"收益统计"`
	Trend  []*TrendItem `json:"trend"  dc:"按天趋势"`
}

// OrderStat 订单统计
type OrderStat struct {
	Count         int     `json:"count"         dc:"订单数"`
	Amount        float64 `json:"amount"        dc:"订单金额"`
	ActiveCount   int     `json:"activeCount"   dc:"进行中订单数"`
	FinishedCount int     `json:"finishedCount" dc:"已完成订单数"`
}

// IncomeStat 收益统计
type IncomeStat struct {
	Count  int               `json:"count"  dc:"产出笔数"`
	Amount float64           `json:"amount" dc:"产出金额"`
	ByType []*IncomeTypeItem `json:"byType" dc:"按类型统计"`
}

// IncomeTypeItem 收益类型统计
type IncomeTypeItem struct {
	Type   string  `json:"type"   dc:"类型"`
	Label  string  `json:"label"  dc:"类型名称"`
	Count  int     `json:"count"  dc:"笔数"`
	Amount float64 `json:"amount" dc:"金额"`
}

// TrendItem 按天趋势
type TrendItem struct {
	Date         string  `json:"date"         dc:"日期"`
	OrderAmount  float64 `json:"orderAmount"  dc:"订单金额"`
	IncomeAmount float64 `json:"incomeAmount" dc:"产出金额"`
}

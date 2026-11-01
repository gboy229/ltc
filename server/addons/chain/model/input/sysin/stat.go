// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// StatInp 充值提现统计
type StatInp struct {
	StartAt *gtime.Time `json:"startAt" dc:"开始时间"`
	EndAt   *gtime.Time `json:"endAt"   dc:"结束时间"`
	Chain   string      `json:"chain"   dc:"链"`
	Symbol  string      `json:"symbol"  dc:"币种"`
}

// StatModel 充值提现统计
type StatModel struct {
	Recharge *RechargeStat `json:"recharge" dc:"充值统计"`
	Withdraw *WithdrawStat `json:"withdraw" dc:"提现统计"`
	Trend    []*TrendItem  `json:"trend"    dc:"按天趋势"`
}

// TrendItem 按天趋势
type TrendItem struct {
	Date           string  `json:"date"           dc:"日期"`
	RechargeAmount float64 `json:"rechargeAmount" dc:"充值金额"`
	WithdrawAmount float64 `json:"withdrawAmount" dc:"提现金额"`
}

// RechargeStat 充值统计
type RechargeStat struct {
	TotalCount  int            `json:"totalCount"  dc:"充值笔数"`
	TotalAmount float64        `json:"totalAmount" dc:"充值总额"`
	List        []*StatDayItem `json:"list"        dc:"按天明细"`
}

// StatDayItem 充值按天统计
type StatDayItem struct {
	Date   string  `json:"date"   dc:"日期"`
	Count  int     `json:"count"  dc:"笔数"`
	Amount float64 `json:"amount" dc:"金额"`
}

// WithdrawStat 提现统计
type WithdrawStat struct {
	TotalCount    int                `json:"totalCount"    dc:"提现笔数"`
	PayAmount     float64            `json:"payAmount"     dc:"申请总额"`
	ReceiveAmount float64            `json:"receiveAmount" dc:"到账总额"`
	Fee           float64            `json:"fee"           dc:"手续费合计"`
	List          []*WithdrawDayItem `json:"list"          dc:"按天明细"`
}

// WithdrawDayItem 提现按天统计
type WithdrawDayItem struct {
	Date          string  `json:"date"          dc:"日期"`
	Count         int     `json:"count"         dc:"笔数"`
	PayAmount     float64 `json:"payAmount"     dc:"申请金额"`
	ReceiveAmount float64 `json:"receiveAmount" dc:"到账金额"`
	Fee           float64 `json:"fee"           dc:"手续费"`
}

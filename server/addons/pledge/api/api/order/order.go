// Package order
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package order

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/model/input/form"
)

// BuyReq 提交报单
type BuyReq struct {
	g.Meta       `path:"/order/buy" method:"post" tags:"质押订单" summary:"提交报单"`
	PlanName     string `json:"planName"     v:"required#请选择投资套餐" dc:"套餐名称"`
	PeriodMonths int    `json:"periodMonths" v:"required|min:1#请选择投资周期" dc:"投资周期（月）"`
}

type BuyRes struct {
	OrderNo   string  `json:"orderNo" dc:"订单号"`
	Amount    float64 `json:"amount"  dc:"报单金额"`
	ToAddress string  `json:"toAddress"  dc:"入金地址"`
	Status    int     `json:"status"  dc:"状态：0=待处理，1=已处理"`
}

// ListReq 报单列表
type ListReq struct {
	g.Meta `path:"/order/list" method:"get" tags:"质押订单" summary:"报单列表"`
	form.PageReq
	Status string `json:"status" dc:"状态筛选：active/completed，不传为全部"`
}

type ListItem struct {
	Id                   string  `json:"id"           orm:"id"            dc:"报单号"`
	Amount               float64 `json:"amount"       orm:"amount"        dc:"入金金额"`
	TodayRelease         float64 `json:"todayRelease" orm:"today_release" dc:"今日释放"`
	TotalRelease         float64 `json:"totalRelease" orm:"total_release" dc:"累计释放"`
	Period               int     `json:"period"               orm:"period"                 dc:"投资周期（月）"`
	InvestDate           string  `json:"investDate"           orm:"invest_date"            dc:"投资日期"`
	Status               string  `json:"status"               orm:"status"                 dc:"状态：active/completed"`
	EarlyExitPenaltyRate float64 `json:"earlyExitPenaltyRate" orm:"early_exit_penalty_rate" dc:"提前退本违约金比例（%）"`
}

type ListRes struct {
	form.PageRes
	List []*ListItem `json:"list" dc:"报单列表"`
}

// StatisticsReq 报单统计
type StatisticsReq struct {
	g.Meta `path:"/order/statistics" method:"get" tags:"质押订单" summary:"报单统计"`
}

type StatisticsRes struct {
	Total          float64 `json:"total"          dc:"入金总额"`
	TodayRelease   float64 `json:"todayRelease"   dc:"今日释放"`
	TotalRelease   float64 `json:"totalRelease"   dc:"累计释放"`
	EarningsTotal  float64 `json:"earningsTotal"  dc:"收益总额（不含本金）"`
	ReferralIncome float64 `json:"referralIncome" dc:"直推"`
	StaticIncome   float64 `json:"staticIncome"   dc:"静态收益总额"`
}

// RefundReq 提前退本
type RefundReq struct {
	g.Meta    `path:"/order/refund" method:"post" tags:"质押订单" summary:"提前退本"`
	OrderNo   string `json:"orderNo"   v:"required#订单号不能为空" dc:"订单号"`
	Nonce     string `json:"nonce"     v:"required#签名挑战值不能为空" dc:"签名挑战值"`
	Signature string `json:"signature" v:"required#钱包签名不能为空" dc:"钱包签名"`
}

type RefundRes struct {
	Amount float64 `json:"amount" dc:"实际退款金额"`
}

// StatusReq 查询报单状态
type StatusReq struct {
	g.Meta  `path:"/order/status" method:"get" tags:"质押订单" summary:"查询报单状态"`
	OrderNo string `json:"orderNo" v:"required#订单号不能为空" dc:"订单号"`
}

type StatusRes struct {
	OrderNo           string  `json:"orderNo"           dc:"订单号"`
	Amount            float64 `json:"amount"            dc:"报单金额"`
	Status            int     `json:"status"            dc:"报单状态：0=待处理，1=已处理"`
	PledgeOrderStatus int     `json:"pledgeOrderStatus" dc:"正式订单状态：0=未生成，1=进行中"`
	PaidAt            string  `json:"paidAt"            dc:"支付确认时间"`
}

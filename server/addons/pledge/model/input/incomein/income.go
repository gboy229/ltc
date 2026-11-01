// Package incomein
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package incomein

// CreateInp 写入一笔收益并入账
type CreateInp struct {
	UserId       uint64  `json:"userId"     dc:"收益用户"`
	OrderId      uint64  `json:"orderId"    dc:"来源订单"`
	FromUserId   uint64  `json:"fromUserId" dc:"来源用户"`
	Type         string  `json:"type"       dc:"收益类型：static/direct/principal"`
	Amount       float64 `json:"amount"     dc:"收益金额"`
	BaseAmount   float64 `json:"baseAmount" dc:"计算基数"`
	Rate         float64 `json:"rate"       dc:"收益率"`
	PeriodNo     int     `json:"periodNo"   dc:"期数"`
	BizSn        string  `json:"bizSn"        dc:"业务单号"`
	Remark       string  `json:"remark"       dc:"备注"`
	AssetBizType string  `json:"assetBizType" dc:"资产流水业务类型，为空时按收益类型映射"`
}

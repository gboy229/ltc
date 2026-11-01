// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"hotgo/addons/chain/consts"
	"hotgo/addons/chain/model/entity"
	"hotgo/internal/library/dict"
	"hotgo/internal/model/input/form"
)

// WithdrawViewInp 获取指定提现记录信息
type WithdrawViewInp struct {
	Id int64 `json:"id" v:"required#提现记录ID不能为空" dc:"提现记录ID"`
}

type WithdrawViewModel struct {
	entity.ChainWithdraw
}

// WithdrawListInp 获取提现记录列表
type WithdrawListInp struct {
	form.PageReq
	Keyword         string        `json:"keyword" dc:"用户关键词：用户ID/钱包地址"`
	Uid             int           `json:"uid" dc:"用户ID"`
	Chain           string        `json:"chain" dc:"链"`
	Symbol          string        `json:"symbol" dc:"币种"`
	Hash            string        `json:"hash" dc:"哈希"`
	FromAddress     string        `json:"fromAddress" dc:"出币地址"`
	ToAddress       string        `json:"toAddress" dc:"用户钱包地址"`
	ContractAddress string        `json:"contractAddress" dc:"合约地址"`
	Status          *int          `json:"status" dc:"状态"`
	AmountMin       float64       `json:"amountMin" dc:"最小申请金额"`
	AmountMax       float64       `json:"amountMax" dc:"最大申请金额"`
	CreatedAt       []*gtime.Time `json:"createdAt" dc:"申请时间"`
	AuditAt         []*gtime.Time `json:"auditAt" dc:"审核时间"`
}

type WithdrawListModel struct {
	entity.ChainWithdraw
}

// WithdrawStatusInp 提现审核
type WithdrawStatusInp struct {
	Id     int64  `json:"id" v:"required#提现记录ID不能为空" dc:"提现记录ID"`
	Status int    `json:"status" dc:"审核状态，1=审核通过，-1=审核拒绝"`
	Remark string `json:"remark" dc:"审核备注"`
}

func (in *WithdrawStatusInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		err = gerror.New("提现记录ID不能为空")
		return
	}
	if !dict.HasOptionKey(consts.ChainWithdrawAuditOptions, in.Status) {
		err = gerror.New("审核状态不合法，只能审核通过或拒绝")
		return
	}
	if gstr.LenRune(in.Remark) > 255 {
		err = gerror.New("审核备注不能超过255个字符")
		return
	}
	return
}

type WithdrawStatusModel struct{}

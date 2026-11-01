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

// ContractEditInp 修改/新增合约
type ContractEditInp struct {
	entity.ChainContract
}

func (in *ContractEditInp) Filter(ctx context.Context) (err error) {
	in.Chain = gstr.ToUpper(gstr.Trim(in.Chain))
	in.Address = gstr.Trim(in.Address)
	in.Symbol = gstr.ToUpper(gstr.Trim(in.Symbol))

	if in.Chain == "" {
		err = gerror.New("链名不能为空")
		return
	}
	if gstr.LenRune(in.Chain) > 20 {
		err = gerror.New("链名不能超过20个字符")
		return
	}
	if in.Address == "" {
		err = gerror.New("合约地址不能为空")
		return
	}
	if gstr.LenRune(in.Address) > 100 {
		err = gerror.New("合约地址不能超过100个字符")
		return
	}
	if in.Symbol == "" {
		err = gerror.New("币种不能为空")
		return
	}
	if gstr.LenRune(in.Symbol) > 30 {
		err = gerror.New("币种不能超过30个字符")
		return
	}
	if in.Decimal < 0 || in.Decimal > 36 {
		err = gerror.New("精度必须在0~36之间")
		return
	}
	if in.Scale < 0 || in.Scale > in.Decimal {
		err = gerror.New("保留小数位必须在0~精度之间")
		return
	}
	if in.MinAmount < 0 || in.MaxAmount < 0 {
		err = gerror.New("提币数量不能为负数")
		return
	}
	if in.MinAmount > 0 && in.MaxAmount > 0 && in.MinAmount > in.MaxAmount {
		err = gerror.New("最大提币数量不能小于最小提币数量")
		return
	}
	if in.MultipleAmount < 0 {
		err = gerror.New("提币倍数不能为负数")
		return
	}
	if in.VerifyAmount < 0 {
		err = gerror.New("提币审核数量不能为负数")
		return
	}
	if in.DayLimit < 0 {
		err = gerror.New("日上限笔数不能为负数")
		return
	}
	if !dict.HasOptionKey(consts.ChainFeeTypeOptions, in.FeeType) {
		err = gerror.New("手续费类型不合法")
		return
	}
	if in.Rate < 0 {
		err = gerror.New("手续费费率不能为负数")
		return
	}
	if in.FeeType == consts.ChainFeeTypePercent && in.Rate > 100 {
		err = gerror.New("百分比手续费费率不能大于100")
		return
	}
	if !dict.HasOptionKey(consts.ChainSwitchOptions, in.RechargeSwitch) {
		err = gerror.New("充币开关值不合法")
		return
	}
	if !dict.HasOptionKey(consts.ChainSwitchOptions, in.WithdrawSwitch) {
		err = gerror.New("提现开关值不合法")
		return
	}
	return
}

type ContractEditModel struct{}

// ContractUpdateFields 修改合约字段过滤，扫块游标等系统字段不在此列
type ContractUpdateFields struct {
	Chain          string  `json:"chain"          description:"链名"`
	Address        string  `json:"address"        description:"合约地址"`
	Symbol         string  `json:"symbol"         description:"币种"`
	Decimal        int     `json:"decimal"        description:"精度"`
	MinAmount      float64 `json:"minAmount"      description:"最小提币数量，0=不限制"`
	MaxAmount      float64 `json:"maxAmount"      description:"最大提币数量，0=不限制"`
	MultipleAmount int     `json:"multipleAmount" description:"提币数量倍数，0=不限制"`
	VerifyAmount   float64 `json:"verifyAmount"   description:"提币审核阈值，0=全部审核"`
	Scale          int     `json:"scale"          description:"金额保留小数位"`
	DayLimit       int     `json:"dayLimit"       description:"日上限笔数，0=不限制"`
	FeeType        int     `json:"feeType"        description:"手续费类型，0=百分比，1=个数"`
	Rate           float64 `json:"rate"           description:"手续费费率：百分比类型填1表示1%，固定类型填币数量"`
	RechargeSwitch int     `json:"rechargeSwitch" description:"充币开关：0=关闭，1=开启"`
	WithdrawSwitch int     `json:"withdrawSwitch" description:"提现开关：0=关闭，1=开启"`
}

// ContractInsertFields 新增合约字段过滤，扫块游标等系统字段不在此列
type ContractInsertFields struct {
	Chain          string  `json:"chain"          description:"链名"`
	Address        string  `json:"address"        description:"合约地址"`
	Symbol         string  `json:"symbol"         description:"币种"`
	Decimal        int     `json:"decimal"        description:"精度"`
	MinAmount      float64 `json:"minAmount"      description:"最小提币数量，0=不限制"`
	MaxAmount      float64 `json:"maxAmount"      description:"最大提币数量，0=不限制"`
	MultipleAmount int     `json:"multipleAmount" description:"提币数量倍数，0=不限制"`
	VerifyAmount   float64 `json:"verifyAmount"   description:"提币审核阈值，0=全部审核"`
	Scale          int     `json:"scale"          description:"金额保留小数位"`
	DayLimit       int     `json:"dayLimit"       description:"日上限笔数，0=不限制"`
	FeeType        int     `json:"feeType"        description:"手续费类型，0=百分比，1=个数"`
	Rate           float64 `json:"rate"           description:"手续费费率：百分比类型填1表示1%，固定类型填币数量"`
	RechargeSwitch int     `json:"rechargeSwitch" description:"充币开关：0=关闭，1=开启"`
	WithdrawSwitch int     `json:"withdrawSwitch" description:"提现开关：0=关闭，1=开启"`
}

// ContractDeleteInp 删除合约
type ContractDeleteInp struct {
	Ids []int64 `json:"ids" v:"required#请选择要删除的合约" dc:"合约ID"`
}

func (in *ContractDeleteInp) Filter(ctx context.Context) (err error) {
	if len(in.Ids) == 0 {
		err = gerror.New("请选择要删除的合约")
	}
	return
}

type ContractDeleteModel struct{}

// ContractViewInp 获取指定合约信息
type ContractViewInp struct {
	Id int64 `json:"id" v:"required#合约ID不能为空" dc:"合约ID"`
}

// ContractViewModel 合约详情（不返回扫块游标等系统字段）
type ContractViewModel struct {
	Id             uint64      `json:"id"             orm:"id"              dc:"合约ID"`
	Chain          string      `json:"chain"          orm:"chain"           dc:"链名"`
	Address        string      `json:"address"        orm:"address"         dc:"合约地址"`
	Symbol         string      `json:"symbol"         orm:"symbol"          dc:"币种"`
	Decimal        int         `json:"decimal"        orm:"decimal"         dc:"精度"`
	MinAmount      float64     `json:"minAmount"      orm:"min_amount"      dc:"最小提币数量"`
	MaxAmount      float64     `json:"maxAmount"      orm:"max_amount"      dc:"最大提币数量"`
	MultipleAmount int         `json:"multipleAmount" orm:"multiple_amount" dc:"提币数量倍数"`
	VerifyAmount   float64     `json:"verifyAmount"   orm:"verify_amount"   dc:"提币审核阈值"`
	Scale          int         `json:"scale"          orm:"scale"           dc:"金额保留小数位"`
	DayLimit       int         `json:"dayLimit"       orm:"day_limit"       dc:"日上限笔数"`
	FeeType        int         `json:"feeType"        orm:"fee_type"        dc:"手续费类型"`
	Rate           float64     `json:"rate"           orm:"rate"            dc:"手续费费率"`
	RechargeSwitch int         `json:"rechargeSwitch" orm:"recharge_switch" dc:"充币开关"`
	WithdrawSwitch int         `json:"withdrawSwitch" orm:"withdraw_switch" dc:"提现开关"`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      dc:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      dc:""`
}

// ContractListInp 获取合约列表
type ContractListInp struct {
	form.PageReq
	Chain   string `json:"chain" dc:"链名"`
	Symbol  string `json:"symbol" dc:"币种"`
	Address string `json:"address" dc:"合约地址"`
}

// ContractListModel 合约列表（不返回扫块游标等系统字段）
type ContractListModel struct {
	Id             uint64      `json:"id"             orm:"id"              dc:"合约ID"`
	Chain          string      `json:"chain"          orm:"chain"           dc:"链名"`
	Address        string      `json:"address"        orm:"address"         dc:"合约地址"`
	Symbol         string      `json:"symbol"         orm:"symbol"          dc:"币种"`
	Decimal        int         `json:"decimal"        orm:"decimal"         dc:"精度"`
	MinAmount      float64     `json:"minAmount"      orm:"min_amount"      dc:"最小提币数量"`
	MaxAmount      float64     `json:"maxAmount"      orm:"max_amount"      dc:"最大提币数量"`
	MultipleAmount int         `json:"multipleAmount" orm:"multiple_amount" dc:"提币数量倍数"`
	VerifyAmount   float64     `json:"verifyAmount"   orm:"verify_amount"   dc:"提币审核阈值"`
	Scale          int         `json:"scale"          orm:"scale"           dc:"金额保留小数位"`
	DayLimit       int         `json:"dayLimit"       orm:"day_limit"       dc:"日上限笔数"`
	FeeType        int         `json:"feeType"        orm:"fee_type"        dc:"手续费类型"`
	Rate           float64     `json:"rate"           orm:"rate"            dc:"手续费费率"`
	RechargeSwitch int         `json:"rechargeSwitch" orm:"recharge_switch" dc:"充币开关"`
	WithdrawSwitch int         `json:"withdrawSwitch" orm:"withdraw_switch" dc:"提现开关"`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      dc:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      dc:""`
}

// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import (
	"hotgo/internal/library/dict"
	"hotgo/internal/model"
)

const AddressCacheKey = "chain:wallet_addresses"

// WalletUidCacheKey 用户钱包地址->用户ID，与 users 插件 walletcache.HashKey 保持一致
const WalletUidCacheKey = "users:wallet:uid"
const AddressCacheListKey = "chain:wallet_addresses_list"

// QueueChainScanTransfer 扫块 Transfer 事件队列
const QueueChainScanTransfer = "chain:scan:transfer"

// 提现状态
const (
	ChainWithdrawStatusPending     = 0  // 审核中
	ChainWithdrawStatusAudited     = 1  // 已审核
	ChainWithdrawStatusTransferred = 2  // 已打币
	ChainWithdrawStatusFailed      = 3  // 失败
	ChainWithdrawStatusFinished    = 4  // 完成
	ChainWithdrawStatusCanceled    = -1 // 取消
	ChainWithdrawStatusSending     = 5  // 打币中（已抢占，正在链上转账）
)

// 手续费类型
const (
	ChainFeeTypePercent = 0 // 百分比
	ChainFeeTypeCount   = 1 // 个数
)

// 开关
const (
	ChainSwitchOff = 0 // 关闭
	ChainSwitchOn  = 1 // 开启
)

// 链驱动
const (
	ChainDriverEVM = "EVM" // EVM 兼容链
)

// ChainWithdrawStatusOptions 提现状态选项
var ChainWithdrawStatusOptions = []*model.Option{
	dict.GenWarningOption(ChainWithdrawStatusPending, "审核中"),
	dict.GenPrimaryOption(ChainWithdrawStatusAudited, "已审核"),
	dict.GenInfoOption(ChainWithdrawStatusSending, "打币中"),
	dict.GenSuccessOption(ChainWithdrawStatusTransferred, "已打币"),
	dict.GenErrorOption(ChainWithdrawStatusFailed, "失败"),
	dict.GenSuccessOption(ChainWithdrawStatusFinished, "完成"),
	dict.GenDefaultOption(ChainWithdrawStatusCanceled, "取消"),
}

// ChainWithdrawAuditOptions 提现审核选项，后台仅允许审核通过或拒绝
var ChainWithdrawAuditOptions = []*model.Option{
	dict.GenSuccessOption(ChainWithdrawStatusAudited, "审核通过"),
	dict.GenErrorOption(ChainWithdrawStatusCanceled, "审核拒绝"),
}

// 充值类型
const (
	ChainRechargeTypeOrder = "order" // 订单
	ChainRechargeTypeGas   = "gas"   // Gas
)

// ChainRechargeTypeOptions 充值类型选项
var ChainRechargeTypeOptions = []*model.Option{
	dict.GenPrimaryOption(ChainRechargeTypeOrder, "订单"),
	dict.GenWarningOption(ChainRechargeTypeGas, "Gas"),
}

// ChainFeeTypeOptions 手续费类型选项
var ChainFeeTypeOptions = []*model.Option{
	dict.GenDefaultOption(ChainFeeTypePercent, "百分比"),
	dict.GenInfoOption(ChainFeeTypeCount, "个数"),
}

// ChainSwitchOptions 开关选项
var ChainSwitchOptions = []*model.Option{
	dict.GenDefaultOption(ChainSwitchOff, "关闭"),
	dict.GenSuccessOption(ChainSwitchOn, "开启"),
}

// ChainDriverOptions 链驱动选项
var ChainDriverOptions = []*model.Option{
	dict.GenHashOption(ChainDriverEVM, "EVM"),
}

func init() {
	dict.RegisterEnums("ChainWithdrawStatusOptions", "提现状态选项", ChainWithdrawStatusOptions)
	dict.RegisterEnums("ChainWithdrawAuditOptions", "提现审核选项", ChainWithdrawAuditOptions)
	dict.RegisterEnums("ChainRechargeTypeOptions", "充值类型选项", ChainRechargeTypeOptions)
	dict.RegisterEnums("ChainFeeTypeOptions", "手续费类型选项", ChainFeeTypeOptions)
	dict.RegisterEnums("ChainSwitchOptions", "开关选项", ChainSwitchOptions)
	dict.RegisterEnums("ChainDriverOptions", "链驱动选项", ChainDriverOptions)
}

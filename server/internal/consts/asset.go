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

func init() {
	dict.RegisterEnums("AssetChangeTypeOptions", "资产变动类型选项", AssetChangeTypeOptions)
	dict.RegisterEnums("AssetBizTypeOptions", "资产业务类型选项", AssetBizTypeOptions)
}

// 资产账户标识
const (
	AssetAccountUSDT = "USDT" // USDT 账户
)

// 资产账务动作
const (
	AssetChangeTypeIncr     = "incr"     // 增加
	AssetChangeTypeDecr     = "decr"     // 扣减
	AssetChangeTypeFreeze   = "freeze"   // 冻结
	AssetChangeTypeUnfreeze = "unfreeze" // 解冻
	AssetChangeTypeSettle   = "settle"   // 扣减冻结
)

// AssetChangeTypeOptions 资产变动类型选项
var AssetChangeTypeOptions = []*model.Option{
	dict.GenSuccessOption(AssetChangeTypeIncr, "增加"),
	dict.GenWarningOption(AssetChangeTypeDecr, "扣减"),
	dict.GenInfoOption(AssetChangeTypeFreeze, "冻结"),
	dict.GenInfoOption(AssetChangeTypeUnfreeze, "解冻"),
	dict.GenErrorOption(AssetChangeTypeSettle, "扣减冻结"),
}

// 资产来源插件
const (
	AssetAddonSystem = "system" // 系统/后台
	AssetAddonChain  = "chain"  // 链插件
	AssetAddonPledge = "pledge" // 质押插件
)

// 内置业务类型，插件可自行定义扩展
const (
	AssetBizTypeChainRecharge     = "recharge"         // 链上充值
	AssetBizTypeChainWithdraw     = "withdraw"         // 链上提现
	AssetBizTypePledgeEarlyRefund = "early_refund"     // 质押提前退本
	AssetBizTypeStaticIncome      = "static_income"    // 静态收益
	AssetBizTypeDirectIncome      = "direct_income"    // 直推收益
	AssetBizTypeDiffIncome        = "diff_income"      // 级差收益
	AssetBizTypeLotIncome         = "lot_income"       // 手数收益
	AssetBizTypeTeamIncome        = "team_income"      // 团队收益
	AssetBizTypePrincipalReturn   = "principal_return" // 本金返还
)

// AssetBizTypeOptions 资产业务类型选项
var AssetBizTypeOptions = []*model.Option{
	dict.GenSuccessOption(AssetBizTypeChainRecharge, "链上充值"),
	dict.GenWarningOption(AssetBizTypeChainWithdraw, "链上提现"),
	dict.GenWarningOption(AssetBizTypePledgeEarlyRefund, "提前退本"),
	dict.GenSuccessOption(AssetBizTypeStaticIncome, "静态收益"),
	dict.GenSuccessOption(AssetBizTypeDirectIncome, "直推收益"),
	dict.GenSuccessOption(AssetBizTypeDiffIncome, "级差收益"),
	dict.GenSuccessOption(AssetBizTypeLotIncome, "手数收益"),
	dict.GenSuccessOption(AssetBizTypeTeamIncome, "团队收益"),
	dict.GenInfoOption(AssetBizTypePrincipalReturn, "本金返还"),
}

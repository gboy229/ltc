// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainContract is the golang structure of table hg_addon_chain_contract for DAO operations like Where/Data.
type ChainContract struct {
	g.Meta         `orm:"table:hg_addon_chain_contract, do:true"`
	Id             interface{} //
	Chain          interface{} // 链名
	Address        interface{} // 合约地址
	Symbol         interface{} // 币种
	Decimal        interface{} // 精度
	MinAmount      interface{} //
	MaxAmount      interface{} //
	MultipleAmount interface{} // 倍数
	VerifyAmount   interface{} //
	Scale          interface{} // 保留几位小数
	DayLimit       interface{} // 日上限笔数
	FeeType        interface{} // 手续费类型，0=百分比，1=个数
	Rate           interface{} //
	RechargeSwitch interface{} // 充币开关
	WithdrawSwitch interface{} // 提现开关
	CreatedAt      *gtime.Time //
	UpdatedAt      *gtime.Time //
	ScanBlock      interface{} // 已扫块高
	ScanBlockHash  interface{} // 已扫块哈希
}

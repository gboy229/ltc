// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainRecharge is the golang structure of table hg_addon_chain_recharge for DAO operations like Where/Data.
type ChainRecharge struct {
	g.Meta          `orm:"table:hg_addon_chain_recharge, do:true"`
	Id              interface{} //
	Type            interface{} //
	Uid             interface{} // 用户ID
	Chain           interface{} // 链
	Symbol          interface{} // 币种
	FromAddress     interface{} // 转出地址
	ToAddress       interface{} // 转入地址
	ContractAddress interface{} // 合约地址
	Amount          interface{} // 充值数量
	BlockHeight     interface{} // 块高
	Hash            interface{} // 交易哈希
	LogIndex        interface{} // 日志序号
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}

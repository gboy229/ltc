// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainWithdraw is the golang structure of table hg_addon_chain_withdraw for DAO operations like Where/Data.
type ChainWithdraw struct {
	g.Meta          `orm:"table:hg_addon_chain_withdraw, do:true"`
	Id              interface{} //
	Uid             interface{} // 用户ID
	Chain           interface{} // 链
	Symbol          interface{} // 币种
	ContractAddress interface{} // 合约地址
	FromAddress     interface{} // 转出地址
	ToAddress       interface{} // 转入地址
	PayAmount       interface{} // 付款数量，含手续费
	ReceiveAmount   interface{} // 收款数量
	Fee             interface{} // 手续费
	Rate            interface{} // 手续费费率
	Hash            interface{} // 哈希
	Status          interface{} // 状态，0=审核中，1=已审核，2=已打币，3=失败，4=完成，-1=取消
	Version         interface{} // 乐观锁
	AuditBy         interface{} // 审核人ID
	AuditAt         *gtime.Time // 审核时间
	Remark          interface{} // 审核备注
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
}

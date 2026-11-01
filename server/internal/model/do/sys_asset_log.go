// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAssetLog is the golang structure of table hg_sys_asset_log for DAO operations like Where/Data.
type SysAssetLog struct {
	g.Meta          `orm:"table:hg_sys_asset_log, do:true"`
	Id              interface{} // 流水ID
	UserId          interface{} // 用户ID（hg_addon_users_user.id）
	AddonName       interface{} // æ¥æºæ’ä»¶
	Account         interface{} // 资产账户标识，如 USDT
	ChangeType      interface{} // 账务动作：incr=增加，decr=扣减，freeze=冻结，unfreeze=解冻，settle=扣减冻结
	BizType         interface{} // 业务类型：chain_recharge、chain_withdraw、cms_reward 等
	BizId           interface{} // 业务单号
	BizSn           interface{} // 业务外部单号（如链上tx hash），便于排查
	BeforeAvailable interface{} // 变动前可用余额
	ChangeAvailable interface{} // 可用余额变动值（正负）
	AfterAvailable  interface{} // 变动后可用余额
	BeforeFrozen    interface{} // 变动前冻结金额
	ChangeFrozen    interface{} // 冻结金额变动值（正负）
	AfterFrozen     interface{} // 变动后冻结金额
	Remark          interface{} // 备注
	OperatorId      interface{} // 操作人ID（后台管理员，0=系统）
	CreatedAt       *gtime.Time // 创建时间
}

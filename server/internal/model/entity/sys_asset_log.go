// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAssetLog is the golang structure for table sys_asset_log.
type SysAssetLog struct {
	Id              uint64      `json:"id"              orm:"id"               description:"流水ID"`
	UserId          uint64      `json:"userId"          orm:"user_id"          description:"用户ID（hg_addon_users_user.id）"`
	AddonName       string      `json:"addonName"       orm:"addon_name"       description:"æ¥æºæ’ä»¶"`
	Account         string      `json:"account"         orm:"account"          description:"资产账户标识，如 USDT"`
	ChangeType      string      `json:"changeType"      orm:"change_type"      description:"账务动作：incr=增加，decr=扣减，freeze=冻结，unfreeze=解冻，settle=扣减冻结"`
	BizType         string      `json:"bizType"         orm:"biz_type"         description:"业务类型：chain_recharge、chain_withdraw、cms_reward 等"`
	BizId           uint64      `json:"bizId"           orm:"biz_id"           description:"业务单号"`
	BizSn           string      `json:"bizSn"           orm:"biz_sn"           description:"业务外部单号（如链上tx hash），便于排查"`
	BeforeAvailable float64     `json:"beforeAvailable" orm:"before_available" description:"变动前可用余额"`
	ChangeAvailable float64     `json:"changeAvailable" orm:"change_available" description:"可用余额变动值（正负）"`
	AfterAvailable  float64     `json:"afterAvailable"  orm:"after_available"  description:"变动后可用余额"`
	BeforeFrozen    float64     `json:"beforeFrozen"    orm:"before_frozen"    description:"变动前冻结金额"`
	ChangeFrozen    float64     `json:"changeFrozen"    orm:"change_frozen"    description:"冻结金额变动值（正负）"`
	AfterFrozen     float64     `json:"afterFrozen"     orm:"after_frozen"     description:"变动后冻结金额"`
	Remark          string      `json:"remark"          orm:"remark"           description:"备注"`
	OperatorId      int64       `json:"operatorId"      orm:"operator_id"      description:"操作人ID（后台管理员，0=系统）"`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:"创建时间"`
}

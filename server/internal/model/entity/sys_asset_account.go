// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAssetAccount is the golang structure for table sys_asset_account.
type SysAssetAccount struct {
	Id        uint64      `json:"id"        orm:"id"         description:"账户ID"`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"用户ID（hg_addon_users_user.id）"`
	Account   string      `json:"account"   orm:"account"    description:"资产账户标识，如 USDT"`
	Available float64     `json:"available" orm:"available"  description:"可用余额"`
	Frozen    float64     `json:"frozen"    orm:"frozen"     description:"冻结金额"`
	Version   int         `json:"version"   orm:"version"    description:"乐观锁版本号"`
	Status    int         `json:"status"    orm:"status"     description:"状态：1=正常，2=禁用"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`
}

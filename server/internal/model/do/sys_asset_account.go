// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAssetAccount is the golang structure of table hg_sys_asset_account for DAO operations like Where/Data.
type SysAssetAccount struct {
	g.Meta    `orm:"table:hg_sys_asset_account, do:true"`
	Id        interface{} // 账户ID
	UserId    interface{} // 用户ID（hg_addon_users_user.id）
	Account   interface{} // 资产账户标识，如 USDT
	Available interface{} // 可用余额
	Frozen    interface{} // 冻结金额
	Version   interface{} // 乐观锁版本号
	Status    interface{} // 状态：1=正常，2=禁用
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}

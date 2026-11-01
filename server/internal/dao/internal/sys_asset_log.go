// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysAssetLogDao is the data access object for table hg_sys_asset_log.
type SysAssetLogDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns SysAssetLogColumns // columns contains all the column names of Table for convenient usage.
}

// SysAssetLogColumns defines and stores column names for table hg_sys_asset_log.
type SysAssetLogColumns struct {
	Id              string // 流水ID
	UserId          string // 用户ID（hg_addon_users_user.id）
	AddonName       string // æ¥æºæ’ä»¶
	Account         string // 资产账户标识，如 USDT
	ChangeType      string // 账务动作：incr=增加，decr=扣减，freeze=冻结，unfreeze=解冻，settle=扣减冻结
	BizType         string // 业务类型：chain_recharge、chain_withdraw、cms_reward 等
	BizId           string // 业务单号
	BizSn           string // 业务外部单号（如链上tx hash），便于排查
	BeforeAvailable string // 变动前可用余额
	ChangeAvailable string // 可用余额变动值（正负）
	AfterAvailable  string // 变动后可用余额
	BeforeFrozen    string // 变动前冻结金额
	ChangeFrozen    string // 冻结金额变动值（正负）
	AfterFrozen     string // 变动后冻结金额
	Remark          string // 备注
	OperatorId      string // 操作人ID（后台管理员，0=系统）
	CreatedAt       string // 创建时间
}

// sysAssetLogColumns holds the columns for table hg_sys_asset_log.
var sysAssetLogColumns = SysAssetLogColumns{
	Id:              "id",
	UserId:          "user_id",
	AddonName:       "addon_name",
	Account:         "account",
	ChangeType:      "change_type",
	BizType:         "biz_type",
	BizId:           "biz_id",
	BizSn:           "biz_sn",
	BeforeAvailable: "before_available",
	ChangeAvailable: "change_available",
	AfterAvailable:  "after_available",
	BeforeFrozen:    "before_frozen",
	ChangeFrozen:    "change_frozen",
	AfterFrozen:     "after_frozen",
	Remark:          "remark",
	OperatorId:      "operator_id",
	CreatedAt:       "created_at",
}

// NewSysAssetLogDao creates and returns a new DAO object for table data access.
func NewSysAssetLogDao() *SysAssetLogDao {
	return &SysAssetLogDao{
		group:   "default",
		table:   "hg_sys_asset_log",
		columns: sysAssetLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysAssetLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SysAssetLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SysAssetLogDao) Columns() SysAssetLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SysAssetLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SysAssetLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SysAssetLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

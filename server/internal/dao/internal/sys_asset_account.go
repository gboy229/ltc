// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysAssetAccountDao is the data access object for table hg_sys_asset_account.
type SysAssetAccountDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns SysAssetAccountColumns // columns contains all the column names of Table for convenient usage.
}

// SysAssetAccountColumns defines and stores column names for table hg_sys_asset_account.
type SysAssetAccountColumns struct {
	Id        string // 账户ID
	UserId    string // 用户ID（hg_addon_users_user.id）
	Account   string // 资产账户标识，如 USDT
	Available string // 可用余额
	Frozen    string // 冻结金额
	Version   string // 乐观锁版本号
	Status    string // 状态：1=正常，2=禁用
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// sysAssetAccountColumns holds the columns for table hg_sys_asset_account.
var sysAssetAccountColumns = SysAssetAccountColumns{
	Id:        "id",
	UserId:    "user_id",
	Account:   "account",
	Available: "available",
	Frozen:    "frozen",
	Version:   "version",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewSysAssetAccountDao creates and returns a new DAO object for table data access.
func NewSysAssetAccountDao() *SysAssetAccountDao {
	return &SysAssetAccountDao{
		group:   "default",
		table:   "hg_sys_asset_account",
		columns: sysAssetAccountColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysAssetAccountDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SysAssetAccountDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *SysAssetAccountDao) Columns() SysAssetAccountColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SysAssetAccountDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SysAssetAccountDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SysAssetAccountDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChainConfigDao is the data access object for table hg_addon_chain_config.
type ChainConfigDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns ChainConfigColumns // columns contains all the column names of Table for convenient usage.
}

// ChainConfigColumns defines and stores column names for table hg_addon_chain_config.
type ChainConfigColumns struct {
	Id              string //
	Driver          string // 驱动
	Chain           string // 网络
	RpcUrl          string // rpc地址
	ChainId         string // 链id
	BrowserUrl      string // 浏览器地址
	TransferAddress string // 转账地址
	TransferPrivate string // 转账私钥
	CreatedAt       string //
	UpdatedAt       string //
	Confirmations   string // 确认数
}

// chainConfigColumns holds the columns for table hg_addon_chain_config.
var chainConfigColumns = ChainConfigColumns{
	Id:              "id",
	Driver:          "driver",
	Chain:           "chain",
	RpcUrl:          "rpc_url",
	ChainId:         "chain_id",
	BrowserUrl:      "browser_url",
	TransferAddress: "transfer_address",
	TransferPrivate: "transfer_private",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	Confirmations:   "confirmations",
}

// NewChainConfigDao creates and returns a new DAO object for table data access.
func NewChainConfigDao() *ChainConfigDao {
	return &ChainConfigDao{
		group:   "default",
		table:   "hg_addon_chain_config",
		columns: chainConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ChainConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ChainConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ChainConfigDao) Columns() ChainConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ChainConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ChainConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ChainConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

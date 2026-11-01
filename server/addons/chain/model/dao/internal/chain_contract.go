// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChainContractDao is the data access object for table hg_addon_chain_contract.
type ChainContractDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns ChainContractColumns // columns contains all the column names of Table for convenient usage.
}

// ChainContractColumns defines and stores column names for table hg_addon_chain_contract.
type ChainContractColumns struct {
	Id             string //
	Chain          string // 链名
	Address        string // 合约地址
	Symbol         string // 币种
	Decimal        string // 精度
	MinAmount      string //
	MaxAmount      string //
	MultipleAmount string // 倍数
	VerifyAmount   string //
	Scale          string // 保留几位小数
	DayLimit       string // 日上限笔数
	FeeType        string // 手续费类型，0=百分比，1=个数
	Rate           string //
	RechargeSwitch string // 充币开关
	WithdrawSwitch string // 提现开关
	CreatedAt      string //
	UpdatedAt      string //
	ScanBlock      string // 已扫块高
	ScanBlockHash  string // 已扫块哈希
}

// chainContractColumns holds the columns for table hg_addon_chain_contract.
var chainContractColumns = ChainContractColumns{
	Id:             "id",
	Chain:          "chain",
	Address:        "address",
	Symbol:         "symbol",
	Decimal:        "decimal",
	MinAmount:      "min_amount",
	MaxAmount:      "max_amount",
	MultipleAmount: "multiple_amount",
	VerifyAmount:   "verify_amount",
	Scale:          "scale",
	DayLimit:       "day_limit",
	FeeType:        "fee_type",
	Rate:           "rate",
	RechargeSwitch: "recharge_switch",
	WithdrawSwitch: "withdraw_switch",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	ScanBlock:      "scan_block",
	ScanBlockHash:  "scan_block_hash",
}

// NewChainContractDao creates and returns a new DAO object for table data access.
func NewChainContractDao() *ChainContractDao {
	return &ChainContractDao{
		group:   "default",
		table:   "hg_addon_chain_contract",
		columns: chainContractColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ChainContractDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ChainContractDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ChainContractDao) Columns() ChainContractColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ChainContractDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ChainContractDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ChainContractDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChainRechargeDao is the data access object for table hg_addon_chain_recharge.
type ChainRechargeDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns ChainRechargeColumns // columns contains all the column names of Table for convenient usage.
}

// ChainRechargeColumns defines and stores column names for table hg_addon_chain_recharge.
type ChainRechargeColumns struct {
	Id              string //
	Type            string //
	Uid             string // 用户ID
	Chain           string // 链
	Symbol          string // 币种
	FromAddress     string // 转出地址
	ToAddress       string // 转入地址
	ContractAddress string // 合约地址
	Amount          string // 充值数量
	BlockHeight     string // 块高
	Hash            string // 交易哈希
	LogIndex        string // 日志序号
	CreatedAt       string //
	UpdatedAt       string //
}

// chainRechargeColumns holds the columns for table hg_addon_chain_recharge.
var chainRechargeColumns = ChainRechargeColumns{
	Id:              "id",
	Type:            "type",
	Uid:             "uid",
	Chain:           "chain",
	Symbol:          "symbol",
	FromAddress:     "from_address",
	ToAddress:       "to_address",
	ContractAddress: "contract_address",
	Amount:          "amount",
	BlockHeight:     "block_height",
	Hash:            "hash",
	LogIndex:        "log_index",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewChainRechargeDao creates and returns a new DAO object for table data access.
func NewChainRechargeDao() *ChainRechargeDao {
	return &ChainRechargeDao{
		group:   "default",
		table:   "hg_addon_chain_recharge",
		columns: chainRechargeColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ChainRechargeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ChainRechargeDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ChainRechargeDao) Columns() ChainRechargeColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ChainRechargeDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ChainRechargeDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ChainRechargeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

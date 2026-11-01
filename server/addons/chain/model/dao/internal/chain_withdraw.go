// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ChainWithdrawDao is the data access object for table hg_addon_chain_withdraw.
type ChainWithdrawDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns ChainWithdrawColumns // columns contains all the column names of Table for convenient usage.
}

// ChainWithdrawColumns defines and stores column names for table hg_addon_chain_withdraw.
type ChainWithdrawColumns struct {
	Id              string //
	Uid             string // 用户ID
	Chain           string // 链
	Symbol          string // 币种
	ContractAddress string // 合约地址
	FromAddress     string // 转出地址
	ToAddress       string // 转入地址
	PayAmount       string // 付款数量，含手续费
	ReceiveAmount   string // 收款数量
	Fee             string // 手续费
	Rate            string // 手续费费率
	Hash            string // 哈希
	Status          string // 状态，0=审核中，1=已审核，2=已打币，3=失败，4=完成，-1=取消
	Version         string // 乐观锁
	AuditBy         string // 审核人ID
	AuditAt         string // 审核时间
	Remark          string // 审核备注
	CreatedAt       string //
	UpdatedAt       string //
}

// chainWithdrawColumns holds the columns for table hg_addon_chain_withdraw.
var chainWithdrawColumns = ChainWithdrawColumns{
	Id:              "id",
	Uid:             "uid",
	Chain:           "chain",
	Symbol:          "symbol",
	ContractAddress: "contract_address",
	FromAddress:     "from_address",
	ToAddress:       "to_address",
	PayAmount:       "pay_amount",
	ReceiveAmount:   "receive_amount",
	Fee:             "fee",
	Rate:            "rate",
	Hash:            "hash",
	Status:          "status",
	Version:         "version",
	AuditBy:         "audit_by",
	AuditAt:         "audit_at",
	Remark:          "remark",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewChainWithdrawDao creates and returns a new DAO object for table data access.
func NewChainWithdrawDao() *ChainWithdrawDao {
	return &ChainWithdrawDao{
		group:   "default",
		table:   "hg_addon_chain_withdraw",
		columns: chainWithdrawColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *ChainWithdrawDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *ChainWithdrawDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *ChainWithdrawDao) Columns() ChainWithdrawColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *ChainWithdrawDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *ChainWithdrawDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *ChainWithdrawDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

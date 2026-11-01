// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PledgeIncomeLogDao is the data access object for table hg_addon_pledge_income_log.
type PledgeIncomeLogDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns PledgeIncomeLogColumns // columns contains all the column names of Table for convenient usage.
}

// PledgeIncomeLogColumns defines and stores column names for table hg_addon_pledge_income_log.
type PledgeIncomeLogColumns struct {
	Id         string //
	UserId     string // 收益用户
	OrderId    string // 来源报单
	FromUserId string // 来源用户（直推/团队）
	Type       string // static/direct/indirect/team/...
	Amount     string // 收益金额
	BaseAmount string // 计算基数
	Rate       string // 收益率
	PeriodNo   string // 第几期
	Status     string // 1=已入账 2=已撤销
	CreatedAt  string //
	UpdatedAt  string //
}

// pledgeIncomeLogColumns holds the columns for table hg_addon_pledge_income_log.
var pledgeIncomeLogColumns = PledgeIncomeLogColumns{
	Id:         "id",
	UserId:     "user_id",
	OrderId:    "order_id",
	FromUserId: "from_user_id",
	Type:       "type",
	Amount:     "amount",
	BaseAmount: "base_amount",
	Rate:       "rate",
	PeriodNo:   "period_no",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewPledgeIncomeLogDao creates and returns a new DAO object for table data access.
func NewPledgeIncomeLogDao() *PledgeIncomeLogDao {
	return &PledgeIncomeLogDao{
		group:   "default",
		table:   "hg_addon_pledge_income_log",
		columns: pledgeIncomeLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PledgeIncomeLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PledgeIncomeLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PledgeIncomeLogDao) Columns() PledgeIncomeLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PledgeIncomeLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PledgeIncomeLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PledgeIncomeLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

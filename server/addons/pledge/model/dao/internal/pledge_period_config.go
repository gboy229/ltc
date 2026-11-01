// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PledgePeriodConfigDao is the data access object for table hg_addon_pledge_period_config.
type PledgePeriodConfigDao struct {
	table   string                    // table is the underlying table name of the DAO.
	group   string                    // group is the database configuration group name of current DAO.
	columns PledgePeriodConfigColumns // columns contains all the column names of Table for convenient usage.
}

// PledgePeriodConfigColumns defines and stores column names for table hg_addon_pledge_period_config.
type PledgePeriodConfigColumns struct {
	Id                string //
	PeriodMonths      string //
	MonthlyRewardInfo string //
	MonthlyRewardRate string //
	DirectRewardRate  string //
	Status            string //
	CreatedAt         string //
	UpdatedAt         string //
}

// pledgePeriodConfigColumns holds the columns for table hg_addon_pledge_period_config.
var pledgePeriodConfigColumns = PledgePeriodConfigColumns{
	Id:                "id",
	PeriodMonths:      "period_months",
	MonthlyRewardInfo: "monthly_reward_info",
	MonthlyRewardRate: "monthly_reward_rate",
	DirectRewardRate:  "direct_reward_rate",
	Status:            "status",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
}

// NewPledgePeriodConfigDao creates and returns a new DAO object for table data access.
func NewPledgePeriodConfigDao() *PledgePeriodConfigDao {
	return &PledgePeriodConfigDao{
		group:   "default",
		table:   "hg_addon_pledge_period_config",
		columns: pledgePeriodConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PledgePeriodConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PledgePeriodConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PledgePeriodConfigDao) Columns() PledgePeriodConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PledgePeriodConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PledgePeriodConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PledgePeriodConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

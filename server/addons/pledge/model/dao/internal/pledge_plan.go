// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PledgePlanDao is the data access object for table hg_addon_pledge_plan.
type PledgePlanDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns PledgePlanColumns // columns contains all the column names of Table for convenient usage.
}

// PledgePlanColumns defines and stores column names for table hg_addon_pledge_plan.
type PledgePlanColumns struct {
	Id                   string //
	Name                 string // å¥—é¤åç§°
	Amount               string // å…¥é‡‘é‡‘é¢
	PeriodMonths         string // 投资周期（月）
	DailyLots            string // 每日交易手数
	RewardPerLot         string // 每手金额（USDT）
	EarlyExitPenaltyRate string // 提前退本违约金比例（%）
	Status               string // 状态：1=启用，2=禁用
	CreatedAt            string // 创建时间
	UpdatedAt            string // 更新时间
}

// pledgePlanColumns holds the columns for table hg_addon_pledge_plan.
var pledgePlanColumns = PledgePlanColumns{
	Id:                   "id",
	Name:                 "name",
	Amount:               "amount",
	PeriodMonths:         "period_months",
	DailyLots:            "daily_lots",
	RewardPerLot:         "reward_per_lot",
	EarlyExitPenaltyRate: "early_exit_penalty_rate",
	Status:               "status",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewPledgePlanDao creates and returns a new DAO object for table data access.
func NewPledgePlanDao() *PledgePlanDao {
	return &PledgePlanDao{
		group:   "default",
		table:   "hg_addon_pledge_plan",
		columns: pledgePlanColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PledgePlanDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PledgePlanDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PledgePlanDao) Columns() PledgePlanColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PledgePlanDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PledgePlanDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PledgePlanDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

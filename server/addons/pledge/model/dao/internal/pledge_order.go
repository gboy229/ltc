// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PledgeOrderDao is the data access object for table hg_addon_pledge_order.
type PledgeOrderDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns PledgeOrderColumns // columns contains all the column names of Table for convenient usage.
}

// PledgeOrderColumns defines and stores column names for table hg_addon_pledge_order.
type PledgeOrderColumns struct {
	Id                   string // 主键ID
	OrderNo              string // 订单号（与报单号一致）
	UserId               string // 用户ID
	PlanId               string // 套餐ID
	Amount               string // 订单金额
	PeriodMonths         string //
	TodayIncome          string // 今天收益
	IncIncome            string // 累计收益
	CycleNum             string // 总天数
	IncIncomeNum         string // 累计收益次数
	EarlyExitPenaltyRate string // 提前退本违约金比例（%）
	RechargeId           string // 关联充值记录ID
	Status               string // 状态：1=进行中
	IsEarlyRefund        string // 是否提前退本：0=否，1=是
	SettleAt             string // 最后结算时间
	CreatedAt            string // 创建时间
	UpdatedAt            string // 更新时间
}

// pledgeOrderColumns holds the columns for table hg_addon_pledge_order.
var pledgeOrderColumns = PledgeOrderColumns{
	Id:                   "id",
	OrderNo:              "order_no",
	UserId:               "user_id",
	PlanId:               "plan_id",
	Amount:               "amount",
	PeriodMonths:         "period_months",
	TodayIncome:          "today_income",
	IncIncome:            "inc_income",
	CycleNum:             "cycle_num",
	IncIncomeNum:         "inc_income_num",
	EarlyExitPenaltyRate: "early_exit_penalty_rate",
	RechargeId:           "recharge_id",
	Status:               "status",
	IsEarlyRefund:        "is_early_refund",
	SettleAt:             "settle_at",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewPledgeOrderDao creates and returns a new DAO object for table data access.
func NewPledgeOrderDao() *PledgeOrderDao {
	return &PledgeOrderDao{
		group:   "default",
		table:   "hg_addon_pledge_order",
		columns: pledgeOrderColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PledgeOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PledgeOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PledgeOrderDao) Columns() PledgeOrderColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PledgeOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PledgeOrderDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PledgeOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

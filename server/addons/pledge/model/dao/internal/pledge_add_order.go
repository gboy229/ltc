// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// PledgeAddOrderDao is the data access object for table hg_addon_pledge_add_order.
type PledgeAddOrderDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns PledgeAddOrderColumns // columns contains all the column names of Table for convenient usage.
}

// PledgeAddOrderColumns defines and stores column names for table hg_addon_pledge_add_order.
type PledgeAddOrderColumns struct {
	Id        string // 主键ID
	OrderNo   string // 公开订单号
	UserId    string // 用户ID
	PlanId    string // 套餐ID
	Amount    string // 报单金额
	Status    string // 状态：0=待处理，1=已处理
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// pledgeAddOrderColumns holds the columns for table hg_addon_pledge_add_order.
var pledgeAddOrderColumns = PledgeAddOrderColumns{
	Id:        "id",
	OrderNo:   "order_no",
	UserId:    "user_id",
	PlanId:    "plan_id",
	Amount:    "amount",
	Status:    "status",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewPledgeAddOrderDao creates and returns a new DAO object for table data access.
func NewPledgeAddOrderDao() *PledgeAddOrderDao {
	return &PledgeAddOrderDao{
		group:   "default",
		table:   "hg_addon_pledge_add_order",
		columns: pledgeAddOrderColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *PledgeAddOrderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *PledgeAddOrderDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *PledgeAddOrderDao) Columns() PledgeAddOrderColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *PledgeAddOrderDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *PledgeAddOrderDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *PledgeAddOrderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

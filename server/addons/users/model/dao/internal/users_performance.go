// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersPerformanceDao is the data access object for table hg_addon_users_performance.
type UsersPerformanceDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns UsersPerformanceColumns // columns contains all the column names of Table for convenient usage.
}

// UsersPerformanceColumns defines and stores column names for table hg_addon_users_performance.
type UsersPerformanceColumns struct {
	UserId         string // 用户ID
	PersonalAmount string // 个人业绩
	DirectAmount   string // 直推业绩
	TeamAmount     string // 团队业绩
	UpdatedAt      string // 更新时间
}

// usersPerformanceColumns holds the columns for table hg_addon_users_performance.
var usersPerformanceColumns = UsersPerformanceColumns{
	UserId:         "user_id",
	PersonalAmount: "personal_amount",
	DirectAmount:   "direct_amount",
	TeamAmount:     "team_amount",
	UpdatedAt:      "updated_at",
}

// NewUsersPerformanceDao creates and returns a new DAO object for table data access.
func NewUsersPerformanceDao() *UsersPerformanceDao {
	return &UsersPerformanceDao{
		group:   "default",
		table:   "hg_addon_users_performance",
		columns: usersPerformanceColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UsersPerformanceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UsersPerformanceDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UsersPerformanceDao) Columns() UsersPerformanceColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UsersPerformanceDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UsersPerformanceDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UsersPerformanceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

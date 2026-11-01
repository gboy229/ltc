// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersTreeDao is the data access object for table hg_addon_users_tree.
type UsersTreeDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns UsersTreeColumns // columns contains all the column names of Table for convenient usage.
}

// UsersTreeColumns defines and stores column names for table hg_addon_users_tree.
type UsersTreeColumns struct {
	Id    string //
	Uid   string // uid
	Pid   string // pid
	Level string // level
}

// usersTreeColumns holds the columns for table hg_addon_users_tree.
var usersTreeColumns = UsersTreeColumns{
	Id:    "id",
	Uid:   "uid",
	Pid:   "pid",
	Level: "level",
}

// NewUsersTreeDao creates and returns a new DAO object for table data access.
func NewUsersTreeDao() *UsersTreeDao {
	return &UsersTreeDao{
		group:   "default",
		table:   "hg_addon_users_tree",
		columns: usersTreeColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UsersTreeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UsersTreeDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UsersTreeDao) Columns() UsersTreeColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UsersTreeDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UsersTreeDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UsersTreeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

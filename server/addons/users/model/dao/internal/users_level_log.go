// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersLevelLogDao is the data access object for table hg_addon_users_level_log.
type UsersLevelLogDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns UsersLevelLogColumns // columns contains all the column names of Table for convenient usage.
}

// UsersLevelLogColumns defines and stores column names for table hg_addon_users_level_log.
type UsersLevelLogColumns struct {
	Id             string //
	UserId         string // 用户ID
	FromLevelId    string // 变更前等级ID
	ToLevelId      string // 变更后等级ID
	PersonalAmount string // 变更时个人业绩
	TeamAmount     string // 变更时团队业绩
	CreatedAt      string // 创建时间
}

// usersLevelLogColumns holds the columns for table hg_addon_users_level_log.
var usersLevelLogColumns = UsersLevelLogColumns{
	Id:             "id",
	UserId:         "user_id",
	FromLevelId:    "from_level_id",
	ToLevelId:      "to_level_id",
	PersonalAmount: "personal_amount",
	TeamAmount:     "team_amount",
	CreatedAt:      "created_at",
}

// NewUsersLevelLogDao creates and returns a new DAO object for table data access.
func NewUsersLevelLogDao() *UsersLevelLogDao {
	return &UsersLevelLogDao{
		group:   "default",
		table:   "hg_addon_users_level_log",
		columns: usersLevelLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UsersLevelLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UsersLevelLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UsersLevelLogDao) Columns() UsersLevelLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UsersLevelLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UsersLevelLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UsersLevelLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

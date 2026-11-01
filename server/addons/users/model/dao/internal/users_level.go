// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersLevelDao is the data access object for table hg_addon_users_level.
type UsersLevelDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns UsersLevelColumns // columns contains all the column names of Table for convenient usage.
}

// UsersLevelColumns defines and stores column names for table hg_addon_users_level.
type UsersLevelColumns struct {
	Id             string // 主键ID
	Name           string // 等级名称，例如 LV1、新手
	PersonalAmount string // 个人业绩
	TeamAmount     string // 团队业绩
	DiffReward     string // 级差奖励
	Remark         string // 备注说明
	Status         string // 状态：1=启用，2=禁用
	Sort           string // 排序值，越大越靠前
	CreatedAt      string // 创建时间
	UpdatedAt      string // 更新时间
}

// usersLevelColumns holds the columns for table hg_addon_users_level.
var usersLevelColumns = UsersLevelColumns{
	Id:             "id",
	Name:           "name",
	PersonalAmount: "personal_amount",
	TeamAmount:     "team_amount",
	DiffReward:     "diff_reward",
	Remark:         "remark",
	Status:         "status",
	Sort:           "sort",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
}

// NewUsersLevelDao creates and returns a new DAO object for table data access.
func NewUsersLevelDao() *UsersLevelDao {
	return &UsersLevelDao{
		group:   "default",
		table:   "hg_addon_users_level",
		columns: usersLevelColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UsersLevelDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UsersLevelDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UsersLevelDao) Columns() UsersLevelColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UsersLevelDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UsersLevelDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UsersLevelDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

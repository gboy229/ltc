// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersUserAccountDao is the data access object for table hg_addon_users_user_account.
type UsersUserAccountDao struct {
	table   string                  // table is the underlying table name of the DAO.
	group   string                  // group is the database configuration group name of current DAO.
	columns UsersUserAccountColumns // columns contains all the column names of Table for convenient usage.
}

// UsersUserAccountColumns defines and stores column names for table hg_addon_users_user_account.
type UsersUserAccountColumns struct {
	Id         string // 凭证ID
	UserId     string // 用户ID
	Type       string // 凭证类型:username,mobile,email,wallet,wechat,google...
	Account    string // 凭证值(手机号/邮箱/用户名/钱包地址/openid)
	Credential string // 凭证附加信息(可选,如钱包签名参数等)
	Verified   string // 是否已验证(手机/邮箱/钱包通常需验证后才可登录)
	Status     string // 状态
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
	DeletedAt  string // 删除时间
}

// usersUserAccountColumns holds the columns for table hg_addon_users_user_account.
var usersUserAccountColumns = UsersUserAccountColumns{
	Id:         "id",
	UserId:     "user_id",
	Type:       "type",
	Account:    "account",
	Credential: "credential",
	Verified:   "verified",
	Status:     "status",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewUsersUserAccountDao creates and returns a new DAO object for table data access.
func NewUsersUserAccountDao() *UsersUserAccountDao {
	return &UsersUserAccountDao{
		group:   "default",
		table:   "hg_addon_users_user_account",
		columns: usersUserAccountColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UsersUserAccountDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UsersUserAccountDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UsersUserAccountDao) Columns() UsersUserAccountColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UsersUserAccountDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UsersUserAccountDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UsersUserAccountDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UsersUserDao is the data access object for table hg_addon_users_user.
type UsersUserDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns UsersUserColumns // columns contains all the column names of Table for convenient usage.
}

// UsersUserColumns defines and stores column names for table hg_addon_users_user.
type UsersUserColumns struct {
	Id                   string // 用户ID
	Uid                  string // 对外唯一用户号
	Pid                  string // 上级用户ID(邀请关系)
	LevelId              string // 会员等级ID
	LockLevel            string // 是否锁定等级
	Nickname             string // 昵称
	Avatar               string // 头像
	Password             string // 登录密码(bcrypt哈希,自带盐)
	TradePassword        string // 交易密码(bcrypt哈希,自带盐)
	AuthenticatorKey     string // TOTP验证器密钥
	AuthenticatorEnabled string // 是否开启验证器
	AuthStatus           string // 实名认证状态
	Status               string // 状态
	Extend               string // 扩展资料(零散资料兜底,如性别/生日/地区等)
	LastLoginIp          string // 最后登录IP
	LastLoginTime        string // 最后登录时间
	CreatedAt            string // 创建时间
	UpdatedAt            string // 更新时间
	DeletedAt            string // 删除时间
}

// usersUserColumns holds the columns for table hg_addon_users_user.
var usersUserColumns = UsersUserColumns{
	Id:                   "id",
	Uid:                  "uid",
	Pid:                  "pid",
	LevelId:              "level_id",
	LockLevel:            "lock_level",
	Nickname:             "nickname",
	Avatar:               "avatar",
	Password:             "password",
	TradePassword:        "trade_password",
	AuthenticatorKey:     "authenticator_key",
	AuthenticatorEnabled: "authenticator_enabled",
	AuthStatus:           "auth_status",
	Status:               "status",
	Extend:               "extend",
	LastLoginIp:          "last_login_ip",
	LastLoginTime:        "last_login_time",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
	DeletedAt:            "deleted_at",
}

// NewUsersUserDao creates and returns a new DAO object for table data access.
func NewUsersUserDao() *UsersUserDao {
	return &UsersUserDao{
		group:   "default",
		table:   "hg_addon_users_user",
		columns: usersUserColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UsersUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UsersUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UsersUserDao) Columns() UsersUserColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UsersUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UsersUserDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UsersUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// CmsArticleDao is the data access object for table hg_addon_cms_article.
type CmsArticleDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns CmsArticleColumns // columns contains all the column names of Table for convenient usage.
}

// CmsArticleColumns defines and stores column names for table hg_addon_cms_article.
type CmsArticleColumns struct {
	Id          string // 主键ID
	ArticleNo   string // 前台公开ID
	Title       string // 标题
	Type        string // 类型：1=公告
	Cover       string // 封面图
	Content     string // 内容
	Sort        string // 排序，越大越靠前
	Status      string // 状态：1=正常，2=禁用
	IsRecommend string // 是否推荐：1=推荐，0=不推荐
	CreatedBy   string // 创建人
	UpdatedBy   string // 修改人
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// cmsArticleColumns holds the columns for table hg_addon_cms_article.
var cmsArticleColumns = CmsArticleColumns{
	Id:          "id",
	ArticleNo:   "article_no",
	Title:       "title",
	Type:        "type",
	Cover:       "cover",
	Content:     "content",
	Sort:        "sort",
	Status:      "status",
	IsRecommend: "is_recommend",
	CreatedBy:   "created_by",
	UpdatedBy:   "updated_by",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewCmsArticleDao creates and returns a new DAO object for table data access.
func NewCmsArticleDao() *CmsArticleDao {
	return &CmsArticleDao{
		group:   "default",
		table:   "hg_addon_cms_article",
		columns: cmsArticleColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CmsArticleDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *CmsArticleDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *CmsArticleDao) Columns() CmsArticleColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *CmsArticleDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CmsArticleDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CmsArticleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

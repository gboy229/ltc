// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// CmsArticle is the golang structure of table hg_addon_cms_article for DAO operations like Where/Data.
type CmsArticle struct {
	g.Meta      `orm:"table:hg_addon_cms_article, do:true"`
	Id          interface{} // 主键ID
	ArticleNo   interface{} // 前台公开ID
	Title       interface{} // 标题
	Type        interface{} // 类型：1=公告
	Cover       interface{} // 封面图
	Content     interface{} // 内容
	Sort        interface{} // 排序，越大越靠前
	Status      interface{} // 状态：1=正常，2=禁用
	IsRecommend interface{} // 是否推荐：1=推荐，0=不推荐
	CreatedBy   interface{} // 创建人
	UpdatedBy   interface{} // 修改人
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}

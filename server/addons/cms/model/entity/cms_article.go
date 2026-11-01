// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// CmsArticle is the golang structure for table cms_article.
type CmsArticle struct {
	Id          uint64      `json:"id"          orm:"id"           description:"主键ID"`
	ArticleNo   string      `json:"articleNo"   orm:"article_no"   description:"前台公开ID"`
	Title       string      `json:"title"       orm:"title"        description:"标题"`
	Type        int         `json:"type"        orm:"type"         description:"类型：1=公告"`
	Cover       string      `json:"cover"       orm:"cover"        description:"封面图"`
	Content     string      `json:"content"     orm:"content"      description:"内容"`
	Sort        int         `json:"sort"        orm:"sort"         description:"排序，越大越靠前"`
	Status      int         `json:"status"      orm:"status"       description:"状态：1=正常，2=禁用"`
	IsRecommend int         `json:"isRecommend" orm:"is_recommend" description:"是否推荐：1=推荐，0=不推荐"`
	CreatedBy   int64       `json:"createdBy"   orm:"created_by"   description:"创建人"`
	UpdatedBy   int64       `json:"updatedBy"   orm:"updated_by"   description:"修改人"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:"更新时间"`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:"删除时间"`
}

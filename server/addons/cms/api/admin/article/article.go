// Package article
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package article

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/cms/model/input/sysin"
	"hotgo/internal/model/input/form"
)

// ListReq 查询公告列表
type ListReq struct {
	g.Meta `path:"/article/list" method:"get" tags:"公告" summary:"获取公告列表"`
	sysin.ArticleListInp
}

type ListRes struct {
	List []*sysin.ArticleListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// ViewReq 获取指定公告
type ViewReq struct {
	g.Meta `path:"/article/view" method:"get" tags:"公告" summary:"获取指定公告"`
	sysin.ArticleViewInp
}

type ViewRes struct {
	*sysin.ArticleViewModel
}

// EditReq 修改/新增公告
type EditReq struct {
	g.Meta `path:"/article/edit" method:"post" tags:"公告" summary:"修改/新增公告"`
	sysin.ArticleEditInp
}

type EditRes struct{}

// DeleteReq 删除公告
type DeleteReq struct {
	g.Meta `path:"/article/delete" method:"post" tags:"公告" summary:"删除公告"`
	sysin.ArticleDeleteInp
}

type DeleteRes struct{}

// MaxSortReq 获取公告最大排序
type MaxSortReq struct {
	g.Meta `path:"/article/maxSort" method:"get" tags:"公告" summary:"获取公告最大排序"`
	sysin.ArticleMaxSortInp
}

type MaxSortRes struct {
	*sysin.ArticleMaxSortModel
}

// StatusReq 更新公告状态
type StatusReq struct {
	g.Meta `path:"/article/status" method:"post" tags:"公告" summary:"更新公告状态"`
	sysin.ArticleStatusInp
}

type StatusRes struct{}

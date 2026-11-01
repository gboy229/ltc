// Package consts
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package consts

import (
	"hotgo/internal/library/dict"
	"hotgo/internal/model"
)

func init() {
	dict.RegisterEnums("CmsArticleTypeOptions", "内容文章类型选项", CmsArticleTypeOptions)
	dict.RegisterEnums("CmsArticleRecommendOptions", "公告推荐选项", CmsArticleRecommendOptions)
}

// 内容文章类型
const (
	CmsArticleTypeNotice = 1 // 公告
)

// 公告首页推荐状态
const (
	CmsArticleRecommendNo  = 0 // 不推荐
	CmsArticleRecommendYes = 1 // 推荐
)

// CmsArticleTypeOptions 内容文章类型选项
var CmsArticleTypeOptions = []*model.Option{
	dict.GenInfoOption(CmsArticleTypeNotice, "公告"),
}

// CmsArticleRecommendOptions 公告推荐选项
var CmsArticleRecommendOptions = []*model.Option{
	dict.GenSuccessOption(CmsArticleRecommendYes, "推荐"),
	dict.GenDefaultOption(CmsArticleRecommendNo, "不推荐"),
}

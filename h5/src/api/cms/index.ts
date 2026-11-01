import { http } from '@/utils/http'

/** 前台公告摘要（不包含数据库主键等内部字段） */
export interface CmsArticleSummary {
  id: string
  title: string
  /** 发布日期，格式：YYYY-MM-DD */
  createdAt: string
}

/** 前台公告详情 */
export interface CmsArticleDetail extends CmsArticleSummary {
  content: string
}

export interface ArticleListRes {
  list: CmsArticleSummary[]
  page: number
  pageSize: number
  pageCount: number
  totalCount: number
}

export interface ArticleRecommendRes {
  list: CmsArticleSummary[]
}

/** 公告列表（公开接口，无需登录） */
export function getArticleList(
  params: { page?: number; pageSize?: number; title?: string } = {},
) {
  return http.get<ArticleListRes>('/cms/article/list', params, { silent: true })
}

/** 首页推荐公告（不分页，返回全部推荐公告） */
export function getRecommendArticles() {
  return http.get<ArticleRecommendRes>('/cms/article/recommend', undefined, {
    silent: true,
  })
}

/** 公告详情（公开接口，无需登录，article_no 走路径参数） */
export function getArticleView(articleNo: string) {
  return http.get<CmsArticleDetail>(
    `/cms/article/view/${encodeURIComponent(articleNo)}`,
    undefined,
    { silent: true },
  )
}

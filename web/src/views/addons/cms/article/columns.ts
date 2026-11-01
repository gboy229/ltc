import { renderOptionTag } from '@/utils';

export const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '公告标题',
    key: 'title',
    width: 260,
  },
  {
    title: '排序',
    key: 'sort',
    width: 80,
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return renderOptionTag('sys_normal_disable', row.status);
    },
  },
  {
    title: '推荐首页',
    key: 'isRecommend',
    width: 100,
    render(row) {
      return renderOptionTag('CmsArticleRecommendOptions', row.isRecommend);
    },
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 180,
  },
];

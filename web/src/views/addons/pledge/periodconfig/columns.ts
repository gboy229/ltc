import { renderOptionTag } from '@/utils';

export const columns = [
  { title: 'ID', key: 'id', width: 80 },
  {
    title: '投资周期',
    key: 'periodMonths',
    width: 120,
    render(row) {
      return renderOptionTag('PledgePeriodOptions', row.periodMonths);
    },
  },
  { title: '月化收益区间', key: 'monthlyRewardInfo', width: 160 },
  {
    title: '静态月化率',
    key: 'monthlyRewardRate',
    width: 130,
    render(row) {
      return `${Number(row.monthlyRewardRate || 0).toFixed(2)}%`;
    },
  },
  {
    title: '直推收益率',
    key: 'directRewardRate',
    width: 130,
    render(row) {
      return `${Number(row.directRewardRate || 0).toFixed(2)}%`;
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return renderOptionTag('sys_normal_disable', row.status);
    },
  },
  { title: '创建时间', key: 'createdAt', width: 180 },
  { title: '更新时间', key: 'updatedAt', width: 180 },
];

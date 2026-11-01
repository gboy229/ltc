import { renderOptionTag } from '@/utils';

export const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '等级名称',
    key: 'name',
    width: 120,
  },
  {
    title: '个人业绩',
    key: 'personalAmount',
    width: 120,
  },
  {
    title: '团队业绩',
    key: 'teamAmount',
    width: 120,
  },
  {
    title: '级差奖励',
    key: 'diffReward',
    width: 120,
    render(row) {
      return `${Number(row.diffReward || 0).toFixed(2)}%`;
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
];

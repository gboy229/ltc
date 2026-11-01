import { renderOptionTag } from '@/utils';

export const columns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '用户ID', key: 'userId', width: 110 },
  { title: '钱包地址', key: 'walletAddress', width: 220 },
  {
    title: '类型',
    key: 'type',
    width: 110,
    render(row) {
      return renderOptionTag('PledgeIncomeTypeOptions', row.type);
    },
  },
  {
    title: '收益金额',
    key: 'amount',
    width: 130,
    render(row) {
      return Number(row.amount || 0).toFixed(8);
    },
  },
  {
    title: '收益率',
    key: 'rate',
    width: 100,
    render(row) {
      return `${Number(row.rate || 0).toFixed(4)}%`;
    },
  },
  { title: '创建时间', key: 'createdAt', width: 180 },
];

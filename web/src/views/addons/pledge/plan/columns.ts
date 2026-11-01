import { renderOptionTag } from '@/utils';

export const columns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '套餐名称', key: 'name', width: 120 },
  {
    title: '投资金额(USDT)',
    key: 'amount',
    width: 140,
    render(row) {
      return Number(row.amount || 0).toFixed(2);
    },
  },
  {
    title: '投资周期',
    key: 'periodMonths',
    width: 110,
    render(row) {
      return renderOptionTag('PledgePeriodOptions', row.periodMonths);
    },
  },
  { title: '交易手数(日)', key: 'dailyLots', width: 120 },
  {
    title: '每手金额(U)',
    key: 'rewardPerLot',
    width: 120,
    render(row) {
      return Number(row.rewardPerLot || 0).toFixed(2);
    },
  },
  {
    title: '提前退本违约金',
    key: 'earlyExitPenaltyRate',
    width: 150,
    render(row) {
      return `${Number(row.earlyExitPenaltyRate || 0).toFixed(2)}%`;
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
];

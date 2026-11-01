export const columns = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '订单号', key: 'orderNo', width: 220 },
  { title: '用户ID', key: 'userId', width: 100 },
  { title: '钱包地址', key: 'walletAddress', width: 220 },
  {
    title: '入金金额',
    key: 'amount',
    width: 120,
    render(row) {
      return Number(row.amount || 0).toFixed(2);
    },
  },
  { title: '周期(月)', key: 'periodMonths', width: 90 },
  { title: '总天数', key: 'cycleNum', width: 90 },
  { title: '已结算天数', key: 'incIncomeNum', width: 110 },
  {
    title: '累计收益',
    key: 'incIncome',
    width: 120,
    render(row) {
      return Number(row.incIncome || 0).toFixed(4);
    },
  },
  {
    title: '订单状态',
    key: 'status',
    width: 100,
    render(row) {
      return row.status === 1 ? '进行中' : '已完成';
    },
  },
  {
    title: '提前退本',
    key: 'isEarlyRefund',
    width: 100,
    render(row) {
      return row.isEarlyRefund === 1 ? '是' : '否';
    },
  },
  { title: '最后结算', key: 'settleAt', width: 180 },
  { title: '创建时间', key: 'createdAt', width: 180 },
];

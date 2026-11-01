import { renderOptionTag } from '@/utils';

export const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 70,
  },
  {
    title: '链',
    key: 'chain',
    width: 90,
  },
  {
    title: '币种',
    key: 'symbol',
    width: 90,
  },
  {
    title: '合约地址',
    key: 'address',
    width: 260,
    ellipsis: { tooltip: true },
  },
  {
    title: '精度',
    key: 'decimal',
    width: 80,
  },
  {
    title: '提币限额',
    key: 'amountLimit',
    width: 160,
    render(row) {
      const min = row.minAmount > 0 ? row.minAmount : '不限';
      const max = row.maxAmount > 0 ? row.maxAmount : '不限';
      return `${min} ~ ${max}`;
    },
  },
  {
    title: '手续费',
    key: 'fee',
    width: 120,
    render(row) {
      return row.feeType === 0 ? `${row.rate}%` : `${row.rate} ${row.symbol}`;
    },
  },
  {
    title: '充币',
    key: 'rechargeSwitch',
    width: 90,
    render(row) {
      return renderOptionTag('ChainSwitchOptions', row.rechargeSwitch);
    },
  },
  {
    title: '提币',
    key: 'withdrawSwitch',
    width: 90,
    render(row) {
      return renderOptionTag('ChainSwitchOptions', row.withdrawSwitch);
    },
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 180,
  },
];

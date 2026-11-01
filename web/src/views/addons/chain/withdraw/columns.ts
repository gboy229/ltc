import { h } from 'vue';
import { NTag } from 'naive-ui';
import { renderOptionTag } from '@/utils';

function formatAmount(value: any) {
  const amount = Number(value);
  if (!Number.isFinite(amount)) {
    return value;
  }
  return amount.toLocaleString('en-US', { maximumFractionDigits: 8 });
}

export const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '用户ID',
    key: 'uid',
    width: 110,
    render(row) {
      return h(
        NTag,
        { size: 'small', bordered: false, type: 'info' },
        { default: () => `#${row.uid}` }
      );
    },
  },
  {
    title: '用户钱包',
    key: 'toAddress',
    width: 210,
    ellipsis: { tooltip: true },
  },
  {
    title: '出币地址',
    key: 'fromAddress',
    width: 210,
    ellipsis: { tooltip: true },
  },
  {
    title: '币种',
    key: 'symbol',
    width: 120,
    render(row) {
      return `${row.chain} / ${row.symbol}`;
    },
  },
  {
    title: '申请金额',
    key: 'payAmount',
    width: 150,
    render(row) {
      return formatAmount(row.payAmount);
    },
  },
  {
    title: '到账金额',
    key: 'receiveAmount',
    width: 160,
    render(row) {
      return h('span', { class: 'receive-text' }, formatAmount(row.receiveAmount));
    },
  },
  {
    title: '手续费',
    key: 'fee',
    width: 120,
    render(row) {
      return formatAmount(row.fee);
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 110,
    render(row) {
      return renderOptionTag('ChainWithdrawStatusOptions', row.status);
    },
  },
  {
    title: '申请时间',
    key: 'createdAt',
    width: 170,
  },
];

import { h } from 'vue';
import { NTag } from 'naive-ui';
import { renderOptionTag } from '@/utils';

// 金额展示：去掉多余的 0，并加千分位
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
    key: 'fromAddress',
    width: 210,
    ellipsis: { tooltip: true },
  },
  {
    title: '充值地址',
    key: 'toAddress',
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
    title: '充值数量',
    key: 'amount',
    width: 170,
    render(row) {
      return h('span', { class: 'amount-text' }, `+${formatAmount(row.amount)} ${row.symbol}`);
    },
  },
  {
    title: '类型',
    key: 'type',
    width: 90,
    render(row) {
      return renderOptionTag('ChainRechargeTypeOptions', row.type);
    },
  },
  {
    title: '交易哈希',
    key: 'hash',
    width: 220,
    ellipsis: { tooltip: true },
  },
  {
    title: '充值时间',
    key: 'createdAt',
    width: 170,
  },
];

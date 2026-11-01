import { FormSchema } from '@/components/Form';
import { useDictStore } from '@/store/modules/dict';
import { renderOptionTag } from '@/utils';

const dict = useDictStore();

export const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '用户ID',
    key: 'userId',
    width: 100,
  },
  {
    title: '用户UID',
    key: 'uid',
    width: 140,
  },
  {
    title: '钱包地址',
    key: 'walletAddress',
    width: 360,
    render(row) {
      return row.walletAddress || '—';
    },
  },
  {
    title: '资产账户',
    key: 'account',
    width: 100,
  },
  {
    title: '可用余额',
    key: 'available',
    width: 140,
    render(row) {
      return Number(row.available || 0).toFixed(2);
    },
  },
  {
    title: '冻结金额',
    key: 'frozen',
    width: 140,
    render(row) {
      return Number(row.frozen || 0).toFixed(2);
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
  {
    title: '更新时间',
    key: 'updatedAt',
    width: 180,
  },
];

export const schemas: FormSchema[] = [
  {
    field: 'userId',
    component: 'NInput',
    label: '用户ID',
    componentProps: {
      placeholder: '请输入用户ID',
    },
  },
  {
    field: 'uid',
    component: 'NInput',
    label: '用户UID',
    componentProps: {
      placeholder: '请输入用户UID',
    },
  },
  {
    field: 'walletAddress',
    component: 'NInput',
    label: '钱包地址',
    componentProps: {
      placeholder: '请输入钱包地址',
    },
  },
];

export function loadOptions() {
  dict.loadOptions(['sys_normal_disable']);
}

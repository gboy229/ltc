import { FormSchema } from '@/components/Form';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export const schemas: FormSchema[] = [
  {
    field: 'userId',
    component: 'NInputNumber',
    label: '用户ID',
    componentProps: { placeholder: '请输入用户ID', showButton: false },
  },
  {
    field: 'walletAddress',
    component: 'NInput',
    label: '钱包地址',
    componentProps: { placeholder: '请输入钱包地址' },
  },
  {
    field: 'orderNo',
    component: 'NInput',
    label: '订单号',
    componentProps: { placeholder: '请输入订单号' },
  },
  {
    field: 'status',
    component: 'NSelect',
    label: '订单状态',
    componentProps: {
      placeholder: '请选择状态',
      options: [
        { label: '进行中', value: 1 },
        { label: '已完成', value: 2 },
      ],
      clearable: true,
    },
  },
  {
    field: 'isEarlyRefund',
    component: 'NSelect',
    label: '提前退本',
    componentProps: {
      placeholder: '请选择',
      options: [
        { label: '否', value: 0 },
        { label: '是', value: 1 },
      ],
      clearable: true,
    },
  },
];

export function loadOptions() {
  dict.loadOptions(['sys_normal_disable']);
}

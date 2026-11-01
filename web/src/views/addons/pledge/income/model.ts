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
    field: 'orderId',
    component: 'NInputNumber',
    label: '订单ID',
    componentProps: { placeholder: '请输入订单ID', showButton: false },
  },
  {
    field: 'type',
    component: 'NSelect',
    label: '收益类型',
    componentProps: {
      placeholder: '请选择类型',
      options: dict.getOption('PledgeIncomeTypeOptions'),
      clearable: true,
    },
  },
];

export function loadOptions() {
  dict.loadOptions(['PledgeIncomeTypeOptions']);
}

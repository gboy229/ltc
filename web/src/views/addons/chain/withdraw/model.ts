import { ref } from 'vue';
import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  id = 0;
  uid = 0;
  chain = '';
  symbol = '';
  contractAddress = '';
  fromAddress = '';
  toAddress = '';
  payAmount = 0;
  receiveAmount = 0;
  fee = 0;
  rate = 0;
  hash = '';
  status = 0;
  version = 0;
  auditBy = 0;
  auditAt = '';
  remark = '';
  createdAt = '';
  updatedAt = '';

  constructor(state?: Partial<State>) {
    if (state) {
      Object.assign(this, state);
    }
  }
}

export function newState(state: State | Record<string, any> | null): State {
  if (state !== null) {
    if (state instanceof State) {
      return cloneDeep(state);
    }
    return new State(state);
  }
  return new State();
}

// 表格搜索表单：默认全部展开
export const schemas = ref<FormSchema[]>([
  {
    field: 'keyword',
    component: 'NInput',
    label: '用户',
    componentProps: {
      placeholder: '用户ID/钱包地址',
      clearable: true,
    },
  },
  {
    field: 'chain',
    component: 'NInput',
    label: '链',
    componentProps: {
      placeholder: '请输入链名，如 bsc',
      clearable: true,
    },
  },
  {
    field: 'symbol',
    component: 'NInput',
    label: '币种',
    componentProps: {
      placeholder: '请输入币种，如 usdt',
      clearable: true,
    },
  },
  {
    field: 'status',
    component: 'NSelect',
    label: '状态',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择状态',
      options: dict.getOption('ChainWithdrawStatusOptions'),
      clearable: true,
    },
  },
  {
    field: 'toAddress',
    component: 'NInput',
    label: '用户钱包',
    componentProps: {
      placeholder: '请输入用户钱包地址（转入）',
      clearable: true,
    },
  },
  {
    field: 'fromAddress',
    component: 'NInput',
    label: '出币地址',
    componentProps: {
      placeholder: '请输入出币地址（转出）',
      clearable: true,
    },
  },
  {
    field: 'contractAddress',
    component: 'NInput',
    label: '合约地址',
    componentProps: {
      placeholder: '请输入合约地址',
      clearable: true,
    },
  },
  {
    field: 'hash',
    component: 'NInput',
    label: '交易哈希',
    componentProps: {
      placeholder: '请输入交易哈希',
      clearable: true,
    },
  },
  {
    field: 'amountMin',
    component: 'NInputNumber',
    label: '最小金额',
    componentProps: {
      placeholder: '最小申请金额',
      precision: 8,
      min: 0,
      clearable: true,
    },
  },
  {
    field: 'amountMax',
    component: 'NInputNumber',
    label: '最大金额',
    componentProps: {
      placeholder: '最大申请金额',
      precision: 8,
      min: 0,
      clearable: true,
    },
  },
  {
    field: 'createdAt',
    component: 'NDatePicker',
    label: '申请时间',
    componentProps: {
      type: 'datetimerange',
      clearable: true,
      shortcuts: defRangeShortcuts(),
    },
  },
  {
    field: 'auditAt',
    component: 'NDatePicker',
    label: '审核时间',
    componentProps: {
      type: 'datetimerange',
      clearable: true,
      shortcuts: defRangeShortcuts(),
    },
  },
]);

// 加载字典数据选项
export function loadOptions() {
  dict.loadOptions(['ChainWithdrawStatusOptions', 'ChainWithdrawAuditOptions']);
}

import { ref } from 'vue';
import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  id = 0;
  type = '';
  uid = 0;
  chain = '';
  symbol = '';
  fromAddress = '';
  toAddress = '';
  contractAddress = '';
  amount = 0;
  blockHeight = 0;
  hash = '';
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
    field: 'type',
    component: 'NSelect',
    label: '类型',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择充值类型',
      options: dict.getOption('ChainRechargeTypeOptions'),
      clearable: true,
    },
  },
  {
    field: 'fromAddress',
    component: 'NInput',
    label: '用户钱包',
    componentProps: {
      placeholder: '请输入用户钱包地址（转出）',
      clearable: true,
    },
  },
  {
    field: 'toAddress',
    component: 'NInput',
    label: '充值地址',
    componentProps: {
      placeholder: '请输入充值地址（转入）',
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
      placeholder: '最小充值金额',
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
      placeholder: '最大充值金额',
      precision: 8,
      min: 0,
      clearable: true,
    },
  },
  {
    field: 'blockHeightMin',
    component: 'NInputNumber',
    label: '起始块高',
    componentProps: {
      placeholder: '起始区块高度',
      min: 0,
      precision: 0,
      clearable: true,
    },
  },
  {
    field: 'blockHeightMax',
    component: 'NInputNumber',
    label: '结束块高',
    componentProps: {
      placeholder: '结束区块高度',
      min: 0,
      precision: 0,
      clearable: true,
    },
  },
  {
    field: 'createdAt',
    component: 'NDatePicker',
    label: '充值时间',
    componentProps: {
      type: 'datetimerange',
      clearable: true,
      shortcuts: defRangeShortcuts(),
    },
  },
]);

// 加载字典数据选项
export function loadOptions() {
  dict.loadOptions(['ChainRechargeTypeOptions']);
}

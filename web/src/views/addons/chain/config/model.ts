import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  id = 0;
  driver = 'EVM';
  chain = '';
  rpcUrl = '';
  chainId = 0;
  browserUrl = '';
  confirmations = 15;
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

// 表单验证规则
export const rules = {
  driver: {
    required: true,
    trigger: ['blur', 'change'],
    message: '请选择驱动',
  },
  chain: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入网络',
  },
  chainId: {
    required: true,
    trigger: ['blur', 'change'],
    validator: (_rule, value) => {
      return value > 0 ? true : new Error('请输入链ID');
    },
  },
  rpcUrl: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入rpc地址',
  },
  browserUrl: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入浏览器地址',
  },
};

// 表格搜索表单
export const schemas: FormSchema[] = [
  {
    field: 'driver',
    component: 'NSelect',
    label: '驱动',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择驱动',
      options: dict.getOption('ChainDriverOptions'),
      clearable: true,
    },
  },
  {
    field: 'chain',
    component: 'NInput',
    label: '网络',
    componentProps: {
      placeholder: '请输入网络',
    },
  },
];

// 加载字典数据选项
export function loadOptions() {
  dict.loadOptions(['ChainDriverOptions']);
}

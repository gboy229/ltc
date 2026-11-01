import { FormSchema } from '@/components/Form';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  id = 0;
  chain = '';
  address = '';
  symbol = '';
  decimal = 18;
  minAmount = 0;
  maxAmount = 0;
  multipleAmount = 0;
  verifyAmount = 0;
  scale = 8;
  dayLimit = 0;
  feeType = 0;
  rate = 0;
  rechargeSwitch = 1;
  withdrawSwitch = 0;

  constructor(state?: Partial<State>) {
    if (state) {
      Object.assign(this, state);
    }
  }
}

// 只回填表单声明过的字段，避免把 createdAt、扫块游标等系统字段带入提交
export function newState(state: State | Record<string, any> | null): State {
  if (state === null) {
    return new State();
  }

  const target = new State();
  Object.keys(target).forEach((key) => {
    if (state[key] !== undefined && state[key] !== null) {
      target[key] = state[key];
    }
  });
  return target;
}

// 表单验证规则
export const rules = {
  chain: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入链名',
  },
  address: {
    required: true,
    trigger: ['blur', 'input'],
    max: 100,
    message: '请输入合约地址，长度不能超过 100 个字符',
  },
  symbol: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入币种',
  },
  decimal: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    min: 0,
    max: 36,
    message: '精度必须在 0~36 之间',
  },
  scale: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    min: 0,
    max: 36,
    message: '保留小数位必须在 0~36 之间',
  },
  minAmount: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    min: 0,
    message: '最小提币数量不能为负数',
  },
  maxAmount: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    min: 0,
    message: '最大提币数量不能为负数',
  },
  multipleAmount: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    min: 0,
    message: '提币倍数不能为负数',
  },
  verifyAmount: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    min: 0,
    message: '审核阈值不能为负数',
  },
  dayLimit: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    min: 0,
    message: '日上限笔数不能为负数',
  },
  feeType: {
    required: true,
    trigger: ['blur', 'change'],
    message: '请选择手续费类型',
  },
  rate: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    min: 0,
    message: '手续费费率不能为负数',
  },
  rechargeSwitch: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    message: '请选择充币开关',
  },
  withdrawSwitch: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    message: '请选择提现开关',
  },
};

// 表格搜索表单
export const schemas: FormSchema[] = [
  {
    field: 'chain',
    component: 'NInput',
    label: '链名',
    componentProps: {
      placeholder: '请输入链名',
    },
  },
  {
    field: 'symbol',
    component: 'NInput',
    label: '币种',
    componentProps: {
      placeholder: '请输入币种',
    },
  },
  {
    field: 'address',
    component: 'NInput',
    label: '合约地址',
    componentProps: {
      placeholder: '请输入合约地址',
    },
  },
];

// 加载字典数据选项
export function loadOptions() {
  dict.loadOptions(['ChainFeeTypeOptions', 'ChainSwitchOptions']);
}

import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  id = 0;
  name = '';
  personalAmount = 0;
  teamAmount = 0;
  diffReward = 0;
  status = 1;
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
  name: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入等级名称',
  },
  personalAmount: {
    required: true,
    trigger: ['blur', 'change'],
    validator: (_rule, value) => {
      return typeof value === 'number' && value >= 0 ? true : new Error('请输入个人业绩');
    },
  },
  teamAmount: {
    required: true,
    trigger: ['blur', 'change'],
    validator: (_rule, value) => {
      return typeof value === 'number' && value >= 0 ? true : new Error('请输入团队业绩');
    },
  },
  diffReward: {
    required: true,
    trigger: ['blur', 'change'],
    validator: (_rule, value) => {
      return typeof value === 'number' && value >= 0 ? true : new Error('请输入级差奖励');
    },
  },
  status: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    message: '请选择状态',
  },
};

// 表格搜索表单
export const schemas: FormSchema[] = [
  {
    field: 'name',
    component: 'NInput',
    label: '等级名称',
    componentProps: {
      placeholder: '请输入等级名称',
    },
  },
  {
    field: 'status',
    component: 'NSelect',
    label: '状态',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择状态',
      options: dict.getOption('sys_normal_disable'),
    },
  },
];

// 加载字典数据选项
export function loadOptions() {
  dict.loadOptions(['sys_normal_disable']);
}

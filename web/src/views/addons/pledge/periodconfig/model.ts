import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  id = 0;
  periodMonths = 3;
  monthlyRewardInfo = '';
  monthlyRewardRate = 0;
  directRewardRate = 0;
  status = 1;
  createdAt = '';
  updatedAt = '';

  constructor(state?: Partial<State>) {
    if (state) Object.assign(this, state);
  }
}

export function newState(state: State | Record<string, any> | null): State {
  if (state !== null) {
    if (state instanceof State) return cloneDeep(state);
    return new State(state);
  }
  return new State();
}

export const rules = {
  periodMonths: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    message: '请选择投资周期',
  },
  monthlyRewardRate: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    message: '请输入月化收益率',
  },
  directRewardRate: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    message: '请输入直推收益率',
  },
};

export const schemas: FormSchema[] = [
  {
    field: 'periodMonths',
    component: 'NSelect',
    label: '投资周期',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择周期',
      options: dict.getOption('PledgePeriodOptions'),
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
      options: dict.getOption('sys_normal_disable'),
      clearable: true,
    },
  },
];

export function loadOptions() {
  dict.loadOptions(['sys_normal_disable', 'PledgePeriodOptions']);
}

import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { useDictStore } from '@/store/modules/dict';

const dict = useDictStore();

export class State {
  id = 0;
  title = '';
  type = 1;
  cover = '';
  content = '';
  sort = 0;
  status = 1;
  isRecommend = 1;
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
  title: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入公告标题',
  },
  content: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入公告内容',
  },
};

// 表格搜索表单
export const schemas: FormSchema[] = [
  {
    field: 'title',
    component: 'NInput',
    label: '公告标题',
    componentProps: {
      placeholder: '请输入公告标题',
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
  {
    field: 'isRecommend',
    component: 'NSelect',
    label: '推荐首页',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择是否推荐',
      options: dict.getOption('CmsArticleRecommendOptions'),
    },
  },
];

// 加载字典数据选项
export function loadOptions() {
  dict.loadOptions(['sys_normal_disable', 'CmsArticleRecommendOptions']);
}

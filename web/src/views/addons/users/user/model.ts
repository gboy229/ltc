import { computed } from 'vue';
import { cloneDeep } from 'lodash-es';
import { FormSchema } from '@/components/Form';
import { useDictStore } from '@/store/modules/dict';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { Option } from '@/api/addons/users/level';

const dict = useDictStore();

// 认证状态搜索选项（含全部）
const authStatusSearchOptions = computed(() => [
  { label: '全部', value: -1 },
  ...(dict.getOptionUnRef('UserAuthStatusOptions') || []),
]);

export class State {
  id = 0;
  uid = '';
  pid = 0;
  levelId = 0;
  lockLevel = 0;
  nickname = '';
  avatar = '';
  authStatus = 0;
  status = 1;

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
  nickname: {
    required: true,
    trigger: ['blur', 'input'],
    message: '请输入昵称',
  },
  levelId: {
    required: true,
    trigger: ['blur', 'change'],
    validator: (_rule, value) => {
      return value > 0 ? true : new Error('请选择会员等级');
    },
  },
  status: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    message: '请选择账号状态',
  },
  lockLevel: {
    type: 'number',
    required: true,
    trigger: ['blur', 'change'],
    message: '请选择是否锁定等级',
  },
};

// 表格搜索表单
export const schemas: FormSchema[] = [
  {
    field: 'uid',
    component: 'NInput',
    label: '用户号',
    componentProps: {
      placeholder: '请输入用户号',
    },
  },
  {
    field: 'nickname',
    component: 'NInput',
    label: '昵称',
    componentProps: {
      placeholder: '请输入昵称',
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
  {
    field: 'parentWallet',
    component: 'NInput',
    label: '推荐钱包',
    componentProps: {
      placeholder: '请输入推荐人钱包地址',
    },
  },
  {
    field: 'teamWallet',
    component: 'NInput',
    label: '团队钱包',
    componentProps: {
      placeholder: '请输入团队成员钱包地址',
    },
  },
  {
    field: 'levelId',
    component: 'NSelect',
    label: '会员等级',
    defaultValue: null,
    componentProps: {
      placeholder: '请选择等级',
      options: dict.getOption('usersLevel'),
    },
  },
  {
    field: 'authStatus',
    component: 'NSelect',
    label: '认证状态',
    defaultValue: -1,
    componentProps: {
      placeholder: '请选择认证状态',
      options: authStatusSearchOptions,
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
    field: 'created_at',
    component: 'NDatePicker',
    label: '注册时间',
    componentProps: {
      type: 'datetimerange',
      clearable: true,
      shortcuts: defRangeShortcuts(),
    },
  },
];

// 加载字典数据选项
export function loadOptions() {
  dict.loadOptions([
    'sys_normal_disable',
    'UserAuthTypeOptions',
    'UserAuthStatusOptions',
    'UserVerifiedOptions',
    'UserLockLevelOptions',
  ]);

  Option().then((res) => {
    const opts = (res?.list || []).map((i) => dict.genOption(i.id, i.name));
    dict.setOption('usersLevel', opts);
  });
}

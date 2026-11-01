import { h } from 'vue';
import { NAvatar, NButton, NTag, NText, NTooltip } from 'naive-ui';
import { FormSchema } from '@/components/Form';
import { useDictStore } from '@/store/modules/dict';
import { defRangeShortcuts } from '@/utils/dateUtil';
import { renderOptionTag } from '@/utils';

const dict = useDictStore();

// 金额格式化：最多保留8位小数，去掉尾部多余0，至少保留2位
function money(value: any): string {
  const n = Number(value || 0);
  let s = n.toFixed(8).replace(/0+$/, '').replace(/\.$/, '');
  const parts = s.split('.');
  if (parts.length === 1) {
    s += '.00';
  } else if (parts[1].length === 1) {
    s += '0';
  }
  return s;
}

// 带符号金额：正数加 + 号
function signedMoney(value: any): string {
  const n = Number(value || 0);
  return n > 0 ? `+${money(value)}` : money(value);
}

// 复制文本
function copyText(text: string) {
  if (!text) return;
  const $message = window['$message'];
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(text).then(() => $message?.success('复制成功'));
    return;
  }
  const input = document.createElement('input');
  input.value = text;
  document.body.appendChild(input);
  input.select();
  document.execCommand('copy');
  document.body.removeChild(input);
  $message?.success('复制成功');
}

// 字典标签，未知值回退展示原值
function optionTagOrRaw(type: string, value: any) {
  if (value === null || value === undefined || value === '') return '—';
  const label = dict.getLabel(type, value);
  if (label) return renderOptionTag(type, value);
  return h(
    NTag,
    { size: 'small', bordered: false, type: 'default' },
    { default: () => String(value) }
  );
}

// 带说明的一行：灰色标签 + 内容
function tagLine(label: string, node: any) {
  return h(
    'div',
    { style: { display: 'flex', alignItems: 'center', gap: '4px', lineHeight: '22px' } },
    [h('span', { style: { color: '#8c8c8c', fontSize: '12px', flexShrink: 0 } }, label), node]
  );
}

// 变动金额：正绿负红，带符号
function changeText(value: any) {
  const n = Number(value || 0);
  return h(
    NText,
    { type: n > 0 ? 'success' : n < 0 ? 'error' : 'default', strong: true },
    { default: () => signedMoney(value) }
  );
}

// 余额变动单元格：变动值 + 变动前→变动后
function balanceChangeCell(changeValue: any, before: any, after: any) {
  return h('div', { style: { lineHeight: '20px' } }, [
    h('div', [changeText(changeValue)]),
    h(
      'div',
      { style: { color: '#8c8c8c', fontSize: '12px' } },
      `${money(before)} → ${money(after)}`
    ),
  ]);
}

// 可复制文本（超出省略，悬浮看完整）
function copyable(text: string, maxWidth = '150px') {
  if (!text) return '—';
  return h('div', { style: { display: 'flex', alignItems: 'center', gap: '4px' } }, [
    h(
      NTooltip,
      { trigger: 'hover' },
      {
        trigger: () =>
          h(
            'span',
            {
              style: {
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
                maxWidth,
                fontSize: '12px',
              },
            },
            text
          ),
        default: () => text,
      }
    ),
    h(
      NButton,
      { size: 'tiny', text: true, type: 'primary', onClick: () => copyText(text) },
      { default: () => '复制' }
    ),
  ]);
}

export const columns = [
  {
    title: '用户',
    key: 'userInfo',
    width: 200,
    render(row) {
      return h('div', { style: { display: 'flex', alignItems: 'center', gap: '8px' } }, [
        h(
          NAvatar,
          { circle: true, size: 'small' },
          { default: () => (row.nickname || 'U').substring(0, 1) }
        ),
        h('div', { style: { minWidth: 0 } }, [
          h('div', { style: { fontWeight: 600, lineHeight: '20px' } }, row.nickname || '—'),
          h('div', { style: { color: '#8c8c8c', lineHeight: '18px' } }, `UID：${row.uid || '—'}`),
        ]),
      ]);
    },
  },
  {
    title: '钱包地址',
    key: 'walletAddress',
    width: 200,
    render(row) {
      const wallet = row.walletAddress;
      if (!wallet) return '—';
      return h('div', { style: { display: 'flex', alignItems: 'center', gap: '4px' } }, [
        h(
          'span',
          { style: { wordBreak: 'break-all', whiteSpace: 'normal', lineHeight: '1.4' } },
          wallet
        ),
        h(
          NButton,
          { size: 'tiny', text: true, type: 'primary', onClick: () => copyText(wallet) },
          { default: () => '复制' }
        ),
      ]);
    },
  },
  {
    title: '账务类型',
    key: 'changeType',
    width: 150,
    render(row) {
      return h('div', [
        tagLine('动作', renderOptionTag('AssetChangeTypeOptions', row.changeType)),
        tagLine('业务', optionTagOrRaw('AssetBizTypeOptions', row.bizType)),
      ]);
    },
  },
  {
    title: '币种',
    key: 'account',
    width: 90,
    render(row) {
      return h(
        NTag,
        { size: 'small', bordered: false, type: 'info' },
        { default: () => row.account || '—' }
      );
    },
  },
  {
    title: '可用余额变动',
    key: 'available',
    width: 170,
    render(row) {
      return balanceChangeCell(row.changeAvailable, row.beforeAvailable, row.afterAvailable);
    },
  },
  {
    title: '冻结余额变动',
    key: 'frozen',
    width: 150,
    render(row) {
      if (Number(row.changeFrozen || 0) === 0) return '—';
      return balanceChangeCell(row.changeFrozen, row.beforeFrozen, row.afterFrozen);
    },
  },
  {
    title: '业务单号',
    key: 'bizSn',
    width: 180,
    render(row) {
      const sn = row.bizSn;
      const bizId = row.bizId;
      if (!sn && !bizId) return '—';
      return h('div', { style: { lineHeight: '20px' } }, [
        sn ? h('div', [copyable(sn, '120px')]) : null,
        bizId
          ? h('div', { style: { color: '#8c8c8c', fontSize: '12px' } }, `业务ID：${bizId}`)
          : null,
      ]);
    },
  },
  {
    title: '备注',
    key: 'remark',
    width: 150,
    render(row) {
      return row.remark || '—';
    },
  },
  {
    title: '发生时间',
    key: 'createdAt',
    width: 170,
  },
];

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
    field: 'walletAddress',
    component: 'NInput',
    label: '钱包地址',
    componentProps: {
      placeholder: '请输入钱包地址',
    },
  },
  {
    field: 'changeType',
    component: 'NSelect',
    label: '动作',
    componentProps: {
      placeholder: '请选择动作',
      options: dict.getOption('AssetChangeTypeOptions'),
      clearable: true,
    },
  },
  {
    field: 'bizType',
    component: 'NSelect',
    label: '业务类型',
    componentProps: {
      placeholder: '请选择业务类型',
      options: dict.getOption('AssetBizTypeOptions'),
      clearable: true,
    },
  },
  {
    field: 'createdAt',
    component: 'NDatePicker',
    label: '发生时间',
    componentProps: {
      type: 'datetimerange',
      clearable: true,
      shortcuts: defRangeShortcuts(),
    },
  },
  {
    field: 'account',
    component: 'NInput',
    label: '资产账户',
    componentProps: {
      placeholder: '例如 USDT',
    },
  },
  {
    field: 'direction',
    component: 'NSelect',
    label: '资金方向',
    componentProps: {
      placeholder: '请选择方向',
      clearable: true,
      options: [
        { label: '收入', value: 'income' },
        { label: '支出', value: 'expense' },
      ],
    },
  },
];

export function loadOptions() {
  dict.loadOptions(['AssetChangeTypeOptions', 'AssetBizTypeOptions']);
}

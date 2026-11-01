import { h } from 'vue';
import { NAvatar, NButton } from 'naive-ui';
import { renderOptionTag } from '@/utils';

function copyWallet(wallet: string) {
  if (!wallet) return;
  const $message = window['$message'];
  if (navigator.clipboard && window.isSecureContext) {
    navigator.clipboard.writeText(wallet).then(() => $message?.success('复制成功'));
    return;
  }
  const input = document.createElement('input');
  input.value = wallet;
  document.body.appendChild(input);
  input.select();
  document.execCommand('copy');
  document.body.removeChild(input);
  $message?.success('复制成功');
}

function money(value: number | undefined) {
  return Number(value || 0).toFixed(4);
}

function line(label: string, value: any) {
  return h('div', { style: { lineHeight: '20px', whiteSpace: 'nowrap' } }, [
    h('span', { style: { color: '#8c8c8c' } }, label),
    h('span', { style: { color: '#1f1a13', marginLeft: '2px' } }, value),
  ]);
}

export const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 70,
  },
  {
    title: '用户信息',
    key: 'userInfo',
    width: 230,
    render(row) {
      return h('div', { style: { display: 'flex', alignItems: 'center', gap: '8px' } }, [
        row.avatar
          ? h(NAvatar, { circle: true, size: 'small', src: row.avatar })
          : h(NAvatar, { circle: true, size: 'small' }, { default: () => (row.nickname || 'U').substring(0, 1) }),
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
    width: 220,
    render(row) {
      const wallet = row.walletAddress;
      if (!wallet) return '—';
      return h('div', { style: { display: 'flex', alignItems: 'center', gap: '4px' } }, [
        h('span', { style: { wordBreak: 'break-all', whiteSpace: 'normal', lineHeight: '1.4' } }, wallet),
        h(
          NButton,
          { size: 'tiny', text: true, type: 'primary', onClick: () => copyWallet(wallet) },
          { default: () => '复制' }
        ),
      ]);
    },
  },
  {
    title: '关系与等级',
    key: 'relation',
    width: 190,
    render(row) {
      const parent = row.pid > 0 ? row.parentNickname || row.parentUid || '#' + row.pid : '—';
      return h('div', [
        line('等级：', row.levelName || '—'),
        line('推荐UID：', row.parentUid || '—'),
        line('推荐人：', parent),
        line('锁定等级：', row.lockLevel === 1 ? '是' : '否'),
      ]);
    },
  },
  {
    title: '业绩',
    key: 'performance',
    width: 230,
    render(row) {
      return h('div', [
        line('个人：', money(row.personalAmount)),
        line('直推：', money(row.directAmount)),
        line('团队：', money(row.teamAmount)),
      ]);
    },
  },
  {
    title: '状态',
    key: 'statusInfo',
    width: 130,
    render(row) {
      return h('div', { style: { display: 'flex', flexDirection: 'column', gap: '2px' } }, [
        h('div', [h('span', { style: { color: '#8c8c8c', marginRight: '4px' } }, '实名'), renderOptionTag('UserAuthStatusOptions', row.authStatus)]),
        h('div', [h('span', { style: { color: '#8c8c8c', marginRight: '4px' } }, '状态'), renderOptionTag('sys_normal_disable', row.status)]),
      ]);
    },
  },
  {
    title: '登录信息',
    key: 'loginInfo',
    width: 220,
    render(row) {
      return h('div', [
        line('IP：', row.lastLoginIp || '—'),
        line('时间：', row.lastLoginTime || '—'),
      ]);
    },
  },
  {
    title: '注册时间',
    key: 'createdAt',
    width: 180,
  },
];

<template>
  <div class="console">
    <!-- 顶部数据卡片 -->
    <n-grid cols="1 s:2 m:2 l:4 xl:4 2xl:4" responsive="screen" :x-gap="12" :y-gap="12">
      <n-gi v-for="card in cards" :key="card.label">
        <n-card
          :bordered="false"
          size="small"
          class="stat-card"
          :class="{ 'stat-card--link': card.click }"
          @click="card.click && router.push(card.click)"
        >
          <n-statistic :label="card.label" :value="card.value">
            <template v-if="card.suffix" #suffix>
              <span class="stat-suffix">{{ card.suffix }}</span>
            </template>
          </n-statistic>
          <div class="stat-sub">{{ card.sub }}</div>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 趋势 -->
    <n-card :bordered="false" class="mt-4">
      <template #header>
        <span>数据趋势</span>
      </template>
      <template #header-extra>
        <n-radio-group v-model:value="rangeDays" size="small" @update:value="loadTrend">
          <n-radio-button :value="7">近7天</n-radio-button>
          <n-radio-button :value="30">近30天</n-radio-button>
        </n-radio-group>
      </template>
      <div ref="chartRef" style="height: 340px"></div>
    </n-card>

    <!-- 快捷入口 -->
    <n-card :bordered="false" class="mt-4" title="快捷入口">
      <n-grid cols="2 s:4 m:4 l:6" responsive="screen" :x-gap="12" :y-gap="12">
        <n-gi v-for="item in shortcuts" :key="item.title">
          <n-card
            size="small"
            hoverable
            :bordered="false"
            class="shortcut-card"
            @click="router.push(item.path)"
          >
            <div class="shortcut-inner">
              <n-icon :size="28" :color="item.color">
                <component :is="item.icon" />
              </n-icon>
              <span class="shortcut-title">{{ item.title }}</span>
            </div>
          </n-card>
        </n-gi>
      </n-grid>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue';
  import type { Ref } from 'vue';
  import { useRouter } from 'vue-router';
  import { format } from 'date-fns';
  import {
    UsergroupAddOutlined,
    AccountBookOutlined,
    WalletOutlined,
    SendOutlined,
    BarChartOutlined,
    SettingOutlined,
    SwapOutlined,
    PieChartOutlined,
  } from '@vicons/antd';
  import { List as UserList } from '@/api/addons/users/user';
  import { Stat as AssetStat } from '@/api/addons/users/asset';
  import { Stat as PledgeStat } from '@/api/addons/pledge/stat';
  import { Stat as ChainStat } from '@/api/addons/chain/stat';
  import { List as WithdrawList } from '@/api/addons/chain/withdraw';
  import { useECharts } from '@/hooks/web/useECharts';

  const router = useRouter();
  const chartRef = ref<HTMLDivElement | null>(null);
  const { setOptions } = useECharts(chartRef as Ref<HTMLDivElement>);

  const rangeDays = ref(7);

  const overview = reactive({
    userTotal: 0,
    userToday: 0,
    orderCount: 0,
    orderAmount: 0,
    incomeAmount: 0,
    incomeCount: 0,
    rechargeAmount: 0,
    rechargeCount: 0,
    withdrawReceive: 0,
    withdrawFee: 0,
    withdrawCount: 0,
    totalAvailable: 0,
    totalFrozen: 0,
    pendingWithdraw: 0,
  });

  const cards = computed(() => [
    {
      label: '用户总数',
      value: overview.userTotal,
      suffix: '人',
      sub: `今日新增 ${overview.userToday}`,
    },
    {
      label: '今日质押订单',
      value: overview.orderCount,
      suffix: '笔',
      sub: `订单金额 ${money(overview.orderAmount)}`,
    },
    {
      label: '今日充值',
      value: money(overview.rechargeAmount),
      suffix: 'USDT',
      sub: `${overview.rechargeCount} 笔`,
    },
    {
      label: '今日提现',
      value: money(overview.withdrawReceive),
      suffix: 'USDT',
      sub: `手续费 ${money(overview.withdrawFee)}`,
    },
    {
      label: '今日收益产出',
      value: money(overview.incomeAmount),
      suffix: 'USDT',
      sub: `${overview.incomeCount} 笔`,
    },
    {
      label: '资产总额',
      value: money(overview.totalAvailable),
      suffix: 'USDT',
      sub: `冻结 ${money(overview.totalFrozen)}`,
    },
    {
      label: '待审核提现',
      value: overview.pendingWithdraw,
      suffix: '笔',
      sub: '点击前往审核',
      click: '/chain/withdraw',
    },
  ]);

  const shortcuts = [
    { title: '用户管理', path: '/users/user', icon: UsergroupAddOutlined, color: '#69c0ff' },
    { title: '订单管理', path: '/pledge/order', icon: AccountBookOutlined, color: '#b37feb' },
    { title: '收益流水', path: '/pledge/income', icon: PieChartOutlined, color: '#5cdbd3' },
    { title: '充值记录', path: '/chain/recharge', icon: WalletOutlined, color: '#18a058' },
    { title: '提现审核', path: '/chain/withdraw', icon: SendOutlined, color: '#f0a020' },
    { title: '质押统计', path: '/pledge/stat', icon: BarChartOutlined, color: '#2080f0' },
    { title: '资产统计', path: '/users/asset/stat', icon: SwapOutlined, color: '#ff9c6e' },
    { title: '网络配置', path: '/chain/config', icon: SettingOutlined, color: '#ffc069' },
  ];

  function money(val: any) {
    return Number(val || 0).toFixed(2);
  }

  function toStr(ts: number) {
    return format(new Date(ts), 'yyyy-MM-dd HH:mm:ss');
  }

  function todayRange(): [number, number] {
    const now = new Date();
    const start = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime();
    return [start, now.getTime()];
  }

  function rangeToTs(days: number): [number, number] {
    const now = new Date();
    const start =
      new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime() - (days - 1) * 86400000;
    return [start, now.getTime()];
  }

  async function loadOverview() {
    const [start, end] = todayRange();
    const [userTotalRes, userTodayRes, assetRes, pledgeRes, chainRes, pendingRes] =
      await Promise.all([
        UserList({ page: 1, pageSize: 1 }),
        UserList({ page: 1, pageSize: 1, createdAt: [start, end] }),
        AssetStat(),
        PledgeStat({ startAt: toStr(start), endAt: toStr(end) }),
        ChainStat({ startAt: toStr(start), endAt: toStr(end) }),
        WithdrawList({ status: 0, page: 1, pageSize: 5 }),
      ]);

    overview.userTotal = userTotalRes?.totalCount || 0;
    overview.userToday = userTodayRes?.totalCount || 0;

    const order = pledgeRes?.order || {};
    overview.orderCount = order.count || 0;
    overview.orderAmount = order.amount || 0;

    const income = pledgeRes?.income || {};
    overview.incomeAmount = income.amount || 0;
    overview.incomeCount = income.count || 0;

    const recharge = chainRes?.recharge || {};
    overview.rechargeAmount = recharge.totalAmount || 0;
    overview.rechargeCount = recharge.totalCount || 0;

    const withdraw = chainRes?.withdraw || {};
    overview.withdrawReceive = withdraw.receiveAmount || 0;
    overview.withdrawFee = withdraw.fee || 0;
    overview.withdrawCount = withdraw.totalCount || 0;

    overview.totalAvailable = assetRes?.totalAvailable || 0;
    overview.totalFrozen = assetRes?.totalFrozen || 0;
    overview.pendingWithdraw = pendingRes?.totalCount || 0;
  }

  async function loadTrend() {
    const [start, end] = rangeToTs(rangeDays.value);
    const [chainRes, pledgeRes] = await Promise.all([
      ChainStat({ startAt: toStr(start), endAt: toStr(end) }),
      PledgeStat({ startAt: toStr(start), endAt: toStr(end) }),
    ]);
    renderTrend(chainRes?.trend || [], pledgeRes?.trend || []);
  }

  function renderTrend(chainTrend: any[], pledgeTrend: any[]) {
    const dates = new Set<string>();
    chainTrend.forEach((d) => dates.add(d.date));
    pledgeTrend.forEach((d) => dates.add(d.date));
    const dateArr = Array.from(dates).sort();

    const chainMap = new Map(chainTrend.map((d) => [d.date, d]));
    const pledgeMap = new Map(pledgeTrend.map((d) => [d.date, d]));

    const rechargeSeries = dateArr.map((d) => Number(chainMap.get(d)?.rechargeAmount || 0));
    const withdrawSeries = dateArr.map((d) => Number(chainMap.get(d)?.withdrawAmount || 0));
    const incomeSeries = dateArr.map((d) => Number(pledgeMap.get(d)?.incomeAmount || 0));

    setOptions({
      tooltip: { trigger: 'axis' },
      legend: { data: ['充值金额', '提现金额', '收益产出'] },
      grid: { left: '3%', right: '3%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: dateArr },
      yAxis: { type: 'value' },
      series: [
        {
          name: '充值金额',
          type: 'line',
          smooth: true,
          data: rechargeSeries,
          itemStyle: { color: '#18a058' },
        },
        {
          name: '提现金额',
          type: 'line',
          smooth: true,
          data: withdrawSeries,
          itemStyle: { color: '#f0a020' },
        },
        {
          name: '收益产出',
          type: 'line',
          smooth: true,
          data: incomeSeries,
          itemStyle: { color: '#2080f0' },
        },
      ],
    });
  }

  onMounted(async () => {
    await loadOverview();
    await loadTrend();
  });
</script>

<style lang="less" scoped>
  .stat-card {
    .stat-suffix {
      margin-left: 4px;
      font-size: 12px;
      color: #909399;
    }

    .stat-sub {
      margin-top: 6px;
      font-size: 12px;
      color: #909399;
    }

    &--link {
      cursor: pointer;
    }
  }

  .shortcut-card {
    .shortcut-inner {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8px;
    }

    .shortcut-title {
      font-size: 13px;
      color: #606266;
    }
  }
</style>

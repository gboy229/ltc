<template>
  <div>
    <n-card :bordered="false" class="proCard" title="数据统计">
      <!-- 筛选区 -->
      <n-form inline :show-feedback="false" label-placement="left" label-width="auto">
        <n-form-item label="日期">
          <n-date-picker
            v-model:value="query.date"
            type="datetimerange"
            clearable
            :shortcuts="defRangeShortcuts()"
            style="width: 340px"
          />
        </n-form-item>
        <n-form-item label="链">
          <n-input
            v-model:value="query.chain"
            placeholder="如 bsc"
            clearable
            style="width: 160px"
          />
        </n-form-item>
        <n-form-item label="币种">
          <n-input
            v-model:value="query.symbol"
            placeholder="如 usdt"
            clearable
            style="width: 160px"
          />
        </n-form-item>
        <n-form-item>
          <n-button type="primary" @click="loadData">查询</n-button>
        </n-form-item>
      </n-form>

      <!-- 充值统计（收入） -->
      <n-card :bordered="false" class="mt-4 stat-block stat-block--income">
        <template #header>
          <span class="stat-title stat-title--income">充值统计</span>
        </template>
        <n-grid :cols="2" :x-gap="16" responsive="screen">
          <n-gi>
            <n-statistic
              label="充值笔数"
              :value="recharge.totalCount"
              :value-style="incomeValueStyle"
            />
          </n-gi>
          <n-gi>
            <n-statistic
              label="充值总额"
              :value="formatAmount(recharge.totalAmount)"
              :value-style="incomeValueStyle"
            />
          </n-gi>
        </n-grid>
      </n-card>

      <!-- 提现统计（支出） -->
      <n-card :bordered="false" class="mt-4 stat-block stat-block--expense">
        <template #header>
          <span class="stat-title stat-title--expense">提现统计</span>
        </template>
        <n-grid :cols="4" :x-gap="16" responsive="screen">
          <n-gi>
            <n-statistic
              label="提现笔数"
              :value="withdraw.totalCount"
              :value-style="expenseValueStyle"
            />
          </n-gi>
          <n-gi>
            <n-statistic
              label="申请总额"
              :value="formatAmount(withdraw.payAmount)"
              :value-style="expenseValueStyle"
            />
          </n-gi>
          <n-gi>
            <n-statistic
              label="到账总额"
              :value="formatAmount(withdraw.receiveAmount)"
              :value-style="expenseValueStyle"
            />
          </n-gi>
          <n-gi>
            <n-statistic
              label="手续费收入"
              :value="formatAmount(withdraw.fee)"
              :value-style="expenseValueStyle"
            />
          </n-gi>
        </n-grid>
      </n-card>

      <!-- 按天趋势 -->
      <n-card title="按天趋势（充值额 vs 提现到账额）" :bordered="false" class="mt-4">
        <div ref="chartRef" style="height: 360px"></div>
      </n-card>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, reactive, ref } from 'vue';
  import type { Ref } from 'vue';
  import { format } from 'date-fns';
  import { Stat } from '@/api/addons/chain/stat';
  import { defRangeShortcuts } from '@/utils/dateUtil';
  import { useECharts } from '@/hooks/web/useECharts';

  const chartRef = ref<HTMLDivElement | null>(null);
  const { setOptions } = useECharts(chartRef as Ref<HTMLDivElement>);

  // 收入/支出配色
  const incomeValueStyle = { color: '#18a058' };
  const expenseValueStyle = { color: '#f0a020' };

  const query = reactive({
    date: defaultRange(),
    chain: '',
    symbol: '',
  });

  const recharge = reactive({
    totalCount: 0,
    totalAmount: 0,
    list: [] as any[],
  });

  const withdraw = reactive({
    totalCount: 0,
    payAmount: 0,
    receiveAmount: 0,
    fee: 0,
    list: [] as any[],
  });

  function defaultRange(): [number, number] {
    const end = new Date();
    const start =
      new Date(end.getFullYear(), end.getMonth(), end.getDate()).getTime() - 6 * 86400000;
    return [start, end.getTime()];
  }

  function formatAmount(val: any) {
    return Number(val || 0).toFixed(2);
  }

  function toStr(ts: number) {
    return format(new Date(ts), 'yyyy-MM-dd HH:mm:ss');
  }

  function loadData() {
    const [start, end] = query.date || [];
    Stat({
      startAt: start ? toStr(start) : undefined,
      endAt: end ? toStr(end) : undefined,
      chain: query.chain || undefined,
      symbol: query.symbol || undefined,
    }).then((res) => {
      Object.assign(recharge, res.recharge || {});
      Object.assign(withdraw, res.withdraw || {});
      renderChart();
    });
  }

  function renderChart() {
    const dates = new Set<string>();
    (recharge.list || []).forEach((d) => dates.add(d.date));
    (withdraw.list || []).forEach((d) => dates.add(d.date));
    const dateArr = Array.from(dates).sort();

    const rechargeMap = new Map((recharge.list || []).map((d) => [d.date, d.amount]));
    const withdrawMap = new Map((withdraw.list || []).map((d) => [d.date, d.receiveAmount]));
    const rechargeSeries = dateArr.map((d) => rechargeMap.get(d) || 0);
    const withdrawSeries = dateArr.map((d) => withdrawMap.get(d) || 0);

    setOptions({
      tooltip: { trigger: 'axis' },
      legend: { data: ['充值额', '提现到账额'] },
      grid: { left: '3%', right: '3%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: dateArr },
      yAxis: { type: 'value' },
      series: [
        {
          name: '充值额',
          type: 'bar',
          data: rechargeSeries,
          barMaxWidth: 30,
          itemStyle: { color: '#18a058' },
        },
        {
          name: '提现到账额',
          type: 'bar',
          data: withdrawSeries,
          barMaxWidth: 30,
          itemStyle: { color: '#f0a020' },
        },
      ],
    });
  }

  onMounted(() => {
    loadData();
  });
</script>

<style lang="less" scoped>
  .mt-4 {
    margin-top: 16px;
  }

  .stat-block {
    border-left: 3px solid transparent;
  }

  .stat-block--income {
    border-left-color: #18a058;
    background-color: #f7fdf9;
  }

  .stat-block--expense {
    border-left-color: #f0a020;
    background-color: #fdfaf3;
  }

  .stat-title {
    font-weight: 600;
  }

  .stat-title--income {
    color: #18a058;
  }

  .stat-title--expense {
    color: #f0a020;
  }
</style>

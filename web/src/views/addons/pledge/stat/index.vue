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
        <n-form-item>
          <n-button type="primary" @click="loadData">查询</n-button>
        </n-form-item>
      </n-form>

      <!-- 订单统计（入金） -->
      <n-card :bordered="false" class="mt-4 stat-block stat-block--order">
        <template #header>
          <span class="stat-title stat-title--order">订单统计（入金）</span>
        </template>
        <n-grid :cols="4" :x-gap="16" responsive="screen">
          <n-gi>
            <n-statistic label="订单数" :value="order.count" :value-style="orderValueStyle" />
          </n-gi>
          <n-gi>
            <n-statistic
              label="订单金额"
              :value="formatAmount(order.amount)"
              :value-style="orderValueStyle"
            />
          </n-gi>
          <n-gi>
            <n-statistic label="进行中" :value="order.activeCount" :value-style="orderValueStyle" />
          </n-gi>
          <n-gi>
            <n-statistic
              label="已完成"
              :value="order.finishedCount"
              :value-style="orderValueStyle"
            />
          </n-gi>
        </n-grid>
      </n-card>

      <!-- 收益统计（产出） -->
      <n-card :bordered="false" class="mt-4 stat-block stat-block--income">
        <template #header>
          <span class="stat-title stat-title--income">收益统计（产出）</span>
        </template>
        <n-grid :cols="2" :x-gap="16" responsive="screen">
          <n-gi>
            <n-statistic
              label="产出金额"
              :value="formatAmount(income.amount, 4)"
              :value-style="incomeValueStyle"
            />
          </n-gi>
          <n-gi>
            <n-statistic label="产出笔数" :value="income.count" :value-style="incomeValueStyle" />
          </n-gi>
        </n-grid>

        <n-table :bordered="false" :single-line="false" size="small" class="mt-3">
          <thead>
            <tr>
              <th>收益类型</th>
              <th>笔数</th>
              <th>金额</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in income.byType" :key="item.type">
              <td>{{ item.label }}</td>
              <td>{{ item.count }}</td>
              <td>{{ formatAmount(item.amount, 4) }}</td>
            </tr>
          </tbody>
        </n-table>
      </n-card>

      <!-- 按天趋势 -->
      <n-card title="按天趋势（订单金额 vs 产出金额）" :bordered="false" class="mt-4">
        <div ref="chartRef" style="height: 360px"></div>
      </n-card>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, reactive, ref } from 'vue';
  import type { Ref } from 'vue';
  import { format } from 'date-fns';
  import { Stat } from '@/api/addons/pledge/stat';
  import { defRangeShortcuts } from '@/utils/dateUtil';
  import { useECharts } from '@/hooks/web/useECharts';

  const chartRef = ref<HTMLDivElement | null>(null);
  const { setOptions } = useECharts(chartRef as Ref<HTMLDivElement>);

  // 入金/产出配色
  const orderValueStyle = { color: '#18a058' };
  const incomeValueStyle = { color: '#f0a020' };

  const query = reactive({
    date: null as [number, number] | null,
  });

  const order = reactive({
    count: 0,
    amount: 0,
    activeCount: 0,
    finishedCount: 0,
  });

  const income = reactive({
    count: 0,
    amount: 0,
    byType: [] as any[],
  });

  const trend = ref<any[]>([]);

  function formatAmount(val: any, digits = 2) {
    return Number(val || 0).toFixed(digits);
  }

  function toStr(ts: number) {
    return format(new Date(ts), 'yyyy-MM-dd HH:mm:ss');
  }

  function loadData() {
    const [start, end] = query.date || [];
    Stat({
      startAt: start ? toStr(start) : undefined,
      endAt: end ? toStr(end) : undefined,
    }).then((res) => {
      Object.assign(order, res.order || {});
      Object.assign(income, res.income || {});
      trend.value = res.trend || [];
      renderChart();
    });
  }

  function renderChart() {
    const dateArr = trend.value.map((d) => d.date);
    const orderSeries = trend.value.map((d) => d.orderAmount || 0);
    const incomeSeries = trend.value.map((d) => d.incomeAmount || 0);

    setOptions({
      tooltip: { trigger: 'axis' },
      legend: { data: ['订单金额', '产出金额'] },
      grid: { left: '3%', right: '3%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: dateArr },
      yAxis: { type: 'value' },
      series: [
        {
          name: '订单金额',
          type: 'bar',
          data: orderSeries,
          barMaxWidth: 30,
          itemStyle: { color: '#18a058' },
        },
        {
          name: '产出金额',
          type: 'bar',
          data: incomeSeries,
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

  .mt-3 {
    margin-top: 12px;
  }

  .stat-block {
    border-left: 3px solid transparent;
  }

  .stat-block--order {
    border-left-color: #18a058;
    background-color: #f7fdf9;
  }

  .stat-block--income {
    border-left-color: #f0a020;
    background-color: #fdfaf3;
  }

  .stat-title {
    font-weight: 600;
  }

  .stat-title--order {
    color: #18a058;
  }

  .stat-title--income {
    color: #f0a020;
  }
</style>

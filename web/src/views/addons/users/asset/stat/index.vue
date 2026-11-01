<template>
  <div>
    <n-card :bordered="false" class="proCard" title="资产统计">
      <n-grid :cols="3" :x-gap="16" responsive="screen" item-responsive>
        <n-gi v-for="item in statistics" :key="item.label" span="1 s:1 m:1 l:1">
          <n-card :bordered="true" class="stat-card">
            <n-statistic :label="item.label" :value="item.value" />
          </n-card>
        </n-gi>
      </n-grid>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref } from 'vue';
  import { Stat } from '@/api/addons/users/asset';

  interface StatItem {
    label: string;
    value: string;
  }

  const statistics = ref<StatItem[]>([]);

  function formatAmount(val: any) {
    return Number(val || 0).toFixed(2);
  }

  function loadData() {
    Stat({}).then((res) => {
      statistics.value = [
        { label: '账户数', value: String(res.accountCount || 0) },
        { label: '可用余额合计', value: formatAmount(res.totalAvailable) },
        { label: '冻结金额合计', value: formatAmount(res.totalFrozen) },
      ];
    });
  }

  onMounted(() => {
    loadData();
  });
</script>

<style lang="less" scoped>
  .stat-card {
    margin-bottom: 16px;
  }
</style>

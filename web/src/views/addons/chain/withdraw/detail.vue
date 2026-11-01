<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="提现详情"
    :bordered="false"
    :mask-closable="false"
    :style="{ width: dialogWidth }"
  >
    <n-spin :show="loading">
      <n-descriptions v-if="data" :column="1" bordered label-placement="left" size="small">
        <n-descriptions-item label="记录ID">{{ data.id }}</n-descriptions-item>
        <n-descriptions-item label="用户ID">
          <n-tag size="small" :bordered="false" type="info">#{{ data.uid }}</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="链 / 币种">
          {{ data.chain }} / {{ data.symbol }}
        </n-descriptions-item>
        <n-descriptions-item label="状态">{{ statusLabel }}</n-descriptions-item>
        <n-descriptions-item label="申请金额">
          {{ data.payAmount }} {{ data.symbol }}
        </n-descriptions-item>
        <n-descriptions-item label="到账金额">
          <span class="receive-text">{{ data.receiveAmount }} {{ data.symbol }}</span>
        </n-descriptions-item>
        <n-descriptions-item label="手续费">
          {{ data.fee }} {{ data.symbol }}
          <span class="rate-text">（费率 {{ data.rate }}%）</span>
        </n-descriptions-item>
        <n-descriptions-item label="用户钱包">
          <n-text v-if="data.toAddress" copyable>{{ data.toAddress }}</n-text>
          <span v-else>-</span>
        </n-descriptions-item>
        <n-descriptions-item label="出币地址">
          <n-text v-if="data.fromAddress" copyable>{{ data.fromAddress }}</n-text>
          <span v-else>-</span>
        </n-descriptions-item>
        <n-descriptions-item label="合约地址">
          <n-text v-if="data.contractAddress" copyable>{{ data.contractAddress }}</n-text>
          <span v-else>-</span>
        </n-descriptions-item>
        <n-descriptions-item label="交易哈希">
          <n-text v-if="data.hash" copyable>{{ data.hash }}</n-text>
          <span v-else>-</span>
        </n-descriptions-item>
        <n-descriptions-item label="审核时间">{{ data.auditAt || '-' }}</n-descriptions-item>
        <n-descriptions-item label="申请时间">{{ data.createdAt }}</n-descriptions-item>
      </n-descriptions>
    </n-spin>

    <template #footer>
      <n-space justify="end">
        <n-button @click="show = false">关闭</n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import { useDictStore } from '@/store/modules/dict';
  import { View } from '@/api/addons/chain/withdraw';
  import { adaModalWidth } from '@/utils/hotgo';

  const dict = useDictStore();
  const show = ref(false);
  const loading = ref(false);
  const data = ref<any>(null);
  const dialogWidth = ref(adaModalWidth(760));

  const statusLabel = computed(() => {
    if (data.value?.status === undefined || data.value?.status === null) {
      return '-';
    }
    return (
      dict.getLabel('ChainWithdrawStatusOptions', data.value.status) || String(data.value.status)
    );
  });

  function openModal(id: number) {
    show.value = true;
    loading.value = true;
    data.value = null;
    dialogWidth.value = adaModalWidth(760);

    View({ id })
      .then((res) => {
        data.value = res;
      })
      .finally(() => {
        loading.value = false;
      });
  }

  defineExpose({
    openModal,
  });
</script>

<style lang="less" scoped>
  .receive-text {
    color: #18a058;
    font-weight: 600;
  }

  .rate-text {
    color: #909399;
    font-size: 12px;
  }
</style>

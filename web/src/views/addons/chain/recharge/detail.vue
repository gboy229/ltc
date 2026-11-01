<template>
  <n-modal
    v-model:show="show"
    preset="card"
    title="充值详情"
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
        <n-descriptions-item label="充值类型">{{ typeLabel }}</n-descriptions-item>
        <n-descriptions-item label="充值数量">
          <span class="amount-text">+{{ data.amount }} {{ data.symbol }}</span>
        </n-descriptions-item>
        <n-descriptions-item label="用户钱包">
          <n-text v-if="data.fromAddress" copyable>{{ data.fromAddress }}</n-text>
          <span v-else>-</span>
        </n-descriptions-item>
        <n-descriptions-item label="充值地址">
          <n-text v-if="data.toAddress" copyable>{{ data.toAddress }}</n-text>
          <span v-else>-</span>
        </n-descriptions-item>
        <n-descriptions-item label="合约地址">
          <n-text v-if="data.contractAddress" copyable>{{ data.contractAddress }}</n-text>
          <span v-else>-</span>
        </n-descriptions-item>
        <n-descriptions-item label="区块高度">{{ data.blockHeight }}</n-descriptions-item>
        <n-descriptions-item label="交易哈希">
          <n-text v-if="data.hash" copyable>{{ data.hash }}</n-text>
          <span v-else>-</span>
        </n-descriptions-item>
        <n-descriptions-item label="充值时间">{{ data.createdAt }}</n-descriptions-item>
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
  import { View } from '@/api/addons/chain/recharge';
  import { adaModalWidth } from '@/utils/hotgo';

  const dict = useDictStore();
  const show = ref(false);
  const loading = ref(false);
  const data = ref<any>(null);
  const dialogWidth = ref(adaModalWidth(760));

  const typeLabel = computed(() => {
    if (!data.value?.type) {
      return '-';
    }
    return dict.getLabel('ChainRechargeTypeOptions', data.value.type) || data.value.type;
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
  .amount-text {
    color: #18a058;
    font-weight: 600;
  }
</style>

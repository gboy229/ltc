<template>
  <div>
    <n-card :bordered="false" class="proCard" title="资产账户">
      <BasicForm
        @register="register"
        @submit="handleSubmit"
        @reset="handleReset"
        @keyup.enter="handleSubmit"
        ref="searchFormRef"
      />

      <BasicTable
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row) => row.userId"
        ref="actionRef"
        :scroll-x="scrollX"
        :resizeHeightOffset="-20000"
      />
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, ref } from 'vue';
  import { BasicTable } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { AccountList } from '@/api/addons/users/asset';
  import { adaTableScrollX } from '@/utils/hotgo';
  import { columns, loadOptions, schemas } from './model';

  const actionRef = ref();
  const searchFormRef = ref<any>({});
  const scrollX = computed(() => adaTableScrollX(columns, 0));

  const [register] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 80,
    schemas,
  });

  const loadDataTable = async (res) => {
    return await AccountList({ ...res, ...searchFormRef.value?.formModel });
  };

  function handleSubmit(_values: Recordable) {
    actionRef.value.reload();
  }

  function handleReset(_values: Recordable) {
    actionRef.value.reload();
  }

  onMounted(() => {
    loadOptions();
  });
</script>

<style lang="less" scoped></style>

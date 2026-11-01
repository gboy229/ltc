<template>
  <div>
    <n-card :bordered="false" class="proCard" title="报单列表">
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
        :row-key="(row) => row.id"
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
  import { List } from '@/api/addons/pledge/order';
  import { adaTableScrollX } from '@/utils/hotgo';
  import { columns } from './columns';
  import { loadOptions, schemas } from './model';

  const actionRef = ref();
  const searchFormRef = ref<any>({});
  const scrollX = computed(() => adaTableScrollX(columns, 0));

  const [register] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 90,
    schemas,
  });

  const loadDataTable = async (res) => {
    return await List({ ...res, ...searchFormRef.value?.formModel });
  };

  function handleSubmit() {
    actionRef.value.reload();
  }

  function handleReset() {
    actionRef.value.reload();
  }

  onMounted(() => loadOptions());
</script>

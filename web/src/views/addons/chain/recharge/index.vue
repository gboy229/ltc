<template>
  <div>
    <n-card :bordered="false" class="proCard" title="充值记录">
      <BasicForm
        @register="register"
        @submit="handleSubmit"
        @reset="handleReset"
        @keyup.enter="handleSubmit"
        ref="searchFormRef"
        :show-advanced-button="false"
      />

      <BasicTable
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row) => row.id"
        ref="actionRef"
        :actionColumn="actionColumn"
        :scroll-x="scrollX"
        :resizeHeightOffset="-20000"
      />
    </n-card>

    <RechargeDetail ref="detailRef" />
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, onMounted, reactive, ref } from 'vue';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { List } from '@/api/addons/chain/recharge';
  import { adaTableScrollX } from '@/utils/hotgo';
  import { columns } from './columns';
  import { loadOptions, schemas } from './model';
  import RechargeDetail from './detail.vue';

  const actionRef = ref();
  const detailRef = ref();
  const searchFormRef = ref<any>({});

  const actionColumn = reactive({
    width: 90,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record) {
      return h(TableAction as any, {
        style: 'button',
        actions: [
          {
            label: '详情',
            onClick: handleDetail.bind(null, record),
            type: 'primary',
          },
        ],
      });
    },
  });

  const scrollX = computed(() => {
    return adaTableScrollX(columns, actionColumn.width);
  });

  const [register, {}] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 90,
    schemas,
  });

  const loadDataTable = async (res) => {
    return await List({ ...res, ...searchFormRef.value?.formModel });
  };

  function reloadTable() {
    actionRef.value.reload();
  }

  function handleDetail(record: Recordable) {
    detailRef.value.openModal(record.id);
  }

  function handleSubmit(_values: Recordable) {
    reloadTable();
  }

  function handleReset(_values: Recordable) {
    reloadTable();
  }

  onMounted(() => {
    loadOptions();
  });
</script>

<style lang="less" scoped>
  :deep(.amount-text) {
    color: #18a058;
    font-weight: 600;
  }
</style>

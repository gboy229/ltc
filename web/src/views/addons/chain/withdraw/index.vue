<template>
  <div>
    <n-card :bordered="false" class="proCard" title="提现记录">
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

    <n-modal
      v-model:show="statusModal.show"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      title="提现审核"
      :style="{ width: '480px' }"
    >
      <n-form label-placement="left" :label-width="80" class="py-4">
        <n-form-item label="审核结果">
          <n-select
            v-model:value="statusModal.status"
            :options="dict.getOptionUnRef('ChainWithdrawAuditOptions')"
            placeholder="请选择审核结果"
          />
        </n-form-item>
        <n-form-item label="审核备注">
          <n-input
            v-model:value="statusModal.remark"
            type="textarea"
            maxlength="255"
            show-count
            placeholder="选填，审核拒绝时建议填写原因"
          />
        </n-form-item>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="statusModal.show = false">取消</n-button>
          <n-button type="info" :loading="statusLoading" @click="confirmStatus">确定</n-button>
        </n-space>
      </template>
    </n-modal>

    <WithdrawDetail ref="detailRef" />
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, onMounted, reactive, ref } from 'vue';
  import { useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { List, Status } from '@/api/addons/chain/withdraw';
  import { adaTableScrollX } from '@/utils/hotgo';
  import { useDictStore } from '@/store/modules/dict';
  import { columns } from './columns';
  import { loadOptions, schemas } from './model';
  import WithdrawDetail from './detail.vue';

  const message = useMessage();
  const dict = useDictStore();
  const actionRef = ref();
  const detailRef = ref();
  const searchFormRef = ref<any>({});
  const statusLoading = ref(false);
  const statusModal = reactive({
    show: false,
    id: 0,
    status: 1,
    remark: '',
  });

  const actionColumn = reactive({
    width: 160,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record) {
      return h(TableAction as any, {
        style: 'button',
        actions: [
          {
            label: '审核',
            onClick: handleStatusModal.bind(null, record),
            auth: ['/chain/withdraw/status'],
            type: 'primary',
            ifShow: record.status === 0,
          },
          {
            label: '详情',
            onClick: handleDetail.bind(null, record),
            type: 'default',
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

  function handleStatusModal(record: Recordable) {
    statusModal.id = record.id;
    statusModal.status = 1;
    statusModal.remark = '';
    statusModal.show = true;
  }

  function confirmStatus() {
    statusLoading.value = true;
    Status({ id: statusModal.id, status: statusModal.status, remark: statusModal.remark })
      .then((_res) => {
        message.success('操作成功');
        statusModal.show = false;
        reloadTable();
      })
      .finally(() => {
        statusLoading.value = false;
      });
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
  :deep(.receive-text) {
    color: #18a058;
    font-weight: 600;
  }
</style>

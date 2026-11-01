<template>
  <div>
    <n-card :bordered="false" class="proCard" title="用户管理">
      <BasicForm
        @register="register"
        @submit="handleSubmit"
        @reset="handleReset"
        @keyup.enter="handleSubmit"
        ref="searchFormRef"
      >
        <template #statusSlot="{ model, field }">
          <n-input v-model:value="model[field]" />
        </template>
      </BasicForm>

      <BasicTable
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row) => row.id"
        ref="actionRef"
        :actionColumn="actionColumn"
        :scroll-x="scrollX"
        :resizeHeightOffset="-20000"
      >
        <template #tableTitle>
          <n-button
            type="primary"
            @click="addTable"
            class="min-left-space"
            v-if="hasPermission(['/users/user/edit'])"
          >
            <template #icon>
              <n-icon>
                <PlusOutlined />
              </n-icon>
            </template>
            添加用户
          </n-button>
        </template>
      </BasicTable>
    </n-card>

    <Edit ref="editRef" @reload-table="reloadTable" />
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, onMounted, reactive, ref } from 'vue';
  import { useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { List, Status } from '@/api/addons/users/user';
  import { PlusOutlined } from '@vicons/antd';
  import { adaTableScrollX } from '@/utils/hotgo';
  import { usePermission } from '@/hooks/web/usePermission';
  import { columns } from './columns';
  import { loadOptions, schemas } from './model';
  import Edit from './edit.vue';

  const { hasPermission } = usePermission();
  const message = useMessage();
  const actionRef = ref();
  const searchFormRef = ref<any>({});
  const editRef = ref();

  const actionColumn = reactive({
    width: 150,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record) {
      return h(TableAction as any, {
        style: 'button',
        actions: [
          {
            label: '设为禁用',
            onClick: handleStatus.bind(null, record, 2),
            ifShow: () => {
              return record.status === 1;
            },
            auth: ['/users/user/status'],
          },
          {
            label: '设为启用',
            onClick: handleStatus.bind(null, record, 1),
            ifShow: () => {
              return record.status === 2;
            },
            auth: ['/users/user/status'],
          },
          {
            label: '编辑',
            onClick: handleEdit.bind(null, record),
            auth: ['/users/user/edit'],
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
    labelWidth: 80,
    showAdvancedButton: false,
    schemas,
  });

  function addTable() {
    editRef.value.openModal(null);
  }

  const loadDataTable = async (res) => {
    return await List({ ...res, ...searchFormRef.value?.formModel });
  };

  function reloadTable() {
    actionRef.value.reload();
  }

  function handleEdit(record: Recordable) {
    editRef.value.openModal(record);
  }

  function handleStatus(record: Recordable, status: number) {
    Status({ id: record.id, status }).then((_res) => {
      message.success('操作成功');
      setTimeout(() => {
        reloadTable();
      });
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

<style lang="less" scoped></style>

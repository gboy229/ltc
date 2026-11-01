<template>
  <div>
    <n-card :bordered="false" class="proCard" title="公告管理">
      <BasicForm
        @register="register"
        @submit="handleSubmit"
        @reset="handleReset"
        @keyup.enter="handleSubmit"
        ref="searchFormRef"
      />

      <BasicTable
        :openChecked="true"
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row) => row.id"
        ref="actionRef"
        :actionColumn="actionColumn"
        @update:checked-row-keys="onCheckedRow"
        :scroll-x="scrollX"
        :resizeHeightOffset="-20000"
      >
        <template #tableTitle>
          <n-button
            type="primary"
            @click="addTable"
            class="min-left-space"
            v-if="hasPermission(['/cms/article/edit'])"
          >
            <template #icon>
              <n-icon>
                <PlusOutlined />
              </n-icon>
            </template>
            添加公告
          </n-button>

          <n-button
            type="error"
            @click="batchDelete"
            :disabled="batchDeleteDisabled"
            class="min-left-space"
            v-if="hasPermission(['/cms/article/delete'])"
          >
            <template #icon>
              <n-icon>
                <DeleteOutlined />
              </n-icon>
            </template>
            批量删除
          </n-button>
        </template>
      </BasicTable>
    </n-card>

    <Edit ref="editRef" @reload-table="reloadTable" />
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, onMounted, reactive, ref } from 'vue';
  import { useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { Delete, List, Status } from '@/api/addons/cms/article';
  import { DeleteOutlined, PlusOutlined } from '@vicons/antd';
  import { adaTableScrollX } from '@/utils/hotgo';
  import { usePermission } from '@/hooks/web/usePermission';
  import { columns } from './columns';
  import { loadOptions, schemas } from './model';
  import Edit from './edit.vue';

  const { hasPermission } = usePermission();
  const message = useMessage();
  const actionRef = ref();
  const dialog = useDialog();
  const searchFormRef = ref<any>({});
  const editRef = ref();
  const batchDeleteDisabled = ref(true);
  const checkedIds = ref([]);

  const actionColumn = reactive({
    width: 200,
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
            auth: ['/cms/article/status'],
          },
          {
            label: '设为启用',
            onClick: handleStatus.bind(null, record, 1),
            ifShow: () => {
              return record.status === 2;
            },
            auth: ['/cms/article/status'],
          },
          {
            label: '编辑',
            onClick: handleEdit.bind(null, record),
            auth: ['/cms/article/edit'],
            type: 'primary',
          },
          {
            label: '删除',
            onClick: handleDelete.bind(null, record),
            auth: ['/cms/article/delete'],
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
    schemas,
  });

  function addTable() {
    editRef.value.openModal(null);
  }

  const loadDataTable = async (res) => {
    return await List({ ...res, ...searchFormRef.value?.formModel });
  };

  function onCheckedRow(rowKeys) {
    batchDeleteDisabled.value = rowKeys.length <= 0;
    checkedIds.value = rowKeys;
  }

  function reloadTable() {
    actionRef.value.reload();
  }

  function handleEdit(record: Recordable) {
    editRef.value.openModal(record);
  }

  function handleDelete(record: Recordable) {
    dialog.warning({
      title: '警告',
      content: '你确定要删除？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        Delete({ id: record.id }).then((_res) => {
          message.success('操作成功');
          reloadTable();
        });
      },
    });
  }

  function batchDelete() {
    dialog.warning({
      title: '警告',
      content: '你确定要删除？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        Delete({ id: checkedIds.value }).then((_res) => {
          message.success('操作成功');
          reloadTable();
        });
      },
    });
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

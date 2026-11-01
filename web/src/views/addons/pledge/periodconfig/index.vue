<template>
  <div>
    <n-card :bordered="false" class="proCard" title="周期配置">
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
        :actionColumn="actionColumn"
        :scroll-x="scrollX"
        :resizeHeightOffset="-20000"
      >
        <template #tableTitle>
          <n-button
            type="primary"
            @click="addTable"
            class="min-left-space"
            v-if="hasPermission(['/pledge/periodConfig/edit'])"
          >
            <template #icon>
              <n-icon><PlusOutlined /></n-icon>
            </template>
            添加周期配置
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
  import { Delete, List, Status } from '@/api/addons/pledge/periodconfig';
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

  const actionColumn = reactive({
    width: 180,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record) {
      return h(TableAction as any, {
        style: 'button',
        actions: [
          {
            label: record.status === 1 ? '设为禁用' : '设为启用',
            onClick: handleStatus.bind(null, record, record.status === 1 ? 2 : 1),
            auth: ['/pledge/periodConfig/status'],
          },
          {
            label: '编辑',
            onClick: handleEdit.bind(null, record),
            auth: ['/pledge/periodConfig/edit'],
            type: 'primary',
          },
          {
            label: '删除',
            onClick: handleDelete.bind(null, record),
            auth: ['/pledge/periodConfig/delete'],
          },
        ],
      });
    },
  });

  const scrollX = computed(() => adaTableScrollX(columns, actionColumn.width));

  const [register] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 90,
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

  function handleSubmit() {
    reloadTable();
  }

  function handleReset() {
    reloadTable();
  }

  function handleEdit(record: Recordable) {
    editRef.value.openModal(record);
  }

  function handleStatus(record: Recordable, status: number) {
    Status({ id: record.id, status }).then(() => {
      message.success('操作成功');
      reloadTable();
    });
  }

  function handleDelete(record: Recordable) {
    dialog.warning({
      title: '警告',
      content: '你确定要删除？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        Delete({ ids: [record.id] }).then(() => {
          message.success('操作成功');
          reloadTable();
        });
      },
    });
  }

  onMounted(() => loadOptions());
</script>

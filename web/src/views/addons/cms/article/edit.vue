<template>
  <div>
    <n-drawer v-model:show="showModal" :width="dialogWidth" :mask-closable="false">
      <n-drawer-content closable>
        <template #header>
          {{ formValue.id > 0 ? '编辑公告 #' + formValue.id : '添加公告' }}
        </template>
        <n-scrollbar style="max-height: 87vh" class="pr-5">
          <n-spin :show="loading" description="请稍候...">
            <n-form
              :model="formValue"
              :rules="rules"
              ref="formRef"
              :label-placement="settingStore.isMobile ? 'top' : 'left'"
              :label-width="80"
              class="py-4"
            >
              <n-form-item label="公告标题" path="title">
                <n-input v-model:value="formValue.title" placeholder="请输入公告标题" />
              </n-form-item>

              <n-form-item label="公告内容" path="content">
                <Editor style="height: 360px" v-model:value="formValue.content" />
              </n-form-item>

              <n-form-item label="推荐首页" path="isRecommend">
                <div class="flex items-center gap-2">
                  <n-switch
                    v-model:value="formValue.isRecommend"
                    :checked-value="1"
                    :unchecked-value="0"
                  />
                  <span class="text-xs text-gray-400">开启后显示在 C 端首页公告滚动区</span>
                </div>
              </n-form-item>

              <n-grid x-gap="24" :cols="2">
                <n-gi>
                  <n-form-item label="排序" path="sort">
                    <n-input-number
                      style="width: 100%"
                      v-model:value="formValue.sort"
                      clearable
                    />
                  </n-form-item>
                </n-gi>
                <n-gi>
                  <n-form-item label="状态" path="status">
                    <n-radio-group v-model:value="formValue.status" name="status">
                      <n-radio-button
                        v-for="s in statusOptions"
                        :key="s.value"
                        :value="s.value"
                        :label="s.label"
                      />
                    </n-radio-group>
                  </n-form-item>
                </n-gi>
              </n-grid>
            </n-form>
          </n-spin>
        </n-scrollbar>
        <template #footer>
          <n-space>
            <n-button @click="closeForm">取消</n-button>
            <n-button type="primary" :loading="formBtnLoading" @click="confirmForm">确定</n-button>
          </n-space>
        </template>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useProjectSettingStore } from '@/store/modules/projectSetting';
  import { useMessage } from 'naive-ui';
  import Editor from '@/components/Editor/editor.vue';
  import { statusOptions } from '@/enums/optionsiEnum';
  import { Edit, MaxSort } from '@/api/addons/cms/article';
  import { adaModalWidth } from '@/utils/hotgo';
  import { State, newState, rules } from './model';

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const settingStore = useProjectSettingStore();
  const loading = ref(false);
  const showModal = ref(false);
  const formValue = ref<State>(newState(null));
  const formRef = ref<any>({});
  const formBtnLoading = ref(false);
  const dialogWidth = ref(adaModalWidth(800));

  function confirmForm(e) {
    e.preventDefault();
    formBtnLoading.value = true;
    formRef.value.validate((errors) => {
      if (!errors) {
        Edit(formValue.value).then((_res) => {
          message.success('操作成功');
          setTimeout(() => {
            closeForm();
            emit('reloadTable');
          });
        });
      } else {
        message.error('请填写完整信息');
      }
      formBtnLoading.value = false;
    });
  }

  // 关闭弹窗
  function closeForm() {
    showModal.value = false;
    loading.value = false;
  }

  // 打开弹窗
  function openModal(state: State | null) {
    showModal.value = true;
    dialogWidth.value = adaModalWidth(800);

    // 新增
    if (!state || state.id < 1) {
      formValue.value = newState(null);
      loading.value = true;
      MaxSort()
        .then((res) => {
          formValue.value.sort = res.sort;
        })
        .finally(() => {
          loading.value = false;
        });
      return;
    }

    // 编辑
    formValue.value = newState(state);
  }

  defineExpose({
    openModal,
  });
</script>

<style lang="less" scoped></style>

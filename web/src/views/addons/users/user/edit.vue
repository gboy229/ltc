<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      :title="formValue.id > 0 ? '编辑用户 #' + formValue.id : '添加用户'"
      :style="{ width: dialogWidth }"
    >
      <n-form
        :model="formValue"
        :rules="rules"
        ref="formRef"
        :label-placement="settingStore.isMobile ? 'top' : 'left'"
        :label-width="100"
        class="py-4"
      >
        <n-grid cols="1 s:1 m:2 l:2 xl:2 2xl:2" responsive="screen">
          <n-gi>
            <n-form-item label="昵称" path="nickname">
              <n-input v-model:value="formValue.nickname" placeholder="请输入昵称" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="会员等级" path="levelId">
              <n-select
                v-model:value="formValue.levelId"
                :options="dict.getOptionUnRef('usersLevel')"
                placeholder="请选择等级"
                clearable
              />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="状态" path="status">
              <n-radio-group v-model:value="formValue.status" name="status">
                <n-radio-button
                  v-for="s in dict.getOptionUnRef('sys_normal_disable')"
                  :key="s.value"
                  :value="s.value"
                  :label="s.label"
                />
              </n-radio-group>
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="锁定等级" path="lockLevel">
              <n-radio-group v-model:value="formValue.lockLevel" name="lockLevel">
                <n-radio-button
                  v-for="s in dict.getOptionUnRef('UserLockLevelOptions')"
                  :key="s.value"
                  :value="s.value"
                  :label="s.label"
                />
              </n-radio-group>
              <template #feedback>
                <span class="form-tip">锁定后等级只会升不会降</span>
              </template>
            </n-form-item>
          </n-gi>
        </n-grid>
      </n-form>

      <template #action>
        <n-space>
          <n-button @click="closeForm">取消</n-button>
          <n-button type="info" :loading="formBtnLoading" @click="confirmForm">确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useMessage } from 'naive-ui';
  import { useDictStore } from '@/store/modules/dict';
  import { useProjectSettingStore } from '@/store/modules/projectSetting';
  import { Edit } from '@/api/addons/users/user';
  import { adaModalWidth } from '@/utils/hotgo';
  import { State, newState, rules } from './model';

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const settingStore = useProjectSettingStore();
  const dict = useDictStore();
  const showModal = ref(false);
  const formValue = ref<State>(newState(null));
  const formRef = ref<any>({});
  const formBtnLoading = ref(false);
  const dialogWidth = ref(adaModalWidth(700));

  function confirmForm(e) {
    e.preventDefault();
    formBtnLoading.value = true;
    formRef.value.validate((errors) => {
      if (!errors) {
        Edit(formValue.value)
          .then(() => {
            message.success('操作成功');
            emit('reloadTable');
            closeForm();
          })
          .finally(() => {
            formBtnLoading.value = false;
          });
      } else {
        message.error('请填写完整信息');
        formBtnLoading.value = false;
      }
    });
  }

  // 关闭弹窗
  function closeForm() {
    showModal.value = false;
  }

  // 打开弹窗
  function openModal(state: State | null) {
    showModal.value = true;
    formValue.value = newState(state);
  }

  defineExpose({
    openModal,
  });
</script>

<style lang="less" scoped></style>

<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      :title="formValue.id > 0 ? '编辑等级 #' + formValue.id : '添加等级'"
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
            <n-form-item label="等级名称" path="name">
              <n-input v-model:value="formValue.name" placeholder="请输入等级名称" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="个人业绩" path="personalAmount">
              <n-input-number
                v-model:value="formValue.personalAmount"
                style="width: 100%"
                clearable
              />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="团队业绩" path="teamAmount">
              <n-input-number v-model:value="formValue.teamAmount" style="width: 100%" clearable />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="级差奖励(%)" path="diffReward">
              <n-input-number
                v-model:value="formValue.diffReward"
                style="width: 100%"
                :min="0"
                :precision="2"
                placeholder="例如 10"
                clearable
              />
              <template #feedback>
                <span class="form-tip">直接填写百分比数值，例如 10 表示 10%</span>
              </template>
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
  import { useDictStore } from '@/store/modules/dict';
  import { useProjectSettingStore } from '@/store/modules/projectSetting';
  import { useMessage } from 'naive-ui';
  import { Edit } from '@/api/addons/users/level';
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
  const dialogWidth = ref(adaModalWidth(840));

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
  }

  // 打开弹窗
  function openModal(state: State | null) {
    showModal.value = true;
    dialogWidth.value = adaModalWidth(840);

    // 新增
    if (!state || state.id < 1) {
      formValue.value = newState(null);
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

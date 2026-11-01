<template>
  <n-modal
    v-model:show="showModal"
    :mask-closable="false"
    :show-icon="false"
    preset="dialog"
    :title="formValue.id > 0 ? '编辑周期配置 #' + formValue.id : '添加周期配置'"
    :style="{ width: dialogWidth }"
  >
    <n-form
      :model="formValue"
      :rules="rules"
      ref="formRef"
      :label-placement="settingStore.isMobile ? 'top' : 'left'"
      :label-width="120"
      class="py-4"
    >
      <n-grid cols="1 s:1 m:2 l:2 xl:2 2xl:2" responsive="screen" :x-gap="16">
        <n-gi>
          <n-form-item label="投资周期" path="periodMonths">
            <n-select
              v-model:value="formValue.periodMonths"
              :options="dict.getOptionUnRef('PledgePeriodOptions')"
              placeholder="请选择投资周期"
            />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="月化收益区间" path="monthlyRewardInfo">
            <n-input
              v-model:value="formValue.monthlyRewardInfo"
              :maxlength="64"
              placeholder="例如 8%～10%"
            />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="静态月化率(%)" path="monthlyRewardRate">
            <n-input-number
              v-model:value="formValue.monthlyRewardRate"
              class="w-full"
              :min="0"
              :max="100"
              :precision="2"
              placeholder="例如 9"
            />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="直推收益率(%)" path="directRewardRate">
            <n-input-number
              v-model:value="formValue.directRewardRate"
              class="w-full"
              :min="0"
              :max="100"
              :precision="2"
              placeholder="例如 6"
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
      </n-grid>
    </n-form>

    <template #action>
      <n-space>
        <n-button @click="closeForm">取消</n-button>
        <n-button type="info" :loading="formBtnLoading" @click="confirmForm">确定</n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useDictStore } from '@/store/modules/dict';
  import { useProjectSettingStore } from '@/store/modules/projectSetting';
  import { useMessage } from 'naive-ui';
  import { Edit } from '@/api/addons/pledge/periodconfig';
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
  const dialogWidth = ref(adaModalWidth(760));

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
          .finally(() => (formBtnLoading.value = false));
      } else {
        message.error('请填写完整信息');
        formBtnLoading.value = false;
      }
    });
  }

  function closeForm() {
    showModal.value = false;
  }

  function openModal(state: State | null) {
    showModal.value = true;
    dialogWidth.value = adaModalWidth(760);
    formValue.value = !state || state.id < 1 ? newState(null) : newState(state);
  }

  defineExpose({ openModal });
</script>

<style lang="less" scoped>
  .w-full {
    width: 100%;
  }
</style>

<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      :title="formValue.id > 0 ? '编辑网络配置 #' + formValue.id : '添加网络配置'"
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
        <n-grid cols="1 s:1 m:2 l:2 xl:2 2xl:2" responsive="screen">
          <n-gi>
            <n-form-item label="驱动" path="driver">
              <n-select
                v-model:value="formValue.driver"
                :options="dict.getOptionUnRef('ChainDriverOptions')"
                placeholder="请选择驱动"
              />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="网络" path="chain">
              <n-input v-model:value="formValue.chain" placeholder="请输入网络" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="链ID" path="chainId">
              <n-input-number v-model:value="formValue.chainId" class="w-full" clearable />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="RPC地址" path="rpcUrl">
              <n-input v-model:value="formValue.rpcUrl" placeholder="请输入rpc地址" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="浏览器地址" path="browserUrl">
              <n-input v-model:value="formValue.browserUrl" placeholder="请输入浏览器地址" />
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
  import { Edit } from '@/api/addons/chain/chainconfig';
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

<style lang="less" scoped>
  .w-full {
    width: 100%;
  }
</style>

<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      :title="formValue.id > 0 ? '编辑合约 #' + formValue.id : '添加合约'"
      :style="{ width: dialogWidth }"
    >
      <n-form
        :model="formValue"
        :rules="rules"
        ref="formRef"
        :label-placement="settingStore.isMobile ? 'top' : 'left'"
        :label-width="130"
        class="py-2"
      >
        <n-tabs v-model:value="activeTab" type="line" animated display-directive="show">
          <n-tab-pane name="basic" tab="基础信息">
            <n-grid cols="1 s:1 m:2 l:2 xl:2 2xl:2" responsive="screen" :x-gap="16">
              <n-gi>
                <n-form-item label="链名" path="chain">
                  <n-select
                    v-model:value="formValue.chain"
                    :options="chainOptions"
                    :loading="chainLoading"
                    filterable
                    placeholder="请选择链"
                  />
                  <template #feedback>
                    <span class="form-tip">
                      只能选择已在「网络配置」中维护的链；如需新增请先到网络配置添加
                    </span>
                  </template>
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="币种" path="symbol">
                  <n-input
                    v-model:value="formValue.symbol"
                    :maxlength="30"
                    placeholder="例如 USDT"
                  />
                  <template #feedback>
                    <span class="form-tip">币种符号，用于流水和手续费展示，建议使用大写</span>
                  </template>
                </n-form-item>
              </n-gi>
              <n-gi :span="2">
                <n-form-item label="合约地址" path="address">
                  <n-input
                    v-model:value="formValue.address"
                    :maxlength="100"
                    placeholder="请输入代币合约地址"
                  />
                  <template #feedback>
                    <span class="form-tip">
                      代币合约地址格式由所选链的驱动校验（EVM 为 0x 地址，Tron 为 T 开头地址等）
                    </span>
                  </template>
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="精度" path="decimal">
                  <n-input-number
                    v-model:value="formValue.decimal"
                    class="w-full"
                    :min="0"
                    :max="36"
                    :precision="0"
                    placeholder="例如 18"
                  />
                  <template #feedback>
                    <span class="form-tip">
                      必须与合约 decimals() 一致。USDT 常见 6 或 18，多数 ERC20 为
                      18；填错会导致金额换算错误
                    </span>
                  </template>
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="保留小数位" path="scale">
                  <n-input-number
                    v-model:value="formValue.scale"
                    class="w-full"
                    :min="0"
                    :max="36"
                    :precision="0"
                    placeholder="例如 6"
                  />
                  <template #feedback>
                    <span class="form-tip"
                      >金额展示/计算保留的小数位，建议不超过精度，常用 2~8</span
                    >
                  </template>
                </n-form-item>
              </n-gi>
            </n-grid>
          </n-tab-pane>

          <n-tab-pane name="limit" tab="提币规则">
            <n-grid cols="1 s:1 m:2 l:2 xl:2 2xl:2" responsive="screen" :x-gap="16">
              <n-gi>
                <n-form-item label="最小提币数量" path="minAmount">
                  <n-input-number
                    v-model:value="formValue.minAmount"
                    class="w-full"
                    :min="0"
                    :precision="8"
                    placeholder="0 表示不限制"
                  />
                  <template #feedback>
                    <span class="form-tip"
                      >单笔最小提币数量，按币种单位填写，不是 wei；0 表示不限制</span
                    >
                  </template>
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="最大提币数量" path="maxAmount">
                  <n-input-number
                    v-model:value="formValue.maxAmount"
                    class="w-full"
                    :min="0"
                    :precision="8"
                    placeholder="0 表示不限制"
                  />
                  <template #feedback>
                    <span class="form-tip">单笔最大提币数量；0 表示不限制</span>
                  </template>
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="提币数量倍数" path="multipleAmount">
                  <n-input-number
                    v-model:value="formValue.multipleAmount"
                    class="w-full"
                    :min="0"
                    :precision="0"
                    placeholder="0 表示不限制"
                  />
                  <template #feedback>
                    <span class="form-tip"
                      >提币数量必须是该值的整数倍；0 表示不限制，例如填 10 则只能提
                      10、20、30...</span
                    >
                  </template>
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="审核阈值" path="verifyAmount">
                  <n-input-number
                    v-model:value="formValue.verifyAmount"
                    class="w-full"
                    :min="0"
                    :precision="8"
                    placeholder="0 表示全部审核"
                  />
                  <template #feedback>
                    <span class="form-tip"
                      >提币数量 ≥ 此值时需要人工审核；0 表示所有提现都需要审核</span
                    >
                  </template>
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="日上限笔数" path="dayLimit">
                  <n-input-number
                    v-model:value="formValue.dayLimit"
                    class="w-full"
                    :min="0"
                    :precision="0"
                    placeholder="0 表示不限制"
                  />
                  <template #feedback>
                    <span class="form-tip">每个用户每天最多提币笔数；0 表示不限制</span>
                  </template>
                </n-form-item>
              </n-gi>
            </n-grid>

            <n-alert type="info" :bordered="false" class="form-alert">
              提币数量、审核阈值都按币种单位填写，例如 USDT 填 100 表示 100 USDT；0
              表示不限制或全部审核。
            </n-alert>
          </n-tab-pane>

          <n-tab-pane name="fee" tab="手续费">
            <n-form-item label="手续费类型" path="feeType">
              <n-radio-group v-model:value="formValue.feeType" name="feeType">
                <n-radio-button
                  v-for="s in dict.getOptionUnRef('ChainFeeTypeOptions')"
                  :key="s.value"
                  :value="s.value"
                  :label="s.label"
                />
              </n-radio-group>
              <template #feedback>
                <span class="form-tip"
                  >百分比：按提币数量比例收取；固定个数：每笔收取固定数量的代币</span
                >
              </template>
            </n-form-item>

            <n-form-item label="手续费费率" path="rate">
              <n-input-number
                v-model:value="formValue.rate"
                class="w-full"
                :min="0"
                :max="formValue.feeType === 0 ? 100 : null"
                :precision="8"
                :placeholder="
                  formValue.feeType === 0 ? '例如 1 表示 1%' : '例如 1 表示每笔 1 个代币'
                "
              />
              <template #feedback>
                <span class="form-tip">{{ rateTip }}</span>
              </template>
            </n-form-item>

            <n-alert type="warning" :bordered="false" class="form-alert">
              {{ feeExample }}
            </n-alert>
          </n-tab-pane>

          <n-tab-pane name="switch" tab="功能开关">
            <n-form-item label="充币开关" path="rechargeSwitch">
              <n-radio-group v-model:value="formValue.rechargeSwitch" name="rechargeSwitch">
                <n-radio-button
                  v-for="s in dict.getOptionUnRef('ChainSwitchOptions')"
                  :key="s.value"
                  :value="s.value"
                  :label="s.label"
                />
              </n-radio-group>
              <template #feedback>
                <span class="form-tip">开启后才会扫描并处理该币种的充值</span>
              </template>
            </n-form-item>

            <n-form-item label="提现开关" path="withdrawSwitch">
              <n-radio-group v-model:value="formValue.withdrawSwitch" name="withdrawSwitch">
                <n-radio-button
                  v-for="s in dict.getOptionUnRef('ChainSwitchOptions')"
                  :key="s.value"
                  :value="s.value"
                  :label="s.label"
                />
              </n-radio-group>
              <template #feedback>
                <span class="form-tip">开启后用户才能发起该币种提现，建议配置确认无误后再开启</span>
              </template>
            </n-form-item>
          </n-tab-pane>
        </n-tabs>
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
  import { computed, onMounted, ref } from 'vue';
  import { useDictStore } from '@/store/modules/dict';
  import { useProjectSettingStore } from '@/store/modules/projectSetting';
  import { useMessage } from 'naive-ui';
  import { Edit } from '@/api/addons/chain/contract';
  import { Options } from '@/api/addons/chain/chainconfig';
  import { adaModalWidth } from '@/utils/hotgo';
  import { State, newState, rules } from './model';

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const settingStore = useProjectSettingStore();
  const dict = useDictStore();
  const showModal = ref(false);
  const activeTab = ref('basic');
  const formValue = ref<State>(newState(null));
  const formRef = ref<any>({});
  const formBtnLoading = ref(false);
  const dialogWidth = ref(adaModalWidth(960));
  const chainOptions = ref<Array<{ label: string; value: string }>>([]);
  const chainLoading = ref(false);

  const rateTip = computed(() => {
    if (formValue.value.feeType === 0) {
      return '按百分比收取：填 1 表示 1%，填 0.5 表示 0.5%；手续费 = 提币数量 × 费率 ÷ 100';
    }
    return `按固定数量收取：填 1 表示每笔收取 1 个 ${formValue.value.symbol || '代币'}；手续费 = 费率`;
  });

  const feeExample = computed(() => {
    const symbol = formValue.value.symbol || '代币';
    if (formValue.value.feeType === 0) {
      return `示例：提币 100 ${symbol}，费率填 1，则手续费 = 100 × 1 ÷ 100 = 1 ${symbol}。`;
    }
    return `示例：费率填 1，则每笔固定收取 1 个 ${symbol}。`;
  });

  async function loadChainOptions() {
    chainLoading.value = true;
    try {
      const res = await Options();
      const list = (res.list || []) as Array<{ label: string; value: string }>;
      if (formValue.value.chain && !list.some((item) => item.value === formValue.value.chain)) {
        list.unshift({ label: formValue.value.chain, value: formValue.value.chain });
      }
      chainOptions.value = list;
    } finally {
      chainLoading.value = false;
    }
  }

  function confirmForm(e) {
    e.preventDefault();

    const value = formValue.value;
    if (value.scale > value.decimal) {
      message.error('保留小数位不能大于代币精度');
      return;
    }
    if (value.maxAmount > 0 && value.minAmount > 0 && value.maxAmount < value.minAmount) {
      message.error('最大提币数量不能小于最小提币数量');
      return;
    }
    if (value.feeType === 0 && value.rate > 100) {
      message.error('百分比手续费费率不能大于 100');
      return;
    }

    formBtnLoading.value = true;
    formRef.value.validate((errors) => {
      if (errors) {
        message.error('请检查表单填写是否完整、格式是否正确');
        formBtnLoading.value = false;
        return;
      }

      Edit(formValue.value)
        .then((_res) => {
          message.success('操作成功');
          closeForm();
          emit('reloadTable');
        })
        .finally(() => {
          formBtnLoading.value = false;
        });
    });
  }

  // 关闭弹窗
  function closeForm() {
    showModal.value = false;
  }

  // 打开弹窗
  function openModal(state: State | null) {
    showModal.value = true;
    activeTab.value = 'basic';
    dialogWidth.value = adaModalWidth(960);
    formValue.value = newState(state);
    loadChainOptions();
  }

  onMounted(() => {
    loadChainOptions();
  });

  defineExpose({
    openModal,
  });
</script>

<style lang="less" scoped>
  .w-full {
    width: 100%;
  }

  .form-tip {
    color: #909399;
    font-size: 12px;
    line-height: 1.6;
  }

  .form-alert {
    margin-top: 12px;
  }
</style>

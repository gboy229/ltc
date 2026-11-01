<template>
  <div>
    <n-grid cols="24 300:1 600:24" :x-gap="12">
      <n-grid-item span="6">
        <n-card :bordered="false" size="small" class="proCard">
          <n-thing
            class="thing-cell"
            v-for="item in groups"
            :key="item.group"
            :class="{ 'thing-cell-on': activeGroup === item.group }"
            @click="switchGroup(item)"
          >
            <template #header>{{ item.name }}</template>
            <template #description>{{ item.desc }}</template>
          </n-thing>
        </n-card>
      </n-grid-item>
      <n-grid-item span="18">
        <n-card :bordered="false" size="small" :title="activeGroupName" class="proCard">
          <n-spin :show="show" description="请稍候...">
            <n-form label-placement="top" :model="formValue">
              <n-form-item v-for="item in currentItems" :key="item.key">
                <template #label>
                  <div class="item-label">
                    <span class="item-name">{{ item.name }}</span>
                    <span v-if="item.tip" class="item-tip">{{ item.tip }}</span>
                  </div>
                </template>
                <n-switch v-if="item.type === 'bool'" v-model:value="formValue[item.key]" />
                <n-input-number
                  v-else-if="isNumber(item.type)"
                  v-model:value="formValue[item.key]"
                  :show-button="false"
                  class="w-full"
                />
                <n-date-picker
                  v-else-if="isDate(item.type)"
                  v-model:value="formValue[item.key]"
                  :type="item.type"
                  class="w-full"
                />
                <n-input
                  v-else-if="item.type === 'textarea'"
                  v-model:value="formValue[item.key]"
                  type="textarea"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                />
                <n-dynamic-tags v-else-if="isSlice(item.type)" v-model:value="formValue[item.key]" />
                <n-input v-else v-model:value="formValue[item.key]" />
              </n-form-item>
              <div class="mt-10">
                <n-space>
                  <n-button type="primary" @click="formSubmit">保存更新</n-button>
                </n-space>
              </div>
            </n-form>
          </n-spin>
        </n-card>
      </n-grid-item>
    </n-grid>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, reactive, ref } from 'vue';
  import { useMessage } from 'naive-ui';
  import { getConfigItems, updateConfig } from '@/api/addons/pledge/config';

  const message = useMessage();
  const show = ref(false);
  const groups = ref<any[]>([]);
  const items = ref<any[]>([]);
  const activeGroup = ref('');
  const formValue = reactive<any>({});

  const currentItems = computed(() => items.value.filter((i) => i.group === activeGroup.value));
  const activeGroupName = computed(() => {
    const g = groups.value.find((i) => i.group === activeGroup.value);
    return g ? g.name : '';
  });

  function isNumber(t: string) {
    return [
      'int',
      'int8',
      'int16',
      'int32',
      'int64',
      'uint',
      'uint8',
      'uint16',
      'uint32',
      'uint64',
      'float32',
      'float64',
    ].includes(t);
  }

  function isDate(t: string) {
    return t === 'date' || t === 'datetime';
  }

  function isSlice(t: string) {
    return t.startsWith('[]');
  }

  function switchGroup(g: any) {
    activeGroup.value = g.group;
    buildForm();
  }

  function buildForm() {
    Object.keys(formValue).forEach((k) => delete formValue[k]);
    for (const item of currentItems.value) {
      formValue[item.key] = item.value;
    }
  }

  function formSubmit() {
    updateConfig({ group: activeGroup.value, list: { ...formValue } }).then(() => {
      message.success('更新成功');
      load();
    });
  }

  function load() {
    show.value = true;
    getConfigItems()
      .then((res) => {
        groups.value = res.groups || [];
        items.value = res.items || [];
        if (groups.value.length && !activeGroup.value) {
          activeGroup.value = groups.value[0].group;
        }
        buildForm();
      })
      .finally(() => {
        show.value = false;
      });
  }

  onMounted(() => {
    load();
  });
</script>

<style lang="less" scoped>
  .thing-cell {
    margin: 0 -16px 10px;
    padding: 5px 16px;

    &:hover {
      background: #f3f3f3;
      cursor: pointer;
    }
  }

  .thing-cell-on {
    background: #f0faff;
    color: #2d8cf0;

    ::v-deep(.n-thing-main .n-thing-header .n-thing-header__title) {
      color: #2d8cf0;
    }

    &:hover {
      background: #f0faff;
    }
  }

  .w-full {
    width: 100%;
  }

  .item-label {
    display: flex;
    flex-direction: column;
    gap: 2px;
    width: 100%;
  }

  .item-name {
    font-weight: 600;
    color: #1f2329;
    line-height: 1.5;
  }

  .item-tip {
    font-weight: 400;
    font-size: 12px;
    color: #909399;
    line-height: 1.4;
  }
</style>

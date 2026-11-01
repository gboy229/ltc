<template>
  <div class="view-account">
    <div class="view-account-header"></div>
    <div :style="containerCSS">
      <n-card :bordered="false">
        <header class="justify-between">
          <n-space justify="center">
            <div></div>
            <img src="~@/assets/images/logo.png" class="account-logo" alt="" />
            <n-gradient-text type="primary" :size="26">{{ projectName }}</n-gradient-text>
            <div></div>
          </n-space>
        </header>
        <main class="pt-24px">
          <div class="pt-18px">
            <LoginFrom />
          </div>
        </main>
      </n-card>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import LoginFrom from './login/index.vue';
  import { useUserStore } from '@/store/modules/user';

  const userStore = useUserStore();
  const projectName = computed(() => userStore.loginConfig?.projectName);

  const containerCSS = computed(() => {
    const val = document.body.clientWidth;
    return val <= 720
      ? {}
      : {
          flex: `1`,
          padding: `62px 12px`,
          'max-width': `484px`,
          'min-width': '320px',
          margin: '0 auto',
        };
  });
</script>

<style lang="less" scoped>
  .view-account {
    display: flex;
    flex-direction: column;
    height: 100vh;
    overflow: auto;
  }

  @media (min-width: 768px) {
    .view-account {
      background-image: url('../../assets/images/login.svg');
      background-repeat: no-repeat;
      background-position: 50%;
      background-size: 100%;
    }
  }

  .pt-24px {
    padding-top: 24px;
  }

  .pt-18px {
    padding-top: 18px;
  }

  .account-logo {
    width: 42px;
    height: 42px;
  }
</style>

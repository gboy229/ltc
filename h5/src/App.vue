<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'

import AuthGate from '@/components/AuthGate.vue'
import H5Shell from '@/components/H5Shell.vue'
import RegisterModal from '@/components/RegisterModal.vue'
import WalletSheet from '@/components/WalletSheet.vue'
import { useChainStore } from '@/stores/chain'

const chain = useChainStore()

onMounted(() => {
  // 应用启动加载一次支付链配置，BuyModal 直接复用缓存
  void chain.load()
})
</script>

<template>
  <H5Shell>
    <AuthGate>
      <RouterView v-slot="{ Component }">
        <Transition name="page" mode="out-in">
          <component :is="Component" />
        </Transition>
      </RouterView>
    </AuthGate>

    <WalletSheet />
    <RegisterModal />
  </H5Shell>
</template>

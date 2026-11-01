<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import logo from '@/assets/images/logo.png'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const user = useUserStore()

const needsAuth = computed(() => Boolean(route.meta.requiresAuth))
const blocked = computed(
  () => needsAuth.value && !(user.initialized && user.isLoggedIn),
)
const loading = computed(() => needsAuth.value && !user.initialized)

onMounted(() => {
  // 全局只初始化一次，公共页面也会尝试静默恢复会话
  void user.bootstrap()
})

async function connect() {
  await user.ensureLogin()
}

function goHome() {
  void router.push('/')
}
</script>

<template>
  <template v-if="!blocked">
    <slot />
  </template>

  <div
    v-else
    class="flex flex-1 flex-col items-center justify-center px-6 pb-[calc(1.5rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))]"
  >
    <div v-if="loading" class="flex flex-col items-center gap-3 text-ink-500">
      <span class="h-8 w-8 animate-spin rounded-full border-2 border-[#e3cda0] border-t-[#c9954a]" />
      <span class="text-[13px]">正在恢复钱包登录...</span>
    </div>

    <section
      v-else
      class="w-full max-w-[320px] rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-5 py-7 text-center shadow-gold-soft backdrop-blur-[2px]"
    >
      <img :src="logo" alt="LTC" class="mx-auto h-16 w-16 object-contain" />

      <h2 class="mt-3 text-[20px] font-black tracking-tight text-ink-900">
        连接钱包后查看
      </h2>
      <p class="mt-2 text-[12.5px] leading-5 text-ink-500">
        该页面需要登录，连接钱包并完成签名后即可继续
      </p>

      <button
        type="button"
        class="mt-5 flex h-11 w-full items-center justify-center rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[16px] font-semibold text-white shadow-gold-btn transition active:scale-[0.99] disabled:opacity-60"
        :disabled="user.connecting"
        @click="connect"
      >
        {{ user.connecting ? '连接中...' : '连接钱包' }}
      </button>

      <button
        type="button"
        class="mt-2.5 h-10 w-full rounded-full border border-[#e7d2aa] bg-[#fdf8ed] text-[14px] font-semibold text-[#5a4a35] transition active:scale-[0.99]"
        @click="goHome"
      >
        返回首页
      </button>
    </section>
  </div>
</template>

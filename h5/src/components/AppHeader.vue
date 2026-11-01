<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'

import icon2 from '@/assets/images/icon2.png'
import icon3 from '@/assets/images/icon3.png'
import icon4 from '@/assets/images/icon4.png'
import logo from '@/assets/images/logo.png'
import { useSideNavStore } from '@/stores/sideNav'
import { useUserStore } from '@/stores/user'
import { copyText } from '@/utils/clipboard'
import { toast } from '@/utils/toast'

const router = useRouter()
const sideNav = useSideNavStore()
const user = useUserStore()
const { address, shortAddress, connecting, isLoggedIn } = storeToRefs(user)

async function handleWallet() {
  if (isLoggedIn.value && address.value) {
    const copied = await copyText(address.value)
    toast(copied ? '钱包地址已复制' : '复制失败，请手动复制', copied ? 'success' : 'error')
    return
  }

  await user.connect()
}
</script>

<template>
  <header class="grid grid-cols-[1fr_auto_1fr] items-center gap-2 max-[359px]:gap-1">
    <button
      type="button"
      class="flex items-center gap-1.5 justify-self-start transition active:opacity-70 max-[359px]:gap-1"
      aria-label="返回首页"
      @click="router.push('/')"
    >
      <img
        :src="logo"
        alt="LTC"
        class="h-10 w-9 object-contain max-[359px]:h-9 max-[359px]:w-8"
      />
      <span
        class="text-[24px] font-black leading-none tracking-[-0.03em] text-[#17130f] max-[359px]:text-[22px]"
      >
        LTC
      </span>
    </button>

    <button
      type="button"
      class="flex h-7 w-[150px] min-w-0 max-w-full items-center gap-1 justify-self-center rounded-full border border-[#d8c094] bg-[#fdf8ea]/80 px-2.5 text-ink-700 shadow-[0_1px_4px_rgba(160,115,45,0.08)] backdrop-blur-sm transition active:bg-[#f8efdb] max-[359px]:w-[132px] max-[359px]:gap-0.5 max-[359px]:px-1.5"
      :aria-label="isLoggedIn ? '复制钱包地址' : '连接钱包'"
      @click="handleWallet"
    >
      <img :src="icon4" alt="" class="h-4 w-4 shrink-0 object-contain" />
      <span
        class="min-w-0 flex-1 truncate text-center text-[12px] font-medium tracking-tight max-[359px]:text-[11px]"
      >
        {{ connecting ? '连接中...' : isLoggedIn ? shortAddress : '连接钱包' }}
      </span>
      <svg
        v-if="isLoggedIn"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.6"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="h-3 w-3 shrink-0 text-[#6b6358] max-[359px]:h-2.5 max-[359px]:w-2.5"
        aria-hidden="true"
      >
        <rect x="9" y="9" width="11" height="11" rx="2.5" />
        <path d="M5 15V5.5A2.5 2.5 0 0 1 7.5 3H17" />
      </svg>
      <svg
        v-else
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.8"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="h-3 w-3 shrink-0 text-[#6b6358] max-[359px]:h-2.5 max-[359px]:w-2.5"
        aria-hidden="true"
      >
        <path d="m9 6 6 6-6 6" />
      </svg>
    </button>

    <div class="flex items-center gap-4 justify-self-end max-[359px]:gap-3">
      <button
        type="button"
        class="relative flex h-5 w-5 items-center justify-center"
        aria-label="查看公告"
        @click="router.push('/notice')"
      >
        <img :src="icon2" alt="" class="h-5 w-5 object-contain" />
      </button>
      <button
        type="button"
        class="flex h-[18px] w-[18px] items-center justify-center"
        aria-label="打开菜单"
        aria-controls="side-nav-drawer"
        :aria-expanded="sideNav.isOpen"
        @click="sideNav.openNav()"
      >
        <img :src="icon3" alt="" class="h-[18px] w-[18px] object-contain" />
      </button>
    </div>
  </header>
</template>

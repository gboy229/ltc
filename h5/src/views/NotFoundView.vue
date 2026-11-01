<script setup lang="ts">
import { RouterLink } from 'vue-router'

import logo from '@/assets/images/logo.png'
</script>

<template>
  <div
    class="flex flex-1 flex-col items-center justify-center px-8 pb-[calc(2rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))] text-center"
  >
    <img :src="logo" alt="LTC" class="h-14 w-12 object-contain opacity-90" />

    <p
      class="mt-5 bg-gradient-to-r from-[#936027] via-[#d9a75d] to-[#936027] bg-clip-text text-6xl font-black tracking-tight text-transparent"
    >
      404
    </p>
    <h1 class="mt-4 text-[18px] font-bold text-ink-900">页面不存在</h1>
    <p class="mt-2 max-w-[240px] text-[13px] leading-5 text-ink-500">
      你访问的页面可能已被移除，或者链接地址有误。
    </p>

    <RouterLink
      to="/"
      class="mt-7 rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] px-7 py-3 text-[14px] font-semibold text-white shadow-gold-btn transition active:scale-[0.98]"
    >
      返回首页
    </RouterLink>
  </div>
</template>

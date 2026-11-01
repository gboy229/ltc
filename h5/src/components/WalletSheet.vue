<script setup lang="ts">
import { storeToRefs } from 'pinia'

import { useUserStore } from '@/stores/user'
import { copyText } from '@/utils/clipboard'
import { toast } from '@/utils/toast'
import {
  openWallet,
  WALLET_LIST,
  type DiscoveredWallet,
  type WalletId,
} from '@/web3/wallet'

const user = useUserStore()
const { sheetVisible, wallets, connecting } = storeToRefs(user)

async function pick(wallet: DiscoveredWallet) {
  await user.connect(wallet)
}

function install(id: WalletId) {
  openWallet(id)
}

async function copyLink() {
  const copied = await copyText(window.location.href)
  toast(
    copied ? '链接已复制，请在钱包 App 中打开' : '复制失败，请手动复制',
    copied ? 'success' : 'error',
  )
}

function close() {
  user.closeSheet()
}
</script>

<template>
  <Transition
    enter-active-class="transition-opacity duration-200 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition-opacity duration-150 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="sheetVisible"
      class="fixed inset-0 z-[130]"
      role="dialog"
      aria-modal="true"
      aria-label="连接钱包"
    >
      <div class="absolute inset-0 bg-[#3b2f22]/50" @click="close" />

      <section
        class="absolute inset-x-0 bottom-0 rounded-t-[24px] border-t border-[#e7d2a8] bg-[#fffdf8] px-4 pb-[calc(1.25rem+env(safe-area-inset-bottom))] pt-2.5 shadow-[0_-18px_50px_-20px_rgba(80,50,12,0.45)]"
      >
        <div class="mx-auto h-1 w-10 rounded-full bg-[#e2d4bd]" />

        <header class="relative mt-3">
          <h2 class="text-center text-[19px] font-black tracking-tight text-ink-900">连接钱包</h2>
          <p class="mt-1 text-center text-[12.5px] leading-5 text-ink-500">
            {{
              wallets.length
                ? '选择钱包并签名，首次登录将自动创建账户'
                : '未检测到 Web3 钱包，请在钱包 App 的内置浏览器中打开本页'
            }}
          </p>
          <button
            type="button"
            class="absolute right-0 top-0 flex h-7 w-7 items-center justify-center rounded-full text-[#7d6a51] transition active:bg-black/5"
            aria-label="关闭"
            @click="close"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.7"
              stroke-linecap="round"
              class="h-[18px] w-[18px]"
              aria-hidden="true"
            >
              <path d="M6 6l12 12M18 6 6 18" />
            </svg>
          </button>
        </header>

        <div v-if="wallets.length" class="mt-4 space-y-2">
          <button
            v-for="wallet in wallets"
            :key="wallet.rdns || wallet.id"
            type="button"
            class="flex w-full items-center gap-3 rounded-2xl border border-[#e9d9ba] bg-[#fffdf9] px-3 py-2.5 text-left transition active:scale-[0.99] active:bg-[#fbf4e6] disabled:opacity-60"
            :disabled="connecting"
            @click="pick(wallet)"
          >
            <img
              v-if="wallet.icon"
              :src="wallet.icon"
              alt=""
              class="h-9 w-9 shrink-0 rounded-xl object-contain"
            />
            <span
              v-else
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl text-[12px] font-bold text-white"
              :class="wallet.colorClass"
            >
              {{ wallet.short }}
            </span>
            <span class="min-w-0 flex-1">
              <span class="block truncate text-[15px] font-semibold text-ink-900">
                {{ wallet.name }}
              </span>
              <span class="mt-0.5 block text-[11.5px] text-ink-500">已检测到，点击连接</span>
            </span>
            <span
              class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-[#f2d493] to-[#c9954a] text-white shadow-sm"
            >
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-3.5 w-3.5"
                aria-hidden="true"
              >
                <path d="m9 6 6 6-6 6" />
              </svg>
            </span>
          </button>
        </div>

        <div v-else class="mt-4">
          <div class="grid grid-cols-3 gap-2">
            <button
              v-for="wallet in WALLET_LIST"
              :key="wallet.id"
              type="button"
              class="flex flex-col items-center gap-1.5 rounded-2xl border border-[#e9d9ba] bg-[#fffdf9] px-1 py-3 transition active:scale-[0.98] active:bg-[#fbf4e6]"
              @click="install(wallet.id)"
            >
              <span
                class="flex h-9 w-9 items-center justify-center rounded-xl text-[12px] font-bold text-white"
                :class="wallet.colorClass"
              >
                {{ wallet.short }}
              </span>
              <span class="text-[11.5px] font-medium text-ink-700">{{ wallet.name }}</span>
            </button>
          </div>

          <p class="mt-3 rounded-xl bg-[#fbf4e6]/80 px-3 py-2 text-[11.5px] leading-5 text-ink-600">
            已安装钱包？请复制当前页面链接，在钱包 App 的「DApp 浏览器」中打开，即可自动拉起钱包登录。
          </p>

          <button
            type="button"
            class="mt-2 h-8 w-full rounded-full border border-[#e2cda2] bg-[#fffdf8] text-[12.5px] font-medium text-[#a97227] transition active:bg-[#fbf2df]"
            @click="copyLink"
          >
            复制当前页面链接
          </button>
        </div>

        <button
          type="button"
          class="mt-4 h-10 w-full rounded-full border border-[#e7d2aa] bg-[#fdf8ed] text-[14px] font-semibold text-[#5a4a35] transition active:scale-[0.99]"
          @click="close"
        >
          稍后再说
        </button>
      </section>
    </div>
  </Transition>
</template>

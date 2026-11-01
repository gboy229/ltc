<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { formatAmount } from '@/utils/format'
import { useRouter } from 'vue-router'

import icon11 from '@/assets/images/icon11.png'
import icon12 from '@/assets/images/icon12.png'
import iconIn from '@/assets/images/icon_in.png'
import iconOut from '@/assets/images/icon_out.png'
import zcBg from '@/assets/images/zc_bg1.png'
import AppHeader from '@/components/AppHeader.vue'
import PageHeading from '@/components/PageHeading.vue'
import {
  getAssetAccount,
  getAssetLogs,
  type AssetLogItem,
  type AssetLogType,
} from '@/api/asset'

const router = useRouter()

type TabKey = 'expense' | 'income'

const activeTab = ref<TabKey>('income')
const account = ref({ total: 0, available: 0, frozen: 0 })
const logs = ref<AssetLogItem[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = 15
const hasMore = ref(false)
const sentinel = ref<HTMLElement | null>(null)

let observer: IntersectionObserver | null = null
let requestSeq = 0

async function loadAccount() {
  try {
    account.value = await getAssetAccount()
  } catch {
    account.value = { total: 0, available: 0, frozen: 0 }
  }
}

async function loadLogs(reset = false, targetPage = page.value) {
  if (!reset && (loading.value || loadingMore.value)) {
    return
  }

  const seq = ++requestSeq
  const currentPage = reset ? 1 : targetPage
  if (reset) {
    loading.value = true
    loadingMore.value = false
    hasMore.value = false
    logs.value = []
  } else {
    loadingMore.value = true
  }

  try {
    const res = await getAssetLogs({
      page: currentPage,
      pageSize,
      type: activeTab.value as AssetLogType,
    })
    if (seq !== requestSeq) return

    const list = res.list ?? []
    logs.value = reset ? list : [...logs.value, ...list]
    page.value = res.page ?? currentPage
    hasMore.value = (res.pageCount ?? 0) > page.value
  } catch {
    if (seq !== requestSeq) return
    if (reset) {
      logs.value = []
      page.value = 1
    }
    hasMore.value = false
  } finally {
    if (seq === requestSeq) {
      loading.value = false
      loadingMore.value = false
      void nextTick(maybeLoadMore)
    }
  }
}

function maybeLoadMore() {
  if (loading.value || loadingMore.value || !hasMore.value) return
  const el = sentinel.value
  if (!el) return
  if (el.getBoundingClientRect().top <= window.innerHeight + 120) {
    void loadLogs(false, page.value + 1)
  }
}

function setupObserver() {
  observer?.disconnect()
  if (!sentinel.value || typeof IntersectionObserver === 'undefined') return

  observer = new IntersectionObserver(
    (entries) => {
      if (entries.some((entry) => entry.isIntersecting)) {
        maybeLoadMore()
      }
    },
    { rootMargin: '120px 0px' },
  )
  observer.observe(sentinel.value)
}

function loadMore() {
  if (!hasMore.value || loadingMore.value) return
  void loadLogs(false, page.value + 1)
}

watch(activeTab, () => {
  void loadLogs(true)
})

watch(sentinel, setupObserver)

onMounted(() => {
  void loadAccount()
  void loadLogs(true)
  window.addEventListener('scroll', maybeLoadMore, { passive: true })
})

onBeforeUnmount(() => {
  observer?.disconnect()
  window.removeEventListener('scroll', maybeLoadMore)
})

function signedAmount(record: AssetLogItem) {
  return `${record.type === 'income' ? '+' : '-'}${formatAmount(record.amount)}`
}

function goBack() {
  const state = window.history.state as { back?: string | null } | null
  if (state?.back) {
    router.back()
    return
  }
  void router.push('/')
}
</script>

<template>
  <div
    class="flex flex-1 flex-col px-4 pb-[calc(1.5rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))]"
  >
    <AppHeader />

    <PageHeading
      title="我的资产"
      subtitle="资产一览 · 收益可视 · 随时可查"
      @back="goBack"
    />

    <section
      class="relative mt-3 min-h-[164px] overflow-hidden rounded-2xl border border-[#e8d3a9] p-4 max-[359px]:min-h-[154px]"
    >
      <img
        :src="zcBg"
        alt=""
        class="pointer-events-none absolute inset-0 h-full w-full object-cover object-right"
      />
      <div class="relative z-10 flex h-full flex-col justify-between">
        <div>
          <div class="flex items-center gap-2">
            <img :src="icon11" alt="" class="h-7 w-7 object-contain" />
            <span class="text-[18px] font-bold text-[#2a1c10]">USDT 总资产</span>
          </div>
          <div class="mt-1 flex items-baseline gap-1.5 text-[#a45f1a]">
            <span
              class="text-[34px] font-black leading-none tracking-tight tabular-nums max-[359px]:text-[30px]"
            >
              {{ formatAmount(account.total) }}
            </span>
            <span class="text-[12px] font-bold">USDT</span>
          </div>
        </div>

        <div class="mt-3 grid grid-cols-2 border-t border-[#e3cda0]/80 pt-2">
          <div>
            <p class="text-[12.5px] text-ink-500">可用余额</p>
            <p
              class="mt-0.5 whitespace-nowrap text-[17px] font-bold tabular-nums text-[#1f1a13] max-[359px]:text-[15px]"
            >
              {{ formatAmount(account.available) }}
              <span class="text-[11px] font-medium text-ink-500 max-[359px]:text-[10px]">USDT</span>
            </p>
          </div>
          <div class="border-l border-[#e3cda0]/80 pl-3">
            <p class="text-[12.5px] text-ink-500">冻结</p>
            <p
              class="mt-0.5 whitespace-nowrap text-[17px] font-bold tabular-nums text-[#1f1a13] max-[359px]:text-[15px]"
            >
              {{ formatAmount(account.frozen) }}
              <span class="text-[11px] font-medium text-ink-500 max-[359px]:text-[10px]">USDT</span>
            </p>
          </div>
        </div>
      </div>
    </section>

    <button
      type="button"
      class="mt-3 flex w-full items-center gap-3 rounded-2xl border border-[#d9b26e] bg-gradient-to-r from-[#e7c07a] via-[#d8ad65] to-[#f0d49b] px-4 py-3 text-left shadow-gold-soft transition active:scale-[0.99]"
      @click="router.push('/withdraw')"
    >
      <img :src="icon12" alt="" class="h-9 w-9 shrink-0 object-contain" />
      <span class="min-w-0 flex-1">
        <span class="block text-[18px] font-bold text-[#5a3a12]">在线提币</span>
        <span class="mt-0.5 block text-[12px] text-[#7b5a2a]">在线提币到链上钱包</span>
      </span>
      <svg
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="1.8"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="h-5 w-5 shrink-0 text-[#5a3a12]"
        aria-hidden="true"
      >
        <path d="m9 6 6 6-6 6" />
      </svg>
    </button>

    <section
      class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 p-3 shadow-gold-soft"
    >
      <div class="grid grid-cols-2 rounded-full border border-[#e7d2aa] p-0.5">
        <button
          type="button"
          class="flex h-9 items-center justify-center gap-1.5 rounded-full text-[13.5px] font-medium transition"
          :class="
            activeTab === 'income'
              ? 'bg-gradient-to-r from-[#e6bf7f] via-[#d5a65c] to-[#efd09a] font-semibold text-white shadow-sm'
              : 'text-ink-700'
          "
          @click="activeTab = 'income'"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.7"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-4 w-4"
            aria-hidden="true"
          >
            <path d="M4 19V9M10 19V5M16 19v-7M22 19H2" />
          </svg>
          收入明细
        </button>
        <button
          type="button"
          class="flex h-9 items-center justify-center gap-1.5 rounded-full text-[13.5px] font-medium transition"
          :class="
            activeTab === 'expense'
              ? 'bg-gradient-to-r from-[#e6bf7f] via-[#d5a65c] to-[#efd09a] font-semibold text-white shadow-sm'
              : 'text-ink-700'
          "
          @click="activeTab = 'expense'"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.7"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-4 w-4"
            aria-hidden="true"
          >
            <path d="M6 3h12v18l-3-2-3 2-3-2-3 2V3Z" />
            <path d="M9 8h6M9 12h5" />
          </svg>
          支出明细
        </button>
      </div>

      <div class="divide-y divide-[#f1e5cc]">
        <div v-if="loading" class="py-8 text-center text-[13px] text-ink-500">加载中...</div>

        <div
          v-else-if="!logs.length"
          class="py-8 text-center text-[13px] text-ink-500"
        >
          暂无资产流水
        </div>

        <div
          v-for="(record, index) in logs"
          :key="index"
          class="flex items-center gap-2.5 py-2.5 max-[359px]:gap-2"
        >
          <img
            :src="record.type === 'expense' ? iconOut : iconIn"
            alt=""
            class="h-8 w-8 shrink-0 object-contain"
          />
          <div class="min-w-0 flex-1">
            <p class="text-[14px] font-medium text-ink-900">{{ record.remark }}</p>
            <p class="whitespace-nowrap text-[11.5px] text-ink-500 max-[359px]:text-[11px]">{{ record.date }}</p>
          </div>
          <p
            class="whitespace-nowrap text-[15px] font-bold tabular-nums max-[359px]:text-[14px]"
            :class="record.type === 'expense' ? 'text-[#d14a3a]' : 'text-[#a45f1a]'"
          >
            {{ signedAmount(record) }} USDT
          </p>
        </div>
      </div>

      <div ref="sentinel" class="h-px" aria-hidden="true" />

      <button
        v-if="hasMore"
        type="button"
        class="mt-2 flex h-9 w-full items-center justify-center rounded-full border border-[#d7b57b] bg-[#fffdf8]/60 text-[13.5px] font-medium text-[#a97227] transition active:bg-[#fbf2df]"
        :disabled="loadingMore"
        @click="loadMore"
      >
        {{ loadingMore ? '加载中...' : '查看更多' }}
      </button>
      <p v-else-if="logs.length" class="mt-2 text-center text-[12px] text-ink-400">
        没有更多了
      </p>
    </section>
  </div>
</template>

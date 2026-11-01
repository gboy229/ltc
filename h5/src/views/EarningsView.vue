<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { formatAmount } from '@/utils/format'
import { useRouter } from 'vue-router'

import icon5 from '@/assets/images/icon5.png'
import icon7 from '@/assets/images/icon7.png'
import icon8 from '@/assets/images/icon8.png'
import syBg1 from '@/assets/images/sy_bg1.png'
import syBg2 from '@/assets/images/sy_bg2.png'
import AppHeader from '@/components/AppHeader.vue'
import PageHeading from '@/components/PageHeading.vue'
import {
  getIncomeLogs,
  getOrderStatistics,
  type IncomeLogItem,
  type IncomeLogType,
} from '@/api/pledge'

const router = useRouter()

type TabKey = 'all' | IncomeLogType

const tabs = [
  { key: 'all', label: '全部' },
  { key: 'static', label: '静态收益' },
  { key: 'direct', label: '直推收益' },
  { key: 'diff', label: '级差收益' },
  { key: 'lot', label: '手数收益' },
  { key: 'principal', label: '本金返还' },
] as const

const activeTab = ref<TabKey>('all')
const records = ref<IncomeLogItem[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = 15
const hasMore = ref(false)
const sentinel = ref<HTMLElement | null>(null)

let observer: IntersectionObserver | null = null
let requestSeq = 0

const statistics = ref({
  total: 0,
  todayRelease: 0,
  totalRelease: 0,
  earningsTotal: 0,
  referralIncome: 0,
  staticIncome: 0,
})

async function loadStatistics() {
  try {
    const res = await getOrderStatistics()
    statistics.value = {
      total: res.total ?? 0,
      todayRelease: res.todayRelease ?? 0,
      totalRelease: res.totalRelease ?? 0,
      earningsTotal: res.earningsTotal ?? 0,
      referralIncome: res.referralIncome ?? 0,
      staticIncome: res.staticIncome ?? 0,
    }
  } catch {
    // 统计失败时保持 0，不阻塞页面
  }
}

async function loadRecords(reset = false, targetPage = page.value) {
  if (!reset && (loading.value || loadingMore.value)) {
    return
  }

  const seq = ++requestSeq
  const currentPage = reset ? 1 : targetPage
  if (reset) {
    loading.value = true
    loadingMore.value = false
    hasMore.value = false
    records.value = []
  } else {
    loadingMore.value = true
  }

  try {
    const type = activeTab.value === 'all' ? undefined : activeTab.value
    const res = await getIncomeLogs({ page: currentPage, pageSize, type })
    if (seq !== requestSeq) return

    const list = res.list ?? []
    records.value = reset ? list : [...records.value, ...list]
    page.value = res.page ?? currentPage
    hasMore.value = (res.pageCount ?? 0) > page.value
  } catch {
    if (seq !== requestSeq) return
    if (reset) {
      records.value = []
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
    void loadRecords(false, page.value + 1)
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

function goBack() {
  const state = window.history.state as { back?: string | null } | null
  if (state?.back) {
    router.back()
    return
  }
  void router.push('/')
}

function loadMore() {
  if (!hasMore.value || loadingMore.value) return
  void loadRecords(false, page.value + 1)
}

watch(activeTab, () => {
  void loadRecords(true)
})

watch(sentinel, setupObserver)

onMounted(() => {
  void loadStatistics()
  void loadRecords(true)
  window.addEventListener('scroll', maybeLoadMore, { passive: true })
})

onBeforeUnmount(() => {
  observer?.disconnect()
  window.removeEventListener('scroll', maybeLoadMore)
})
</script>

<template>
  <div
    class="flex flex-1 flex-col px-4 pb-[calc(1.5rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))]"
  >
    <AppHeader />

    <PageHeading
      title="收益明细"
      subtitle="收支留痕 · 账目清晰 · 有据可查"
      @back="goBack"
    />

    <section
      class="relative mt-3 min-h-[148px] overflow-hidden rounded-2xl border border-[#e8d3a9] p-4 max-[359px]:min-h-[140px]"
    >
      <img
        :src="syBg1"
        alt=""
        class="pointer-events-none absolute inset-0 h-full w-full object-cover object-right"
      />
      <div class="relative z-10 flex h-full flex-col justify-between">
        <div>
          <div class="flex items-center gap-1.5">
            <span class="h-4 w-1 shrink-0 rounded-full bg-[#c99a4a]" />
            <span class="text-[17px] font-bold text-[#2a1c10]">收益总额</span>
          </div>
          <div class="mt-1 flex items-baseline gap-1.5 text-[#a45f1a]">
            <span
              class="text-[34px] font-black leading-none tracking-tight tabular-nums max-[359px]:text-[30px]"
            >
              {{ formatAmount(statistics.earningsTotal) }}
            </span>
            <span class="text-[12px] font-bold">USDT</span>
          </div>
        </div>

        <div class="mt-3 grid grid-cols-2 border-t border-[#e3cda0]/80 pt-2">
          <div>
            <p class="text-[12.5px] text-ink-500">今日释放收益</p>
            <p
              class="mt-0.5 whitespace-nowrap text-[17px] font-bold tabular-nums text-[#1f1a13] max-[359px]:text-[15px]"
            >
              +{{ formatAmount(statistics.todayRelease) }}
              <span class="text-[11px] font-medium text-ink-500 max-[359px]:text-[10px]">USDT</span>
            </p>
          </div>
          <div class="border-l border-[#e3cda0]/80 pl-3">
            <p class="text-[12.5px] text-ink-500">累计释放收益</p>
            <p
              class="mt-0.5 whitespace-nowrap text-[17px] font-bold tabular-nums text-[#1f1a13] max-[359px]:text-[15px]"
            >
              +{{ formatAmount(statistics.totalRelease) }}
              <span class="text-[11px] font-medium text-ink-500 max-[359px]:text-[10px]">USDT</span>
            </p>
          </div>
        </div>
      </div>
    </section>

    <section
      class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 p-3.5 shadow-gold-soft"
    >
      <div class="flex items-center gap-2">
        <img :src="icon5" alt="" class="h-6 w-8 object-contain" />
        <h2 class="text-[17px] font-bold text-ink-900">入金报单</h2>
      </div>
      <div class="mt-1.5 divide-y divide-[#f1e5cc] text-[14px]">
        <div class="flex items-center justify-between gap-3 py-2">
          <span class="text-ink-500">入金报单金额</span>
          <span class="whitespace-nowrap font-medium tabular-nums text-[#a45f1a]">
            {{ formatAmount(statistics.total,2) }} USDT
          </span>
        </div>
      </div>
    </section>

    <div class="mt-3 grid grid-cols-2 gap-3">
      <section
        class="relative min-h-[104px] overflow-hidden rounded-2xl border border-[#e8d3a9] p-3"
      >
        <img
          :src="syBg2"
          alt=""
          class="pointer-events-none absolute inset-0 h-full w-full object-cover"
        />
        <div class="relative z-10">
          <img :src="icon8" alt="" class="h-6 w-8 object-contain" />
          <p class="mt-1.5 text-[13px] text-ink-500">直推收益总额</p>
          <p class="mt-1 text-[19px] font-bold tabular-nums text-[#a45f1a] max-[359px]:text-[18px]">
            {{ formatAmount(statistics.referralIncome) }}
            <span class="text-[10.5px] font-medium text-ink-500">USDT</span>
          </p>
        </div>
      </section>

      <section
        class="relative min-h-[104px] overflow-hidden rounded-2xl border border-[#e8d3a9] p-3"
      >
        <img
          :src="syBg2"
          alt=""
          class="pointer-events-none absolute inset-0 h-full w-full object-cover"
        />
        <div class="relative z-10">
          <img :src="icon7" alt="" class="h-6 w-8 object-contain" />
          <p class="mt-1.5 text-[13px] text-ink-500">静态收益总额</p>
          <p class="mt-1 text-[19px] font-bold tabular-nums text-[#a45f1a] max-[359px]:text-[18px]">
            {{ formatAmount(statistics.staticIncome) }}
            <span class="text-[10.5px] font-medium text-ink-500">USDT</span>
          </p>
        </div>
      </section>
    </div>

    <section
      class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 p-3 shadow-gold-soft"
    >
      <div class="grid grid-cols-6 gap-1 rounded-full border border-[#e7d2aa] p-0.5">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          type="button"
          class="h-9 rounded-full text-[10px] font-medium transition"
          :class="
            activeTab === tab.key
              ? 'bg-gradient-to-r from-[#e6bf7f] via-[#d5a65c] to-[#efd09a] font-semibold text-white shadow-sm'
              : 'text-ink-700'
          "
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <div v-if="loading" class="py-10 text-center text-[13px] text-ink-500">加载中...</div>

      <div
        v-else-if="!records.length"
        class="py-10 text-center text-[13px] text-ink-500"
      >
        暂无收益记录
      </div>

      <div v-else class="divide-y divide-[#f1e5cc]">
        <div
          v-for="(record, index) in records"
          :key="index"
          class="flex items-center gap-2.5 py-2.5"
        >
          <span
            class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-[#fbf4e4]"
          >
            <img :src="icon5" alt="" class="h-5 w-7 object-contain" />
          </span>
          <div class="min-w-0 flex-1">
            <p class="text-[14px] font-medium text-ink-900">{{ record.label }}</p>
            <p class="text-[11.5px] text-ink-500">{{ record.date }}</p>
          </div>
          <p class="whitespace-nowrap text-[15px] font-bold tabular-nums text-[#a45f1a]">
            +{{ formatAmount(record.amount) }} USDT
          </p>
        </div>
      </div>

      <div ref="sentinel" class="h-px" aria-hidden="true" />

      <button
        v-if="hasMore"
        type="button"
        class="mt-2 flex h-9 w-full items-center justify-center gap-1 rounded-full border border-[#d7b57b] bg-[#fffdf8]/60 text-[13.5px] font-medium text-[#a97227] transition active:bg-[#fbf2df]"
        :disabled="loadingMore"
        @click="loadMore"
      >
        {{ loadingMore ? '加载中...' : '查看更多' }}
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.8"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="h-3 w-3"
          aria-hidden="true"
        >
          <path d="m9 6 6 6-6 6" />
        </svg>
      </button>
      <p v-else-if="records.length" class="mt-2 text-center text-[12px] text-ink-400">
        没有更多了
      </p>
    </section>
  </div>
</template>

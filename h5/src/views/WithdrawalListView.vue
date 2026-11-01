<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import icon5 from '@/assets/images/icon5.png'
import AppHeader from '@/components/AppHeader.vue'
import PageHeading from '@/components/PageHeading.vue'
import {
  getWithdrawalList,
  type WithdrawalRecord,
  type WithdrawalStatus,
} from '@/api/chain/withdrawal'
import { formatAmount } from '@/utils/format'
import { copyText } from '@/utils/clipboard'
import { toast } from '@/utils/toast'

const router = useRouter()

type TabKey = 'all' | 'processing' | 'success'

const tabs = [
  { key: 'all', label: '全部' },
  { key: 'processing', label: '提币中' },
  { key: 'success', label: '已完成' },
] as const

const activeTab = ref<TabKey>('all')
const records = ref<WithdrawalRecord[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = 15
const hasMore = ref(false)
const sentinel = ref<HTMLElement | null>(null)

let observer: IntersectionObserver | null = null
let requestSeq = 0

function statusLabel(status: WithdrawalStatus) {
  if (status === 'success') return '完成'
  if (status === 'failed') return '失败'
  return '处理中'
}

function statusClass(status: WithdrawalStatus) {
  if (status === 'success') return 'bg-[#eef8f0] text-[#3a9e5a]'
  if (status === 'failed') return 'bg-[#f2f2f2] text-[#8c8c8c]'
  return 'bg-[#fbf4e4] text-[#c08a2f]'
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
    const status = activeTab.value === 'all' ? undefined : activeTab.value
    const res = await getWithdrawalList({ page: currentPage, pageSize, status })
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
  void router.push('/withdraw')
}

async function copyHash(hash: string) {
  if (!hash) return
  const ok = await copyText(hash)
  toast(ok ? '已复制' : '复制失败，请手动复制', ok ? 'success' : 'error')
}

watch(activeTab, () => {
  void loadRecords(true)
})

watch(sentinel, setupObserver)

onMounted(() => {
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
      title="提币记录"
      subtitle="链上提币 · 审核进度 · 全程留痕"
      @back="goBack"
    />

    <div class="mt-3 grid grid-cols-3 gap-2">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        class="h-9 rounded-full border text-[13px] font-medium transition active:scale-[0.98]"
        :class="
          activeTab === tab.key
            ? 'border-[#d2a25a] bg-gradient-to-r from-[#e6bf7f] via-[#d5a65c] to-[#efd09a] font-semibold text-white shadow-sm'
            : 'border-[#e5d5b7] bg-[#fffdf9] text-ink-700'
        "
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </div>

    <div class="mt-3 space-y-3">
      <div v-if="loading" class="py-10 text-center text-[13px] text-ink-500">加载中...</div>

      <div
        v-else-if="!records.length"
        class="rounded-2xl border border-dashed border-[#e5d5b7] bg-[#fffdf8]/70 py-12 text-center text-[13px] text-ink-500"
      >
        暂无提币记录
      </div>

      <article
        v-for="record in records"
        :key="record.id"
        class="rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-3.5 py-3 shadow-gold-soft"
      >
        <div class="flex items-center gap-2">
          <img :src="icon5" alt="" class="h-6 w-8 shrink-0 object-contain" />
          <span class="rounded-full bg-[#fbf4e4] px-2 py-0.5 text-[11px] font-medium text-[#8a5c22]">
            {{ record.network }} · {{ record.symbol }}
          </span>
          <span
            class="ml-auto rounded-full px-2.5 py-0.5 text-[11px] font-medium"
            :class="statusClass(record.status)"
          >
            {{ statusLabel(record.status) }}
          </span>
        </div>

        <p class="mt-2 text-[24px] font-black leading-none tracking-tight tabular-nums text-[#1f1a13]">
          -{{ formatAmount(record.amount) }}
          <span class="text-[12px] font-bold text-ink-500">USDT</span>
        </p>

        <dl class="mt-3 space-y-1.5 text-[12.5px]">
          <div class="flex items-start justify-between gap-3">
            <dt class="shrink-0 text-ink-500">收款地址</dt>
            <dd class="min-w-0 break-all text-right font-medium text-ink-800">
              {{ record.address }}
            </dd>
          </div>
          <div class="flex items-start justify-between gap-3">
            <dt class="shrink-0 text-ink-500">交易哈希</dt>
            <dd class="flex min-w-0 items-start justify-end gap-1.5 text-right font-medium text-ink-800">
              <span class="min-w-0 break-all text-[11px] leading-4">{{ record.hash || '--' }}</span>
              <button
                v-if="record.hash"
                type="button"
                class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-md text-[#a97227] transition active:bg-[#fbf2df]"
                aria-label="复制交易哈希"
                @click="copyHash(record.hash)"
              >
                <svg
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="1.7"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  class="h-3.5 w-3.5"
                  aria-hidden="true"
                >
                  <rect x="9" y="9" width="11" height="11" rx="2.5" />
                  <path d="M5 15V5.5A2.5 2.5 0 0 1 7.5 3H17" />
                </svg>
              </button>
            </dd>
          </div>
          <div class="flex items-center justify-between gap-3">
            <dt class="text-ink-500">申请时间</dt>
            <dd class="text-right font-medium text-ink-800">{{ record.date }}</dd>
          </div>
        </dl>

        <div class="mt-2.5 grid grid-cols-2 gap-2 border-t border-[#f1e5cc] pt-2.5 text-[12px]">
          <div class="flex items-center justify-between gap-2">
            <span class="text-ink-500">手续费</span>
            <span class="font-semibold tabular-nums text-[#d14a3a]">
              {{ formatAmount(record.fee) }}
            </span>
          </div>
          <div class="flex items-center justify-between gap-2 border-l border-[#f1e5cc] pl-3">
            <span class="text-ink-500">到账</span>
            <span class="font-semibold tabular-nums text-[#a45f1a]">
              {{ formatAmount(record.receiveAmount) }}
            </span>
          </div>
        </div>
      </article>
    </div>

    <div ref="sentinel" class="h-1" aria-hidden="true" />

    <div v-if="loadingMore" class="flex items-center justify-center gap-2 py-4 text-[12px] text-ink-500">
      <span class="h-4 w-4 animate-spin rounded-full border-2 border-[#e3cda0] border-t-[#c9954a]" />
      加载中...
    </div>
    <div v-else-if="records.length && !hasMore" class="py-4 text-center text-[12px] text-ink-400">
      没有更多了
    </div>
    <div v-else-if="hasMore" class="py-4 text-center text-[12px] text-ink-400">
      上拉加载更多
    </div>
  </div>
</template>

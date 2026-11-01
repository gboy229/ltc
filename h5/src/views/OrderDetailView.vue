<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { formatAmount } from '@/utils/format'
import { useRouter } from 'vue-router'

import bdBg from '@/assets/images/bd_bg1.png'
import icon5 from '@/assets/images/icon5.png'
import AppHeader from '@/components/AppHeader.vue'
import EarlyRefundModal from '@/components/EarlyRefundModal.vue'
import PageHeading from '@/components/PageHeading.vue'
import {
  getOrderList,
  getOrderStatistics,
  refundPledgeOrder,
  type PledgeOrderItem,
  type PledgeOrderStatus,
} from '@/api/pledge'
import { toast } from '@/utils/toast'
import { ApiError } from '@/utils/http'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const user = useUserStore()

type TabKey = 'all' | PledgeOrderStatus

const tabs = [
  { key: 'all', label: '全部' },
  { key: 'active', label: '进行中' },
  { key: 'completed', label: '已结束' },
] as const

const activeTab = ref<TabKey>('all')
const orders = ref<PledgeOrderItem[]>([])
const loading = ref(false)
const refunding = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = 10
const hasMore = ref(false)
const sentinel = ref<HTMLElement | null>(null)
const statistics = ref({ total: 0, todayRelease: 0, totalRelease: 0 })
const refundOrder = ref<PledgeOrderItem | null>(null)
const showRefundModal = ref(false)

const filteredOrders = computed(() => orders.value)

let observer: IntersectionObserver | null = null
let requestSeq = 0


async function loadStatistics() {
  try {
    const res = await getOrderStatistics()
    statistics.value = {
      total: res.total ?? 0,
      todayRelease: res.todayRelease ?? 0,
      totalRelease: res.totalRelease ?? 0,
    }
  } catch {
    statistics.value = { total: 0, todayRelease: 0, totalRelease: 0 }
  }
}

function statusLabel(status: PledgeOrderStatus) {
  return status === 'active' ? '进行中' : '已完成'
}

async function loadOrders(reset = false, targetPage = page.value) {
  if (!reset && (loading.value || loadingMore.value)) {
    return
  }

  const seq = ++requestSeq
  const currentPage = reset ? 1 : targetPage
  if (reset) {
    loading.value = true
    loadingMore.value = false
    hasMore.value = false
    orders.value = []
  } else {
    loadingMore.value = true
  }

  try {
    const status = activeTab.value === 'all' ? undefined : activeTab.value
    const res = await getOrderList({
      page: currentPage,
      pageSize,
      status,
    })
    if (seq !== requestSeq) {
      return
    }

    const list = res.list ?? []
    orders.value = reset ? list : [...orders.value, ...list]
    page.value = res.page ?? currentPage
    hasMore.value = (res.pageCount ?? 0) > page.value
  } catch {
    if (seq !== requestSeq) {
      return
    }
    if (reset) {
      orders.value = []
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
  if (loading.value || loadingMore.value || !hasMore.value) {
    return
  }

  const el = sentinel.value
  if (!el) {
    return
  }
  if (el.getBoundingClientRect().top <= window.innerHeight + 120) {
    void loadOrders(false, page.value + 1)
  }
}

function setupObserver() {
  observer?.disconnect()
  if (!sentinel.value || typeof IntersectionObserver === 'undefined') {
    return
  }

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

watch(sentinel, setupObserver)

watch(activeTab, () => {
  void loadOrders(true)
})

onMounted(() => {
  void loadOrders(true)
  void loadStatistics()
  window.addEventListener('scroll', maybeLoadMore, { passive: true })
})

onBeforeUnmount(() => {
  observer?.disconnect()
  window.removeEventListener('scroll', maybeLoadMore)
})

function goBack() {
  const state = window.history.state as { back?: string | null } | null
  if (state?.back) {
    router.back()
    return
  }
  void router.push('/')
}

function requestRefund(order: PledgeOrderItem) {
  refundOrder.value = order
  showRefundModal.value = true
}

async function handleRefundConfirm() {
  if (!refundOrder.value || refunding.value) {
    return
  }

  refunding.value = true
  try {
    // 让用户明确感知到提交过程
    await new Promise((resolve) => window.setTimeout(resolve, 2000))
    const proof = await user.signAction()
    const res = await refundPledgeOrder({
      orderNo: refundOrder.value.id,
      ...proof,
    })
    showRefundModal.value = false
    toast(`提前退本成功，到账 ${formatAmount(res.amount)} USDT`, 'success')
    await Promise.all([loadOrders(true), loadStatistics()])
  } catch (error) {
    if (!(error instanceof ApiError) && error instanceof Error) {
      toast(error.message, 'info')
    }
  } finally {
    refunding.value = false
  }
}
</script>

<template>
  <div
    class="flex flex-1 flex-col px-4 pb-[calc(1.5rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))]"
  >
    <AppHeader />

    <PageHeading
      title="报单明细"
      subtitle="每笔留痕 · 进度可查 · 台账清晰"
      @back="goBack"
    />

    <section
      class="relative mt-3 min-h-[172px] overflow-hidden rounded-2xl border border-[#e8d3a9] p-4 max-[359px]:min-h-[156px]"
    >
      <img
        :src="bdBg"
        alt=""
        class="pointer-events-none absolute inset-0 h-full w-full object-cover"
      />
      <div class="relative z-10 flex h-full flex-col justify-between">
        <div>
          <div class="flex items-center gap-2">
            <img :src="icon5" alt="" class="h-6 w-8 object-contain" />
            <span class="text-[17px] font-bold text-[#2a1c10]">入金总额</span>
          </div>
          <div class="mt-1 flex items-baseline gap-1.5 text-[#a45f1a]">
            <span class="text-[34px] font-black leading-none tracking-tight tabular-nums max-[359px]:text-[30px]">
              {{ formatAmount(statistics.total,2) }}
            </span>
            <span class="text-[12px] font-bold">USDT</span>
          </div>
        </div>

        <div class="mt-3 grid grid-cols-2 border-t border-[#e3cda0]/80 pt-2">
          <div>
            <p class="text-[12.5px] text-ink-500">今日释放</p>
            <p
              class="mt-0.5 whitespace-nowrap text-[17px] font-bold tabular-nums text-[#1f1a13] max-[359px]:text-[15px]"
            >
              +{{ formatAmount(statistics.todayRelease) }}
              <span class="text-[11px] font-medium text-ink-500 max-[359px]:text-[10px]">USDT</span>
            </p>
          </div>
          <div class="border-l border-[#e3cda0]/80 pl-3">
            <p class="text-[12.5px] text-ink-500">累计释放</p>
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
      <div v-if="loading" class="py-8 text-center text-[13px] text-ink-500">加载中...</div>

      <div
        v-else-if="!filteredOrders.length"
        class="rounded-2xl border border-dashed border-[#e5d5b7] bg-[#fffdf8]/70 py-10 text-center text-[13px] text-ink-500"
      >
        暂无报单记录
      </div>

      <article
        v-for="order in filteredOrders"
        :key="order.id"
        class="rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-3.5 py-3 shadow-gold-soft"
      >
        <div class="flex items-center gap-2">
          <img :src="icon5" alt="" class="h-6 w-8 shrink-0 object-contain" />
          <p class="min-w-0 flex-1 text-[17px] font-bold text-ink-900 max-[359px]:text-[16px]">
            金额：<span class="text-[#a45f1a]">{{ order.amount.toFixed(2) }} U</span>
          </p>
          <span
            class="shrink-0 rounded-full px-2.5 py-0.5 text-[11.5px] font-medium text-white"
            :class="
              order.status === 'active'
                ? 'bg-gradient-to-r from-[#e7c07a] to-[#c9954a]'
                : 'bg-[#a99a83]'
            "
          >
            {{ statusLabel(order.status) }}
          </span>
        </div>

        <dl class="mt-2.5 space-y-2 text-[14px]">
          <div class="flex items-center justify-between gap-3">
            <dt class="text-ink-500">今日释放</dt>
            <dd class="text-right font-medium text-ink-800">
              {{ formatAmount(order.todayRelease) }} U
            </dd>
          </div>
          <div class="flex items-center justify-between gap-3">
            <dt class="text-ink-500">累计释放</dt>
            <dd class="text-right font-medium text-ink-800">
              {{ formatAmount(order.totalRelease) }} U
            </dd>
          </div>
          <div class="flex items-center justify-between gap-3">
            <dt class="text-ink-500">投资周期</dt>
            <dd class="text-right font-medium text-ink-800">{{ order.period }} 个月</dd>
          </div>
          <div class="flex items-center justify-between gap-3">
            <dt class="text-ink-500">投资时间</dt>
            <dd class="text-right font-medium text-ink-800">{{ order.investDate }}</dd>
          </div>
        </dl>

        <template v-if="order.status === 'active'">
          <button
            type="button"
            class="mt-3 flex h-10 w-full items-center justify-center gap-2 rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[15px] font-semibold text-white shadow-gold-btn transition active:scale-[0.99]"
            @click="requestRefund(order)"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.7"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-5 w-5"
              aria-hidden="true"
            >
              <circle cx="12" cy="12" r="8.5" />
              <path d="m8.5 12 2.5 2.5 4.5-5" />
            </svg>
            申请提前退本
          </button>

          <div
            class="mt-2 flex items-start gap-1.5 rounded-xl bg-[#fbf4e4]/90 px-2.5 py-2 text-[11.5px] leading-4 text-[#8a7456]"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.7"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="mt-0.5 h-3.5 w-3.5 shrink-0 text-[#c69b4f]"
              aria-hidden="true"
            >
              <path d="M9 18h6M10 21h4M12 3a6 6 0 0 0-3.5 10.9c.5.4.8 1 .8 1.6V16h5.4v-.5c0-.6.3-1.2.8-1.6A6 6 0 0 0 12 3Z" />
            </svg>
            <span>退本金额 = 入金金额 × {{ 100 - order.earlyExitPenaltyRate }}% - 累计已释放收益，请谨慎操作</span>
          </div>
        </template>
      </article>
    </div>

    <div ref="sentinel" class="h-1" aria-hidden="true" />

    <div v-if="loadingMore" class="flex items-center justify-center gap-2 py-4 text-[12px] text-ink-500">
      <span class="h-4 w-4 animate-spin rounded-full border-2 border-[#e3cda0] border-t-[#c9954a]" />
      加载中...
    </div>
    <div v-else-if="orders.length && !hasMore" class="py-4 text-center text-[12px] text-ink-400">
      没有更多了
    </div>
    <div v-else-if="hasMore" class="py-4 text-center text-[12px] text-ink-400">
      上拉加载更多
    </div>

    <EarlyRefundModal
      v-model:show="showRefundModal"
      :amount="refundOrder?.amount ?? 0"
      :released="refundOrder?.totalRelease ?? 0"
      :penalty-rate="refundOrder?.earlyExitPenaltyRate ?? 0"
      :submitting="refunding"
      @confirm="handleRefundConfirm"
    />
  </div>
</template>

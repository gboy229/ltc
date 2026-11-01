<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import banner from '@/assets/images/banner.jpg'
import icon5 from '@/assets/images/icon5.png'
import AppHeader from '@/components/AppHeader.vue'
import BuyModal from '@/components/BuyModal.vue'
import NoticeTicker from '@/components/NoticeTicker.vue'
import { getRecommendArticles, type CmsArticleSummary } from '@/api/cms'
import { getOrderList, getPlans, type PledgePlan } from '@/api/pledge'
import { formatAmount } from '@/utils/format'
import { toast } from '@/utils/toast'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const user = useUserStore()

const plans = ref<PledgePlan[]>([])
const homeNotices = ref<CmsArticleSummary[]>([])
const loading = ref(true)

const selectedPlanName = ref('')
const selectedPeriod = ref(0)
const showBuyModal = ref(false)

const recentOrders = ref<Array<{ id: string; amount: number; period: number; time: string }>>([])

interface PledgeTier {
  name: string
  amount: number
}

/** 当前周期下的套餐 */
const periodPlans = computed(() =>
  plans.value.filter(
    (plan) => !selectedPeriod.value || plan.periodMonths === selectedPeriod.value,
  ),
)

/** 挡位由当前周期的套餐反推：没有该周期套餐的挡位不展示 */
const tiers = computed<PledgeTier[]>(() => {
  const map = new Map<string, PledgeTier>()
  for (const plan of periodPlans.value) {
    if (!map.has(plan.name)) {
      map.set(plan.name, {
        name: plan.name,
        amount: plan.amount,
      })
    }
  }
  return Array.from(map.values())
})

const periods = computed(() =>
  Array.from(new Set(plans.value.map((plan) => plan.periodMonths))).sort((a, b) => a - b),
)

const maxMonthlyRate = computed(() => {
  if (!periodPlans.value.length) {
    return 0
  }
  return Math.max(...periodPlans.value.map((plan) => plan.monthlyRewardRate))
})

const selectedPlan = computed(
  () =>
    plans.value.find(
      (plan) =>
        plan.name === selectedPlanName.value &&
        plan.periodMonths === selectedPeriod.value,
    ) ?? null,
)

const buyTotalReturn = computed(() => selectedPlan.value?.totalReturn ?? 0)

function formatRate(rate: number) {
  return Number.isInteger(rate)
    ? rate.toFixed(0)
    : rate.toFixed(2).replace(/0+$/, '').replace(/\.$/, '')
}

function pickDefaults() {
  const availablePeriods = periods.value
  if (!availablePeriods.includes(selectedPeriod.value)) {
    selectedPeriod.value = availablePeriods.includes(3) ? 3 : (availablePeriods[0] ?? 0)
  }

  // 周期确定后，再从该周期可用挡位中选择
  const tierList = tiers.value
  if (!tierList.some((tier) => tier.name === selectedPlanName.value)) {
    const preferred =
      tierList.find((tier) => tier.name.toUpperCase() === 'V3') ?? tierList[0] ?? null
    selectedPlanName.value = preferred?.name ?? ''
  }
}

async function loadRecentOrders() {
  if (!user.isLoggedIn) {
    recentOrders.value = []
    return
  }
  try {
    const res = await getOrderList({ page: 1, pageSize: 5 })
    recentOrders.value = (res.list ?? []).map((order) => ({
      id: order.id,
      amount: order.amount,
      period: order.period,
      time: order.investDate || '',
    }))
  } catch {
    recentOrders.value = []
  }
}

async function loadHome() {
  loading.value = true

  const [planResult, articleResult] = await Promise.allSettled([
    getPlans(),
    getRecommendArticles(),
  ])

  if (planResult.status === 'fulfilled') {
    plans.value = planResult.value.list ?? []
  }
  if (articleResult.status === 'fulfilled') {
    homeNotices.value = articleResult.value.list ?? []
  }

  await loadRecentOrders()
  pickDefaults()
  loading.value = false
}

function openNotice(notice: CmsArticleSummary) {
  void router.push(`/notice/${notice.id}`)
}

function selectTier(name: string) {
  selectedPlanName.value = name

  if (
    plans.value.length &&
    !plans.value.some(
      (plan) => plan.name === name && plan.periodMonths === selectedPeriod.value,
    )
  ) {
    const fallback = plans.value.find((plan) => plan.name === name)
    if (fallback) {
      selectedPeriod.value = fallback.periodMonths
    }
  }
}

function selectPeriod(period: number) {
  const previousTiers = tiers.value
  const currentIndex = previousTiers.findIndex(
    (tier) => tier.name === selectedPlanName.value,
  )

  selectedPeriod.value = period

  // 新周期下当前挡位不可用时，优先切到相邻位置的可用挡位
  const nextTiers = tiers.value
  if (!nextTiers.some((tier) => tier.name === selectedPlanName.value)) {
    const index = currentIndex >= 0 ? Math.min(currentIndex, nextTiers.length - 1) : 0
    selectedPlanName.value = nextTiers[index]?.name ?? ''
  }
}

function openBuy() {
  if (!selectedPlan.value) {
    toast('当前方案暂不可用，请稍后再试', 'info')
    return
  }
  showBuyModal.value = true
}


onMounted(() => {
  void loadHome()
})

watch(
  () => user.isLoggedIn,
  () => {
    void loadRecentOrders()
  },
)
</script>

<template>
  <div
    class="flex flex-1 flex-col px-4 pb-[calc(1.25rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))]"
  >
    <AppHeader />

    <img
      :src="banner"
      alt="黄金AI量化交易模型，AI驱动 · 智能交易 · 稳健增值"
      class="mt-2 block w-full rounded-xl border border-[#ead6ac] object-cover shadow-[0_6px_18px_rgba(166,120,44,0.1)]"
    />

    <NoticeTicker :notices="homeNotices" @select="openNotice" />

    <section
      class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 p-3.5 shadow-gold-soft backdrop-blur-[2px]"
    >
      <div class="flex items-center gap-2">
        <img :src="icon5" alt="" class="h-6 w-8 object-contain" />
        <h2 class="text-[18px] font-bold tracking-tight text-ink-900">入金报单</h2>
        <span
          class="ml-1 inline-flex h-6 items-center rounded-full bg-gradient-to-r from-[#edcd92] to-[#c89045] px-2.5 text-[12px] font-medium text-white shadow-sm"
        >
          最低1,000 USDT
        </span>
      </div>

      <div
        class="mt-2.5 grid grid-cols-3 rounded-xl bg-gradient-to-r from-[#fbf4e3] via-[#fefbf4] to-[#fbf4e3] py-2.5"
      >
        <div class="flex flex-col items-center justify-center">
          <p class="text-[22px] font-bold leading-none text-[#9c6528]">
            {{ tiers.length }}<span class="ml-0.5 text-[13px] font-semibold">挡</span>
          </p>
          <p class="mt-1 text-[12px] leading-none text-ink-500">投资挡位</p>
        </div>
        <div class="flex flex-col items-center justify-center border-l border-[#e2cda2]">
          <p class="text-[22px] font-bold leading-none text-[#9c6528]">
            {{ periods.length }}<span class="ml-0.5 text-[13px] font-semibold">期</span>
          </p>
          <p class="mt-1 text-[12px] leading-none text-ink-500">投资周期</p>
        </div>
        <div class="flex flex-col items-center justify-center border-l border-[#e2cda2]">
          <p class="text-[22px] font-bold leading-none text-[#9c6528]">
            {{ formatRate(maxMonthlyRate) }}<span class="ml-0.5 text-[13px] font-semibold">%</span>
          </p>
          <p class="mt-1 text-[12px] leading-none text-ink-500">最高月化</p>
        </div>
      </div>

      <p class="mt-2.5 text-[13.5px] font-semibold text-ink-700">选择挡位</p>
      <div v-if="tiers.length" class="mt-1.5 grid grid-cols-5 gap-1.5 max-[359px]:gap-1">
        <button
          v-for="tier in tiers"
          :key="tier.name"
          type="button"
          class="flex h-9 flex-col items-center justify-center rounded-lg border text-center transition active:scale-[0.98]"
          :class="
            selectedPlanName === tier.name
              ? 'border-[#d5a45a] bg-gradient-to-br from-[#f4d89c] via-[#e1b064] to-[#c08a40] text-white shadow-[0_2px_6px_rgba(177,124,46,0.28)]'
              : 'border-[#e5d5b7] bg-[#fffdf9] text-ink-800'
          "
          @click="selectTier(tier.name)"
        >
          <span class="text-[12px] font-semibold leading-none">{{ tier.name }}</span>
          <span
            class="mt-0.5 text-[12px] leading-none tabular-nums max-[359px]:text-[11px]"
            :class="selectedPlanName === tier.name ? 'text-white/90' : 'text-ink-500'"
          >
            {{ tier.amount.toFixed(2) }}
          </span>
        </button>
      </div>
      <p v-else class="mt-1.5 rounded-lg border border-[#e5d5b7] py-3 text-center text-[12.5px] text-ink-500">
        {{ loading ? '方案加载中...' : '暂无可选挡位' }}
      </p>

      <p class="mt-2.5 text-[13.5px] font-semibold text-ink-700">选择周期</p>
      <div v-if="periods.length" class="mt-1.5 grid grid-cols-3 gap-2">
        <button
          v-for="period in periods"
          :key="period"
          type="button"
          class="h-9 rounded-full border text-[13px] font-medium transition active:scale-[0.98]"
          :class="
            selectedPeriod === period
              ? 'border-[#cf9a4a] bg-gradient-to-r from-[#e1b26a] to-[#c89045] font-semibold text-white shadow-sm'
              : 'border-[#e5d5b7] bg-[#fffdf9] text-ink-700'
          "
          @click="selectPeriod(period)"
        >
          {{ period }}个月
        </button>
      </div>
      <p v-else class="mt-1.5 rounded-lg border border-[#e5d5b7] py-3 text-center text-[12.5px] text-ink-500">
        {{ loading ? '周期加载中...' : '暂无可选周期' }}
      </p>

      <button
        type="button"
        class="mt-3 flex h-11 w-full items-center justify-center gap-1.5 rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[17px] font-semibold text-white transition"
        :class="selectedPlan ? 'shadow-gold-btn active:scale-[0.99]' : 'opacity-60'"
        :disabled="!selectedPlan"
        @click="openBuy"
      >
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.6"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="h-[18px] w-[18px]"
          aria-hidden="true"
        >
          <path
            d="M4 7.5A2.5 2.5 0 0 1 6.5 5h11A2.5 2.5 0 0 1 20 7.5v1.7a2.8 2.8 0 0 0 0 5.6v1.7a2.5 2.5 0 0 1-2.5 2.5h-11A2.5 2.5 0 0 1 4 16.5v-1.7a2.8 2.8 0 0 0 0-5.6V7.5Z"
          />
          <path d="M9 9h2.5M9 12h4M9 15h2.5" />
        </svg>
        立即报单
      </button>
    </section>

    <section
      id="order-details"
      class="mt-3 scroll-mt-4 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 p-3.5 shadow-gold-soft backdrop-blur-[2px]"
    >
      <div class="flex items-center gap-2">
        <img :src="icon5" alt="" class="h-6 w-8 object-contain" />
        <h2 class="text-[18px] font-bold tracking-tight text-ink-900">报单明细</h2>
      </div>

      <div v-if="recentOrders.length" class="mt-1 divide-y divide-[#efe2c8]">
        <div
          v-for="record in recentOrders"
          :key="record.id"
          class="flex items-center gap-2 py-2 max-[359px]:gap-1.5"
        >
          <span
            class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-[#fbf4e4] max-[359px]:h-6 max-[359px]:w-6"
          >
            <img
              :src="icon5"
              alt=""
              class="h-5 w-7 object-contain max-[359px]:h-4 max-[359px]:w-6"
            />
          </span>
          <div class="min-w-0 flex-1">
            <p class="text-[14px] font-medium leading-[18px] text-ink-800 max-[359px]:text-[13.5px]">
              入金金额
            </p>
            <p class="text-[12px] leading-4 text-ink-500 max-[359px]:text-[11.5px]">
              投资周期：{{ record.period }}个月
            </p>
          </div>
          <div class="shrink-0 text-right">
            <p class="text-[17px] font-bold leading-[19px] text-[#a45e1a] max-[359px]:text-[16px]">
              {{ formatAmount(record.amount) }} U
            </p>
            <p class="text-[12px] leading-4 text-ink-500 max-[359px]:text-[11px]">
              时间：{{ record.time }}
            </p>
          </div>
        </div>
      </div>

      <div v-else class="py-6 text-center">
        <p class="text-[13px] text-ink-500">暂无报单记录</p>
        <p class="mt-1 text-[11.5px] text-ink-400">完成首笔入金后，报单进度将在这里展示</p>
      </div>

      <button
        v-if="recentOrders.length"
        type="button"
        class="mt-2 flex h-9 w-full items-center justify-center gap-1 rounded-full border border-[#d7b57b] bg-[#fffdf8]/60 text-[13.5px] font-medium text-[#a97227] transition active:bg-[#fbf2df]"
        @click="router.push('/orders')"
      >
        查看更多
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
    </section>

    <BuyModal
      v-model:show="showBuyModal"
      :tier="selectedPlan?.name ?? ''"
      :amount="selectedPlan?.amount ?? 0"
      :period="selectedPlan?.periodMonths ?? 0"
      :monthly-rate="selectedPlan?.monthlyRewardInfo"
      :total-return="buyTotalReturn"
    />
  </div>
</template>

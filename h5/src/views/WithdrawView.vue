<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { formatAmount } from '@/utils/format'
import { useRouter } from 'vue-router'

import icon5 from '@/assets/images/icon5.png'
import wBg from '@/assets/images/w_bg1.png'
import AppHeader from '@/components/AppHeader.vue'
import PageHeading from '@/components/PageHeading.vue'
import { getAssetAccount } from '@/api/asset'
import {
  applyWithdrawal,
  getWithdrawalList,
  getWithdrawalOptions,
  type WithdrawalOption,
  type WithdrawalRecord,
  type WithdrawalStatus,
} from '@/api/chain/withdrawal'
import { toast } from '@/utils/toast'
import { ApiError } from '@/utils/http'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const user = useUserStore()

const options = ref<WithdrawalOption[]>([])
const selectedChain = ref('')
const amountInput = ref('')
const available = ref(0)
const records = ref<WithdrawalRecord[]>([])
const submitting = ref(false)
const loadingRecords = ref(false)
const showConfirm = ref(false)
const confirmInfo = ref<{
  chain: string
  network: string
  symbol: string
  amount: number
  fee: number
  receiveAmount: number
} | null>(null)
const showSuccess = ref(false)
const successInfo = ref<{
  network: string
  symbol: string
  amount: number
  fee: number
  receiveAmount: number
} | null>(null)
const recordsSection = ref<HTMLElement | null>(null)

const currentOption = computed(
  () => options.value.find((item) => item.chain === selectedChain.value) ?? options.value[0] ?? null,
)

const minimum = computed(() => {
  const value = currentOption.value?.minAmount ?? 0
  return value > 0 ? value : 1
})

const numericAmount = computed(() => {
  const value = Number.parseFloat(amountInput.value)
  return Number.isFinite(value) && value > 0 ? value : 0
})

const fee = computed(() => {
  const option = currentOption.value
  if (!option) {
    return 0
  }
  if (option.feeType === 0) {
    return Math.round(numericAmount.value * option.rate) / 100
  }
  return option.rate
})

const receivedAmount = computed(() => Math.max(numericAmount.value - fee.value, 0))

function statusLabel(status: WithdrawalStatus) {
  if (status === 'success') {
    return '成功'
  }
  if (status === 'failed') {
    return '失败'
  }
  return '处理中'
}

function statusColor(record: WithdrawalRecord) {
  if (record.status === 'success') {
    return 'text-[#3a9e5a]'
  }
  if (record.status === 'failed') {
    return 'text-[#8c8c8c]'
  }
  return 'text-[#c08a2f]'
}

function onAmountInput(event: Event) {
  const input = event.target as HTMLInputElement
  const cleaned = input.value.replace(/\D/g, '')
  amountInput.value = cleaned
  input.value = cleaned
}

async function loadOptions() {
  try {
    const res = await getWithdrawalOptions()
    options.value = res.list ?? []
    if (!options.value.some((item) => item.chain === selectedChain.value)) {
      selectedChain.value = options.value[0]?.chain ?? ''
    }
  } catch {
    options.value = []
  }
}

async function loadAccount() {
  try {
    const account = await getAssetAccount()
    available.value = account.available ?? 0
  } catch {
    available.value = 0
  }
}

async function loadRecords() {
  loadingRecords.value = true
  try {
    const res = await getWithdrawalList({ page: 1, pageSize: 5 })
    records.value = res.list ?? []
  } catch {
    records.value = []
  } finally {
    loadingRecords.value = false
  }
}

function confirmWithdraw() {
  const option = currentOption.value
  if (!option) {
    toast('暂未配置可用的提币网络', 'error')
    return
  }
  if (numericAmount.value < minimum.value) {
    toast(`最少提币 ${formatAmount(minimum.value)} ${option.symbol}`, 'error')
    return
  }
  if (option.maxAmount > 0 && numericAmount.value > option.maxAmount) {
    toast(`单笔最多提币 ${formatAmount(option.maxAmount)} ${option.symbol}`, 'error')
    return
  }
  if (numericAmount.value > available.value) {
    toast('可用余额不足', 'error')
    return
  }

  confirmInfo.value = {
    chain: option.chain,
    network: option.network,
    symbol: option.symbol,
    amount: numericAmount.value,
    fee: fee.value,
    receiveAmount: receivedAmount.value,
  }
  showConfirm.value = true
}

function closeConfirm() {
  if (submitting.value) {
    return
  }
  showConfirm.value = false
}

async function submitWithdraw() {
  const info = confirmInfo.value
  if (!info || submitting.value) {
    return
  }

  submitting.value = true
  try {
    // 让用户明确感知到提交过程
    await new Promise((resolve) => window.setTimeout(resolve, 2000))
    const proof = await user.signAction()
    await applyWithdrawal({
      chain: info.chain,
      symbol: info.symbol,
      amount: info.amount,
      ...proof,
    })
    successInfo.value = {
      network: info.network,
      symbol: info.symbol,
      amount: info.amount,
      fee: info.fee,
      receiveAmount: info.receiveAmount,
    }
    amountInput.value = ''
    await Promise.all([loadAccount(), loadRecords()])
    showConfirm.value = false
    showSuccess.value = true
  } catch (error) {
    if (!(error instanceof ApiError) && error instanceof Error) {
      toast(error.message, 'info')
    }
  } finally {
    submitting.value = false
  }
}

function goBack() {
  const state = window.history.state as { back?: string | null } | null
  if (state?.back) {
    router.back()
    return
  }
  void router.push('/my-assets')
}

function closeSuccess() {
  showSuccess.value = false
}

async function viewRecords() {
  showSuccess.value = false
  await nextTick()
  recordsSection.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

function viewAllRecords() {
  void router.push('/withdrawals')
}

onMounted(() => {
  void loadOptions()
  void loadAccount()
  void loadRecords()
})
</script>

<template>
  <div
    class="flex flex-1 flex-col px-4 pb-[calc(1.5rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))]"
  >
    <AppHeader />

    <PageHeading
      title="在线提币"
      subtitle="链上转账 · 极速放行 · 资产直达"
      @back="goBack"
    />

    <section
      class="relative mt-3 min-h-[92px] overflow-hidden rounded-2xl max-[359px]:min-h-[88px]"
    >
      <img
        :src="wBg"
        alt=""
        class="pointer-events-none absolute inset-0 h-full w-full object-cover object-left"
      />
      <div class="relative z-10 flex h-full min-h-[92px] flex-col justify-center pl-[36%] pr-4 max-[359px]:min-h-[88px]">
        <p class="text-[16px] font-medium text-[#6b512f]">USDT 可用余额</p>
        <p
          class="mt-0.5 whitespace-nowrap text-[34px] font-black leading-none tracking-tight text-[#a45f1a] max-[359px]:text-[28px]"
        >
          {{ formatAmount(available) }}
          <span class="text-[12px] font-bold">USDT</span>
        </p>
      </div>
    </section>

    <section
      class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 p-3.5 shadow-gold-soft"
    >
      <div class="flex items-center gap-2">
        <img :src="icon5" alt="" class="h-6 w-7 object-contain" />
        <h2 class="text-[17px] font-bold text-ink-900">选择链</h2>
      </div>

      <div class="mt-2.5 grid grid-cols-3 gap-2">
        <button
          v-for="chain in options"
          :key="chain.chain"
          type="button"
          class="flex h-9 items-center justify-center gap-1.5 rounded-full border text-[12.5px] font-medium transition"
          :class="
            selectedChain === chain.chain
              ? 'border-[#d2a25a] bg-gradient-to-r from-[#e6bf7f] via-[#d5a65c] to-[#efd09a] font-semibold text-white shadow-sm'
              : 'border-[#e5d5b7] bg-[#fffdf8] text-ink-700'
          "
          @click="selectedChain = chain.chain"
        >
          <svg
            v-if="chain.network === 'TRC20'"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-3.5 w-3.5"
            aria-hidden="true"
          >
            <path d="m21 3-7.5 18-3-7.5L3 10.5 21 3Z" />
            <path d="m21 3-10.5 10.5" />
          </svg>
          <svg
            v-else-if="chain.network === 'ERC20'"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linejoin="round"
            class="h-3.5 w-3.5"
            aria-hidden="true"
          >
            <path d="M12 2 5.5 12 12 16l6.5-4L12 2Z" />
            <path d="M5.5 13.5 12 22l6.5-8.5L12 17.5l-6.5-4Z" />
          </svg>
          <svg
            v-else
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linejoin="round"
            class="h-3.5 w-3.5"
            aria-hidden="true"
          >
            <path d="m12 2 4 4-4 4-4-4 4-4Z" />
            <path d="m12 14 4 4-4 4-4-4 4-4Z" />
            <path d="m18 8 4 4-4 4-4-4 4-4Z" />
            <path d="m6 8 4 4-4 4-4-4 4-4Z" />
          </svg>
          {{ chain.network }}
        </button>
      </div>

      <div class="mt-2 flex items-start gap-1.5 text-[11px] leading-4 text-[#d14a3a]">
        <svg
          viewBox="0 0 24 24"
          fill="currentColor"
          class="mt-0.5 h-3.5 w-3.5 shrink-0"
          aria-hidden="true"
        >
          <path d="M12 2 1 21h22L12 2Zm0 4.4L19.3 19H4.7L12 6.4ZM11 10v4h2v-4h-2Zm0 5v2h2v-2h-2Z" />
        </svg>
        <span>请确保提币网络与接收方网络一致，否则资产可能丢失</span>
      </div>

      <div class="mt-3 flex items-center gap-2">
        <img :src="icon5" alt="" class="h-5 w-6 object-contain" />
        <h2 class="text-[17px] font-bold text-ink-900">提币数量</h2>
      </div>

      <div
        class="mt-2 flex h-12 items-center rounded-xl border border-[#e7d2aa] bg-[#fffdf8] px-3"
      >
        <input
          :value="amountInput"
          type="text"
          inputmode="numeric"
          :placeholder="`最少 ${formatAmount(minimum)} USDT`"
          class="min-w-0 flex-1 bg-transparent text-[15px] text-ink-900 outline-none placeholder:text-[#b3a38c]"
          @input="onAmountInput"
        />
        <span class="shrink-0 text-[12px] font-medium text-ink-500">USDT</span>
      </div>

      <div class="mt-2 flex items-center justify-between gap-3 text-[14px]">
        <div class="flex items-center gap-1.5">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.7"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-4 w-4 text-[#c09a55]"
            aria-hidden="true"
          >
            <circle cx="12" cy="12" r="9" />
            <path d="m8 12 2.5 2.5L16 9" />
          </svg>
          <span class="text-ink-500">到账数量</span>
        </div>
        <span class="whitespace-nowrap font-bold tabular-nums text-[#a45f1a]">
          {{ formatAmount(receivedAmount) }} USDT
        </span>
      </div>

      <div class="mt-2 flex items-center justify-between gap-3 text-[14px]">
        <div class="flex items-center gap-1.5">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.7"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-4 w-4 text-[#c09a55]"
            aria-hidden="true"
          >
            <circle cx="12" cy="12" r="7" />
            <path d="M12 8v8M9.5 10h5M9.5 14h5" />
          </svg>
          <span class="text-ink-500">网络手续费</span>
        </div>
        <span class="whitespace-nowrap font-bold tabular-nums text-[#d14a3a]">
          {{ formatAmount(fee) }} USDT
        </span>
      </div>

      <button
        type="button"
        class="mt-3 flex h-12 w-full items-center justify-center gap-2 rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[16px] font-semibold text-white shadow-gold-btn transition active:scale-[0.99] disabled:opacity-60"
        :disabled="submitting"
        @click="confirmWithdraw"
      >
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.8"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="h-5 w-5"
          aria-hidden="true"
        >
          <circle cx="12" cy="12" r="8.5" />
          <path d="m8.5 12 2.5 2.5 4.5-5" />
        </svg>
        {{ submitting ? '提交中...' : '确认提币' }}
      </button>
    </section>

    <section
      ref="recordsSection"
      class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 p-3.5 shadow-gold-soft"
    >
      <div class="flex items-center gap-2">
        <img :src="icon5" alt="" class="h-6 w-7 object-contain" />
        <h2 class="text-[17px] font-bold text-ink-900">提币记录</h2>
        <button
          type="button"
          class="ml-auto flex items-center gap-0.5 text-[12.5px] font-medium text-[#a97227]"
          @click="viewAllRecords"
        >
          查看全部
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
      </div>

      <div class="mt-1 divide-y divide-[#f1e5cc]">
        <div v-if="loadingRecords" class="py-8 text-center text-[13px] text-ink-500">加载中...</div>

        <div v-else-if="!records.length" class="py-8 text-center text-[13px] text-ink-500">
          暂无提币记录
        </div>

        <div
          v-for="record in records"
          :key="record.id"
          class="flex items-center gap-1.5 py-2.5 max-[359px]:gap-1"
        >
          <div class="flex w-[68px] shrink-0 items-center gap-1 max-[359px]:w-[62px]">
            <span
              class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full max-[359px]:h-5 max-[359px]:w-5"
              :class="
                record.status === 'success'
                  ? 'bg-[#3a9e5a] text-white'
                  : record.status === 'failed'
                    ? 'bg-[#9a9a9a] text-white'
                    : 'bg-[#f4e3bd] text-[#b47e1f]'
              "
            >
              <svg
                v-if="record.status === 'success'"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-3.5 w-3.5"
                aria-hidden="true"
              >
                <path d="m5 12 4.5 4.5L19 7" />
              </svg>
              <svg
                v-else-if="record.status === 'failed'"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.2"
                stroke-linecap="round"
                class="h-3.5 w-3.5"
                aria-hidden="true"
              >
                <path d="m7 7 10 10M17 7 7 17" />
              </svg>
              <svg
                v-else
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-3.5 w-3.5"
                aria-hidden="true"
              >
                <circle cx="12" cy="12" r="8.5" />
                <path d="M12 7v5l3 2" />
              </svg>
            </span>
            <span
              class="whitespace-nowrap text-[11.5px] font-medium max-[359px]:text-[11px]"
              :class="statusColor(record)"
            >
              {{ statusLabel(record.status) }}
            </span>
          </div>

          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-1">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-3 w-3 shrink-0 text-[#a99a83] max-[359px]:h-2.5 max-[359px]:w-2.5"
                aria-hidden="true"
              >
                <path d="M10.5 13.5a4 4 0 0 0 5.7 0l2.3-2.3a4 4 0 0 0-5.7-5.7l-1.3 1.3" />
                <path d="M13.5 10.5a4 4 0 0 0-5.7 0l-2.3 2.3a4 4 0 0 0 5.7 5.7l1.3-1.3" />
              </svg>
              <span class="truncate text-[12px] text-ink-700 max-[359px]:text-[11.5px]">{{ record.address }}</span>
            </div>
            <div class="mt-0.5 flex items-center gap-1">
              <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.8"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="h-3 w-3 shrink-0 text-[#a99a83] max-[359px]:h-2.5 max-[359px]:w-2.5"
                aria-hidden="true"
              >
                <circle cx="12" cy="12" r="8.5" />
                <path d="M12 7v5l3 2" />
              </svg>
              <span class="whitespace-nowrap text-[11px] text-ink-500">{{ record.date }}</span>
            </div>
          </div>

          <span
            class="shrink-0 whitespace-nowrap text-[14px] font-semibold tabular-nums max-[359px]:text-[12.5px]"
            :class="statusColor(record)"
          >
            -{{ formatAmount(record.amount) }} USDT
          </span>
        </div>
      </div>
    </section>

    <Transition
      enter-active-class="transition-opacity duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="showConfirm && confirmInfo"
        class="fixed inset-0 z-[135] flex items-center justify-center px-5"
        role="dialog"
        aria-modal="true"
        aria-label="确认提币"
      >
        <div class="absolute inset-0 bg-[#3b2f22]/55" @click="closeConfirm" />

        <div
          class="relative w-full max-w-[320px] rounded-[24px] border border-[#e7d2a8] bg-[#fffdf8] px-5 pb-5 pt-6 text-center shadow-[0_24px_60px_-20px_rgba(80,50,12,0.45)]"
        >
          <h2 class="text-[22px] font-black tracking-tight text-ink-900">确认提币？</h2>
          <p class="mt-1 text-[12.5px] leading-5 text-ink-500">
            请确认提币信息，确认后将提交后台审核
          </p>

          <div class="mt-4 rounded-2xl bg-[#fbf5e8] px-3.5 py-3 text-left text-[13.5px]">
            <div class="flex items-center justify-between py-1">
              <span class="text-ink-500">提币网络</span>
              <span class="font-semibold text-ink-900">{{ confirmInfo.network }}</span>
            </div>
            <div class="flex items-center justify-between py-1">
              <span class="text-ink-500">提币数量</span>
              <span class="font-semibold tabular-nums text-ink-900">
                {{ formatAmount(confirmInfo.amount) }} {{ confirmInfo.symbol }}
              </span>
            </div>
            <div class="flex items-center justify-between py-1">
              <span class="text-ink-500">手续费</span>
              <span class="font-semibold tabular-nums text-[#d14a3a]">
                -{{ formatAmount(confirmInfo.fee) }} {{ confirmInfo.symbol }}
              </span>
            </div>
            <div class="mt-1 flex items-end justify-between border-t border-[#eadcc0] pt-2.5">
              <span class="text-[15px] font-bold text-ink-900">预计到账</span>
              <span class="text-[22px] font-black tabular-nums text-[#a45f1a]">
                {{ formatAmount(confirmInfo.receiveAmount) }}
                <span class="text-[11px] font-bold">USDT</span>
              </span>
            </div>
          </div>

          <div class="mt-3 grid grid-cols-2 gap-3">
            <button
              type="button"
              class="h-11 rounded-full border border-[#e7d2aa] bg-[#fdf8ed] text-[15px] font-semibold text-[#5a4a35] transition active:scale-[0.98] disabled:opacity-60"
              :disabled="submitting"
              @click="closeConfirm"
            >
              取消
            </button>
            <button
              type="button"
              class="h-11 rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[15px] font-semibold text-white shadow-gold-btn transition active:scale-[0.98] disabled:opacity-60"
              :disabled="submitting"
              @click="submitWithdraw"
            >
              <span v-if="submitting" class="flex items-center justify-center gap-2">
                <span
                  class="h-4 w-4 animate-spin rounded-full border-2 border-white/40 border-t-white"
                />
                提交中...
              </span>
              <span v-else>确认提交</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <Transition
      enter-active-class="transition-opacity duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="showSuccess && successInfo"
        class="fixed inset-0 z-[140] flex items-center justify-center px-5"
        role="dialog"
        aria-modal="true"
        aria-label="提币申请已提交"
      >
        <div class="absolute inset-0 bg-[#3b2f22]/55" @click="closeSuccess" />

        <div
          class="relative w-full max-w-[320px] rounded-[24px] border border-[#e7d2a8] bg-[#fffdf8] px-5 pb-5 pt-6 text-center shadow-[0_24px_60px_-20px_rgba(80,50,12,0.45)]"
        >
          <span
            class="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-[#eef8f0] text-[#3a9e5a]"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-8 w-8"
              aria-hidden="true"
            >
              <path d="m5 12 4.5 4.5L19 7" />
            </svg>
          </span>

          <h2 class="mt-3 text-[22px] font-black tracking-tight text-ink-900">
            提币申请已提交
          </h2>
          <div class="mt-4 rounded-2xl bg-[#fbf5e8] px-3.5 py-3 text-left text-[13.5px]">
            <div class="flex items-center justify-between py-1">
              <span class="text-ink-500">提币网络</span>
              <span class="font-semibold text-ink-900">{{ successInfo.network }}</span>
            </div>
            <div class="flex items-center justify-between py-1">
              <span class="text-ink-500">提币数量</span>
              <span class="font-semibold tabular-nums text-ink-900">
                {{ formatAmount(successInfo.amount) }} {{ successInfo.symbol }}
              </span>
            </div>
            <div class="flex items-center justify-between py-1">
              <span class="text-ink-500">手续费</span>
              <span class="font-semibold tabular-nums text-[#d14a3a]">
                -{{ formatAmount(successInfo.fee) }} {{ successInfo.symbol }}
              </span>
            </div>
            <div class="mt-1 flex items-end justify-between border-t border-[#eadcc0] pt-2.5">
              <span class="text-[15px] font-bold text-ink-900">预计到账</span>
              <span class="text-[22px] font-black tabular-nums text-[#a45f1a]">
                {{ formatAmount(successInfo.receiveAmount) }}
                <span class="text-[11px] font-bold">USDT</span>
              </span>
            </div>
          </div>

          <button
            type="button"
            class="mt-4 h-11 w-full rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[15px] font-semibold text-white shadow-gold-btn transition active:scale-[0.99]"
            @click="viewRecords"
          >
            查看提币记录
          </button>
          <button
            type="button"
            class="mt-2.5 h-10 w-full rounded-full border border-[#e7d2aa] bg-[#fdf8ed] text-[14px] font-semibold text-[#5a4a35] transition active:scale-[0.99]"
            @click="closeSuccess"
          >
            继续提币
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

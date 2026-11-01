<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'

import buyBg from '@/assets/images/buy_bg.png'
import icon6 from '@/assets/images/icon6.png'
import type { ChainConfig, ChainContract } from '@/api/chain'
import { buyPlan, getBuyOrderStatus, type BuyOrderStatusRes } from '@/api/pledge'
import { useChainStore } from '@/stores/chain'
import { useUserStore } from '@/stores/user'
import { formatAmount } from '@/utils/format'
import { toast } from '@/utils/toast'
import { ensureWalletChain, formatTokenAmount, getErc20Balance, sendErc20Transfer } from '@/web3/wallet'

const props = withDefaults(
  defineProps<{
    show: boolean
    tier: string
    amount: number
    period: number
    monthlyRate?: string
    totalReturn: number
  }>(),
  {
    monthlyRate: '8%-10%',
  },
)

const emit = defineEmits<{
  'update:show': [boolean]
}>()

type Phase = 'idle' | 'creating' | 'paying' | 'confirming' | 'success' | 'failed'

const user = useUserStore()
const chainStore = useChainStore()
const router = useRouter()
const cancelButton = ref<HTMLButtonElement | null>(null)
const phase = ref<Phase>('idle')
const orderNo = ref('')
const txHash = ref('')
const errorMessage = ref('')
const orderStatus = ref<BuyOrderStatusRes | null>(null)
const payAmount = ref(0)
const toAddress = ref('')
const balance = ref('--')
const balanceLoading = ref(false)
const chains = computed(() => chainStore.chains)
const selectedChain = ref<ChainConfig | null>(null)

let pollTimer: number | undefined

const busy = computed(
  () => phase.value === 'creating' || phase.value === 'paying' || phase.value === 'confirming',
)
const title = computed(() => {
  if (phase.value === 'success') return '支付成功'
  if (phase.value === 'confirming') return '支付确认中'
  return '确认入金'
})
const buttonText = computed(() => {
  switch (phase.value) {
    case 'creating':
      return '创建订单中...'
    case 'paying':
      return '请在钱包确认'
    case 'confirming':
      return '确认中...'
    case 'failed':
      return orderNo.value ? '重新支付' : '确认支付'
    default:
      return '确认支付'
  }
})
const displayAmount = computed(() => payAmount.value || props.amount)
const selectedContract = computed<ChainContract | null>(() =>
  selectedChain.value ? getPayContract(selectedChain.value) : null,
)
const walletOnSelectedChain = computed(() => {
  if (!selectedChain.value || !user.chainId) {
    return false
  }
  return Number.parseInt(user.chainId, 16) === selectedChain.value.chainId
})
const balanceText = computed(() => {
  if (balanceLoading.value) return '查询中...'
  if (!user.provider || !user.address) return '--'
  if (!walletOnSelectedChain.value) return '--'
  return balance.value
})

function syncSelectedChainFromBackend() {
  if (!chains.value.length) {
    selectedChain.value = null
    return
  }

  const current = selectedChain.value
  const matched = current
    ? chains.value.find((item) => item.chain === current.chain)
    : null
  selectedChain.value = matched ?? chains.value[0]
}

watch(chains, syncSelectedChainFromBackend, { immediate: true })

function stopPolling() {
  if (pollTimer !== undefined) {
    window.clearTimeout(pollTimer)
    pollTimer = undefined
  }
}

function startPolling() {
  stopPolling()
  if (!orderNo.value || !props.show) return
  phase.value = 'confirming'
  pollTimer = window.setTimeout(pollStatus, 3000)
}

async function pollStatus() {
  if (!orderNo.value) return

  try {
    const status = await getBuyOrderStatus(orderNo.value)
    orderStatus.value = status
    if (status.status === 1) {
      finishSuccess()
      return
    }
    if (status.status === 2) {
      phase.value = 'failed'
      errorMessage.value = '报单已被新报单替代，请重新提交'
      stopPolling()
      return
    }
  } catch {
    // 网络异常时保持轮询，由用户关闭弹窗停止
  }

  pollTimer = window.setTimeout(pollStatus, 3000)
}

function getPayContract(chain: ChainConfig): ChainContract | null {
  return chain.contracts.find((item) => item.symbol.toUpperCase() === 'USDT') ?? chain.contracts[0] ?? null
}

async function loadPayConfig() {
  if (!chains.value.length) {
    await chainStore.load()
  }
  if (!chains.value.length) {
    throw new Error('暂未配置支付链，请联系客服')
  }

  syncSelectedChainFromBackend()
  if (!selectedContract.value) {
    throw new Error('暂未配置支付合约，请联系客服')
  }
  return true
}

function ensureSelectedChain() {
  syncSelectedChainFromBackend()
}

async function refreshBalance() {
  const provider = user.provider
  const contract = selectedContract.value
  if (!props.show || !provider || !user.address || !contract || !walletOnSelectedChain.value) {
    balance.value = '--'
    return
  }

  balanceLoading.value = true
  try {
    const raw = await getErc20Balance(provider, user.address, contract.address)
    balance.value = formatTokenAmount(raw, contract.decimals, 2)
  } catch {
    balance.value = '--'
  } finally {
    balanceLoading.value = false
  }
}

let autoSwitchChainId = 0

async function ensureWalletOnBackendChain() {
  if (!props.show || !user.provider || !selectedChain.value || walletOnSelectedChain.value) {
    return
  }
  if (autoSwitchChainId === selectedChain.value.chainId) {
    return
  }

  autoSwitchChainId = selectedChain.value.chainId
  try {
    await ensureWalletChain(user.provider, selectedChain.value)
    await refreshBalance()
  } catch {
    balance.value = '--'
  }
}

watch(
  [() => props.show, () => user.address, () => user.chainId, selectedChain],
  () => {
    if (props.show) {
      void ensureWalletOnBackendChain()
    }
    void refreshBalance()
  },
  { immediate: true },
)

async function refreshChainsIfNeeded() {
  if (!user.chainId || !chainStore.chains.length) return
  const walletChainId = Number.parseInt(user.chainId, 16)
  if (!chainStore.chains.some((item) => item.chainId === walletChainId)) {
    await chainStore.load(true)
  }
}

async function selectChain(chain: ChainConfig) {
  if (busy.value) return
  selectedChain.value = chain
  errorMessage.value = ''

  if (!user.provider) {
    return
  }

  try {
    await ensureWalletChain(user.provider, chain)
    await refreshBalance()
  } catch (error) {
    phase.value = 'failed'
    errorMessage.value = error instanceof Error ? error.message : '切换钱包网络失败'
  }
}

async function refreshExistingOrder(): Promise<boolean> {
  if (!orderNo.value) {
    return false
  }
  try {
    const status = await getBuyOrderStatus(orderNo.value)
    orderStatus.value = status
    if (status.status === 1) {
      finishSuccess()
      return true
    }
    if (status.status === 2) {
      phase.value = 'failed'
      errorMessage.value = '报单已被新报单替代，请重新提交'
      stopPolling()
      return true
    }
  } catch {
    // 查询失败时保持原状态，由用户重试
  }
  return false
}

async function createOrReuseOrder() {
  if (orderNo.value) {
    return orderNo.value
  }

  const result = await buyPlan({
    planName: props.tier,
    periodMonths: props.period,
  })
  orderNo.value = result.orderNo
  payAmount.value = result.amount
  toAddress.value = result.toAddress
  return result.orderNo
}

async function payOrder(orderNumber: string) {
  await loadPayConfig()
  const chain = selectedChain.value
  const contract = selectedContract.value
  const provider = user.provider

  if (!chain || !contract) {
    throw new Error('支付配置异常')
  }
  if (!provider) {
    throw new Error('请先连接 Web3 钱包')
  }
  if (!user.address) {
    throw new Error('未获取到钱包地址')
  }
  if (!toAddress.value) {
    throw new Error('未获取到入金地址，请重新提交报单')
  }

  phase.value = 'paying'
  await ensureWalletChain(provider, chain)

  const amount = displayAmount.value.toFixed(2)
  const hash = await sendErc20Transfer(
    provider,
    user.address,
    contract.address,
    toAddress.value,
    amount,
    contract.decimals,
  )

  txHash.value = hash
  phase.value = 'confirming'
  startPolling()
  return orderNumber
}

async function confirmBuy() {
  if (busy.value) return

  errorMessage.value = ''

  if (orderNo.value && (await refreshExistingOrder())) {
    return
  }

  phase.value = orderNo.value ? 'paying' : 'creating'

  try {
    const loggedIn = await user.ensureLogin()
    if (!loggedIn) {
      throw new Error('请先连接钱包并登录')
    }

    await loadPayConfig()
    const currentOrderNo = await createOrReuseOrder()
    await payOrder(currentOrderNo)
  } catch (error) {
    phase.value = 'failed'
    errorMessage.value = error instanceof Error ? error.message : '支付失败，请重试'
  }
}

function finishSuccess() {
  stopPolling()
  phase.value = 'success'
  toast('报单成功', 'success')
  emit('update:show', false)
  orderNo.value = ''
  txHash.value = ''
  orderStatus.value = null
  payAmount.value = 0
  toAddress.value = ''
  void router.push('/orders')
}

function close() {
  stopPolling()
  if (phase.value === 'success') {
    phase.value = 'idle'
    orderNo.value = ''
    txHash.value = ''
    orderStatus.value = null
    payAmount.value = 0
    toAddress.value = ''
  }
  autoSwitchChainId = 0
  emit('update:show', false)
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    close()
  }
}

watch(
  () => props.show,
  async (show) => {
    document.body.classList.toggle('h5-buy-open', show)

    if (show) {
      window.addEventListener('keydown', onKeydown)
      if (orderNo.value && phase.value === 'confirming') {
        startPolling()
      }
      ensureSelectedChain()
      await refreshChainsIfNeeded()
      ensureSelectedChain()
      await ensureWalletOnBackendChain()
      await nextTick()
      cancelButton.value?.focus()
      return
    }

    stopPolling()
    window.removeEventListener('keydown', onKeydown)
  },
)

onBeforeUnmount(() => {
  stopPolling()
  document.body.classList.remove('h5-buy-open')
  window.removeEventListener('keydown', onKeydown)
})
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
      v-if="show"
      class="fixed inset-0 z-[120] overflow-y-auto"
      role="dialog"
      aria-modal="true"
      aria-label="确认入金"
    >
      <div class="fixed inset-0 bg-[#3b2f22]/50" @click="close" />

      <div
        class="relative flex min-h-full items-center justify-center px-4 py-6"
        @click.self="close"
      >
        <div
          class="relative w-full max-w-[320px] rounded-[24px] border border-[#e7d2a8] bg-[#fffdf8] px-4 pb-5 pt-[92px] shadow-[0_24px_60px_-20px_rgba(80,50,12,0.45)]"
        >
          <img
            :src="buyBg"
            alt=""
            class="pointer-events-none absolute -top-5 left-1/2 w-[120px] -translate-x-1/2 object-contain"
          />

          <div class="flex items-center gap-2">
            <span class="h-px flex-1 bg-gradient-to-r from-transparent to-[#d9ba7d]" />
            <span class="h-1.5 w-1.5 shrink-0 rotate-45 rounded-[1px] bg-[#c99a4a]" />
            <h2 class="shrink-0 text-[24px] font-black tracking-wide text-[#1f1a13]">
              {{ title }}
            </h2>
            <span class="h-1.5 w-1.5 shrink-0 rotate-45 rounded-[1px] bg-[#c99a4a]" />
            <span class="h-px flex-1 bg-gradient-to-l from-transparent to-[#d9ba7d]" />
          </div>

          <div class="mt-4 rounded-2xl bg-[#fbf5e8] px-3.5 py-3">
            <div class="space-y-2 text-[14px]">
              <div class="flex items-center justify-between gap-3">
                <span class="text-ink-500">入金挡位</span>
                <span class="font-medium text-ink-900">{{ tier }}</span>
              </div>
              <div class="flex items-center justify-between gap-3">
                <span class="text-ink-500">入金金额</span>
                <span class="font-medium tabular-nums text-ink-900">
                  {{ formatAmount(orderStatus?.amount ?? displayAmount) }} USDT
                </span>
              </div>
              <div class="flex items-center justify-between gap-3">
                <span class="text-ink-500">投资周期</span>
                <span class="font-medium text-ink-900">{{ period }}个月</span>
              </div>
              <div class="flex items-center justify-between gap-3">
                <span class="text-ink-500">预计月化</span>
                <span class="font-medium text-ink-900">{{ monthlyRate }}</span>
              </div>
              <div class="flex items-center justify-between gap-3 pt-0.5">
                <span class="text-ink-500">预计总收益</span>
                <span class="text-[17px] font-bold tabular-nums text-[#a45f1a]">
                  {{ formatAmount(totalReturn) }} USDT
                </span>
              </div>
            </div>
          </div>

          <div class="mt-4 flex items-center gap-2">
            <span class="h-px flex-1 bg-gradient-to-r from-transparent to-[#d9ba7d]" />
            <span class="h-1.5 w-1.5 shrink-0 rotate-45 rounded-[1px] bg-[#c99a4a]" />
            <h3 class="shrink-0 text-[16px] font-bold text-[#1f1a13]">支付方式</h3>
            <span class="h-1.5 w-1.5 shrink-0 rotate-45 rounded-[1px] bg-[#c99a4a]" />
            <span class="h-px flex-1 bg-gradient-to-l from-transparent to-[#d9ba7d]" />
          </div>

          <div
            class="mt-2.5 flex items-center gap-2.5 rounded-2xl border border-[#e7d2aa] bg-[#fffdf8] px-3 py-2.5"
          >
            <img :src="icon6" alt="" class="h-9 w-9 shrink-0 object-contain" />
            <span class="min-w-0 flex-1">
              <span class="block text-[15px] font-semibold text-[#1f1a13]">USDT 余额</span>
              <span class="mt-0.5 block text-[11.5px] text-ink-500">
                {{ selectedChain?.chainName || '未选择网络' }}
              </span>
            </span>
            <span class="shrink-0 text-right text-[15px] font-bold tabular-nums text-[#a45f1a]">
              {{ balanceText }}
              <span class="text-[10.5px] font-medium text-ink-500">USDT</span>
            </span>
          </div>

          <div v-if="chains.length > 1" class="mt-2">
            <p class="text-[12px] text-ink-500">选择支付网络</p>
            <div class="mt-1.5 grid grid-cols-2 gap-2">
              <button
                v-for="chain in chains"
                :key="chain.chain"
                type="button"
                class="h-8 rounded-full border text-[12px] font-medium transition disabled:opacity-60"
                :class="
                  selectedChain?.chain === chain.chain
                    ? 'border-[#cf9a4a] bg-gradient-to-r from-[#e1b26a] to-[#c89045] text-white'
                    : 'border-[#e5d5b7] bg-[#fffdf9] text-ink-700'
                "
                :disabled="busy"
                @click="selectChain(chain)"
              >
                {{ chain.chainName }}
              </button>
            </div>
          </div>

          <div v-else-if="selectedChain" class="mt-2 text-[12px] text-ink-500">
            当前网络：{{ selectedChain.chainName }}
          </div>

          <div
            v-if="phase === 'creating'"
            class="mt-3 rounded-xl bg-[#fbf4e6]/80 px-3 py-2 text-[12px] text-ink-600"
          >
            正在创建报单，请稍候...
          </div>

          <div
            v-else-if="phase === 'paying'"
            class="mt-3 rounded-xl bg-[#fbf4e6]/80 px-3 py-2 text-[12px] text-ink-600"
          >
            请在钱包中确认支付，支付金额必须为
            <span class="font-semibold text-[#a45f1a]">{{ formatAmount(displayAmount) }} USDT</span>
          </div>

          <div
            v-else-if="phase === 'confirming'"
            class="mt-5 flex flex-col items-center justify-center gap-2 py-2 text-[12px] text-ink-500"
          >
            <span
              class="h-9 w-9 animate-spin rounded-full border-[3px] border-[#e3cda0] border-t-[#c9954a]"
            />
            <span>支付确认中，请稍候...</span>
          </div>

          <div
            v-else-if="phase === 'failed' && errorMessage"
            class="mt-3 rounded-xl border border-[#f0cccc] bg-[#fff5f5] px-3 py-2 text-[12px] text-[#b33a3a]"
          >
            {{ errorMessage }}
          </div>

          <div
            v-else
            class="mt-3 flex items-start justify-center gap-1.5 text-center text-[11.5px] leading-4 text-[#2f7bd6]"
          >
            <svg
              viewBox="0 0 24 24"
              fill="currentColor"
              class="mt-0.5 h-3.5 w-3.5 shrink-0"
              aria-hidden="true"
            >
              <path
                d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm0 5.2a1.2 1.2 0 1 1 0 2.4 1.2 1.2 0 0 1 0-2.4ZM13.3 17h-2.6v-1.2h.7v-3.2h-.7v-1.2h2v4.4h.6V17Z"
              />
            </svg>
            <span>点击确认后将拉起您的 Web3 钱包进行签名支付</span>
          </div>

          <div v-if="phase !== 'success'" class="mt-4 grid grid-cols-2 gap-3">
            <button
              ref="cancelButton"
              type="button"
              class="h-11 rounded-full border border-[#e7d2aa] bg-[#fdf8ed] text-[15px] font-semibold text-[#5a4a35] outline-none transition focus-visible:ring-2 focus-visible:ring-[#c99a4a]/60 focus-visible:ring-offset-2 focus-visible:ring-offset-[#fffdf8] active:scale-[0.98]"
              @click="close"
            >
              {{ phase === 'confirming' ? '关闭' : '取消' }}
            </button>
            <button
              type="button"
              class="h-11 rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[15px] font-semibold text-white shadow-gold-btn outline-none transition focus-visible:ring-2 focus-visible:ring-[#c99a4a]/60 focus-visible:ring-offset-2 focus-visible:ring-offset-[#fffdf8] active:scale-[0.98] disabled:opacity-60"
              :disabled="busy"
              @click="confirmBuy"
            >
              {{ buttonText }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

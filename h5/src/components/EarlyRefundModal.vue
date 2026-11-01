<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { formatAmount } from '@/utils/format'

import endBg from '@/assets/images/end_bg.png'

const props = defineProps<{
  show: boolean
  amount: number
  released: number
  penaltyRate: number
  submitting: boolean
}>()

const emit = defineEmits<{
  'update:show': [boolean]
  confirm: []
}>()

const closeButton = ref<HTMLButtonElement | null>(null)

const penalty = computed(() => props.amount * (props.penaltyRate / 100))
const actual = computed(() => Math.max(props.amount - penalty.value - props.released, 0))


function close() {
  if (props.submitting) {
    return
  }
  emit('update:show', false)
}

function confirmRefund() {
  if (props.submitting) {
    return
  }
  emit('confirm')
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    close()
  }
}

watch(
  () => props.show,
  async (show) => {
    document.body.classList.toggle('h5-refund-open', show)

    if (show) {
      window.addEventListener('keydown', onKeydown)
      await nextTick()
      closeButton.value?.focus()
      return
    }

    window.removeEventListener('keydown', onKeydown)
  },
)

onBeforeUnmount(() => {
  document.body.classList.remove('h5-refund-open')
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
      aria-label="确认提前退本"
    >
      <div class="fixed inset-0 bg-[#3b2f22]/55" @click="close" />

      <div
        class="relative flex min-h-full items-center justify-center px-4 py-8"
        @click.self="close"
      >
        <div
          class="relative w-full max-w-[320px] rounded-[24px] bg-[#fffdf8] px-4 pb-5 pt-[54px] shadow-[0_24px_60px_-20px_rgba(80,50,12,0.45)]"
        >
          <img
            :src="endBg"
            alt=""
            class="pointer-events-none absolute -top-4 left-1/2 w-[190px] -translate-x-1/2 object-contain"
          />

          <button
            ref="closeButton"
            type="button"
            class="absolute right-3 top-3 flex h-9 w-9 items-center justify-center rounded-full text-[#8a6a3f] outline-none transition hover:bg-black/5 focus-visible:ring-2 focus-visible:ring-[#c99a4a]/60"
            aria-label="关闭"
            @click="close"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.8"
              stroke-linecap="round"
              class="h-6 w-6"
              aria-hidden="true"
            >
              <path d="m6 6 12 12M18 6 6 18" />
            </svg>
          </button>

          <h2 class="text-center text-[26px] font-black tracking-tight text-[#1f1a13]">
            确认提前退本？
          </h2>
          <p class="mt-2 text-center text-[13.5px] leading-5 text-ink-500">
            提前退本将扣除 {{ penaltyRate }}% 违约金，请确认是否继续
          </p>

          <div class="mt-4 rounded-2xl bg-[#fbf5e8] px-3.5 py-3 text-[14px]">
            <div class="flex items-center justify-between gap-3 py-1">
              <span class="text-ink-500">入金金额</span>
              <span class="font-medium tabular-nums text-ink-900">
                {{ formatAmount(amount) }} USDT
              </span>
            </div>
            <div class="flex items-center justify-between gap-3 py-1">
              <span class="text-ink-500">违约金（{{ penaltyRate }}%）</span>
              <span class="font-medium tabular-nums text-[#d14a3a]">
                -{{ formatAmount(penalty) }} USDT
              </span>
            </div>
            <div class="flex items-center justify-between gap-3 py-1">
              <span class="text-ink-500">累计已释放</span>
              <span class="font-medium tabular-nums text-ink-900">
                -{{ formatAmount(released) }} USDT
              </span>
            </div>

            <div
              class="mt-2 flex items-end justify-between gap-3 border-t border-[#eadcc0] pt-3"
            >
              <span class="text-[17px] font-bold text-[#1f1a13]">实际到账</span>
              <span
                class="whitespace-nowrap text-[27px] font-black leading-none tracking-tight tabular-nums text-[#a45f1a] max-[359px]:text-[24px]"
              >
                {{ formatAmount(actual) }}
                <span class="text-[12px] font-bold">USDT</span>
              </span>
            </div>
          </div>

          <p class="mt-2 text-center text-[11.5px] leading-4 text-[#8a7456]">
            实际到账 = 入金金额 × {{ 100 - penaltyRate }}% - 累计已释放收益
          </p>

          <div class="mt-4 grid grid-cols-2 gap-3">
            <button
              type="button"
              class="h-11 rounded-full border border-[#e7d2aa] bg-[#fdf8ed] text-[15px] font-semibold text-[#5a4a35] outline-none transition focus-visible:ring-2 focus-visible:ring-[#c99a4a]/60 focus-visible:ring-offset-2 focus-visible:ring-offset-[#fffdf8] active:scale-[0.98] disabled:opacity-60"
              :disabled="submitting"
              @click="close"
            >
              取消
            </button>
            <button
              type="button"
              class="flex h-11 items-center justify-center rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[15px] font-semibold text-white shadow-gold-btn outline-none transition focus-visible:ring-2 focus-visible:ring-[#c99a4a]/60 focus-visible:ring-offset-2 focus-visible:ring-offset-[#fffdf8] active:scale-[0.98] disabled:opacity-60"
              :disabled="submitting"
              @click="confirmRefund"
            >
              <span v-if="submitting" class="flex items-center gap-2">
                <span
                  class="h-4 w-4 animate-spin rounded-full border-2 border-white/40 border-t-white"
                />
                提交中...
              </span>
              <span v-else>确认退本</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

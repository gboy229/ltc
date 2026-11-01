<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

import icon1 from '@/assets/images/icon1.png'
import type { CmsArticleSummary } from '@/api/cms'

const props = defineProps<{
  notices: CmsArticleSummary[]
}>()

const emit = defineEmits<{
  select: [notice: CmsArticleSummary]
}>()

const activeIndex = ref(0)
const current = computed(() => props.notices[activeIndex.value] ?? null)
const hasMultiple = computed(() => props.notices.length > 1)

let timer: number | undefined

function stopTimer() {
  if (timer !== undefined) {
    window.clearInterval(timer)
    timer = undefined
  }
}

function startTimer() {
  stopTimer()
  if (!hasMultiple.value) {
    return
  }

  timer = window.setInterval(() => {
    if (props.notices.length < 2) {
      return
    }
    activeIndex.value = (activeIndex.value + 1) % props.notices.length
  }, 3000)
}

function selectCurrent() {
  if (current.value) {
    emit('select', current.value)
  }
}

function formatMonthDay(value?: string) {
  if (!value) {
    return ''
  }
  const match = value.match(/\d{4}-(\d{2})-(\d{2})/)
  return match ? `${match[1]}/${match[2]}` : ''
}

watch(
  () => props.notices,
  () => {
    activeIndex.value = 0
    startTimer()
  },
  { deep: true },
)

onMounted(startTimer)
onBeforeUnmount(stopTimer)
</script>

<template>
  <div
    class="relative mt-2 h-8 w-full overflow-hidden rounded-xl border border-[#e7d1a6] bg-[#fffdf6]/90 shadow-[0_2px_8px_rgba(166,120,44,0.06)]"
    aria-live="polite"
    aria-atomic="true"
  >
    <Transition
      mode="out-in"
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="translate-y-full opacity-0"
      enter-to-class="translate-y-0 opacity-100"
      leave-active-class="transition duration-300 ease-in"
      leave-from-class="translate-y-0 opacity-100"
      leave-to-class="-translate-y-full opacity-0"
    >
      <button
        v-if="current"
        :key="current.id"
        type="button"
        class="flex h-8 w-full items-center gap-2 px-3 text-left transition active:bg-[#fbf3e1]"
        @click="selectCurrent"
        @mouseenter="stopTimer"
        @mouseleave="startTimer"
        @touchstart.passive="stopTimer"
        @touchend.passive="startTimer"
      >
        <img :src="icon1" alt="" class="h-4 w-5 shrink-0 object-contain" />
        <span class="shrink-0 text-[14px] font-bold text-ink-900">公告</span>
        <span class="h-3 w-px shrink-0 bg-[#d9c092]" />
        <span class="min-w-0 flex-1 truncate text-[12.5px] text-ink-700">
          {{ current.title }}
        </span>
        <span class="shrink-0 text-[12px] text-ink-500">
          {{ formatMonthDay(current.createdAt) }}
        </span>
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.8"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="h-3 w-3 shrink-0 text-ink-500"
          aria-hidden="true"
        >
          <path d="m9 6 6 6-6 6" />
        </svg>
      </button>

      <div
        v-else
        key="__empty__"
        class="flex h-8 w-full items-center gap-2 px-3 text-[12.5px] text-ink-500"
      >
        <img :src="icon1" alt="" class="h-4 w-5 shrink-0 object-contain" />
        <span class="shrink-0 text-[14px] font-bold text-ink-900">公告</span>
        <span class="h-3 w-px shrink-0 bg-[#d9c092]" />
        <span class="min-w-0 flex-1 truncate">暂无公告</span>
      </div>
    </Transition>
  </div>
</template>

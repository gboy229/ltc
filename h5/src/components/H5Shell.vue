<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'

import SideNav from './SideNav.vue'
import { useSideNavStore } from '@/stores/sideNav'

const DESIGN_WIDTH = 375
const PC_CANVAS_WIDTH = 400
const DESKTOP_SCALE = PC_CANVAS_WIDTH / DESIGN_WIDTH
const MIN_FRAME_HEIGHT = 877
const DESKTOP_QUERY = '(min-width: 768px)'

const shell = ref<HTMLElement | null>(null)
const { isOpen: sideNavOpen } = storeToRefs(useSideNavStore())
const isDesktop = ref(
  typeof window !== 'undefined' && window.matchMedia(DESKTOP_QUERY).matches,
)
const frameHeight = ref(MIN_FRAME_HEIGHT)

const frameStyle = computed(() =>
  isDesktop.value
    ? {
        width: `${DESIGN_WIDTH * DESKTOP_SCALE}px`,
        height: `${frameHeight.value * DESKTOP_SCALE}px`,
      }
    : undefined,
)

const shellStyle = computed(() =>
  isDesktop.value ? { transform: `scale(${DESKTOP_SCALE})` } : undefined,
)

let observer: ResizeObserver | undefined
let raf = 0

function updateLayout() {
  const element = shell.value
  if (!element) {
    return
  }

  const desktop = window.matchMedia(DESKTOP_QUERY).matches
  isDesktop.value = desktop

  if (!desktop) {
    frameHeight.value = MIN_FRAME_HEIGHT
    return
  }

  const content = element.firstElementChild as HTMLElement | null
  frameHeight.value = Math.max(content?.offsetHeight ?? 0, MIN_FRAME_HEIGHT)
}

function scheduleUpdate() {
  cancelAnimationFrame(raf)
  raf = requestAnimationFrame(() => {
    void nextTick(updateLayout)
  })
}

onMounted(() => {
  scheduleUpdate()
  window.addEventListener('resize', scheduleUpdate)

  const content = shell.value?.firstElementChild
  if (content && 'ResizeObserver' in window) {
    observer = new ResizeObserver(scheduleUpdate)
    observer.observe(content)
  }
})

onBeforeUnmount(() => {
  cancelAnimationFrame(raf)
  window.removeEventListener('resize', scheduleUpdate)
  observer?.disconnect()
})
</script>

<template>
  <div class="h5-viewport">
    <div class="h5-frame" :style="frameStyle">
      <div ref="shell" class="h5-shell" :style="shellStyle">
        <div class="relative z-10 flex flex-1 flex-col" :inert="sideNavOpen ? true : undefined">
          <slot />
        </div>

        <SideNav />
      </div>
    </div>

    <div
      id="h5-toast-root"
      class="pointer-events-none fixed inset-x-0 top-0 z-[9999] flex flex-col items-center gap-2 px-4 pt-[calc(env(safe-area-inset-top)+1rem)] md:pt-7"
    />
  </div>
</template>

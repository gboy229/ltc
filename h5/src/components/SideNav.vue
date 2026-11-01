<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoute, useRouter } from 'vue-router'

import centerBg from '@/assets/images/side/center_bg.png'
import levelBg from '@/assets/images/side/level_bg.png'
import sideBg from '@/assets/images/side/bg.png'
import sideLogo from '@/assets/images/side/logo.png'
import nav1 from '@/assets/images/side/nav1.png'
import nav2 from '@/assets/images/side/nav2.png'
import nav3 from '@/assets/images/side/nav3.png'
import nav4 from '@/assets/images/side/nav4.png'
import nav5 from '@/assets/images/side/nav5.png'
import icon4 from '@/assets/images/icon4.png'
import { useSideNavStore } from '@/stores/sideNav'
import { useUserStore } from '@/stores/user'
import { shortenAddress } from '@/utils/address'
import { copyText } from '@/utils/clipboard'
import { toast } from '@/utils/toast'

const user = useUserStore()
const { address, uid: userUid, levelName, isLoggedIn } = storeToRefs(user)

const sideNav = useSideNavStore()
const { isOpen } = storeToRefs(sideNav)
const route = useRoute()
const router = useRouter()

const closeButton = ref<HTMLButtonElement | null>(null)
const uid = computed(() => userUid.value || '--')
const level = computed(() => levelName.value || '--')
const shortAddress = computed(() => (address.value ? shortenAddress(address.value) : '未连接'))

const SWIPE_THRESHOLD = 60
let touchStartX = 0
let touchStartY = 0
let touchStartTime = 0
let trackingTouch = false

const menuItems = [
  { key: 'home', label: '首页', icon: sideLogo },
  { key: 'orders', label: '报单明细', icon: nav1 },
  { key: 'community', label: '我的社区', icon: nav2 },
  { key: 'assets', label: '我的资产', icon: nav3 },
  { key: 'notice', label: '平台公告', icon: nav4 },
  { key: 'earnings', label: '收益明细', icon: nav5 },
] as const

type MenuItem = (typeof menuItems)[number]

async function copyAddress() {
  if (!address.value) {
    toast('请先连接钱包', 'info')
    return
  }
  const copied = await copyText(address.value)
  toast(copied ? '钱包地址已复制' : '复制失败，请手动复制', copied ? 'success' : 'error')
}

async function handleSwitchWallet() {
  sideNav.closeNav()
  await user.switchWallet()
}

async function handleLogout() {
  sideNav.closeNav()
  await user.logout()
}

async function handleItem(item: MenuItem) {
  sideNav.closeNav()

  if (item.key === 'home') {
    if (route.path !== '/') {
      await router.push('/')
    }
    return
  }

  if (item.key === 'orders') {
    if (route.path !== '/orders') {
      await router.push('/orders')
    }
    return
  }

  if (item.key === 'notice') {
    if (route.path !== '/notice') {
      await router.push('/notice')
    }
    return
  }

  if (item.key === 'earnings') {
    if (route.path !== '/earnings') {
      await router.push('/earnings')
    }
    return
  }

  if (item.key === 'community') {
    if (route.path !== '/community') {
      await router.push('/community')
    }
    return
  }

  if (item.key === 'assets') {
    if (route.path !== '/my-assets') {
      await router.push('/my-assets')
    }
    return
  }
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    sideNav.closeNav()
  }
}

function onTouchStart(event: TouchEvent) {
  if (event.touches.length !== 1) {
    return
  }

  const touch = event.touches[0]
  touchStartX = touch.clientX
  touchStartY = touch.clientY
  touchStartTime = Date.now()
  trackingTouch = true
}

function onTouchEnd(event: TouchEvent) {
  if (!trackingTouch) {
    return
  }
  trackingTouch = false

  const touch = event.changedTouches[0]
  if (!touch || Date.now() - touchStartTime > 800) {
    return
  }

  const deltaX = touch.clientX - touchStartX
  const deltaY = touch.clientY - touchStartY

  if (Math.abs(deltaX) < SWIPE_THRESHOLD || Math.abs(deltaX) < Math.abs(deltaY) * 1.2) {
    return
  }

  if (deltaX < 0 && !isOpen.value) {
    sideNav.openNav()
    return
  }

  if (deltaX > 0 && isOpen.value) {
    sideNav.closeNav()
  }
}

function onTouchCancel() {
  trackingTouch = false
}

watch(isOpen, async (open) => {
  document.body.classList.toggle('h5-nav-open', open)

  if (open) {
    if (window.matchMedia('(min-width: 768px)').matches) {
      window.scrollTo({ top: 0 })
    }
    window.addEventListener('keydown', onKeydown)
    await nextTick()
    closeButton.value?.focus()
    return
  }

  window.removeEventListener('keydown', onKeydown)
})

watch(
  () => route.fullPath,
  () => sideNav.closeNav(),
)

onMounted(() => {
  window.addEventListener('touchstart', onTouchStart, { passive: true })
  window.addEventListener('touchend', onTouchEnd, { passive: true })
  window.addEventListener('touchcancel', onTouchCancel, { passive: true })
})

onBeforeUnmount(() => {
  document.body.classList.remove('h5-nav-open')
  window.removeEventListener('keydown', onKeydown)
  window.removeEventListener('touchstart', onTouchStart)
  window.removeEventListener('touchend', onTouchEnd)
  window.removeEventListener('touchcancel', onTouchCancel)
})
</script>

<template>
  <div
    class="fixed inset-0 z-[100]"
    :class="isOpen ? 'pointer-events-auto' : 'pointer-events-none'"
    :aria-hidden="!isOpen"
    :inert="!isOpen ? true : undefined"
  >
    <div
      class="absolute inset-0 bg-[#3b2f22]/45 transition-opacity duration-300"
      :class="isOpen ? 'opacity-100' : 'opacity-0'"
      @click="sideNav.closeNav()"
    />

    <aside
      class="absolute right-0 top-0 flex h-full w-[76.5%] flex-col overflow-hidden border-l border-[#e6d2aa] bg-[#fdf9ef] transition-transform duration-300 ease-out"
      :class="isOpen ? 'translate-x-0' : 'translate-x-full'"
      id="side-nav-drawer"
      role="dialog"
      aria-modal="true"
      aria-label="侧边导航"
    >
      <img
        :src="sideBg"
        alt=""
        class="pointer-events-none absolute inset-0 h-full w-full object-cover object-top"
        aria-hidden="true"
      />

      <div
        class="relative z-10 flex h-full flex-col overflow-y-auto overscroll-contain px-3.5 pb-[calc(1rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))]"
      >
        <header class="flex shrink-0 items-center gap-1.5 max-[359px]:gap-1">
          <img
            :src="sideLogo"
            alt="LTC"
            class="h-11 w-11 shrink-0 object-contain max-[359px]:h-9 max-[359px]:w-9"
          />
          <span
            class="shrink-0 text-[22px] font-black leading-none tracking-[-0.03em] text-[#17130f] max-[359px]:text-[20px]"
          >
            LTC
          </span>

          <button
            type="button"
            class="ml-auto flex h-7 min-w-0 flex-1 max-w-[127px] items-center gap-1 rounded-full border border-[#d8c094] bg-[#fdf8ea]/80 px-2 text-ink-700 shadow-[0_1px_4px_rgba(160,115,45,0.08)] backdrop-blur-sm transition active:bg-[#f8efdb] max-[359px]:max-w-[96px] max-[359px]:gap-0.5 max-[359px]:px-1.5"
            aria-label="复制钱包地址"
            @click="copyAddress"
          >
            <img
              :src="icon4"
              alt=""
              class="h-3.5 w-3.5 shrink-0 object-contain max-[359px]:h-3 max-[359px]:w-3"
            />
            <span
              class="min-w-0 truncate text-[11px] font-medium tracking-tight max-[359px]:text-[10px]"
            >
              {{ shortAddress }}
            </span>
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.6"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-3 w-3 shrink-0 text-[#6b6358] max-[359px]:h-2.5 max-[359px]:w-2.5"
              aria-hidden="true"
            >
              <rect x="9" y="9" width="11" height="11" rx="2.5" />
              <path d="M5 15V5.5A2.5 2.5 0 0 1 7.5 3H17" />
            </svg>
          </button>

          <button
            ref="closeButton"
            type="button"
            class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-[#7d6a51] outline-none transition focus-visible:ring-2 focus-visible:ring-[#c9a15c]/60 focus-visible:ring-offset-2 focus-visible:ring-offset-[#fdf9ef] active:bg-black/5 max-[359px]:h-6 max-[359px]:w-6"
            aria-label="关闭侧边导航"
            @click="sideNav.closeNav()"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.7"
              stroke-linecap="round"
              class="h-[18px] w-[18px] max-[359px]:h-4 max-[359px]:w-4"
              aria-hidden="true"
            >
              <path d="M6 6l12 12M18 6 6 18" />
            </svg>
          </button>
        </header>

        <section
          class="relative mt-3 aspect-[527/195] shrink-0 overflow-hidden rounded-2xl border border-[#e8d1a3]"
        >
          <img
            :src="centerBg"
            alt=""
            class="pointer-events-none absolute inset-0 h-full w-full object-cover"
            aria-hidden="true"
          />
          <div
            class="relative flex h-full items-center gap-3 px-4 max-[359px]:gap-2 max-[359px]:px-3"
          >
            <img
              :src="sideLogo"
              alt=""
              class="h-16 w-16 shrink-0 object-contain max-[359px]:h-14 max-[359px]:w-14"
            />
            <div class="min-w-0">
              <p
                class="text-[18px] font-bold leading-tight tracking-tight text-ink-900 max-[359px]:text-[16px]"
              >
                UID:{{ uid }}
              </p>
              <div class="relative mt-1.5 w-[76px] max-[359px]:w-[70px]">
                <img :src="levelBg" alt="" class="block h-auto w-full" />
                <span
                  class="absolute left-[60%] top-1/2 -translate-x-1/2 -translate-y-1/2 text-[12px] font-bold text-white drop-shadow-[0_1px_1px_rgba(90,58,18,0.55)] max-[359px]:text-[11px]"
                >
                  {{ level }}
                </span>
              </div>
            </div>
          </div>
        </section>

        <nav
          class="mt-3 divide-y divide-[#f1e4cb] overflow-hidden rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 backdrop-blur-[2px]"
        >
          <button
            v-for="item in menuItems"
            :key="item.key"
            type="button"
            class="flex w-full items-center gap-3 px-3 py-3.5 text-left transition active:bg-[#fbf3e3] max-[359px]:gap-2 max-[359px]:px-2.5 max-[359px]:py-3"
            @click="handleItem(item)"
          >
            <img
              :src="item.icon"
              alt=""
              class="h-10 w-10 shrink-0 object-contain max-[359px]:h-9 max-[359px]:w-9"
            />
            <span
              class="min-w-0 flex-1 text-[15.5px] font-semibold text-ink-900 max-[359px]:text-[14px]"
            >
              {{ item.label }}
            </span>
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.8"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-4 w-4 shrink-0 text-[#c7a05f] max-[359px]:h-3.5 max-[359px]:w-3.5"
              aria-hidden="true"
            >
              <path d="m9 6 6 6-6 6" />
            </svg>
          </button>
        </nav>

        <div v-if="isLoggedIn" class="mt-3 grid shrink-0 grid-cols-2 gap-2">
          <button
            type="button"
            class="h-9 rounded-full border border-[#d7b57b] bg-[#fffdf8]/70 text-[13px] font-semibold text-[#8a5c22] transition active:scale-[0.98]"
            @click="handleSwitchWallet"
          >
            切换钱包
          </button>
          <button
            type="button"
            class="h-9 rounded-full border border-[#e2cda2] bg-[#fdf8ed] text-[13px] font-semibold text-[#6b5a42] transition active:scale-[0.98]"
            @click="handleLogout"
          >
            退出登录
          </button>
        </div>
      </div>
    </aside>
  </div>
</template>

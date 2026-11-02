<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { formatAmount } from '@/utils/format'

import icon5 from '@/assets/images/icon5.png'
import icon8 from '@/assets/images/icon8.png'
import icon9 from '@/assets/images/icon9.png'
import icon10 from '@/assets/images/icon10.png'
import myBg1 from '@/assets/images/my_bg1.png'
import myBg2 from '@/assets/images/my_bg2.png'
import levelBg from '@/assets/images/side/level_bg.png'
import AppHeader from '@/components/AppHeader.vue'
import PageHeading from '@/components/PageHeading.vue'
import { useUserStore } from '@/stores/user'
import { getTeamPush, type TeamMember } from '@/api/users/team'
import { getOrderStatistics } from '@/api/pledge'
import { copyText } from '@/utils/clipboard'
import { toast } from '@/utils/toast'

const router = useRouter()
const user = useUserStore()
const { uid, levelName } = storeToRefs(user)

const members = ref<TeamMember[]>([])
const directCount = ref(0)
const income = ref({
  direct: 0,
  static: 0,
})

async function loadMembers() {
  try {
    const res = await getTeamPush({ page: 1, pageSize: 5 })
    members.value = res.list ?? []
    directCount.value = res.totalCount ?? members.value.length
  } catch {
    members.value = []
    directCount.value = 0
  }
}

async function loadIncome() {
  try {
    const res = await getOrderStatistics()
    income.value = {
      direct: res.referralIncome ?? 0,
      static: res.staticIncome ?? 0,
    }
  } catch {
    // 保持 0
  }
}

onMounted(() => {
  void loadMembers()
  void loadIncome()
})

function goBack() {
  const state = window.history.state as { back?: string | null } | null
  if (state?.back) {
    router.back()
    return
  }
  void router.push('/')
}

// 邀请链接：打开站点首页并通过 ?uid= 记录邀请人，注册时自动绑定
const inviteLink = computed(() => {
  if (!uid.value) {
    return ''
  }
  const link = new URL(import.meta.env.BASE_URL || '/', window.location.origin)
  link.searchParams.set('uid', uid.value)
  return link.toString()
})

async function copyInviteLink() {
  if (!inviteLink.value) return
  const copied = await copyText(inviteLink.value)
  toast(copied ? '邀请链接已复制' : '复制失败，请手动复制', copied ? 'success' : 'error')
}

async function copyInviteCode() {
  if (!uid.value) return
  const copied = await copyText(uid.value)
  toast(copied ? '邀请码已复制' : '复制失败，请手动复制', copied ? 'success' : 'error')
}

function viewAllMembers() {
  if (directCount.value === 0) {
    toast('暂无直推成员', 'info')
    return
  }
  void router.push('/community/team')
}
</script>

<template>
  <div
    class="flex flex-1 flex-col px-4 pb-[calc(1.5rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))]"
  >
    <AppHeader />

    <PageHeading
      title="我的社区"
      subtitle="汇聚同好 · 携手同行 · 收益共享"
      @back="goBack"
    />

    <section
      class="relative mt-3 min-h-[112px] overflow-hidden rounded-2xl border border-[#e8d3a9] p-4"
    >
      <img
        :src="myBg1"
        alt=""
        class="pointer-events-none absolute inset-0 h-full w-full object-cover object-right"
      />
      <div class="relative z-10 flex items-center gap-3">
        <img :src="icon9" alt="" class="h-[72px] w-[82px] shrink-0 object-contain" />
        <div class="min-w-0">
          <p class="text-[17px] font-bold text-[#2a1c10]">
            UID：<span class="tabular-nums">{{ uid }}</span>
          </p>
          <div class="relative -ml-2 mt-1.5 h-[37px] w-[108px]">
            <img :src="levelBg" alt="" class="absolute inset-0 h-full w-full object-fill" />
            <span
              class="absolute inset-y-0 left-[32%] right-[5%] z-10 flex items-center justify-center text-[17px] font-bold tracking-wide text-white"
            >
              {{ levelName }}
            </span>
          </div>
        </div>
      </div>
    </section>

    <div class="mt-3 grid grid-cols-2 gap-3">
      <section
        class="min-h-[94px] rounded-2xl border border-[#e8d3a9] bg-gradient-to-br from-[#fffdf8] to-[#f8ecd6] p-3"
      >
        <div class="flex items-center gap-2">
          <img :src="icon8" alt="" class="h-8 w-8 shrink-0 object-contain" />
          <p class="text-[14px] text-ink-500">直推收益</p>
        </div>
        <p
          class="mt-2 whitespace-nowrap text-[20px] font-bold tabular-nums text-[#a45f1a] max-[359px]:text-[18px]"
        >
          +{{ formatAmount(income.direct) }}
          <span class="text-[10.5px] font-medium text-ink-500">USDT</span>
        </p>
      </section>

      <section
        class="min-h-[94px] rounded-2xl border border-[#e8d3a9] bg-gradient-to-br from-[#fffdf8] to-[#f8ecd6] p-3"
      >
        <div class="flex items-center gap-2">
          <img :src="icon5" alt="" class="h-8 w-8 shrink-0 object-contain" />
          <p class="text-[14px] text-ink-500">静态收益</p>
        </div>
        <p
          class="mt-2 whitespace-nowrap text-[20px] font-bold tabular-nums text-[#a45f1a] max-[359px]:text-[18px]"
        >
          +{{ formatAmount(income.static) }}
          <span class="text-[10.5px] font-medium text-ink-500">USDT</span>
        </p>
      </section>
    </div>

    <section
      class="mt-3 overflow-hidden rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 shadow-gold-soft"
    >
      <div class="relative h-[100px] overflow-hidden">
        <img
          :src="myBg2"
          alt=""
          class="pointer-events-none absolute inset-0 h-full w-full object-cover object-left"
        />
        <div
          class="pointer-events-none absolute inset-0 bg-gradient-to-l from-[#fffdf8]/85 via-[#fffdf8]/30 to-transparent"
        />
        <div
          class="relative z-10 flex h-full flex-col items-end justify-center pr-4 text-right"
        >
          <p class="text-[24px] font-black tracking-tight text-[#1f1a13] max-[359px]:text-[22px]">
            邀请好友
          </p>
          <p class="text-[21px] font-black tracking-tight text-[#b27a2d] max-[359px]:text-[19px]">
            共享收益
          </p>
          <span class="mt-1 h-0.5 w-6 rounded-full bg-[#c99a4a]" />
        </div>
      </div>

      <div class="space-y-2 p-3">
        <div class="flex items-center gap-2 rounded-xl bg-[#fbf5e8] px-2.5 py-2 max-[359px]:gap-1.5 max-[359px]:px-2">
          <span
            class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-[#f3e3c4] text-[#b7863c] max-[359px]:h-6 max-[359px]:w-6"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.8"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-3.5 w-3.5 max-[359px]:h-3 max-[359px]:w-3"
              aria-hidden="true"
            >
              <path d="M10.5 13.5a4 4 0 0 0 5.7 0l2.3-2.3a4 4 0 0 0-5.7-5.7l-1.3 1.3" />
              <path d="M13.5 10.5a4 4 0 0 0-5.7 0l-2.3 2.3a4 4 0 0 0 5.7 5.7l1.3-1.3" />
            </svg>
          </span>
          <span class="shrink-0 text-[13px] text-ink-500 max-[359px]:text-[12px]">邀请链接：</span>
          <span class="min-w-0 flex-1 truncate text-[13.5px] font-medium text-ink-900 max-[359px]:text-[12.5px]">
            {{ inviteLink }}
          </span>
          <button
            type="button"
            class="flex h-7 shrink-0 items-center gap-1 rounded-full bg-gradient-to-r from-[#e7c07a] to-[#c9954a] px-2.5 text-[11.5px] font-medium text-white shadow-sm transition active:scale-[0.98] max-[359px]:h-6 max-[359px]:gap-0.5 max-[359px]:px-2 max-[359px]:text-[11px]"
            @click="copyInviteLink"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.8"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-3 w-3 max-[359px]:h-2.5 max-[359px]:w-2.5"
              aria-hidden="true"
            >
              <rect x="8" y="8" width="11" height="11" rx="2" />
              <path d="M16 8V6a2 2 0 0 0-2-2H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h2" />
            </svg>
            复制
          </button>
        </div>

        <div class="flex items-center gap-2 rounded-xl bg-[#fbf5e8] px-2.5 py-2 max-[359px]:gap-1.5 max-[359px]:px-2">
          <span
            class="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-[#f3e3c4] text-[#b7863c] max-[359px]:h-6 max-[359px]:w-6"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.8"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-3.5 w-3.5 max-[359px]:h-3 max-[359px]:w-3"
              aria-hidden="true"
            >
              <path d="M10.5 13.5a4 4 0 0 0 5.7 0l2.3-2.3a4 4 0 0 0-5.7-5.7l-1.3 1.3" />
              <path d="M13.5 10.5a4 4 0 0 0-5.7 0l-2.3 2.3a4 4 0 0 0 5.7 5.7l1.3-1.3" />
            </svg>
          </span>
          <span class="shrink-0 text-[13px] text-ink-500 max-[359px]:text-[12px]">邀请码：</span>
          <span class="min-w-0 flex-1 truncate text-[13.5px] font-medium text-ink-900 max-[359px]:text-[12.5px]">
            {{ uid }}
          </span>
          <button
            type="button"
            class="flex h-7 shrink-0 items-center gap-1 rounded-full bg-gradient-to-r from-[#e7c07a] to-[#c9954a] px-2.5 text-[11.5px] font-medium text-white shadow-sm transition active:scale-[0.98] max-[359px]:h-6 max-[359px]:gap-0.5 max-[359px]:px-2 max-[359px]:text-[11px]"
            @click="copyInviteCode"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.8"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="h-3 w-3 max-[359px]:h-2.5 max-[359px]:w-2.5"
              aria-hidden="true"
            >
              <rect x="8" y="8" width="11" height="11" rx="2" />
              <path d="M16 8V6a2 2 0 0 0-2-2H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h2" />
            </svg>
            复制
          </button>
        </div>
      </div>
    </section>

    <section
      class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 p-3.5 shadow-gold-soft"
    >
      <div class="flex items-center gap-2">
        <svg
          viewBox="0 0 24 24"
          fill="currentColor"
          class="h-5 w-5 text-[#c79a4a]"
          aria-hidden="true"
        >
          <path
            d="M9 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8Zm7 0a3 3 0 1 0 0-6 3 3 0 0 0 0 6ZM2 20c0-3.3 3.1-6 7-6s7 2.7 7 6H2Zm14.5-6c2.8.4 5.5 2.6 5.5 6h-4.2c0-2.2-.6-4.2-1.3-6Z"
          />
        </svg>
        <h2 class="text-[19px] font-bold text-ink-900">我的直推</h2>
        <button
          type="button"
          class="ml-auto flex items-center gap-0.5 text-[12.5px] font-medium text-[#a97227]"
          @click="viewAllMembers"
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

      <div
        class="mt-2 divide-y divide-[#f1e5cc] rounded-xl border border-[#f1e5cc] bg-[#fffdf8]/80 px-3"
      >
        <div
          v-for="member in members"
          :key="`${member.uid}-${member.levelName}`"
          class="flex items-center gap-2.5 py-2.5 max-[359px]:gap-2"
        >
          <img :src="icon10" alt="" class="h-9 w-9 shrink-0 object-contain max-[359px]:h-8 max-[359px]:w-8" />
          <div class="min-w-0 flex-1">
            <p class="text-[14px] font-medium text-ink-900">UID：{{ member.uid }}</p>
            <p class="whitespace-nowrap text-[11.5px] text-ink-500 max-[359px]:text-[11px]">注册时间：{{ member.registeredAt }}</p>
          </div>
          <div class="relative h-[26px] w-[76px] shrink-0 max-[359px]:h-[24px] max-[359px]:w-[70px]">
            <img :src="levelBg" alt="" class="absolute inset-0 h-full w-full object-fill" />
            <span
              class="absolute inset-y-0 left-[32%] right-[5%] z-10 flex items-center justify-center text-[12px] font-bold tracking-wide text-white max-[359px]:text-[11.5px]"
            >
              {{ member.levelName }}
            </span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

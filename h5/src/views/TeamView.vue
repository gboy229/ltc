<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import icon10 from '@/assets/images/icon10.png'
import levelBg from '@/assets/images/side/level_bg.png'
import AppHeader from '@/components/AppHeader.vue'
import PageHeading from '@/components/PageHeading.vue'
import { getTeamPush, type TeamMember } from '@/api/users/team'

const router = useRouter()

const members = ref<TeamMember[]>([])
const total = ref(0)
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = 15
const hasMore = ref(false)
const sentinel = ref<HTMLElement | null>(null)

let observer: IntersectionObserver | null = null
let requestSeq = 0

async function loadMembers(reset = false, targetPage = page.value) {
  if (!reset && (loading.value || loadingMore.value)) {
    return
  }

  const seq = ++requestSeq
  const currentPage = reset ? 1 : targetPage
  if (reset) {
    loading.value = true
    loadingMore.value = false
    hasMore.value = false
    members.value = []
  } else {
    loadingMore.value = true
  }

  try {
    const res = await getTeamPush({ page: currentPage, pageSize })
    if (seq !== requestSeq) return

    const list = res.list ?? []
    members.value = reset ? list : [...members.value, ...list]
    total.value = res.totalCount ?? members.value.length
    page.value = res.page ?? currentPage
    hasMore.value = (res.pageCount ?? 0) > page.value
  } catch {
    if (seq !== requestSeq) return
    if (reset) {
      members.value = []
      total.value = 0
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
    void loadMembers(false, page.value + 1)
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

function loadMore() {
  if (!hasMore.value || loadingMore.value) return
  void loadMembers(false, page.value + 1)
}

function goBack() {
  const state = window.history.state as { back?: string | null } | null
  if (state?.back) {
    router.back()
    return
  }
  void router.push('/community')
}

onMounted(() => {
  void loadMembers(true)
  window.addEventListener('scroll', maybeLoadMore, { passive: true })
  setupObserver()
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
      title="我的直推"
      subtitle="团队成员 · 关系清晰 · 实时查看"
      @back="goBack"
    />

    <section class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-4 py-3 shadow-gold-soft">
      <p class="text-[13px] text-ink-500">直推总人数</p>
      <p class="mt-1 text-[28px] font-black leading-none tabular-nums text-[#a45f1a]">
        {{ total }}
        <span class="text-[12px] font-medium text-ink-500">人</span>
      </p>
    </section>

    <section class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 p-3 shadow-gold-soft">
      <div v-if="loading" class="py-10 text-center text-[13px] text-ink-500">加载中...</div>

      <div
        v-else-if="!members.length"
        class="rounded-xl border border-dashed border-[#e5d5b7] bg-[#fffdf8]/70 py-12 text-center text-[13px] text-ink-500"
      >
        暂无直推成员
      </div>

      <div v-else class="divide-y divide-[#f1e5cc] rounded-xl border border-[#f1e5cc] bg-[#fffdf8]/80 px-3">
        <div
          v-for="member in members"
          :key="member.uid"
          class="flex items-center gap-2.5 py-2.5 max-[359px]:gap-2"
        >
          <img :src="icon10" alt="" class="h-9 w-9 shrink-0 object-contain max-[359px]:h-8 max-[359px]:w-8" />
          <div class="min-w-0 flex-1">
            <p class="text-[14px] font-medium text-ink-900">UID：{{ member.uid }}</p>
            <p class="whitespace-nowrap text-[11.5px] text-ink-500 max-[359px]:text-[11px]">
              注册时间：{{ member.registeredAt || '--' }}
            </p>
          </div>
          <div class="relative h-[26px] w-[76px] shrink-0 max-[359px]:h-[24px] max-[359px]:w-[70px]">
            <img :src="levelBg" alt="" class="absolute inset-0 h-full w-full object-fill" />
            <span
              class="absolute inset-y-0 left-[32%] right-[5%] z-10 flex items-center justify-center text-[12px] font-bold tracking-wide text-white max-[359px]:text-[11.5px]"
            >
              {{ member.levelName || '--' }}
            </span>
          </div>
        </div>
      </div>

      <div ref="sentinel" class="h-px" aria-hidden="true" />

      <button
        v-if="hasMore"
        type="button"
        class="mt-2 flex h-9 w-full items-center justify-center rounded-full border border-[#d7b57b] bg-[#fffdf8]/60 text-[13.5px] font-medium text-[#a97227] transition active:bg-[#fbf2df]"
        :disabled="loadingMore"
        @click="loadMore"
      >
        {{ loadingMore ? '加载中...' : '查看更多' }}
      </button>
      <p v-else-if="members.length" class="mt-2 text-center text-[12px] text-ink-400">
        没有更多了
      </p>
    </section>
  </div>
</template>

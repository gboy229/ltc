<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

import ggBg from '@/assets/images/gg_bg1.png'
import AppHeader from '@/components/AppHeader.vue'
import NoticeCategoryIcon from '@/components/NoticeCategoryIcon.vue'
import PageHeading from '@/components/PageHeading.vue'
import { getArticleList, type CmsArticleSummary } from '@/api/cms'

const router = useRouter()

const notices = ref<CmsArticleSummary[]>([])
const loading = ref(true)
const loadingMore = ref(false)
const page = ref(1)
const pageSize = 15
const hasMore = ref(false)
const sentinel = ref<HTMLElement | null>(null)

let observer: IntersectionObserver | null = null
let requestSeq = 0

function goBack() {
  const state = window.history.state as { back?: string | null } | null
  if (state?.back) {
    router.back()
    return
  }
  void router.push('/')
}

function openNotice(id: string) {
  void router.push(`/notice/${id}`)
}

function formatDate(value: string) {
  return value ? value.slice(0, 10) : ''
}

async function loadNotices(reset = false, targetPage = page.value) {
  if (!reset && (loading.value || loadingMore.value)) {
    return
  }

  const seq = ++requestSeq
  const currentPage = reset ? 1 : targetPage
  if (reset) {
    loading.value = true
    loadingMore.value = false
    hasMore.value = false
    notices.value = []
  } else {
    loadingMore.value = true
  }

  try {
    const res = await getArticleList({ page: currentPage, pageSize })
    if (seq !== requestSeq) return

    const list = res.list ?? []
    notices.value = reset ? list : [...notices.value, ...list]
    page.value = res.page ?? currentPage
    hasMore.value = (res.pageCount ?? 0) > page.value
  } catch {
    if (seq !== requestSeq) return
    if (reset) {
      notices.value = []
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
    void loadNotices(false, page.value + 1)
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
  void loadNotices(false, page.value + 1)
}

onMounted(() => {
  void loadNotices(true)
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
      title="平台公告"
      subtitle="平台动态 · 最新资讯 · 一览掌握"
      @back="goBack"
    />

    <section
      class="relative mt-3 aspect-[685/183] overflow-hidden rounded-2xl border border-[#e8d3a9]"
    >
      <img :src="ggBg" alt="" class="absolute inset-0 h-full w-full object-cover" />
      <div
        class="absolute inset-0 flex flex-col items-center justify-center px-10 text-center"
      >
        <h2
          class="text-[24px] font-black tracking-tight text-[#2a1c10] max-[359px]:text-[21px]"
        >
          平台最新公告
        </h2>
        <p class="mt-1 text-[12px] text-[#7c6138] max-[359px]:text-[11px]">
          及时了解 LTC 平台动态与重要通知
        </p>
      </div>
    </section>

    <div v-if="loading" class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-4 py-10 text-center text-[13px] text-ink-500">
      公告加载中...
    </div>

    <div v-else-if="!notices.length" class="mt-3 rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-4 py-10 text-center text-[13px] text-ink-500">
      暂无公告
    </div>

    <div v-else class="mt-3 space-y-2.5">
      <button
        v-for="notice in notices"
        :key="notice.id"
        type="button"
        class="relative w-full rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-3.5 py-2 text-left shadow-gold-soft transition active:scale-[0.995]"
        @click="openNotice(notice.id)"
      >
        <span
          class="inline-flex h-5 items-center gap-1 rounded-full bg-gradient-to-r from-[#e7c07a] to-[#c9954a] px-2.5 text-[11.5px] font-medium text-white shadow-sm max-[359px]:text-[11px]"
        >
          <NoticeCategoryIcon name="system" />
          公告
        </span>

        <div class="mt-1.5 flex items-center gap-2">
          <svg
            viewBox="0 0 24 24"
            fill="currentColor"
            class="h-3.5 w-3.5 shrink-0 text-[#c79a4a]"
            aria-hidden="true"
          >
            <path d="M12 2l2.2 7.8L22 12l-7.8 2.2L12 22l-2.2-7.8L2 12l7.8-2.2L12 2Z" />
          </svg>
          <span
            class="min-w-0 flex-1 text-[15.5px] font-semibold leading-[20px] text-[#1f1a13] max-[359px]:text-[14.5px]"
          >
            {{ notice.title }}
          </span>
          <span
            class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-[#f2d493] to-[#c9954a] text-white shadow-sm"
          >
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
          </span>
        </div>

        <p class="mt-1 pl-[22px] text-[12px] text-ink-500">{{ formatDate(notice.createdAt) }}</p>
      </button>
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
    <p v-else-if="notices.length" class="mt-2 text-center text-[12px] text-ink-400">
      没有更多了
    </p>
  </div>
</template>

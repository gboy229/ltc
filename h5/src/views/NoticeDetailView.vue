<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import ggBg from '@/assets/images/gg_bg2.png'
import AppHeader from '@/components/AppHeader.vue'
import PageHeading from '@/components/PageHeading.vue'
import { getArticleView, type CmsArticleDetail } from '@/api/cms'

const route = useRoute()
const router = useRouter()

const notice = ref<CmsArticleDetail | null>(null)
const loading = ref(true)
const error = ref('')

const noticeId = computed(() => String(route.params.id ?? ''))
const safeContent = computed(() => sanitizeHtml(notice.value?.content ?? ''))

function goBack() {
  const state = window.history.state as { back?: string | null } | null
  if (state?.back) {
    router.back()
    return
  }
  void router.push('/notice')
}

function formatDateTime(value?: string) {
  if (!value) {
    return ''
  }
  const match = value.match(/^(\d{4}-\d{2}-\d{2})[ T](\d{2}:\d{2})/)
  return match ? `${match[1]} ${match[2]}` : value
}

/** 仅保留安全标签/属性，避免公告正文注入脚本 */
function sanitizeHtml(html: string): string {
  if (!html) {
    return ''
  }

  const doc = new DOMParser().parseFromString(html, 'text/html')
  doc
    .querySelectorAll('script,style,iframe,object,embed,form,input,button,textarea,select,link,meta')
    .forEach((node) => node.remove())

  doc.querySelectorAll('*').forEach((el) => {
    for (const attr of Array.from(el.attributes)) {
      const name = attr.name.toLowerCase()
      const value = attr.value.trim().toLowerCase()
      if (
        name.startsWith('on') ||
        name === 'srcdoc' ||
        ((name === 'href' || name === 'src') && value.startsWith('javascript:'))
      ) {
        el.removeAttribute(attr.name)
      }
    }

    const tag = el.tagName.toLowerCase()
    if (tag === 'a') {
      el.setAttribute('target', '_blank')
      el.setAttribute('rel', 'noopener noreferrer')
    }
    if (tag === 'img') {
      el.setAttribute('loading', 'lazy')
    }
  })

  return doc.body.innerHTML
}

async function loadNotice() {
  const id = noticeId.value.trim()
  if (!id) {
    error.value = '公告不存在'
    loading.value = false
    return
  }

  loading.value = true
  error.value = ''

  try {
    const data = await getArticleView(id)
    if (!data?.id) {
      throw new Error('empty article')
    }
    notice.value = data
  } catch {
    notice.value = null
    error.value = '公告不存在或已下架'
  } finally {
    loading.value = false
  }
}

onMounted(loadNotice)
watch(noticeId, loadNotice)
</script>

<template>
  <div
    class="flex flex-1 flex-col px-4 pb-[calc(1.5rem+env(safe-area-inset-bottom))] pt-[calc(1rem+env(safe-area-inset-top))]"
  >
    <AppHeader />

    <PageHeading
      title="公告详情"
      subtitle="平台动态 · 最新资讯 · 一览掌握"
      @back="goBack"
    />

    <section
      class="relative mt-3 min-h-[150px] overflow-hidden rounded-2xl border border-[#e8d3a9] px-4 py-4"
    >
      <img
        :src="ggBg"
        alt=""
        class="pointer-events-none absolute inset-0 h-full w-full object-cover object-right"
      />
      <div class="relative z-10">
        <span
          class="inline-flex h-5 items-center rounded-full bg-gradient-to-r from-[#e7c07a] to-[#c9954a] px-2.5 text-[11.5px] font-medium text-white shadow-sm"
        >
          公告
        </span>

        <h2
          class="mt-3 text-[22px] font-black leading-[1.28] tracking-tight text-[#2a1c10] max-[359px]:text-[24px]"
        >
          {{ notice?.title || '公告详情' }}
        </h2>

        <div
          v-if="notice"
          class="mt-2 flex items-center gap-1.5 text-[13.5px] text-[#7c6a50]"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.7"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-4 w-4 shrink-0"
            aria-hidden="true"
          >
            <circle cx="12" cy="12" r="8.5" />
            <path d="M12 7v5l3 2" />
          </svg>
          <span>{{ formatDateTime(notice.createdAt) }}</span>
        </div>
      </div>
    </section>

    <article
      v-if="loading"
      class="mt-3 flex min-h-[318px] items-center justify-center rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-4 pb-6 pt-4 text-[13px] text-ink-500 shadow-gold-soft"
    >
      公告加载中...
    </article>

    <article
      v-else-if="!notice"
      class="mt-3 flex min-h-[318px] flex-col items-center justify-center rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-4 pb-6 pt-4 text-center shadow-gold-soft"
    >
      <p class="text-[14px] text-ink-600">{{ error || '公告不存在或已下架' }}</p>
      <button
        type="button"
        class="mt-4 h-9 rounded-full border border-[#d7b57b] bg-[#fffdf8]/60 px-5 text-[13px] font-medium text-[#a97227] transition active:bg-[#fbf2df]"
        @click="router.push('/notice')"
      >
        返回公告列表
      </button>
    </article>

    <article
      v-else
      class="mt-3 flex min-h-[318px] flex-col rounded-2xl border border-[#e8d3a9] bg-[#fffdf8]/95 px-4 pb-6 pt-4 shadow-gold-soft max-[359px]:min-h-[290px]"
    >
      <div
        class="notice-content text-[15px] leading-[27px] text-[#2a2118] max-[359px]:text-[14px] max-[359px]:leading-[25px]"
        v-html="safeContent"
      />
    </article>
  </div>
</template>

<style scoped>
.notice-content :deep(p) {
  margin: 0 0 14px;
}

.notice-content :deep(p:last-child) {
  margin-bottom: 0;
}

.notice-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 12px;
}

.notice-content :deep(a) {
  color: #a45f1a;
  text-decoration: underline;
}

.notice-content :deep(ul),
.notice-content :deep(ol) {
  margin: 0 0 14px;
  padding-left: 20px;
}

.notice-content :deep(ul) {
  list-style: disc;
}

.notice-content :deep(ol) {
  list-style: decimal;
}
</style>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useUserStore } from '@/stores/user'

const user = useUserStore()
const {
  registerVisible,
  registerStep,
  registerInviteUid,
  registerSubmitting,
  registerError,
} = storeToRefs(user)

async function confirm() {
  if (registerStep.value === 'bind') {
    await user.confirmBindInvite()
    return
  }
  await user.confirmRegister()
}
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
      v-if="registerVisible"
      class="fixed inset-0 z-[140] flex items-center justify-center bg-[#3b2f22]/50 px-5"
      role="dialog"
      aria-modal="true"
      aria-label="注册确认"
    >
      <section
        class="w-full max-w-[340px] rounded-2xl border border-[#e8d3a9] bg-[#fffdf8] p-5 shadow-[0_18px_50px_-20px_rgba(80,50,12,0.45)]"
      >
        <header class="relative">
          <h2 class="pr-8 text-[19px] font-black tracking-tight text-ink-900">
            {{ registerStep === 'bind' ? '绑定邀请码' : '确认注册' }}
          </h2>
          <button
            type="button"
            class="absolute right-0 top-0 flex h-7 w-7 items-center justify-center rounded-full text-[#7d6a51] transition active:bg-black/5"
            aria-label="关闭"
            :disabled="registerSubmitting"
            @click="user.cancelRegister()"
          >
            <svg
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="1.7"
              stroke-linecap="round"
              class="h-[18px] w-[18px]"
              aria-hidden="true"
            >
              <path d="M6 6l12 12M18 6 6 18" />
            </svg>
          </button>
        </header>

        <template v-if="registerStep === 'confirm'">
          <p class="mt-3 text-[13.5px] leading-6 text-ink-500">
            当前没有绑定邀请码，是否确认注册？
          </p>
          <p class="mt-1 text-[12px] leading-5 text-ink-400">
            绑定邀请码后可建立社区关系，不影响注册。
          </p>
        </template>

        <template v-else>
          <p class="mt-3 text-[13px] leading-5 text-ink-500">
            请输入邀请码后完成注册。
          </p>
          <label class="mt-4 block text-[12.5px] font-medium text-ink-600">邀请码</label>
          <input
            v-model="registerInviteUid"
            placeholder="请输入邀请码"
            class="mt-1.5 h-11 w-full rounded-xl border border-[#e7d2aa] bg-[#fdf8ed] px-3 text-[14px] text-ink-900 outline-none transition focus:border-[#c9954a]"
            :disabled="registerSubmitting"
          />
        </template>

        <p v-if="registerError" class="mt-3 text-[12.5px] leading-5 text-[#d14a3a]">
          {{ registerError }}
        </p>

        <div v-if="registerStep === 'confirm'" class="mt-5 grid grid-cols-2 gap-2">
          <button
            type="button"
            class="h-10 rounded-full border border-[#d7b57b] bg-[#fffdf8] text-[13.5px] font-semibold text-[#8a5c22] transition active:bg-[#fbf2df]"
            :disabled="registerSubmitting"
            @click="user.goBind()"
          >
            去绑定
          </button>
          <button
            type="button"
            class="h-10 rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[13.5px] font-semibold text-white shadow-gold-btn transition active:scale-[0.99] disabled:opacity-60"
            :disabled="registerSubmitting"
            @click="confirm"
          >
            {{ registerSubmitting ? '处理中...' : '确认注册' }}
          </button>
        </div>

        <div v-else class="mt-5 grid grid-cols-2 gap-2">
          <button
            type="button"
            class="h-10 rounded-full border border-[#d7b57b] bg-[#fffdf8] text-[13.5px] font-semibold text-[#8a5c22] transition active:bg-[#fbf2df]"
            :disabled="registerSubmitting"
            @click="user.cancelRegister()"
          >
            取消
          </button>
          <button
            type="button"
            class="h-10 rounded-full bg-gradient-to-r from-[#bd8842] via-[#d6aa68] to-[#eccd96] text-[13.5px] font-semibold text-white shadow-gold-btn transition active:scale-[0.99] disabled:opacity-60"
            :disabled="registerSubmitting"
            @click="confirm"
          >
            {{ registerSubmitting ? '处理中...' : '确认绑定并注册' }}
          </button>
        </div>
      </section>
    </div>
  </Transition>
</template>

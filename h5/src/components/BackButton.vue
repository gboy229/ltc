<script setup lang="ts">
defineEmits<{
  click: []
}>()
</script>

<template>
  <button
    type="button"
    class="relative flex h-9 items-center justify-center rounded-full text-[#2a1c10] transition after:absolute after:-inset-1.5 after:content-[''] active:bg-black/5"
    aria-label="返回"
    @click="$emit('click')"
  >
    <svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="1.8"
      stroke-linecap="round"
      stroke-linejoin="round"
      class="h-6 w-6"
      aria-hidden="true"
    >
      <path d="m15 6-6 6 6 6" />
    </svg>
  </button>
</template>

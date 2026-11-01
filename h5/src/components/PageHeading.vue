<script setup lang="ts">
import BackButton from './BackButton.vue'

defineProps<{
  title: string
  subtitle?: string
}>()

defineEmits<{
  back: []
}>()
</script>

<template>
  <div class="mt-3">
    <div class="flex items-center gap-1">
      <BackButton @click="$emit('back')" />
      <h1
        class="text-[24px] font-black leading-none tracking-tight text-[#1f1a13] max-[359px]:text-[22px]"
      >
        {{ title }}
      </h1>
    </div>
    <p v-if="subtitle" class="mt-1.5 text-[14px] text-ink-500 max-[359px]:text-[13px]">
      {{ subtitle }}
    </p>
  </div>
</template>

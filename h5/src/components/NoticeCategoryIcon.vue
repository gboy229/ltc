<script setup lang="ts">
type CategoryIconName = 'system' | 'activity' | 'maintenance' | 'community' | 'rules'

defineProps<{
  name: CategoryIconName
}>()
</script>

<template>
  <svg
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="1.8"
    stroke-linecap="round"
    stroke-linejoin="round"
    class="h-3 w-3 shrink-0"
    aria-hidden="true"
  >
    <template v-if="name === 'system'">
      <circle cx="12" cy="12" r="3" />
      <path
        d="M12 2v3M12 19v3M2 12h3M19 12h3M4.9 4.9 7 7M17 17l2.1 2.1M19.1 4.9 17 7M7 17l-2.1 2.1"
      />
    </template>
    <template v-else-if="name === 'activity'">
      <rect x="4" y="9" width="16" height="11" rx="2" />
      <path
        d="M4 13h16M12 9v11M8.5 9C7 9 6 8 6 6.8S7 4.5 8.2 4.5c2 0 3.8 4.5 3.8 4.5M15.5 9C17 9 18 8 18 6.8s-1-2.3-2.2-2.3c-2 0-3.8 4.5-3.8 4.5"
      />
    </template>
    <template v-else-if="name === 'maintenance'">
      <path
        d="M14.7 6.3a4 4 0 0 0-5.6 5.6L3 18l3 3 6.1-6.1a4 4 0 0 0 5.6-5.6l-2.8 2.8-2.8-2.8 2.6-2.8Z"
      />
    </template>
    <template v-else-if="name === 'community'">
      <circle cx="9" cy="8" r="3" />
      <path d="M3 20a6 6 0 0 1 12 0M16 5a3 3 0 0 1 0 6M18 20a6 6 0 0 0-3-5.2" />
    </template>
    <template v-else>
      <path d="M5 4h14v16H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2Z" />
      <path d="M8 8h7M8 12h5" />
    </template>
  </svg>
</template>

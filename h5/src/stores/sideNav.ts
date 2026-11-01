import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useSideNavStore = defineStore('sideNav', () => {
  const isOpen = ref(false)

  function openNav() {
    isOpen.value = true
  }

  function closeNav() {
    isOpen.value = false
  }

  function toggleNav() {
    isOpen.value = !isOpen.value
  }

  return { isOpen, openNav, closeNav, toggleNav }
})

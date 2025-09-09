<template>
  <transition name="fade">
    <div v-if="state.show" class="fixed bottom-6 left-1/2 -translate-x-1/2 z-50">
      <div :class="['px-4 py-2 rounded shadow text-white', colorCls]">{{ state.message }}</div>
    </div>
  </transition>
</template>

<script>
import { defineComponent, computed } from 'vue'
import { toastState } from '../utils/toast'

export default defineComponent({
  name: 'Toast',
  setup() {
    const state = toastState
    const colorCls = computed(() => {
      switch (state.kind) {
        case 'error': return 'bg-red-600'
        case 'success': return 'bg-green-600'
        default: return 'bg-gray-800'
      }
    })
    return { state, colorCls }
  },
})
</script>

<style>
.fade-enter-active,.fade-leave-active{ transition: opacity .2s }
.fade-enter-from,.fade-leave-to{ opacity: 0 }
</style>




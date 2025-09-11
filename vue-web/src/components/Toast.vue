<template>
  <teleport to="body">
    <transition name="toast-slide">
      <div v-if="state.show" class="fixed top-6 right-6 z-[9999] max-w-sm">
        <div :class="['px-4 py-3 rounded-xl backdrop-blur-sm border shadow-lg flex items-center gap-3', colorCls]">
          <!-- 图标 -->
          <div class="flex-shrink-0">
            <svg v-if="state.kind === 'error'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            <svg v-else-if="state.kind === 'success'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <svg v-else-if="state.kind === 'warning'" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
            </svg>
          </div>
          
          <!-- 消息文本 -->
          <div class="flex-1 text-sm font-medium leading-relaxed">{{ state.message }}</div>
          
          <!-- 关闭按钮 -->
          <button 
            @click="closeToast" 
            class="flex-shrink-0 ml-2 p-1 rounded-lg hover:bg-black/5 transition-colors duration-200"
            :class="closeButtonCls"
          >
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
      </div>
    </transition>
  </teleport>
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
        case 'error': 
          return 'bg-red-50/95 border-red-200/80 text-red-800'
        case 'success': 
          return 'bg-green-50/95 border-green-200/80 text-green-800'
        case 'warning': 
          return 'bg-amber-50/95 border-amber-200/80 text-amber-800'
        default: 
          return 'bg-gray-50/95 border-gray-200/80 text-gray-800'
      }
    })
    
    const closeButtonCls = computed(() => {
      switch (state.kind) {
        case 'error': return 'text-red-600 hover:text-red-800'
        case 'success': return 'text-green-600 hover:text-green-800'
        case 'warning': return 'text-amber-600 hover:text-amber-800'
        default: return 'text-gray-600 hover:text-gray-800'
      }
    })
    
    const closeToast = () => {
      state.show = false
    }
    
    return { 
      state, 
      colorCls, 
      closeButtonCls, 
      closeToast 
    }
  },
})
</script>

<style scoped>
.toast-slide-enter-active,
.toast-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.toast-slide-enter-from {
  opacity: 0;
  transform: translateX(100%) scale(0.95);
}

.toast-slide-leave-to {
  opacity: 0;
  transform: translateX(100%) scale(0.95);
}

.toast-slide-enter-to,
.toast-slide-leave-from {
  opacity: 1;
  transform: translateX(0) scale(1);
}
</style>




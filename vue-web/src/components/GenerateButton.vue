<template>
  <button
    class="px-4 py-2 text-white bg-indigo-600 rounded-md  hover:bg-indigo-500 disabled:opacity-50"
    :disabled="disabled || loading"
    @click="$emit('click')"
  >
    <span v-if="loading">生成中...</span>
    <span v-else>{{ text }}</span>
  </button>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'GenerateButton',
  props: {
    text: { type: String, default: '开始生成' },
    loading: { type: Boolean, default: false },
    disabled: { type: Boolean, default: false },
  },
  emits: ['click'],
})
</script>


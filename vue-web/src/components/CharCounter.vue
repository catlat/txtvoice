<template>
  <div class="text-sm text-gray-500">
    <span>字符数：{{ count }}</span>
    <span class="ml-2 text-gray-400">（以服务端计数为准）</span>
  </div>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'CharCounter',
  props: { count: { type: Number, required: true } },
})
</script>





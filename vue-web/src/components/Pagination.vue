<template>
  <nav class="flex items-center justify-between select-none" aria-label="Pagination">
    <!-- Prev -->
    <button
      class="px-3 py-2 text-sm font-medium rounded-md border transition-colors duration-200"
      :class="canPrev ? 'bg-white text-gray-800 border-gray-300 hover:bg-gray-50' : 'bg-white text-gray-400 border-gray-200 cursor-not-allowed'"
      :disabled="!canPrev"
      @click="$emit('update:page', currentPage - 1)"
      aria-label="上一页"
    >
      上一页
    </button>

    <!-- Page Numbers -->
    <div class="flex items-center gap-1">
      <template v-if="totalPages > 0">
        <button
          v-if="currentPage > 3"
          class="px-3 py-2 text-sm rounded-md border transition-colors duration-200"
          :class="btnCls(1)"
          @click="$emit('update:page', 1)"
        >1</button>
        <span v-if="currentPage > 4" class="px-2 text-gray-400">…</span>

        <button
          v-for="p in visiblePages"
          :key="p"
          class="px-3 py-2 text-sm rounded-md border transition-colors duration-200"
          :class="btnCls(p)"
          @click="$emit('update:page', p)"
        >{{ p }}</button>

        <span v-if="currentPage < totalPages - 3" class="px-2 text-gray-400">…</span>
        <button
          v-if="currentPage < totalPages - 2"
          class="px-3 py-2 text-sm rounded-md border transition-colors duration-200"
          :class="btnCls(totalPages)"
          @click="$emit('update:page', totalPages)"
        >{{ totalPages }}</button>
      </template>
    </div>

    <!-- Next -->
    <button
      class="px-3 py-2 text-sm font-medium rounded-md border transition-colors duration-200"
      :class="canNext ? 'bg-white text-gray-800 border-gray-300 hover:bg-gray-50' : 'bg-white text-gray-400 border-gray-200 cursor-not-allowed'"
      :disabled="!canNext"
      @click="$emit('update:page', currentPage + 1)"
      aria-label="下一页"
    >
      下一页
    </button>
  </nav>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Pagination',
  props: {
    page: { type: Number, default: 1 },
    pageSize: { type: Number, default: 10 },
    total: { type: Number, default: 0 },
    totalPages: { type: Number, default: 0 },
    hasPrev: { type: Boolean, default: false },
    hasNext: { type: Boolean, default: false },
    range: { type: Number, default: 2 },
  },
  emits: ['update:page'],
  computed: {
    currentPage() { return this.page || 1 },
    canPrev() { return this.hasPrev || this.currentPage > 1 },
    canNext() {
      if (this.totalPages > 0) return this.hasNext || this.currentPage < this.totalPages
      // 当后端不返回 totalPages 时，保守允许下一页由外部控制
      return this.hasNext
    },
    visiblePages() {
      const totalPages = this.totalPages
      const cur = this.currentPage
      const r = this.range
      if (totalPages <= 0) return [cur]
      const start = Math.max(1, cur - r)
      const end = Math.min(totalPages, cur + r)
      const list = []
      for (let i = start; i <= end; i++) list.push(i)
      return list
    },
  },
  methods: {
    btnCls(p) {
      const isActive = p === this.currentPage
      return isActive
        ? 'bg-gray-900 text-white border-gray-900'
        : 'bg-white text-gray-800 border-gray-300 hover:bg-gray-50'
    },
  },
})
</script>



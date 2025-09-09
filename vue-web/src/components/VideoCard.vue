<template>
  <div class="group relative flex items-start gap-3 md:gap-4 rounded-[12px] border border-gray-200 bg-white p-3 md:p-4 shadow-sm hover:shadow transition">
    <!-- 删除按钮（卡片右上角，贴边框外） -->
    <button
      class="absolute -top-2 -right-2 z-10 inline-flex items-center justify-center w-6 h-6 rounded-full border border-gray-300 bg-white text-gray-500 hover:text-gray-700 hover:bg-gray-50 shadow"
      title="删除"
      @click="$emit('delete')"
    >
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-3.5 h-3.5" fill="currentColor"><path d="M6 7a1 1 0 0 1 1-1h10a1 1 0 1 1 0 2H7a1 1 0 0 1-1-1m2 3a1 1 0 0 1 1 1v6a1 1 0 1 1-2 0v-6a1 1 0 0 1 1-1m4 0a1 1 0 0 1 1 1v6a1 1 0 1 1-2 0v-6a1 1 0 0 1 1-1m5 0a1 1 0 0 1 1 1v6a1 1 0 1 1-2 0v-6a1 1 0 0 1 1-1M9 4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2h3a1 1 0 1 1 0 2H6a1 1 0 1 1 0-2z"/></svg>
    </button>
    <div class="shrink-0">
      <a v-if="youtubeUrl" :href="youtubeUrl" target="_blank" rel="noopener" class="block">
        <img :src="thumbnail" alt="thumb" class="w-16 h-16 md:w-20 md:h-20 rounded-lg object-cover" />
      </a>
      <img v-else :src="thumbnail" alt="thumb" class="w-16 h-16 md:w-20 md:h-20 rounded-lg object-cover" />
    </div>
    <div class="flex-1 min-w-0">
      <h3 class="text-sm md:text-base font-medium text-gray-900 leading-snug line-clamp-2">
        <template v-if="resolvedUrl">
          <a :href="resolvedUrl" target="_blank" rel="noopener" class="hover:underline">
            {{ title }}
          </a>
        </template>
        <template v-else>
          {{ title }}
        </template>
      </h3>

      <div class="mt-2 flex gap-2 text-xs text-gray-500">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4 text-gray-400 shrink-0" fill="currentColor" aria-hidden="true"><path d="M12 2a5 5 0 1 1-5 5 5.006 5.006 0 0 1 5-5m0 9a9 9 0 0 0-9 9 1 1 0 0 0 1 1h16a1 1 0 0 0 1-1 9 9 0 0 0-9-9"/></svg>
        <span class="flex-1 whitespace-normal break-words leading-tight" :title="displayAuthor">{{ displayAuthor }}</span>
      </div>

      <div class="mt-3 flex items-center justify-between">
        <div class="flex items-center gap-4 text-xs text-gray-500">
          <div class="inline-flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4" fill="currentColor"><path d="M3 10a5 5 0 0 1 9.9-1h-2.1a3 3 0 1 0-5.8 1zm8.1 1h1.8a5 5 0 1 1 0 2h-1.8A3 3 0 1 1 5 13H3.2a5 5 0 1 0 9.9-2z"/></svg>
            <span>{{ displayViews }}</span>
          </div>
          <div class="inline-flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4" fill="currentColor"><path d="M12 2a10 10 0 1 0 10 10A10.011 10.011 0 0 0 12 2m1 11h-4a1 1 0 0 1 0-2h3V7a1 1 0 0 1 2 0z"/></svg>
            <span>{{ displayDuration }}</span>
          </div>
          <div v-if="displayDate" class="inline-flex items-center gap-1">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4" fill="currentColor"><path d="M6 2a1 1 0 0 1 1 1v1h10V3a1 1 0 0 1 2 0v1h1a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h1V3a1 1 0 1 1 2 0v1zm14 6H4v10h16z"/></svg>
            <span>{{ displayDate }}</span>
          </div>
        </div>
      </div>

      <div class="mt-2 text-[11px] text-gray-400 truncate" v-if="id">ID：{{ id }}</div>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'VideoCard',
  props: {
    thumbnail: String,
    title: String,
    author: String,
    views: [String, Number],
    publishedAt: String,
    publishDate: String,
    duration: String,
    durationSec: [String, Number],
    id: String,
    url: String,
  },
  emits: ['fetch', 'delete'],
  computed: {
    authorInitial() {
      const text = (this.displayAuthor || '').trim()
      return text ? text.slice(0, 1).toUpperCase() : 'X'
    },
    displayAuthor() {
      return (this.author && String(this.author)) || '未知作者'
    },
    youtubeUrl() {
      if (this.url) return this.url
      if (this.id) return `https://www.youtube.com/watch?v=${this.id}`
      return ''
    },
    resolvedUrl() {
      return this.youtubeUrl || ''
    },
    displayViews() {
      const views = this.views
      if (views === undefined || views === null || views === '') {
        return '—'
      }
      
      const numViews = Number(views)
      if (isNaN(numViews)) {
        return String(views)
      }
      
      // 格式化大数字显示
      if (numViews >= 100000000) {
        return (numViews / 100000000).toFixed(1) + '亿次观看'
      } else if (numViews >= 10000) {
        return (numViews / 10000).toFixed(1) + '万次观看'
      } else {
        return numViews.toLocaleString('zh-CN') + '次观看'
      }
    },
    displayDuration() {
      if (this.durationSec !== undefined && this.durationSec !== null && this.durationSec !== '') {
        const sec = Number(this.durationSec)
        if (!isNaN(sec)) {
          const h = Math.floor(sec / 3600)
          const m = Math.floor((sec % 3600) / 60)
          const s = Math.floor(sec % 60)
          const mm = h > 0 ? String(m).padStart(2, '0') : String(m)
          const ss = String(s).padStart(2, '0')
          return h > 0 ? `${h}:${mm}:${ss}` : `${m}分${ss}秒`
        }
      }
      return this.duration || '—'
    },
    displayDate() {
      const dateStr = this.publishDate || this.publishedAt || ''
      if (!dateStr) return ''
      
      // 如果是YYYYMMDD格式，转换为更友好的显示
      if (/^\d{8}$/.test(dateStr)) {
        const year = dateStr.substring(0, 4)
        const month = dateStr.substring(4, 6)
        const day = dateStr.substring(6, 8)
        return `${year}年${parseInt(month)}月${parseInt(day)}日`
      }
      
      // 如果是ISO日期格式，也进行转换
      if (dateStr.includes('-') || dateStr.includes('T')) {
        try {
          const date = new Date(dateStr)
          if (!isNaN(date.getTime())) {
            return date.toLocaleDateString('zh-CN', {
              year: 'numeric',
              month: 'long', 
              day: 'numeric'
            })
          }
        } catch (e) {
          // 转换失败，返回原始字符串
        }
      }
      
      return dateStr
    },
  },
})
</script>





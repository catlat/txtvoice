<template>
  <main>
    <div class="px-4 py-10 sm:px-0">
      <!-- Page Header -->
      <div class="max-w-4xl mx-auto text-center mb-8">
        <div class="text-3xl font-bold tracking-tight text-gray-900">合成历史</div>
        <div class="mt-3 text-lg text-gray-500">查看您的文本转语音历史记录</div>
      </div>

      <!-- Main Content -->
      <div class="max-w-4xl mx-auto">
        <div class="bg-white rounded-xl border border-gray-200 shadow-[0_0_16px_0_rgba(0,0,0,0.06)] overflow-hidden">
          <!-- Loading State -->
          <div v-if="loading" class="p-8 text-center">
            <div class="flex items-center justify-center gap-3">
              <Spinner /> 
              <span class="text-gray-600">正在加载合成历史...</span>
            </div>
          </div>

          <!-- Error State -->
          <div v-else-if="error" class="p-8 text-center">
            <div class="text-red-600 bg-red-50 rounded-lg p-4">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6 mx-auto mb-2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
              </svg>
              {{ error }}
            </div>
          </div>

          <!-- Content -->
          <div v-else-if="items.length" class="divide-y divide-gray-100">
            <div v-for="(it, idx) in items" :key="idx" class="p-6 hover:bg-gray-50 transition-colors duration-200">
              <div class="flex items-start justify-between gap-4">
                <div class="flex-1 min-w-0">
                  <!-- Date -->
                  <div class="flex items-center gap-2 mb-2">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 text-gray-400">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                    </svg>
                    <span class="text-sm text-gray-500">{{ formatDate(it.created_at) }}</span>
                  </div>
                  
                  <!-- Text Preview -->
                  <div class="text-gray-900 leading-relaxed">
                    <p class="line-clamp-3">{{ it.text_preview || '无文本内容' }}</p>
                  </div>
                  
                  <!-- Character Count & Speaker -->
                  <div class="mt-2 text-xs text-gray-500 flex items-center gap-3">
                    <span>{{ it.char_count || 0 }} 字符</span>
                    <span v-if="it.speaker">语音：{{ it.speaker }}</span>
                  </div>
                </div>

                <!-- Audio Player -->
                <div v-if="it.audio_url" class="flex-shrink-0">
                  <div class="bg-gray-50 rounded-lg p-3 min-w-0">
                    <AudioPlayer :src="it.audio_url" />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="p-12 text-center">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-16 h-16 mx-auto mb-4 text-gray-300">
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.114 5.636a9 9 0 0 1 0 12.728M16.463 8.288a5.25 5.25 0 0 1 0 7.424M6.75 8.25l4.72-4.72a.75.75 0 0 1 1.28.53v15.88a.75.75 0 0 1-1.28.53l-4.72-4.72H4.51c-.88 0-1.59-.79-1.59-1.75v-4.5c0-.96.71-1.75 1.59-1.75h2.24Z" />
            </svg>
            <div class="text-xl font-medium text-gray-900 mb-2">暂无合成历史</div>
            <div class="text-gray-500">开始您的第一次文本转语音吧</div>
          </div>

          <!-- Pagination -->
          <div v-if="items.length" class="px-6 py-4 bg-gray-50 border-t border-gray-100 flex items-center justify-between">
            <button 
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200" 
              :disabled="page<=1" 
              @click="changePage(page-1)"
            >
              上一页
            </button>
            <div class="text-sm text-gray-600">第 {{ page }} 页</div>
            <button 
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200" 
              :disabled="items.length<size" 
              @click="changePage(page+1)"
            >
              下一页
            </button>
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<script>
import { defineComponent } from 'vue'
import * as historyApi from '../api/history'
import AudioPlayer from '../components/AudioPlayer.vue'
import Spinner from '../components/Spinner.vue'

export default defineComponent({
  components: { AudioPlayer, Spinner },
  data() {
    return { loading: false, error: '', items: [], page: 1, size: 20 }
  },
  created() { this.fetch() },
  methods: {
    async fetch() {
      this.loading = true; this.error = ''
      try {
        const res = await historyApi.listTTS({ page: this.page, size: this.size })
        const data = (res && res.data) || {}
        this.items = (Array.isArray(data.items) ? data.items : (Array.isArray(res && res.items) ? res.items : []))
      } catch (e) { this.error = e && e.message ? e.message : '加载失败'; try { const { toast } = require('../utils/toast'); toast(this.error, 'error') } catch (e) {} }
      finally { this.loading = false }
    },
    changePage(p) { this.page = p; this.fetch() },
    formatDate(dateStr) {
      if (!dateStr) return '未知时间'
      try {
        const date = new Date(dateStr)
        return date.toLocaleString('zh-CN', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit'
        })
      } catch (e) {
        return dateStr
      }
    },
  },
})
</script>




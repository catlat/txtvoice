<template>
  <main>
    <div class="px-4 py-10 sm:px-0">
      <!-- Page Header -->
      <div class="max-w-4xl mx-auto text-center mb-8">
        <div class="text-3xl font-bold tracking-tight text-gray-900">合成历史</div>
        <div class="mt-3 text-lg text-gray-500">查看您的文本转语音历史记录</div>
      </div>

      <!-- Main Content -->
      <!-- Login Section (if not logged in) -->
      <LoginRequired v-if="!token" title="未登录" description="请登录后查看语音合成历史记录" />

      <div v-else class="max-w-4xl mx-auto">
        <div>
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
          <div v-else-if="items.length" class="p-3 md:p-4">
            <div class="grid grid-cols-1 gap-3 md:gap-4">
              <div v-for="(it, idx) in items" :key="idx" class="bg-white rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition-shadow duration-200">
                <div class="p-4 md:p-5">
                  <!-- Header Information - 紧凑单行布局 -->
                  <div class="mb-3">
                    <!-- 第一行：时间、字符数、音色 -->
                    <div class="flex items-center justify-between text-xs md:text-sm text-gray-500 mb-2">
                      <div class="flex items-center gap-3 md:gap-4">
                        <div class="flex items-center gap-1">
                          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3.5 h-3.5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
                          </svg>
                          <span>{{ formatDate(it.created_at) }}</span>
                        </div>
                        <div class="flex items-center gap-1">
                          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3.5 h-3.5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 0 1 .865-.501 48.172 48.172 0 0 0 3.423-.379c1.584-.233 2.707-1.627 2.707-3.227V6.741c0-1.6-1.123-2.994-2.707-3.227A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.514C3.373 3.747 2.25 5.14 2.25 6.741v6.018Z" />
                          </svg>
                          <span>{{ it.char_count || 0 }} 字符</span>
                        </div>
                        <div v-if="it.speaker" class="flex items-center gap-1 text-gray-600">
                          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-3.5 h-3.5">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M19.114 5.636a9 9 0 0 1 0 12.728M16.463 8.288a5.25 5.25 0 0 1 0 7.424M6.75 8.25l4.72-4.72a.75.75 0 0 1 1.28.53v15.88a.75.75 0 0 1-1.28.53l-4.72-4.72H4.51c-.88 0-1.59-.79-1.59-1.75v-4.5c0-.96.71-1.75 1.59-1.75h2.24Z" />
                          </svg>
                          <span class="font-medium">{{ getVoiceName(it.speaker) }}</span>
                        </div>
                      </div>
                    </div>
                    
                    <!-- 文本预览 - 可展开折叠 -->
                    <div class="text-gray-900 leading-relaxed">
                      <div 
                        class="cursor-pointer select-none"
                        @click="toggleTextExpand(idx)"
                      >
                        <p 
                          :class="[
                            'text-sm md:text-base transition-all duration-200',
                            expandedTexts[idx] ? '' : 'line-clamp-2'
                          ]"
                        >
                          {{ it.text_preview || '无文本内容' }}
                        </p>
                        <div 
                          v-if="isTextLong(it.text_preview)" 
                          class="flex items-center gap-1 mt-1 text-xs text-gray-500 hover:text-gray-700"
                        >
                          <span>{{ expandedTexts[idx] ? '收起' : '展开' }}</span>
                          <svg 
                            xmlns="http://www.w3.org/2000/svg" 
                            fill="none" 
                            viewBox="0 0 24 24" 
                            stroke-width="1.5" 
                            stroke="currentColor" 
                            :class="[
                              'w-3 h-3 transition-transform duration-200',
                              expandedTexts[idx] ? 'rotate-180' : ''
                            ]"
                          >
                            <path stroke-linecap="round" stroke-linejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
                          </svg>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Audio Player Row -->
                  <div v-if="it.audio_url" class="pt-3 md:pt-4 border-t border-gray-100">
                    <AudioPlayer 
                      :src="it.audio_url" 
                      :show-download="true"
                      :show-speed-control="true"
                      :label="it.text_preview"
                    />
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
          <div v-if="items.length" class="px-6 py-4 bg-gray-50 border-t border-gray-100">
            <div class="flex items-center justify-center mb-4 text-sm text-gray-600">
              <span v-if="pagination.total > 0">
                共 {{ pagination.total }} 条记录，第 {{ page }} / {{ pagination.total_pages }} 页
              </span>
              <span v-else>第 {{ page }} 页</span>
            </div>
            <Pagination
              :page="page"
              :page-size="size"
              :total="pagination.total"
              :total-pages="pagination.total_pages"
              :has-prev="pagination.has_prev"
              :has-next="pagination.has_next"
              @update:page="changePage"
            />
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<script>
import { defineComponent } from 'vue'
import { getToken } from '../utils/auth'
import * as historyApi from '../api/history'
import AudioPlayer from '../components/AudioPlayer.vue'
import Spinner from '../components/Spinner.vue'
import Pagination from '../components/Pagination.vue'
import LoginRequired from '../components/LoginRequired.vue'

export default defineComponent({
  components: { AudioPlayer, Spinner, Pagination, LoginRequired },
  data() {
    return { 
      loading: false, 
      error: '', 
      items: [], 
      page: 1, 
      size: 10,
      token: getToken() || '',
      pagination: {
        total: 0,
        total_pages: 0,
        has_next: false,
        has_prev: false
      },
      expandedTexts: {}, // 跟踪展开的文本
      voiceMap: {} // 音色映射表
    }
  },
  created() { 
    this.loadVoices()
    if (this.token) this.fetch() 
  },
  methods: {
    async fetch() {
      this.loading = true; this.error = ''
      try {
        const res = await historyApi.listTTS({ page: this.page, size: this.size })
        const data = (res && res.data) || res || {}
        
        // 处理数据列表
        this.items = (Array.isArray(data.items) ? data.items : [])
        
        // 处理分页信息
        if (data.pagination) {
          this.pagination = {
            total: data.pagination.total || 0,
            total_pages: data.pagination.total_pages || 0,
            has_next: data.pagination.has_next || false,
            has_prev: data.pagination.has_prev || false
          }
        } else {
          // 兼容旧版本API，使用传统方式判断
          this.pagination = {
            total: 0,
            total_pages: 0,
            has_next: this.items.length >= this.size,
            has_prev: this.page > 1
          }
        }
      } catch (e) { 
        this.error = e && e.message ? e.message : '加载失败'
        try { 
          const { toast } = require('../utils/toast')
          toast(this.error, 'error')
        } catch (e) {}
      }
      finally { this.loading = false }
    },
    changePage(p) { 
      if (p < 1 || (this.pagination.total_pages > 0 && p > this.pagination.total_pages)) return
      this.page = p
      // 重置展开状态
      this.expandedTexts = {}
      this.fetch()
    },
    getVisiblePages() {
      const current = this.page
      const total = this.pagination.total_pages
      if (total <= 0) return [current]
      
      const range = 2 // 当前页左右各显示2页
      const start = Math.max(1, current - range)
      const end = Math.min(total, current + range)
      
      const pages = []
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
      return pages
    },
    async loadVoices() {
      try {
        const response = await fetch('/voices.json')
        const voicesData = await response.json()
        
        // 构建音色映射表
        const voiceMap = {}
        Object.keys(voicesData).forEach(providerKey => {
          const provider = voicesData[providerKey]
          if (provider.voice_options && Array.isArray(provider.voice_options)) {
            provider.voice_options.forEach(voice => {
              if (voice.voice_config && Array.isArray(voice.voice_config)) {
                voice.voice_config.forEach(config => {
                  if (config.params && config.params.voice_type) {
                    voiceMap[config.params.voice_type] = voice.name
                  }
                })
              }
            })
          }
        })
        
        this.voiceMap = voiceMap
      } catch (e) {
        console.warn('加载音色数据失败:', e)
      }
    },
    getVoiceName(speaker) {
      if (!speaker) return '未知音色'
      return this.voiceMap[speaker] || speaker
    },
    toggleTextExpand(index) {
      // Vue 3 响应式更新
      this.expandedTexts = {
        ...this.expandedTexts,
        [index]: !this.expandedTexts[index]
      }
    },
    isTextLong(text) {
      if (!text) return false
      // 判断文本是否超过两行（大约80个字符）
      return text.length > 80
    },
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




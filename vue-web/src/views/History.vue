<template>
  <main>
    <div class="px-4 py-10 sm:px-0">
      <!-- Page Header -->
      <div class="max-w-4xl mx-auto text-center mb-8">
        <div class="text-3xl font-bold tracking-tight text-gray-900">历史记录</div>
        <div class="mt-3 text-lg text-gray-500">查看您的视频处理历史</div>
      </div>

      <!-- Main Content -->
      <div class="max-w-4xl mx-auto">
        <div class="bg-white rounded-xl border border-gray-200 shadow-[0_0_16px_0_rgba(0,0,0,0.06)] overflow-hidden">
          <!-- Loading State -->
          <div v-if="loading" class="p-8 text-center">
            <div class="flex items-center justify-center gap-3">
              <Spinner /> 
              <span class="text-gray-600">正在加载历史记录...</span>
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
            <div v-for="(v, idx) in items" :key="idx" class="p-6 hover:bg-gray-50 transition-colors duration-200">
              <VideoCard
                :thumbnail="v.thumbnail || v.thumbnail_url || ''"
                :title="v.title || ''"
                :author="v.author || v.channel || ''"
                :views="v.views ? `${v.views} 次观看` : ''"
                :publishedAt="v.published_at || ''"
                :duration="v.duration || ''"
                @fetch="goDetail(v)"
                class="cursor-pointer hover:shadow-md transition-shadow duration-200"
              />
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="p-12 text-center">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-16 h-16 mx-auto mb-4 text-gray-300">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
            </svg>
            <div class="text-xl font-medium text-gray-900 mb-2">暂无历史记录</div>
            <div class="text-gray-500">开始处理您的第一个视频吧</div>
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

      <!-- Detail Modal -->
      <transition name="fade">
        <div v-if="showDetail" class="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm" @click="showDetail=false"></div>
      </transition>
      <transition name="slide">
        <div v-if="showDetail" class="fixed right-0 top-0 bottom-0 z-50 w-full max-w-2xl bg-white shadow-2xl overflow-auto">
          <div class="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
            <h3 class="text-xl font-semibold text-gray-900">视频详情</h3>
            <button 
              class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-200" 
              @click="showDetail=false"
            >
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div v-if="detail" class="p-6 space-y-6">
            <div class="aspect-video bg-gray-100 rounded-lg overflow-hidden">
              <img :src="detail.thumbnail_url" alt="视频缩略图" class="w-full h-full object-cover">
            </div>
            <div class="space-y-4">
              <div>
                <div class="text-sm font-medium text-gray-500 mb-1">标题</div>
                <div class="text-lg font-medium text-gray-900">{{ detail.title }}</div>
              </div>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <div class="text-sm font-medium text-gray-500 mb-1">频道</div>
                  <div class="text-gray-900">{{ detail.channel_title }}</div>
                </div>
                <div>
                  <div class="text-sm font-medium text-gray-500 mb-1">时长</div>
                  <div class="text-gray-900">{{ detail.duration_sec }}秒</div>
                </div>
              </div>
              <div>
                <div class="text-sm font-medium text-gray-500 mb-1">发布时间</div>
                <div class="text-gray-900">{{ detail.published_at }}</div>
              </div>
              <div v-if="detail.audio_url">
                <div class="text-sm font-medium text-gray-500 mb-2">音频文件</div>
                <a 
                  :href="detail.audio_url" 
                  target="_blank"
                  class="inline-flex items-center gap-2 px-4 py-2 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 transition-colors duration-200"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M13.19 8.688a4.5 4.5 0 0 1 1.242 7.244l-4.5 4.5a4.5 4.5 0 0 1-6.364-6.364l1.757-1.757m13.35-.622 1.757-1.757a4.5 4.5 0 0 0-6.364-6.364l-4.5 4.5a4.5 4.5 0 0 0 1.242 7.244" />
                  </svg>
                  下载音频
                </a>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </div>
  </main>
</template>

<script>
import { defineComponent } from 'vue'
import VideoCard from '../components/VideoCard.vue'
import * as historyApi from '../api/history'
import Spinner from '../components/Spinner.vue'

export default defineComponent({
  components: { VideoCard, Spinner },
  data() {
    return { loading: false, error: '', items: [], page: 1, size: 20, detail: null, showDetail: false }
  },
  created() { this.fetch() },
  methods: {
    async fetch() {
      this.loading = true; this.error = ''
      try {
        const res = await historyApi.listVideos({ page: this.page, size: this.size })
        this.items = (res && (res.items || res.data || [])) || []
      } catch (e) { this.error = e && e.message ? e.message : '加载失败' }
      finally { this.loading = false }
    },
    changePage(p) { this.page = p; this.fetch() },
    async goDetail(v) {
      const site = v.source_site || 'youtube'
      const id = v.video_id || v.id
      this.loading = true; this.error = ''
      try {
        const res = await historyApi.getVideo(site, id)
        this.detail = (res && (res.detail || res)) || null
        this.showDetail = true
      } catch (e) { this.error = e && e.message ? e.message : '加载详情失败' }
      finally { this.loading = false }
    },
  },
})
</script>

<style>
.fade-enter-active,.fade-leave-active{ transition: opacity .2s }
.fade-enter-from,.fade-leave-to{ opacity: 0 }
.slide-enter-active,.slide-leave-active{ transition: transform .2s }
.slide-enter-from,.slide-leave-to{ transform: translateX(100%) }
</style>




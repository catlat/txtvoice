<template>
  <main>
    <div class="px-4 py-6 sm:px-0">
      <div class="p-4 border-4 border-gray-200 border-dashed rounded-lg  text-gray-700">
        <h2 class="text-2xl font-bold">视频信息</h2>
        <p class="mt-2 text-gray-500">根据链接/ID 获取视频信息，并可拉取文本与音频。</p>

        <div class="mt-4 flex items-center gap-2" v-if="loading">
          <Spinner /> <span>加载中...</span>
        </div>
        <div class="mt-4 text-red-600" v-if="error">{{ error }}</div>

        <div class="mt-4" v-if="info">
          <VideoCard
            :thumbnail="infoThumb"
            :title="infoTitle"
            :author="infoAuthor"
            :views="infoViews"
            :publishedAt="infoPublished"
            :duration="infoDuration"
            @fetch="onFetchText"
          />
        </div>

        <div class="mt-4 text-sm text-gray-500" v-else>未找到视频信息，请返回首页重新输入。</div>
      </div>
    </div>
  </main>
</template>

<script>
import { defineComponent } from 'vue'
import VideoCard from '../components/VideoCard.vue'
import * as yt from '../api/yt'
import Spinner from '../components/Spinner.vue'

export default defineComponent({
  components: { VideoCard, Spinner },
  data() {
    return {
      loading: false,
      error: '',
      info: null,
    }
  },
  computed: {
    vid() {
      return this.$route.query.vid || ''
    },
    infoThumb() {
      return (this.info && (this.info.thumbnail || this.info.thumbnail_url || this.info.thumb)) || ''
    },
    infoTitle() {
      return (this.info && (this.info.title || this.info.name)) || ''
    },
    infoAuthor() {
      return (this.info && (this.info.author || this.info.channel || this.info.uploader)) || ''
    },
    infoViews() {
      const v = this.info && (this.info.views || this.info.view_count)
      return v ? `${v} 次观看` : ''
    },
    infoPublished() {
      return (this.info && (this.info.published_at || this.info.publish_date)) || ''
    },
    infoDuration() {
      return (this.info && (this.info.duration || this.info.length_text)) || ''
    },
  },
  created() {
    this.fetchInfo()
  },
  methods: {
    async fetchInfo() {
      if (!this.vid) return
      this.loading = true
      this.error = ''
      try {
        const res = await yt.info(this.vid)
        // 后端可能返回 { data: {...} } 或直接对象
        this.info = res && (res.data || res)
      } catch (e) {
        this.error = e && e.message ? e.message : '获取视频信息失败'
        try { const { toast } = require('../utils/toast'); toast(this.error, 'error') } catch (e) {}
      } finally {
        this.loading = false
      }
    },
    async onFetchText() {
      if (!this.vid) return
      this.loading = true
      this.error = ''
      try {
        const res = await yt.text(this.vid)
        const payload = res && (res.data || res) || {}
        // 缓存供 Edit 页面使用
        const textEn = payload.original_text || payload.text_en || ''
        const textZh = payload.translated_text || payload.text_zh || ''
        if (textEn) sessionStorage.setItem('edit:text_en', textEn)
        if (textZh) sessionStorage.setItem('edit:text_zh', textZh)
        if (payload.utterances) sessionStorage.setItem('edit:utterances', JSON.stringify(payload.utterances))
        if (this.info) sessionStorage.setItem('edit:video_info', JSON.stringify(this.info))
        sessionStorage.setItem('edit:payload', JSON.stringify(payload))
        this.$router.push({ path: '/edit' })
      } catch (e) {
        this.error = e && e.message ? e.message : '获取文本失败'
        try { const { toast } = require('../utils/toast'); toast(this.error, 'error') } catch (e) {}
      } finally {
        this.loading = false
      }
    },
  },
})
</script>



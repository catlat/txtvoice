<template>
  <div class="rounded-[12px] border border-gray-200 bg-white p-4 shadow-sm space-y-4">
    <!-- 标题 -->
    <div class="flex items-center justify-between">
      <h3 class="text-base font-medium text-gray-900">视频结果信息</h3>
      <div class="text-xs text-gray-500">{{ videoTitle || '视频内容' }}</div>
    </div>

    <!-- 字幕文件区域 -->
    <div class="space-y-3">
      <h4 class="text-sm font-medium text-gray-700">字幕文件</h4>
      
      <!-- Bilibili：只显示一个字幕 -->
      <template v-if="platform === 'bilibili'">
        <div v-if="subtitle" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg border">
          <div class="flex items-center gap-2">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4 text-gray-500" fill="currentColor">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8zm4 18H6V4h7v5h5v11zM8 12h8v2H8zm0 4h5v2H8z"/>
            </svg>
            <div class="min-w-0">
              <div class="text-sm font-medium text-gray-900">字幕.txt</div>
              <div class="text-xs text-gray-500">{{ getCharCount(subtitle) }} 字符</div>
            </div>
          </div>
          <div class="flex items-center gap-1">
            <button 
              @click="viewSubtitle('subtitle', subtitle)" 
              class="px-2 py-1 text-xs text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded"
              title="查看"
            >
              查看
            </button>
            <button 
              @click="downloadSubtitle('字幕.txt', subtitle)" 
              class="px-2 py-1 text-xs text-gray-600 hover:text-gray-800 hover:bg-gray-50 rounded"
              title="下载"
            >
              下载
            </button>
            <button 
              @click="writeToInput(subtitle)" 
              class="px-3 py-1 text-xs bg-black text-white rounded hover:opacity-90"
              title="写入输入框"
            >
              写入
            </button>
          </div>
        </div>
      </template>

      <!-- YouTube：显示中文和英文字幕 -->
      <template v-else>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <!-- 中文字幕文件 -->
          <div v-if="chineseSubtitle" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg border">
            <div class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4 text-gray-500" fill="currentColor">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8zm4 18H6V4h7v5h5v11zM8 12h8v2H8zm0 4h5v2H8z"/>
              </svg>
              <div class="min-w-0">
                <div class="text-sm font-medium text-gray-900">中文字幕.txt</div>
                <div class="text-xs text-gray-500">{{ getCharCount(chineseSubtitle) }} 字符</div>
              </div>
            </div>
            <div class="flex items-center gap-1">
              <button 
                @click="viewSubtitle('zh', chineseSubtitle)" 
                class="px-2 py-1 text-xs text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded"
                title="查看"
              >
                查看
              </button>
              <button 
                @click="downloadSubtitle('中文字幕.txt', chineseSubtitle)" 
                class="px-2 py-1 text-xs text-gray-600 hover:text-gray-800 hover:bg-gray-50 rounded"
                title="下载"
              >
                下载
              </button>
              <button 
                @click="writeToInput(chineseSubtitle)" 
                class="px-3 py-1 text-xs bg-black text-white rounded hover:opacity-90"
                title="写入输入框"
              >
                写入
              </button>
            </div>
          </div>

          <!-- 英文字幕文件 -->
          <div v-if="englishSubtitle" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg border">
            <div class="flex items-center gap-2">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4 text-gray-500" fill="currentColor">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8zm4 18H6V4h7v5h5v11zM8 12h8v2H8zm0 4h5v2H8z"/>
              </svg>
              <div class="min-w-0">
                <div class="text-sm font-medium text-gray-900">英文字幕.txt</div>
                <div class="text-xs text-gray-500">{{ getCharCount(englishSubtitle) }} 字符</div>
              </div>
            </div>
            <div class="flex items-center gap-1">
              <button 
                @click="viewSubtitle('en', englishSubtitle)" 
                class="px-2 py-1 text-xs text-blue-600 hover:text-blue-800 hover:bg-blue-50 rounded"
                title="查看"
              >
                查看
              </button>
              <button 
                @click="downloadSubtitle('英文字幕.txt', englishSubtitle)" 
                class="px-2 py-1 text-xs text-gray-600 hover:text-gray-800 hover:bg-gray-50 rounded"
                title="下载"
              >
                下载
              </button>
              <button 
                @click="writeToInput(englishSubtitle)" 
                class="px-3 py-1 text-xs bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
                title="写入输入框"
              >
                写入
              </button>
            </div>
          </div>
        </div>
      </template>

      <!-- 如果没有字幕 -->
      <div v-if="!hasAnySubtitle" class="text-sm text-gray-500 text-center py-4">
        暂无字幕内容
      </div>
    </div>

    <!-- 音频播放器 -->
    <div v-if="audioUrl" class="space-y-2">
      <h4 class="text-sm font-medium text-gray-700">
        原始音频 
        <span v-if="audioUrl.startsWith('data:audio/')" class="text-xs text-blue-600">({{ audioType.toUpperCase() }} Base64数据)</span>
        <span v-else class="text-xs text-green-600">({{ audioType.toUpperCase() }} URL链接)</span>
      </h4>
      <AudioPlayer 
        :src="audioUrl" 
        :audio-type="audioType"
        :compact="true" 
        :show-actions="false" 
        :show-download="true"
        :label="videoTitle || '视频音频'"
      />
    </div>
    <div v-else class="space-y-2">
      <h4 class="text-sm font-medium text-gray-700">原始音频</h4>
      <div class="text-sm text-gray-500 text-center py-4">
        音频加载中...
      </div>
    </div>

    <!-- 操作按钮区域 -->
    <div class="flex items-center justify-between pt-2 border-t border-gray-100">
      <div class="text-xs text-gray-500">
        {{ hasAnySubtitle ? '字幕已准备就绪' : '等待字幕生成...' }}
      </div>
      <div class="flex items-center gap-2">
        <!-- 自动写入字幕按钮 -->
        <button 
          v-if="primarySubtitle && !autoWritten"
          @click="autoWriteSubtitle" 
          class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm"
        >
          {{ platform === 'bilibili' ? '自动写入字幕' : '自动写入中文' }}
        </button>
        <span v-else-if="autoWritten" class="text-xs text-green-600">✓ 已自动写入</span>
      </div>
    </div>
  </div>

  <!-- 字幕查看模态框 -->
  <div v-if="showModal" class="fixed inset-0 bg-gray-900 bg-opacity-40 flex items-center justify-center z-50 p-4" @click="closeModal">
    <div class="bg-white rounded-lg max-w-2xl w-full max-h-[80vh] flex flex-col" @click.stop>
      <div class="flex items-center justify-between p-4 border-b">
        <h3 class="text-lg font-medium">{{ modalTitle }}</h3>
        <button @click="closeModal" class="text-gray-500 hover:text-gray-700">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <div class="p-4 overflow-auto flex-1">
        <div class="bg-gray-50 p-4 rounded-lg">
          <pre class="text-sm text-gray-700 whitespace-pre-wrap font-mono">{{ modalContent }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import AudioPlayer from './AudioPlayer.vue'

export default defineComponent({
  name: 'VideoResultCard',
  components: { AudioPlayer },
  props: {
    chineseSubtitle: String, 
    englishSubtitle: String, 
    subtitle: String, // 用于Bilibili的统一字幕
    platform: String, // 平台类型：youtube/bilibili
    audioUrl: String,
    audioType: { type: String, default: 'm4a' }, // 音频文件类型
    videoTitle: String,
    autoWrite: { type: Boolean, default: true }, // 是否自动写入字幕
  },
  emits: ['write-to-input'],
  data() {
    return {
      showModal: false,
      modalTitle: '',
      modalContent: '',
      modalType: '', // 'zh' or 'en'
      autoWritten: false
    }
  },
  computed: {
    // 判断是否有任何字幕内容
    hasAnySubtitle() {
      if (this.platform === 'bilibili') {
        return !!this.subtitle
      } else {
        return !!(this.chineseSubtitle || this.englishSubtitle)
      }
    },
    // 用于自动写入的主要字幕内容
    primarySubtitle() {
      if (this.platform === 'bilibili') {
        return this.subtitle
      } else {
        // YouTube优先使用中文字幕
        return this.chineseSubtitle || this.englishSubtitle
      }
    }
  },
  watch: {
    primarySubtitle: {
      handler(newValue) {
        if (newValue && this.autoWrite && !this.autoWritten) {
          // 延迟执行自动写入，确保组件完全加载
          this.$nextTick(() => {
            this.autoWriteSubtitle()
          })
        }
      },
      immediate: true
    }
  },
  methods: {
    // 获取字符数
    getCharCount(text) {
      return text ? String(text).length : 0
    },

    // 查看字幕
    viewSubtitle(type, content) {
      this.modalType = type
      if (type === 'zh') {
        this.modalTitle = '中文字幕内容'
      } else if (type === 'en') {
        this.modalTitle = '英文字幕内容'
      } else {
        this.modalTitle = '字幕内容'
      }
      this.modalContent = content || ''
      this.showModal = true
    },

    // 关闭模态框
    closeModal() {
      this.showModal = false
      this.modalTitle = ''
      this.modalContent = ''
      this.modalType = ''
    },

    // 下载字幕
    downloadSubtitle(filename, content) {
      if (!content) return
      
      const blob = new Blob([content], { type: 'text/plain;charset=utf-8' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = filename
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    },

    // 写入输入框
    writeToInput(content) {
      if (!content) return
      this.$emit('write-to-input', content)
    },

    // 自动写入字幕
    autoWriteSubtitle() {
      if (this.primarySubtitle && !this.autoWritten) {
        this.writeToInput(this.primarySubtitle)
        this.autoWritten = true
      }
    }
  }
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-clamp: 2;
}
</style>

<template>
  <main>
    <div class="px-4 pt-12 pb-8 sm:px-0 md:pt-16 md:pb-10">
      <div class="max-w-4xl mx-auto text-center">
        <div class="flex justify-center mb-5 md:mb-7">
          <img src="/logo.svg" alt="vvHub" class="h-28 sm:h-32 md:h-36 lg:h-40 w-auto" />
        </div>
        <div class="mt-1 text-xl sm:text-2xl md:text-3xl lg:text-4xl text-gray-900 font-bold tracking-tight leading-tight">
          一键生成，用你的声音演绎所有内容
        </div>
        <div class="mt-3 text-sm sm:text-base md:text-lg text-gray-600 leading-relaxed max-w-2xl mx-auto">
          视频、音频、文字，都能变成你的声音
        </div>
      </div>

      <div class="max-w-3xl mx-auto mt-10 md:mt-12">
        <div 
          class="flex w-full flex-col items-center rounded-[12px] border p-4 shadow-[0_0_16px_0_rgba(0,0,0,0.06)] transition-all duration-200"
          :class="isDragOver ? 'border-blue-400 border-2 bg-blue-50' : 'border-gray-200'"
          @dragover.prevent="onDragOver"
          @dragenter.prevent="onDragEnter" 
          @dragleave.prevent="onDragLeave"
          @drop.prevent="onDrop"
        >
          <!-- 拖拽遮罩层 -->
          <div v-if="isDragOver" class="absolute inset-0 z-10 flex items-center justify-center bg-blue-50/80 rounded-[12px] border-2 border-blue-400 border-dashed">
            <div class="text-center">
              <div class="text-blue-600 text-xl font-medium mb-2">📁 释放文件以导入内容</div>
              <div class="text-blue-500 text-sm">支持 .txt、.md、.docx 文件</div>
            </div>
          </div>

          <form class="relative w-full" @submit.prevent="onCreate">
            <div class="transition-all duration-300 ease-in-out">
              <div class="grid gap-2">
                <div class="relative">
                  <textarea
                    ref="textareaRef"
                    class="flex w-full bg-transparent text-base outline-none md:text-sm min-h-20 max-h-96 resize-none rounded-none border-none p-0 pr-12 pb-8 text-gray-800 placeholder:text-gray-400 shadow-none overflow-y-auto"
                    v-model="mainText"
                    placeholder="贴上链接或输入文字，即刻生成你的声音"
                    @input="adjustTextareaHeight"
                  />
                  <div class="pointer-events-none absolute bottom-1 right-2 text-[11px] leading-none text-gray-400 bg-white/70 rounded px-1">
                    {{ normalizeAndCount(mainText).count }}
                  </div>
                </div>
                <div class="text-xs text-gray-500">
                  支持 YouTube / 哔哩哔哩，亦可拖拽上传文本
                </div>
              </div>
            </div>

            <!-- 工具栏：左中右三段布局 -->
            <div class="mt-2 w-full">
              <div class="flex items-center justify-between gap-3">
                <!-- 左：文件、音色 -->
                <div class="flex items-center gap-2">
                  <button type="button" class="size-9 inline-flex items-center justify-center text-gray-800 hover:opacity-80" title="上传文件" @click="triggerUpload">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="size-5"><path fill="currentColor" d="M11.434 10.358a.75.75 0 0 1 1.004.052l2.585 2.585a.75.75 0 0 1-1.061 1.06l-1.305-1.304v4.606a.75.75 0 0 1-1.5 0v-4.606l-1.12 1.12a.75.75 0 0 1-1.06-1.06l2.4-2.4z"/><path fill="currentColor" fill-rule="evenodd" d="M13.586 2c.464 0 .91.185 1.237.513l5.164 5.164.117.128c.255.311.396.703.396 1.11V19.25A2.75 2.75 0 0 1 17.75 22H6.25a2.75 2.75 0 0 1-2.75-2.75V4.75A2.75 2.75 0 0 1 6.25 2zM6.25 3.5C5.56 3.5 5 4.06 5 4.75v14.5c0 .69.56 1.25 1.25 1.25h11.5c.69 0 1.25-.56 1.25-1.25v-10h-4.25A1.75 1.75 0 0 1 13 7.5v-4zm8.25 4c0 .138.112.25.25.25h3.19L14.5 4.31z" clip-rule="evenodd"/></svg>
                  </button>
                  <!-- 音色按钮放在左侧（上传后面） -->
                  <div class="ml-1">
                    <VoiceSelector v-model="speaker" />
                  </div>
                </div>


                <!-- 右：创作（靠最右） -->
                <div class="flex items-center gap-2">
                  <button 
                    type="submit" 
                    :disabled="loading || inlineLoading"
                    class="inline-flex items-center gap-2 px-6 py-2 bg-black text-white rounded-md hover:opacity-90 active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-black/20"
                  >
                    <svg 
                      v-if="!loading && !inlineLoading" 
                      xmlns="http://www.w3.org/2000/svg" 
                      viewBox="0 0 24 24" 
                      class="w-4 h-4" 
                      fill="currentColor"
                    >
                      <path d="M8 5v14l11-7z"/>
                    </svg>
                    <svg 
                      v-else 
                      class="animate-spin h-4 w-4" 
                      viewBox="0 0 24 24" 
                      fill="none"
                    >
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                    </svg>
                    <span>{{ (loading || inlineLoading) ? '创作中...' : '创作' }}</span>
                  </button>
                </div>
              </div>
            </div>

            <input ref="fileInput" id="file-input" type="file" accept=".txt,.md,.docx" class="sr-only" @change="onFileChange" />
          </form>

            <!-- 内联视频卡片 -->
            <div v-if="showInlineCard" class="w-full mt-4">
              <!-- 加载状态 -->
              <div v-if="inlineLoading" class="animate-pulse">
                <div class="flex items-start gap-3 md:gap-4 rounded-[12px] border border-gray-200 bg-white p-3 md:p-4">
                  <div class="w-16 h-16 md:w-20 md:h-20 rounded-lg bg-gray-100"></div>
                  <div class="flex-1">
                    <div class="h-4 bg-gray-100 rounded w-3/4 mb-2"></div>
                    <div class="h-3 bg-gray-100 rounded w-1/2 mb-3"></div>
                    <div class="flex items-center gap-3">
                      <div class="h-3 w-20 bg-gray-100 rounded"></div>
                      <div class="h-3 w-16 bg-gray-100 rounded"></div>
            </div>
          </div>
        </div>
      </div>

              <!-- 视频卡片 -->
              <div v-else-if="inlineVideo" class="relative space-y-3">
                <VideoCard
                  :thumbnail="inlineVideo.thumbnail_url"
                  :title="inlineVideo.title"
                  :author="inlineVideo.author || inlineVideo.channel_title"
                  :views="inlineVideo.views"
                  :publishDate="inlineVideo.publish_date || inlineVideo.published_at"
                  :durationSec="inlineVideo.duration_sec"
                  :id="inlineVideo.id || inlineVideo.video_id"
                  :url="inlineVideo.url"
                  @delete="clearInlineCard"
                />

                <!-- 右下角：获取字幕按钮（显眼、带图标与动态效果） -->
                <button
                  v-if="!inlineSubsEn && !inlineSubsZh"
                  @click="fetchInlineSubtitles"
                  :disabled="inlineSubsLoading"
                  class="absolute bottom-2 right-2 z-10 inline-flex items-center gap-2 px-4 py-2 rounded-full bg-black text-white shadow hover:opacity-90 active:scale-[0.98] transition-all duration-200 disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-black/20"
                >
                  <svg v-if="!inlineSubsLoading" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4" fill="currentColor"><path d="M4 5a3 3 0 0 0-3 3v8a3 3 0 0 0 3 3h16a3 3 0 0 0 3-3V8a3 3 0 0 0-3-3zm0 2h16a1 1 0 0 1 1 1v4H3V8a1 1 0 0 1 1-1m-1 7h8v2H3zm10 0h8v2h-8z"/></svg>
                  <svg v-else class="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path></svg>
                  <span>{{ inlineSubsLoading ? '获取中...' : '获取字幕' }}</span>
                </button>
                
                <!-- 视频操作区域结束 -->
            </div>
            
            <!-- 视频结果信息卡片 -->
            <div v-if="showResultCard" class="w-full mt-4">
              <VideoResultCard
                :chinese-subtitle="inlineSubsZh"
                :english-subtitle="inlineSubsEn"
                :audio-url="resultAudioUrl"
                :audio-type="resultAudioType"
                :video-title="inlineVideo?.title"
                :auto-write="true"
                @write-to-input="handleWriteToInput"
              />
            </div>

              <!-- 错误状态 -->
              <div v-if="inlineError" class="p-3 bg-red-50 border border-red-200 rounded-md">
                <div class="flex items-center gap-2">
                  <svg class="w-4 h-4 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
                  </svg>
                  <span class="text-sm text-red-700">{{ inlineError }}</span>
                </div>
              </div>
            </div>

            <!-- 合成状态与播放器 -->
            <div class="w-full mt-4">
              <div v-if="loading" class="text-sm text-gray-500">处理中...</div>
              <div v-if="audioUrl" class="mt-2">
                <AudioPlayer 
                  :src="audioUrl" 
                  :compact="true" 
                  :show-actions="false" 
                  :show-download="true"
                  :label="mainText"
                />
              </div>
            </div>
          </div>
        </div>

    </div>
  </main>
</template>

<script>
import { defineComponent, ref, onMounted, watch, nextTick } from 'vue'
import { parseTxt, parseMd, parseDocx } from '../utils/file'
import * as yt from '../api/yt'
import * as tts from '../api/tts'
import VideoCard from '../components/VideoCard.vue'
import VideoResultCard from '../components/VideoResultCard.vue'
import VoiceSelector from '../components/VoiceSelector.vue'
import AudioPlayer from '../components/AudioPlayer.vue'
import { normalizeAndCount } from '../utils/text'

export default defineComponent({
  components: { VideoCard, VideoResultCard, VoiceSelector, AudioPlayer },
  setup() {
    const mainText = ref('')
    const fileInput = ref(null)
    const textareaRef = ref(null)
    const loading = ref(false)
    const error = ref('')
    const speaker = ref('rec_A')
    const audioUrl = ref('')

    // inline video card states
    const inlineVideo = ref(null)
    const inlineLoading = ref(false)
    const inlineError = ref('')
    const inlineSubsLoading = ref(false)
    const inlineSubsEn = ref('')
    const inlineSubsZh = ref('')
    const showInlineCard = ref(false)
    
    // video result card states
    const showResultCard = ref(false)
    const resultAudioUrl = ref('')
    const resultAudioType = ref('')

    // drag and drop states
    const isDragOver = ref(false)
    const dragCounter = ref(0) // 用于处理嵌套元素的dragenter/dragleave

    function triggerUpload() { fileInput.value && fileInput.value.click() }

    // 自动调整textarea高度
    function adjustTextareaHeight() {
      const textarea = textareaRef.value
      if (!textarea) return
      
      // 重置高度到最小值，然后根据内容调整
      textarea.style.height = 'auto'
      const scrollHeight = textarea.scrollHeight
      
      // 设置最小高度80px (5rem = 20 * 4px)，最大高度384px (24rem = 96 * 4px)
      const minHeight = 80
      const maxHeight = 384
      
      if (scrollHeight <= minHeight) {
        textarea.style.height = `${minHeight}px`
      } else if (scrollHeight >= maxHeight) {
        textarea.style.height = `${maxHeight}px`
      } else {
        textarea.style.height = `${scrollHeight}px`
      }
    }

    // 拖拽事件处理
    function onDragEnter(e) {
      e.preventDefault()
      dragCounter.value++
      if (dragCounter.value === 1) {
        isDragOver.value = true
      }
    }

    function onDragOver(e) {
      e.preventDefault()
      // 设置拖拽效果
      e.dataTransfer.dropEffect = 'copy'
    }

    function onDragLeave(e) {
      e.preventDefault()
      dragCounter.value--
      if (dragCounter.value === 0) {
        isDragOver.value = false
      }
    }

    async function onDrop(e) {
      e.preventDefault()
      isDragOver.value = false
      dragCounter.value = 0

      const files = Array.from(e.dataTransfer.files)
      if (files.length === 0) return

      // 只处理第一个文件
      const file = files[0]
      
      // 验证文件类型
      const name = (file.name || '').toLowerCase()
      if (!(name.endsWith('.txt') || name.endsWith('.md') || name.endsWith('.markdown') || name.endsWith('.docx'))) {
        try { 
          const { toast } = require('../utils/toast')
          toast('仅支持 .txt、.md 或 .docx 文件', 'error') 
        } catch (err) {}
        return
      }

      // 验证文件大小
      if (file.size && file.size > 5 * 1024 * 1024) { // >5MB
        try { 
          const { toast } = require('../utils/toast')
          toast('文件过大，建议不超过 5MB', 'warning') 
        } catch (err) {}
        return
      }

      // 处理文件内容
      await handleDroppedFile(file)
    }

    // 处理拖拽的文件
    async function handleDroppedFile(file) {
      try {
        const name = (file.name || '').toLowerCase()
        let text = ''
        
        if (name.endsWith('.txt')) {
          text = await parseTxt(file)
        } else if (name.endsWith('.docx')) {
          text = await parseDocx(file)
        } else if (name.endsWith('.md') || name.endsWith('.markdown')) {
          text = await parseMd(file)
        }
        
        const value = (text || '').trim()
        if (!value) {
          try { 
            const { toast } = require('../utils/toast')
            toast('文件内容为空', 'error') 
          } catch (err) {}
          return
        }
        
        // 写入输入框并调整高度
        mainText.value = value
        nextTick(() => {
          adjustTextareaHeight()
        })
        
        try { 
          const { toast } = require('../utils/toast')
          toast('文件内容已导入', 'success') 
        } catch (err) {}
        
      } catch (err) {
        const errorMsg = err && err.message ? err.message : '读取文件失败'
        try { 
          const { toast } = require('../utils/toast')
          toast(errorMsg, 'error') 
        } catch (e2) {}
      }
    }

    async function onFileChange(e) {
      const file = e.target && e.target.files && e.target.files[0]
      // 预先缓存 input 元素，避免异步后 e.target 为空
      const inputEl = (e && e.target) ? e.target : (fileInput && fileInput.value) || null
      if (!file) return
      try {
        const name = (file.name || '').toLowerCase()
        if (!(name.endsWith('.txt') || name.endsWith('.md') || name.endsWith('.markdown') || name.endsWith('.docx'))) {
          try { const { toast } = require('../utils/toast'); toast('仅支持 .txt、.md 或 .docx 文件', 'error') } catch (err) {}
          return
        }
        if (file.size && file.size > 5 * 1024 * 1024) { // >5MB
          try { const { toast } = require('../utils/toast'); toast('文件过大，建议不超过 5MB', 'warning') } catch (err) {}
        }
        let text = ''
        if (name.endsWith('.txt')) {
          text = await parseTxt(file)
        } else if (name.endsWith('.docx')) {
          text = await parseDocx(file)
        } else {
          text = await parseMd(file)
        }
        const value = (text || '').trim()
        if (!value) {
          try { const { toast } = require('../utils/toast'); toast('文件内容为空', 'error') } catch (err) {}
          return
        }
        mainText.value = value
      } catch (err) {
        const errorMsg = err && err.message ? err.message : '读取文件失败'
        try { const { toast } = require('../utils/toast'); toast(errorMsg, 'error') } catch (e2) {}
      } finally {
        // 重置 input，避免同一文件无法再次触发 change
        try { if (inputEl) inputEl.value = '' } catch (resetErr) {}
      }
    }


    // 视频链接检测函数
    function detectVideoType(text) {
      const trimmed = (text || '').trim()
      
      // YouTube链接正则：支持多种格式
      const youtubeRegex = /^(https?:\/\/)?(www\.)?(youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/|youtube\.com\/v\/)([\w\-]{6,})([\S]*)?$/i
      if (youtubeRegex.test(trimmed)) {
        return { type: 'youtube', url: trimmed }
      }
      
      // 哔哩哔哩链接正则：支持多种格式和参数
      const bilibiliRegex = /^(https?:\/\/)?(www\.)?(bilibili\.com\/video\/|b23\.tv\/)([\w\-]+)([\S]*)?$/i
      if (bilibiliRegex.test(trimmed)) {
        return { type: 'bilibili', url: trimmed }
      }
      
      return { type: 'text', url: null }
    }

    // 保持向后兼容
    function isYouTubeUrl(text) {
      return detectVideoType(text).type === 'youtube'
    }

    // 检测是否为任何视频链接
    function isVideoUrl(text) {
      const videoType = detectVideoType(text)
      return videoType.type !== 'text'
    }

    // 内联预览视频信息
    async function previewInlineVideo(url) {
      inlineLoading.value = true
      inlineError.value = ''
      try {
        const info = await yt.info(url.trim())
        inlineVideo.value = info && (info.data || info)
        showInlineCard.value = true
        // 链接转卡片后清空输入框
        mainText.value = ''
      } catch (e) {
        inlineError.value = e && e.message ? e.message : '获取视频信息失败'
      } finally {
        inlineLoading.value = false
      }
    }

    // 获取字幕
    async function fetchInlineSubtitles() {
      if (!inlineVideo.value) return
      inlineSubsLoading.value = true
      inlineError.value = ''
      try {
        const res = await yt.text(inlineVideo.value.id || inlineVideo.value.video_id || inlineVideo.value.url || '')
        const payload = res && (res.data || res) || {}
        inlineSubsEn.value = payload.original_text || payload.text_en || ''
        inlineSubsZh.value = payload.translated_text || payload.text_zh || ''
        // 优先使用audio_data（base64数据），否则使用audio_url
        resultAudioUrl.value = payload.audio_data || payload.audio_url || ''
        resultAudioType.value = payload.audio_type || 'm4a'
        
        // 调试日志：检查音频数据格式
        if (payload.audio_data) {
          console.log('使用base64音频数据:', payload.audio_data.substring(0, 50) + '...')
          console.log('音频类型:', payload.audio_type)
        } else if (payload.audio_url) {
          console.log('使用音频URL:', payload.audio_url)
        }
        
        // 获取字幕成功后显示结果卡片
        if (inlineSubsEn.value || inlineSubsZh.value) {
          showResultCard.value = true
        }
      } catch (e) {
        inlineError.value = e && e.message ? e.message : '获取字幕失败'
      } finally {
        inlineSubsLoading.value = false
      }
    }

    // 写入字幕到输入框
    function writeSubtitleToInput(subtitle) {
      if (!subtitle) return
      mainText.value = String(subtitle)
      nextTick(() => {
        adjustTextareaHeight()
      })
    }
    
    // 从结果卡片写入输入框
    function handleWriteToInput(content) {
      writeSubtitleToInput(content)
    }

    // 清除内联卡片
    function clearInlineCard() {
      showInlineCard.value = false
      inlineVideo.value = null
      inlineError.value = ''
      inlineSubsEn.value = ''
      inlineSubsZh.value = ''
      showResultCard.value = false
      resultAudioUrl.value = ''
      resultAudioType.value = ''
      mainText.value = ''
    }

    // 更换链接
    function changeLink() {
      showInlineCard.value = false
      inlineVideo.value = null
      inlineError.value = ''
      inlineSubsEn.value = ''
      inlineSubsZh.value = ''
      showResultCard.value = false
      resultAudioUrl.value = ''
      resultAudioType.value = ''
      // 保持输入框焦点，方便用户输入新链接
    }

    async function onCreate() {
      const text = (mainText.value || '').trim()
      if (!text) { 
        try { const { toast } = require('../utils/toast'); toast('请输入文本或视频链接', 'error') } catch (e) {} 
        return 
      }
      
      // 检测视频链接类型
      const videoType = detectVideoType(text)
      if (videoType.type !== 'text') {
        // 如果是视频链接，预览视频信息
        await previewInlineVideo(videoType.url)
        return
      }
      
      // 普通文本，走原来的TTS流程
      loading.value = true; error.value = ''; audioUrl.value = ''
      try {
        const res = await tts.synthesize({ text, speaker: speaker.value })
        const data = res && (res.data || res)
        audioUrl.value = data && (data.audio_url || data.url) || ''
        if (!audioUrl.value) throw new Error('未返回音频地址')
      } catch (e) {
        error.value = e && e.message ? e.message : '合成失败'
        try { const { toast } = require('../utils/toast'); toast(error.value, 'error') } catch (e) {}
      } finally { loading.value = false }
    }

    // 监听文本内容变化，自动调整高度
    watch(mainText, () => {
      nextTick(() => {
        adjustTextareaHeight()
      })
    })

    // 组件挂载后初始化textarea高度
    onMounted(() => {
      nextTick(() => {
        adjustTextareaHeight()
      })
    })

    return { 
      mainText, fileInput, textareaRef, triggerUpload, adjustTextareaHeight, onFileChange, onCreate, loading, error, speaker, audioUrl, normalizeAndCount,
      // 拖拽相关
      isDragOver, onDragEnter, onDragOver, onDragLeave, onDrop,
      // 内联卡片相关
      inlineVideo, inlineLoading, inlineError, inlineSubsLoading, inlineSubsEn, inlineSubsZh, showInlineCard, previewInlineVideo, fetchInlineSubtitles, writeSubtitleToInput, clearInlineCard, changeLink, 
      // 结果卡片相关
      showResultCard, resultAudioUrl, resultAudioType, handleWriteToInput,
      // 视频链接检测相关
      detectVideoType, isYouTubeUrl, isVideoUrl
    }
  },
})
</script>


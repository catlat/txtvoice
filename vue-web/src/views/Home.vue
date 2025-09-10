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
                  <div v-if="showInlineCard" class="relative group">
                    <div class="flex items-center gap-3 rounded-lg border border-gray-200 bg-white p-2 pr-8 shadow-sm">
                      <img :src="inlineVideo?.thumbnail_url" alt="thumb" class="w-12 h-12 rounded object-cover" />
                      <div class="flex-1 min-w-0">
                        <div class="text-sm font-medium text-gray-900 truncate">{{ inlineVideo?.title }}</div>
                        <div class="text-xs text-gray-500 truncate">{{ inlineVideo?.author || inlineVideo?.channel_title }}</div>
                        <div class="mt-0.5 flex items-center gap-3 text-[11px] text-gray-500">
                          <span v-if="inlineVideo?.duration_sec" class="inline-flex items-center gap-1">
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-3.5 h-3.5" fill="currentColor"><path d="M12 22a10 10 0 1 1 0-20 10 10 0 0 1 0 20m0-18a8 8 0 1 0 0 16 8 8 0 0 0 0-16m.75 4.5v4.19l3.31 1.91a.75.75 0 1 1-.75 1.3l-3.69-2.13A.75.75 0 0 1 11.25 13V8.5a.75.75 0 1 1 1.5 0z"/></svg>
                            {{ formatDuration(inlineVideo?.duration_sec) }}
                          </span>
                          <span v-if="inlineVideo?.publish_date || inlineVideo?.published_at" class="inline-flex items-center gap-1">
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-3.5 h-3.5" fill="currentColor"><path d="M7 2a1 1 0 0 1 1 1v1h8V3a1 1 0 1 1 2 0v1h1a2 2 0 0 1 2 2v2H2V6a2 2 0 0 1 2-2h1V3a1 1 0 1 1 2 0v1z"/><path d="M2 10h20v8a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2z"/></svg>
                            {{ formatDate(inlineVideo?.publish_date || inlineVideo?.published_at) }}
                          </span>
                        </div>
                      </div>
                    </div>
                    <!-- 右上角悬浮删除按钮（X），带动态效果 -->
                    <button type="button" @click="clearInlineCard" class="absolute -top-2 -right-2 inline-flex items-center justify-center w-6 h-6 rounded-full bg-white/90 shadow hover:bg-white transition-all duration-200 hover:rotate-90 hover:scale-110" title="删除">
                      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-3.5 h-3.5" fill="currentColor"><path d="M18.3 5.71a1 1 0 0 0-1.41 0L12 10.59 7.11 5.7A1 1 0 0 0 5.7 7.11L10.59 12l-4.9 4.89a1 1 0 1 0 1.41 1.42L12 13.41l4.89 4.9a1 1 0 0 0 1.42-1.41L13.41 12l4.9-4.89a1 1 0 0 0-.01-1.4z"/></svg>
                    </button>
                  </div>
                  <textarea
                    v-else
                    ref="textareaRef"
                    class="flex w-full bg-transparent text-base outline-none md:text-sm min-h-20 max-h-96 resize-none rounded-none border-none p-0 pr-12 pb-8 text-gray-800 placeholder:text-gray-400 shadow-none overflow-y-auto"
                    v-model="mainText"
                    placeholder="贴上链接或输入文字，即刻生成你的声音"
                    @input="adjustTextareaHeight"
                  />
                  <div v-if="!showInlineCard" class="pointer-events-none absolute bottom-1 right-2 text-[11px] leading-none text-gray-400 bg-white/70 rounded px-1">
                    {{ normalizeAndCount(mainText).count }}
                  </div>
                </div>
                <div v-if="!showInlineCard" class="text-xs text-gray-500">
                  支持 YouTube / 哔哩哔哩，亦可拖拽上传文本
                </div>
              </div>
            </div>

            <!-- 工具栏：左中右三段布局 -->
            <div class="mt-2 w-full">
              <div class="flex items-center justify-between gap-3">
                <!-- 左：文件、音色（卡片显示时保持占位，固定右侧按钮位置） -->
                <div class="flex items-center gap-2" :class="showInlineCard ? 'invisible' : ''">
                  <button type="button" class="size-9 inline-flex items-center justify-center text-gray-800 hover:opacity-80" title="上传文件" @click="triggerUpload">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="size-5"><path fill="currentColor" d="M11.434 10.358a.75.75 0 0 1 1.004.052l2.585 2.585a.75.75 0 0 1-1.061 1.06l-1.305-1.304v4.606a.75.75 0 0 1-1.5 0v-4.606l-1.12 1.12a.75.75 0 0 1-1.06-1.06l2.4-2.4z"/><path fill="currentColor" fill-rule="evenodd" d="M13.586 2c.464 0 .91.185 1.237.513l5.164 5.164.117.128c.255.311.396.703.396 1.11V19.25A2.75 2.75 0 0 1 17.75 22H6.25a2.75 2.75 0 0 1-2.75-2.75V4.75A2.75 2.75 0 0 1 6.25 2zM6.25 3.5C5.56 3.5 5 4.06 5 4.75v14.5c0 .69.56 1.25 1.25 1.25h11.5c.69 0 1.25-.56 1.25-1.25v-10h-4.25A1.75 1.75 0 0 1 13 7.5v-4zm8.25 4c0 .138.112.25.25.25h3.19L14.5 4.31z" clip-rule="evenodd"/></svg>
                  </button>
                  <button ref="voiceBtnRef" type="button" :class="voiceBtnClass" @click="toggleVoicePanel" title="选择音色">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4 mr-1" :class="voiceIconClass" fill="currentColor"><path d="M7 4a1 1 0 0 1 1 1v14a1 1 0 1 1-2 0V5a1 1 0 0 1 1-1zm10 0a1 1 0 0 1 1 1v14a1 1 0 1 1-2 0V5a1 1 0 0 1 1-1zM12 7a1 1 0 0 1 1 1v8a1 1 0 1 1-2 0V8a1 1 0 0 1 1-1z"/></svg>
                    <span class="truncate max-w-[10rem]">{{ currentVoiceLabel }}</span>
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="inline w-4 h-4 ml-1 text-gray-500">
                      <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.24 4.5a.75.75 0 01-1.08 0l-4.24-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
                    </svg>
                  </button>
                </div>


                <!-- 右：创作（靠最右） -->
                <div class="flex items-center gap-2">
                  <button 
                    type="submit" 
                    :disabled="loading || inlineLoading || actionLoading"
                    class="inline-flex items-center gap-2 px-6 py-2 bg-black text-white rounded-md hover:opacity-90 active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-black/20"
                  >
                    <svg 
                      v-if="!loading && !inlineLoading && !actionLoading" 
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
                    <span>{{ actionLabel }}</span>
                  </button>
                </div>
              </div>
            </div>
            <!-- 工具栏下方的音色内联面板 -->
            <div v-if="showVoicePanel" class="mt-2" ref="voicePanelRef">
              <VoiceSelector v-model="speaker" :inline="true" :emit-on-mount="false" @selected="onVoiceSelected" />
            </div>

            <input ref="fileInput" id="file-input" type="file" accept=".txt,.md,.docx" class="sr-only" @change="onFileChange" />
          </form>

          <!-- 隐藏的音色解析器：用于页面刷新后自动同步默认音色的名称与性别到工具栏按钮 -->
          <VoiceSelector v-model="speaker" :inline="true" :emit-on-mount="true" style="display:none" @selected="onVoiceSelected" />

            <!-- 下方独立视频卡片与错误信息（已整合到输入框内，仅保留错误提示） -->
            <div v-if="inlineError" class="w-full mt-4 p-3 bg-red-50 border border-red-200 rounded-md">
              <div class="flex items-center gap-2">
                <svg class="w-4 h-4 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
                </svg>
                <span class="text-sm text-red-700">{{ inlineError }}</span>
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
import { defineComponent, ref, onMounted, onUnmounted, watch, nextTick, computed } from 'vue'
import { parseTxt, parseMd, parseDocx } from '../utils/file'
import * as yt from '../api/yt'
import * as tts from '../api/tts'
import VideoCard from '../components/VideoCard.vue'
import VoiceSelector from '../components/VoiceSelector.vue'
import AudioPlayer from '../components/AudioPlayer.vue'
import { normalizeAndCount } from '../utils/text'

export default defineComponent({
  components: { VideoCard, VoiceSelector, AudioPlayer },
  setup() {
    const mainText = ref('')
    const fileInput = ref(null)
    const textareaRef = ref(null)
    const loading = ref(false)
    const error = ref('')
    const speaker = ref('rec_A')
    const audioUrl = ref('')
    const voicePanelRef = ref(null)
    const voiceBtnRef = ref(null)

    // inline video card states
    const inlineVideo = ref(null)
    const inlineLoading = ref(false)
    const inlineError = ref('')
    const showInlineCard = ref(false)
    // 标记：视频字幕已写回，待完成创作（TTS）
    const videoFlowReady = ref(false)
    // 按钮级加载（用于继续创作/完成创作），不影响视频卡片骨架屏
    const actionLoading = ref(false)
    
    // video result card states（已移除）

    // voice selector inline panel toggle state
    const showVoicePanel = ref(false)
    function toggleVoicePanel() { showVoicePanel.value = !showVoicePanel.value }
    const selectedVoiceMeta = ref({ label: '', gender: '' })
    const currentVoiceLabel = computed(() => selectedVoiceMeta.value.label || '音色')
    const voiceBtnClass = computed(() => {
      const base = 'px-2 py-2 rounded-md text-sm bg-transparent border-0 focus:outline-none focus:ring-1 inline-flex items-center hover:bg-gray-50'
      const info = getCurrentVoiceInfo()
      if (!info) return base + ' text-gray-800 focus:ring-gray-300'
      const g = (info.gender || '').toLowerCase()
      if (g.includes('女') || g.includes('female')) return base + ' text-pink-600 focus:ring-pink-300'
      if (g.includes('男') || g.includes('male')) return base + ' text-blue-600 focus:ring-blue-300'
      return base + ' text-gray-700 focus:ring-gray-300'
    })
    const voiceIconClass = computed(() => {
      const info = getCurrentVoiceInfo()
      if (!info) return 'text-gray-700'
      const g = (info.gender || '').toLowerCase()
      if (g.includes('女') || g.includes('female')) return 'text-pink-600'
      if (g.includes('男') || g.includes('male')) return 'text-blue-600'
      return 'text-gray-700'
    })
    function getCurrentVoiceInfo() {
      // 从 VoiceSelector 的默认集合或远端集合无法在此直接获取，这里仅做占位：仅以 speaker 值显示
      return selectedVoiceMeta.value
    }

    function onVoiceSelected(payload) {
      try {
        if (payload && typeof payload === 'object') {
          selectedVoiceMeta.value = { label: payload.label || '', gender: payload.gender || '' }
        }
      } catch (_) {}
      showVoicePanel.value = false
    }
    async function preloadVoiceMeta() {
      try {
        const res = await fetch('/voices.json', { cache: 'no-store' })
        if (!res.ok) return
        const json = await res.json()
        const list = []
        const pushFrom = (item) => {
          if (!item) return
          const value = item.value || (item.voice_config && item.voice_config[0] && item.voice_config[0].params && item.voice_config[0].params.voice_type) || ''
          const label = item.name || item.label || ''
          const gender = item.gender || item.sex || ''
          if (value && label) list.push({ value, label, gender })
        }
        if (Array.isArray(json)) {
          json.forEach(pushFrom)
        } else if (json && typeof json === 'object') {
          Object.keys(json).forEach(k => {
            const grp = json[k]
            if (grp && Array.isArray(grp.voice_options)) grp.voice_options.forEach(pushFrom)
          })
        }
        const found = list.find(v => v.value === speaker.value)
        if (found) selectedVoiceMeta.value = { label: found.label, gender: found.gender }
      } catch (_) { /* ignore */ }
    }

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


    // 检查是否有字幕内容
    // 已不再需要检测卡片内字幕存在（改为主按钮流）

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
        const videoType = detectVideoType(url.trim())
        const info = await yt.info(url.trim(), videoType.type)
        const videoInfo = info && (info.data || info)
        // 保存原始输入用于后续平台检测
        if (videoInfo) {
          videoInfo.original_input = url.trim()
          videoInfo.detected_platform = videoType.type
        }
        inlineVideo.value = videoInfo
        showInlineCard.value = true
        // 链接转卡片后清空输入框
        mainText.value = ''
      } catch (e) {
        inlineError.value = e && e.message ? e.message : '获取视频信息失败'
      } finally {
        inlineLoading.value = false
      }
    }

    // 获取字幕（优化：直接写回输入框，不显示结果卡片/音频）
    async function fetchInlineSubtitles() {
      if (!inlineVideo.value) return
      actionLoading.value = true
      inlineError.value = ''
      try {
        const videoUrl = inlineVideo.value.id || inlineVideo.value.video_id || inlineVideo.value.url || ''
        // 使用已保存的平台类型或重新检测
        const platform = inlineVideo.value.detected_platform || 
                         detectVideoType(inlineVideo.value.original_input || videoUrl).type
        const res = await yt.text(videoUrl, platform)
        const payload = res && (res.data || res) || {}
        const zh = payload && payload.translated_text || ''
        if (!zh) throw new Error('未返回字幕')
        // 直接写回输入框
        writeSubtitleToInput(zh)
        // 清理视频卡片，并标记可“完成创作”
        showInlineCard.value = false
        inlineVideo.value = null
        videoFlowReady.value = true
      } catch (e) {
        inlineError.value = e && e.message ? e.message : '获取字幕失败'
      } finally {
        actionLoading.value = false
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
    
    // 从结果卡片写入输入框（已移除）

    // 清除内联卡片
    function clearInlineCard() {
      showInlineCard.value = false
      inlineVideo.value = null
      inlineError.value = ''
      // 清理旧字幕状态（已移除具体字段）
      mainText.value = ''
    }

    // 更换交互已移除（保留右上角删除）

    async function onCreate() {
      const text = (mainText.value || '').trim()
      // 若存在视频卡片，执行“继续创作”：获取字幕
      if (showInlineCard.value) {
        await fetchInlineSubtitles()
        return
      }
      if (!text) { 
        try { const { toast } = require('../utils/toast'); toast('请输入文本或视频链接', 'error') } catch (e) {} 
        return 
      }
      // 检测视频链接类型
      const videoType = detectVideoType(text)
      if (videoType.type !== 'text') {
        // 如果是视频链接，预览视频信息（按钮进入“继续创作”态）
        await previewInlineVideo(videoType.url)
        return
      }
      // 普通文本 或 “完成创作”阶段：进行TTS合成
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
      // 输入被清空时，恢复为“开始创作”
      if (!((mainText.value || '').trim())) {
        videoFlowReady.value = false
      }
    })

    // 文档点击（点面板外关闭）与初始化textarea
    const onDocumentClick = (e) => {
      if (!showVoicePanel.value) return
      const panel = voicePanelRef && voicePanelRef.value
      const btn = voiceBtnRef && voiceBtnRef.value
      const t = e && e.target
      if (panel && panel.contains(t)) return
      if (btn && btn.contains(t)) return
      showVoicePanel.value = false
    }
    onMounted(() => {
      nextTick(() => {
        adjustTextareaHeight()
      })
      document.addEventListener('click', onDocumentClick, true)
      preloadVoiceMeta()
    })
    onUnmounted(() => {
      document.removeEventListener('click', onDocumentClick, true)
    })

    const actionLabel = computed(() => {
      if (loading.value || inlineLoading.value || actionLoading.value) return '创作中...'
      if (showInlineCard.value) return '继续创作'
      const hasText = ((mainText.value || '').trim().length > 0)
      if (hasText && !!audioUrl.value) return '重新创作'
      if (videoFlowReady.value) return '完成创作'
      return '开始创作'
    })

    // 展示用格式化
    function formatDuration(sec) {
      const s = Number(sec) || 0
      const h = Math.floor(s / 3600)
      const m = Math.floor((s % 3600) / 60)
      const ss = s % 60
      if (h > 0) return `${h}:${String(m).padStart(2,'0')}:${String(ss).padStart(2,'0')}`
      return `${m}:${String(ss).padStart(2,'0')}`
    }
    function formatDate(input) {
      const s = (input || '').toString().trim()
      if (!s) return ''
      // 简化处理：若是ISO日期，显示YYYY-MM-DD；否则原样
      const d = new Date(s)
      if (!isNaN(d.getTime())) {
        const y = d.getFullYear()
        const mm = String(d.getMonth()+1).padStart(2,'0')
        const dd = String(d.getDate()).padStart(2,'0')
        return `${y}-${mm}-${dd}`
      }
      return s
    }
    // 已移除播放量显示

    return { 
      mainText, fileInput, textareaRef, triggerUpload, adjustTextareaHeight, onFileChange, onCreate, loading, error, speaker, audioUrl, normalizeAndCount,
      // 拖拽相关
      isDragOver, onDragEnter, onDragOver, onDragLeave, onDrop,
      // 内联卡片相关
      inlineVideo, inlineLoading, inlineError, showInlineCard, previewInlineVideo, fetchInlineSubtitles, writeSubtitleToInput, clearInlineCard, actionLabel, actionLoading, formatDuration, formatDate,
      // 视频链接检测相关
      detectVideoType, isYouTubeUrl, isVideoUrl,
      // 音色面板开关
      showVoicePanel, toggleVoicePanel, voicePanelRef, voiceBtnRef, currentVoiceLabel, voiceBtnClass, voiceIconClass, onVoiceSelected,
    }
  },
})
</script>


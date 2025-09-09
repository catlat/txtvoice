<template>
  <div class="flex items-center gap-3 w-full">
    <!-- Compact custom player -->
    <template v-if="compact">
      <div class="flex items-center gap-2 w-full rounded-xl border border-gray-200 bg-white px-2 py-2 shadow-sm">
      <button type="button" class="w-10 h-10 inline-flex items-center justify-center rounded-full bg-white text-gray-900 border border-gray-300 hover:bg-gray-100 transition shrink-0" @click="toggle" aria-label="toggle">
        <svg v-if="!playing" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="size-5" fill="currentColor"><path d="M5 3.879v16.242a1 1 0 0 0 1.555.832l12.121-8.12a1 1 0 0 0 0-1.664L6.555 3.047A1 1 0 0 0 5 3.879Z"/></svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="size-5" fill="currentColor"><path d="M6 4.5A1.5 1.5 0 0 1 7.5 3h1A1.5 1.5 0 0 1 10 4.5v15A1.5 1.5 0 0 1 8.5 21h-1A1.5 1.5 0 0 1 6 19.5zM14 4.5A1.5 1.5 0 0 1 15.5 3h1A1.5 1.5 0 0 1 18 4.5v15A1.5 1.5 0 0 1 16.5 21h-1A1.5 1.5 0 0 1 14 19.5z"/></svg>
      </button>
      <div class="flex-1 min-w-0">
        <!-- text row above progress -->
        <div class="h-5 mb-1 flex items-center justify-center text-[12px] text-gray-800">
          <div class="relative w-full max-w-[80%] overflow-hidden text-center tracking-wide">
            <div class="inline-block whitespace-nowrap" :class="marqueeActive ? 'ap-marquee' : ''" :style="marqueeActive ? marqueeStyle : null">{{ labelToShow }}</div>
          </div>
        </div>
        <!-- progress -->
        <div class="relative h-6 select-none">
          <input type="range" min="0" :max="duration" step="0.01" :value="currentTime" @input="onSeek" class="ap-range w-full h-6 appearance-none bg-transparent cursor-pointer">
          <div class="absolute left-0 right-0 top-1/2 -translate-y-1/2 h-[3px] bg-gray-300 rounded-full"></div>
          <div class="absolute left-0 top-1/2 -translate-y-1/2 h-[3px] rounded-full bg-gray-900" :style="{ width: progressPercent + '%' }"></div>
        </div>
      </div>
      <div class="w-12 text-[10px] text-right text-gray-600 tabular-nums shrink-0">{{ timeText }}</div>
      <button v-if="showDownload" type="button" class="px-3 py-1 text-xs bg-black text-white rounded-md hover:opacity-90 shrink-0" @click="download">下载</button>
      <audio ref="audioEl" :src="src" class="hidden"></audio>
      </div>
    </template>

    <!-- Default native controls -->
    <template v-else>
      <audio ref="audioEl" :src="src" controls class="flex-1"></audio>
      <template v-if="showActions">
        <button class="px-3 py-1 text-sm text-white bg-gray-700 rounded  hover:bg-gray-600" @click="download">下载</button>
        <button class="px-3 py-1 text-sm text-white bg-gray-700 rounded  hover:bg-gray-600" @click="copy">复制链接</button>
      </template>
    </template>
  </div>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'AudioPlayer',
  props: {
    src: { type: String, required: true },
    audioType: { type: String, default: '' }, // 音频文件类型，如 mp3, m4a, wav 等
    label: { type: String, default: '' },
    compact: { type: Boolean, default: false },
    showActions: { type: Boolean, default: true },
    showDownload: { type: Boolean, default: false },
  },
  data() {
    return {
      playing: false,
      currentTime: 0,
      duration: 0,
    }
  },
  computed: {
    progressPercent() { return this.duration ? Math.min(100, (this.currentTime / this.duration) * 100) : 0 },
    timeText() {
      const fmt = (s) => {
        if (!isFinite(s)) return '0:00'
        const m = Math.floor(s / 60); const ss = Math.floor(s % 60)
        return `${m}:${ss < 10 ? '0' + ss : ss}`
      }
      return `${fmt(this.currentTime)}/${fmt(this.duration)}`
    },
    defaultLabel() { return this.src || '' },
    labelToShow() { return this.label || this.defaultLabel },
    labelShouldScroll() { return this.labelToShow.length > 18 },
    marqueeActive() { return this.labelShouldScroll && this.playing },
    marqueeStyle() {
      // 更慢的滚动：基于时长 x1.5，且最低 16s/轮，最高 40s/轮
      const dur = this.duration || 16
      const baseSec = Math.max(16, Math.min(40, Math.ceil(dur * 1.5)))
      return { animationDuration: `${baseSec}s` }
    },
  },
  mounted() {
    const a = this.$refs.audioEl
    if (!a) return
    a.crossOrigin = 'anonymous'
    a.preload = 'auto'
    a.addEventListener('timeupdate', () => { this.currentTime = a.currentTime })
    a.addEventListener('loadedmetadata', () => { this.duration = a.duration })
    a.addEventListener('play', () => { this.playing = true })
    a.addEventListener('pause', () => { this.playing = false })
  },
  methods: {
    toggle() {
      const a = this.$refs.audioEl
      if (!a) return
      if (a.paused) a.play().catch(()=>{})
      else a.pause()
    },
    onSeek(e) {
      const a = this.$refs.audioEl
      if (!a) return
      const v = Number(e.target.value)
      try { a.currentTime = v } catch (err) {}
    },
    download() {
      if (!this.src) return
      const a = document.createElement('a')
      a.href = this.src
      
      // 使用后端返回的音频类型，或从 URL 推断
      if (this.audioType) {
        a.download = `audio.${this.audioType}`
      } else if (this.src.startsWith('data:audio/')) {
        const mimeMatch = this.src.match(/data:audio\/([^;]+)/)
        const extension = mimeMatch ? mimeMatch[1] : 'm4a'
        a.download = `audio.${extension}`
      } else {
        a.download = ''
      }
      
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
    },
    async copy() {
      try {
        // 对于 base64 数据，只复制简短说明
        if (this.src.startsWith('data:audio/')) {
          const format = this.audioType || 'unknown'
          await navigator.clipboard.writeText(`Base64 音频数据 (${format} 格式)`)
        } else {
          await navigator.clipboard.writeText(this.src)
        }
      } catch (e) {
        // noop
      }
    },
  },
})
</script>

<style>
@keyframes ap-marquee { 0% { transform: translateX(0) } 100% { transform: translateX(-100%) } }
.ap-marquee { animation-name: ap-marquee; animation-timing-function: linear; animation-iteration-count: infinite; padding-right: 100%; }
/* prettier thumb for range */
.ap-range::-webkit-slider-thumb { -webkit-appearance: none; appearance: none; width: 14px; height: 14px; border-radius: 9999px; background: #111; border: 2px solid #fff; box-shadow: 0 0 0 1px rgba(0,0,0,.1); cursor: pointer }
.ap-range::-moz-range-thumb { width: 14px; height: 14px; border-radius: 9999px; background: #111; border: 2px solid #fff; box-shadow: 0 0 0 1px rgba(0,0,0,.1); cursor: pointer; border: none }
.ap-range::-moz-range-track { background: transparent; border: none }
</style>





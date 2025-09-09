<template>
  <div class="w-full">
    <div class="flex items-center w-full gap-3">
      <media-player
        class="w-full rounded-xl border border-gray-200 bg-white shadow-sm"
        v-bind:src.prop="vidSrc"
        v-bind:title.prop="labelToShow"
        view-type="audio"
        lang="zh-CN"
        crossorigin
      >
        <media-provider></media-provider>
        <media-audio-layout
          class="text-gray-800 [--media-accent:#2563eb] [--media-controls-color:#1f2937] [--media-button-hover-bg:rgb(0_0_0_/_0.06)] [--media-slider-track-bg:rgb(0_0_0_/_0.12)] [--media-slider-track-progress-bg:rgb(0_0_0_/_0.3)] [--media-slider-thumb-border:1px_solid_#d1d5db] [--media-slider-thumb-bg:#ffffff] [--media-slider-track-height:6px] [--media-slider-thumb-size:14px]"
          v-bind:translations.prop="zhCNTranslations"
          v-bind:colorScheme.prop="'light'"
        ></media-audio-layout>
      </media-player>
      <button v-if="showDownload" type="button" class="px-3 py-1 text-xs bg-black text-white rounded-md hover:opacity-90 shrink-0 dark:bg-white dark:text-black" @click="download">下载</button>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import 'vidstack/bundle'
import 'vidstack/icons'

export default defineComponent({
  name: 'AudioPlayer',
  props: {
    src: { type: String, required: true },
    audioType: { type: String, default: '' },
    label: { type: String, default: '' },
    compact: { type: Boolean, default: false },
    showActions: { type: Boolean, default: true },
    showDownload: { type: Boolean, default: false },
  },
  computed: {
    vidSrc() {
      return this.src
    },
    defaultLabel() { return this.src || '' },
    labelToShow() { return this.label || this.defaultLabel },
  },
  data() {
    return {
      zhCNTranslations: {
        Play: '播放',
        Pause: '暂停',
        Replay: '重播',
        Continue: '继续',
        Mute: '静音',
        Unmute: '取消静音',
        'Enter Fullscreen': '进入全屏',
        'Exit Fullscreen': '退出全屏',
        Fullscreen: '全屏',
        'Enter PiP': '进入画中画',
        'Exit PiP': '退出画中画',
        PiP: '画中画',
        'Seek Forward': '快进',
        'Seek Backward': '快退',
        'Skip To Live': '跳到直播',
        LIVE: '直播',
        Download: '下载',
        Captions: '字幕',
        'Closed-Captions On': '字幕开',
        'Closed-Captions Off': '字幕关',
        'Captions look like this': '字幕示例效果',
        Chapters: '章节',
        'Caption Styles': '字幕样式',
        Font: '字体',
        Text: '文本',
        'Text Background': '文本背景',
        'Display Background': '显示背景',
        Reset: '重置',
        Track: '音轨',
        Boost: '增强',
        Accessibility: '无障碍',
        Audio: '音频',
        Default: '默认',
        Playback: '播放',
        Speed: '速度',
        Normal: '正常',
        Quality: '画质',
        Auto: '自动',
        Settings: '设置',
        Volume: '音量',
        Seek: '跳转',
        'Current time': '当前时间',
        Duration: '总时长',
      },
    }
  },
  methods: {
    download() {
      if (!this.src) return
      const a = document.createElement('a')
      a.href = this.src
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
  },
})
</script>

<style>
/* 样式主要由 Tailwind 与 Vidstack 主题控制。 */
/* 关闭标题滚动（marquee）效果，并隐藏重复标题副本 */
:where(.vds-audio-layout .vds-title.vds-marquee .vds-title-text) {
  animation: none !important;
  transform: none !important;
}
:where(.vds-audio-layout .vds-title.vds-marquee .vds-title-text:nth-child(2)) {
  display: none !important;
}
/* 完全隐藏标题区域 */
:where(.vds-audio-layout .vds-title) {
  display: none !important;
}

/* 始终显示时间进度条（暂停前也可见，固定宽度） */
:where(.vds-audio-layout .vds-time-slider) {
  opacity: 1 !important;
  max-width: 100% !important;
  transform: none !important;
  visibility: visible !important;
}

/* 播放时为进度填充添加轻微动态效果（兼容不支持时自动回退为纯色） */
@keyframes vds-track-move {
  0% { background-position: 0 0 }
  100% { background-position: -200% 0 }
}
:where([data-playing] .vds-audio-layout .vds-time-slider .vds-slider-track-fill) {
  background-image: linear-gradient(
    90deg,
    var(--media-brand) 0%,
    var(--media-brand) 40%,
    rgba(37, 99, 235, 0.7) 60%,
    var(--media-brand) 100%
  ) !important;
  background-size: 200% 100% !important;
  animation: vds-track-move 2.2s linear infinite !important;
}
</style>





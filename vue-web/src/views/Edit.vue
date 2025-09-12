<template>
  <main>
    <div class="px-4 py-6 sm:px-0">
      <div class="grid gap-6 md:grid-cols-2">
        <div class="p-4 border-4 border-gray-200 border-dashed rounded-lg  text-gray-700">
          <h3 class="text-lg font-semibold">英文原文（只读）</h3>
          <div class="mt-3">
            <TextEditor :model-value="textEn" readonly placeholder="英文文本" />
          </div>
        </div>
        <div class="p-4 border-4 border-gray-200 border-dashed rounded-lg  text-gray-700">
          <h3 class="text-lg font-semibold">中文编辑</h3>
          <div class="mt-3">
            <TextEditor v-model="textZh" placeholder="中文可编辑文本" />
          </div>
          <div class="flex items-center justify-between mt-2">
            <CharCounter :count="charCountZh" />
            <div class="flex items-center gap-2">
              <!-- 我的声音VIP按钮 -->
              <button 
                type="button" 
                :class="useMyVoice ? 'bg-gradient-to-r from-gray-800 to-gray-900 text-amber-400 border-2 border-amber-400 shadow-lg' : 'bg-gray-100 text-gray-700 border border-gray-200 hover:bg-gray-200'"
                class="px-3 py-2 rounded-lg text-sm inline-flex items-center transition-all duration-300 relative" 
                @click="toggleMyVoice" 
                title="使用我的声音 - VIP专享">
                <!-- VIP皇冠图标 -->
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" :class="useMyVoice ? 'text-amber-400' : 'text-amber-500'" class="w-4 h-4 mr-2" fill="currentColor">
                  <path d="M5 16L3 6l5.5 4L12 4l3.5 6L21 6l-2 10H5zm2.7-2h8.6l.9-4.4-2.4 1.6L12 8l-2.8 3.2-2.4-1.6L7.7 14z"/>
                </svg>
                <span class="font-medium">我的声音</span>
              </button>
              <VoiceSelector v-model="speaker" :disabled="useMyVoice" />
            </div>
          </div>
          <div class="mt-3 flex items-center gap-2">
            <GenerateButton :loading="loading" :disabled="!textZh" @click="onSynthesize" />
            <Spinner v-if="loading" />
          </div>
        </div>
      </div>

      <div class="mt-6" v-if="audioUrl">
        <AudioPlayer :src="audioUrl" />
      </div>
      <div class="mt-2 text-red-600" v-if="error">{{ error }}</div>
    </div>

    <!-- 登录弹框 -->
    <LoginModal v-model="showLogin" @success="showLogin = false" />
  </main>
</template>

<script>
import { defineComponent } from 'vue'
import TextEditor from '../components/TextEditor.vue'
import CharCounter from '../components/CharCounter.vue'
import VoiceSelector from '../components/VoiceSelector.vue'
import GenerateButton from '../components/GenerateButton.vue'
import AudioPlayer from '../components/AudioPlayer.vue'
import * as tts from '../api/tts'
import { normalizeAndCount } from '../utils/text'
import Spinner from '../components/Spinner.vue'
import LoginModal from '../components/LoginModal.vue'
import { getToken } from '../utils/auth'

export default defineComponent({
  components: { TextEditor, CharCounter, VoiceSelector, GenerateButton, AudioPlayer, Spinner, LoginModal },
  data() {
    return {
      textEn: sessionStorage.getItem('edit:text_en') || '',
      textZh: sessionStorage.getItem('edit:text_zh') || '',
      speaker: 'rec_A',
      loading: false,
      error: '',
      audioUrl: '',
      useMyVoice: true, // 默认使用我的声音
      showLogin: false,
    }
  },
  computed: {
    charCountZh() {
      return normalizeAndCount(this.textZh).count
    },
  },
  methods: {
    toggleMyVoice() {
      this.useMyVoice = !this.useMyVoice
    },
    async onSynthesize() {
      if (!getToken()) { this.showLogin = true; return }
      this.loading = true
      this.error = ''
      this.audioUrl = ''
      try {
        const payload = { text: this.textZh, speaker: this.speaker, use_my_voice: this.useMyVoice }
        const res = await tts.synthesize(payload)
        const data = res && (res.data || res)
        const url = data && (data.url || data.audio_url)
        if (url) {
          this.audioUrl = url
        } else {
          throw new Error('未返回音频地址')
        }
      } catch (e) {
        this.error = e && e.message ? e.message : '合成失败'
        try { const { toast } = require('../utils/toast'); toast(this.error, 'error') } catch (e) {}
      } finally {
        this.loading = false
      }
    },
  },
})
</script>


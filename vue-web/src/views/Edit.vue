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
            <VoiceSelector v-model="speaker" />
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

export default defineComponent({
  components: { TextEditor, CharCounter, VoiceSelector, GenerateButton, AudioPlayer, Spinner },
  data() {
    return {
      textEn: sessionStorage.getItem('edit:text_en') || '',
      textZh: sessionStorage.getItem('edit:text_zh') || '',
      speaker: 'rec_A',
      loading: false,
      error: '',
      audioUrl: '',
    }
  },
  computed: {
    charCountZh() {
      return normalizeAndCount(this.textZh).count
    },
  },
  methods: {
    async onSynthesize() {
      this.loading = true
      this.error = ''
      this.audioUrl = ''
      try {
        const payload = { text: this.textZh, speaker: this.speaker }
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


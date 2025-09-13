<template>
  <!-- Inline panel mode -->
  <div v-if="inline" class="w-full" ref="root">
    <div class="rounded-lg border border-gray-200 bg-white shadow-sm overflow-hidden">
      <div class="p-3 border-b border-gray-50">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-sm font-medium text-gray-900">选择音色</h3>
        </div>
        <div class="relative mb-3">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="size-4 text-gray-400" fill="currentColor">
              <path fill-rule="evenodd" d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z" clip-rule="evenodd" />
            </svg>
          </div>
          <input v-model="keyword" class="w-full pl-10 pr-4 py-2 text-sm border border-gray-200 rounded-lg outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500" placeholder="搜索音色名称、场景或特征..." />
          <button v-if="keyword" class="absolute inset-y-0 right-0 pr-3 flex items-center" @click="keyword=''">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="size-4 text-gray-400 hover:text-gray-600" fill="currentColor">
              <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
            </svg>
          </button>
        </div>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <select v-model="selectedScene" class="px-2 py-1 text-xs border border-gray-200 rounded-md bg-white text-gray-700 focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none">
              <option value="全部">场景</option>
              <option v-for="s in sceneOptions" :key="s" :value="s">{{ s }}</option>
            </select>
            <select v-model="genderFilter" class="px-2 py-1 text-xs border border-gray-200 rounded-md bg-white text-gray-700 focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none">
              <option value="全部">性别</option>
              <option v-for="g in genderOptions" :key="g" :value="g">{{ g }}</option>
            </select>
            <select v-model="ageFilter" class="px-2 py-1 text-xs border border-gray-200 rounded-md bg-white text-gray-700 focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none">
              <option value="全部">年龄</option>
              <option v-for="a in ageOptions" :key="a" :value="a">{{ a }}</option>
            </select>
            <button v-if="hasActiveFilters" class="px-2 py-1 text-xs text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-md transition" @click="resetFilters">
              清空筛选
            </button>
          </div>
          <div class="text-xs text-gray-400">共 {{ filteredOptions.length }} 项</div>
        </div>
      </div>
      <div class="max-h-80 overflow-auto p-2">
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
          <div v-for="opt in filteredOptions" :key="opt.value" class="relative border rounded-[calc(var(--radius,8px)+4px)] border-gray-200 p-3 hover:bg-gray-50 transition cursor-pointer" :class="selectedCardClass(opt)" @click="choose(opt)">
            <div v-if="localValue===opt.value" class="absolute -top-1 -left-1 size-4 rounded-full flex items-center justify-center ring-2 ring-white z-10" :class="selectedBadgeClass(opt)">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="size-3 text-white" fill="currentColor">
                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
              </svg>
            </div>
            <div class="flex items-start justify-between gap-2 pt-0.5 pl-0.5">
              <div class="min-w-0 flex-1">
                <div class="text-sm truncate flex items-center gap-1" :class="{'font-medium': selectedValue===opt.value}">
                  <span class="truncate">{{ opt.label }}</span>
                  <span v-if="opt.tags && opt.tags.includes('上新')" class="shrink-0 inline-flex items-center px-1 text-[10px] rounded bg-pink-100 text-pink-600">新</span>
                </div>
                <div class="text-[10px] text-gray-400 truncate mt-0.5" v-if="opt.scene || opt.gender || opt.age">
                  <span v-if="opt.scene">{{ opt.scene }}</span>
                  <span v-if="opt.gender"> · {{ opt.gender }}</span>
                  <span v-if="opt.age"> · {{ opt.age }}</span>
                </div>
              </div>
              <button v-if="opt.previewUrl" type="button" class="shrink-0 size-7 rounded-full flex items-center justify-center hover:scale-110 transition shadow-sm" 
                :class="getGenderButtonClass(opt.gender)" @click.stop="togglePreview(opt)" :title="'试听 ' + opt.label">
                <svg v-if="playingValue !== opt.value" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="size-3.5 text-white ml-0.5" fill="currentColor">
                  <path d="M5 3.879v16.242a1 1 0 0 0 1.555.832l12.121-8.12a1 1 0 0 0 0-1.664L6.555 3.047A1 1 0 0 0 5 3.879Z"/>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="size-3.5 text-white" fill="currentColor">
                  <path d="M6 4.5A1.5 1.5 0 0 1 7.5 3h1A1.5 1.5 0 0 1 10 4.5v15A1.5 1.5 0 0 1 8.5 21h-1A1.5 1.5 0 0 1 6 19.5zM14 4.5A1.5 1.5 0 0 1 15.5 3h1A1.5 1.5 0 0 1 18 4.5v15A1.5 1.5 0 0 1 16.5 21h-1A1.5 1.5 0 0 1 14 19.5z"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
        <div v-if="!filteredOptions.length" class="px-3 py-6 text-center text-xs text-gray-400">无匹配结果</div>
      </div>
      
    </div>
  </div>

  <!-- Dropdown mode (default) -->
  <div v-else class="relative inline-block" ref="root">
    <button type="button" class="px-2 py-2 rounded-md text-sm bg-transparent border-0 focus:outline-none focus:ring-1 focus:ring-gray-300 hover:bg-gray-50 inline-flex items-center" @click="toggle">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-4 h-4 mr-1 text-gray-700" fill="currentColor"><path d="M7 4a1 1 0 0 1 1 1v14a1 1 0 1 1-2 0V5a1 1 0 0 1 1-1zm10 0a1 1 0 0 1 1 1v14a1 1 0 1 1-2 0V5a1 1 0 0 1 1-1zM12 7a1 1 0 0 1 1 1v8a1 1 0 1 1-2 0V8a1 1 0 0 1 1-1z"/></svg>
      {{ currentLabel }}
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="inline w-4 h-4 ml-1 text-gray-500">
        <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.24 4.5a.75.75 0 01-1.08 0l-4.24-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
      </svg>
    </button>
    <div v-show="open" class="absolute left-0 mt-1 w-[34rem] rounded-lg border border-gray-200 bg-white shadow-lg z-50 overflow-hidden">
      <div class="p-3 border-b border-gray-50">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-sm font-medium text-gray-900">选择音色</h3>
          <button class="size-6 rounded-full hover:bg-gray-100 flex items-center justify-center text-gray-400 hover:text-gray-600" @click="open=false">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="size-4" fill="currentColor">
              <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
            </svg>
          </button>
        </div>
        <div class="relative mb-3">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="size-4 text-gray-400" fill="currentColor">
              <path fill-rule="evenodd" d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z" clip-rule="evenodd" />
            </svg>
          </div>
          <input v-model="keyword" class="w-full pl-10 pr-4 py-2 text-sm border border-gray-200 rounded-lg outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500" placeholder="搜索音色名称、场景或特征..." />
          <button v-if="keyword" class="absolute inset-y-0 right-0 pr-3 flex items-center" @click="keyword=''">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" class="size-4 text-gray-400 hover:text-gray-600" fill="currentColor">
              <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
            </svg>
          </button>
        </div>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <select v-model="selectedScene" class="px-2 py-1 text-xs border border-gray-200 rounded-md bg-white text-gray-700 focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none">
              <option value="全部">场景</option>
              <option v-for="s in sceneOptions" :key="s" :value="s">{{ s }}</option>
            </select>
            <select v-model="genderFilter" class="px-2 py-1 text-xs border border-gray-200 rounded-md bg-white text-gray-700 focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none">
              <option value="全部">性别</option>
              <option v-for="g in genderOptions" :key="g" :value="g">{{ g }}</option>
            </select>
            <select v-model="ageFilter" class="px-2 py-1 text-xs border border-gray-200 rounded-md bg-white text-gray-700 focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none">
              <option value="全部">年龄</option>
              <option v-for="a in ageOptions" :key="a" :value="a">{{ a }}</option>
            </select>
            <button v-if="hasActiveFilters" class="px-2 py-1 text-xs text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-md transition" @click="resetFilters">
              清空筛选
            </button>
          </div>
          <div class="text-xs text-gray-400">共 {{ filteredOptions.length }} 项</div>
        </div>
      </div>
      <div class="max-h-80 overflow-auto p-2">
        <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
          <div v-for="opt in filteredOptions" :key="opt.value" class="relative border rounded-[calc(var(--radius,8px)+4px)] border-gray-200 p-3 hover:bg-gray-50 transition cursor-pointer" :class="selectedCardClass(opt)" @click="choose(opt)">
            <div v-if="localValue===opt.value" class="absolute -top-1 -left-1 size-4 rounded-full flex items-center justify-center ring-2 ring-white z-10" :class="selectedBadgeClass(opt)">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="size-3 text-white" fill="currentColor">
                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
              </svg>
            </div>
            <div class="flex items-start justify-between gap-2 pt-0.5 pl-0.5">
              <div class="min-w-0 flex-1">
                <div class="text-sm truncate flex items-center gap-1" :class="{'font-medium': selectedValue===opt.value}">
                  <span class="truncate">{{ opt.label }}</span>
                  <span v-if="opt.tags && opt.tags.includes('上新')" class="shrink-0 inline-flex items-center px-1 text-[10px] rounded bg-pink-100 text-pink-600">新</span>
                </div>
                <div class="text-[10px] text-gray-400 truncate mt-0.5" v-if="opt.scene || opt.gender || opt.age">
                  <span v-if="opt.scene">{{ opt.scene }}</span>
                  <span v-if="opt.gender"> · {{ opt.gender }}</span>
                  <span v-if="opt.age"> · {{ opt.age }}</span>
                </div>
              </div>
              <button v-if="opt.previewUrl" type="button" class="shrink-0 size-7 rounded-full flex items-center justify-center hover:scale-110 transition shadow-sm" 
                :class="getGenderButtonClass(opt.gender)" @click.stop="togglePreview(opt)" :title="'试听 ' + opt.label">
                <svg v-if="playingValue !== opt.value" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="size-3.5 text-white ml-0.5" fill="currentColor">
                  <path d="M5 3.879v16.242a1 1 0 0 0 1.555.832l12.121-8.12a1 1 0 0 0 0-1.664L6.555 3.047A1 1 0 0 0 5 3.879Z"/>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="size-3.5 text-white" fill="currentColor">
                  <path d="M6 4.5A1.5 1.5 0 0 1 7.5 3h1A1.5 1.5 0 0 1 10 4.5v15A1.5 1.5 0 0 1 8.5 21h-1A1.5 1.5 0 0 1 6 19.5zM14 4.5A1.5 1.5 0 0 1 15.5 3h1A1.5 1.5 0 0 1 18 4.5v15A1.5 1.5 0 0 1 16.5 21h-1A1.5 1.5 0 0 1 14 19.5z"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
        <div v-if="!filteredOptions.length" class="px-3 py-6 text-center text-xs text-gray-400">无匹配结果</div>
      </div>
      
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import AudioPlayer from './AudioPlayer.vue'

const defaultOptions = [
  { label: '温柔女神', value: 'ICL_zh_female_wenrounvshen_239eff5e8ffa_tob', lang: '中文', scene: '通用', gender: '女', emotions: [] },
  { label: '阳光青年', value: 'zh_male_yangguangqingnian_moon_bigtts', lang: '中文', scene: '通用', gender: '男', emotions: [] },
  { label: '清新女声', value: 'zh_female_qingxinnvsheng_mars_bigtts', lang: '中文', scene: '通用', gender: '女', emotions: [] },
  { label: '深夜播客', value: 'zh_male_shenyeboke_moon_bigtts', lang: '中文', scene: '播客', gender: '男', emotions: [] },
  { label: 'Vivi', value: 'zh_female_vv_mars_bigtts', lang: '中文', scene: '通用', gender: '女', emotions: [] },
  { label: 'Candice(美)', value: 'en_female_candice_emo_v2_mars_bigtts', lang: '英语', scene: '通用', gender: '女', emotions: ['深情','愤怒','ASMR','对话','兴奋','愉悦','中性','温暖'] },
  { label: 'Corey(英)', value: 'en_male_corey_emo_v2_mars_bigtts', lang: '英语', scene: '通用', gender: '男', emotions: ['权威'] },
  { label: 'Serena(美)', value: 'en_female_skye_emo_v2_mars_bigtts', lang: '英语', scene: '通用', gender: '女', emotions: [] },
  { label: 'Glen(美)', value: 'en_male_glen_emo_v2_mars_bigtts', lang: '英语', scene: '通用', gender: '男', emotions: [] },
]

export default defineComponent({
  name: 'VoiceSelector',
  components: { AudioPlayer },
  props: {
    modelValue: { type: String, default: '' },
    options: { type: Array, default: () => defaultOptions },
    inline: { type: Boolean, default: false },
    emitOnMount: { type: Boolean, default: false },
    // 允许从外部注入一些置顶的选项（例如“我的声音”）
    prependOptions: { type: Array, default: () => [] },
  },
  emits: ['update:modelValue','selected'],
  data() {
    return {
      open: false,
      keyword: '',
      localValue: this.modelValue || (this.options[0] && this.options[0].value),
      allOptions: ([]).concat(this.prependOptions || [], this.options || []),
      selectedScene: '全部',
      genderFilter: '全部',
      ageFilter: '全部',
      playingValue: '',
      previewingValue: '',
      audio: null,
    }
  },
  computed: {
    currentLabel() {
      const found = (this.allOptions || []).find(o => o.value === this.localValue)
      return found ? found.label : '选择音色'
    },
    currentGenderClass() {
      const found = (this.allOptions || []).find(o => o.value === this.localValue)
      const g = found && found.gender ? String(found.gender).toLowerCase() : ''
      if (g.includes('女') || g.includes('female')) return 'text-pink-600'
      if (g.includes('男') || g.includes('male')) return 'text-blue-600'
      return 'text-gray-700'
    },
    selectedOption() {
      const key = this.previewingValue || this.localValue
      return (this.allOptions || []).find(o => o.value === key) || null
    },
    selectedValue() {
      return this.previewingValue || this.localValue
    },
    previewText() {
      return this.selectedOption && this.selectedOption.demoText || ''
    },
    isPlayingCurrent() {
      return !!(this.playingValue && this.selectedOption && this.playingValue === this.selectedOption.value)
    },
    hasActiveFilters() {
      return this.keyword || this.selectedScene !== '全部' || this.genderFilter !== '全部' || this.ageFilter !== '全部'
    },
    sceneOptions() {
      const set = new Set()
      const list = this.getFacetSource('scene')
      list.forEach(o => { if (o.scene) set.add(o.scene) })
      return Array.from(set)
    },
    genderOptions() {
      const set = new Set()
      const list = this.getFacetSource('gender')
      list.forEach(o => { if (o.gender) set.add(o.gender) })
      return Array.from(set)
    },
    ageOptions() {
      const set = new Set()
      const list = this.getFacetSource('age')
      list.forEach(o => { if (o.age) set.add(o.age) })
      return Array.from(set)
    },
    filteredOptions() {
      let list = this.allOptions
      // Scene filter
      if (this.selectedScene && this.selectedScene !== '全部') list = list.filter(o => (o.scene || '') === this.selectedScene)
      // Gender / Age
      if (this.genderFilter && this.genderFilter !== '全部') list = list.filter(o => (o.gender || '') === this.genderFilter)
      if (this.ageFilter && this.ageFilter !== '全部') list = list.filter(o => (o.age || '') === this.ageFilter)

      // Keyword
      if (this.keyword) {
        const k = this.keyword.toLowerCase()
        list = list.filter(o => (
          (o.label && o.label.toLowerCase().includes(k)) ||
          (o.value && String(o.value).toLowerCase().includes(k)) ||
          (o.scene && String(o.scene).toLowerCase().includes(k)) ||
          (o.gender && String(o.gender).toLowerCase().includes(k)) ||
          (o.age && String(o.age).toLowerCase().includes(k))
        ))
      }
      return list
    },
    
  },
  watch: {
    modelValue(v) {
      if (v !== this.localValue) this.localValue = v
    },
    
  },
  mounted() {
    if (!this.inline) document.addEventListener('click', this.onOutside, true)
    this.tryLoadRemote()
    if (this.emitOnMount) this.emitSelectedForCurrent()
  },
  beforeUnmount() {
    if (!this.inline) document.removeEventListener('click', this.onOutside, true)
    try { if (this.audio) { this.audio.pause(); this.audio.src = '' } } catch (e) {}
  },
  methods: {
    getFacetSource(excludeKey) {
      let list = this.allOptions
      // 关键词过滤始终参与（更贴近“真实可见内容”）
      if (this.keyword) {
        const k = this.keyword.toLowerCase()
        list = list.filter(o => (
          (o.label && o.label.toLowerCase().includes(k)) ||
          (o.value && String(o.value).toLowerCase().includes(k)) ||
          (o.scene && String(o.scene).toLowerCase().includes(k)) ||
          (o.gender && String(o.gender).toLowerCase().includes(k)) ||
          (o.age && String(o.age).toLowerCase().includes(k))
        ))
      }
      if (excludeKey !== 'scene' && this.selectedScene && this.selectedScene !== '全部') {
        list = list.filter(o => (o.scene || '') === this.selectedScene)
      }
      if (excludeKey !== 'gender' && this.genderFilter && this.genderFilter !== '全部') {
        list = list.filter(o => (o.gender || '') === this.genderFilter)
      }
      if (excludeKey !== 'age' && this.ageFilter && this.ageFilter !== '全部') {
        list = list.filter(o => (o.age || '') === this.ageFilter)
      }
      return list
    },
    emitSelectedForCurrent() {
      const found = (this.allOptions || []).find(o => o.value === this.localValue)
      if (found) this.$emit('selected', { value: found.value, label: found.label, gender: found.gender })
    },
    selectedBadgeClass(opt) {
      const g = (opt && opt.gender ? String(opt.gender).toLowerCase() : '')
      if (g.includes('女') || g.includes('female')) return 'bg-pink-500'
      if (g.includes('男') || g.includes('male')) return 'bg-blue-500'
      return 'bg-gray-500'
    },
    toggle() { 
      if (this.inline) return
      this.open = !this.open 
      // 打开时如果还没有预览项，默认选中当前值对应的音色
      if (this.open && !this.previewingValue && this.localValue) {
        this.previewingValue = this.localValue
      }
    },
    switchScene(s) { this.selectedScene = s },
    onOutside(e) { if (this.inline) return; const root = this.$refs.root; if (this.open && root && !root.contains(e.target)) this.open = false },
    select(opt) { this.previewingValue = (opt && opt.value) || '' },
    preview(opt) { this.previewingValue = (opt && opt.value) || '' },
    choose(opt) {
      if (!opt) return
      // 点击卡片即选中并提交
      this.localValue = opt.value
      this.previewingValue = ''
      this.$emit('update:modelValue', this.localValue)
      this.$emit('selected', { value: opt.value, label: opt.label, gender: opt.gender })
      // 关闭下拉模式
      if (!this.inline) this.open = false
    },
    selectedCardClass(opt) {
      const isSelected = this.localValue === (opt && opt.value)
      if (!isSelected) return ''
      const g = (opt && opt.gender ? String(opt.gender).toLowerCase() : '')
      if (g.includes('女') || g.includes('female')) return 'ring-2 ring-pink-500 bg-pink-50 border-pink-300'
      if (g.includes('男') || g.includes('male')) return 'ring-2 ring-blue-500 bg-blue-50 border-blue-300'
      return 'ring-2 ring-gray-500 bg-gray-50 border-gray-300'
    },
    resetFilters() { this.keyword=''; this.selectedScene='全部'; this.genderFilter='全部'; this.ageFilter='全部' },
    togglePreview(opt) {
      if (!opt || !opt.previewUrl) return
      if (!this.audio) this.audio = new Audio()
      this.audio.crossOrigin = 'anonymous'
      this.audio.preload = 'auto'
      if (this.playingValue === opt.value) {
        this.audio.pause(); this.playingValue = ''; return
      }
      this.playingValue = opt.value
      this.audio.src = opt.previewUrl
      try { this.audio.currentTime = 0 } catch (e) {}
      try { this.audio.load() } catch (e) {}
      this.audio.onended = () => { this.playingValue = '' }
      this.audio.play().catch(() => { this.playingValue = '' })
    },
    togglePreviewCurrent() {
      const opt = this.selectedOption
      if (!opt) return
      this.togglePreview(opt)
    },
    confirmUse() {
      if (!this.previewingValue) return
      this.localValue = this.previewingValue
      this.$emit('update:modelValue', this.localValue)
      this.open = false
    },
    getGenderButtonClass(gender) {
      const g = (gender || '').toLowerCase()
      if (g.includes('女') || g.includes('female')) return 'bg-pink-500 hover:bg-pink-600'
      if (g.includes('男') || g.includes('male')) return 'bg-blue-500 hover:bg-blue-600'
      return 'bg-gray-500 hover:bg-gray-600'
    },
    playAndSelect(opt) {},
    async tryLoadRemote() {
      // 支持从 /voices.normalized.json 全量导入（放在 public 目录）
      try {
        const res = await fetch('/voices.normalized.json', { cache: 'no-store' })
        if (!res.ok) return
        const json = await res.json()
        const normalized = this.normalizeVoices(json)
        if (normalized.length) {
          // 远端数据前拼接外部注入项
          this.allOptions = ([]).concat(this.prependOptions || [], normalized)
          if (!this.allOptions.find(o => o.value === this.localValue)) {
            this.localValue = this.allOptions[0].value
            this.$emit('update:modelValue', this.localValue)
          }
          if (!this.previewingValue && this.localValue) {
            this.previewingValue = this.localValue
          }
          if (this.emitOnMount) this.emitSelectedForCurrent()
        }
      } catch (e) { /* 静默失败，使用默认集合 */ }
    },
    normalizeVoices(json) {
      // 适配两种：1) 直接数组；2) 按集分组后的对象（含 voice_options）
      const out = []
      const seenValues = new Set()
      const coerceUrl = (u) => (u && String(u).trim().length ? String(u).trim() : '')
      const pickPreferredMode = (item, links) => {
        // 优先 speak，其次 item.mode 中声明的第一个，否则 links 的第一个 key
        if (links && typeof links === 'object') {
          if (links.speak || links.Speak) return links.speak ? 'speak' : 'Speak'
          if (item && Array.isArray(item.mode) && item.mode.length) return item.mode[0]
          const keys = Object.keys(links)
          return keys && keys.length ? keys[0] : ''
        }
        return (item && Array.isArray(item.mode) && item.mode[0]) || 'speak'
      }
      const pickPreviewFromLinks = (item) => {
        const links = item && item.audio_links
        if (!links || typeof links !== 'object') return ''
        const modeKey = pickPreferredMode(item, links)
        const candidate = links[modeKey] || links.default || Object.values(links)[0]
        if (!candidate) return ''
        if (typeof candidate === 'string') return coerceUrl(candidate)
        return coerceUrl(candidate.url || candidate.target || candidate.source)
      }
      const pushFrom = (item) => {
        if (!item) return
        const name = item.name || item.label
        const scene = (item.category && item.category[0] && item.category[0].level1) || item.scene || ''
        const gender = item.gender || item.sex || ''
        const age = item.age || ''
        const voiceType = item.value || item.voice_type || (item.voice_config && item.voice_config[0] && item.voice_config[0].params && item.voice_config[0].params.voice_type) || ''
        let previewUrl = coerceUrl(item.previewUrl || item.trial_url || item.audio_url || item.audio)
        if (!previewUrl) previewUrl = pickPreviewFromLinks(item)
        const demoText = (Array.isArray(item.voice_config) && item.voice_config[0] && item.voice_config[0].text) || item.text || ''
        const tags = Array.isArray(item.labels) ? item.labels.slice(0) : (Array.isArray(item.tags) ? item.tags.slice(0) : [])
        if (!name || !voiceType) return
        if (seenValues.has(voiceType)) return
        seenValues.add(voiceType)
        out.push({ label: name, value: voiceType, scene, gender, age, previewUrl, tags, demoText })
      }
      if (Array.isArray(json)) {
        json.forEach(pushFrom)
      } else if (json && typeof json === 'object') {
        Object.keys(json).forEach(k => {
          const grp = json[k]
          if (grp && Array.isArray(grp.voice_options)) grp.voice_options.forEach(pushFrom)
        })
      }
      return out
    },
  },
})
</script>





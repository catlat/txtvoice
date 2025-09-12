<template>
  <main>
    <div class="px-4 py-10 sm:px-0">
      <!-- Page Header -->
      <div class="max-w-4xl mx-auto text-center mb-8">
        <div class="text-3xl font-bold tracking-tight text-gray-800">账号管理</div>
        <div class="mt-3 text-base text-gray-500">管理您的账号信息和套餐</div>
      </div>

      <!-- Main Content -->
      <div class="max-w-4xl mx-auto space-y-6">
        <!-- Login Section (if not logged in) -->
        <LoginRequired v-if="!token" title="未登录" description="请登录后查看账号信息与套餐详情" />

        <!-- Account Info (if logged in) -->
        <div v-else class="space-y-6">
          <!-- User Overview -->
          <div class="bg-white/90 backdrop-blur rounded-xl border border-gray-200 shadow-sm p-6">
            <div class="flex items-center justify-between mb-6">
              <h3 class="text-xl font-semibold text-gray-800">账号概览</h3>
              <button 
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors duration-200" 
                @click="doLogout"
              >
                退出登录
              </button>
            </div>
            
            <!-- User Profile -->
            <div class="flex items-center gap-4 mb-6">
              <div class="w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center border border-gray-200">
                <span class="text-gray-700 text-xl font-medium">{{ identity.charAt(0).toUpperCase() }}</span>
              </div>
              <div>
                <div class="text-lg font-semibold text-gray-800">{{ identity }}</div>
                <div class="text-sm text-gray-500">用户标识</div>
                <div v-if="profileData && profileData.created_at" class="text-xs text-gray-400 mt-1">
                  注册于 {{ formatDate(profileData.created_at) }}
                </div>
              </div>
              <div class="ml-auto">
                <button
                  class="px-3 py-1.5 text-sm border border-gray-300 rounded-md text-gray-700 hover:bg-gray-100 transition-colors"
                  @click="togglePasswordForm"
                >
                  {{ showPasswordForm ? '收起' : '修改密码' }}
                </button>
              </div>
            </div>

            <!-- 修改密码（折叠区域） -->
            <div v-if="showPasswordForm" class="mt-2 max-w-md border border-gray-200 rounded-lg p-4 bg-white" autocomplete="off">
              <div class="space-y-3">
                <label class="block">
                  <span class="block mb-1 text-sm text-gray-600">新密码</span>
                  <input
                    v-model="newPassword"
                    :readonly="antiFill"
                    @focus="antiFill = false"
                    type="password"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md bg-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-300 focus:border-transparent"
                    placeholder="至少 6 位"
                    name="new-password"
                    autocomplete="new-password"
                    inputmode="text"
                    autocapitalize="off"
                    autocorrect="off"
                    spellcheck="false"
                  />
                </label>
                <label class="block">
                  <span class="block mb-1 text-sm text-gray-600">确认新密码</span>
                  <input
                    v-model="confirmPassword"
                    :readonly="antiFill"
                    @focus="antiFill = false"
                    type="password"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md bg-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-gray-300 focus:border-transparent"
                    placeholder="再次输入新密码"
                    name="confirm-password"
                    autocomplete="new-password"
                    inputmode="text"
                    autocapitalize="off"
                    autocorrect="off"
                    spellcheck="false"
                  />
                </label>
                <div v-if="passwordMismatch" class="text-xs text-red-600">两次输入不一致</div>
                <div class="flex items-center gap-3 pt-1">
                  <button
                    class="px-4 py-2 rounded-md bg-gray-200 text-gray-800 hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
                    :disabled="!canSubmit || isSubmitting"
                    @click="doChangePassword"
                  >保存</button>
                  <button
                    class="px-4 py-2 rounded-md border border-gray-300 text-gray-700 hover:bg-gray-100"
                    @click="cancelPasswordChange"
                  >取消</button>
                </div>
              </div>
            </div>

            
            
            <!-- Package Summary -->
            <div v-if="packages.length" class="border-t pt-6 mt-6">
              <h4 class="text-lg font-medium text-gray-800 mb-4">我的套餐</h4>
              <div class="space-y-6">
                <div v-for="(p, idx) in packages" :key="idx" class="bg-white/80 backdrop-blur rounded-xl border border-gray-200 p-6 shadow-sm">
                  <!-- 套餐头部 -->
                  <div class="flex items-start justify-between mb-6">
                    <div>
                      <h5 class="text-lg font-semibold text-gray-800">{{ p.package_name || ('套餐 #' + (p.package_id || idx + 1)) }}</h5>
                      <div class="text-sm text-gray-600 mt-1 flex items-center gap-2">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                          <path stroke-linecap="round" stroke-linejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5a2.25 2.25 0 0 1 2.25-2.25m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5a2.25 2.25 0 0 1 2.25-2.25V18.75M9 10.5h6" />
                        </svg>
                        到期: {{ formatDate(p.expire_at) }}
                      </div>
                    </div>
                    <div class="text-xs bg-gray-100 text-gray-700 px-3 py-1.5 rounded-full border border-gray-200">
                      有效
                    </div>
                  </div>
                  
                  <!-- 语音识别使用情况 -->
                  <div class="mb-6">
                    <div class="flex items-center justify-between mb-2">
                      <div class="flex items-center gap-2">
                        <div class="w-3 h-3 rounded-full bg-sky-400"></div>
                        <span class="text-sm font-medium text-gray-700">语音识别</span>
                      </div>
                      <div class="text-sm text-gray-600">
                        {{ formatNumber(getUsedASR(p)) }} / {{ formatNumber(getTotalASR(p)) }}
                      </div>
                    </div>
                    <div class="w-full bg-gray-100 rounded-full h-2.5">
                      <div 
                        class="bg-gray-400 h-2.5 rounded-full transition-all duration-500"
                        :style="{ width: getUsagePercentage('asr', p) + '%' }"
                      ></div>
                    </div>
                    <div class="mt-1 flex justify-between text-xs text-gray-500">
                      <span>余额: {{ formatNumber(computeASR(p).remain) }}</span>
                      <span>{{ getUsagePercentage('asr', p) }}%</span>
                    </div>
                  </div>
                  
                  <!-- 语音合成使用情况 -->
                  <div class="mb-2">
                    <div class="flex items-center justify-between mb-2">
                      <div class="flex items-center gap-2">
                        <div class="w-3 h-3 rounded-full bg-indigo-400"></div>
                        <span class="text-sm font-medium text-gray-700">语音合成</span>
                      </div>
                      <div class="text-sm text-gray-600">
                        {{ formatNumber(getUsedTTS(p)) }} / {{ formatNumber(getTotalTTS(p)) }}
                      </div>
                    </div>
                    <div class="w-full bg-gray-100 rounded-full h-2.5">
                      <div 
                        class="bg-gray-400 h-2.5 rounded-full transition-all duration-500"
                        :style="{ width: getUsagePercentage('tts', p) + '%' }"
                      ></div>
                    </div>
                    <div class="mt-1 flex justify-between text-xs text-gray-500">
                      <span>余额: {{ formatNumber(computeTTS(p).remain) }}</span>
                      <span>{{ getUsagePercentage('tts', p) }}%</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- No Package State -->
            <div v-else class="border-t pt-6 text-center py-8">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 mx-auto mb-3 text-gray-300">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-2.25-1.313M21 7.5v2.25m0-2.25l-2.25 1.313M3 7.5l2.25-1.313M3 7.5l2.25 1.313M3 7.5v2.25m9 3l2.25-1.313L15 12.75l-2.25-1.313M9 12.75l-2.25-1.313L9 10.5l2.25 1.313m0 0L11.25 9.75l2.25 1.313M15 12.75l-2.25-1.313M15 12.75V15m0 0l2.25 1.313M15 15H12.75m3.75 0l-3.75-3.75" />
              </svg>
              <div class="text-gray-500 mb-2">暂无有效套餐</div>
              <div class="text-sm text-gray-400">请联系管理员购买服务套餐</div>
            </div>
          </div>

          <!-- Usage Statistics -->
          <div class="bg-white/90 backdrop-blur rounded-xl border border-gray-200 shadow-sm p-6">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h3 class="text-xl font-semibold text-gray-800">使用统计</h3>
                <div class="text-sm text-gray-500 mt-1">近30天使用概览</div>
              </div>
              <button @click="refresh" class="text-sm text-gray-600 hover:text-gray-700 flex items-center gap-1">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0 3.181 3.183a8.25 8.25 0 0 0 13.803-3.7M4.031 9.865a8.25 8.25 0 0 1 13.803-3.7l3.181 3.182m0-4.991v4.99" />
                </svg>
                刷新
              </button>
            </div>
            
            <!-- 简化的统计总览 -->
            <div v-if="usageDays.length" class="space-y-4">
              <!-- 总计统计 -->
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div class="bg-white border border-gray-200 rounded-lg p-4 text-center">
                  <div class="text-2xl font-bold text-gray-800">{{ formatNumber(totalAsrChars) }}</div>
                  <div class="text-gray-500 text-sm">语音识别总计</div>
                </div>
                <div class="bg-white border border-gray-200 rounded-lg p-4 text-center">
                  <div class="text-2xl font-bold text-gray-800">{{ formatNumber(totalTtsChars) }}</div>
                  <div class="text-gray-500 text-sm">语音合成总计</div>
                </div>
                <div class="bg-white border border-gray-200 rounded-lg p-4 text-center">
                  <div class="text-2xl font-bold text-gray-800">{{ totalRequests }}</div>
                  <div class="text-gray-500 text-sm">请求总数</div>
                </div>
              </div>
              
              <!-- 活跃天数概览 -->
              <div v-if="activeDays.length > 0" class="bg-gray-50 rounded-lg p-4">
                <div class="flex items-center justify-between mb-3">
                  <h4 class="font-medium text-gray-800">活跃概览</h4>
                  <span class="text-sm text-gray-500">{{ activeDays.length }} 个活跃天</span>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 text-sm">
                  <div v-for="day in activeDays.slice(0, 6)" :key="day.date" class="bg-white/80 rounded-lg p-3 border border-gray-100">
                    <div class="font-medium text-gray-800">{{ formatDate(day.date) }}</div>
                    <div class="flex gap-2 text-xs text-gray-600 mt-1">
                      <span v-if="day.asr_chars > 0" class="bg-gray-100 text-gray-700 px-2 py-1 rounded border border-gray-200">识别{{ formatNumber(day.asr_chars) }}</span>
                      <span v-if="day.tts_chars > 0" class="bg-gray-100 text-gray-700 px-2 py-1 rounded border border-gray-200">合成{{ formatNumber(day.tts_chars) }}</span>
                    </div>
                  </div>
                </div>
                <div v-if="activeDays.length > 6" class="mt-3 text-center">
                  <span class="text-xs text-gray-500">还有 {{ activeDays.length - 6 }} 天有使用记录</span>
                </div>
              </div>
            </div>
            
            <!-- 无数据状态 -->
            <div v-else class="text-center py-8">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 mx-auto mb-3 text-gray-300">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" />
              </svg>
              <div class="text-gray-500 mb-2">暂无使用记录</div>
              <div class="text-sm text-gray-400">开始使用语音服务后，这里将显示您的使用统计</div>
              <div v-if="error" class="text-red-500 text-sm mt-2">{{ error }}</div>
            </div>
          </div>
        </div>

        <!-- Error Display -->
        <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
          <div class="flex items-center gap-2 text-red-700">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z" />
            </svg>
            {{ error }}
          </div>
        </div>
      </div>
    </div>
  </main>
</template>

<script>
import { defineComponent } from 'vue'
import * as account from '../api/account'
import { getToken, clearToken, getIdentity, clearIdentity } from '../utils/auth'
import LoginRequired from '../components/LoginRequired.vue'

export default defineComponent({
  components: { LoginRequired },
  data() {
    return {
      identity: getIdentity() || '',
      
      newPassword: '',
      confirmPassword: '',
      showPasswordForm: false,
      isSubmitting: false,
      antiFill: true,
      token: getToken() || '',
      profileData: null,
      packages: [],
      usageDays: [],
      error: '',
    }
  },
  computed: {
    // 计算总计数据
    totalAsrChars() {
      return this.usageDays.reduce((sum, day) => sum + (day.asr_chars || 0), 0)
    },
    totalTtsChars() {
      return this.usageDays.reduce((sum, day) => sum + (day.tts_chars || 0), 0)
    },
    totalRequests() {
      return this.usageDays.reduce((sum, day) => sum + (day.requests || 0), 0)
    },
    // 有使用记录的天数
    activeDays() {
      return this.usageDays.filter(day => 
        (day.asr_chars || 0) > 0 || 
        (day.tts_chars || 0) > 0 || 
        (day.requests || 0) > 0
      )
    },
    passwordMismatch() {
      if (!this.showPasswordForm) return false
      if (!this.confirmPassword) return false
      return this.newPassword !== this.confirmPassword
    },
    canSubmit() {
      return (
        this.newPassword &&
        this.confirmPassword &&
        this.newPassword.length >= 6 &&
        !this.passwordMismatch
      )
    }
  },
  created() {
    if (this.identity && this.token) this.refresh()
  },
  methods: {
    togglePasswordForm() {
      this.showPasswordForm = !this.showPasswordForm
      if (!this.showPasswordForm) {
        this.newPassword = ''
        this.confirmPassword = ''
        this.antiFill = true
      }
    },
    cancelPasswordChange() {
      this.showPasswordForm = false
      this.newPassword = ''
      this.confirmPassword = ''
      this.antiFill = true
    },
    
    async doLogout() {
      try { await account.logout(this.token); const { toast } = require('../utils/toast'); toast('已退出', 'success') } catch (e) {}
      clearToken(); clearIdentity(); this.token = ''; this.identity = ''
      this.profileData = null; this.packages = []; this.usageDays = []
    },
    async doChangePassword() {
      this.error = ''
      if (!this.canSubmit) return
      this.isSubmitting = true
      try {
        await account.changePassword(this.newPassword)
        this.newPassword = ''
        this.confirmPassword = ''
        this.showPasswordForm = false
        this.antiFill = true
        try { const { toast } = require('../utils/toast'); toast('密码已更新', 'success') } catch (e) {}
      } catch (e) {
        this.error = e && e.message ? e.message : '修改密码失败'
        try { const { toast } = require('../utils/toast'); toast(this.error, 'error') } catch (e2) {}
      } finally { this.isSubmitting = false }
    },
    async refresh() {
      try {
        const [prof, packs, usage] = await Promise.all([
          account.profile(), account.packages(), account.usage({})
        ])
        
        this.profileData = prof || null
        this.packages = this.extractPackagesData(packs)
        this.usageDays = this.extractUsageData(usage)
        
        // 调试日志
        console.log('Account API responses:', { prof, packs, usage })
        console.log('Extracted usage days:', this.usageDays.length)
        
      } catch (e) { 
        console.error('Account refresh error:', e)
        this.error = e && e.message ? e.message : '加载失败'
        try { 
          const { toast } = require('../utils/toast')
          toast(this.error, 'error') 
        } catch (e) {}
      }
    },
    extractPackagesData(packs) {
      // 处理套餐数据的多种可能结构
      if (Array.isArray(packs)) return packs
      if (packs && packs.data && Array.isArray(packs.data.items)) return packs.data.items
      if (packs && Array.isArray(packs.items)) return packs.items
      if (packs && Array.isArray(packs.data)) return packs.data
      return []
    },
    extractUsageData(usage) {
      // 修复数据提取逻辑 - 正确处理API返回结构
      let days = []
      
      if (usage && usage.data && Array.isArray(usage.data.days)) {
        days = usage.data.days
      } else if (usage && Array.isArray(usage.days)) {
        days = usage.days
      } else if (Array.isArray(usage)) {
        days = usage
      }
      
      console.log('Raw usage data:', usage)
      console.log('Extracted days:', days)
      
      return days
    },
    // 判断某天是否有使用记录
    hasUsage(day) {
      return (day.asr_chars || 0) > 0 || (day.tts_chars || 0) > 0 || (day.requests || 0) > 0
    },
    // 格式化数字显示
    formatNumber(num) {
      if (!num || num === 0) return '0'
      if (num >= 1000000) {
        return (num / 1000000).toFixed(1) + 'M'
      }
      if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'K'
      }
      return num.toString()
    },
    // 规范化 ASR 配额数据（总量/已用/余额）
    computeASR(pkg) {
      let total = Number(pkg.quota_asr_chars) || 0
      let remain = Number(pkg.remain_asr_chars)
      remain = Number.isFinite(remain) ? Math.max(0, remain) : 0
      let used = Number(pkg.used_asr_chars)
      used = Number.isFinite(used) ? Math.max(0, used) : NaN
      if (!total) {
        const baseUsed = Number.isFinite(used) ? used : 0
        if (baseUsed > 0 || remain > 0) total = baseUsed + remain
      }
      if (!Number.isFinite(used)) used = Math.max(total - remain, 0)
      // clamp
      remain = Math.min(Math.max(remain, 0), total)
      used = Math.min(Math.max(used, 0), total)
      return { total, used, remain }
    },
    // 规范化 TTS 配额数据（总量/已用/余额）
    computeTTS(pkg) {
      let total = Number(pkg.quota_tts_chars) || 0
      let remain = Number(pkg.remain_tts_chars)
      remain = Number.isFinite(remain) ? Math.max(0, remain) : 0
      let used = Number(pkg.used_tts_chars)
      used = Number.isFinite(used) ? Math.max(0, used) : NaN
      if (!total) {
        const baseUsed = Number.isFinite(used) ? used : 0
        if (baseUsed > 0 || remain > 0) total = baseUsed + remain
      }
      if (!Number.isFinite(used)) used = Math.max(total - remain, 0)
      // clamp
      remain = Math.min(Math.max(remain, 0), total)
      used = Math.min(Math.max(used, 0), total)
      return { total, used, remain }
    },
    // 套餐使用情况计算（仅分开显示）
    getTotalASR(pkg) { return this.computeASR(pkg).total },
    getTotalTTS(pkg) { return this.computeTTS(pkg).total },
    getUsedASR(pkg) { return this.computeASR(pkg).used },
    getUsedTTS(pkg) { return this.computeTTS(pkg).used },
    getUsagePercentage(type, pkg) {
      const info = type === 'asr' ? this.computeASR(pkg) : this.computeTTS(pkg)
      if (!info.total) return 0
      const pct = Math.round((info.used / info.total) * 100)
      return Math.min(Math.max(pct, 0), 100)
    },
    formatDate(dateStr) {
      if (!dateStr) return '未知时间'
      try {
        const date = new Date(dateStr)
        return date.toLocaleDateString('zh-CN', {
          year: 'numeric',
          month: '2-digit',
          day: '2-digit'
        })
      } catch (e) {
        return dateStr
      }
    },
  },
})
</script>




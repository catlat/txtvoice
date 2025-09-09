<template>
  <main>
    <div class="px-4 py-10 sm:px-0">
      <!-- Page Header -->
      <div class="max-w-4xl mx-auto text-center mb-8">
        <div class="text-3xl font-bold tracking-tight text-gray-900">账号管理</div>
        <div class="mt-3 text-lg text-gray-500">管理您的账号信息和套餐</div>
      </div>

      <!-- Main Content -->
      <div class="max-w-4xl mx-auto space-y-6">
        <!-- Login Section (if not logged in) -->
        <div v-if="!token" class="bg-white rounded-xl border border-gray-200 shadow-[0_0_16px_0_rgba(0,0,0,0.06)] p-6">
          <div class="text-center">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 mx-auto mb-4 text-gray-400">
              <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
            </svg>
            <h3 class="text-lg font-semibold text-gray-900 mb-2">请先登录</h3>
            <p class="text-gray-500 mb-6">登录后查看您的账号信息和套餐详情</p>
            
            <div class="max-w-sm mx-auto space-y-4">
              <input 
                v-model="identity" 
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent" 
                placeholder="输入手机号或标识" 
              />
              <button 
                class="w-full px-4 py-3 bg-gradient-to-r from-purple-500 to-pink-500 text-white font-medium rounded-lg hover:from-purple-600 hover:to-pink-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200" 
                @click="doLogin" 
                :disabled="!identity.trim()"
              >
                登录
              </button>
            </div>
          </div>
        </div>

        <!-- Account Info (if logged in) -->
        <div v-else class="space-y-6">
          <!-- User Profile -->
          <div class="bg-white rounded-xl border border-gray-200 shadow-[0_0_16px_0_rgba(0,0,0,0.06)] p-6">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-xl font-semibold text-gray-900">账号信息</h3>
              <button 
                class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors duration-200" 
                @click="doLogout"
              >
                退出登录
              </button>
            </div>
            
            <div v-if="profileData" class="space-y-4">
              <div class="flex items-center gap-3">
                <div class="w-12 h-12 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center">
                  <span class="text-white text-lg font-medium">{{ identity.charAt(0).toUpperCase() }}</span>
                </div>
                <div>
                  <div class="font-medium text-gray-900">{{ identity }}</div>
                  <div class="text-sm text-gray-500">用户标识</div>
                </div>
              </div>
              
              <!-- Debug Info (collapsible) -->
              <details class="mt-4">
                <summary class="cursor-pointer text-sm text-gray-500 hover:text-gray-700">详细信息</summary>
                <pre class="mt-2 p-4 bg-gray-50 rounded-lg text-xs text-gray-600 overflow-auto">{{ JSON.stringify(profileData, null, 2) }}</pre>
              </details>
            </div>
          </div>

          <!-- Packages -->
          <div class="bg-white rounded-xl border border-gray-200 shadow-[0_0_16px_0_rgba(0,0,0,0.06)] p-6">
            <h3 class="text-xl font-semibold text-gray-900 mb-4">我的套餐</h3>
            
            <div v-if="packages.length" class="grid gap-4 md:grid-cols-2">
              <div v-for="(p, idx) in packages" :key="idx" class="border border-gray-200 rounded-lg p-4 hover:border-purple-200 hover:bg-purple-50/50 transition-all duration-200">
                <div class="flex items-start justify-between mb-3">
                  <div>
                    <h4 class="font-medium text-gray-900">{{ p.package_name || ('套餐 #' + (p.package_id || idx + 1)) }}</h4>
                    <div class="text-sm text-gray-500 mt-1">有效期至 {{ formatDate(p.expire_at) }}</div>
                  </div>
                  <div class="text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full">有效</div>
                </div>
                
                <div class="space-y-2">
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">语音识别</span>
                    <span class="text-sm font-medium text-gray-900">{{ p.remain_asr_chars || 0 }} 字符</span>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">语音合成</span>
                    <span class="text-sm font-medium text-gray-900">{{ p.remain_tts_chars || 0 }} 字符</span>
                  </div>
                </div>
              </div>
            </div>
            
            <div v-else class="text-center py-8">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 mx-auto mb-3 text-gray-300">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 7.5l-2.25-1.313M21 7.5v2.25m0-2.25l-2.25 1.313M3 7.5l2.25-1.313M3 7.5l2.25 1.313M3 7.5v2.25m9 3l2.25-1.313L15 12.75l-2.25-1.313M9 12.75l-2.25-1.313L9 10.5l2.25 1.313m0 0L11.25 9.75l2.25 1.313M15 12.75l-2.25-1.313M15 12.75V15m0 0l2.25 1.313M15 15H12.75m3.75 0l-3.75-3.75" />
              </svg>
              <div class="text-gray-500">暂无套餐信息</div>
            </div>
          </div>

          <!-- Usage Statistics -->
          <div class="bg-white rounded-xl border border-gray-200 shadow-[0_0_16px_0_rgba(0,0,0,0.06)] p-6">
            <h3 class="text-xl font-semibold text-gray-900 mb-4">使用统计</h3>
            <div class="text-sm text-gray-500 mb-4">近30天使用情况</div>
            
            <div v-if="usageDays.length" class="overflow-x-auto">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-200">
                    <th class="text-left py-2 px-2 font-medium text-gray-700">日期</th>
                    <th class="text-right py-2 px-2 font-medium text-gray-700">语音识别</th>
                    <th class="text-right py-2 px-2 font-medium text-gray-700">语音合成</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100">
                  <tr v-for="(d, idx) in usageDays.slice(0, 10)" :key="idx" class="hover:bg-gray-50">
                    <td class="py-2 px-2 text-gray-900">{{ formatDate(d.date) }}</td>
                    <td class="py-2 px-2 text-right text-gray-600">{{ d.asr_chars || 0 }}</td>
                    <td class="py-2 px-2 text-right text-gray-600">{{ d.tts_chars || 0 }}</td>
                  </tr>
                </tbody>
              </table>
              
              <div v-if="usageDays.length > 10" class="mt-4 text-center">
                <button class="text-sm text-purple-600 hover:text-purple-700">查看全部记录</button>
              </div>
            </div>
            
            <div v-else class="text-center py-8">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 mx-auto mb-3 text-gray-300">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 0 1 3 19.875v-6.75ZM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V8.625ZM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 0 1-1.125-1.125V4.125Z" />
              </svg>
              <div class="text-gray-500">暂无使用记录</div>
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
import { getToken, setToken, clearToken, setIdentity, getIdentity, clearIdentity } from '../utils/auth'

export default defineComponent({
  data() {
    return {
      identity: getIdentity() || '',
      token: getToken() || '',
      profileData: null,
      packages: [],
      usageDays: [],
      error: '',
    }
  },
  created() {
    if (this.identity && this.token) this.refresh()
  },
  methods: {
    async doLogin() {
      this.error = ''
      try {
        const res = await account.loginSimple(this.identity)
        const data = res && (res.data || res)
        if (data && data.token) {
          setToken(data.token)
          setIdentity(this.identity)
          this.token = data.token
          await this.refresh()
        } else {
          throw new Error('未返回 token')
        }
      } catch (e) { this.error = e && e.message ? e.message : '登录失败'; try { const { toast } = require('../utils/toast'); toast(this.error, 'error') } catch (e) {} }
      finally { try { if (this.token) { const { toast } = require('../utils/toast'); toast('登录成功', 'success') } } catch (e) {} }
    },
    async doLogout() {
      try { await account.logout(this.token); const { toast } = require('../utils/toast'); toast('已退出', 'success') } catch (e) {}
      clearToken(); clearIdentity(); this.token = ''; this.identity = ''
      this.profileData = null; this.packages = []; this.usageDays = []
    },
    async refresh() {
      try {
        const [prof, packs, usage] = await Promise.all([
          account.profile(), account.packages(), account.usage({})
        ])
        this.profileData = prof || null
        this.packages = (packs && (packs.items || [])) || []
        this.usageDays = (usage && (usage.days || [])) || []
      } catch (e) { this.error = e && e.message ? e.message : '加载失败'; try { const { toast } = require('../utils/toast'); toast(this.error, 'error') } catch (e) {} }
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




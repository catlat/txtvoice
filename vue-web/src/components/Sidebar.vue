<template>
  <div class="fixed inset-y-0 left-0 z-10 h-screen w-64 transition-[left,right,width] duration-200 ease-linear md:flex border-r border-gray-200 bg-white flex flex-col">
    <!-- Header -->
    <div class="flex h-12 items-center justify-between px-5 py-0 border-b border-gray-100">
      <div class="flex items-center gap-2">
        <img src="/logo.svg" alt="vvHub" class="h-10 w-auto text-gray-900" />
      </div>
    </div>

    <!-- Navigation -->
    <div class="flex min-h-0 flex-1 flex-col gap-2 overflow-auto p-4">
      <div class="flex w-full min-w-0 flex-col gap-1">
        <router-link 
          v-for="(menu, index) in menus" 
          :key="index" 
          :to="menu.to"
          class="flex w-full items-center gap-3 rounded-md px-3 py-2.5 text-sm transition-all duration-200 hover:bg-gray-50 hover:text-gray-900"
          :class="{
            'bg-gray-100 text-gray-900 font-medium shadow-sm': $route.path === menu.to,
            'text-gray-600': $route.path !== menu.to
          }"
        >
          <component :is="menu.icon" class="w-5 h-5 flex-shrink-0" />
          <span class="truncate">{{ menu.text }}</span>
        </router-link>
      </div>
    </div>

    <!-- Footer -->
    <div class="p-3 border-t border-gray-100">
      <!-- 登录区域 -->
      <div v-if="!token" class="space-y-3">
        <input 
          v-model="identity" 
          class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent" 
          placeholder="请输入手机号或标识" 
        />
        <button 
          @click="onLogin" 
          :disabled="!identity.trim()"
          class="w-full px-3 py-2 bg-gradient-to-r from-purple-500 to-pink-500 text-white text-sm font-medium rounded-md transition-all duration-200 hover:from-purple-600 hover:to-pink-600 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          登录
        </button>
      </div>

      <!-- 已登录：头像 + 名称 + 隐藏操作 -->
      <div v-else class="relative" ref="accountBox">
        <button class="w-full flex items-center gap-3 rounded-lg px-2 py-2 hover:bg-gray-50 transition-colors" @click="open = !open">
          <div class="w-8 h-8 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center text-white text-sm font-semibold">
            {{ lastDigit(identity) }}
          </div>
          <div class="min-w-0 flex-1 text-left">
            <div class="truncate text-sm font-medium text-gray-900">{{ maskPhone(identity) }}</div>
          </div>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 text-gray-500">
            <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.24 4.5a.75.75 0 01-1.08 0l-4.24-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
          </svg>
        </button>
        <div v-show="open" class="absolute bottom-12 left-2 right-2 rounded-lg border border-gray-200 bg-white shadow-lg overflow-hidden">
          <button @click="onLogout" class="w-full px-3 py-2 text-left text-sm text-gray-700 hover:bg-gray-50">
            退出登录
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import { getIdentity, getToken, setIdentity, setToken, clearIdentity, clearToken } from '../utils/auth'
import * as account from '../api/account'

// Icon components
const HomeIcon = {
  template: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
    <path stroke-linecap="round" stroke-linejoin="round" d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />
  </svg>`
}

const HistoryIcon = {
  template: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
    <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
  </svg>`
}

const SpeakerIcon = {
  template: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
    <path stroke-linecap="round" stroke-linejoin="round" d="M19.114 5.636a9 9 0 0 1 0 12.728M16.463 8.288a5.25 5.25 0 0 1 0 7.424M6.75 8.25l4.72-4.72a.75.75 0 0 1 1.28.53v15.88a.75.75 0 0 1-1.28.53l-4.72-4.72H4.51c-.88 0-1.59-.79-1.59-1.75v-4.5c0-.96.71-1.75 1.59-1.75h2.24Z" />
  </svg>`
}

const UserIcon = {
  template: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor">
    <path stroke-linecap="round" stroke-linejoin="round" d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
  </svg>`
}

export default defineComponent({
  name: 'Sidebar',
  components: { HomeIcon, HistoryIcon, SpeakerIcon, UserIcon },
  data: () => ({
    menus: [
      { text: '首页', to: '/', icon: 'HomeIcon' },
      { text: '历史记录', to: '/history', icon: 'HistoryIcon' },
      { text: '合成历史', to: '/history/tts', icon: 'SpeakerIcon' },
      { text: '账号管理', to: '/account', icon: 'UserIcon' },
    ],
    identity: getIdentity() || '',
    token: getToken() || '',
    userInfo: null,
    open: false,
  }),
  async created() {
    if (this.token) {
      await this.loadUserInfo()
    }
    document.addEventListener('click', this.handleOutside, true)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.handleOutside, true)
  },
  watch: {
    token: {
      handler(newToken) {
        if (newToken) {
          this.loadUserInfo()
        } else {
          this.userInfo = null
        }
      }
    }
  },
  methods: {
    handleOutside(e) {
      if (!this.open) return
      const box = this.$refs.accountBox
      if (box && !box.contains(e.target)) this.open = false
    },
    lastDigit(phone) {
      if (!phone) return 'U'
      const digits = String(phone).replace(/\D/g, '')
      return digits ? digits.slice(-1) : 'U'
    },
    maskPhone(phone) {
      if (!phone) return '用户'
      const digits = String(phone).replace(/\D/g, '')
      if (digits.length < 7) return phone
      return digits.replace(/(\d{3})\d{6}(\d{2,})/, '$1******$2')
    },
    async loadUserInfo() {
      try {
        const res = await account.profile()
        const data = res && (res.data || res)
        if (data) {
          this.userInfo = { ...data, identity: this.identity }
        }
      } catch (e) {
        console.error('Failed to load user info:', e)
      }
    },
    async onLogin() {
      try {
        const res = await account.loginSimple(this.identity)
        const data = res && (res.data || res)
        if (data && data.token) {
          setToken(data.token); setIdentity(this.identity)
          this.token = data.token
          this.$emit('login', { token: data.token, identity: this.identity })
          await this.loadUserInfo()
        }
      } catch (e) {}
    },
    async onLogout() {
      try { await account.logout(this.token) } catch (e) {}
      clearToken(); clearIdentity(); 
      this.token = ''; this.identity = ''; this.userInfo = null; this.open = false
      this.$emit('logout')
    },
  },
})
</script>



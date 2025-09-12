<template>
  <div :class="[
    'fixed inset-y-0 left-0 z-10 h-screen transition-all duration-300 ease-in-out md:flex border-r border-gray-200 bg-white flex flex-col',
    collapsed ? 'w-16' : 'w-64'
  ]">
    <!-- Header -->
    <div class="flex h-12 items-center justify-between border-b border-gray-100" :class="collapsed ? 'px-2' : 'px-5'">
      <div class="flex items-center gap-2" v-if="!collapsed">
        <img src="/logo.svg" alt="vvHub" class="h-10 w-auto text-gray-900" />
      </div>
      <!-- 收缩/展开按钮 -->
      <button 
        @click="toggleCollapse"
        class="flex items-center justify-center w-8 h-8 rounded-md hover:bg-gray-100 transition-colors duration-200 text-gray-600 hover:text-gray-900"
        :class="collapsed ? 'mx-auto' : 'ml-auto'"
      >
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
          <path v-if="!collapsed" stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12" />
          <path v-else stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5M3.75 17.25h16.5" />
        </svg>
      </button>
    </div>

    <!-- Navigation -->
    <div class="flex min-h-0 flex-1 flex-col gap-2 overflow-auto" :class="collapsed ? 'p-2' : 'p-4'">
      <div class="flex w-full min-w-0 flex-col gap-1">
        <router-link 
          v-for="(menu, index) in menus" 
          :key="index" 
          :to="menu.to"
          :title="collapsed ? menu.text : ''"
          class="flex items-center rounded-md text-sm transition-all duration-200 hover:bg-gray-50 hover:text-gray-900 relative group"
          :class="[
            collapsed ? 'w-12 h-12 justify-center' : 'w-full gap-3 px-3 py-2.5',
            {
              'bg-gray-100 text-gray-900 font-medium shadow-sm': $route.path === menu.to,
              'text-gray-600': $route.path !== menu.to
            }
          ]"
        >
          <component :is="menu.icon" class="w-5 h-5 flex-shrink-0" />
          <span v-if="!collapsed" class="truncate">{{ menu.text }}</span>
          
          <!-- Tooltip for collapsed state -->
          <div v-if="collapsed" class="absolute left-full ml-2 px-2 py-1 bg-gray-900 text-white text-sm rounded-md opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 whitespace-nowrap z-50">
            {{ menu.text }}
            <div class="absolute top-1/2 left-0 transform -translate-y-1/2 -translate-x-full">
              <div class="w-0 h-0 border-t-4 border-b-4 border-r-4 border-transparent border-r-gray-900"></div>
            </div>
          </div>
        </router-link>
      </div>
    </div>

    <!-- Footer -->
    <div class="border-t border-gray-100" :class="collapsed ? 'p-2' : 'p-3'">
      <!-- 登录区域 -->
      <div v-if="!token && !collapsed" class="space-y-3">
        <input 
          v-model="identity" 
          class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent" 
          placeholder="请输入手机号" 
        />
        <input 
          v-model="password" type="password"
          class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent" 
          placeholder="请输入密码" 
        />
        <button 
          @click="onLogin" 
          :disabled="!identity.trim() || !password"
          class="w-full px-3 py-2 bg-gray-900 text-white text-sm font-medium rounded-md transition-colors duration-200 hover:bg-black disabled:opacity-50 disabled:cursor-not-allowed"
        >
          登录
        </button>
      </div>

      <!-- 收缩状态下的登录按钮 -->
      <div v-if="!token && collapsed" class="flex justify-center">
        <button 
          @click="collapsed = false" 
          title="点击展开登录"
          class="w-12 h-12 rounded-lg bg-gray-900 text-white flex items-center justify-center hover:bg-black transition-colors duration-200"
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
          </svg>
        </button>
      </div>

      <!-- 已登录：头像 + 名称 + 隐藏操作 -->
      <div v-else-if="token" class="relative" ref="accountBox">
        <button 
          :class="[
            'flex items-center rounded-lg hover:bg-gray-50 transition-colors group',
            collapsed ? 'w-12 h-12 justify-center' : 'w-full gap-3 px-2 py-2'
          ]" 
          @click="open = !open"
          :title="collapsed ? maskPhone(identity) : ''"
        >
          <div class="w-8 h-8 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center text-white text-sm font-semibold">
            {{ lastDigit(identity) }}
          </div>
          <div v-if="!collapsed" class="min-w-0 flex-1 text-left">
            <div class="truncate text-sm font-medium text-gray-900">{{ maskPhone(identity) }}</div>
          </div>
          <svg v-if="!collapsed" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" class="w-4 h-4 text-gray-500">
            <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 11.168l3.71-3.938a.75.75 0 111.08 1.04l-4.24 4.5a.75.75 0 01-1.08 0l-4.24-4.5a.75.75 0 01.02-1.06z" clip-rule="evenodd" />
          </svg>
          
          <!-- 收缩状态下的tooltip -->
          <div v-if="collapsed" class="absolute left-full ml-2 px-2 py-1 bg-gray-900 text-white text-sm rounded-md opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 whitespace-nowrap z-50">
            {{ maskPhone(identity) }}
            <div class="absolute top-1/2 left-0 transform -translate-y-1/2 -translate-x-full">
              <div class="w-0 h-0 border-t-4 border-b-4 border-r-4 border-transparent border-r-gray-900"></div>
            </div>
          </div>
        </button>
        <div v-show="open" :class="[
          'absolute bottom-12 rounded-lg border border-gray-200 bg-white shadow-lg overflow-hidden',
          collapsed ? 'left-0 right-0' : 'left-2 right-2'
        ]">
          <button @click="onLogout" class="w-full px-3 py-2 text-left text-sm text-gray-700 hover:bg-gray-50">
            退出登录
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, h } from 'vue'
import { getIdentity, getToken, setIdentity, setToken, clearIdentity, clearToken } from '../utils/auth'
import { authEvents, emitAuthLogout } from '../utils/events'
import * as account from '../api/account'

// Icon components (render functions to avoid runtime template compiler dependency)
const HomeIcon = defineComponent({
  name: 'HomeIcon',
  setup(_, { attrs }) {
    return () => h('svg', {
      xmlns: 'http://www.w3.org/2000/svg',
      fill: 'none',
      viewBox: '0 0 24 24',
      'stroke-width': '1.5',
      stroke: 'currentColor',
      ...attrs,
    }, [
      h('path', {
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        d: 'm2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25'
      })
    ])
  }
})

const HistoryIcon = defineComponent({
  name: 'HistoryIcon',
  setup(_, { attrs }) {
    return () => h('svg', {
      xmlns: 'http://www.w3.org/2000/svg',
      fill: 'none',
      viewBox: '0 0 24 24',
      'stroke-width': '1.5',
      stroke: 'currentColor',
      ...attrs,
    }, [
      h('path', {
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        d: 'M12 6v6h4.5m4.5 0a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z'
      })
    ])
  }
})

const SpeakerIcon = defineComponent({
  name: 'SpeakerIcon',
  setup(_, { attrs }) {
    return () => h('svg', {
      xmlns: 'http://www.w3.org/2000/svg',
      fill: 'none',
      viewBox: '0 0 24 24',
      'stroke-width': '1.5',
      stroke: 'currentColor',
      ...attrs,
    }, [
      h('path', {
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        d: 'M19.114 5.636a9 9 0 0 1 0 12.728M16.463 8.288a5.25 5.25 0 0 1 0 7.424M6.75 8.25l4.72-4.72a.75.75 0 0 1 1.28.53v15.88a.75.75 0 0 1-1.28.53l-4.72-4.72H4.51c-.88 0-1.59-.79-1.59-1.75v-4.5c0-.96.71-1.75 1.59-1.75h2.24Z'
      })
    ])
  }
})

const UserIcon = defineComponent({
  name: 'UserIcon',
  setup(_, { attrs }) {
    return () => h('svg', {
      xmlns: 'http://www.w3.org/2000/svg',
      fill: 'none',
      viewBox: '0 0 24 24',
      'stroke-width': '1.5',
      stroke: 'currentColor',
      ...attrs,
    }, [
      h('path', {
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        d: 'M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z'
      })
    ])
  }
})

export default defineComponent({
  name: 'Sidebar',
  components: { HomeIcon, HistoryIcon, SpeakerIcon, UserIcon },
  data: () => ({
    menus: [
      { text: '首页', to: '/', icon: HomeIcon },
      { text: '合成历史', to: '/history/tts', icon: HistoryIcon },
      { text: '账号管理', to: '/account', icon: UserIcon },
    ],
    identity: getIdentity() || '',
    password: '',
    token: getToken() || '',
    userInfo: null,
    open: false,
    collapsed: localStorage.getItem('sidebar-collapsed') === 'true',
  }),
  async created() {
    if (this.token) {
      await this.loadUserInfo()
    }
    document.addEventListener('click', this.handleOutside, true)
    // 订阅全局登录/登出事件，保持左下角与弹框状态同步
    try {
      authEvents.addEventListener('auth:login', this.onAuthLogin)
      authEvents.addEventListener('auth:logout', this.onAuthLogout)
    } catch (e) {}
    // 发射初始收缩状态
    this.$emit('sidebar-toggle', this.collapsed)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.handleOutside, true)
    try {
      authEvents.removeEventListener('auth:login', this.onAuthLogin)
      authEvents.removeEventListener('auth:logout', this.onAuthLogout)
    } catch (e) {}
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
    },
    collapsed: {
      handler(newValue) {
        // 持久化收缩状态
        localStorage.setItem('sidebar-collapsed', String(newValue))
        // 发射状态变化事件
        this.$emit('sidebar-toggle', newValue)
      }
    }
  },
  methods: {
    onAuthLogin(e) {
      const detail = (e && e.detail) || {}
      const t = detail.token || getToken()
      const id = detail.identity || getIdentity()
      if (t) this.token = t
      if (id) this.identity = id
      this.loadUserInfo()
    },
    onAuthLogout() {
      this.token = ''
      this.identity = ''
      this.userInfo = null
    },
    toggleCollapse() {
      this.collapsed = !this.collapsed
    },
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
        const res = await account.login(this.identity, this.password)
        const data = res && (res.data || res)
        if (data && data.token) {
          setToken(data.token); setIdentity(this.identity)
          this.token = data.token
          this.password = ''
          this.$emit('login', { token: data.token, identity: this.identity })
          await this.loadUserInfo()
        }
      } catch (e) {}
    },
    async onLogout() {
      try { await account.logout(this.token) } catch (e) {}
      clearToken(); clearIdentity(); 
      this.token = ''; this.identity = ''; this.password = ''; this.userInfo = null; this.open = false
      try { emitAuthLogout() } catch (e) {}
      this.$emit('logout')
    },
  },
})
</script>




<template>
  <div class="flex">
    <Sidebar @login="onLoginSuccess" @logout="logout" @sidebar-toggle="onSidebarToggle" />
    <div :class="[
      'flex-1 min-h-screen bg-gray-50 transition-all duration-300 ease-in-out',
      sidebarCollapsed ? 'ml-16' : 'ml-64'
    ]">
      <div class="max-w-screen-xl mx-auto py-6 px-6">
        <router-view />
      </div>
      <LoginModal v-model="showLogin" @success="onLoginSuccess" />
      <Toast />
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import { authEvents } from './utils/events'
import { clearToken, clearIdentity } from './utils/auth'
import LoginModal from './components/LoginModal.vue'
import Toast from './components/Toast.vue'
import Sidebar from './components/Sidebar.vue'

export default defineComponent({
  components: { LoginModal, Toast, Sidebar },
  data: () => ({
    identity: '',
    token: '',
    showLogin: false,
    sidebarCollapsed: localStorage.getItem('sidebar-collapsed') === 'true',
  }),
  created() {
    try {
      const { getIdentity, getToken } = require('./utils/auth')
      this.identity = getIdentity() || ''
      this.token = getToken() || ''
    } catch (e) {}
    try {
      authEvents.addEventListener('auth:login', this.onAuthLogin)
      authEvents.addEventListener('auth:logout', this.onAuthLogout)
    } catch (e) {}
  },
  beforeUnmount() {
    try {
      authEvents.removeEventListener('auth:login', this.onAuthLogin)
      authEvents.removeEventListener('auth:logout', this.onAuthLogout)
    } catch (e) {}
  },
  methods: {
    async login() { this.showLogin = true },
    async logout() {
      try {
        const account = require('./api/account')
        await account.logout(this.token)
      } catch (e) {}
      try { clearToken(); clearIdentity(); } catch (e) {}
      this.token = ''
      this.identity = ''
      try { this.$router.replace({ path: '/' }) } catch (e) {}
      // 轻量刷新，确保全局状态与后端 cookie 清理一致
      try { location.reload() } catch (e) {}
    },
    onAuthLogin(e) {
      const detail = (e && e.detail) || {}
      const t = detail.token
      const id = detail.identity
      if (t) this.token = t
      if (id) this.identity = id
      this.showLogin = false
    },
    onAuthLogout() {
      this.token = ''
      this.identity = ''
      try { this.$router.replace({ path: '/' }) } catch (e) {}
      try { location.reload() } catch (e) {}
    },
    onLoginSuccess({ token, identity }) {
      this.token = token
      this.identity = identity
    },
    onSidebarToggle(collapsed) {
      this.sidebarCollapsed = collapsed
    },
  },
})
</script>
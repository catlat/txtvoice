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
  },
  methods: {
    async login() { this.showLogin = true },
    async logout() {
      try {
        const account = require('./api/account')
        await account.logout(this.token)
      } catch (e) {}
      try {
        const { clearToken, clearIdentity } = require('./utils/auth')
        clearToken(); clearIdentity();
      } catch (e) {}
      this.token = ''
      this.identity = ''
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
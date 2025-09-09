<template>
  <div v-if="modelValue" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
    <div class="w-full max-w-sm p-4 bg-white rounded shadow">
      <h3 class="text-lg font-semibold">登录</h3>
      <div class="mt-3">
        <input v-model="identity" class="w-full px-3 py-2 border rounded" placeholder="输入手机号/标识" />
      </div>
      <div class="flex justify-end gap-2 mt-4">
        <button class="px-3 py-2 bg-gray-200 rounded" @click="$emit('update:modelValue', false)">取消</button>
        <button class="px-3 py-2 text-white bg-indigo-600 rounded" :disabled="!identity" @click="onLogin">登录</button>
      </div>
      <div class="mt-2 text-red-600 text-sm" v-if="error">{{ error }}</div>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import * as account from '../api/account'
import { setToken, setIdentity } from '../utils/auth'

export default defineComponent({
  name: 'LoginModal',
  props: { modelValue: { type: Boolean, default: false } },
  emits: ['update:modelValue', 'success'],
  data: () => ({ identity: '', error: '' }),
  methods: {
    async onLogin() {
      this.error = ''
      try {
        const res = await account.loginSimple(this.identity)
        const data = res && (res.data || res)
        if (data && data.token) {
          setToken(data.token); setIdentity(this.identity)
          this.$emit('success', { token: data.token, identity: this.identity })
          this.$emit('update:modelValue', false)
        } else {
          throw new Error('未返回 token')
        }
      } catch (e) { this.error = e && e.message ? e.message : '登录失败' }
    },
  },
})
</script>




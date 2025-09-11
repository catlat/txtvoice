import { reactive } from 'vue'

export const toastState = reactive({ show: false, message: '', kind: 'info' })

export function toast(message, kind = 'info', duration = 3000) {
  toastState.message = String(message || '')
  toastState.kind = kind
  toastState.show = true
  if (toast._timer) clearTimeout(toast._timer)
  toast._timer = setTimeout(() => {
    toastState.show = false
  }, duration)
}




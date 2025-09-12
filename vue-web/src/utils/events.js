export const authEvents = new EventTarget()

export function emitAuthLogin(token, identity) {
  try {
    authEvents.dispatchEvent(new CustomEvent('auth:login', { detail: { token, identity } }))
  } catch (_) {}
}

export function emitAuthLogout() {
  try {
    authEvents.dispatchEvent(new CustomEvent('auth:logout'))
  } catch (_) {}
}

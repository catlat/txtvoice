const TOKEN_KEY = 'auth:token'
const IDENTITY_KEY = 'auth:identity'

export function getToken() {
  return localStorage.getItem(TOKEN_KEY) || ''
}

export function setToken(token) {
  if (token) localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

export function getIdentity() {
  return localStorage.getItem(IDENTITY_KEY) || ''
}

export function setIdentity(identity) {
  if (identity) localStorage.setItem(IDENTITY_KEY, identity)
}

export function clearIdentity() {
  localStorage.removeItem(IDENTITY_KEY)
}




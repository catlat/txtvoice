import { getToken, getIdentity } from '../utils/auth'

// const BASE = 'http://43.133.180.227:9005/api'
const BASE = 'http://127.0.0.1:9005/api'
//  const BASE = '/api'
function resolveUrl(path) {
  if (path.startsWith('http')) return path
  return `${BASE}${path}`
}

function extractBizCode(payload) {
  if (!payload || typeof payload !== 'object') return undefined
  // 常见业务码字段：status / code / errno
  return payload.status ?? payload.code ?? payload.errno
}

function extractBizMessage(payload, fallback) {
  if (!payload || typeof payload !== 'object') return fallback
  return payload.msg || payload.message || fallback
}

export async function request(path, options = {}) {
  const {
    method = 'GET',
    data,
    headers = {},
    credentials = 'include',
  } = options
  const token = getToken()
  const identity = getIdentity()
  const init = { method, headers: { 'Content-Type': 'application/json', ...headers }, credentials }
  if (token) init.headers['token'] = token
  if (identity) {
    // 追加 identity 查询参数（GET/POST均追加，后端从 query 读取）
    if (path.includes('?')) path += `&identity=${encodeURIComponent(identity)}`
    else path += `?identity=${encodeURIComponent(identity)}`
  }
  if (data !== undefined) init.body = typeof data === 'string' ? data : JSON.stringify(data)
  try {
    const res = await fetch(resolveUrl(path), init)
    const contentType = res.headers.get('content-type') || ''
    const isJson = contentType.includes('application/json')
    const payload = isJson ? await res.json().catch(() => ({})) : await res.text()

    // HTTP 层错误
    if (!res.ok) {
      const message = extractBizMessage(payload, res.statusText || '请求失败')
      try { const { toast } = require('../utils/toast'); toast(message, 'error') } catch (e) {}
      throw new Error(message)
    }

    // 业务层错误（即使 HTTP 200）
    if (isJson && payload && typeof payload === 'object') {
      const bizCode = extractBizCode(payload)
      // 约定：0 / 200 / 20000 / undefined 视为成功，其它非零为错误
      if (bizCode !== undefined && bizCode !== 0 && bizCode !== 200 && bizCode !== 20000) {
        const message = extractBizMessage(payload, `请求失败(${bizCode})`)
        try { const { toast } = require('../utils/toast'); toast(message, 'error') } catch (e) {}
        throw new Error(message)
      }
    }

    return payload
  } catch (err) {
    const msg = err && err.message ? err.message : '网络错误，请稍后再试'
    try { const { toast } = require('../utils/toast'); toast(msg, 'error') } catch (e) {}
    throw err
  }
}

export function get(path, params) {
  let url = path
  if (params && Object.keys(params).length) {
    const qs = new URLSearchParams(params).toString()
    url = `${path}?${qs}`
  }
  return request(url)
}

export function post(path, data) {
  return request(path, { method: 'POST', data })
}

export function postForm(path, formObj) {
  const token = getToken()
  const identity = getIdentity()
  const params = new URLSearchParams()
  for (const k in formObj || {}) {
    if (formObj[k] !== undefined && formObj[k] !== null) params.append(k, String(formObj[k]))
  }
  const headers = { 'Content-Type': 'application/x-www-form-urlencoded' }
  if (token) headers['token'] = token
  let url = path
  if (identity) url += (url.includes('?') ? '&' : '?') + `identity=${encodeURIComponent(identity)}`
  return request(url, { method: 'POST', data: params.toString(), headers })
}






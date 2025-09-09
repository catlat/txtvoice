import { post } from './http'

export function synthesize(payload) {
  // 后端从 query 读取 identity，这里仍走 JSON body
  return post('/tts/synthesize', payload)
}




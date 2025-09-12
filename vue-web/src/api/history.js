import { get } from './http'

export function listTTS(params) {
  return get('/history/tts', params)
}





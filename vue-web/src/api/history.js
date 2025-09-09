import { get } from './http'

export function listVideos(params) {
  return get('/history/videos', params)
}

export function getVideo(site, id) {
  return get(`/history/video/${encodeURIComponent(site)}/${encodeURIComponent(id)}`)
}

export function listTTS(params) {
  return get('/history/tts', params)
}





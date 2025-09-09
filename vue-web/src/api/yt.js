import { post } from './http'

export function info(idOrUrl, platform) {
  const payload = { id_or_url: idOrUrl }
  if (platform) payload.platform = platform
  return post('/yt/info', payload)
}

export function text(idOrUrl, platform) {
  const payload = { id_or_url: idOrUrl }
  if (platform) payload.platform = platform
  return post('/yt/text', payload)
}





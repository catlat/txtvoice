import { post } from './http'

export function info(idOrUrl) {
  return post('/yt/info', { id_or_url: idOrUrl })
}

export function text(idOrUrl) {
  return post('/yt/text', { id_or_url: idOrUrl })
}





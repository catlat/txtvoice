export function normalizeAndCount(text) {
  const normalized = (text || '').replace(/\r\n?/g, '\n').trim()
  const count = [...normalized].length
  return { text: normalized, count }
}





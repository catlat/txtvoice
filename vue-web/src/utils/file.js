export async function parseTxt(file) {
  const text = await file.text()
  return text
}

export async function parseMd(file) {
  const text = await file.text()
  // very naive markdown removal: strip code blocks and inline marks
  let t = text
    .replace(/```[\s\S]*?```/g, '')
    .replace(/^\s{0,3}#{1,6}\s+/gm, '')
    .replace(/\!\[[^\]]*\]\([^)]*\)/g, '')
    // keep link text, drop url
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/\*\*([^*]+)\*\*/g, '$1')
    .replace(/\*([^*]+)\*/g, '$1')
  return t
}

export async function parseDocx(file) {
  try {
    // 动态导入 mammoth 浏览器版本
    const mod = await import('mammoth/mammoth.browser')
    const mammoth = mod && mod.default ? mod.default : mod

    // 将文件转换为 ArrayBuffer
    const arrayBuffer = await file.arrayBuffer()

    // 优先转为 HTML，便于保留段落、列表等块级信息
    const result = await mammoth.convertToHtml({ arrayBuffer })
    const html = (result && result.value) || ''

    // 将 HTML 映射为纯文本，保留段落/换行/列表
    let text = html
      // 列表项前缀
      .replace(/<li[^>]*>/gi, '\n- ')
      // 段内换行
      .replace(/<br\s*\/?>(?=\s*<)/gi, '\n')
      .replace(/<br\s*\/?>(?!\s*<)/gi, '\n')
      // 块级元素结束添加换行
      .replace(/<\/(p|h[1-6]|div)>/gi, '\n\n')
      .replace(/<\/(ul|ol)>/gi, '\n')
      // 移除其余标签
      .replace(/<[^>]+>/g, '')
      // 解码常见实体
      .replace(/&nbsp;/g, ' ')
      .replace(/&amp;/g, '&')
      .replace(/&lt;/g, '<')
      .replace(/&gt;/g, '>')
      // 规整空行
      .replace(/\u00a0/g, ' ')
      .replace(/[ \t]+\n/g, '\n')
      .replace(/\n{3,}/g, '\n\n')
      .trim()

    return text
  } catch (error) {
    console.error('解析Word文档失败:', error)
    throw new Error('无法解析Word文档，请确保文件格式正确')
  }
}





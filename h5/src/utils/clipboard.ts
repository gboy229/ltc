/**
 * 复制文本，优先使用 Clipboard API，失败时回退到 textarea + execCommand。
 * 返回是否复制成功。
 */
export async function copyText(text: string): Promise<boolean> {
  if (!text) {
    return false
  }

  try {
    if (navigator.clipboard?.writeText) {
      try {
        await navigator.clipboard.writeText(text)
        return true
      } catch {
        // 浏览器拒绝授权时继续走 textarea 回退
      }
    }

    const textarea = document.createElement('textarea')
    textarea.value = text
    textarea.setAttribute('readonly', '')
    textarea.style.position = 'fixed'
    textarea.style.top = '0'
    textarea.style.left = '0'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.focus()
    textarea.select()
    const copied = document.execCommand('copy')
    textarea.remove()
    return copied
  } catch {
    return false
  }
}

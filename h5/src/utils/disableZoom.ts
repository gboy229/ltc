/**
 * 禁止页面缩放。
 *
 * viewport 的 user-scalable=no 在 iOS Safari 上会被忽略，
 * 需要额外拦截非标准的 gesture 事件；部分 Android 浏览器
 * 还需要阻止多指 touchmove，桌面端则要拦截 Ctrl + 滚轮。
 */
const GESTURE_EVENTS = ['gesturestart', 'gesturechange', 'gestureend'] as const

let installed = false

function preventDefault(event: Event): void {
  event.preventDefault()
}

export function disablePageZoom(): void {
  if (installed || typeof window === 'undefined') {
    return
  }
  installed = true

  // iOS Safari：双指缩放触发的是非标准 gesture 事件
  GESTURE_EVENTS.forEach((name) => {
    document.addEventListener(name, preventDefault, { passive: false })
  })

  // Android 等：双指触摸移动时阻止默认的缩放行为
  document.addEventListener(
    'touchmove',
    (event) => {
      if (event.touches.length > 1) {
        event.preventDefault()
      }
    },
    { passive: false },
  )

  // 桌面端：阻止 Ctrl + 滚轮（含触控板双指捏合）缩放页面
  window.addEventListener(
    'wheel',
    (event) => {
      if (event.ctrlKey) {
        event.preventDefault()
      }
    },
    { passive: false },
  )
}

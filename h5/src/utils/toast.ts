type ToastType = 'info' | 'success' | 'error'

const TOAST_ROOT_ID = 'h5-toast-root'

const typeClass: Record<ToastType, string> = {
  info: 'bg-[#3b2a18]/95 text-white',
  success: 'bg-[#986127]/95 text-white',
  error: 'bg-rose-600/95 text-white',
}

const containerClass =
  'pointer-events-none fixed inset-x-0 top-0 z-[9999] flex flex-col items-center gap-2 px-4 pt-[calc(env(safe-area-inset-top)+1rem)]'

let container: HTMLElement | undefined

function getContainer(): HTMLElement {
  // H5Shell 内预置了 toast 容器；PC 端会跟随手机画布缩放
  const host = document.getElementById(TOAST_ROOT_ID)
  if (host) {
    return host
  }

  if (container && document.body.contains(container)) {
    return container
  }

  container = document.createElement('div')
  container.className = containerClass
  document.body.appendChild(container)
  return container
}

/**
 * 轻量全局提示，避免为初始化引入 UI 组件库。
 */
export function toast(message: string, type: ToastType = 'info', duration = 2200): void {
  if (!message) {
    return
  }

  const el = document.createElement('div')
  el.className = [
    'pointer-events-auto max-w-xs rounded-xl px-4 py-2.5 text-center text-sm font-medium shadow-lg shadow-slate-900/10 backdrop-blur transition-opacity duration-200 md:text-[15px]',
    'opacity-100',
    typeClass[type],
  ].join(' ')
  el.setAttribute('role', 'status')
  el.textContent = message

  getContainer().appendChild(el)

  window.setTimeout(() => {
    el.classList.remove('opacity-100')
    el.classList.add('opacity-0')
    window.setTimeout(() => el.remove(), 220)
  }, duration)
}

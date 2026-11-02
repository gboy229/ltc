import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './App.vue'
import router from './router'
import { clearAuthStorage } from './utils/auth'
import { disablePageZoom } from './utils/disableZoom'

import './styles/index.css'

// 构建版本标记：线上已在控制台核对版本时使用（window.__H5_BUILD__）
const BUILD_TIME = typeof __BUILD_TIME__ === 'string' ? __BUILD_TIME__ : 'dev'
if (typeof window !== 'undefined') {
  window.__H5_BUILD__ = BUILD_TIME
  console.info('[h5] build', BUILD_TIME)
}

// H5 作为 App 使用，禁止浏览器双指/滚轮缩放页面
disablePageZoom()

// 切换后端环境后，旧 token / 钱包地址属于旧环境，自动清理
const API_BASE = import.meta.env.VITE_API_BASE_URL || '/api'
const API_BASE_KEY = 'h5_api_base'
try {
  const savedBase = localStorage.getItem(API_BASE_KEY)
  if (savedBase && savedBase !== API_BASE) {
    clearAuthStorage()
  }
  localStorage.setItem(API_BASE_KEY, API_BASE)
} catch {
  // 隐私模式或存储不可用时忽略
}

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.mount('#app')

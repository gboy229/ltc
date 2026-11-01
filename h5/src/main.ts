import { createPinia } from 'pinia'
import { createApp } from 'vue'

import App from './App.vue'
import router from './router'
import { clearAuthStorage } from './utils/auth'

import './styles/index.css'

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

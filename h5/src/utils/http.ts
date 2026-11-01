import axios, { type AxiosError, type AxiosRequestConfig } from 'axios'

import { getToken, removeToken } from './auth'
import { toast } from './toast'

declare module 'axios' {
  export interface AxiosRequestConfig {
    /** 静默请求：失败时不弹出全局错误提示 */
    silent?: boolean
  }
}

/** 登录失效事件，由 user store 监听后清理本地会话 */
export const UNAUTHORIZED_EVENT = 'h5:unauthorized'

export interface ApiResult<T = unknown> {
  code: number
  message: string
  data: T
  timestamp?: number
  traceID?: string
}

export class ApiError extends Error {
  readonly code: number

  constructor(code: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.code = code
  }
}

/** 业务成功码 */
const SUCCESS_CODE = 0
/** GoFrame gcode.CodeNotAuthorized，登录失效 */
const UNAUTHORIZED_CODE = 61

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15_000,
  headers: {
    'Content-Type': 'application/json;charset=UTF-8',
  },
})

instance.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

instance.interceptors.response.use(
  (response) => {
    const result = response.data as ApiResult | undefined

    // 非标准响应（文件流等）直接返回原始数据
    if (!result || typeof result.code !== 'number') {
      return response.data
    }

    if (result.code === SUCCESS_CODE) {
      return result.data
    }

    if (result.code === UNAUTHORIZED_CODE) {
      removeToken()
      window.dispatchEvent(new Event(UNAUTHORIZED_EVENT))
    }

    const message = result.message || '请求失败'
    if (!response.config.silent) {
      toast(message, 'error')
    }
    return Promise.reject(new ApiError(result.code, message))
  },
  (error: AxiosError<ApiResult>) => {
    if (axios.isCancel(error)) {
      return Promise.reject(error)
    }

    let message = '网络异常，请稍后重试'

    if (error.code === 'ECONNABORTED' || error.message.includes('timeout')) {
      message = '请求超时，请稍后重试'
    } else if (error.response?.data?.message) {
      message = error.response.data.message
    } else if (error.response?.status === 401) {
      removeToken()
      window.dispatchEvent(new Event(UNAUTHORIZED_EVENT))
      message = '登录已失效，请重新登录'
    }

    if (!error.config?.silent) {
      toast(message, 'error')
    }
    return Promise.reject(error)
  },
)

function request<T>(config: AxiosRequestConfig): Promise<T> {
  // 响应拦截器已将统一响应解包为 data，这里需要显式断言返回类型
  return instance.request(config) as unknown as Promise<T>
}

/**
 * 统一请求入口：自动携带 token、解包响应 data、集中提示错误。
 * 业务接口请放在 src/api 下按模块封装，不要在页面里直接调用 axios。
 */
export const http = {
  get<T>(url: string, params?: object, config?: AxiosRequestConfig): Promise<T> {
    return request<T>({ url, method: 'GET', params, ...config })
  },
  post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return request<T>({ url, method: 'POST', data, ...config })
  },
  put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    return request<T>({ url, method: 'PUT', data, ...config })
  },
  delete<T>(url: string, params?: object, config?: AxiosRequestConfig): Promise<T> {
    return request<T>({ url, method: 'DELETE', params, ...config })
  },
}

export default http

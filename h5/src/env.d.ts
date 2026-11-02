/// <reference types="vite/client" />

/** 构建时间（vite define 注入），用于确认线上运行的前端版本 */
declare const __BUILD_TIME__: string

interface Window {
  /** 构建时间（main.ts 运行时写入，便于在控制台核对版本） */
  __H5_BUILD__?: string
}

interface ImportMetaEnv {
  /** 应用标题 */
  readonly VITE_APP_TITLE: string
  /** 接口基础路径 */
  readonly VITE_API_BASE_URL: string
  /** 开发环境接口代理目标 */
  readonly VITE_PROXY_TARGET: string
  /** 开发服务器端口 */
  readonly VITE_PORT: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

/// <reference types="vite/client" />

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

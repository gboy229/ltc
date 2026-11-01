import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    /** 页面标题，用于 document.title */
    title?: string
    /** 是否需要登录后访问 */
    requiresAuth?: boolean
  }
}

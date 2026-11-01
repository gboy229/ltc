import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const DEFAULT_TITLE = import.meta.env.VITE_APP_TITLE || 'H5'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomeView.vue'),
    meta: { title: '首页' },
  },
  {
    path: '/orders',
    name: 'orders',
    component: () => import('@/views/OrderDetailView.vue'),
    meta: { title: '报单明细', requiresAuth: true },
  },
  {
    path: '/earnings',
    name: 'earnings',
    component: () => import('@/views/EarningsView.vue'),
    meta: { title: '收益明细', requiresAuth: true },
  },
  {
    path: '/community',
    name: 'community',
    component: () => import('@/views/CommunityView.vue'),
    meta: { title: '我的社区', requiresAuth: true },
  },
  {
    path: '/community/team',
    name: 'community-team',
    component: () => import('@/views/TeamView.vue'),
    meta: { title: '我的直推', requiresAuth: true },
  },
  {
    path: '/my-assets',
    name: 'assets',
    component: () => import('@/views/AssetsView.vue'),
    meta: { title: '我的资产', requiresAuth: true },
  },
  {
    path: '/withdraw',
    name: 'withdraw',
    component: () => import('@/views/WithdrawView.vue'),
    meta: { title: '在线提币', requiresAuth: true },
  },
  {
    path: '/withdrawals',
    name: 'withdrawals',
    component: () => import('@/views/WithdrawalListView.vue'),
    meta: { title: '提币记录', requiresAuth: true },
  },
  {
    path: '/notice',
    name: 'notice',
    component: () => import('@/views/NoticeView.vue'),
    meta: { title: '平台公告' },
  },
  {
    path: '/notice/:id',
    name: 'notice-detail',
    component: () => import('@/views/NoticeDetailView.vue'),
    meta: { title: '公告详情' },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundView.vue'),
    meta: { title: '页面不存在' },
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} - ${DEFAULT_TITLE}` : DEFAULT_TITLE
})

export default router

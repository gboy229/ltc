import { http } from '@/utils/http'

/** 投资套餐（前台展示所需字段，挡位由套餐列表反推） */
export interface PledgePlan {
  name: string
  amount: number
  periodMonths: number
  monthlyRewardInfo: string
  monthlyRewardRate: number
  totalReturn: number
}

export interface ListRes<T> {
  list: T[]
}

/** 获取启用的投资套餐（返回全部启用套餐） */
export function getPlans() {
  return http.get<ListRes<PledgePlan>>('/pledge/plan/list')
}

export interface BuyOrderRes {
  orderNo: string
  amount: number
  toAddress: string
  status: number
}

export interface BuyOrderStatusRes {
  orderNo: string
  amount: number
  status: number
  pledgeOrderStatus: number
  paidAt: string
}

export type PledgeOrderStatus = 'active' | 'completed'

export interface PledgeOrderItem {
  id: string
  amount: number
  todayRelease: number
  totalRelease: number
  period: number
  investDate: string
  status: PledgeOrderStatus
  earlyExitPenaltyRate: number
}

export interface PledgeOrderListRes {
  list: PledgeOrderItem[]
  page: number
  pageSize: number
  pageCount: number
  totalCount: number
}

export interface PledgeOrderStatistics {
  total: number
  todayRelease: number
  totalRelease: number
  earningsTotal: number
  referralIncome: number
  staticIncome: number
}

/** 报单统计 */
export function getOrderStatistics() {
  return http.get<PledgeOrderStatistics>('/pledge/order/statistics')
}

/** 报单列表（正式订单） */
export function getOrderList(params?: {
  page?: number
  pageSize?: number
  status?: PledgeOrderStatus
}) {
  return http.get<PledgeOrderListRes>('/pledge/order/list', params)
}

/** 提交报单（后端根据 planName + periodMonths 解析套餐） */
export function buyPlan(params: { planName: string; periodMonths: number }) {
  return http.post<BuyOrderRes>('/pledge/order/buy', params)
}

export interface RefundOrderRes {
  amount: number
}

/** 提前退本 */
export function refundPledgeOrder(params: {
  orderNo: string
  nonce: string
  signature: string
}) {
  return http.post<RefundOrderRes>('/pledge/order/refund', params)
}

/** 查询报单处理状态 */
export function getBuyOrderStatus(orderNo: string) {
  return http.get<BuyOrderStatusRes>(
    '/pledge/order/status',
    { orderNo },
    { silent: true },
  )
}

export type IncomeLogType = 'static' | 'direct' | 'diff' | 'lot' | 'indirect' | 'team' | 'principal'

export interface IncomeLogItem {
  type: IncomeLogType
  label: string
  date: string
  amount: number
}

export interface IncomeLogListRes {
  list: IncomeLogItem[]
  page: number
  pageSize: number
  pageCount: number
  totalCount: number
}

/** 当前用户收益流水 */
export function getIncomeLogs(params?: {
  page?: number
  pageSize?: number
  type?: IncomeLogType
}) {
  return http.get<IncomeLogListRes>('/pledge/income/list', params)
}

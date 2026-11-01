import { http } from '@/utils/http'

export interface AssetAccount {
  total: number
  available: number
  frozen: number
}

export type AssetLogType = 'income' | 'expense'

export interface AssetLogItem {
  type: AssetLogType
  remark: string
  date: string
  amount: number
}

export interface AssetLogListRes {
  list: AssetLogItem[]
  page: number
  pageSize: number
  pageCount: number
  totalCount: number
}

/** 当前用户资产账户 */
export function getAssetAccount() {
  return http.get<AssetAccount>('/users/asset/account')
}

/** 当前用户资产流水 */
export function getAssetLogs(params: {
  page?: number
  pageSize?: number
  type: AssetLogType
}) {
  return http.get<AssetLogListRes>('/users/asset/logs', params)
}

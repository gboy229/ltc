import { http } from '@/utils/http'

export interface WithdrawalOption {
  chain: string
  network: string
  symbol: string
  decimals: number
  feeType: number
  rate: number
  minAmount: number
  maxAmount: number
}

export interface WithdrawalOptionRes {
  list: WithdrawalOption[]
}

export type WithdrawalStatus = 'processing' | 'success' | 'failed'

export interface WithdrawalRecord {
  id: string
  chain: string
  network: string
  symbol: string
  status: WithdrawalStatus
  address: string
  hash: string
  date: string
  amount: number
  fee: number
  receiveAmount: number
}

export interface WithdrawalListRes {
  page: number
  pageSize: number
  pageCount: number
  totalCount: number
  list: WithdrawalRecord[]
}

export interface ApplyWithdrawalParams {
  chain: string
  symbol: string
  amount: number
  nonce: string
  signature: string
}

export function getWithdrawalOptions() {
  return http.get<WithdrawalOptionRes>('/chain/withdrawal/options')
}

export function applyWithdrawal(params: ApplyWithdrawalParams) {
  return http.post<Record<string, never>>('/chain/withdrawal/apply', params)
}

export function getWithdrawalList(params?: {
  page?: number
  pageSize?: number
  status?: 'processing' | 'success'
}) {
  return http.get<WithdrawalListRes>('/chain/withdrawal/list', params)
}

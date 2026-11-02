import { http } from '@/utils/http'

export interface ChainNativeCurrency {
  name: string
  symbol: string
  decimals: number
}

export interface ChainContract {
  symbol: string
  address: string
  decimals: number
}

export interface ChainConfig {
  chain: string
  chainId: number
  chainName: string
  nativeCurrency: ChainNativeCurrency
  rpcUrls: string[]
  blockExplorerUrls: string[]
  contracts: ChainContract[]
}

export interface ChainConfigRes {
  chains: ChainConfig[]
}

/** 获取后台配置的支付链 */
export function getChainConfig() {
  return http.get<ChainConfigRes>('/chain/config', undefined, { silent: true })
}

export interface ChainBalanceRes {
  chain: string
  symbol: string
  address: string
  contract: string
  decimals: number
  amount: number
}

/** 服务端查询钱包代币余额（避免浏览器无法访问 RPC） */
export function getChainBalance(params: { address: string; chain?: string; symbol?: string }) {
  return http.get<ChainBalanceRes>('/chain/balance', params, { silent: true })
}

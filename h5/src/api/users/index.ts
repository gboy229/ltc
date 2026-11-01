import { http } from '@/utils/http'

export interface UserInfo {
  uid: string
  nickname: string
  avatar: string
  /** 以下字段后端后续扩展，当前可能为空 */
  address?: string
  walletAddress?: string
  levelName?: string
}

export interface NonceRes {
  nonce: string
}

export interface LoginRes {
  token: string
  expires: number
}

export interface WalletLoginPayload {
  type: 'wallet'
  account: string
  signature: string
  nonce: string
  /** 推荐人用户号，首次注册时绑定 */
  uid?: string
}

export interface InfoRes {
  user: UserInfo
}

export interface AuthStatusRes {
  registered: boolean
  isFirstUser: boolean
  forceInvite: boolean
}

export interface AuthStatusPayload {
  type: 'wallet' | 'username' | 'mobile' | 'email'
  account: string
}

/** 获取钱包签名挑战值 */
export function getNonce(address: string, purpose: 'login' | 'action' = 'login') {
  return http.get<NonceRes>('/users/auth/nonce', { address, purpose })
}

/** 钱包签名登录（首次登录自动注册） */
export function loginWithWallet(payload: WalletLoginPayload) {
  return http.post<LoginRes>('/users/auth/login', payload)
}

/** 查询账号是否已注册，用于注册前判断 */
export function getAuthStatus(payload: AuthStatusPayload) {
  return http.post<AuthStatusRes>('/users/auth/status', payload)
}

/** 当前登录用户信息 */
export function getProfile(silent = false) {
  return http.get<InfoRes>('/users/center/info', undefined, { silent })
}

/** 退出登录 */
export function logout(silent = true) {
  return http.post<Record<string, never>>('/users/auth/logout', undefined, { silent })
}

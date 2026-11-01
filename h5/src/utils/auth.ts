const TOKEN_KEY = 'h5_token'
const ADDRESS_KEY = 'h5_wallet_address'
const WALLET_KEY = 'h5_wallet_id'
const INVITE_UID_KEY = 'h5_invite_uid'

function read(key: string): string {
  try {
    return localStorage.getItem(key) || ''
  } catch {
    // 隐私模式或存储不可用时静默失败
    return ''
  }
}

function write(key: string, value: string): void {
  try {
    localStorage.setItem(key, value)
  } catch {
    // 忽略存储异常
  }
}

function remove(key: string): void {
  try {
    localStorage.removeItem(key)
  } catch {
    // 忽略存储异常
  }
}

export function getToken(): string {
  return read(TOKEN_KEY)
}

export function setToken(token: string): void {
  write(TOKEN_KEY, token)
}

export function removeToken(): void {
  remove(TOKEN_KEY)
}

/** 最近一次登录使用的钱包地址（用于刷新后免签名恢复会话） */
export function getAddress(): string {
  return read(ADDRESS_KEY)
}

export function setAddress(address: string): void {
  write(ADDRESS_KEY, address.toLowerCase())
}

export function removeAddress(): void {
  remove(ADDRESS_KEY)
}

/** 最近一次使用的钱包标识，用于多钱包环境下优先选择 */
export function getWalletId(): string {
  return read(WALLET_KEY)
}

export function setWalletId(id: string): void {
  write(WALLET_KEY, id)
}

export function removeWalletId(): void {
  remove(WALLET_KEY)
}

/** 邀请人用户号，从 URL ?uid= 记录，登录时上报给后端 */
export function getInviteUid(): string {
  const fromUrl = new URLSearchParams(window.location.search).get('uid')?.trim() || ''
  if (fromUrl) {
    write(INVITE_UID_KEY, fromUrl)
    return fromUrl
  }
  return read(INVITE_UID_KEY)
}

/** 清空本地登录相关数据（切换后端环境/强制重新登录时使用） */
export function clearAuthStorage(): void {
  removeToken()
  removeAddress()
  removeWalletId()
  remove(INVITE_UID_KEY)
}

/**
 * Web3 钱包适配层（EIP-1193 / EIP-6963）。
 *
 * 目标：兼容 MetaMask、TokenPocket、OKX、Bitget、Trust、imToken 等
 * 主流 EVM 钱包，并在移动端钱包内置浏览器中自动完成连接与签名登录。
 */

export type WalletId =
  | 'metamask'
  | 'tokenpocket'
  | 'okx'
  | 'bitget'
  | 'trust'
  | 'imtoken'
  | 'injected'

type ProviderListener = (...args: unknown[]) => void

export interface Eip1193Provider {
  request(args: { method: string; params?: unknown[] }): Promise<unknown>
  on?(event: string, listener: ProviderListener): void
  removeListener?(event: string, listener: ProviderListener): void
  isMetaMask?: boolean
  isTokenPocket?: boolean
  isOkxWallet?: boolean
  isOKExWallet?: boolean
  isBitKeep?: boolean
  isBitget?: boolean
  isTrust?: boolean
  isTrustWallet?: boolean
  isImToken?: boolean
  isImtoken?: boolean
  providers?: Eip1193Provider[]
  [key: string]: unknown
}

export interface WalletMeta {
  id: WalletId
  name: string
  /** 图标兜底展示的字母 */
  short: string
  /** 图标兜底背景色（Tailwind 类，必须保持字面量以便被扫描） */
  colorClass: string
  downloadUrl?: string
}

export interface DiscoveredWallet extends WalletMeta {
  provider: Eip1193Provider
  /** EIP-6963 提供的 data URI 图标 */
  icon?: string
  rdns?: string
}

export class WalletError extends Error {
  readonly code: 'NO_PROVIDER' | 'NO_ACCOUNT' | 'REJECTED' | 'SIGN_FAILED' | 'UNKNOWN'

  constructor(
    code: WalletError['code'],
    message: string,
  ) {
    super(message)
    this.name = 'WalletError'
    this.code = code
  }
}

export const WALLET_CATALOG: Record<WalletId, WalletMeta> = {
  metamask: {
    id: 'metamask',
    name: 'MetaMask',
    short: 'M',
    colorClass: 'bg-[#f6851b]',
    downloadUrl: 'https://metamask.io/download/',
  },
  tokenpocket: {
    id: 'tokenpocket',
    name: 'TokenPocket',
    short: 'TP',
    colorClass: 'bg-[#2980fe]',
    downloadUrl: 'https://www.tokenpocket.pro/en/download/app',
  },
  okx: {
    id: 'okx',
    name: 'OKX Wallet',
    short: 'OKX',
    colorClass: 'bg-[#111111]',
    downloadUrl: 'https://www.okx.com/web3',
  },
  bitget: {
    id: 'bitget',
    name: 'Bitget Wallet',
    short: 'BG',
    colorClass: 'bg-[#00f0ff] text-[#0b1c2c]',
    downloadUrl: 'https://web3.bitget.com/',
  },
  trust: {
    id: 'trust',
    name: 'Trust Wallet',
    short: 'TW',
    colorClass: 'bg-[#3375bb]',
    downloadUrl: 'https://trustwallet.com/download',
  },
  imtoken: {
    id: 'imtoken',
    name: 'imToken',
    short: 'iT',
    colorClass: 'bg-[#11c4d1]',
    downloadUrl: 'https://token.im/download',
  },
  injected: {
    id: 'injected',
    name: '浏览器钱包',
    short: 'W',
    colorClass: 'bg-[#c9954a]',
  },
}

/** 供“未检测到钱包”弹层展示的安装列表 */
export const WALLET_LIST: WalletMeta[] = [
  WALLET_CATALOG.metamask,
  WALLET_CATALOG.tokenpocket,
  WALLET_CATALOG.okx,
  WALLET_CATALOG.bitget,
  WALLET_CATALOG.trust,
  WALLET_CATALOG.imtoken,
]

const WALLET_RANK: Record<WalletId, number> = {
  metamask: 0,
  tokenpocket: 1,
  okx: 2,
  bitget: 3,
  trust: 4,
  imtoken: 5,
  injected: 9,
}

const DEEP_LINKS: Partial<Record<WalletId, (url: string) => string>> = {
  metamask: (url) => `https://metamask.app.link/dapp/${url.replace(/^https?:\/\//, '')}`,
  tokenpocket: (url) =>
    `tpoutside://pull.activity?param=${encodeURIComponent(JSON.stringify({ url }))}`,
  okx: (url) => `okx://wallet/dapp/url?dappUrl=${encodeURIComponent(url)}`,
  bitget: (url) => `bitkeep://bkconnect?dappUrl=${encodeURIComponent(url)}`,
  trust: (url) => `trust://open_url?coin_id=60&url=${encodeURIComponent(url)}`,
  imtoken: (url) => `imtokenv2://navigate?action=openDapp&url=${encodeURIComponent(url)}`,
}

declare global {
  interface Window {
    ethereum?: Eip1193Provider
    okxwallet?: Eip1193Provider
    trustwallet?: Eip1193Provider
    imToken?: Eip1193Provider
    tokenpocket?: Eip1193Provider
    bitkeep?: { ethereum?: Eip1193Provider }
  }
}

function isProvider(value: unknown): value is Eip1193Provider {
  return (
    typeof value === 'object' &&
    value !== null &&
    typeof (value as { request?: unknown }).request === 'function'
  )
}

function asStringArray(value: unknown): string[] {
  return Array.isArray(value) ? value.filter((item): item is string => typeof item === 'string') : []
}

function asString(value: unknown): string {
  return typeof value === 'string' ? value : ''
}

function detectWalletId(provider: Eip1193Provider, rdns?: string): WalletId {
  if (provider.isTokenPocket) {
    return 'tokenpocket'
  }
  if (provider.isOkxWallet || provider.isOKExWallet) {
    return 'okx'
  }
  if (provider.isBitKeep || provider.isBitget) {
    return 'bitget'
  }
  if (provider.isTrust || provider.isTrustWallet) {
    return 'trust'
  }
  if (provider.isImToken || provider.isImtoken) {
    return 'imtoken'
  }
  if (provider.isMetaMask) {
    return 'metamask'
  }

  const key = (rdns || '').toLowerCase()
  if (key.includes('tokenpocket')) return 'tokenpocket'
  if (key.includes('okx') || key.includes('okex')) return 'okx'
  if (key.includes('bitget') || key.includes('bitkeep')) return 'bitget'
  if (key.includes('trust')) return 'trust'
  if (key.includes('imtoken')) return 'imtoken'
  if (key.includes('metamask')) return 'metamask'
  return 'injected'
}

interface Eip6963ProviderInfo {
  uuid: string
  name: string
  icon: string
  rdns: string
}

interface Eip6963Detail {
  info: Eip6963ProviderInfo
  provider: Eip1193Provider
}

/**
 * 发现可用钱包。
 * 同时扫描注入式全局变量与 EIP-6963 公告，按对象去重。
 */
export async function discoverWallets(timeoutMs = 500): Promise<DiscoveredWallet[]> {
  if (typeof window === 'undefined') {
    return []
  }

  const found = new Map<Eip1193Provider, DiscoveredWallet>()

  const addProvider = (
    provider: unknown,
    options: { rdns?: string; icon?: string; name?: string } = {},
  ): void => {
    if (!isProvider(provider)) {
      return
    }

    const existing = found.get(provider as Eip1193Provider)
    if (existing) {
      if (options.icon) existing.icon = options.icon
      if (options.name) existing.name = options.name
      if (options.rdns) existing.rdns = options.rdns
      return
    }

    const p = provider as Eip1193Provider
    const id = detectWalletId(p, options.rdns)
    const meta = WALLET_CATALOG[id]
    found.set(p, {
      ...meta,
      id,
      name: options.name || meta.name,
      icon: options.icon,
      rdns: options.rdns,
      provider: p,
    })
  }

  // 1. 传统注入式钱包
  const injected = window.ethereum
  if (injected) {
    if (Array.isArray(injected.providers) && injected.providers.length > 0) {
      for (const provider of injected.providers) {
        addProvider(provider)
      }
    } else {
      addProvider(injected)
    }
  }
  addProvider(window.okxwallet)
  addProvider(window.trustwallet)
  addProvider(window.imToken)
  addProvider(window.tokenpocket)
  addProvider(window.bitkeep?.ethereum)

  // 2. EIP-6963 多钱包标准
  const onAnnounce = (event: Event): void => {
    const detail = (event as CustomEvent<Eip6963Detail>).detail
    if (!detail?.provider) {
      return
    }
    addProvider(detail.provider, {
      rdns: detail.info?.rdns,
      icon: detail.info?.icon,
      name: detail.info?.name,
    })
  }

  window.addEventListener('eip6963:announceProvider', onAnnounce)
  window.dispatchEvent(new Event('eip6963:requestProvider'))
  await new Promise<void>((resolve) => {
    window.setTimeout(resolve, timeoutMs)
  })
  window.removeEventListener('eip6963:announceProvider', onAnnounce)

  return Array.from(found.values()).sort((a, b) => WALLET_RANK[a.id] - WALLET_RANK[b.id])
}

/** 直接读取已授权账户，不弹窗（用于静默恢复会话） */
export async function getAuthorizedAccounts(provider: Eip1193Provider): Promise<string[]> {
  try {
    const accounts = asStringArray(await provider.request({ method: 'eth_accounts' }))
    return accounts.map((account) => account.toLowerCase())
  } catch {
    return []
  }
}

/** 主动请求账户授权，会拉起钱包 */
export async function requestAccount(wallet: DiscoveredWallet): Promise<string> {
  try {
    const accounts = asStringArray(
      await wallet.provider.request({ method: 'eth_requestAccounts' }),
    )
    const account = accounts[0]?.toLowerCase()
    if (!account) {
      throw new WalletError('NO_ACCOUNT', '未获取到钱包地址，请在钱包中创建或导入账户')
    }
    return account
  } catch (error) {
    throw normalizeWalletError(error, '连接钱包失败')
  }
}

export async function getChainId(provider: Eip1193Provider): Promise<string> {
  try {
    return asString(await provider.request({ method: 'eth_chainId' }))
  } catch {
    return ''
  }
}

function toHexMessage(message: string): string {
  const bytes = new TextEncoder().encode(message)
  return `0x${Array.from(bytes)
    .map((byte) => byte.toString(16).padStart(2, '0'))
    .join('')}`
}

/** personal_sign 签名，消息按 EIP-191 以十六进制传入，兼容主流钱包 */
export async function signMessage(
  wallet: DiscoveredWallet,
  address: string,
  message: string,
): Promise<string> {
  try {
    const signature = await wallet.provider.request({
      method: 'personal_sign',
      params: [toHexMessage(message), address],
    })
    if (typeof signature !== 'string' || !signature) {
      throw new WalletError('SIGN_FAILED', '钱包未返回签名结果')
    }
    return normalizeSignature(signature)
  } catch (error) {
    throw normalizeWalletError(error, '签名失败')
  }
}

/**
 * 部分钱包返回 v=0/1，统一成以太坊常见的 27/28，兼容后端签名恢复。
 */
function normalizeSignature(signature: string): string {
  const hex = signature.startsWith('0x') ? signature.slice(2) : signature
  if (hex.length !== 130) {
    return signature
  }

  const v = Number.parseInt(hex.slice(128, 130), 16)
  if (v === 0 || v === 1) {
    return `0x${hex.slice(0, 128)}${(v + 27).toString(16).padStart(2, '0')}`
  }
  return signature.startsWith('0x') ? signature : `0x${signature}`
}

export interface WalletEventHandlers {
  accountsChanged?: (accounts: string[]) => void
  chainChanged?: (chainId: string) => void
  disconnect?: () => void
}

/** 监听钱包账户/网络变化，返回解绑函数 */
export function bindWalletEvents(
  provider: Eip1193Provider,
  handlers: WalletEventHandlers,
): () => void {
  const onAccountsChanged: ProviderListener = (...args) => {
    handlers.accountsChanged?.(asStringArray(args[0]).map((account) => account.toLowerCase()))
  }
  const onChainChanged: ProviderListener = (...args) => {
    handlers.chainChanged?.(asString(args[0]))
  }
  const onDisconnect: ProviderListener = () => {
    handlers.disconnect?.()
  }

  provider.on?.('accountsChanged', onAccountsChanged)
  provider.on?.('chainChanged', onChainChanged)
  provider.on?.('disconnect', onDisconnect)

  return () => {
    provider.removeListener?.('accountsChanged', onAccountsChanged)
    provider.removeListener?.('chainChanged', onChainChanged)
    provider.removeListener?.('disconnect', onDisconnect)
  }
}

/** 判断是否移动端（用于选择 deep link 还是下载页） */
export function isMobileBrowser(): boolean {
  if (typeof navigator === 'undefined') {
    return false
  }
  return /Android|iPhone|iPad|iPod|Mobile/i.test(navigator.userAgent)
}

/** 在未检测到钱包时，尝试用 deep link 唤起钱包 App，否则打开下载页 */
export function openWallet(id: WalletId): void {
  if (typeof window === 'undefined') {
    return
  }

  const builder = DEEP_LINKS[id]
  if (builder && isMobileBrowser()) {
    window.location.href = builder(window.location.href)
    return
  }

  const downloadUrl = WALLET_CATALOG[id].downloadUrl
  if (downloadUrl) {
    window.open(downloadUrl, '_blank', 'noopener,noreferrer')
  }
}

/** 将钱包/底层错误统一转换为可展示的 WalletError */
export function normalizeWalletError(error: unknown, fallback = '操作失败'): WalletError {
  if (error instanceof WalletError) {
    return error
  }

  const code =
    typeof (error as { code?: unknown })?.code === 'number'
      ? (error as { code: number }).code
      : undefined
  const rawMessage = typeof (error as { message?: unknown })?.message === 'string'
    ? (error as { message: string }).message
    : ''

  if (
    code === 4001 ||
    /reject|denied|cancel|用户取消|用户拒绝/i.test(rawMessage)
  ) {
    return new WalletError('REJECTED', '您已取消钱包授权')
  }

  if (/no provider|provider not found|not installed/i.test(rawMessage)) {
    return new WalletError('NO_PROVIDER', '未检测到 Web3 钱包，请使用钱包 App 打开')
  }

  return new WalletError('UNKNOWN', rawMessage || fallback)
}

/** 钱包链配置（与后台 /chain/config 对应） */
export interface WalletChainConfig {
  chainId: number
  chainName: string
  nativeCurrency: {
    name: string
    symbol: string
    decimals: number
  }
  rpcUrls: string[]
  blockExplorerUrls: string[]
}

/** 确保钱包切换到目标链，钱包没有该链时自动添加 */
export async function ensureWalletChain(
  provider: Eip1193Provider,
  chain: WalletChainConfig,
): Promise<void> {
  const chainId = `0x${chain.chainId.toString(16)}`

  try {
    await provider.request({
      method: 'wallet_switchEthereumChain',
      params: [{ chainId }],
    })
    return
  } catch (error) {
    const raw = error as { code?: number; message?: string; data?: { originalError?: { code?: number } } }
    const code = raw?.code ?? raw?.data?.originalError?.code
    const message = String(raw?.message || '')

    // 4902：钱包没有该链，需要先 addEthereumChain
    if (code !== 4902 && !/unrecognized chain|try adding the chain|not added/i.test(message)) {
      throw normalizeWalletError(error, '切换钱包网络失败')
    }
  }

  try {
    await provider.request({
      method: 'wallet_addEthereumChain',
      params: [
        {
          chainId,
          chainName: chain.chainName,
          nativeCurrency: chain.nativeCurrency,
          rpcUrls: chain.rpcUrls,
          blockExplorerUrls: chain.blockExplorerUrls,
        },
      ],
    })
  } catch (error) {
    throw normalizeWalletError(error, '添加钱包网络失败')
  }
}

/** 将十进制金额字符串转为链上最小单位 */
export function toBaseUnits(amount: string, decimals: number): string {
  const [intPart = '0', fracPart = ''] = amount.split('.')
  const frac = (fracPart + '0'.repeat(decimals)).slice(0, decimals)
  const value = `${intPart}${frac}`.replace(/^0+(?=\d)/, '')
  return value || '0'
}

/** 发起 ERC20 transfer，返回 txHash */
export async function sendErc20Transfer(
  provider: Eip1193Provider,
  from: string,
  contract: string,
  to: string,
  amount: string,
  decimals: number,
): Promise<string> {
  const toWord = to.toLowerCase().replace(/^0x/, '').padStart(64, '0')
  // toBaseUnits 返回十进制字符串，calldata 必须使用十六进制
  const amountWord = BigInt(toBaseUnits(amount, decimals)).toString(16).padStart(64, '0')
  const data = `0xa9059cbb${toWord}${amountWord}`

  try {
    const txHash = await provider.request({
      method: 'eth_sendTransaction',
      params: [
        {
          from,
          to: contract,
          value: '0x0',
          data,
        },
      ],
    })
    if (typeof txHash !== 'string' || !txHash) {
      throw new WalletError('UNKNOWN', '钱包未返回交易哈希')
    }
    return txHash
  } catch (error) {
    throw normalizeWalletError(error, '发起支付失败')
  }
}

/** 查询 ERC20 余额，返回最小单位 */
export async function getErc20Balance(
  provider: Eip1193Provider,
  walletAddress: string,
  contract: string,
): Promise<bigint> {
  const addressWord = walletAddress.toLowerCase().replace(/^0x/, '').padStart(64, '0')
  const data = `0x70a08231${addressWord}`

  try {
    const result = await provider.request({
      method: 'eth_call',
      params: [{ to: contract, data }, 'latest'],
    })
    if (typeof result !== 'string' || !result.startsWith('0x')) {
      throw new WalletError('UNKNOWN', '查询 USDT 余额失败')
    }
    if (result === '0x') {
      return 0n
    }
    return BigInt(result)
  } catch (error) {
    throw normalizeWalletError(error, '查询 USDT 余额失败')
  }
}

/** 将最小单位格式化为固定小数位，不经过 Number，避免精度丢失 */
export function formatTokenAmount(raw: bigint, decimals: number, digits = 2): string {
  if (decimals <= 0) {
    return raw.toString()
  }

  const base = 10n ** BigInt(decimals)
  const integer = raw / base
  const fraction = (raw % base).toString().padStart(decimals, '0').slice(0, digits).padEnd(digits, '0')
  return `${integer.toString()}.${fraction}`
}

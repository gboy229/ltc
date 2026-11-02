import { computed, ref, shallowRef } from 'vue'
import { defineStore } from 'pinia'

import { getAuthStatus, getNonce, getProfile, loginWithWallet, logout as logoutApi } from '@/api/users'
import { shortenAddress } from '@/utils/address'
import {
  getAddress,
  getInviteUid,
  getToken,
  getWalletId,
  removeAddress,
  removeToken,
  removeWalletId,
  setAddress,
  setToken,
  setWalletId,
} from '@/utils/auth'
import { ApiError, UNAUTHORIZED_EVENT } from '@/utils/http'
import { toast } from '@/utils/toast'
import {
  bindWalletEvents,
  discoverWallets,
  getAuthorizedAccounts,
  getChainId,
  normalizeWalletError,
  requestAccount,
  signMessage,
  type DiscoveredWallet,
  type Eip1193Provider,
  type WalletId,
} from '@/web3/wallet'

export const useUserStore = defineStore('user', () => {
  const token = ref(getToken())
  const address = ref(getAddress())
  const uid = ref('')
  const nickname = ref('')
  const avatar = ref('')
  const levelName = ref('')
  const walletId = ref<WalletId | ''>((getWalletId() as WalletId) || '')
  const walletName = ref('')
  const chainId = ref('')

  const wallets = ref<DiscoveredWallet[]>([])
  const currentWallet = shallowRef<DiscoveredWallet | null>(null)
  const provider = shallowRef<Eip1193Provider | null>(null)

  const connecting = ref(false)
  const initialized = ref(false)
  const sheetVisible = ref(false)

  type RegisterStep = 'idle' | 'confirm' | 'bind'
  const registerVisible = ref(false)
  const registerStep = ref<RegisterStep>('idle')
  const registerInviteUid = ref('')
  const registerSubmitting = ref(false)
  const registerError = ref('')

  const isLoggedIn = computed(() => Boolean(token.value && address.value))
  const shortAddress = computed(() => (address.value ? shortenAddress(address.value) : ''))
  const hasProvider = computed(() => Boolean(provider.value))

  let authenticating = false
  let bootstrapPromise: Promise<void> | null = null
  let unbindEvents: (() => void) | null = null
  let pendingWallet: DiscoveredWallet | null = null
  let pendingAccount = ''

  /** 接口/网络错误已由 http 拦截器提示，这里避免重复 toast */
  function isRequestError(error: unknown): boolean {
    if (error instanceof ApiError) {
      return true
    }
    return Boolean((error as { isAxiosError?: boolean } | null)?.isAxiosError)
  }

  function expireSession(): void {
    token.value = ''
    uid.value = ''
    nickname.value = ''
    avatar.value = ''
    levelName.value = ''
    removeToken()
    resetRegisterState()
  }

  function disconnectLocal(): void {
    unbindEvents?.()
    unbindEvents = null
    currentWallet.value = null
    provider.value = null

    token.value = ''
    address.value = ''
    uid.value = ''
    nickname.value = ''
    avatar.value = ''
    levelName.value = ''
    chainId.value = ''

    removeToken()
    removeAddress()
    removeWalletId()
    resetRegisterState()
  }

  function attachWallet(wallet: DiscoveredWallet): void {
    currentWallet.value = wallet
    provider.value = wallet.provider
    walletId.value = wallet.id
    walletName.value = wallet.name
    setWalletId(wallet.id)

    unbindEvents?.()
    unbindEvents = bindWalletEvents(wallet.provider, {
      accountsChanged: (accounts) => {
        const next = accounts[0] ?? ''
        if (!next) {
          disconnectLocal()
          toast('钱包已断开连接', 'info')
          return
        }
        if (next === address.value) {
          return
        }

        address.value = next
        setAddress(next)
        expireSession()
        toast('钱包账户已切换，请重新签名登录', 'info')
        if (currentWallet.value) {
          void authenticate(currentWallet.value)
        }
      },
      chainChanged: (id) => {
        chainId.value = id
      },
      disconnect: () => {
        disconnectLocal()
        toast('钱包已断开连接', 'info')
      },
    })
  }

  async function loadProfile(silent: boolean): Promise<boolean> {
    try {
      const { user } = await getProfile(silent)
      uid.value = user.uid || ''
      nickname.value = user.nickname || ''
      avatar.value = user.avatar || ''
      levelName.value = user.levelName || ''
      if (user.address) {
        address.value = user.address.toLowerCase()
        setAddress(user.address)
      }
      return true
    } catch (error) {
      if (error instanceof ApiError && (error.code === 61 || error.code === 401)) {
        expireSession()
      }
      return false
    }
  }

  function resetRegisterState(): void {
    registerVisible.value = false
    registerStep.value = 'idle'
    registerInviteUid.value = ''
    registerSubmitting.value = false
    registerError.value = ''
    pendingWallet = null
    pendingAccount = ''
  }

  function openRegister(forceInvite: boolean): void {
    const inviteUid = getInviteUid()
    registerInviteUid.value = inviteUid
    registerStep.value = inviteUid || forceInvite ? 'bind' : 'confirm'
    registerError.value = ''
    registerVisible.value = true
  }

  async function completeWalletLogin(
    wallet: DiscoveredWallet,
    account: string,
    inviteUid: string,
  ): Promise<boolean> {
    const { nonce } = await getNonce(account)
    const signature = await signMessage(wallet, account, nonce)
    const result = await loginWithWallet({
      type: 'wallet',
      account,
      signature,
      nonce,
      uid: inviteUid || undefined,
    })

    token.value = result.token
    setToken(result.token)
    address.value = account
    setAddress(account)
    chainId.value = await getChainId(wallet.provider)

    const ok = await loadProfile(false)
    if (ok) {
      toast('钱包登录成功', 'success')
    }
    return ok
  }

  function registerErrorMessage(error: unknown): string {
    if (error instanceof Error && error.message) {
      return error.message
    }
    return '注册失败，请重试'
  }

  async function submitRegistration(inviteUid: string): Promise<boolean> {
    if (!pendingWallet || !pendingAccount || registerSubmitting.value) {
      return false
    }

    registerSubmitting.value = true
    registerError.value = ''
    try {
      const ok = await completeWalletLogin(pendingWallet, pendingAccount, inviteUid)
      if (!ok) {
        registerError.value = '注册成功但获取用户信息失败，请重试'
        return false
      }
      resetRegisterState()
      return true
    } catch (error) {
      registerError.value = registerErrorMessage(error)
      return false
    } finally {
      registerSubmitting.value = false
    }
  }

  async function confirmRegister(): Promise<boolean> {
    return submitRegistration('')
  }

  async function confirmBindInvite(): Promise<boolean> {
    const uid = registerInviteUid.value.trim()
    if (!uid) {
      registerError.value = '请输入邀请码'
      return false
    }
    return submitRegistration(uid)
  }

  function goBind(): void {
    registerStep.value = 'bind'
    registerError.value = ''
    registerInviteUid.value = ''
  }

  function cancelRegister(): void {
    resetRegisterState()
  }

  async function authenticate(wallet: DiscoveredWallet): Promise<boolean> {
    if (authenticating) {
      return false
    }

    authenticating = true
    connecting.value = true
    resetRegisterState()

    try {
      attachWallet(wallet)

      const account = await requestAccount(wallet)
      const status = await getAuthStatus({ type: 'wallet', account })

      // 已注册或系统第一个账号，直接走原登录/注册逻辑，不弹注册框
      if (status.registered || status.isFirstUser) {
        return await completeWalletLogin(wallet, account, '')
      }

      pendingWallet = wallet
      pendingAccount = account
      openRegister(status.forceInvite)
      return false
    } catch (error) {
      // 接口错误已由 http 拦截器统一提示，避免重复弹 toast
      if (!isRequestError(error)) {
        const walletError = normalizeWalletError(error, '钱包登录失败')
        toast(walletError.message, walletError.code === 'REJECTED' ? 'info' : 'error')
      }
      return false
    } finally {
      authenticating = false
      connecting.value = false
      initialized.value = true
    }
  }

  /** 对敏感操作进行钱包签名，返回后端校验所需的 nonce 和签名 */
  async function signAction(): Promise<{ nonce: string; signature: string }> {
    if (!currentWallet.value || !address.value) {
      throw new Error('请先连接钱包')
    }
    const { nonce } = await getNonce(address.value, 'action')
    const signature = await signMessage(currentWallet.value, address.value, nonce)
    return { nonce, signature }
  }

  async function refreshWallets(): Promise<DiscoveredWallet[]> {
    try {
      wallets.value = await discoverWallets()
    } catch {
      wallets.value = []
    }
    return wallets.value
  }

  /** 重新读取钱包当前链 ID（部分钱包切换网络后不会派发 chainChanged） */
  async function syncChainId(): Promise<string> {
    if (!provider.value) {
      chainId.value = ''
      return ''
    }
    chainId.value = await getChainId(provider.value)
    return chainId.value
  }

  /**
   * 应用启动时自动连接钱包并登录。
   * - 已有有效会话：静默恢复，不弹签名
   * - 无会话：自动拉起钱包请求授权 + 签名登录
   * - 无钱包环境：展示钱包选择/下载引导
   */
  async function bootstrap(): Promise<void> {
    if (initialized.value) {
      return
    }
    if (bootstrapPromise) {
      return bootstrapPromise
    }

    bootstrapPromise = (async () => {
      const savedToken = token.value
      const savedAddress = address.value
      const available = await refreshWallets()

      const target =
        available.find((wallet) => wallet.id === walletId.value) ?? available[0] ?? null

      if (!target) {
        if (savedToken && savedAddress) {
          await loadProfile(true)
        } else {
          sheetVisible.value = true
        }
        initialized.value = true
        return
      }

      if (savedToken && savedAddress) {
        const accounts = await getAuthorizedAccounts(target.provider)
        if (accounts[0] === savedAddress) {
          attachWallet(target)
          chainId.value = await getChainId(target.provider)
          const ok = await loadProfile(true)
          if (ok) {
            initialized.value = true
            return
          }
        }
      }

      await authenticate(target)
      initialized.value = true
    })().finally(() => {
      bootstrapPromise = null
    })

    return bootstrapPromise
  }

  /** 手动连接：指定钱包时直接连接，否则自动选择/弹出选择面板 */
  async function connect(wallet?: DiscoveredWallet): Promise<boolean> {
    if (wallet) {
      sheetVisible.value = false
      if (wallet.id === walletId.value && isLoggedIn.value && provider.value === wallet.provider) {
        return true
      }
      return authenticate(wallet)
    }

    if (connecting.value) {
      return false
    }

    if (isLoggedIn.value && provider.value) {
      return true
    }

    const available = await refreshWallets()
    if (available.length === 1) {
      return authenticate(available[0])
    }

    sheetVisible.value = true
    return false
  }

  /** 确保已登录，供下单等需要身份的操作调用 */
  async function ensureLogin(): Promise<boolean> {
    if (isLoggedIn.value) {
      return true
    }
    await connect()
    return isLoggedIn.value
  }

  function openSheet(): void {
    sheetVisible.value = true
    if (!wallets.value.length) {
      void refreshWallets()
    }
  }

  function closeSheet(): void {
    sheetVisible.value = false
  }

  /** 切换钱包：退出当前会话，清空本地身份，重新选择钱包 */
  async function switchWallet(): Promise<void> {
    try {
      await logoutApi()
    } catch {
      // 服务端注销失败不影响本地切换
    }

    // 尽量让 MetaMask 等钱包重新弹出账户选择
    if (provider.value) {
      try {
        await provider.value.request({
          method: 'wallet_revokePermissions',
          params: [{ eth_accounts: {} }],
        })
      } catch {
        // 部分钱包不支持 revoke，用户可在钱包内手动切换账户
      }
    }

    disconnectLocal()
    await refreshWallets()
    sheetVisible.value = true
    toast('请选择要连接的钱包', 'info')
  }

  async function logout(): Promise<void> {
    try {
      await logoutApi()
    } catch {
      // 服务端注销失败不影响本地退出
    }
    disconnectLocal()
    toast('已退出登录', 'success')
  }

  if (typeof window !== 'undefined') {
    window.addEventListener(UNAUTHORIZED_EVENT, () => {
      expireSession()
    })
  }

  return {
    token,
    address,
    uid,
    nickname,
    avatar,
    levelName,
    walletId,
    walletName,
    chainId,
    wallets,
    provider,
    connecting,
    initialized,
    sheetVisible,
    isLoggedIn,
    shortAddress,
    hasProvider,
    registerVisible,
    registerStep,
    registerInviteUid,
    registerSubmitting,
    registerError,
    confirmRegister,
    confirmBindInvite,
    goBind,
    cancelRegister,
    bootstrap,
    connect,
    ensureLogin,
    signAction,
    refreshWallets,
    syncChainId,
    openSheet,
    closeSheet,
    switchWallet,
    logout,
  }
})

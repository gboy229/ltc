import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { getChainConfig, type ChainConfig, type ChainContract } from '@/api/chain'

export const useChainStore = defineStore('chain', () => {
  const chains = ref<ChainConfig[]>([])
  const loading = ref(false)
  const loaded = ref(false)
  const error = ref('')

  let inflight: Promise<ChainConfig[]> | null = null

  const defaultChain = computed<ChainConfig | null>(() => chains.value[0] ?? null)

  function getContract(chain: ChainConfig | null): ChainContract | null {
    if (!chain) {
      return null
    }
    return (
      chain.contracts.find((item) => item.symbol.toUpperCase() === 'USDT') ??
      chain.contracts[0] ??
      null
    )
  }

  const defaultContract = computed<ChainContract | null>(() => getContract(defaultChain.value))

  /** 应用启动时加载一次，后续页面直接复用 */
  async function load(force = false): Promise<ChainConfig[]> {
    if (loaded.value && !force) {
      return chains.value
    }
    if (inflight) {
      return inflight
    }

    loading.value = true
    error.value = ''

    inflight = getChainConfig()
      .then((res) => {
        chains.value = res.chains ?? []
        loaded.value = chains.value.length > 0
        return chains.value
      })
      .catch((err) => {
        error.value = err instanceof Error ? err.message : '加载支付链失败'
        return []
      })
      .finally(() => {
        loading.value = false
        inflight = null
      })

    return inflight
  }

  return {
    chains,
    loading,
    loaded,
    error,
    defaultChain,
    defaultContract,
    getContract,
    load,
  }
})

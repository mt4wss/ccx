import { ref } from 'vue'
import { useAdminApi } from '@/composables/useAdminApi'

// Module-level singletons
const fuzzyMode = ref(false)
const stripBillingHeader = ref(false)
const fuzzyLoading = ref(false)
const stripLoading = ref(false)

export function useConsoleSettings() {
  const api = useAdminApi()

  async function fetchFuzzyMode() {
    const data = await api.get<{ enabled: boolean }>('/api/settings/fuzzy-mode')
    fuzzyMode.value = data.enabled
  }

  async function setFuzzyMode(enabled: boolean) {
    fuzzyLoading.value = true
    try {
      await api.put('/api/settings/fuzzy-mode', { enabled })
      fuzzyMode.value = enabled
    } finally {
      fuzzyLoading.value = false
    }
  }

  async function fetchStripBillingHeader() {
    const data = await api.get<{ enabled: boolean }>('/api/settings/strip-billing-header')
    stripBillingHeader.value = data.enabled
  }

  async function setStripBillingHeader(enabled: boolean) {
    stripLoading.value = true
    try {
      await api.put('/api/settings/strip-billing-header', { enabled })
      stripBillingHeader.value = enabled
    } finally {
      stripLoading.value = false
    }
  }

  return {
    fuzzyMode,
    stripBillingHeader,
    fuzzyLoading,
    stripLoading,
    fetchFuzzyMode,
    setFuzzyMode,
    fetchStripBillingHeader,
    setStripBillingHeader,
  }
}

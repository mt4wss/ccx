import { ref } from 'vue'
import { useAdminApi } from '@/composables/useAdminApi'
import type {
  CapabilitySnapshot,
  CapabilityTestJob,
  CapabilityTestJobStartResponse,
} from '@/services/admin-api'

// Module-level singletons
const activeJob = ref<CapabilityTestJob | null>(null)
const snapshot = ref<CapabilitySnapshot | null>(null)
const loading = ref(false)
const polling = ref(false)
const error = ref('')
let pollTimer: ReturnType<typeof setInterval> | undefined

export function useCapabilityTests() {
  const api = useAdminApi()

  function clearError() {
    error.value = ''
  }

  async function startTest(
    channelType: string,
    channelId: number,
    options?: { targetProtocols?: string[]; models?: string[]; rpm?: number },
  ) {
    loading.value = true
    clearError()
    try {
      const resp = await api.post<CapabilityTestJobStartResponse>(
        `/api/${channelType}/channels/${channelId}/capability-test`,
        options,
      )
      activeJob.value = resp.job || null
      // 启动轮询
      if (resp.jobId) {
        startPolling(channelType, channelId, resp.jobId)
      }
      return resp
    } catch (e) {
      error.value = e instanceof Error ? e.message : String(e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function fetchSnapshot(channelType: string, channelId: number) {
    try {
      snapshot.value = await api.get<CapabilitySnapshot>(
        `/api/${channelType}/channels/${channelId}/capability-snapshot`,
      )
    } catch {
      // snapshot 可能不存在，静默
    }
  }

  async function fetchJobStatus(channelType: string, channelId: number, jobId: string) {
    const job = await api.get<CapabilityTestJob>(
      `/api/${channelType}/channels/${channelId}/capability-test/${jobId}`,
    )
    activeJob.value = job
    // 如果已完成或失败，停止轮询
    if (job.status === 'completed' || job.status === 'failed' || job.status === 'cancelled') {
      stopPolling()
    }
    return job
  }

  function startPolling(channelType: string, channelId: number, jobId: string) {
    stopPolling()
    polling.value = true
    pollTimer = setInterval(async () => {
      try {
        await fetchJobStatus(channelType, channelId, jobId)
      } catch (e) {
        stopPolling()
        error.value = e instanceof Error ? e.message : String(e)
      }
    }, 2000)
  }

  function stopPolling() {
    polling.value = false
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = undefined
    }
  }

  async function cancelTest(channelType: string, channelId: number, jobId: string) {
    await api.del(`/api/${channelType}/channels/${channelId}/capability-test/${jobId}`)
    stopPolling()
    if (activeJob.value) {
      activeJob.value.status = 'cancelled'
      activeJob.value.lifecycle = 'cancelled'
    }
  }

  async function retryModel(channelType: string, channelId: number, jobId: string) {
    await api.post(`/api/${channelType}/channels/${channelId}/capability-test/${jobId}/retry`)
    // 重新开始轮询
    startPolling(channelType, channelId, jobId)
  }

  function reset() {
    stopPolling()
    activeJob.value = null
    snapshot.value = null
    error.value = ''
    loading.value = false
  }

  return {
    activeJob,
    snapshot,
    loading,
    polling,
    error,
    clearError,
    startTest,
    fetchSnapshot,
    fetchJobStatus,
    cancelTest,
    retryModel,
    reset,
  }
}

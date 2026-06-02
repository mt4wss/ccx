<script setup lang="ts">
import { ref, computed, watch, onBeforeUnmount } from 'vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Skeleton } from '@/components/ui/skeleton'
import { Loader2, X, Play, Square, RefreshCw, CheckCircle2, XCircle, Clock } from 'lucide-vue-next'
import { useCapabilityTests } from '@/composables/useCapabilityTests'
import { useLanguage } from '@/composables/useLanguage'
import type {
  CapabilityTestJob,
  CapabilitySnapshot,
  CapabilityProtocolJobResult,
  CapabilityModelJobResult,
} from '@/services/admin-api'

interface Props {
  open: boolean
  channelType: string
  channelId: number
  channelName: string
}

const props = defineProps<Props>()
const emit = defineEmits<{ (e: 'close'): void }>()

const { tf } = useLanguage()
const {
  activeJob,
  snapshot,
  loading,
  polling,
  error,
  startTest,
  fetchSnapshot,
  cancelTest,
  retryModel,
  reset,
} = useCapabilityTests()

const isStarting = ref(false)

// 加载 snapshot
watch(() => props.open, async (isOpen) => {
  if (isOpen) {
    await fetchSnapshot(props.channelType, props.channelId)
  } else {
    reset()
  }
})

const currentJob = computed(() => activeJob.value)
const progress = computed(() => currentJob.value?.progress)

const progressPercent = computed(() => {
  if (!progress.value || !progress.value.totalModels) return 0
  return Math.round((progress.value.completedModels / progress.value.totalModels) * 100)
})

const isActive = computed(() => {
  const s = currentJob.value?.status
  return s === 'running' || s === 'queued'
})

async function handleStart() {
  isStarting.value = true
  try {
    await startTest(props.channelType, props.channelId)
  } finally {
    isStarting.value = false
  }
}

async function handleCancel() {
  if (!currentJob.value) return
  await cancelTest(props.channelType, props.channelId, currentJob.value.jobId)
}

async function handleRetryModel(model: string) {
  if (!currentJob.value) return
  await retryModel(props.channelType, props.channelId, currentJob.value.jobId)
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}

function getStatusIcon(status: string) {
  if (status === 'success') return CheckCircle2
  if (status === 'failed') return XCircle
  if (status === 'running' || status === 'queued') return Loader2
  return Clock
}

function getStatusColor(status: string) {
  if (status === 'success') return 'text-emerald-500'
  if (status === 'failed') return 'text-rose-500'
  if (status === 'running') return 'text-blue-500'
  return 'text-muted-foreground'
}

// 从 snapshot 或 job 中获取协议结果
const protocolResults = computed((): CapabilityProtocolJobResult[] => {
  if (currentJob.value?.tests?.length) return currentJob.value.tests
  if (snapshot.value?.tests?.length) return snapshot.value.tests
  return []
})

onBeforeUnmount(() => {
  reset()
})
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center"
        @keydown="onKeyDown"
      >
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="emit('close')" />

        <div class="relative z-10 w-[90vw] max-w-3xl max-h-[85vh] rounded-2xl border border-border bg-card shadow-2xl flex flex-col">
          <!-- Header -->
          <div class="flex items-center justify-between p-4 border-b border-border shrink-0">
            <div class="flex items-center gap-2">
              <Play class="w-4 h-4 text-primary" />
              <h3 class="text-sm font-semibold">
                {{ tf('console.capability.title', '能力测试') }}: {{ channelName }}
              </h3>
              <Badge v-if="currentJob?.runMode && currentJob.runMode !== 'fresh'" variant="secondary" class="text-[10px]">
                {{ currentJob.runMode }}
              </Badge>
            </div>
            <Button variant="ghost" size="icon-sm" @click="emit('close')">
              <X class="w-4 h-4" />
            </Button>
          </div>

          <!-- Body -->
          <div class="flex-1 min-h-0 overflow-hidden p-4 space-y-4">
            <!-- Error -->
            <div v-if="error" class="text-sm text-destructive bg-destructive/10 rounded-lg p-3">
              {{ error }}
            </div>

            <!-- 无任务状态 -->
            <div v-if="!currentJob && !loading" class="text-center py-8 space-y-4">
              <div v-if="protocolResults.length > 0">
                <p class="text-sm text-muted-foreground mb-4">
                  {{ tf('console.capability.lastResults', '上次测试结果') }}
                </p>
              </div>
              <div v-else>
                <p class="text-sm text-muted-foreground mb-4">
                  {{ tf('console.capability.noResults', '尚未进行能力测试') }}
                </p>
              </div>
              <Button :disabled="isStarting" @click="handleStart">
                <Loader2 v-if="isStarting" class="w-4 h-4 mr-2 animate-spin" />
                <Play v-else class="w-4 h-4 mr-2" />
                {{ tf('console.capability.start', '开始测试') }}
              </Button>
            </div>

            <!-- 进行中状态 -->
            <div v-if="isActive && currentJob" class="space-y-3">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <Loader2 class="w-4 h-4 animate-spin text-primary" />
                  <span class="text-sm font-medium">
                    {{ tf('console.capability.running', '测试进行中...') }}
                  </span>
                </div>
                <Button variant="destructive" size="sm" @click="handleCancel">
                  <Square class="w-3.5 h-3.5 mr-1.5" />
                  {{ tf('console.capability.cancel', '取消') }}
                </Button>
              </div>
              <Progress :model-value="progressPercent" />
              <p class="text-xs text-muted-foreground text-center">
                {{ progress?.completedModels || 0 }} / {{ progress?.totalModels || 0 }} {{ tf('console.capability.models', '模型') }}
              </p>
            </div>

            <!-- 结果列表 -->
            <div v-if="protocolResults.length > 0" class="space-y-3">
              <h4 class="text-sm font-medium">
                {{ tf('console.capability.protocolResults', '协议结果') }}
              </h4>
              <ScrollArea class="max-h-[40vh]">
                <div class="space-y-2 pr-4">
                  <div
                    v-for="proto in protocolResults"
                    :key="proto.protocol"
                    class="border border-border rounded-lg p-3 space-y-2"
                  >
                    <div class="flex items-center justify-between">
                      <div class="flex items-center gap-2">
                        <component
                          :is="getStatusIcon(proto.status)"
                          class="w-4 h-4"
                          :class="[getStatusColor(proto.status), { 'animate-spin': proto.status === 'running' }]"
                        />
                        <span class="text-sm font-medium">{{ proto.protocol }}</span>
                        <Badge
                          :variant="proto.success ? 'default' : 'destructive'"
                          class="text-[10px]"
                        >
                          {{ proto.status }}
                        </Badge>
                      </div>
                      <div class="flex items-center gap-2 text-xs text-muted-foreground">
                        <span v-if="proto.latency">{{ proto.latency }}ms</span>
                        <span v-if="proto.testedModel">model: {{ proto.testedModel }}</span>
                        <Badge v-if="proto.streamingSupported" variant="secondary" class="text-[10px]">
                          SSE
                        </Badge>
                      </div>
                    </div>

                    <!-- 模型结果 -->
                    <div v-if="proto.modelResults?.length" class="space-y-1 pl-6">
                      <div
                        v-for="model in proto.modelResults"
                        :key="model.model"
                        class="flex items-center justify-between text-xs py-1"
                      >
                        <div class="flex items-center gap-1.5">
                          <component
                            :is="getStatusIcon(model.status)"
                            class="w-3 h-3"
                            :class="[getStatusColor(model.status), { 'animate-spin': model.status === 'running' }]"
                          />
                          <span class="font-mono">{{ model.model }}</span>
                          <span v-if="model.actualModel && model.actualModel !== model.model" class="text-muted-foreground">
                            → {{ model.actualModel }}
                          </span>
                        </div>
                        <div class="flex items-center gap-2">
                          <span v-if="model.latency" class="text-muted-foreground">{{ model.latency }}ms</span>
                          <span v-if="model.error" class="text-destructive truncate max-w-[200px]" :title="model.error">
                            {{ model.error }}
                          </span>
                          <Button
                            v-if="model.status === 'failed'"
                            variant="ghost"
                            size="icon-sm"
                            @click="handleRetryModel(model.model)"
                          >
                            <RefreshCw class="w-3 h-3" />
                          </Button>
                        </div>
                      </div>
                    </div>

                    <!-- 协议级错误 -->
                    <p v-if="proto.error" class="text-xs text-destructive pl-6">
                      {{ proto.error }}
                    </p>
                  </div>
                </div>
              </ScrollArea>
            </div>

            <!-- 兼容协议 -->
            <div v-if="currentJob?.compatibleProtocols?.length || snapshot?.compatibleProtocols?.length" class="flex items-center gap-2 flex-wrap">
              <span class="text-xs text-muted-foreground">
                {{ tf('console.capability.compatible', '兼容协议') }}:
              </span>
              <Badge
                v-for="proto in (currentJob?.compatibleProtocols || snapshot?.compatibleProtocols || [])"
                :key="proto"
                variant="secondary"
                class="text-[10px]"
              >
                {{ proto }}
              </Badge>
            </div>

            <!-- 总耗时 -->
            <div v-if="currentJob?.totalDuration || snapshot?.totalDuration" class="text-xs text-muted-foreground text-right">
              {{ tf('console.capability.duration', '总耗时') }}:
              {{ ((currentJob?.totalDuration || snapshot?.totalDuration || 0) / 1000).toFixed(1) }}s
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>

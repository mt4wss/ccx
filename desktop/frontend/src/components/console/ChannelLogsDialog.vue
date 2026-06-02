<script setup lang="ts">
import { ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Skeleton } from '@/components/ui/skeleton'
import { Badge } from '@/components/ui/badge'
import { Loader2, RefreshCw, X } from 'lucide-vue-next'
import { useAdminApi } from '@/composables/useAdminApi'
import { useLanguage } from '@/composables/useLanguage'
import type { ChannelLogEntry, ChannelLogsResponse } from '@/services/admin-api'

interface Props {
  open: boolean
  channelType: string
  channelId: number
  channelName: string
}

const props = defineProps<Props>()
const emit = defineEmits<{ (e: 'close'): void }>()

const { tf } = useLanguage()
const api = useAdminApi()

const logs = ref<ChannelLogEntry[]>([])
const loading = ref(false)
const error = ref('')

async function fetchLogs() {
  loading.value = true
  error.value = ''
  try {
    const data = await api.get<ChannelLogsResponse>(
      `/api/${props.channelType}/channels/${props.channelId}/logs`
    )
    logs.value = data.logs || []
  } catch (e) {
    error.value = e instanceof Error ? e.message : String(e)
  } finally {
    loading.value = false
  }
}

watch(() => props.open, (isOpen) => {
  if (isOpen) fetchLogs()
})

function formatTime(ts: string) {
  try {
    return new Date(ts).toLocaleTimeString()
  } catch {
    return ts
  }
}

function formatDuration(ms: number) {
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(1)}s`
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center"
        @keydown="onKeyDown"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="emit('close')" />

        <!-- Dialog -->
        <div class="relative z-10 w-[90vw] max-w-4xl max-h-[85vh] rounded-2xl border border-border bg-card shadow-2xl flex flex-col">
          <!-- Header -->
          <div class="flex items-center justify-between p-4 border-b border-border shrink-0">
            <div class="flex items-center gap-2">
              <h3 class="text-sm font-semibold">
                {{ tf('console.logs.title', '频道日志') }}: {{ channelName }}
              </h3>
              <Badge variant="secondary" class="text-[10px]">
                {{ logs.length }} {{ tf('console.logs.entries', '条') }}
              </Badge>
            </div>
            <div class="flex items-center gap-2">
              <Button variant="ghost" size="icon-sm" :disabled="loading" @click="fetchLogs">
                <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
              </Button>
              <Button variant="ghost" size="icon-sm" @click="emit('close')">
                <X class="w-4 h-4" />
              </Button>
            </div>
          </div>

          <!-- Body -->
          <div class="flex-1 min-h-0 overflow-hidden">
            <div v-if="loading && logs.length === 0" class="p-4 space-y-2">
              <Skeleton v-for="i in 5" :key="i" class="h-10 w-full" />
            </div>

            <div v-else-if="error" class="p-4 text-sm text-destructive">
              {{ error }}
            </div>

            <div v-else-if="logs.length === 0" class="p-8 text-center text-sm text-muted-foreground">
              {{ tf('console.logs.empty', '暂无日志') }}
            </div>

            <ScrollArea v-else class="h-full">
              <table class="w-full text-xs">
                <thead class="sticky top-0 bg-card border-b border-border">
                  <tr class="text-left text-muted-foreground">
                    <th class="p-2 font-medium">{{ tf('console.logs.time', '时间') }}</th>
                    <th class="p-2 font-medium">{{ tf('console.logs.model', '模型') }}</th>
                    <th class="p-2 font-medium">{{ tf('console.logs.statusCode', '状态码') }}</th>
                    <th class="p-2 font-medium">{{ tf('console.logs.duration', '耗时') }}</th>
                    <th class="p-2 font-medium">{{ tf('console.logs.key', 'Key') }}</th>
                    <th class="p-2 font-medium">{{ tf('console.logs.baseUrl', 'URL') }}</th>
                    <th class="p-2 font-medium">{{ tf('console.logs.error', '错误') }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="log in logs"
                    :key="log.requestId"
                    class="border-b border-border/50 hover:bg-accent/50"
                  >
                    <td class="p-2 font-mono text-[10px]">{{ formatTime(log.timestamp) }}</td>
                    <td class="p-2 font-mono">{{ log.model }}</td>
                    <td class="p-2">
                      <Badge
                        :variant="log.success ? 'default' : 'destructive'"
                        class="text-[10px]"
                      >
                        {{ log.statusCode }}
                      </Badge>
                    </td>
                    <td class="p-2 font-mono">{{ formatDuration(log.durationMs) }}</td>
                    <td class="p-2 font-mono text-[10px] text-muted-foreground">{{ log.keyMask }}</td>
                    <td class="p-2 font-mono text-[10px] text-muted-foreground max-w-[200px] truncate">
                      {{ log.baseUrl }}
                    </td>
                    <td class="p-2 text-destructive max-w-[150px] truncate" :title="log.errorInfo">
                      {{ log.errorInfo }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </ScrollArea>
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

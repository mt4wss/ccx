<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Play, Square, RotateCcw, Globe, ExternalLink, RefreshCw } from 'lucide-vue-next'
import type { DesktopStatus } from '@/types'
import { useLanguage } from '@/composables/useLanguage'

defineProps<{
  status: DesktopStatus
  loading: boolean
}>()

const emit = defineEmits<{
  start: []
  stop: []
  restart: []
  openWebUI: []
  openBrowser: []
  refresh: []
}>()

const { t } = useLanguage()
</script>

<template>
  <div class="flex flex-wrap items-center justify-between gap-3 bg-glass dark:bg-glass-dark border border-border p-3.5 rounded-xl select-none shrink-0">
    <!-- 主运行操作组合 (启动, 停止, 重启) -->
    <div class="flex flex-wrap items-center gap-2.5">
      <!-- 启动按钮 (带 Shimmer 渐变) -->
      <Button
        size="sm"
        :disabled="loading || status.running"
        @click="emit('start')"
        class="bg-blue-600 hover:bg-blue-500 text-white font-semibold shadow-md active:scale-95 hover:scale-[1.02] transition-all duration-300 btn-shimmer disabled:opacity-30 border border-blue-500/10 cursor-pointer"
      >
        <Play class="w-3.5 h-3.5 mr-1.5 fill-white" />
        {{ t('actions.start') }}
      </Button>

      <!-- 停止按钮 -->
      <Button
        size="sm"
        variant="secondary"
        :disabled="loading || !status.running || status.attached"
        @click="emit('stop')"
        class="bg-secondary hover:bg-secondary/80 text-foreground border border-border active:scale-95 hover:scale-[1.02] transition-all duration-300 disabled:opacity-20 cursor-pointer"
      >
        <Square class="w-3.5 h-3.5 mr-1.5 fill-current" />
        {{ t('actions.stop') }}
      </Button>

      <!-- 重启按钮 -->
      <Button
        size="sm"
        variant="secondary"
        :disabled="loading || status.attached"
        @click="emit('restart')"
        class="bg-secondary hover:bg-secondary/80 text-foreground border border-border active:scale-95 hover:scale-[1.02] transition-all duration-300 disabled:opacity-20 cursor-pointer"
      >
        <RotateCcw class="w-3.5 h-3.5 mr-1.5" />
        {{ t('actions.restart') }}
      </Button>
    </div>

    <!-- 附属功能组合 (浏览器, 刷新, 内嵌) -->
    <div class="flex flex-wrap items-center gap-2">
      <!-- 内嵌 Web UI -->
      <Button
        size="sm"
        variant="outline"
        :disabled="loading"
        @click="emit('openWebUI')"
        class="bg-background/40 border border-border hover:bg-secondary hover:text-foreground hover:border-border text-muted-foreground active:scale-95 transition-all duration-200 cursor-pointer"
      >
        <Globe class="w-3.5 h-3.5 mr-1.5 text-blue-500/80" />
        {{ t('actions.openWebUI') }}
      </Button>

      <!-- 浏览器中打开 -->
      <Button
        size="sm"
        variant="outline"
        :disabled="loading"
        @click="emit('openBrowser')"
        class="bg-background/40 border border-border hover:bg-secondary hover:text-foreground hover:border-border text-muted-foreground active:scale-95 transition-all duration-200 cursor-pointer"
      >
        <ExternalLink class="w-3.5 h-3.5 mr-1.5 text-emerald-500/80" />
        {{ t('actions.openBrowser') }}
      </Button>

      <!-- 刷新状态 -->
      <Button
        size="sm"
        variant="ghost"
        :disabled="loading"
        @click="emit('refresh')"
        class="text-muted-foreground hover:text-foreground hover:bg-secondary p-2 rounded-lg cursor-pointer transition-colors"
        :title="t('actions.refreshStatus')"
      >
        <RefreshCw class="w-3.5 h-3.5" :class="loading ? 'animate-spin' : ''" />
      </Button>
    </div>
  </div>
</template>

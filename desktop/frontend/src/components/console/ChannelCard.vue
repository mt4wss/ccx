<script setup lang="ts">
import { computed } from 'vue'
import type {
  Channel,
  ChannelMetrics,
  ChannelRecentActivity,
} from '@/services/admin-api'
import { expandSparseSegments } from '@/services/admin-api'
import { useLanguage } from '@/composables/useLanguage'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Activity,
  AlertTriangle,
  Ban,
  CheckCircle2,
  ChevronRight,
  Edit3,
  ExternalLink,
  Eye,
  Gauge,
  Key,
  MoreVertical,
  Pause,
  Play,
  RotateCcw,
  Sparkles,
  Terminal,
  Timer,
  Trash2,
  XCircle,
  Zap,
} from 'lucide-vue-next'

const props = defineProps<{
  channel: Channel
  metrics?: ChannelMetrics
  activity?: ChannelRecentActivity
}>()

const emit = defineEmits<{
  edit: []
  delete: []
  ping: []
  keys: []
  logs: []
  capability: []
  status: []
  resume: []
  promote: []
}>()

const { t, tf } = useLanguage()

// ── Service Type ──────────────────────────────────────

const serviceTypeBadgeVariant = computed(() => {
  const map: Record<string, string> = {
    openai: 'default',
    claude: 'secondary',
    gemini: 'outline',
    responses: 'default',
  }
  return (map[props.channel.serviceType] || 'default') as 'default' | 'secondary' | 'outline' | 'destructive'
})

const serviceTypeColors = computed(() => {
  const map: Record<string, { bg: string; border: string; text: string }> = {
    openai: { bg: 'bg-blue-500/10 dark:bg-blue-500/15', border: 'border-blue-500/20', text: 'text-blue-700 dark:text-blue-400' },
    claude: { bg: 'bg-orange-500/10 dark:bg-orange-500/15', border: 'border-orange-500/20', text: 'text-orange-700 dark:text-orange-400' },
    gemini: { bg: 'bg-purple-500/10 dark:bg-purple-500/15', border: 'border-purple-500/20', text: 'text-purple-700 dark:text-purple-400' },
    responses: { bg: 'bg-emerald-500/10 dark:bg-emerald-500/15', border: 'border-emerald-500/20', text: 'text-emerald-700 dark:text-emerald-400' },
  }
  return map[props.channel.serviceType] || map.openai
})

// ── Status ────────────────────────────────────────────

const isSuspended = computed(() => props.channel.status === 'suspended')
const isDisabled = computed(() => props.channel.status === 'disabled')
const isActive = computed(() => !isSuspended.value && !isDisabled.value)

const statusConfig = computed(() => {
  if (isDisabled.value) {
    return {
      label: tf('console.channelStatus.disabled', 'Disabled'),
      bgClass: 'bg-red-500/10 dark:bg-red-500/15',
      borderClass: 'border-red-500/20',
      textClass: 'text-red-700 dark:text-red-400',
      glowClass: 'shadow-glow-red',
      icon: XCircle,
    }
  }
  if (isSuspended.value) {
    return {
      label: tf('console.channelStatus.suspended', 'Suspended'),
      bgClass: 'bg-amber-500/10 dark:bg-amber-500/15',
      borderClass: 'border-amber-500/20',
      textClass: 'text-amber-700 dark:text-amber-400',
      glowClass: 'shadow-glow-orange',
      icon: Pause,
    }
  }
  return {
    label: tf('console.channelStatus.active', 'Active'),
    bgClass: 'bg-emerald-500/10 dark:bg-emerald-500/15',
    borderClass: 'border-emerald-500/20',
    textClass: 'text-emerald-700 dark:text-emerald-400',
    glowClass: 'shadow-glow-green',
    icon: CheckCircle2,
  }
})

// ── Key Count ─────────────────────────────────────────

const totalKeyCount = computed(() => props.channel.apiKeys?.length ?? 0)
const disabledKeyCount = computed(() => props.channel.disabledApiKeys?.length ?? 0)

// ── Latency ───────────────────────────────────────────

const latencyLevel = computed(() => {
  const ms = props.metrics?.latency ?? props.channel.latency
  if (ms == null) return null
  if (ms < 200) return { label: 'fast', bg: 'bg-emerald-500/10 dark:bg-emerald-500/15', border: 'border-emerald-500/20', text: 'text-emerald-700 dark:text-emerald-400' }
  if (ms < 500) return { label: 'ok', bg: 'bg-blue-500/10 dark:bg-blue-500/15', border: 'border-blue-500/20', text: 'text-blue-700 dark:text-blue-400' }
  if (ms < 1000) return { label: 'fair', bg: 'bg-amber-500/10 dark:bg-amber-500/15', border: 'border-amber-500/20', text: 'text-amber-700 dark:text-amber-400' }
  return { label: 'slow', bg: 'bg-red-500/10 dark:bg-red-500/15', border: 'border-red-500/20', text: 'text-red-700 dark:text-red-400' }
})

const latencyDisplay = computed(() => {
  const ms = props.metrics?.latency ?? props.channel.latency
  if (ms == null) return null
  return `${Math.round(ms)}ms`
})

// ── Metrics ───────────────────────────────────────────

const formattedSuccessRate = computed(() => {
  if (!props.metrics) return null
  return `${(props.metrics.successRate * 100).toFixed(1)}%`
})

const formattedErrorRate = computed(() => {
  if (!props.metrics) return null
  return `${(props.metrics.errorRate * 100).toFixed(1)}%`
})

// ── Circuit Breaker ───────────────────────────────────

const circuitState = computed(() => props.metrics?.circuitState)

const circuitDisplay = computed(() => {
  if (!circuitState.value || circuitState.value === 'closed') return null
  if (circuitState.value === 'open') {
    return {
      label: tf('console.circuit.open', 'Circuit Open'),
      icon: AlertTriangle,
      bgClass: 'bg-red-500/10 dark:bg-red-500/15',
      borderClass: 'border-red-500/20',
      textClass: 'text-red-700 dark:text-red-400',
    }
  }
  return {
    label: tf('console.circuit.halfOpen', 'Half-Open'),
    icon: AlertTriangle,
    bgClass: 'bg-amber-500/10 dark:bg-amber-500/15',
    borderClass: 'border-amber-500/20',
    textClass: 'text-amber-700 dark:text-amber-400',
  }
})

// ── Promotion ─────────────────────────────────────────

const isPromoted = computed(() => {
  if (!props.channel.promotionUntil) return false
  return new Date(props.channel.promotionUntil).getTime() > Date.now()
})

// ── Activity Sparkline ────────────────────────────────

const sparklineSegments = computed(() => {
  if (!props.activity) return []
  const segs = expandSparseSegments(props.activity)
  const maxReq = Math.max(1, ...segs.map(s => s.requestCount))
  return segs.map(s => {
    if (s.requestCount === 0) return { height: 0, color: 'bg-transparent' }
    const ratio = s.requestCount / maxReq
    const hasFailure = s.failureCount > 0
    const height = Math.max(2, Math.round(ratio * 100))
    let color: string
    if (hasFailure && s.failureCount >= s.successCount) {
      color = 'bg-red-400/70 dark:bg-red-500/60'
    } else if (hasFailure) {
      color = 'bg-amber-400/70 dark:bg-amber-500/60'
    } else {
      color = 'bg-emerald-400/70 dark:bg-emerald-500/60'
    }
    return { height, color }
  })
})

const hasActivity = computed(() =>
  props.activity && sparklineSegments.value.some(s => s.height > 0),
)
</script>

<template>
  <TooltipProvider :delay-duration="200">
    <div
      :class="[
        'group relative flex flex-col gap-4 rounded-xl border p-5 transition-all duration-300',
        'bg-glass dark:bg-glass-dark bg-glass-hover dark:bg-glass-hover-dark',
        statusConfig.glowClass,
      ]"
    >
      <!-- ── Header: Name + Badges ─────────────────── -->
      <div class="flex items-start justify-between gap-3">
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2.5">
            <h3 class="truncate text-sm font-semibold text-foreground">
              {{ channel.name }}
            </h3>
            <Badge :variant="serviceTypeBadgeVariant" class="shrink-0">
              {{ channel.serviceType.toUpperCase() }}
            </Badge>
            <Tooltip v-if="isPromoted">
              <TooltipTrigger as-child>
                <span
                  class="inline-flex items-center gap-1 shrink-0 rounded-full border border-purple-500/25 bg-purple-500/10 px-2 py-0.5 text-[10px] font-bold text-purple-700 dark:text-purple-400"
                >
                  <Sparkles class="h-3 w-3" />
                  PROMO
                </span>
              </TooltipTrigger>
              <TooltipContent>
                <p>
                  {{ tf('console.promotion.until', 'Promoted until') }}
                  {{ new Date(channel.promotionUntil!).toLocaleString() }}
                </p>
              </TooltipContent>
            </Tooltip>
          </div>
          <p
            v-if="channel.description"
            class="mt-1 text-xs text-muted-foreground line-clamp-2"
          >
            {{ channel.description }}
          </p>
        </div>

        <!-- Status Badge -->
        <div
          :class="[
            'inline-flex items-center gap-1.5 shrink-0 rounded-full border px-2.5 py-1 text-[11px] font-semibold',
            statusConfig.bgClass,
            statusConfig.borderClass,
            statusConfig.textClass,
          ]"
        >
          <component :is="statusConfig.icon" class="h-3 w-3" />
          {{ statusConfig.label }}
        </div>
      </div>

      <!-- ── Key Count + Latency ───────────────────── -->
      <div class="flex items-center gap-3 flex-wrap">
        <!-- Key Count -->
        <Tooltip>
          <TooltipTrigger as-child>
            <div
              class="inline-flex items-center gap-1.5 rounded-lg border border-border bg-secondary/50 px-2.5 py-1 text-xs font-medium text-foreground cursor-default"
            >
              <Key class="h-3 w-3 text-muted-foreground" />
              <span>{{ totalKeyCount }}</span>
              <span
                v-if="disabledKeyCount > 0"
                class="ml-0.5 inline-flex items-center gap-0.5 text-destructive"
              >
                <Ban class="h-3 w-3" />
                {{ disabledKeyCount }}
              </span>
            </div>
          </TooltipTrigger>
          <TooltipContent>
            <p>
              {{ totalKeyCount }} {{ tf('console.keys.active', 'active keys') }}
              <template v-if="disabledKeyCount > 0">
                , {{ disabledKeyCount }} {{ tf('console.keys.disabled', 'disabled') }}
              </template>
            </p>
          </TooltipContent>
        </Tooltip>

        <!-- Latency -->
        <Tooltip v-if="latencyDisplay">
          <TooltipTrigger as-child>
            <div
              :class="[
                'inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1 text-xs font-semibold cursor-default',
                latencyLevel!.bg,
                latencyLevel!.border,
                latencyLevel!.text,
              ]"
            >
              <Timer class="h-3 w-3" />
              {{ latencyDisplay }}
            </div>
          </TooltipTrigger>
          <TooltipContent>
            <p>{{ tf('console.latency.tooltip', 'Response latency') }}</p>
          </TooltipContent>
        </Tooltip>

        <!-- Spacer -->
        <div class="flex-1" />

        <!-- Circuit Breaker -->
        <Tooltip v-if="circuitDisplay">
          <TooltipTrigger as-child>
            <div
              :class="[
                'inline-flex items-center gap-1.5 rounded-lg border px-2.5 py-1 text-[11px] font-bold cursor-default',
                circuitDisplay.bgClass,
                circuitDisplay.borderClass,
                circuitDisplay.textClass,
              ]"
            >
              <component :is="circuitDisplay.icon" class="h-3 w-3" />
              {{ circuitDisplay.label }}
            </div>
          </TooltipTrigger>
          <TooltipContent>
            <p v-if="metrics?.circuitBrokenAt">
              {{ tf('console.circuit.brokenAt', 'Broken at') }}:
              {{ new Date(metrics.circuitBrokenAt).toLocaleString() }}
            </p>
            <p v-if="metrics?.nextRetryAt">
              {{ tf('console.circuit.retryAt', 'Next retry') }}:
              {{ new Date(metrics.nextRetryAt).toLocaleString() }}
            </p>
          </TooltipContent>
        </Tooltip>
      </div>

      <!-- ── Metrics Row ───────────────────────────── -->
      <div
        v-if="metrics"
        class="grid grid-cols-3 gap-3"
      >
        <!-- Request Count -->
        <div class="bg-glass dark:bg-glass-dark rounded-lg border border-border p-3 flex flex-col gap-1.5">
          <div class="flex items-center justify-between">
            <span class="text-[10px] font-bold tracking-wider text-muted-foreground uppercase">
              {{ tf('console.metrics.requests', 'Requests') }}
            </span>
            <div class="rounded-md bg-blue-500/10 border border-blue-500/15 p-1">
              <Activity class="h-3 w-3 text-blue-700 dark:text-blue-400" />
            </div>
          </div>
          <span class="text-lg font-bold font-mono tracking-tight text-foreground">
            {{ metrics.requestCount.toLocaleString() }}
          </span>
        </div>

        <!-- Success Rate -->
        <div class="bg-glass dark:bg-glass-dark rounded-lg border border-border p-3 flex flex-col gap-1.5">
          <div class="flex items-center justify-between">
            <span class="text-[10px] font-bold tracking-wider text-muted-foreground uppercase">
              {{ tf('console.metrics.successRate', 'Success') }}
            </span>
            <div class="rounded-md bg-emerald-500/10 border border-emerald-500/15 p-1">
              <CheckCircle2 class="h-3 w-3 text-emerald-700 dark:text-emerald-400" />
            </div>
          </div>
          <span class="text-lg font-bold font-mono tracking-tight text-foreground">
            {{ formattedSuccessRate }}
          </span>
        </div>

        <!-- Error Rate -->
        <div class="bg-glass dark:bg-glass-dark rounded-lg border border-border p-3 flex flex-col gap-1.5">
          <div class="flex items-center justify-between">
            <span class="text-[10px] font-bold tracking-wider text-muted-foreground uppercase">
              {{ tf('console.metrics.errorRate', 'Errors') }}
            </span>
            <div class="rounded-md bg-red-500/10 border border-red-500/15 p-1">
              <XCircle class="h-3 w-3 text-red-700 dark:text-red-400" />
            </div>
          </div>
          <span class="text-lg font-bold font-mono tracking-tight text-foreground">
            {{ formattedErrorRate }}
          </span>
        </div>
      </div>

      <!-- ── Activity Sparkline ────────────────────── -->
      <div v-if="activity && hasActivity" class="flex flex-col gap-1.5">
        <div class="flex items-center justify-between">
          <span class="text-[10px] font-bold tracking-wider text-muted-foreground uppercase">
            {{ tf('console.activity.title', 'Recent Activity') }}
          </span>
          <div class="flex items-center gap-2 text-[10px] text-muted-foreground font-mono">
            <span v-if="activity.rpm > 0">{{ activity.rpm.toFixed(1) }} rpm</span>
            <span v-if="activity.tpm > 0">{{ activity.tpm.toFixed(0) }} tpm</span>
          </div>
        </div>
        <div
          class="flex items-end gap-px h-10 overflow-hidden rounded-md bg-secondary/30 px-0.5 pb-0.5"
        >
          <div
            v-for="(seg, i) in sparklineSegments"
            :key="i"
            :class="[
              'flex-1 rounded-[1px] transition-all duration-200 min-w-[1px]',
              seg.height > 0 ? seg.color : 'bg-transparent',
            ]"
            :style="{ height: seg.height > 0 ? `${seg.height}%` : '0' }"
          />
        </div>
      </div>

      <!-- ── Actions ───────────────────────────────── -->
      <div class="flex items-center justify-end gap-2 pt-1">
        <Tooltip>
          <TooltipTrigger as-child>
            <Button
              variant="outline"
              size="sm"
              @click="emit('ping')"
            >
              <Gauge class="h-3.5 w-3.5" />
              {{ tf('console.actions.ping', 'Ping') }}
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>{{ tf('console.actions.pingTip', 'Test channel latency') }}</p>
          </TooltipContent>
        </Tooltip>

        <Tooltip v-if="!isDisabled && isSuspended">
          <TooltipTrigger as-child>
            <Button
              variant="success"
              size="sm"
              @click="emit('status')"
            >
              <Play class="h-3.5 w-3.5" />
              {{ tf('console.actions.resume', 'Resume') }}
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>{{ tf('console.actions.resumeTip', 'Set channel to active') }}</p>
          </TooltipContent>
        </Tooltip>

        <Tooltip v-if="!isDisabled && isActive">
          <TooltipTrigger as-child>
            <Button
              variant="outline"
              size="sm"
              @click="emit('status')"
            >
              <Pause class="h-3.5 w-3.5" />
              {{ tf('console.actions.suspend', 'Suspend') }}
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>{{ tf('console.actions.suspendTip', 'Pause channel routing') }}</p>
          </TooltipContent>
        </Tooltip>

        <!-- More Actions Menu -->
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <Button variant="ghost" size="icon-sm">
              <MoreVertical class="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" class="w-48">
            <DropdownMenuLabel>
              {{ tf('console.actions.label', 'Actions') }}
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuItem @click="emit('edit')">
                <Edit3 class="h-4 w-4" />
                {{ tf('console.actions.edit', 'Edit Channel') }}
              </DropdownMenuItem>
              <DropdownMenuItem @click="emit('keys')">
                <Key class="h-4 w-4" />
                {{ tf('console.actions.manageKeys', 'Manage Keys') }}
              </DropdownMenuItem>
              <DropdownMenuItem @click="emit('capability')">
                <Zap class="h-4 w-4" />
                {{ tf('console.actions.capability', 'Capability Test') }}
              </DropdownMenuItem>
              <DropdownMenuItem @click="emit('logs')">
                <Terminal class="h-4 w-4" />
                {{ tf('console.actions.logs', 'View Logs') }}
              </DropdownMenuItem>
              <DropdownMenuItem
                v-if="channel.website"
                as="a"
                :href="channel.website"
                target="_blank"
                rel="noopener"
              >
                <ExternalLink class="h-4 w-4" />
                {{ tf('console.actions.website', 'Visit Website') }}
              </DropdownMenuItem>
            </DropdownMenuGroup>

            <!-- Circuit breaker resume -->
            <template v-if="circuitDisplay">
              <DropdownMenuSeparator />
              <DropdownMenuItem @click="emit('resume')">
                <RotateCcw class="h-4 w-4" />
                {{ tf('console.actions.resetCircuit', 'Reset Circuit Breaker') }}
              </DropdownMenuItem>
            </template>

            <!-- Promotion -->
            <template v-if="!isPromoted && !isDisabled">
              <DropdownMenuSeparator />
              <DropdownMenuItem @click="emit('promote')">
                <Sparkles class="h-4 w-4" />
                {{ tf('console.actions.promote', 'Promote') }}
              </DropdownMenuItem>
            </template>

            <DropdownMenuSeparator />
            <DropdownMenuItem
              variant="destructive"
              @click="emit('delete')"
            >
              <Trash2 class="h-4 w-4" />
              {{ tf('console.actions.delete', 'Delete Channel') }}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  </TooltipProvider>
</template>

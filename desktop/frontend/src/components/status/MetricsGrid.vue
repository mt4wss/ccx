<script setup lang="ts">
import { computed } from 'vue'
import type { DesktopStatus } from '@/types'
import { Server, Clock, GitBranch, ArrowUpRight } from 'lucide-vue-next'

const props = defineProps<{
  status: DesktopStatus
}>()

// 精细运行时长计算
const uptimeDisplay = computed(() => {
  const uptime = props.status.health?.uptime
  if (!uptime) return '——'
  if (uptime < 60) return `${uptime}s`
  const minutes = Math.floor(uptime / 60)
  if (minutes < 60) return `${minutes}m`
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  return `${hours}h ${remainingMinutes}m`
})
</script>

<template>
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 select-none">
    <!-- 1. 网关运行端口 -->
    <div class="bg-glass bg-glass-hover rounded-xl p-4.5 border border-white/[0.03] flex flex-col justify-between group">
      <div class="flex items-center justify-between">
        <span class="text-[11px] font-bold tracking-wider text-slate-500 uppercase">网关端口</span>
        <div class="p-1.5 rounded-lg bg-blue-500/10 border border-blue-500/15 group-hover:bg-blue-500/20 transition-colors">
          <Server class="w-3.5 h-3.5 text-blue-400" />
        </div>
      </div>
      <div class="mt-4 flex items-baseline gap-1.5 min-w-0">
        <span class="text-xl font-bold font-mono tracking-tight text-slate-100 truncate">
          {{ status.port || '——' }}
        </span>
        <span class="text-[9px] font-bold text-emerald-400 bg-emerald-500/10 px-1 py-0.2 rounded border border-emerald-500/15 uppercase tracking-wide shrink-0" v-if="status.running">
          Active
        </span>
      </div>
    </div>

    <!-- 2. 网关运行时长 -->
    <div class="bg-glass bg-glass-hover rounded-xl p-4.5 border border-white/[0.03] flex flex-col justify-between group">
      <div class="flex items-center justify-between">
        <span class="text-[11px] font-bold tracking-wider text-slate-500 uppercase">运行时长</span>
        <div class="p-1.5 rounded-lg bg-emerald-500/10 border border-emerald-500/15 group-hover:bg-emerald-500/20 transition-colors">
          <Clock class="w-3.5 h-3.5 text-emerald-400" />
        </div>
      </div>
      <div class="mt-4 flex items-baseline gap-1.5 min-w-0">
        <span class="text-xl font-bold font-mono tracking-tight text-slate-100 truncate">
          {{ uptimeDisplay }}
        </span>
        <span class="text-[9px] text-slate-500 font-medium shrink-0" v-if="status.health?.uptime">
          uptime
        </span>
      </div>
    </div>

    <!-- 3. 上游渠道数 -->
    <div class="bg-glass bg-glass-hover rounded-xl p-4.5 border border-white/[0.03] flex flex-col justify-between group">
      <div class="flex items-center justify-between">
        <span class="text-[11px] font-bold tracking-wider text-slate-500 uppercase">调度信道</span>
        <div class="p-1.5 rounded-lg bg-indigo-500/10 border border-indigo-500/15 group-hover:bg-indigo-500/20 transition-colors">
          <GitBranch class="w-3.5 h-3.5 text-indigo-400" />
        </div>
      </div>
      <div class="mt-4 flex items-baseline gap-1.5 min-w-0">
        <span class="text-xl font-bold font-mono tracking-tight text-slate-100 truncate">
          {{ status.health?.config?.upstreamCount || 0 }}
        </span>
        <span class="text-[9px] text-indigo-400 bg-indigo-500/10 px-1 py-0.2 rounded border border-indigo-500/15 uppercase tracking-wide shrink-0">
          Channels
        </span>
      </div>
    </div>

    <!-- 4. 网关版本 -->
    <div class="bg-glass bg-glass-hover rounded-xl p-4.5 border border-white/[0.03] flex flex-col justify-between group">
      <div class="flex items-center justify-between">
        <span class="text-[11px] font-bold tracking-wider text-slate-500 uppercase">网关版本</span>
        <div class="p-1.5 rounded-lg bg-slate-500/10 border border-slate-500/15 group-hover:bg-slate-500/20 transition-colors">
          <ArrowUpRight class="w-3.5 h-3.5 text-slate-400" />
        </div>
      </div>
      <div class="mt-4 flex items-baseline gap-1.5 min-w-0">
        <span class="text-xl font-bold font-mono tracking-tight text-slate-100 truncate">
          {{ status.health?.version?.version || 'v0.0.0' }}
        </span>
        <span class="text-[9px] text-slate-500 font-medium shrink-0">
          stable
        </span>
      </div>
    </div>
  </div>
</template>

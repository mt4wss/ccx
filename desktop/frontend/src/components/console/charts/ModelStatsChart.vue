<template>
  <div class="model-stats-chart">
    <!-- Compact summary (top models) -->
    <div v-if="topModels.length" class="flex flex-wrap items-center gap-3 mb-2 text-xs bg-secondary/20 dark:bg-secondary/10 rounded px-2 py-1.5">
      <span v-for="(m, i) in topModels" :key="m.name" class="flex items-center gap-1">
        <span class="w-2 h-2 rounded-full" :style="{ backgroundColor: MODEL_COLORS[i % MODEL_COLORS.length] }" />
        <span class="font-medium">{{ m.name }}</span>
        <span class="text-muted-foreground">{{ formatNumber(m.count) }} 次</span>
      </span>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex items-center justify-center" style="height: 200px">
      <div class="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin" />
    </div>

    <!-- Empty state -->
    <div
      v-else-if="!hasData"
      class="flex flex-col items-center justify-center text-muted-foreground"
      style="height: 200px"
    >
      <div class="text-2xl mb-2 opacity-40">&#x1F4CA;</div>
      <div class="text-xs">暂无模型统计数据</div>
    </div>

    <!-- Chart -->
    <div v-else>
      <VueApexCharts
        ref="chartRef"
        type="area"
        :height="200"
        :options="chartOptions"
        :series="chartSeries"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import VueApexCharts from 'vue3-apexcharts'
import type { ApexOptions } from 'apexcharts'
import { useTheme } from '@/composables/useTheme'
import type { ModelHistoryDataPoint } from '@/services/admin-api'

const props = withDefaults(
  defineProps<{
    data: Record<string, ModelHistoryDataPoint[]>
    duration?: string
    loading?: boolean
  }>(),
  {
    duration: '6h',
    loading: false,
  },
)

const { theme } = useTheme()
const chartRef = ref<InstanceType<typeof VueApexCharts> | null>(null)

const isDark = computed(() => {
  if (theme.value === 'dark') return true
  if (theme.value === 'auto') return window.matchMedia('(prefers-color-scheme: dark)').matches
  return false
})

const textColor = computed(() => (isDark.value ? '#94a3b8' : '#64748b'))
const gridBorder = computed(() => (isDark.value ? 'rgba(255,255,255,0.06)' : 'rgba(0,0,0,0.06)'))

const MODEL_COLORS = [
  '#3b82f6', '#10b981', '#f97316', '#8b5cf6', '#ef4444',
  '#06b6d4', '#ec4899', '#84cc16', '#f59e0b', '#6366f1',
]

const sortedModels = computed(() => {
  if (!props.data) return []
  return Object.entries(props.data)
    .map(([name, points]) => ({
      name,
      points,
      totalRequests: points.reduce((s, p) => s + p.requestCount, 0),
    }))
    .sort((a, b) => b.totalRequests - a.totalRequests)
})

const topModels = computed(() =>
  sortedModels.value.slice(0, 5).map(m => ({ name: m.name, count: m.totalRequests })),
)

const hasData = computed(() => sortedModels.value.some(m => m.totalRequests > 0))

const xLabelFormat = computed(() =>
  props.duration === '7d' || props.duration === '30d' ? 'MM-dd HH:mm' : 'HH:mm',
)

const chartOptions = computed<ApexOptions>(() => ({
  chart: {
    toolbar: { show: false },
    zoom: { enabled: false },
    background: 'transparent',
    fontFamily: 'inherit',
    stacked: true,
    animations: { enabled: true, speed: 400 },
  },
  theme: { mode: isDark.value ? 'dark' : 'light' },
  colors: MODEL_COLORS.slice(0, sortedModels.value.length),
  fill: {
    type: 'gradient',
    gradient: { shadeIntensity: 1, opacityFrom: 0.3, opacityTo: 0.05, stops: [0, 90, 100] },
  },
  dataLabels: { enabled: false },
  stroke: { curve: 'smooth', width: 2 },
  grid: { borderColor: gridBorder.value, padding: { left: 10, right: 10 } },
  xaxis: {
    type: 'datetime',
    labels: {
      datetimeUTC: false,
      format: xLabelFormat.value,
      style: { fontSize: '10px', colors: textColor.value },
    },
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: {
    labels: {
      formatter: (val: number) => Math.round(val).toString(),
      style: { fontSize: '11px', colors: textColor.value },
    },
    min: 0,
  },
  tooltip: {
    x: { format: 'MM-dd HH:mm' },
    y: { formatter: (val: number) => `${Math.round(val)} 次` },
  },
  legend: {
    show: true,
    position: 'top',
    horizontalAlign: 'right',
    fontSize: '11px',
    markers: { size: 4 },
    labels: { colors: textColor.value },
  },
}))

const chartSeries = computed(() => {
  return sortedModels.value.map(m => ({
    name: m.name,
    data: m.points.map(p => ({
      x: new Date(p.timestamp).getTime(),
      y: p.requestCount,
    })),
  }))
})

function formatNumber(num: number): string {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toFixed(0)
}

defineExpose({ chartRef })
</script>

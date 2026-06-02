<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Alert } from '@/components/ui/alert'
import { Skeleton } from '@/components/ui/skeleton'
import { Plus, Search, Wifi, Loader2 } from 'lucide-vue-next'
import { useConsoleChannels } from '@/composables/useConsoleChannels'
import { useLanguage } from '@/composables/useLanguage'
import ChannelStatusBadge from '@/components/console/ChannelStatusBadge.vue'
import ChannelCard from '@/components/console/ChannelCard.vue'
import ChannelEditDialog from '@/components/console/ChannelEditDialog.vue'
import type { ManagedChannelType } from '@/utils/channel-type-api'
import type { Channel, ChannelMetrics, ChannelRecentActivity } from '@/services/admin-api'

interface Props {
  type: ManagedChannelType
}

const props = defineProps<Props>()

const { t, tf } = useLanguage()
const {
  activeTab,
  channelsByType,
  dashboardCache,
  isPingingAll,
  refreshChannels,
  pingChannel,
  pingAllChannels,
  deleteChannel,
  setChannelStatus,
  resumeChannel,
  promoteChannel,
  reorderChannels,
} = useConsoleChannels()

// 切换到当前协议的 tab
watch(() => props.type, (newType) => {
  activeTab.value = newType
}, { immediate: true })

// 频道数据
const channels = computed(() => channelsByType.value[props.type]?.channels || [])
const metrics = computed(() => dashboardCache.value[props.type]?.metrics || [])
const activity = computed(() => dashboardCache.value[props.type]?.recentActivity || [])

// 搜索过滤
const searchQuery = ref('')
const filteredChannels = computed(() => {
  if (!searchQuery.value.trim()) return channels.value
  const q = searchQuery.value.toLowerCase()
  return channels.value.filter(ch =>
    ch.name.toLowerCase().includes(q) ||
    ch.description?.toLowerCase().includes(q) ||
    ch.baseUrl.toLowerCase().includes(q)
  )
})

// Channel metrics map by index
const metricsMap = computed(() => {
  const map = new Map<number, ChannelMetrics>()
  for (const m of metrics.value) {
    map.set(m.channelIndex, m)
  }
  return map
})

// Activity map by index
const activityMap = computed(() => {
  const map = new Map<number, ChannelRecentActivity>()
  for (const a of activity.value) {
    map.set(a.channelIndex, a)
  }
  return map
})

// 加载状态
const loading = computed(() => channels.value.length === 0)

// 操作状态
const actionLoading = ref(false)
const actionError = ref('')

// 对话框状态
const showAddDialog = ref(false)
const editingChannel = ref<Channel | null>(null)

function clearActionError() {
  actionError.value = ''
}

async function handlePing(channelId: number) {
  clearActionError()
  try {
    await pingChannel(channelId)
  } catch (e) {
    actionError.value = e instanceof Error ? e.message : String(e)
  }
}

async function handlePingAll() {
  clearActionError()
  try {
    await pingAllChannels()
  } catch (e) {
    actionError.value = e instanceof Error ? e.message : String(e)
  }
}

async function handleDelete(channelId: number) {
  clearActionError()
  actionLoading.value = true
  try {
    await deleteChannel(channelId)
  } catch (e) {
    actionError.value = e instanceof Error ? e.message : String(e)
  } finally {
    actionLoading.value = false
  }
}

async function handleStatusToggle(channelId: number, currentStatus: string) {
  clearActionError()
  const newStatus = currentStatus === 'active' ? 'suspended' : 'active'
  try {
    await setChannelStatus(channelId, newStatus as 'active' | 'suspended')
  } catch (e) {
    actionError.value = e instanceof Error ? e.message : String(e)
  }
}

async function handleResume(channelId: number) {
  clearActionError()
  try {
    await resumeChannel(channelId)
  } catch (e) {
    actionError.value = e instanceof Error ? e.message : String(e)
  }
}

async function handlePromote(channelId: number, duration: number) {
  clearActionError()
  try {
    await promoteChannel(channelId, duration)
  } catch (e) {
    actionError.value = e instanceof Error ? e.message : String(e)
  }
}

function handleEdit(channel: Channel) {
  editingChannel.value = channel
  showAddDialog.value = true
}

function handleAdd() {
  editingChannel.value = null
  showAddDialog.value = true
}

// 拖拽排序
async function handleReorder(newOrder: number[]) {
  clearActionError()
  try {
    await reorderChannels(newOrder)
  } catch (e) {
    actionError.value = e instanceof Error ? e.message : String(e)
  }
}

// 拖拽状态
const draggedIndex = ref<number | null>(null)

function onDragStart(e: DragEvent, channelIndex: number) {
  draggedIndex.value = channelIndex
  e.dataTransfer?.setData('text/plain', String(channelIndex))
  e.dataTransfer!.effectAllowed = 'move'
}

function onDragOver(e: DragEvent) {
  e.preventDefault()
  e.dataTransfer!.dropEffect = 'move'
}

function onDrop(e: DragEvent, targetIndex: number) {
  e.preventDefault()
  const sourceIndex = draggedIndex.value
  if (sourceIndex === null || sourceIndex === targetIndex) return

  const currentChannels = [...filteredChannels.value]
  const sourcePos = currentChannels.findIndex(c => c.index === sourceIndex)
  const targetPos = currentChannels.findIndex(c => c.index === targetIndex)
  if (sourcePos === -1 || targetPos === -1) return

  // 重新排列
  const [moved] = currentChannels.splice(sourcePos, 1)
  currentChannels.splice(targetPos, 0, moved)
  const newOrder = currentChannels.map(c => c.index)
  handleReorder(newOrder)
  draggedIndex.value = null
}

function onDragEnd() {
  draggedIndex.value = null
}

onMounted(() => {
  activeTab.value = props.type
})
</script>

<template>
  <div class="space-y-4">
    <!-- 错误提示 -->
    <Alert v-if="actionError" variant="destructive" class="mb-4">
      <p class="text-sm">{{ actionError }}</p>
      <template #action>
        <Button variant="ghost" size="sm" @click="clearActionError">
          {{ t('common.cancel') }}
        </Button>
      </template>
    </Alert>

    <!-- 工具栏 -->
    <div class="flex items-center gap-3">
      <div class="relative flex-1 max-w-sm">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <Input
          v-model="searchQuery"
          :placeholder="tf('console.searchChannels', '搜索频道...')"
          class="pl-9"
        />
      </div>
      <Button variant="outline" size="sm" :disabled="isPingingAll" @click="handlePingAll">
        <Wifi v-if="!isPingingAll" class="w-3.5 h-3.5 mr-1.5" />
        <Loader2 v-else class="w-3.5 h-3.5 mr-1.5 animate-spin" />
        {{ tf('console.pingAll', '批量测速') }}
      </Button>
      <Button size="sm" @click="handleAdd">
        <Plus class="w-3.5 h-3.5 mr-1.5" />
        {{ tf('console.addChannel', '添加频道') }}
      </Button>
    </div>

    <!-- 频道列表 -->
    <div v-if="loading" class="space-y-3">
      <Skeleton v-for="i in 3" :key="i" class="h-24 w-full rounded-lg" />
    </div>

    <div v-else-if="filteredChannels.length === 0" class="text-center py-12">
      <p class="text-sm text-muted-foreground">
        {{ searchQuery
          ? tf('console.noSearchResults', '没有匹配的频道')
          : tf('console.noChannels', '暂无频道，点击上方按钮添加')
        }}
      </p>
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="channel in filteredChannels"
        :key="channel.index"
        draggable="true"
        class="transition-all"
        :class="{ 'opacity-50': draggedIndex === channel.index }"
        @dragstart="onDragStart($event, channel.index)"
        @dragover="onDragOver"
        @drop="onDrop($event, channel.index)"
        @dragend="onDragEnd"
      >
        <!-- ChannelCard 占位，等 ChannelCard agent 完成后替换 -->
        <ChannelCard
          :channel="channel"
          :metrics="metricsMap.get(channel.index)"
          :activity="activityMap.get(channel.index)"
          @edit="handleEdit(channel)"
          @delete="handleDelete(channel.index)"
          @ping="handlePing(channel.index)"
          @status="handleStatusToggle(channel.index, channel.status || 'active')"
          @resume="handleResume(channel.index)"
          @promote="handlePromote(channel.index, 300)"
        />
      </div>
    </div>

    <!-- ChannelEditDialog -->
    <ChannelEditDialog
      v-if="showAddDialog"
      :channel="editingChannel"
      :channel-type="type"
      @close="showAddDialog = false"
      @saved="showAddDialog = false"
    />
  </div>
</template>

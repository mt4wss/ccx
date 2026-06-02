<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Alert } from '@/components/ui/alert'
import { RefreshCw, Search } from 'lucide-vue-next'
import { useConversations } from '@/composables/useConversations'
import { useLanguage } from '@/composables/useLanguage'
import { useStatus } from '@/composables/useStatus'
import ConversationCard from './ConversationCard.vue'
import { ref, computed } from 'vue'

const { status } = useStatus()
const { tf } = useLanguage()
const {
  conversations,
  total,
  channelsByKind,
  overrides,
  loading,
  error,
  activeKind,
  fetchConversations,
  setOverride,
  removeOverride,
} = useConversations()

const searchQuery = ref('')

const filteredConversations = computed(() => {
  if (!searchQuery.value.trim()) return conversations.value
  const q = searchQuery.value.toLowerCase()
  return conversations.value.filter(c =>
    c.title?.toLowerCase().includes(q) ||
    c.id.toLowerCase().includes(q) ||
    c.userId.toLowerCase().includes(q) ||
    c.lastModel.toLowerCase().includes(q)
  )
})

const kinds = computed(() => Object.keys(channelsByKind.value))

function handleRefresh() {
  fetchConversations(activeKind.value || undefined)
}

onMounted(() => {
  if (status.value.running) {
    fetchConversations()
  }
})

watch(() => status.value.running, (running) => {
  if (running) fetchConversations()
})
</script>

<template>
  <div class="space-y-4">
    <!-- 错误提示 -->
    <Alert v-if="error" variant="destructive">
      <p class="text-sm">{{ error }}</p>
    </Alert>

    <!-- 工具栏 -->
    <div class="flex items-center gap-3">
      <div class="relative flex-1 max-w-sm">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <Input
          v-model="searchQuery"
          :placeholder="tf('console.conversations.search', '搜索会话...')"
          class="pl-9"
        />
      </div>

      <Select
        v-if="kinds.length > 0"
        :model-value="activeKind || '__all__'"
        @update:model-value="(v: any) => { activeKind = (!v || v === '__all__') ? '' : String(v); fetchConversations(activeKind || undefined) }"
      >
        <SelectTrigger class="w-[140px]">
          <SelectValue :placeholder="tf('console.conversations.allKinds', '所有类型')" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="__all__">所有类型</SelectItem>
          <SelectItem v-for="kind in kinds" :key="kind" :value="kind">
            {{ kind }}
          </SelectItem>
        </SelectContent>
      </Select>

      <Button variant="outline" size="sm" :disabled="loading" @click="handleRefresh">
        <RefreshCw class="w-3.5 h-3.5 mr-1.5" :class="{ 'animate-spin': loading }" />
        {{ tf('console.conversations.refresh', '刷新') }}
      </Button>

      <span class="text-xs text-muted-foreground">
        {{ tf('console.conversations.total', '共 {count} 个会话', { count: String(total) }) }}
      </span>
    </div>

    <!-- 会话列表 -->
    <div v-if="loading && conversations.length === 0" class="space-y-3">
      <Skeleton v-for="i in 3" :key="i" class="h-24 w-full rounded-lg" />
    </div>

    <div v-else-if="filteredConversations.length === 0" class="text-center py-12">
      <p class="text-sm text-muted-foreground">
        {{ searchQuery
          ? tf('console.conversations.noSearchResults', '没有匹配的会话')
          : tf('console.conversations.empty', '暂无活跃会话')
        }}
      </p>
    </div>

    <div v-else class="space-y-2">
      <ConversationCard
        v-for="conv in filteredConversations"
        :key="conv.id"
        :conversation="conv"
        :override="overrides[conv.id]"
        :channels-by-kind="channelsByKind[conv.kind]"
        @set-override="(id: string) => { /* TODO: open override dialog to select channel sequence */ }"
        @remove-override="removeOverride"
      />
    </div>
  </div>
</template>

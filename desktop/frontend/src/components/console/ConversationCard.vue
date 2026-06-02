<script setup lang="ts">
import { computed } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import type { ConversationInfo, SequenceOverrideInfo, ChannelSequenceEntry } from '@/services/admin-api'

interface Props {
  conversation: ConversationInfo
  override?: SequenceOverrideInfo
  channelsByKind?: ChannelSequenceEntry[]
}

defineProps<Props>()

defineEmits<{
  (e: 'setOverride', id: string): void
  (e: 'removeOverride', id: string): void
}>()

function formatTime(ts: string) {
  try {
    const d = new Date(ts)
    const now = Date.now()
    const diff = now - d.getTime()
    if (diff < 60000) return '刚刚'
    if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
    return d.toLocaleDateString()
  } catch {
    return ts
  }
}
</script>

<template>
  <div class="border border-border rounded-lg p-4 bg-card hover:bg-accent/30 transition-colors space-y-3">
    <div class="flex items-start justify-between gap-4">
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-2 mb-1">
          <span class="text-sm font-medium truncate">{{ conversation.title || conversation.id }}</span>
          <Badge variant="secondary" class="text-[10px] shrink-0">{{ conversation.kind }}</Badge>
          <Badge v-if="override" variant="default" class="text-[10px] shrink-0">Override</Badge>
        </div>
        <div class="text-xs text-muted-foreground space-y-0.5">
          <p>用户: {{ conversation.userId }} | 模型: {{ conversation.lastModel }}</p>
          <p>请求数: {{ conversation.requestCount }} | 频道: {{ conversation.currentChannel }}</p>
          <p>最后活跃: {{ formatTime(conversation.lastActiveAt) }}</p>
        </div>
      </div>
      <div class="flex items-center gap-1.5 shrink-0">
        <Button
          v-if="!override"
          variant="outline"
          size="sm"
          @click="$emit('setOverride', conversation.id)"
        >
          覆盖频道
        </Button>
        <Button
          v-else
          variant="destructive"
          size="sm"
          @click="$emit('removeOverride', conversation.id)"
        >
          移除覆盖
        </Button>
      </div>
    </div>
    <div v-if="override" class="text-xs text-muted-foreground bg-muted/50 rounded p-2">
      <p class="font-medium mb-1">当前覆盖序列:</p>
      <p>{{ override.sequence.map(s => s.channelName).join(' → ') }}</p>
    </div>
  </div>
</template>

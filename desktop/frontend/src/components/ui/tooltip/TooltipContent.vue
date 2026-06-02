<script setup lang="ts">
import { TooltipContent, type TooltipContentProps } from 'reka-ui'
import { cn } from '@/lib/utils'

interface Props extends TooltipContentProps {
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  sideOffset: 4,
})
</script>

<template>
  <TooltipContent
    v-bind="props"
    :class="cn(
      'z-50 overflow-hidden rounded-md border bg-popover px-3 py-1.5 text-sm text-popover-foreground shadow-md animate-in fade-in-0 zoom-in-95',
      props.class
    )"
  >
    <slot />
  </TooltipContent>
</template>

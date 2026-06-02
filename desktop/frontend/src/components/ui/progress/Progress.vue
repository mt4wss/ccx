<script setup lang="ts">
import { cn } from '@/lib/utils'

interface ProgressProps {
  modelValue?: number
  max?: number
  class?: string
}

const props = withDefaults(defineProps<ProgressProps>(), {
  modelValue: 0,
  max: 100,
})

const percentage = computed(() =>
  props.max > 0 ? Math.min(100, (props.modelValue / props.max) * 100) : 0
)

import { computed } from 'vue'
</script>

<template>
  <div
    :class="cn(
      'relative h-2 w-full overflow-hidden rounded-full bg-primary/20',
      props.class
    )"
  >
    <div
      class="h-full w-full flex-1 bg-primary transition-all"
      :style="{ transform: `translateX(-${100 - percentage}%)` }"
    />
  </div>
</template>

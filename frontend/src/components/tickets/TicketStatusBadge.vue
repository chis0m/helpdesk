<template>
  <span
    class="shrink-0 rounded-full font-medium capitalize"
    :class="[sizeClass, ticketStatusBadgeClass(status)]"
  >{{ formatTicketStatusLabel(status) }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TicketStatus } from '@/types/ticket'
import { formatTicketStatusLabel, ticketStatusBadgeClass } from '@/utils/ticket-ui'

const props = withDefaults(
  defineProps<{
    status: TicketStatus
    size?: 'sm' | 'md'
  }>(),
  { size: 'sm' },
)

const sizeClass = computed(() =>
  props.size === 'md' ? 'px-3 py-1 text-sm' : 'px-2.5 py-0.5 text-xs',
)
</script>

<template>
  <span
    class="inline-flex max-w-full items-center truncate rounded-md px-2.5 py-1 text-xs font-semibold leading-none shadow-sm ring-1 ring-inset"
    :class="categoryBadgeClass(category)"
    :title="formatTicketCategoryLabel(category)"
  >
    {{ formatTicketCategoryLabel(category) }}
  </span>
</template>

<script setup lang="ts">
import { categoryBadgeClass, formatTicketCategoryLabel } from '@/utils/ticket-ui'

defineProps<{
  category: string
}>()
</script>

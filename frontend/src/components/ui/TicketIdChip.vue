<template>
  <span
    class="inline-flex items-center rounded-lg bg-gradient-to-b from-white to-neutral-100 px-2 py-0.5 font-mono text-[11px] font-bold tabular-nums text-neutral-800 shadow-sm ring-1 ring-inset ring-neutral-200/90"
  >#{{ id }}</span>
</template>

<script setup lang="ts">
defineProps<{
  id: string | number
}>()
</script>

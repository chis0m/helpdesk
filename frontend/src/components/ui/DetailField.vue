<template>
  <div
    class="rounded-xl border border-[var(--border-subtle)] bg-gradient-to-b from-white to-[var(--surface-muted)]/80 px-3 py-2.5 shadow-sm sm:px-4 sm:py-3"
  >
    <p class="text-[10px] font-bold uppercase tracking-[0.14em] text-[var(--text-muted)]">
      {{ label }}
    </p>
    <p class="mt-1 break-all text-sm font-semibold text-[var(--text-primary)]">
      <slot />
    </p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
}>()
</script>

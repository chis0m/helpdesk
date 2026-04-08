<template>
  <header class="space-y-4 pb-1">
    <div
      v-if="kicker"
      class="flex items-center gap-3"
    >
      <span
        class="inline-flex items-center rounded-full bg-gradient-to-r from-[var(--surface-mint)] to-[var(--surface-mint-hover)] px-3 py-1 text-[10px] font-bold uppercase tracking-[0.16em] text-[var(--brand-green-dark)] shadow-sm ring-1 ring-[var(--brand-green)]/20"
      >
        {{ kicker }}
      </span>
      <span
        class="hidden h-px flex-1 max-w-[120px] bg-gradient-to-r from-[var(--brand-green)]/35 to-transparent sm:block"
        aria-hidden="true"
      />
    </div>
    <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div class="min-w-0">
        <h1 class="text-2xl font-bold tracking-tight text-[var(--text-primary)] sm:text-[1.75rem] sm:leading-tight">
          {{ title }}
        </h1>
        <p
          v-if="description"
          class="mt-2 max-w-2xl text-[15px] leading-relaxed text-[var(--text-secondary)]"
        >
          {{ description }}
        </p>
      </div>
      <div
        v-if="$slots.actions"
        class="shrink-0 sm:pb-0.5"
      >
        <slot name="actions" />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  /** Short eyebrow label (e.g. "Overview"). */
  kicker?: string
  title: string
  description?: string
}>()
</script>

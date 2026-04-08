<template>
  <a
    :href="href"
    class="hd-motion-underline flex items-center gap-3 rounded-xl border border-transparent py-2.5 pl-3 pr-3 text-sm font-semibold"
    :class="isLinkActive
      ? 'border-[var(--brand-green)]/25 bg-[var(--surface-mint)]/70 text-[var(--brand-green-dark)] shadow-sm ring-1 ring-[var(--brand-green)]/15'
      : 'text-[var(--text-secondary)] hover:border-[var(--border-subtle)] hover:bg-white hover:text-[var(--text-primary)] hover:shadow-sm'"
    @click="onClick"
  >
    <component
      :is="icon"
      class="h-5 w-5 shrink-0"
      :class="isLinkActive ? 'text-[var(--brand-green-dark)]' : 'opacity-85'"
    />
    {{ label }}
  </a>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { RouteLocationRaw } from 'vue-router'
import type { Component } from 'vue'

const props = withDefaults(
  defineProps<{
    to: RouteLocationRaw
    label: string
    icon: Component
    /** Highlight only when the path equals this link (not child paths). */
    exact?: boolean
    /** Highlight when `route.path` equals this or starts with `prefix + '/'`. */
    activePrefix?: string
  }>(),
  { exact: false, activePrefix: undefined },
)

const route = useRoute()
const router = useRouter()

const resolvedPath = computed(() => router.resolve(props.to).path)
const href = computed(() => router.resolve(props.to).href)

const isLinkActive = computed(() => {
  if (props.activePrefix) {
    const p = props.activePrefix
    return route.path === p || route.path.startsWith(`${p}/`)
  }
  const base = resolvedPath.value
  if (props.exact) return route.path === base
  return route.path === base || route.path.startsWith(`${base}/`)
})

function onClick(e: MouseEvent) {
  if (e.metaKey || e.ctrlKey || e.shiftKey || e.altKey || e.button !== 0) return
  e.preventDefault()
  router.push(props.to)
}
</script>

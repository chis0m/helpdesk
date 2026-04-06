<template>
  <a
    :href="href"
    class="flex items-center gap-3 rounded-xl px-3 py-2 text-sm font-medium transition"
    :class="isLinkActive
      ? 'bg-[var(--surface-active)] text-[var(--text-primary)]'
      : 'text-[var(--text-secondary)] hover:bg-[var(--surface-hover)]'"
    @click="onClick"
  >
    <component
      :is="icon"
      class="h-5 w-5 shrink-0 opacity-80"
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
  { exact: false },
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

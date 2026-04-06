<template>
  <nav
    class="flex flex-wrap gap-2 border-b border-[var(--border-subtle)] pb-3"
    aria-label="Admin sections"
  >
    <RouterLink
      v-for="link in links"
      :key="link.to"
      :to="link.to"
      class="rounded-full px-4 py-2 text-sm font-medium transition"
      :class="isActive(link)
        ? 'bg-[var(--surface-active)] text-[var(--text-primary)]'
        : 'text-[var(--text-secondary)] hover:bg-[var(--surface-hover)]'"
    >
      {{ link.label }}
    </RouterLink>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { paths } from '@/constants/routes'

const route = useRoute()

const links = computed(() => [
  { to: paths.dashboard.adminUsers, label: 'Users' },
  { to: paths.dashboard.adminStaff, label: 'Staff' },
  { to: paths.dashboard.adminStaffNew, label: 'Create staff' },
])

function isActive(link: { to: string }) {
  return route.path === link.to
}
</script>

<template>
  <div class="overflow-hidden rounded-2xl border border-[var(--border-subtle)] bg-white shadow-sm">
    <div class="border-b border-[var(--border-subtle)] px-4 py-3">
      <p class="text-sm font-medium text-[var(--text-primary)]">
        {{ countLabel }}
      </p>
    </div>
    <div
      v-if="users.length > 0"
      class="overflow-x-auto"
    >
      <table class="w-full min-w-[640px] text-left text-sm">
        <thead>
          <tr class="border-b border-[var(--border-subtle)] bg-[var(--surface-muted)]/60 text-xs font-semibold uppercase tracking-wide text-[var(--text-muted)]">
            <th class="px-4 py-3 font-semibold">
              Name
            </th>
            <th class="px-4 py-3 font-semibold">
              Email
            </th>
            <th class="px-4 py-3 font-semibold">
              ID
            </th>
            <th
              v-if="showOrgColumn"
              class="px-4 py-3 font-semibold"
            >
              Organization
            </th>
            <th class="px-4 py-3 font-semibold">
              Created
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="u in users"
            :key="u.id"
            class="border-b border-[var(--border-subtle)] last:border-0"
          >
            <td class="px-4 py-3 font-medium text-[var(--text-primary)]">
              {{ u.displayName }}
            </td>
            <td class="px-4 py-3 text-[var(--text-secondary)]">
              {{ u.email }}
            </td>
            <td class="px-4 py-3 font-mono text-xs text-[var(--text-muted)]">
              {{ u.id }}
            </td>
            <td
              v-if="showOrgColumn"
              class="px-4 py-3 text-[var(--text-secondary)]"
            >
              {{ u.organization ?? '—' }}
            </td>
            <td class="px-4 py-3 text-[var(--text-muted)]">
              {{ formatDateTime(u.createdAt) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <p
      v-else
      class="px-4 py-8 text-center text-sm text-[var(--text-secondary)]"
    >
      No results.
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { PortalUser } from '@/types/directory-user'
import { formatDateTime } from '@/utils/date-format'

const props = defineProps<{
  users: PortalUser[]
  /** When true, show organization column (typically for customer users). */
  showOrganization?: boolean
}>()

const showOrgColumn = computed(() => props.showOrganization === true)

const countLabel = computed(() => {
  const n = props.users.length
  return `${n} account${n === 1 ? '' : 's'}`
})
</script>

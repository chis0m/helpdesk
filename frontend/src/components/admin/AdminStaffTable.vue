<template>
  <div class="overflow-hidden rounded-2xl border border-[var(--border-subtle)] bg-white shadow-sm">
    <div class="border-b border-[var(--border-subtle)] px-4 py-3">
      <p class="text-sm font-medium text-[var(--text-primary)]">
        {{ countLabel }}
      </p>
    </div>
    <div
      v-if="staff.length > 0"
      class="overflow-x-auto"
    >
      <table class="w-full min-w-[720px] text-left text-sm">
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
            <th class="px-4 py-3 font-semibold">
              Administrator
            </th>
            <th class="px-4 py-3 font-semibold">
              Created
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="s in staff"
            :key="s.id"
            class="border-b border-[var(--border-subtle)] last:border-0"
          >
            <td class="px-4 py-3 font-medium text-[var(--text-primary)]">
              {{ s.displayName }}
            </td>
            <td class="px-4 py-3 text-[var(--text-secondary)]">
              {{ s.email }}
            </td>
            <td class="px-4 py-3 font-mono text-xs text-[var(--text-muted)]">
              {{ s.id }}
            </td>
            <td class="px-4 py-3">
              <span
                class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-semibold"
                :class="s.isAdmin
                  ? 'bg-violet-100 text-violet-900'
                  : 'bg-neutral-100 text-neutral-700'"
              >{{ s.isAdmin ? 'Yes' : 'No' }}</span>
            </td>
            <td class="px-4 py-3 text-[var(--text-muted)]">
              {{ formatDateTime(s.createdAt) }}
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
import type { StaffMember } from '@/types/directory-user'
import { formatDateTime } from '@/utils/date-format'

const props = defineProps<{
  staff: StaffMember[]
}>()

const countLabel = computed(() => {
  const n = props.staff.length
  return `${n} staff member${n === 1 ? '' : 's'}`
})
</script>

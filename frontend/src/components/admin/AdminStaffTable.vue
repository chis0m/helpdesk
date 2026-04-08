<template>
  <div class="overflow-hidden rounded-2xl border border-[var(--border-subtle)] bg-white shadow-sm">
    <div class="border-b border-[var(--border-subtle)] px-4 py-3">
      <p class="text-sm font-medium text-[var(--text-primary)]">
        {{ countLabel }}
      </p>
      <p
        v-if="roleError"
        class="mt-2 rounded-xl border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-900"
        role="alert"
      >
        {{ roleError }}
      </p>
    </div>
    <div
      v-if="staff.length > 0"
      class="overflow-x-auto"
    >
      <table class="w-full min-w-[780px] text-left text-sm">
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
              Role
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
              <select
                v-if="canEditRole(s)"
                class="max-w-[220px] rounded-xl border border-[var(--border-subtle)] bg-white px-2 py-1.5 text-sm text-[var(--text-primary)] outline-none ring-[var(--brand-green)] focus:border-transparent focus:ring-2 disabled:cursor-not-allowed disabled:opacity-60"
                :value="s.role"
                :disabled="savingId === s.id"
                :aria-label="`Role for ${s.email}`"
                @change="onRoleSelectChange(s, $event)"
              >
                <option
                  v-for="opt in roleOptionsForRow(s)"
                  :key="opt.value"
                  :value="opt.value"
                >
                  {{ opt.label }}
                </option>
              </select>
              <span
                v-else
                class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-semibold"
                :class="roleBadgeClass(s.role)"
              >{{ roleLabel(s.role) }}</span>
              <span
                v-if="!canEditRole(s) && isOwnRow(s)"
                class="mt-1 block text-[11px] text-[var(--text-muted)]"
              >You cannot change your own role here.</span>
              <span
                v-else-if="!canEditRole(s) && s.role === 'staff' && actorRole !== 'super_admin'"
                class="mt-1 block text-[11px] text-[var(--text-muted)]"
              >Only a super administrator can promote staff to an administrator role.</span>
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
import { computed, ref } from 'vue'
import type { StaffMember } from '@/types/directory-user'
import type { AdminDirectoryRole } from '@/api/admin'
import { patchAdminUserRole } from '@/api/admin'
import { getSessionCsrfToken } from '@/stores/auth-session'
import { formatDateTime } from '@/utils/date-format'

const props = defineProps<{
  staff: StaffMember[]
  /** Current viewer role from `getAuthUserSnapshot().role` */
  actorRole: string | null | undefined
  /** Current viewer user id; own row is read-only to avoid accidental self-demotion */
  actorUserId: number | null | undefined
}>()

const emit = defineEmits<{
  updated: []
}>()

const savingId = ref<string | null>(null)
const roleError = ref('')

const countLabel = computed(() => {
  const n = props.staff.length
  return `${n} staff member${n === 1 ? '' : 's'}`
})

function roleLabel(role: AdminDirectoryRole): string {
  switch (role) {
    case 'user':
      return 'User'
    case 'staff':
      return 'Staff'
    case 'admin':
      return 'Administrator'
    case 'super_admin':
      return 'Super administrator'
    default:
      return role
  }
}

function roleBadgeClass(role: AdminDirectoryRole): string {
  if (role === 'super_admin')
    return 'bg-amber-100 text-amber-950'
  if (role === 'admin')
    return 'bg-violet-100 text-violet-900'
  return 'bg-neutral-100 text-neutral-700'
}

function isOwnRow(s: StaffMember): boolean {
  const id = props.actorUserId
  if (id == null)
    return false
  return Number(s.id) === id
}

/** Only super administrators may promote staff; staff may only become admin or super_admin (not customers). */
function canEditRole(s: StaffMember): boolean {
  if (props.actorRole !== 'super_admin')
    return false
  if (isOwnRow(s))
    return false
  return s.role === 'staff'
}

function roleOptionsForRow(s: StaffMember): { value: AdminDirectoryRole; label: string }[] {
  if (props.actorRole === 'super_admin' && s.role === 'staff') {
    return [
      { value: 'admin', label: 'Administrator' },
      { value: 'super_admin', label: 'Super administrator' },
    ]
  }
  return [{ value: s.role, label: roleLabel(s.role) }]
}

function onRoleSelectChange(s: StaffMember, ev: Event) {
  const el = ev.target as HTMLSelectElement | null
  if (!el)
    return
  void onRoleChange(s, el.value as AdminDirectoryRole)
}

async function onRoleChange(s: StaffMember, newRole: AdminDirectoryRole) {
  if (newRole === s.role)
    return
  roleError.value = ''
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    roleError.value = 'Your session is missing a security token. Sign out and sign in again, then retry.'
    return
  }
  savingId.value = s.id
  try {
    const res = await patchAdminUserRole(Number(s.id), newRole, csrf)
    if (!res.ok) {
      roleError.value = res.message
      return
    }
    emit('updated')
  }
  finally {
    savingId.value = null
  }
}
</script>

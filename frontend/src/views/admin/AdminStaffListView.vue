<template>
  <div class="space-y-6">
    <header>
      <p class="text-xs font-semibold uppercase tracking-wider text-[var(--brand-green-dark)]">
        Admin
      </p>
      <h1 class="mt-2 text-xl font-semibold tracking-tight text-[var(--text-primary)]">
        Staff
      </h1>
      <p class="mt-2 max-w-2xl text-sm leading-relaxed text-[var(--text-secondary)]">
        Everyone on the internal team. Super administrators can promote <strong class="font-medium text-[var(--text-primary)]">staff</strong>
        to administrator or super administrator. Customers and internal roles are not interchangeable here.
      </p>
    </header>

    <AdminSubnav />

    <div
      v-if="loadError"
      class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-900"
      role="alert"
    >
      {{ loadError }}
    </div>

    <div
      v-else-if="loading"
      class="rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-10 text-center text-sm text-[var(--text-secondary)]"
    >
      Loading staff…
    </div>

    <AdminStaffTable
      v-else
      :staff="staff"
      :actor-role="actorRole"
      :actor-user-id="actorUserId"
      @updated="loadStaff"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { adminUserToStaffMember, fetchAdminStaffDirectory } from '@/api/admin'
import AdminSubnav from '@/components/admin/AdminSubnav.vue'
import AdminStaffTable from '@/components/admin/AdminStaffTable.vue'
import { getAuthUserSnapshot } from '@/stores/auth-session'
import type { StaffMember } from '@/types/directory-user'

const staff = ref<StaffMember[]>([])
const loading = ref(true)
const loadError = ref('')

const actorRole = computed(() => getAuthUserSnapshot()?.role ?? null)
const actorUserId = computed(() => getAuthUserSnapshot()?.user_id ?? null)

async function loadStaff() {
  loading.value = true
  loadError.value = ''
  const result = await fetchAdminStaffDirectory()
  loading.value = false
  if (!result.ok) {
    loadError.value = result.message
    return
  }
  staff.value = result.items.map(adminUserToStaffMember)
}

onMounted(() => {
  void loadStaff()
})
</script>

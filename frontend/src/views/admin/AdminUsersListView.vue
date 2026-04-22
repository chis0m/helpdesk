<template>
  <div class="space-y-6">
    <header>
      <p class="text-xs font-semibold uppercase tracking-wider text-[var(--brand-green-dark)]">
        Admin
      </p>
      <h1 class="mt-2 text-xl font-semibold tracking-tight text-[var(--text-primary)]">
        Users
      </h1>
      <p class="mt-2 max-w-2xl text-sm leading-relaxed text-[var(--text-secondary)]">
        Customer and portal accounts that use {{ appName }}.
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
      Loading users…
    </div>

    <AdminUserTable
      v-else
      :users="users"
      show-organization
    />
  </div>
</template>

<script setup lang="ts">
import { computed, inject, onMounted, ref } from 'vue'
import { adminUserToPortalUser, fetchAdminUsers } from '@/api/admin'
import AdminSubnav from '@/components/admin/AdminSubnav.vue'
import AdminUserTable from '@/components/admin/AdminUserTable.vue'
import { appBrandKey } from '@/stores/app-detail'
import type { PortalUser } from '@/types/directory-user'

const brand = inject(appBrandKey, null)
const appName = computed(() => brand?.appName.value ?? 'SecWeb HelpDesk')

const users = ref<PortalUser[]>([])
const loading = ref(true)
const loadError = ref('')

onMounted(async () => {
  loading.value = true
  loadError.value = ''
  const result = await fetchAdminUsers({ role: 'user', limit: 100, page: 1 })
  loading.value = false
  if (!result.ok) {
    loadError.value = result.message
    return
  }
  users.value = result.data.items.map(adminUserToPortalUser)
})
</script>

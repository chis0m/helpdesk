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
        Customer and portal accounts that use SecWeb Helpdesk.
      </p>
    </header>

    <AdminSubnav />

    <AdminUserTable
      :users="MOCK_PORTAL_USERS"
      show-organization
    />
  </div>
</template>

<script setup lang="ts">
import AdminSubnav from '@/components/admin/AdminSubnav.vue'
import AdminUserTable from '@/components/admin/AdminUserTable.vue'
import { MOCK_PORTAL_USERS } from '@/data/mocks/admin-directory.mock'
</script>

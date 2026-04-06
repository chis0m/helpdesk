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
        Everyone on the internal team.
        <strong class="font-medium text-[var(--text-primary)]">Administrator</strong>
        shows who has elevated admin access.
      </p>
    </header>

    <AdminSubnav />

    <AdminStaffTable :staff="MOCK_STAFF_MEMBERS" />
  </div>
</template>

<script setup lang="ts">
import AdminSubnav from '@/components/admin/AdminSubnav.vue'
import AdminStaffTable from '@/components/admin/AdminStaffTable.vue'
import { MOCK_STAFF_MEMBERS } from '@/data/mocks/admin-directory.mock'
</script>

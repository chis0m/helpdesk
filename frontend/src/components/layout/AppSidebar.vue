<template>
  <aside
    class="flex w-[260px] shrink-0 flex-col border-r border-[var(--border-subtle)] bg-[var(--surface-sidebar)] px-3 py-4"
  >
    <RouterLink
      to="/"
      class="mb-5 flex items-center gap-2 px-3"
    >
      <span
        class="text-xl font-semibold tracking-tight text-[var(--brand-green-dark)]"
      >SecWeb Helpdesk</span>
    </RouterLink>

    <nav class="flex flex-col gap-0.5">
      <SidebarLink
        :to="paths.dashboard.home"
        :icon="IconHome"
        label="Support home"
        exact
      />
      <SidebarLink
        :to="paths.dashboard.tickets"
        :icon="IconTicket"
        label="Tickets"
        :active-prefix="paths.dashboard.tickets"
      />
      <SidebarLink
        :to="paths.dashboard.profile"
        :icon="IconUser"
        label="Profile"
      />
      <div class="mt-2 border-t border-[var(--border-subtle)] pt-2">
        <p class="mb-2 px-3 text-xs font-semibold uppercase tracking-wide text-[var(--text-muted)]">
          Admin
        </p>
        <div class="flex flex-col gap-0.5">
          <SidebarLink
            :to="paths.dashboard.adminUsers"
            :icon="IconGear"
            label="Users"
            exact
          />
          <SidebarLink
            :to="paths.dashboard.adminStaff"
            :icon="IconGear"
            label="Staff"
            exact
          />
          <SidebarLink
            :to="paths.dashboard.adminStaffNew"
            :icon="IconGear"
            label="Create staff"
            exact
          />
        </div>
      </div>
    </nav>

    <div class="mt-auto border-t border-[var(--border-subtle)] pt-3">
      <RouterLink
        to="/login"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left text-sm font-medium text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)]"
      >
        Log out
      </RouterLink>
    </div>
  </aside>
</template>

<script setup lang="ts">
import SidebarLink from '@/components/layout/SidebarLink.vue'
import IconHome from '@/components/icons/IconHome.vue'
import IconTicket from '@/components/icons/IconTicket.vue'
import IconUser from '@/components/icons/IconUser.vue'
import IconGear from '@/components/icons/IconGear.vue'
import { paths } from '@/constants/routes'
</script>

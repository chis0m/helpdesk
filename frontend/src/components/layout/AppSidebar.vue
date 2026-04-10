<template>
  <aside
    class="flex w-[260px] shrink-0 flex-col border-r border-[var(--border-subtle)] bg-[var(--surface-sidebar)] px-3 py-5 shadow-[2px_0_12px_rgba(15,23,42,0.03)]"
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
        :to="profilePath"
        :icon="IconUser"
        label="Profile"
        active-prefix="/dashboard/profile"
      />
      <SidebarLink
        :to="paths.dashboard.sessions"
        :icon="IconGear"
        label="Sessions"
        exact
      />
      <div
        v-if="showAdminNav"
        class="mt-2 border-t border-[var(--border-subtle)] pt-2"
      >
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
      <button
        type="button"
        :disabled="loggingOut"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left text-sm font-medium text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] disabled:cursor-not-allowed disabled:opacity-60"
        @click="onLogout"
      >
        {{ loggingOut ? 'Signing out…' : 'Log out' }}
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SidebarLink from '@/components/layout/SidebarLink.vue'
import IconHome from '@/components/icons/IconHome.vue'
import IconTicket from '@/components/icons/IconTicket.vue'
import IconUser from '@/components/icons/IconUser.vue'
import IconGear from '@/components/icons/IconGear.vue'
import { paths } from '@/constants/routes'
import { getAuthUserSnapshot } from '@/stores/auth-session'
import { isAdminPortalRole } from '@/utils/admin-access'
import { performLogout } from '@/utils/perform-logout'

const router = useRouter()
const route = useRoute()
const loggingOut = ref(false)

async function onLogout() {
  if (loggingOut.value)
    return
  loggingOut.value = true
  try {
    await performLogout(router)
  }
  finally {
    loggingOut.value = false
  }
}

const profilePath = computed(() => {
  const u = getAuthUserSnapshot()
  return u ? paths.dashboard.profile(u.user_id) : paths.login
})

/** `route` keeps this in sync after login/navigation (session snapshot is not reactive). */
const showAdminNav = computed(() => {
  void route.fullPath
  const u = getAuthUserSnapshot()
  return isAdminPortalRole(u?.role ?? null)
})
</script>

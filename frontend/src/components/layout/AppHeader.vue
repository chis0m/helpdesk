<template>
  <header
    class="sticky top-0 z-20 flex h-[56px] shrink-0 items-center justify-between border-b border-[var(--border-subtle)] bg-white/90 px-5 shadow-[var(--shadow-header)] backdrop-blur-md supports-[backdrop-filter]:bg-white/75 lg:px-8"
  >
    <div class="min-w-0">
      <h1 class="truncate text-lg font-bold tracking-tight text-[var(--text-primary)] lg:text-xl">
        {{ title }}
      </h1>
    </div>

    <div class="flex items-center gap-2 sm:gap-3">
      <RouterLink
        :to="paths.dashboard.tickets"
        class="hidden rounded-full border border-[var(--border-subtle)] px-3 py-2 text-xs font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] sm:inline"
      >
        Tickets
      </RouterLink>
      <span
        class="hidden rounded-full bg-[var(--brand-green)] px-3 py-2 text-xs font-semibold text-[var(--text-on-green)] shadow-sm sm:inline"
      >Product support</span>
      <button
        type="button"
        :disabled="loggingOut"
        class="rounded-full border border-[var(--border-subtle)] px-3 py-2 text-xs font-semibold text-[var(--text-secondary)] transition hover:bg-[var(--surface-hover)] disabled:cursor-not-allowed disabled:opacity-60 sm:px-4 sm:text-sm"
        @click="onLogout"
      >
        {{ loggingOut ? 'Signing out…' : 'Log out' }}
      </button>
      <RouterLink
        :to="profilePath"
        class="flex items-center gap-2 rounded-full py-1.5 pl-1.5 pr-2 transition hover:bg-[var(--surface-hover)]"
      >
        <span
          class="flex h-9 w-9 items-center justify-center rounded-full bg-[var(--surface-active)] text-sm font-semibold text-[var(--text-primary)]"
        >{{ initials }}</span>
        <span class="hidden max-w-[160px] truncate text-sm font-medium text-[var(--text-primary)] sm:inline">{{ displayName }}</span>
        <svg
          class="h-4 w-4 text-[var(--text-muted)]"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </RouterLink>
    </div>
  </header>
</template>

<script setup lang="ts">
// VULN-02: Profile link targets `/dashboard/profile/${session user_id}` — user can navigate to any id manually (IDOR demo).
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { paths } from '@/constants/routes'
import { getAuthUserSnapshot } from '@/stores/auth-session'
import { displayFromEmail, initialsFromEmail } from '@/utils/user-display'
import { performLogout } from '@/utils/perform-logout'

defineProps<{
  title: string
}>()

const router = useRouter()
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
  return u ? paths.dashboard.profile : paths.login
})

const displayName = computed(() => {
  const u = getAuthUserSnapshot()
  return u ? displayFromEmail(u.email) : 'Account'
})

const initials = computed(() => {
  const u = getAuthUserSnapshot()
  return u ? initialsFromEmail(u.email) : '?'
})
</script>

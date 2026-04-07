<template>
  <!-- VULN-01: Session list requires authenticated cookies. VULN-05: Revoke actions send `X-CSRF-Token`. -->
  <div class="mx-auto max-w-2xl space-y-6">
    <header>
      <h2 class="text-lg font-semibold text-[var(--text-primary)]">
        Active sessions
      </h2>
      <p class="mt-1 text-sm text-[var(--text-secondary)]">
        Devices signed in to your account. Revoke any you don’t recognize.
      </p>
    </header>

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
      Loading sessions…
    </div>

    <template v-else>
      <div
        v-if="actionError"
        class="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-950"
        role="alert"
      >
        {{ actionError }}
      </div>

      <div class="flex flex-wrap gap-2">
        <button
          type="button"
          class="rounded-full border border-[var(--border-strong)] bg-white px-4 py-2 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--surface-hover)] disabled:opacity-50"
          :disabled="revokingOthers || items.length <= 1"
          @click="onRevokeOthers"
        >
          {{ revokingOthers ? 'Signing out…' : 'Sign out other sessions' }}
        </button>
      </div>

      <ul class="space-y-2">
        <li
          v-for="s in items"
          :key="s.session_id"
          class="rounded-xl border border-[var(--border-subtle)] bg-white px-4 py-3 shadow-sm"
        >
          <div class="flex flex-wrap items-start justify-between gap-2">
            <div class="min-w-0">
              <p class="font-medium text-[var(--text-primary)]">
                {{ s.user_agent ?? 'Unknown client' }}
              </p>
              <p class="mt-0.5 text-xs text-[var(--text-muted)]">
                {{ s.ip ?? '—' }} · {{ formatDateTime(s.created_at) }}
              </p>
              <p
                v-if="s.is_current"
                class="mt-1 text-xs font-semibold text-[var(--brand-green-dark)]"
              >
                This device
              </p>
            </div>
            <button
              v-if="!s.is_current"
              type="button"
              class="shrink-0 rounded-full border border-red-200 bg-red-50 px-3 py-1.5 text-xs font-semibold text-red-900 transition hover:bg-red-100 disabled:opacity-50"
              :disabled="revokingId === s.session_id"
              @click="onRevokeOne(s.session_id)"
            >
              {{ revokingId === s.session_id ? '…' : 'Revoke' }}
            </button>
          </div>
        </li>
      </ul>

      <p
        v-if="items.length === 0"
        class="text-sm text-[var(--text-secondary)]"
      >
        No other sessions.
      </p>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  fetchAuthSessions,
  revokeAuthSession,
  revokeMyOtherSessionsRequest,
  type AuthSessionRow,
} from '@/api/auth'
import { clearSessionRefreshSchedule } from '@/api/session-refresh'
import { paths } from '@/constants/routes'
import { clearAuthSession, getSessionCsrfToken } from '@/stores/auth-session'
import { formatDateTime } from '@/utils/date-format'

const router = useRouter()

const loading = ref(true)
const loadError = ref('')
const items = ref<AuthSessionRow[]>([])
const actionError = ref('')
const revokingId = ref<string | null>(null)
const revokingOthers = ref(false)

async function load() {
  loadError.value = ''
  loading.value = true
  const res = await fetchAuthSessions()
  loading.value = false
  if (!res.ok) {
    loadError.value = res.message
    items.value = []
    return
  }
  items.value = res.items
}

onMounted(() => void load())

function csrfOrBail(): string | null {
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    actionError.value = 'Your session expired. Sign in again.'
    return null
  }
  return csrf
}

async function onRevokeOne(sessionId: string) {
  actionError.value = ''
  const csrf = csrfOrBail()
  if (!csrf)
    return
  revokingId.value = sessionId
  const res = await revokeAuthSession(sessionId, csrf)
  revokingId.value = null
  if (!res.ok) {
    actionError.value = res.message
    return
  }
  if (res.data.logged_out) {
    clearSessionRefreshSchedule()
    clearAuthSession()
    await router.replace(paths.login)
    return
  }
  await load()
}

async function onRevokeOthers() {
  actionError.value = ''
  const csrf = csrfOrBail()
  if (!csrf)
    return
  revokingOthers.value = true
  const res = await revokeMyOtherSessionsRequest(csrf)
  revokingOthers.value = false
  if (!res.ok) {
    actionError.value = res.message
    return
  }
  await load()
}
</script>

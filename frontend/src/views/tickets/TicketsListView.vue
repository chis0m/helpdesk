<template>
  <!-- VULN-07: Search calls `GET /api/tickets/search?q=` — `q` is passed to unsafe SQL on the server (see backend TicketController.Search). -->
  <div class="space-y-4">
    <header class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-base font-semibold text-[var(--text-primary)] sm:text-lg">
          Tickets & search
        </h2>
        <p class="mt-0.5 text-xs text-[var(--text-secondary)] sm:text-sm">
          List your tickets or search keywords (server-side search for coursework demos).
        </p>
      </div>
      <RouterLink
        :to="paths.dashboard.ticketNew"
        class="inline-flex shrink-0 items-center justify-center rounded-full bg-[var(--brand-green)] px-4 py-2 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95"
      >
        New ticket
      </RouterLink>
    </header>

    <div class="relative">
      <span class="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-[var(--text-muted)]">
        <svg
          class="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          stroke-width="1.5"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
          />
        </svg>
      </span>
      <input
        v-model="query"
        type="search"
        placeholder="Search keywords (empty = list your tickets)…"
        class="w-full rounded-full border border-[var(--border-subtle)] bg-white py-2.5 pl-11 pr-4 text-sm text-[var(--text-primary)] shadow-sm outline-none ring-[var(--brand-green)] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2"
      >
    </div>

    <div
      v-if="loadError"
      class="rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-950"
      role="alert"
    >
      {{ loadError }}
    </div>

    <div class="overflow-hidden rounded-xl border border-[var(--border-subtle)] bg-white">
      <p class="border-b border-[var(--border-subtle)] px-3 py-2 text-xs font-medium text-[var(--text-muted)]">
        <span v-if="loading">Loading…</span>
        <span v-else>{{ rows.length }} result{{ rows.length === 1 ? '' : 's' }}{{ searchEcho ? ` for “${searchEcho}”` : '' }}</span>
      </p>
      <ul>
        <li
          v-for="(t, idx) in rows"
          :key="t.ticket_id"
        >
          <div
            v-if="idx > 0"
            class="h-px bg-[var(--border-subtle)]"
          />
          <RouterLink
            :to="paths.dashboard.ticketDetail(String(t.ticket_id))"
            class="flex items-center gap-3 px-3 py-3 transition hover:bg-[var(--surface-hover)]"
          >
            <div
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-[var(--surface-muted)] text-xs font-medium text-[var(--text-secondary)] sm:h-10 sm:w-10 sm:text-sm"
            >
              #{{ t.ticket_id }}
            </div>
            <div class="min-w-0 flex-1">
              <p class="font-medium text-[var(--text-primary)]">
                {{ t.title }}
              </p>
              <p class="text-sm text-[var(--text-muted)]">
                {{ formatTicketCategoryLabel(t.category) }} · {{ reporterLabel(t) }}
              </p>
            </div>
            <TicketStatusBadge :status="t.status" />
            <svg
              class="h-5 w-5 shrink-0 text-[var(--text-muted)]"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M9 5l7 7-7 7"
              />
            </svg>
          </RouterLink>
        </li>
      </ul>
      <p
        v-if="!loading && rows.length === 0"
        class="px-3 py-7 text-center text-sm text-[var(--text-secondary)]"
      >
        {{ emptyMessage }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
// VULN-07: Uses `fetchTicketSearch` / `fetchTicketList` — search path exercises unsafe SQL on the server.
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import TicketStatusBadge from '@/components/tickets/TicketStatusBadge.vue'
import { paths } from '@/constants/routes'
import { fetchTicketList, fetchTicketSearch, type ApiTicketRow } from '@/api/tickets'
import { getAuthUserSnapshot } from '@/stores/auth-session'
import { formatTicketCategoryLabel, reporterDisplayLabel } from '@/utils/ticket-ui'

const query = ref('')
const rows = ref<ApiTicketRow[]>([])
const loading = ref(false)
const loadError = ref('')
const searchEcho = ref('')

const emptyMessage = computed(() => {
  const q = query.value.trim()
  if (q === '')
    return 'No tickets yet.'
  return `No tickets match “${q}”.`
})

function reporterLabel(t: ApiTicketRow): string {
  const s = getAuthUserSnapshot()
  return reporterDisplayLabel(t, s?.user_id ?? null)
}

let debounceTimer: ReturnType<typeof setTimeout> | null = null

async function runLoad() {
  const q = query.value.trim()
  loading.value = true
  loadError.value = ''
  searchEcho.value = ''
  if (q === '') {
    const res = await fetchTicketList({ page: 1, limit: 100 })
    loading.value = false
    if (!res.ok) {
      loadError.value = res.message
      rows.value = []
      return
    }
    rows.value = res.items
    return
  }
  const res = await fetchTicketSearch(q)
  loading.value = false
  if (!res.ok) {
    loadError.value = res.message
    rows.value = []
    return
  }
  rows.value = res.items
  searchEcho.value = res.queryEcho
}

function scheduleLoad() {
  if (debounceTimer)
    clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => void runLoad(), 320)
}

onMounted(() => void runLoad())

watch(
  () => query.value,
  () => scheduleLoad(),
)

onUnmounted(() => {
  if (debounceTimer)
    clearTimeout(debounceTimer)
})
</script>

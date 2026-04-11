<template>
  <!-- VULN-07: Search calls `GET /api/tickets/search?q=` — `q` is passed to unsafe SQL on the server (see backend TicketController.Search). -->
  <div class="space-y-6">
    <PageHeader
      kicker="Tickets"
      title="Find & manage requests"
      description="Search by keyword or browse everything you can access. Newest activity appears first in the list."
    >
      <template #actions>
        <RouterLink
          :to="paths.dashboard.ticketNew"
          class="inline-flex items-center justify-center rounded-xl bg-[var(--brand-green)] px-5 py-2.5 text-sm font-bold text-[var(--text-on-green)] shadow-[var(--shadow-card)] transition hover:brightness-[1.02] active:scale-[0.98]"
        >
          + New ticket
        </RouterLink>
      </template>
    </PageHeader>

    <div>
      <label
        for="ticket-search"
        class="mb-2 block text-xs font-bold uppercase tracking-[0.12em] text-[var(--text-muted)]"
      >Search</label>
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
          id="ticket-search"
          v-model="query"
          type="search"
          placeholder="Keywords… (leave empty to list all your tickets)"
          class="w-full rounded-xl border border-[var(--border-subtle)] bg-white py-3 pl-11 pr-4 text-sm font-medium text-[var(--text-primary)] shadow-sm outline-none ring-[var(--brand-green)] transition placeholder:font-normal placeholder:text-[var(--text-muted)] focus:border-transparent focus:shadow-[var(--shadow-card)] focus:ring-2"
        >
      </div>
    </div>

    <div
      v-if="loadError"
      class="rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm font-medium text-amber-950 shadow-sm"
      role="alert"
    >
      {{ loadError }}
    </div>

    <div class="overflow-hidden rounded-2xl border border-[var(--border-subtle)] bg-white shadow-[var(--shadow-card)]">
      <div class="flex flex-wrap items-center justify-between gap-2 border-b border-[var(--border-subtle)] bg-gradient-to-r from-[var(--surface-muted)]/90 to-white px-4 py-3.5">
        <span class="text-xs font-bold uppercase tracking-wider text-[var(--text-muted)]">Results</span>
        <span
          v-if="loading"
          class="rounded-full bg-white px-2.5 py-0.5 text-xs font-semibold text-[var(--text-secondary)] ring-1 ring-[var(--border-subtle)]"
        >Loading…</span>
        <span
          v-else
          class="rounded-full bg-[var(--brand-green)]/15 px-2.5 py-0.5 text-xs font-bold text-[var(--brand-green-dark)] ring-1 ring-[var(--brand-green)]/30"
        >{{ rows.length }} ticket{{ rows.length === 1 ? '' : 's' }}{{ searchEcho ? ` · “${searchEcho}”` : '' }}</span>
      </div>
      <ul>
        <li
          v-for="(t, idx) in rows"
          :key="t.ticket_uuid"
        >
          <div
            v-if="idx > 0"
            class="h-px bg-[var(--border-subtle)]"
          />
          <RouterLink
            :to="paths.dashboard.ticketDetail(t.ticket_uuid)"
            class="group flex items-start gap-3 px-4 py-4 transition duration-200 hover:bg-gradient-to-r hover:from-[var(--surface-mint)]/25 hover:to-transparent sm:items-center sm:gap-4"
          >
            <div class="shrink-0 pt-0.5">
              <span
                class="inline-flex min-w-[1.75rem] items-center justify-center rounded-lg bg-gradient-to-b from-white to-neutral-100 px-2 py-0.5 font-mono text-[11px] font-bold tabular-nums text-neutral-800 shadow-sm ring-1 ring-inset ring-neutral-200/90"
                :aria-label="`Row ${idx + 1} of ${rows.length}`"
              >{{ idx + 1 }}</span>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-base font-semibold leading-snug text-[var(--text-primary)] group-hover:text-[var(--brand-green-dark)]">
                {{ t.title }}
              </p>
              <div class="mt-2 flex flex-wrap items-center gap-2">
                <CategoryBadge :category="t.category" />
                <span class="text-xs font-medium text-[var(--text-muted)]">·</span>
                <span class="text-xs font-semibold text-[var(--text-secondary)]">{{ reporterLabel(t) }}</span>
              </div>
            </div>
            <div class="flex shrink-0 items-center gap-2">
              <TicketStatusBadge
                :status="t.status"
                size="md"
              />
              <svg
                class="h-5 w-5 text-[var(--text-muted)] transition group-hover:translate-x-0.5 group-hover:text-[var(--brand-green-dark)]"
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
            </div>
          </RouterLink>
        </li>
      </ul>
      <p
        v-if="!loading && rows.length === 0"
        class="px-4 py-12 text-center text-sm font-medium text-[var(--text-secondary)]"
      >
        {{ emptyMessage }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
// VULN-07: Uses `fetchTicketSearch` / `fetchTicketList` — search path exercises unsafe SQL on the server.
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import CategoryBadge from '@/components/ui/CategoryBadge.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import TicketStatusBadge from '@/components/tickets/TicketStatusBadge.vue'
import { paths } from '@/constants/routes'
import { fetchTicketList, fetchTicketSearch, type ApiTicketRow } from '@/api/tickets'
import { getAuthUserSnapshot } from '@/stores/auth-session'
import { reporterDisplayLabel } from '@/utils/ticket-ui'

const query = ref('')
const rows = ref<ApiTicketRow[]>([])
const loading = ref(false)
const loadError = ref('')
const searchEcho = ref('')

const emptyMessage = computed(() => {
  const q = query.value.trim()
  if (q === '')
    return 'No tickets yet — create one to reach support.'
  return `No tickets match “${q}”. Try different keywords.`
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

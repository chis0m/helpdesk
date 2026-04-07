<template>
  <div class="space-y-7">
    <div
      v-if="loadError"
      class="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-950"
      role="alert"
    >
      {{ loadError }}
    </div>

    <!-- Hero metric + actions -->
    <section class="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <p class="text-sm font-medium text-[var(--text-secondary)]">
          Your open requests
        </p>
        <p class="mt-1 flex items-baseline gap-2 text-3xl font-semibold tracking-tight text-[var(--text-primary)] lg:text-4xl">
          {{ loading ? '…' : openTicketCount }}
          <span class="text-base font-normal text-[var(--text-muted)] lg:text-lg">with SecWeb support</span>
        </p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <RouterLink
          :to="paths.dashboard.ticketNew"
          class="inline-flex items-center justify-center rounded-full bg-[var(--brand-green)] px-5 py-2.5 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95"
        >
          New ticket
        </RouterLink>
        <RouterLink
          :to="paths.dashboard.tickets"
          class="inline-flex items-center justify-center rounded-full bg-[var(--surface-mint)] px-5 py-2.5 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--surface-mint-hover)]"
        >
          My requests
        </RouterLink>
        <RouterLink
          :to="paths.dashboard.tickets"
          class="inline-flex items-center justify-center rounded-full border border-[var(--border-strong)] bg-white px-5 py-2.5 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--surface-hover)]"
        >
          Search
        </RouterLink>
      </div>
    </section>

    <!-- Two-up cards -->
    <section class="grid gap-4 lg:grid-cols-2">
      <article
        class="rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-[0_1px_2px_rgba(0,0,0,0.04)]"
      >
        <h2 class="text-base font-semibold text-[var(--text-primary)]">
          Awaiting your reply
        </h2>
        <p class="mt-1 text-sm text-[var(--text-secondary)]">
          When support needs more detail on your SecWeb tickets
        </p>
        <p
          v-if="!loading && ticketPreview.length === 0"
          class="mt-5 text-sm text-[var(--text-muted)]"
        >
          No tickets yet — open one to get started.
        </p>
        <ul
          v-else
          class="mt-5 space-y-3"
        >
          <li
            v-for="t in ticketPreview"
            :key="t.id"
            class="flex items-center justify-between rounded-xl bg-[var(--surface-muted)] px-4 py-3"
          >
            <RouterLink
              :to="paths.dashboard.ticketDetail(t.id)"
              class="min-w-0 flex-1 text-left"
            >
              <p class="truncate font-medium text-[var(--text-primary)]">
                {{ t.title }}
              </p>
              <p class="text-xs text-[var(--text-muted)]">
                {{ formatTicketCategoryLabel(t.category) }} · #{{ t.id }}
              </p>
            </RouterLink>
            <TicketStatusBadge :status="t.status" />
          </li>
        </ul>
        <RouterLink
          :to="paths.dashboard.tickets"
          class="mt-4 inline-flex items-center gap-1 text-sm font-semibold text-[var(--brand-green-dark)] hover:underline"
        >
          See all
          <span aria-hidden="true">→</span>
        </RouterLink>
      </article>

      <article
        class="flex flex-col justify-between rounded-2xl border border-dashed border-[var(--border-strong)] bg-[var(--surface-muted)] p-6"
      >
        <div>
          <h2 class="text-base font-semibold text-[var(--text-primary)]">
            Need something else?
          </h2>
          <p class="mt-2 max-w-sm text-sm leading-relaxed text-[var(--text-secondary)]">
            Open another ticket for a different SecWeb product issue.
          </p>
        </div>
        <div class="mt-8 flex justify-end">
          <RouterLink
            :to="paths.dashboard.ticketNew"
            class="flex h-14 w-14 items-center justify-center rounded-full bg-[var(--brand-green)] text-2xl font-light text-white shadow-md transition hover:brightness-95"
            aria-label="Create ticket"
          >
            +
          </RouterLink>
        </div>
      </article>
    </section>

    <!-- Summary strip -->
    <section
      class="grid gap-4 rounded-2xl bg-[var(--surface-muted)] p-5 sm:grid-cols-3"
    >
      <div
        v-for="s in summaries"
        :key="s.label"
        class="text-center sm:text-left"
      >
        <p class="text-xs font-medium uppercase tracking-wide text-[var(--text-muted)]">
          {{ s.label }}
        </p>
        <p class="mt-1 text-2xl font-semibold text-[var(--text-primary)]">
          {{ loading ? '…' : s.value }}
        </p>
      </div>
    </section>

    <!-- Recent activity -->
    <section>
      <div class="mb-4 flex items-center justify-between">
        <h2 class="text-lg font-semibold text-[var(--text-primary)]">
          Recent activity on your requests
        </h2>
        <RouterLink
          :to="paths.dashboard.tickets"
          class="text-sm font-semibold text-[var(--brand-green-dark)] hover:underline"
        >
          See all
        </RouterLink>
      </div>
      <div
        v-if="!loading && recentGroups.length === 0"
        class="rounded-2xl border border-[var(--border-subtle)] bg-white px-4 py-8 text-center text-sm text-[var(--text-secondary)]"
      >
        No recent ticket updates yet.
      </div>
      <div
        v-else
        class="overflow-hidden rounded-2xl border border-[var(--border-subtle)] bg-white"
      >
        <div
          v-for="(group, idx) in recentGroups"
          :key="group.date"
        >
          <div
            v-if="idx > 0"
            class="h-px bg-[var(--border-subtle)]"
          />
          <p class="px-4 py-3 text-xs font-medium text-[var(--text-muted)]">
            {{ group.date }}
          </p>
          <ul>
            <li
              v-for="row in group.items"
              :key="row.id"
            >
              <RouterLink
                :to="paths.dashboard.ticketDetail(String(row.ticketId))"
                class="flex cursor-pointer items-center gap-4 border-t border-[var(--border-subtle)] px-4 py-4 transition hover:bg-[var(--surface-hover)]"
              >
                <div
                  class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-[var(--surface-muted)] text-[var(--text-secondary)]"
                >
                  <IconTicket class="h-5 w-5" />
                </div>
                <div class="min-w-0 flex-1">
                  <p class="font-medium text-[var(--text-primary)]">
                    {{ row.title }}
                  </p>
                  <p class="text-sm text-[var(--text-muted)]">
                    {{ row.subtitle }}
                  </p>
                </div>
                <div class="hidden shrink-0 text-right sm:block">
                  <p class="text-sm font-semibold text-[var(--brand-green-dark)]">
                    {{ row.badge }}
                  </p>
                  <p class="text-xs text-[var(--text-muted)]">
                    {{ row.hint }}
                  </p>
                </div>
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
        </div>
      </div>
    </section>

    <p class="text-center text-xs text-[var(--text-muted)]">
      SecWeb Helpdesk
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import IconTicket from '@/components/icons/IconTicket.vue'
import TicketStatusBadge from '@/components/tickets/TicketStatusBadge.vue'
import { paths } from '@/constants/routes'
import { fetchTicketList, type ApiTicketRow } from '@/api/tickets'
import type { TicketStatus } from '@/types/ticket'
import { formatDateTime } from '@/utils/date-format'
import { formatTicketCategoryLabel, TICKET_STATUS_DEFINITIONS } from '@/utils/ticket-ui'

const loading = ref(true)
const loadError = ref('')
const items = ref<ApiTicketRow[]>([])
const total = ref(0)

onMounted(async () => {
  loading.value = true
  loadError.value = ''
  const res = await fetchTicketList({ page: 1, limit: 100 })
  loading.value = false
  if (!res.ok) {
    loadError.value = res.message
    return
  }
  items.value = res.items
  total.value = res.pagination.total
})

const openTicketCount = computed(() =>
  items.value.filter(t => t.status === 'open' || t.status === 'in_progress').length,
)

const ticketPreview = computed(() => {
  const sorted = [...items.value].sort(
    (a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime(),
  )
  return sorted.slice(0, 3).map(t => ({
    id: String(t.ticket_id),
    title: t.title,
    category: t.category,
    status: t.status as TicketStatus,
  }))
})

const summaries = computed(() => {
  const rows = items.value
  const resolved = rows.filter(t => t.status === 'resolved' || t.status === 'closed').length
  const inProg = rows.filter(t => t.status === 'in_progress').length
  return [
    { label: 'Total (loaded)', value: String(total.value) },
    { label: 'In progress', value: String(inProg) },
    { label: 'Resolved / closed', value: String(resolved) },
  ]
})

function statusTitle(s: TicketStatus): string {
  return TICKET_STATUS_DEFINITIONS.find(d => d.value === s)?.title ?? s
}

function sameCalendarDay(a: Date, b: Date): boolean {
  return (
    a.getFullYear() === b.getFullYear()
    && a.getMonth() === b.getMonth()
    && a.getDate() === b.getDate()
  )
}

function dayGroupLabel(iso: string): string {
  const d = new Date(iso)
  const now = new Date()
  const yest = new Date(now)
  yest.setDate(yest.getDate() - 1)
  if (sameCalendarDay(d, now))
    return 'Today'
  if (sameCalendarDay(d, yest))
    return 'Yesterday'
  try {
    return new Intl.DateTimeFormat(undefined, { dateStyle: 'medium' }).format(d)
  }
  catch {
    return formatDateTime(iso)
  }
}

const recentGroups = computed(() => {
  const sorted = [...items.value].sort(
    (a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime(),
  )
  const top = sorted.slice(0, 12)
  const map = new Map<string, ApiTicketRow[]>()
  for (const t of top) {
    const label = dayGroupLabel(t.updated_at)
    if (!map.has(label))
      map.set(label, [])
    map.get(label)!.push(t)
  }
  return Array.from(map.entries()).map(([date, rows]) => ({
    date,
    items: rows.map(t => ({
      id: `t-${t.ticket_id}`,
      ticketId: t.ticket_id,
      title: t.title,
      subtitle: `${formatTicketCategoryLabel(t.category)} · Updated ${formatDateTime(t.updated_at)}`,
      badge: statusTitle(t.status as TicketStatus),
      hint: `Ticket #${t.ticket_id}`,
    })),
  }))
})
</script>

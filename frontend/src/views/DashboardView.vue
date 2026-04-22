<template>
  <div class="space-y-8">
    <PageHeader
      kicker="Overview"
      title="Your support queue"
      description="See what’s open, scan recent updates, and open a new request when something breaks."
    />

    <div
      v-if="loadError"
      class="rounded-2xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm font-medium text-amber-950 shadow-sm"
      role="alert"
    >
      {{ loadError }}
    </div>

    <!-- Hero metric + actions -->
    <section class="flex flex-col gap-6 lg:flex-row lg:items-stretch lg:justify-between">
      <div
        class="relative flex min-w-0 flex-1 flex-col justify-center overflow-hidden rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-[var(--shadow-card)]"
      >
        <div
          class="pointer-events-none absolute -right-12 -top-12 h-40 w-40 rounded-full bg-[var(--brand-green)]/[0.12] blur-3xl"
          aria-hidden="true"
        />
        <div
          class="pointer-events-none absolute -bottom-8 -left-8 h-32 w-32 rounded-full bg-[var(--surface-mint)] blur-2xl"
          aria-hidden="true"
        />
        <div class="relative">
          <p class="text-xs font-bold uppercase tracking-[0.14em] text-[var(--text-muted)]">
            Open &amp; in progress
          </p>
          <p class="mt-2 flex flex-wrap items-baseline gap-2">
            <span class="text-4xl font-bold tabular-nums tracking-tight text-[var(--text-primary)] lg:text-5xl">
              {{ loading ? '…' : openTicketCount }}
            </span>
            <span class="text-sm font-medium text-[var(--text-secondary)]">tickets needing attention</span>
          </p>
          <p class="mt-3 text-sm text-[var(--text-muted)]">
            Count includes <strong class="font-semibold text-[var(--text-secondary)]">Open</strong> and
            <strong class="font-semibold text-[var(--text-secondary)]">In progress</strong>.
          </p>
        </div>
      </div>
      <div class="flex shrink-0 flex-col justify-center gap-2 sm:flex-row sm:items-center lg:flex-col lg:items-stretch">
        <RouterLink
          :to="paths.dashboard.ticketNew"
          class="inline-flex items-center justify-center rounded-xl bg-[var(--brand-green)] px-5 py-3 text-sm font-bold text-[var(--text-on-green)] shadow-[var(--shadow-card)] transition hover:brightness-[1.02] active:scale-[0.98]"
        >
          + New ticket
        </RouterLink>
        <RouterLink
          :to="paths.dashboard.tickets"
          class="inline-flex items-center justify-center rounded-xl bg-[var(--surface-mint)] px-5 py-3 text-sm font-semibold text-[var(--brand-green-dark)] shadow-sm ring-1 ring-[var(--brand-green)]/20 transition hover:bg-[var(--surface-mint-hover)] active:scale-[0.99]"
        >
          All my tickets
        </RouterLink>
        <RouterLink
          :to="paths.dashboard.tickets"
          class="inline-flex items-center justify-center rounded-xl border-2 border-[var(--border-strong)] bg-white px-5 py-3 text-sm font-semibold text-[var(--text-primary)] shadow-sm transition hover:border-[var(--text-muted)]/40 hover:bg-white active:scale-[0.99]"
        >
          Search
        </RouterLink>
      </div>
    </section>

    <!-- Two-up cards -->
    <section class="grid gap-4 lg:grid-cols-2">
      <article
        class="hd-lift rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-[var(--shadow-card)]"
      >
        <div class="flex items-start justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-[var(--text-primary)]">
              Latest tickets
            </h2>
            <p class="mt-1 text-sm text-[var(--text-secondary)]">
              Most recently updated — jump back in quickly.
            </p>
          </div>
        </div>
        <p
          v-if="!loading && ticketPreview.length === 0"
          class="mt-6 rounded-xl border border-dashed border-[var(--border-strong)] bg-[var(--surface-muted)]/50 px-4 py-6 text-center text-sm font-medium text-[var(--text-secondary)]"
        >
          No tickets yet — create one to get started.
        </p>
        <ul
          v-else
          class="mt-5 space-y-2"
        >
          <li
            v-for="(t, idx) in ticketPreview"
            :key="t.id"
          >
            <RouterLink
              :to="paths.dashboard.ticketDetail(t.id)"
              class="group flex items-center justify-between gap-3 rounded-xl border border-[var(--border-subtle)] bg-[var(--surface-main)] px-4 py-3 shadow-sm transition hover:border-[var(--brand-green)]/45 hover:bg-white hover:shadow-[var(--shadow-card)]"
            >
              <div class="flex min-w-0 flex-1 items-start gap-3">
                <span
                  class="inline-flex min-w-[1.75rem] shrink-0 items-center justify-center rounded-lg bg-gradient-to-b from-white to-neutral-100 px-2 py-0.5 font-mono text-[11px] font-bold tabular-nums text-neutral-800 shadow-sm ring-1 ring-inset ring-neutral-200/90"
                  :aria-label="`Row ${idx + 1} of ${ticketPreview.length}`"
                >{{ idx + 1 }}</span>
                <div class="min-w-0 flex-1">
                  <p class="truncate font-semibold text-[var(--text-primary)]">
                    {{ t.title }}
                  </p>
                  <div class="mt-1.5 flex flex-wrap items-center gap-1.5">
                    <CategoryBadge :category="t.category" />
                  </div>
                </div>
              </div>
              <TicketStatusBadge
                :status="t.status"
                size="md"
              />
            </RouterLink>
          </li>
        </ul>
        <RouterLink
          :to="paths.dashboard.tickets"
          class="mt-4 inline-flex items-center gap-1 text-sm font-bold text-[var(--brand-green-dark)] hover:underline"
        >
          View all tickets
          <span aria-hidden="true">→</span>
        </RouterLink>
      </article>

      <article
        class="flex flex-col justify-between rounded-2xl border-2 border-dashed border-[var(--border-strong)] bg-gradient-to-br from-[var(--surface-mint)]/40 to-white p-6"
      >
        <div>
          <span class="inline-flex rounded-full bg-white/90 px-2.5 py-1 text-[10px] font-bold uppercase tracking-wider text-[var(--brand-green-dark)] ring-1 ring-[var(--brand-green)]/25">
            Quick action
          </span>
          <h2 class="mt-3 text-lg font-semibold text-[var(--text-primary)]">
            Another issue?
          </h2>
          <p class="mt-2 max-w-sm text-sm leading-relaxed text-[var(--text-secondary)]">
            Start a fresh ticket so support can track it separately from your other requests.
          </p>
        </div>
        <div class="mt-8">
          <RouterLink
            :to="paths.dashboard.ticketNew"
            class="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-[var(--brand-green)] py-3.5 text-sm font-bold text-[var(--text-on-green)] shadow-[var(--shadow-card)] transition hover:brightness-[1.02] active:scale-[0.99]"
          >
            <span class="text-lg leading-none">+</span>
            Create ticket
          </RouterLink>
        </div>
      </article>
    </section>

    <!-- Summary strip -->
    <section>
      <h2 class="mb-3 text-xs font-bold uppercase tracking-[0.14em] text-[var(--text-muted)]">
        At a glance
      </h2>
      <div class="grid gap-3 sm:grid-cols-3">
        <div
          v-for="s in summaries"
          :key="s.label"
          class="hd-lift rounded-2xl border border-[var(--border-subtle)] bg-white px-5 py-4 text-center shadow-[var(--shadow-card)] sm:text-left"
        >
          <p class="text-[11px] font-bold uppercase tracking-wide text-[var(--text-muted)]">
            {{ s.label }}
          </p>
          <p class="mt-2 text-3xl font-bold tabular-nums text-[var(--text-primary)]">
            {{ loading ? '…' : s.value }}
          </p>
        </div>
      </div>
    </section>

    <!-- Recent activity -->
    <section>
      <div class="mb-4 flex flex-wrap items-end justify-between gap-3">
        <div>
          <h2 class="text-lg font-semibold text-[var(--text-primary)]">
            Recent activity
          </h2>
          <p class="mt-0.5 text-sm text-[var(--text-secondary)]">
            Sorted by last update — status shown on the right.
          </p>
        </div>
        <RouterLink
          :to="paths.dashboard.tickets"
          class="rounded-full bg-[var(--surface-mint)] px-4 py-2 text-sm font-bold text-[var(--brand-green-dark)] transition hover:bg-[var(--surface-mint-hover)]"
        >
          All tickets
        </RouterLink>
      </div>
      <div
        v-if="!loading && recentGroups.length === 0"
        class="rounded-2xl border border-dashed border-[var(--border-strong)] bg-white px-4 py-10 text-center text-sm font-medium text-[var(--text-secondary)]"
      >
        No recent updates yet — your tickets will appear here.
      </div>
      <div
        v-else
        class="overflow-hidden rounded-2xl border border-[var(--border-subtle)] bg-white shadow-[var(--shadow-card)]"
      >
        <div
          v-for="(group, idx) in recentGroups"
          :key="group.date"
        >
          <div
            v-if="idx > 0"
            class="h-px bg-[var(--border-subtle)]"
          />
          <div class="flex items-center gap-2 bg-[var(--surface-muted)]/80 px-4 py-2.5">
            <span class="text-xs font-bold uppercase tracking-wider text-[var(--text-muted)]">{{ group.date }}</span>
          </div>
          <ul>
            <li
              v-for="row in group.items"
              :key="row.id"
            >
              <RouterLink
                :to="paths.dashboard.ticketDetail(row.id)"
                class="group flex cursor-pointer items-start gap-3 border-t border-[var(--border-subtle)] px-4 py-4 transition hover:bg-gradient-to-r hover:from-[var(--surface-mint)]/30 hover:to-transparent sm:items-center"
              >
                <span
                  class="inline-flex min-w-[1.75rem] shrink-0 items-center justify-center rounded-lg bg-gradient-to-b from-white to-neutral-100 px-2 py-0.5 font-mono text-[11px] font-bold tabular-nums text-neutral-800 shadow-sm ring-1 ring-inset ring-neutral-200/90"
                  :aria-label="`Row ${row.listIndex} of ${recentActivityTotal}`"
                >{{ row.listIndex }}</span>
                <div class="min-w-0 flex-1">
                  <p class="font-semibold text-[var(--text-primary)]">
                    {{ row.title }}
                  </p>
                  <p class="mt-1 text-sm text-[var(--text-secondary)]">
                    {{ row.subtitle }}
                  </p>
                </div>
                <div class="flex shrink-0 flex-col items-end gap-1">
                  <TicketStatusBadge
                    :status="row.status"
                    size="md"
                  />
                </div>
                <svg
                  class="h-5 w-5 shrink-0 text-[var(--text-muted)] transition group-hover:translate-x-0.5 group-hover:text-[var(--brand-green-dark)]"
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

    <p class="text-center text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
      {{ appName }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, inject, onMounted, ref } from 'vue'
import CategoryBadge from '@/components/ui/CategoryBadge.vue'
import PageHeader from '@/components/ui/PageHeader.vue'
import TicketStatusBadge from '@/components/tickets/TicketStatusBadge.vue'
import { paths } from '@/constants/routes'
import { fetchTicketList, type ApiTicketRow } from '@/api/tickets'
import { appBrandKey } from '@/stores/app-detail'
import type { TicketStatus } from '@/types/ticket'
import { formatDateTime } from '@/utils/date-format'
import { formatTicketCategoryLabel } from '@/utils/ticket-ui'

const brand = inject(appBrandKey, null)
const appName = computed(() => brand?.appName.value ?? 'SecWeb HelpDesk')

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

/** Rows shown in Recent activity (same slice used for numbering 1…n). */
const recentActivityTotal = computed(() => {
  const sorted = [...items.value].sort(
    (a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime(),
  )
  return Math.min(12, sorted.length)
})

const ticketPreview = computed(() => {
  const sorted = [...items.value].sort(
    (a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime(),
  )
  return sorted.slice(0, 3).map(t => ({
    id: t.ticket_uuid,
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
    { label: 'Total tickets', value: String(total.value) },
    { label: 'In progress', value: String(inProg) },
    { label: 'Resolved / closed', value: String(resolved) },
  ]
})

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
  const top = sorted.slice(0, 12).map((t, idx) => ({
    id: t.ticket_uuid,
    listIndex: idx + 1,
    title: t.title,
    subtitle: `${formatTicketCategoryLabel(t.category)} · Updated ${formatDateTime(t.updated_at)}`,
    status: t.status as TicketStatus,
    updated_at: t.updated_at,
  }))
  const map = new Map<string, typeof top>()
  for (const row of top) {
    const label = dayGroupLabel(row.updated_at)
    if (!map.has(label))
      map.set(label, [])
    map.get(label)!.push(row)
  }
  return Array.from(map.entries()).map(([date, rows]) => ({
    date,
    items: rows.map((row) => ({
      ...row,
      updated_at: formatDateTime(row.updated_at),
    })),
  }))
})
</script>

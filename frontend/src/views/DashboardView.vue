<template>
  <div class="space-y-7">
    <!-- Hero metric + actions -->
    <section class="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <p class="text-sm font-medium text-[var(--text-secondary)]">
          Your open requests
        </p>
        <p class="mt-1 flex items-baseline gap-2 text-3xl font-semibold tracking-tight text-[var(--text-primary)] lg:text-4xl">
          {{ openTicketCount }}
          <span class="text-base font-normal text-[var(--text-muted)] lg:text-lg">with SecWeb support</span>
        </p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <button
          type="button"
          class="inline-flex items-center justify-center rounded-full bg-[var(--brand-green)] px-5 py-2.5 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95"
        >
          New ticket
        </button>
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
        <ul class="mt-5 space-y-3">
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
                {{ t.category }} · #{{ t.id }}
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
          <button
            type="button"
            class="flex h-14 w-14 items-center justify-center rounded-full bg-[var(--brand-green)] text-2xl font-light text-white shadow-md transition hover:brightness-95"
            aria-label="Create ticket"
          >
            +
          </button>
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
          {{ s.value }}
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
      <div class="overflow-hidden rounded-2xl border border-[var(--border-subtle)] bg-white">
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
import { computed } from 'vue'
import IconTicket from '@/components/icons/IconTicket.vue'
import TicketStatusBadge from '@/components/tickets/TicketStatusBadge.vue'
import { paths } from '@/constants/routes'
import { allMockTicketsWithOverrides, countOpenMockTickets } from '@/data/mocks/tickets.mock'

const openTicketCount = computed(() => countOpenMockTickets())

const ticketPreview = computed(() => allMockTicketsWithOverrides().slice(0, 3))

const summaries = [
  { label: 'Resolved this week', value: '28' },
  { label: 'Avg. first response', value: '2h 10m' },
  { label: 'SLA at risk', value: '3' },
]

const recentGroups = [
  {
    date: 'Today',
    items: [
      {
        id: 'a1',
        title: 'Reset MFA for contractor account',
        subtitle: 'Updated · Security',
        badge: 'In progress',
        hint: 'Ticket #1050',
      },
      {
        id: 'a2',
        title: 'New hire laptop bundle',
        subtitle: 'Created · Hardware',
        badge: 'Open',
        hint: 'Ticket #1049',
      },
    ],
  },
  {
    date: 'Yesterday',
    items: [
      {
        id: 'b1',
        title: 'Expense tool CSV export blank',
        subtitle: 'Comment added · Billing',
        badge: 'Resolved',
        hint: 'Ticket #1041',
      },
    ],
  },
]

</script>

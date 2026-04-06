<template>
  <div v-if="ticket" class="space-y-4">
    <div>
      <RouterLink
        :to="paths.dashboard.tickets"
        class="inline-flex items-center gap-1 text-xs font-semibold text-[var(--brand-green-dark)] hover:underline sm:text-sm"
      >
        ← Back to tickets
      </RouterLink>
    </div>

    <article class="rounded-xl border border-[var(--border-subtle)] bg-white p-4 shadow-sm">
      <div class="min-w-0">
        <p class="text-xs font-medium uppercase tracking-wide text-[var(--text-muted)]">
          Ticket #{{ ticket.id }}
        </p>
        <h1 class="mt-1 text-xl font-semibold tracking-tight text-[var(--text-primary)] sm:text-2xl">
          {{ ticket.title }}
        </h1>
        <div class="mt-2 flex flex-wrap gap-2 text-xs text-[var(--text-secondary)] sm:text-sm">
          <span>{{ ticket.category }}</span>
          <span aria-hidden="true">·</span>
          <span>Reported by {{ ticket.reporterName }}</span>
          <span aria-hidden="true">·</span>
          <span>Opened {{ formatDateTime(ticket.createdAt) }}</span>
        </div>
      </div>

      <div class="mt-4 border-t border-[var(--border-subtle)] pt-4">
        <h2 class="text-xs font-semibold uppercase tracking-wide text-[var(--text-muted)]">
          Description
        </h2>
        <p class="mt-1.5 whitespace-pre-wrap text-sm leading-snug text-[var(--text-secondary)]">
          {{ ticket.description }}
        </p>
      </div>
      <div class="mt-4 flex flex-wrap gap-x-5 gap-y-1 border-t border-[var(--border-subtle)] pt-4 text-xs sm:text-sm">
        <p>
          <span class="text-[var(--text-muted)]">Assignee</span>
          <span class="ml-2 font-medium text-[var(--text-primary)]">{{ ticket.assigneeName ?? 'Unassigned' }}</span>
        </p>
        <p>
          <span class="text-[var(--text-muted)]">Last updated</span>
          <span class="ml-2 font-medium text-[var(--text-primary)]">{{ formatDateTime(ticket.updatedAt) }}</span>
        </p>
      </div>
    </article>

    <TicketStatusEditor
      :status="ticket.status"
      @update:status="onStatusUpdate"
    />

    <section>
      <h2 class="text-base font-semibold text-[var(--text-primary)]">
        Comments
      </h2>
      <p class="mt-0.5 text-xs text-[var(--text-secondary)] sm:text-sm">
        Conversation between you and SecWeb support.
      </p>

      <ul class="mt-4 space-y-2">
        <li
          v-for="c in allComments"
          :key="c.id"
          class="rounded-xl border border-[var(--border-subtle)] bg-white px-3 py-3 shadow-sm"
        >
          <div class="flex flex-wrap items-center justify-between gap-2">
            <div class="flex items-center gap-2">
              <span class="font-medium text-[var(--text-primary)]">{{ c.authorName }}</span>
              <span
                v-if="c.isStaff"
                class="rounded-full bg-[var(--surface-mint)] px-2 py-0.5 text-xs font-semibold text-[var(--brand-green-dark)]"
              >SecWeb support</span>
            </div>
            <time
              class="text-xs text-[var(--text-muted)]"
              :datetime="c.createdAt"
            >{{ formatDateTime(c.createdAt) }}</time>
          </div>
          <p class="mt-1.5 whitespace-pre-wrap text-sm leading-snug text-[var(--text-secondary)]">
            {{ c.body }}
          </p>
        </li>
      </ul>

      <form
        class="mt-4 rounded-xl border border-dashed border-[var(--border-strong)] bg-[var(--surface-muted)]/50 p-4"
        @submit.prevent="addLocalComment"
      >
        <label
          for="new-comment"
          class="block text-sm font-medium text-[var(--text-primary)]"
        >Add a comment</label>
        <textarea
          id="new-comment"
          v-model="draft"
          rows="2"
          placeholder="Write a reply…"
          class="mt-1.5 w-full rounded-xl border border-[var(--border-subtle)] bg-white px-3 py-2 text-sm text-[var(--text-primary)] shadow-sm outline-none ring-[var(--brand-green)] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2"
        />
        <button
          type="submit"
          class="mt-2 rounded-full bg-[var(--brand-green)] px-4 py-2 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:opacity-50"
          :disabled="!draft.trim()"
        >
          Post comment
        </button>
      </form>
    </section>
  </div>

  <div
    v-else
    class="rounded-xl border border-[var(--border-subtle)] bg-white px-4 py-10 text-center"
  >
    <p class="text-base font-semibold text-[var(--text-primary)]">
      Ticket not found
    </p>
    <p class="mt-2 text-sm text-[var(--text-secondary)]">
      We couldn’t find ticket #{{ route.params.id }}.
    </p>
    <RouterLink
      :to="paths.dashboard.tickets"
      class="mt-4 inline-flex rounded-full bg-[var(--brand-green)] px-4 py-2 text-sm font-semibold text-[var(--text-on-green)]"
    >
      Back to tickets
    </RouterLink>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import TicketStatusEditor from '@/components/tickets/TicketStatusEditor.vue'
import { paths } from '@/constants/routes'
import { getMockTicketById } from '@/data/mocks/tickets.mock'
import { setTicketStatusOverride } from '@/stores/ticket-status-overrides'
import { formatDateTime } from '@/utils/date-format'
import type { TicketComment, TicketStatus } from '@/types/ticket'

const route = useRoute()

const ticket = computed(() => {
  const id = typeof route.params.id === 'string' ? route.params.id : ''
  return id ? getMockTicketById(id) : undefined
})

const localComments = ref<TicketComment[]>([])
const draft = ref('')

watch(
  () => route.params.id,
  () => {
    localComments.value = []
    draft.value = ''
  },
)

function onStatusUpdate(status: TicketStatus) {
  const id = ticket.value?.id
  if (id)
    setTicketStatusOverride(id, status)
}

const allComments = computed(() => {
  if (!ticket.value) return []
  return [...ticket.value.comments, ...localComments.value]
})

let localId = 0
function addLocalComment() {
  const text = draft.value.trim()
  if (!text || !ticket.value) return
  localId += 1
  localComments.value.push({
    id: `local-${localId}`,
    authorName: 'Jordan Lee',
    body: text,
    createdAt: new Date().toISOString(),
    isStaff: false,
  })
  draft.value = ''
}
</script>

<template>
  <div
    v-if="loadingTicket"
    class="rounded-xl border border-[var(--border-subtle)] bg-white px-4 py-10 text-center text-sm text-[var(--text-secondary)]"
  >
    Loading ticket…
  </div>

  <div
    v-else-if="ticket"
    class="space-y-4"
  >
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

      <!-- VULN-03: v-html on description — stored XSS with weak server validation; remediate: {{ }} or sanitize (e.g. DOMPurify). -->
      <div class="mt-4 border-t border-[var(--border-subtle)] pt-4">
        <h2 class="text-xs font-semibold uppercase tracking-wide text-[var(--text-muted)]">
          Description
        </h2>
        <div
          class="ticket-html-content mt-1.5 text-sm leading-snug text-[var(--text-secondary)]"
          v-html="ticket.description"
        />
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

    <!-- VULN-04: Assign/unassign/delete use ticket id from the URL only — backend IDOR completes unauthorized access. -->
    <section
      v-if="isApiTicket"
      class="rounded-xl border border-[var(--border-subtle)] bg-white p-4 shadow-sm"
    >
      <h2 class="text-sm font-semibold text-[var(--text-primary)]">
        Assignment
      </h2>
      <p class="mt-1 text-xs text-[var(--text-secondary)]">
        Set the numeric user id of the assignee (staff directory lists ids), or unassign.
      </p>
      <div
        v-if="assignDeleteError"
        class="mt-3 rounded-xl border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-900"
        role="alert"
      >
        {{ assignDeleteError }}
      </div>
      <div class="mt-3 flex flex-wrap items-end gap-2">
        <div class="min-w-[8rem] flex-1">
          <label
            for="assignee-id"
            class="mb-1 block text-xs font-medium text-[var(--text-muted)]"
          >Assignee user id</label>
          <input
            id="assignee-id"
            v-model="assignUserIdInput"
            type="text"
            inputmode="numeric"
            placeholder="e.g. 2"
            class="w-full rounded-xl border border-[var(--border-subtle)] bg-white px-3 py-2 text-sm text-[var(--text-primary)]"
          >
        </div>
        <button
          type="button"
          class="rounded-full bg-[var(--brand-green)] px-4 py-2 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:opacity-50"
          :disabled="assigning || deleting || !assignUserIdParsed"
          @click="onAssign"
        >
          {{ assigning ? 'Assigning…' : 'Assign' }}
        </button>
        <button
          type="button"
          class="rounded-full border border-[var(--border-strong)] bg-white px-4 py-2 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--surface-hover)] disabled:opacity-50"
          :disabled="assigning || deleting"
          @click="onUnassign"
        >
          {{ assigning ? '…' : 'Unassign' }}
        </button>
      </div>
      <div class="mt-6 border-t border-[var(--border-subtle)] pt-4">
        <button
          type="button"
          class="rounded-full border border-red-300 bg-red-50 px-4 py-2 text-sm font-semibold text-red-900 transition hover:bg-red-100 disabled:opacity-50"
          :disabled="assigning || deleting"
          @click="onDeleteTicket"
        >
          {{ deleting ? 'Deleting…' : 'Delete ticket' }}
        </button>
      </div>
    </section>

    <section>
      <h2 class="text-base font-semibold text-[var(--text-primary)]">
        Comments
      </h2>
      <p class="mt-0.5 text-xs text-[var(--text-secondary)] sm:text-sm">
        Conversation between you and SecWeb support.
      </p>

      <div
        v-if="commentsLoadError"
        class="mt-3 rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-950"
        role="status"
      >
        {{ commentsLoadError }}
      </div>

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
          <!-- VULN-03: v-html on comment body — stored XSS with weak server validation. -->
          <div
            class="ticket-html-content mt-1.5 text-sm leading-snug text-[var(--text-secondary)]"
            v-html="c.body"
          />
        </li>
      </ul>

      <form
        class="mt-4 rounded-xl border border-dashed border-[var(--border-strong)] bg-[var(--surface-muted)]/50 p-4"
        @submit.prevent="onPostComment"
      >
        <div
          v-if="commentSubmitError"
          class="mb-3 rounded-xl border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-900"
          role="alert"
        >
          {{ commentSubmitError }}
        </div>
        <label
          for="new-comment"
          class="block text-sm font-medium text-[var(--text-primary)]"
        >Add a comment</label>
        <textarea
          id="new-comment"
          v-model="draft"
          rows="2"
          maxlength="5000"
          placeholder="Write a reply…"
          class="mt-1.5 w-full rounded-xl border border-[var(--border-subtle)] bg-white px-3 py-2 text-sm text-[var(--text-primary)] shadow-sm outline-none ring-[var(--brand-green)] placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2"
        />
        <p class="mt-1 text-xs text-[var(--text-muted)]">
          {{ draft.length }} / 5000
        </p>
        <button
          type="submit"
          class="mt-2 rounded-full bg-[var(--brand-green)] px-4 py-2 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:opacity-50"
          :disabled="!draft.trim() || postingComment"
        >
          {{ postingComment ? 'Posting…' : 'Post comment' }}
        </button>
      </form>
    </section>
  </div>

  <div
    v-else-if="loadError"
    class="rounded-xl border border-[var(--border-subtle)] bg-white px-4 py-10 text-center"
    role="alert"
  >
    <p class="text-base font-semibold text-[var(--text-primary)]">
      Couldn’t load ticket
    </p>
    <p class="mt-2 text-sm text-[var(--text-secondary)]">
      {{ loadError }}
    </p>
    <RouterLink
      :to="paths.dashboard.tickets"
      class="mt-4 inline-flex rounded-full bg-[var(--brand-green)] px-4 py-2 text-sm font-semibold text-[var(--text-on-green)]"
    >
      Back to tickets
    </RouterLink>
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
// VULN-04: Loads ticket/comments by `route.params.id` only — backend IDOR completes unauthorized access.
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  assignTicket,
  createTicketComment,
  deleteTicket,
  fetchTicket,
  fetchTicketComments,
  mapTicketCommentListItemToTicketComment,
  type ApiTicketRow,
} from '@/api/tickets'
import TicketStatusEditor from '@/components/tickets/TicketStatusEditor.vue'
import { paths } from '@/constants/routes'
import { getMockTicketById } from '@/data/mocks/tickets.mock'
import { getAuthUserSnapshot, getSessionCsrfToken } from '@/stores/auth-session'
import { setTicketStatusOverride } from '@/stores/ticket-status-overrides'
import { formatDateTime } from '@/utils/date-format'
import type { Ticket, TicketComment, TicketStatus } from '@/types/ticket'

const route = useRoute()
const router = useRouter()

const fromApi = ref<ApiTicketRow | null>(null)
const loadingTicket = ref(false)
const loadError = ref('')

function apiRowToTicket(row: ApiTicketRow): Ticket {
  const session = getAuthUserSnapshot()
  const you = session !== null && row.reporter_user_id === session.user_id
  return {
    id: String(row.ticket_id),
    title: row.title,
    description: row.description,
    category: row.category,
    status: row.status,
    createdAt: row.created_at,
    updatedAt: row.updated_at,
    assigneeName:
      row.assigned_user_id != null ? `User #${row.assigned_user_id}` : null,
    reporterName: you ? 'You' : `User #${row.reporter_user_id}`,
    comments: [],
  }
}

const ticket = computed((): Ticket | undefined => {
  const id = typeof route.params.id === 'string' ? route.params.id : ''
  if (!id)
    return undefined
  const mock = getMockTicketById(id)
  if (mock)
    return mock
  if (fromApi.value && String(fromApi.value.ticket_id) === id)
    return apiRowToTicket(fromApi.value)
  return undefined
})

const localComments = ref<TicketComment[]>([])
const apiComments = ref<TicketComment[]>([])
const commentsLoadError = ref('')
const commentSubmitError = ref('')
const postingComment = ref(false)
const draft = ref('')

const assignUserIdInput = ref('')
const assignDeleteError = ref('')
const assigning = ref(false)
const deleting = ref(false)

const isApiTicket = computed(() => {
  const raw = route.params.id
  const id = typeof raw === 'string' ? raw : ''
  if (!id || getMockTicketById(id))
    return false
  return fromApi.value !== null
})

const assignUserIdParsed = computed(() => {
  const n = Number.parseInt(assignUserIdInput.value.trim(), 10)
  return Number.isFinite(n) && n > 0 ? n : null
})

async function refreshApiComments(ticketId: number) {
  commentsLoadError.value = ''
  const res = await fetchTicketComments(ticketId)
  if (!res.ok) {
    commentsLoadError.value = res.message
    apiComments.value = []
    return
  }
  apiComments.value = res.items.map(mapTicketCommentListItemToTicketComment)
}

watch(
  () => route.params.id,
  async (raw) => {
    localComments.value = []
    apiComments.value = []
    commentsLoadError.value = ''
    commentSubmitError.value = ''
    draft.value = ''
    fromApi.value = null
    loadError.value = ''
    assignUserIdInput.value = ''
    assignDeleteError.value = ''
    const id = typeof raw === 'string' ? raw : ''
    if (!id)
      return
    if (getMockTicketById(id))
      return
    const n = Number.parseInt(id, 10)
    if (!Number.isFinite(n) || n <= 0)
      return
    loadingTicket.value = true
    const res = await fetchTicket(n)
    loadingTicket.value = false
    if (!res.ok) {
      loadError.value = res.message
      return
    }
    fromApi.value = res.data
    await refreshApiComments(n)
  },
  { immediate: true },
)

function onStatusUpdate(status: TicketStatus) {
  const id = ticket.value?.id
  if (id)
    setTicketStatusOverride(id, status)
}

function currentTicketNumericId(): number | null {
  const raw = route.params.id
  const id = typeof raw === 'string' ? raw : ''
  const n = Number.parseInt(id, 10)
  return Number.isFinite(n) && n > 0 ? n : null
}

async function onAssign() {
  assignDeleteError.value = ''
  const ticketNum = currentTicketNumericId()
  const uid = assignUserIdParsed.value
  if (ticketNum === null || uid === null) {
    assignDeleteError.value = 'Enter a valid assignee user id.'
    return
  }
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    assignDeleteError.value = 'Your session expired. Sign in again.'
    return
  }
  assigning.value = true
  const res = await assignTicket(ticketNum, { assigned_user_id: uid }, csrf)
  assigning.value = false
  if (!res.ok) {
    assignDeleteError.value = res.message
    return
  }
  fromApi.value = res.data
  assignUserIdInput.value = ''
}

async function onUnassign() {
  assignDeleteError.value = ''
  const ticketNum = currentTicketNumericId()
  if (ticketNum === null)
    return
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    assignDeleteError.value = 'Your session expired. Sign in again.'
    return
  }
  assigning.value = true
  const res = await assignTicket(ticketNum, { unassign: true }, csrf)
  assigning.value = false
  if (!res.ok) {
    assignDeleteError.value = res.message
    return
  }
  fromApi.value = res.data
}

async function onDeleteTicket() {
  assignDeleteError.value = ''
  if (!window.confirm('Delete this ticket permanently? This cannot be undone.'))
    return
  const ticketNum = currentTicketNumericId()
  if (ticketNum === null)
    return
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    assignDeleteError.value = 'Your session expired. Sign in again.'
    return
  }
  deleting.value = true
  const res = await deleteTicket(ticketNum, csrf)
  deleting.value = false
  if (!res.ok) {
    assignDeleteError.value = res.message
    return
  }
  await router.push(paths.dashboard.tickets)
}

const allComments = computed(() => {
  if (!ticket.value) return []
  const id = typeof route.params.id === 'string' ? route.params.id : ''
  if (getMockTicketById(id))
    return [...ticket.value.comments, ...localComments.value]
  return apiComments.value
})

let localId = 0
function addLocalMockComment() {
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

async function onPostComment() {
  commentSubmitError.value = ''
  const text = draft.value.trim()
  if (!text || !ticket.value)
    return
  const id = typeof route.params.id === 'string' ? route.params.id : ''
  if (getMockTicketById(id)) {
    addLocalMockComment()
    return
  }
  const ticketNum = Number.parseInt(id, 10)
  if (!Number.isFinite(ticketNum) || ticketNum <= 0)
    return
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    commentSubmitError.value = 'Your session expired. Sign in again.'
    return
  }
  postingComment.value = true
  const res = await createTicketComment(ticketNum, text, csrf)
  postingComment.value = false
  if (!res.ok) {
    commentSubmitError.value = res.message
    return
  }
  draft.value = ''
  await refreshApiComments(ticketNum)
}
</script>

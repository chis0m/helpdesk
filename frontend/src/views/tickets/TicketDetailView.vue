<template>
  <div
    v-if="loadingTicket"
    class="space-y-4 rounded-2xl border border-[var(--border-subtle)] bg-white p-6 shadow-[var(--shadow-card)]"
    aria-busy="true"
    aria-label="Loading ticket"
  >
    <div class="flex flex-wrap gap-2">
      <div class="hd-skeleton h-6 w-16 rounded-md" />
      <div class="hd-skeleton h-6 w-24 rounded-md" />
      <div class="hd-skeleton h-6 w-20 rounded-md" />
    </div>
    <div class="hd-skeleton h-8 w-4/5 max-w-xl rounded-lg" />
    <div class="hd-skeleton h-4 w-full max-w-lg rounded-md" />
    <div class="hd-skeleton mt-4 h-32 w-full rounded-xl" />
    <p class="text-center text-xs font-semibold uppercase tracking-wider text-[var(--text-muted)]">
      Loading ticket…
    </p>
  </div>

  <div
    v-else-if="ticket"
    class="space-y-5"
  >
    <RouterLink
      :to="paths.dashboard.tickets"
      class="inline-flex items-center gap-2 rounded-full border border-[var(--border-subtle)] bg-white/90 px-3 py-1.5 text-xs font-bold text-[var(--brand-green-dark)] shadow-[var(--shadow-sm)] backdrop-blur-sm transition hover:border-[var(--brand-green)]/30 hover:shadow-[var(--shadow-card)] sm:text-sm"
    >
      <span aria-hidden="true">←</span> Back to tickets
    </RouterLink>

    <article class="overflow-hidden rounded-2xl border border-[var(--border-subtle)] bg-white shadow-[var(--shadow-card)] ring-1 ring-black/[0.03]">
      <div class="relative border-b border-[var(--border-subtle)] bg-gradient-to-br from-[var(--surface-mint)]/40 via-white to-[var(--surface-main)]/30 px-5 py-5 sm:px-6">
        <div
          class="pointer-events-none absolute right-0 top-0 h-32 w-32 rounded-full bg-[var(--brand-green)]/10 blur-2xl"
          aria-hidden="true"
        />
        <div class="relative">
          <div class="flex flex-wrap items-center gap-2">
            <span
              class="inline-flex items-center rounded-lg bg-gradient-to-b from-white to-neutral-100 px-2 py-0.5 font-mono text-[11px] font-bold tabular-nums text-neutral-800 shadow-sm ring-1 ring-inset ring-neutral-200/90"
            >ID: {{ ticket.id }}</span>
            <CategoryBadge :category="ticket.category" />
            <TicketStatusBadge
              :status="ticket.status"
              size="md"
            />
          </div>
          <h1 class="mt-4 text-2xl font-bold leading-tight tracking-tight text-[var(--text-primary)] sm:text-3xl">
            {{ ticket.title }}
          </h1>
        </div>
      </div>

      <!-- Ticket facts: one scannable block (no separate cards) -->
      <div class="border-b border-[var(--border-subtle)] bg-[var(--surface-main)]/50 px-5 py-4 sm:px-6">
        <p class="mb-3 text-[10px] font-bold uppercase tracking-[0.16em] text-[var(--text-muted)]">
          Ticket details
        </p>
        <dl
          class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4"
        >
          <div class="min-w-0">
            <dt class="text-xs text-[var(--text-muted)]">
              Ticket ID
            </dt>
            <dd class="mt-0.5 font-mono text-sm font-semibold tabular-nums text-[var(--text-primary)]">
              {{ ticket.id }}
            </dd>
          </div>
          <div class="min-w-0">
            <dt class="text-xs text-[var(--text-muted)]">
              Reported by
            </dt>
            <dd class="mt-0.5">
              <p class="truncate text-sm font-semibold text-[var(--text-primary)]">
                {{ ticket.reporterName }}
              </p>
              <p class="mt-0.5 break-all font-mono text-[13px] font-medium text-[var(--text-secondary)]">
                {{ ticket.reporterEmail || '—' }}
              </p>
            </dd>
          </div>
          <div class="min-w-0 sm:col-span-2 lg:col-span-1">
            <dt class="text-xs text-[var(--text-muted)]">
              Assignee
            </dt>
            <dd class="mt-0.5">
              <template v-if="ticket.assigneeEmail == null">
                <span class="text-sm font-medium text-[var(--text-primary)]">Unassigned</span>
              </template>
              <template v-else>
                <p class="text-sm font-semibold text-[var(--text-primary)]">
                  {{ ticket.assigneeName }}
                </p>
                <p class="mt-0.5 break-all font-mono text-[13px] font-medium text-[var(--text-primary)]">
                  {{ ticket.assigneeEmail }}
                </p>
              </template>
            </dd>
          </div>
          <div class="min-w-0">
            <dt class="text-xs text-[var(--text-muted)]">
              Opened
            </dt>
            <dd class="mt-0.5 text-sm font-medium text-[var(--text-primary)]">
              {{ formatDateTime(ticket.createdAt) }}
            </dd>
          </div>
          <div class="min-w-0">
            <dt class="text-xs text-[var(--text-muted)]">
              Last updated
            </dt>
            <dd class="mt-0.5 text-sm font-medium text-[var(--text-primary)]">
              {{ formatDateTime(ticket.updatedAt) }}
            </dd>
          </div>
        </dl>
      </div>

      <!-- VULN-03: v-html on description — stored XSS with weak server validation; remediate: {{ }} or sanitize (e.g. DOMPurify). -->
      <div class="px-5 pb-6 pt-5 sm:px-6">
        <h2 class="text-base font-semibold text-[var(--text-primary)]">
          Description
        </h2>
        <p class="mt-1 text-xs text-[var(--text-muted)]">
          Full text as submitted with this ticket.
        </p>
        <div
          class="ticket-html-content mt-4 max-w-none rounded-xl border border-[var(--border-subtle)] bg-white px-5 py-4 text-[15px] leading-relaxed text-[var(--text-secondary)] shadow-[var(--shadow-sm)]"
          v-html="ticket.description"
        />
      </div>
    </article>

    <TicketStatusEditor
      :status="ticket.status"
      :saving="statusSaving"
      :error-message="statusError"
      @update:status="onStatusUpdate"
    />

    <section
      v-if="isApiTicket && isAdminOrSuperAdmin"
      class="rounded-2xl border border-[var(--border-subtle)] bg-white p-5 shadow-[var(--shadow-card)] ring-1 ring-black/[0.03]"
    >
      <span class="text-[10px] font-bold uppercase tracking-[0.14em] text-[var(--text-muted)]">Staff</span>
      <h2 class="mt-1 text-lg font-semibold text-[var(--text-primary)]">
        Assignment
      </h2>
      <p class="mt-1 text-sm text-[var(--text-secondary)]">
        Choose a staff or admin user to assign. To unassign, select the ticket's current assignee in the list, then click Unassign.
      </p>
      <div
        v-if="assignDeleteError"
        class="mt-3 rounded-xl border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-900"
        role="alert"
      >
        {{ assignDeleteError }}
      </div>
      <div class="mt-3 space-y-2">
        <label
          for="assignee-select"
          class="mb-1 block text-xs font-medium text-[var(--text-muted)]"
        >Assignee</label>
        <select
          id="assignee-select"
          v-model="selectedAssigneeId"
          class="w-full max-w-md rounded-xl border border-[var(--border-subtle)] bg-white px-3 py-2 text-sm text-[var(--text-primary)] disabled:opacity-60"
          :disabled="assigning || deleting || assignableLoading"
        >
          <option value="">
            {{ assignableLoading ? 'Loading…' : 'Select user…' }}
          </option>
          <option
            v-for="u in assignableOptions"
            :key="u.user_id"
            :value="String(u.user_id)"
          >
            {{ staffOptionLabel(u) }}
          </option>
        </select>
      </div>
      <div class="mt-3 flex flex-wrap items-end gap-2">
        <button
          type="button"
          class="rounded-full bg-[var(--brand-green)] px-4 py-2 text-sm font-semibold text-[var(--text-on-green)] shadow-sm transition hover:brightness-95 disabled:opacity-50"
          :disabled="assigning || deleting || !canSubmitAssign"
          @click="onAssign"
        >
          {{ assigning ? 'Assigning…' : 'Assign' }}
        </button>
        <button
          type="button"
          class="rounded-full border border-[var(--border-strong)] bg-white px-4 py-2 text-sm font-semibold text-[var(--text-primary)] transition hover:bg-[var(--surface-hover)] disabled:opacity-50"
          :disabled="assigning || deleting || !canUnassign"
          @click="onUnassign"
        >
          {{ assigning ? '…' : 'Unassign' }}
        </button>
      </div>
    </section>

    <section
      v-if="isApiTicket"
      class="rounded-2xl border border-[var(--border-subtle)] bg-white p-5 shadow-[var(--shadow-card)] ring-1 ring-black/[0.03]"
    >
      <div
        v-if="deleteTicketError"
        class="mb-3 rounded-xl border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-900"
        role="alert"
      >
        {{ deleteTicketError }}
      </div>
      <div>
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
      <div class="flex flex-wrap items-end justify-between gap-3">
        <div>
          <p class="text-[10px] font-bold uppercase tracking-[0.14em] text-[var(--text-muted)]">
            Discussion
          </p>
          <h2 class="text-lg font-semibold text-[var(--text-primary)]">
            Comments
          </h2>
          <p class="mt-0.5 text-sm text-[var(--text-secondary)]">
            Newest comments appear at the top. Staff replies are labeled.
          </p>
        </div>
        <span
          v-if="allComments.length"
          class="rounded-full bg-[var(--surface-muted)] px-3 py-1 text-xs font-bold text-[var(--text-secondary)] ring-1 ring-[var(--border-subtle)]"
        >{{ allComments.length }} comment{{ allComments.length === 1 ? '' : 's' }}</span>
      </div>

      <div
        v-if="commentsLoadError"
        class="mt-3 rounded-xl border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-950"
        role="status"
      >
        {{ commentsLoadError }}
      </div>

      <ul
        class="mt-4 space-y-4"
        aria-label="Ticket comments"
      >
        <li
          v-for="c in allComments"
          :key="c.id"
          class="overflow-hidden rounded-xl border border-[var(--border-subtle)] bg-white shadow-[var(--shadow-sm)] ring-1 ring-black/[0.03]"
        >
          <div
            class="flex flex-wrap items-center justify-between gap-2 border-b border-[var(--border-subtle)] bg-[var(--surface-main)]/50 px-4 py-3"
          >
            <div class="flex flex-wrap items-center gap-2">
              <span class="font-semibold text-[var(--text-primary)]">{{ c.authorName }}</span>
              <span
                v-if="c.isStaff"
                class="rounded-full bg-[var(--surface-mint)] px-2.5 py-0.5 text-[10px] font-bold uppercase tracking-wide text-[var(--brand-green-dark)] ring-1 ring-[var(--brand-green)]/25"
              >SecWeb support</span>
              <span
                v-else
                class="rounded-full bg-neutral-100 px-2.5 py-0.5 text-[10px] font-bold uppercase tracking-wide text-neutral-700 ring-1 ring-neutral-200/90"
              >Customer</span>
            </div>
            <time
              class="text-xs font-medium tabular-nums text-[var(--text-muted)]"
              :datetime="c.createdAt"
            >{{ formatDateTime(c.createdAt) }}</time>
          </div>
          <!-- VULN-03: v-html on comment body — stored XSS with weak server validation. -->
          <div
            class="ticket-html-content px-4 py-3 text-sm leading-relaxed text-[var(--text-secondary)]"
            v-html="c.body"
          />
        </li>
      </ul>

      <form
        class="mt-4 rounded-2xl border-2 border-dashed border-[var(--brand-green)]/35 bg-[var(--surface-mint)]/20 p-4 sm:p-5"
        @submit.prevent="onPostComment"
      >
        <div
          v-if="commentSubmitError"
          class="mb-3 rounded-xl border border-red-200 bg-red-50 px-3 py-2 text-sm font-medium text-red-900"
          role="alert"
        >
          {{ commentSubmitError }}
        </div>
        <label
          for="new-comment"
          class="block text-sm font-bold text-[var(--text-primary)]"
        >Add a reply</label>
        <textarea
          id="new-comment"
          v-model="draft"
          rows="3"
          maxlength="5000"
          placeholder="Type your message…"
          class="mt-2 w-full rounded-xl border border-[var(--border-subtle)] bg-white px-3 py-2.5 text-sm text-[var(--text-primary)] shadow-sm outline-none ring-[var(--brand-green)] placeholder:font-normal placeholder:text-[var(--text-muted)] focus:border-transparent focus:ring-2"
        />
        <div class="mt-2 flex flex-wrap items-center justify-between gap-2">
          <p class="text-xs font-medium text-[var(--text-muted)]">
            {{ draft.length }} / 5000
          </p>
          <button
            type="submit"
            class="rounded-xl bg-[var(--brand-green)] px-5 py-2.5 text-sm font-bold text-[var(--text-on-green)] shadow-md transition hover:brightness-95 disabled:opacity-50"
            :disabled="!draft.trim() || postingComment"
          >
            {{ postingComment ? 'Sending…' : 'Post reply' }}
          </button>
        </div>
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
// VULN-02: Loads ticket/comments by `route.params.id` only — backend IDOR completes unauthorized access.
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchAdminStaffDirectory, type AdminUserListItem } from '@/api/admin'
import {
  assignTicket,
  createTicketComment,
  deleteTicket,
  fetchTicket,
  fetchTicketComments,
  mapTicketCommentListItemToTicketComment,
  patchTicketStatus,
  type ApiTicketRow,
} from '@/api/tickets'
import CategoryBadge from '@/components/ui/CategoryBadge.vue'
import TicketStatusEditor from '@/components/tickets/TicketStatusEditor.vue'
import TicketStatusBadge from '@/components/tickets/TicketStatusBadge.vue'
import { paths } from '@/constants/routes'
import { getAuthUserSnapshot, getSessionCsrfToken } from '@/stores/auth-session'
import { formatDateTime } from '@/utils/date-format'
import { assigneeDisplayLabel, reporterDisplayLabel } from '@/utils/ticket-ui'
import type { Ticket, TicketComment, TicketStatus } from '@/types/ticket'

const route = useRoute()
const router = useRouter()

const fromApi = ref<ApiTicketRow | null>(null)
const loadingTicket = ref(false)
const loadError = ref('')

function apiRowToTicket(row: ApiTicketRow): Ticket {
  const session = getAuthUserSnapshot()
  const currentId = session?.user_id ?? null
  const reporterEmail = row.reporter_email?.trim() || null
  const assigneeEmail
    = row.assigned_email != null && String(row.assigned_email).trim() !== ''
      ? String(row.assigned_email).trim()
      : null
  return {
    id: String(row.ticket_id),
    title: row.title,
    description: row.description,
    category: row.category,
    status: row.status,
    createdAt: row.created_at,
    updatedAt: row.updated_at,
    assigneeName: assigneeDisplayLabel(row),
    reporterName: reporterDisplayLabel(row, currentId),
    reporterEmail,
    assigneeEmail,
    comments: [],
  }
}

const ticket = computed((): Ticket | undefined => {
  const id = typeof route.params.id === 'string' ? route.params.id : ''
  if (!id)
    return undefined
  if (fromApi.value && String(fromApi.value.ticket_id) === id)
    return apiRowToTicket(fromApi.value)
  return undefined
})

const apiComments = ref<TicketComment[]>([])
const commentsLoadError = ref('')
const commentSubmitError = ref('')
const postingComment = ref(false)
const draft = ref('')

const assignDeleteError = ref('')
const deleteTicketError = ref('')
const assigning = ref(false)
const deleting = ref(false)
const assignableOptions = ref<AdminUserListItem[]>([])
const assignableLoading = ref(false)
const selectedAssigneeId = ref('')

const isApiTicket = computed(() => fromApi.value !== null)

const isAdminOrSuperAdmin = computed(() => {
  void fromApi.value
  const u = getAuthUserSnapshot()
  return u?.role === 'admin' || u?.role === 'super_admin'
})

const ticketHasAssignee = computed(() => {
  const id = fromApi.value?.assigned_user_id
  return id != null && id !== undefined
})

const selectedAssigneeNumeric = computed(() => {
  const n = Number.parseInt(selectedAssigneeId.value, 10)
  return Number.isFinite(n) && n > 0 ? n : null
})

const canSubmitAssign = computed(() => {
  const id = selectedAssigneeNumeric.value
  if (id === null || selectedAssigneeId.value === '')
    return false
  const cur = fromApi.value?.assigned_user_id
  if (cur != null && cur === id)
    return false
  return true
})

const canUnassign = computed(() => {
  if (!ticketHasAssignee.value)
    return false
  const id = selectedAssigneeNumeric.value
  const cur = fromApi.value?.assigned_user_id
  if (id === null || cur == null)
    return false
  return id === cur
})

function staffOptionLabel(u: AdminUserListItem): string {
  const mid = u.middle_name?.trim()
  const name = [u.first_name, mid, u.last_name].filter(x => x && String(x).length > 0).join(' ')
  const role = u.role.replace(/_/g, ' ')
  return name ? `${name} — ${u.email} (${role})` : `${u.email} (${role})`
}

const statusSaving = ref(false)
const statusError = ref('')

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
    apiComments.value = []
    commentsLoadError.value = ''
    commentSubmitError.value = ''
    draft.value = ''
    fromApi.value = null
    loadError.value = ''
    selectedAssigneeId.value = ''
    assignDeleteError.value = ''
    deleteTicketError.value = ''
    assignableOptions.value = []
    statusError.value = ''
    const id = typeof raw === 'string' ? raw : ''
    if (!id)
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

watch(
  () => [isAdminOrSuperAdmin.value, fromApi.value?.ticket_id] as const,
  async ([isAdm, tid]) => {
    assignableOptions.value = []
    selectedAssigneeId.value = ''
    if (!isAdm || tid == null)
      return
    assignableLoading.value = true
    assignDeleteError.value = ''
    const res = await fetchAdminStaffDirectory()
    assignableLoading.value = false
    if (!res.ok) {
      assignDeleteError.value = res.message
      return
    }
    assignableOptions.value = res.items.filter(u => u.role !== 'super_admin')
  },
  { immediate: true },
)

async function onStatusUpdate(status: TicketStatus) {
  const id = ticket.value?.id
  if (!id)
    return
  // VULN-02: Status PATCH uses ticket id from the URL — backend IDOR completes unauthorized changes.
  statusError.value = ''
  const ticketNum = Number.parseInt(id, 10)
  if (!Number.isFinite(ticketNum) || ticketNum <= 0)
    return
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    statusError.value = 'Your session expired. Sign in again.'
    return
  }
  statusSaving.value = true
  const res = await patchTicketStatus(ticketNum, status, csrf)
  statusSaving.value = false
  if (!res.ok) {
    statusError.value = res.message
    return
  }
  fromApi.value = res.data
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
  const uid = selectedAssigneeNumeric.value
  if (ticketNum === null || uid === null || selectedAssigneeId.value === '') {
    assignDeleteError.value = 'Select a user to assign.'
    return
  }
  const cur = fromApi.value?.assigned_user_id
  if (cur != null && cur === uid) {
    assignDeleteError.value = 'This ticket is already assigned to that user.'
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
  selectedAssigneeId.value = ''
}

async function onUnassign() {
  assignDeleteError.value = ''
  const ticketNum = currentTicketNumericId()
  if (ticketNum === null)
    return
  const cur = fromApi.value?.assigned_user_id
  if (cur == null) {
    assignDeleteError.value = 'This ticket has no assignee.'
    return
  }
  const uid = selectedAssigneeNumeric.value
  if (uid === null || uid !== cur) {
    assignDeleteError.value = 'Select the current assignee in the list before unassigning.'
    return
  }
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    assignDeleteError.value = 'Your session expired. Sign in again.'
    return
  }
  assigning.value = true
  const res = await assignTicket(ticketNum, { unassign: true, assigned_user_id: uid }, csrf)
  assigning.value = false
  if (!res.ok) {
    assignDeleteError.value = res.message
    return
  }
  fromApi.value = res.data
}

async function onDeleteTicket() {
  deleteTicketError.value = ''
  if (!window.confirm('Delete this ticket permanently? This cannot be undone.'))
    return
  const ticketNum = currentTicketNumericId()
  if (ticketNum === null)
    return
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    deleteTicketError.value = 'Your session expired. Sign in again.'
    return
  }
  deleting.value = true
  const res = await deleteTicket(ticketNum, csrf)
  deleting.value = false
  if (!res.ok) {
    deleteTicketError.value = res.message
    return
  }
  await router.push(paths.dashboard.tickets)
}

const allComments = computed(() => {
  const list = [...apiComments.value]
  list.sort(
    (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
  )
  return list
})

async function onPostComment() {
  commentSubmitError.value = ''
  const text = draft.value.trim()
  if (!text || !ticket.value)
    return
  const id = typeof route.params.id === 'string' ? route.params.id : ''
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

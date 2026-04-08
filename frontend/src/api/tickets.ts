// VULN-03: Create ticket sends title/description/category as JSON — weak sanitization server-side; UI may render with v-html (see TicketDetailView).
// VULN-04: Ticket id in URL/path for get/update/assign/delete/comments — no client-side authorization (backend IDOR).
// VULN-05: Mutating ticket calls send `X-CSRF-Token`; weak verification is backend CSRF middleware.
// VULN-07: `GET /api/tickets/search?q=` forwards `q` into unsafe SQL on the server (see backend `TicketController.Search`).
import { apiUrl, CSRF_HEADER, readJson } from './client'
import { fetchWithSessionRefresh } from './session-fetch'
import type { ApiErrorEnvelope, ApiSuccessEnvelope } from './types'
import type { TicketComment } from '@/types/ticket'
import { friendlyHttpError } from '@/utils/auth-error-message'
import { logger } from '@/utils/logger'

export type ApiTicketStatus = 'open' | 'in_progress' | 'resolved' | 'closed'

export interface ApiTicketRow {
  ticket_id: number
  reporter_user_id: number
  /** First + last name (or email) from the API when available. */
  reporter_display_name?: string
  reporter_email?: string
  assigned_user_id: number | null
  /** First + last name (or email) when assigned; omitted/null when unassigned. */
  assigned_display_name?: string | null
  /** Support assignee’s email when assigned; null when unassigned. */
  assigned_email?: string | null
  title: string
  description: string
  category: string
  status: ApiTicketStatus
  created_at: string
  updated_at: string
}

export type CreateTicketBody = {
  title: string
  description: string
  category: string
}

function errorMessage(body: unknown): string {
  if (!body || typeof body !== 'object')
    return 'Request failed'
  const env = body as ApiErrorEnvelope
  if (typeof env.message === 'string' && env.message.length > 0)
    return env.message
  if (typeof env.error === 'string' && env.error.length > 0)
    return env.error
  return 'Request failed'
}

function apiErrorMessage(res: Response, json: unknown): string {
  return friendlyHttpError(res.status, errorMessage(json))
}

export async function createTicket(
  body: CreateTicketBody,
  sessionCsrf: string,
): Promise<{ ok: true; data: ApiTicketRow } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/tickets')
  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify({
      title: body.title.trim(),
      description: body.description.trim(),
      category: body.category.trim(),
    }),
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `POST ${url} → ${res.status}`, {
    title: body.title.trim().slice(0, 80),
    category: body.category.trim(),
    'X-CSRF-Token (session)': sessionCsrf,
  })
  if (!res.ok) {
    logger.debug('api:tickets', 'create ticket error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }

  const env = json as ApiSuccessEnvelope<ApiTicketRow>
  if (!env.data || typeof env.data.ticket_id !== 'number') {
    logger.debug('api:tickets', 'create ticket invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:tickets', 'create ticket success envelope (dev)', json)
  return { ok: true, data: env.data }
}

export async function fetchTicket(
  ticketId: number,
): Promise<{ ok: true; data: ApiTicketRow } | { ok: false; status: number; message: string }> {
  const url = apiUrl(`/api/tickets/${ticketId}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:tickets', 'get ticket error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<ApiTicketRow>
  if (!env.data || typeof env.data.ticket_id !== 'number') {
    logger.debug('api:tickets', 'get ticket invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:tickets', 'get ticket success (dev)', env.data)
  return { ok: true, data: env.data }
}

/** List item from `GET /api/tickets/:id/comments`. */
export interface TicketCommentListItem {
  comment_id: number
  ticket_id: number
  author_user_id: number
  author_email: string
  author_first_name: string
  author_last_name: string
  body: string
  created_at: string
  updated_at: string
}

/** UI convention: treat @secweb.ie as internal staff (matches CA seed emails). */
export function mapTicketCommentListItemToTicketComment(row: TicketCommentListItem): TicketComment {
  const parts = [row.author_first_name?.trim(), row.author_last_name?.trim()].filter(Boolean)
  const authorName = parts.length > 0 ? parts.join(' ') : row.author_email
  return {
    id: String(row.comment_id),
    authorName,
    body: row.body,
    createdAt: row.created_at,
    isStaff: /@secweb\.ie$/i.test(row.author_email),
  }
}

export async function fetchTicketComments(
  ticketId: number,
): Promise<
  | { ok: true; items: TicketCommentListItem[] }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl(`/api/tickets/${ticketId}/comments`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:tickets', 'list comments error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<{
    items?: TicketCommentListItem[] | null
  }>
  const raw = env.data?.items
  // Backend may send `items: null` when the Go slice was nil (JSON null is not an array).
  if (raw === undefined || raw === null) {
    logger.debug('api:tickets', 'list comments success (dev)', { count: 0 })
    return { ok: true, items: [] }
  }
  if (!Array.isArray(raw)) {
    logger.debug('api:tickets', 'list comments invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:tickets', 'list comments success (dev)', { count: raw.length })
  return { ok: true, items: raw }
}

export interface CreatedTicketCommentData {
  comment_id: number
  ticket_id: number
  author_user_id: number
  body: string
  created_at: string
  updated_at: string
}

export async function createTicketComment(
  ticketId: number,
  body: string,
  sessionCsrf: string,
): Promise<
  | { ok: true; data: CreatedTicketCommentData }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl(`/api/tickets/${ticketId}/comments`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify({ body: body.trim() }),
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `POST ${url} → ${res.status}`, {
    'X-CSRF-Token (session)': sessionCsrf,
  })
  if (!res.ok) {
    logger.debug('api:tickets', 'create comment error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<CreatedTicketCommentData>
  if (!env.data || typeof env.data.comment_id !== 'number') {
    logger.debug('api:tickets', 'create comment invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:tickets', 'create comment success envelope (dev)', json)
  return { ok: true, data: env.data }
}

export interface TicketListPagination {
  page: number
  limit: number
  total: number
}

/** `GET /api/tickets` — list scoped by role on the server. */
export async function fetchTicketList(opts?: {
  page?: number
  limit?: number
  status?: ApiTicketStatus
  category?: string
}): Promise<
  | { ok: true; items: ApiTicketRow[]; pagination: TicketListPagination }
  | { ok: false; status: number; message: string }
> {
  const params = new URLSearchParams()
  if (opts?.page != null)
    params.set('page', String(opts.page))
  if (opts?.limit != null)
    params.set('limit', String(opts.limit))
  if (opts?.status)
    params.set('status', opts.status)
  const cat = opts?.category?.trim()
  if (cat)
    params.set('category', cat)
  const qs = params.toString()
  const url = apiUrl(`/api/tickets${qs ? `?${qs}` : ''}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:tickets', 'list tickets error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<{
    items: ApiTicketRow[]
    pagination: TicketListPagination
  }>
  const items = env.data?.items
  const pagination = env.data?.pagination
  if (!Array.isArray(items) || !pagination || typeof pagination.total !== 'number') {
    logger.debug('api:tickets', 'list tickets invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, items, pagination }
}

/** VULN-07: `q` is URL-encoded only; server interpolates into raw SQL (unsafe). */
export async function fetchTicketSearch(q: string): Promise<
  | { ok: true; items: ApiTicketRow[]; queryEcho: string }
  | { ok: false; status: number; message: string }
> {
  const trimmed = q.trim()
  if (!trimmed)
    return { ok: false, status: 400, message: 'Search query is required' }
  const params = new URLSearchParams({ q: trimmed })
  const url = apiUrl(`/api/tickets/search?${params}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `GET ${url} → ${res.status}`, { q: trimmed })
  if (!res.ok) {
    logger.debug('api:tickets', 'search tickets error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<{ items: ApiTicketRow[]; query: string }>
  const items = env.data?.items
  const queryEcho = env.data?.query
  if (!Array.isArray(items) || typeof queryEcho !== 'string') {
    logger.debug('api:tickets', 'search tickets invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, items, queryEcho }
}

/** VULN-04: Assignment by ticket id in path — backend IDOR completes unauthorized access. */
export type AssignTicketBody =
  | { assigned_user_id: number }
  | { unassign: true }

export async function assignTicket(
  ticketId: number,
  body: AssignTicketBody,
  sessionCsrf: string,
): Promise<{ ok: true; data: ApiTicketRow } | { ok: false; status: number; message: string }> {
  const url = apiUrl(`/api/tickets/${ticketId}/assign`)
  const payload =
    'unassign' in body && body.unassign
      ? { unassign: true as const }
      : { assigned_user_id: (body as { assigned_user_id: number }).assigned_user_id }
  const res = await fetchWithSessionRefresh(url, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify(payload),
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `PATCH ${url} → ${res.status}`, {
    'X-CSRF-Token (session)': sessionCsrf,
  })
  if (!res.ok) {
    logger.debug('api:tickets', 'assign ticket error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<ApiTicketRow>
  if (!env.data || typeof env.data.ticket_id !== 'number') {
    logger.debug('api:tickets', 'assign ticket invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

/** VULN-04: Delete by ticket id in path — backend IDOR completes unauthorized access. */
export async function deleteTicket(
  ticketId: number,
  sessionCsrf: string,
): Promise<
  | { ok: true; data: { ticket_id: number } }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl(`/api/tickets/${ticketId}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'DELETE',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `DELETE ${url} → ${res.status}`, {
    'X-CSRF-Token (session)': sessionCsrf,
  })
  if (!res.ok) {
    logger.debug('api:tickets', 'delete ticket error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<{ ticket_id: number }>
  if (!env.data || typeof env.data.ticket_id !== 'number') {
    logger.debug('api:tickets', 'delete ticket invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

/** VULN-03: Optional fields may carry HTML; weak sanitization server-side. VULN-04: ticket id in path (IDOR). */
export type PatchTicketBody = Partial<{
  title: string
  description: string
  category: string
}>

export async function patchTicket(
  ticketId: number,
  body: PatchTicketBody,
  sessionCsrf: string,
): Promise<{ ok: true; data: ApiTicketRow } | { ok: false; status: number; message: string }> {
  const payload: Record<string, string> = {}
  if (body.title !== undefined)
    payload.title = body.title.trim()
  if (body.description !== undefined)
    payload.description = body.description.trim()
  if (body.category !== undefined)
    payload.category = body.category.trim()

  const url = apiUrl(`/api/tickets/${ticketId}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify(payload),
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `PATCH ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:tickets', 'patch ticket error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<ApiTicketRow>
  if (!env.data || typeof env.data.ticket_id !== 'number') {
    logger.debug('api:tickets', 'patch ticket invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

/** VULN-04: Status change by ticket id in path (IDOR). */
export async function patchTicketStatus(
  ticketId: number,
  status: ApiTicketStatus,
  sessionCsrf: string,
): Promise<{ ok: true; data: ApiTicketRow } | { ok: false; status: number; message: string }> {
  const url = apiUrl(`/api/tickets/${ticketId}/status`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify({ status }),
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `PATCH ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:tickets', 'patch ticket status error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<ApiTicketRow>
  if (!env.data || typeof env.data.ticket_id !== 'number') {
    logger.debug('api:tickets', 'patch ticket status invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

/** VULN-03: Comment body persisted with weak sanitization. VULN-04: ticket id in path (IDOR). */
export async function patchTicketComment(
  ticketId: number,
  commentId: number,
  body: string,
  sessionCsrf: string,
): Promise<
  | { ok: true; data: CreatedTicketCommentData }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl(`/api/tickets/${ticketId}/comments/${commentId}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify({ body: body.trim() }),
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `PATCH ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:tickets', 'patch comment error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<CreatedTicketCommentData>
  if (!env.data || typeof env.data.comment_id !== 'number') {
    logger.debug('api:tickets', 'patch comment invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

/** VULN-04: Comment delete by ticket id in path (IDOR). */
export async function deleteTicketComment(
  ticketId: number,
  commentId: number,
  sessionCsrf: string,
): Promise<
  | { ok: true; data: { comment_id: number; ticket_id: number } }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl(`/api/tickets/${ticketId}/comments/${commentId}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'DELETE',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
  })
  const json = await readJson(res)
  logger.debug('api:tickets', `DELETE ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:tickets', 'delete comment error envelope (dev)', json)
    return { ok: false, status: res.status, message: apiErrorMessage(res, json) }
  }
  const env = json as ApiSuccessEnvelope<{ comment_id: number; ticket_id: number }>
  if (!env.data || typeof env.data.comment_id !== 'number') {
    logger.debug('api:tickets', 'delete comment invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

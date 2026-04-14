// SEC-02: `PATCH /api/admin/users/:user_id/role` uses numeric id in the path (API contract; see backend).
// SEC-04: Admin mutating routes send `X-CSRF-Token`; compared to session row on the server (SEC-04 remediated).
import { apiUrl, CSRF_HEADER, readJson } from './client'
import { fetchWithSessionRefresh } from './session-fetch'
import type { ApiErrorEnvelope, ApiSuccessEnvelope } from './types'
import type { PortalUser, StaffMember } from '@/types/directory-user'
import { logger } from '@/utils/logger'

export type AdminDirectoryRole = 'user' | 'staff' | 'admin' | 'super_admin'

export interface AdminUserListItem {
  user_id: number
  user_uuid: string
  email: string
  first_name: string
  last_name: string
  middle_name: string | null
  role: AdminDirectoryRole
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface AdminUsersPagination {
  page: number
  limit: number
  total: number
}

export interface AdminUsersListData {
  items: AdminUserListItem[]
  pagination: AdminUsersPagination
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

function buildQuery(params: { page?: number; limit?: number; role?: AdminDirectoryRole }): string {
  const sp = new URLSearchParams()
  if (params.page != null && params.page > 0)
    sp.set('page', String(params.page))
  if (params.limit != null && params.limit > 0)
    sp.set('limit', String(params.limit))
  if (params.role)
    sp.set('role', params.role)
  const q = sp.toString()
  return q ? `?${q}` : ''
}

export async function fetchAdminUsers(params?: {
  page?: number
  limit?: number
  role?: AdminDirectoryRole
}): Promise<
  | { ok: true; data: AdminUsersListData }
  | { ok: false; status: number; message: string }
> {
  const q = buildQuery(params ?? {})
  const url = apiUrl(`/api/admin/users${q}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:admin', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:admin', 'admin users list error envelope (dev)', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<AdminUsersListData>
  if (!env.data?.items || !env.data.pagination) {
    logger.debug('api:admin', 'admin users list invalid shape (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:admin', 'admin users list success (dev)', {
    count: env.data.items.length,
    pagination: env.data.pagination,
  })
  return { ok: true, data: env.data }
}

/**
 * Internal directory: staff, admin, and super_admin rows (API allows only one `role` filter per request).
 */
export async function fetchAdminStaffDirectory(): Promise<
  | { ok: true; items: AdminUserListItem[] }
  | { ok: false; status: number; message: string }
> {
  const [staffRes, adminRes, superRes] = await Promise.all([
    fetchAdminUsers({ role: 'staff', limit: 100, page: 1 }),
    fetchAdminUsers({ role: 'admin', limit: 100, page: 1 }),
    fetchAdminUsers({ role: 'super_admin', limit: 100, page: 1 }),
  ])

  if (!staffRes.ok)
    return staffRes
  if (!adminRes.ok)
    return adminRes
  if (!superRes.ok)
    return superRes

  const merged = [
    ...staffRes.data.items,
    ...adminRes.data.items,
    ...superRes.data.items,
  ]
  merged.sort((a, b) => a.email.localeCompare(b.email, undefined, { sensitivity: 'base' }))
  return { ok: true, items: merged }
}

function buildDisplayName(row: AdminUserListItem): string {
  const mid = row.middle_name?.trim()
  return [row.first_name, mid, row.last_name].filter(x => x && x.length > 0).join(' ')
}

export function adminUserToPortalUser(row: AdminUserListItem): PortalUser {
  return {
    id: String(row.user_id),
    email: row.email,
    displayName: buildDisplayName(row) || row.email,
    createdAt: row.created_at,
  }
}

export function adminUserToStaffMember(row: AdminUserListItem): StaffMember {
  return {
    id: String(row.user_id),
    email: row.email,
    displayName: buildDisplayName(row) || row.email,
    createdAt: row.created_at,
    isAdmin: row.role === 'admin' || row.role === 'super_admin',
    role: row.role,
  }
}

export type CreateStaffBody = {
  email: string
  password: string
  first_name: string
  last_name: string
  middle_name?: string
  is_active?: boolean
  /** Omit or `staff`. `admin` only when the caller is admin or super_admin (API enforces). */
  role?: 'staff' | 'admin'
}

export interface CreateStaffResponseData {
  user_id: number
  user_uuid: string
  email: string
  role: string
  is_active: boolean
}

export async function createStaffUser(
  body: CreateStaffBody,
  sessionCsrf: string,
): Promise<
  | { ok: true; data: CreateStaffResponseData }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl('/api/admin/staff')
  const payload: Record<string, unknown> = {
    email: body.email.trim(),
    password: body.password,
    first_name: body.first_name.trim(),
    last_name: body.last_name.trim(),
  }
  const mid = body.middle_name?.trim()
  if (mid)
    payload.middle_name = mid
  if (typeof body.is_active === 'boolean')
    payload.is_active = body.is_active
  /** API: optional `staff` | `admin` (default staff). Admin or super_admin callers may use `admin`. */
  payload.role = body.role ?? 'staff'

  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify(payload),
  })
  const json = await readJson(res)
  logger.debug('api:admin', `POST ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:admin', 'create staff error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<CreateStaffResponseData>
  if (!env.data || typeof env.data.user_id !== 'number') {
    logger.debug('api:admin', 'create staff invalid shape', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:admin', 'create staff success envelope (dev)', env.data)
  return { ok: true, data: env.data }
}

export type CreateStaffInviteBody = {
  email: string
  first_name: string
  last_name: string
  middle_name?: string
  /** Omit or `staff`. `admin` only when the caller is admin or super_admin (API enforces). */
  role?: 'staff' | 'admin'
}

export interface CreateStaffInviteResponseData {
  invite_id: number
  email: string
  expires_at_utc: string
  target_role: string
  delivery?: string
  notice?: string
  /** Copy for the invitee when the API returns a URL (e.g. dev / log delivery). */
  invite_url?: string
}

export async function createStaffInvite(
  body: CreateStaffInviteBody,
  sessionCsrf: string,
): Promise<
  | { ok: true; data: CreateStaffInviteResponseData }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl('/api/admin/invites/staff')
  const payload: Record<string, string> = {
    email: body.email.trim(),
    first_name: body.first_name.trim(),
    last_name: body.last_name.trim(),
  }
  const mid = body.middle_name?.trim()
  if (mid)
    payload.middle_name = mid
  payload.role = body.role ?? 'staff'

  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify(payload),
  })
  const json = await readJson(res)
  logger.debug('api:admin', `POST ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:admin', 'create staff invite error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<CreateStaffInviteResponseData>
  if (!env.data || typeof env.data.invite_id !== 'number') {
    logger.debug('api:admin', 'create staff invite invalid shape', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

export async function patchAdminUserRole(
  userId: number,
  role: AdminDirectoryRole,
  sessionCsrf: string,
): Promise<
  | { ok: true; data: { user_id: number; user_uuid: string; role: AdminDirectoryRole } }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl(`/api/admin/users/${userId}/role`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify({ role }),
  })
  const json = await readJson(res)
  logger.debug('api:admin', `PATCH ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:admin', 'patch user role error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<{ user_id: number; user_uuid: string; role: AdminDirectoryRole }>
  if (!env.data || typeof env.data.user_id !== 'number') {
    logger.debug('api:admin', 'patch user role invalid shape', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

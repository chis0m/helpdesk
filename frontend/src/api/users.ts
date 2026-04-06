// VULN-02: GET/PATCH `/api/users/:id` by numeric id from the UI — no client-side actor-vs-target check (backend IDOR).
// VULN-05: PATCH sends `X-CSRF-Token`; weak verification is backend CSRF middleware.
import { apiUrl, CSRF_HEADER, readJson } from './client'
import type { ApiErrorEnvelope, ApiSuccessEnvelope } from './types'
import { logger } from '@/utils/logger'

export interface UserProfileData {
  user_id: number
  user_uuid: string
  email: string
  first_name: string
  last_name: string
  middle_name: string | null
  role: string
  is_active: boolean
}

export type PatchUserBody = {
  email?: string
  first_name?: string
  last_name?: string
  middle_name?: string | null
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

export async function fetchUser(
  userId: number,
): Promise<{ ok: true; data: UserProfileData } | { ok: false; status: number; message: string }> {
  const url = apiUrl(`/api/users/${userId}`)
  const res = await fetch(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:users', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:users', 'GET user error envelope (dev)', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<UserProfileData>
  if (!env.data) {
    logger.debug('api:users', 'GET user missing data (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:users', 'GET user success data (dev)', env.data)
  return { ok: true, data: env.data }
}

export async function patchUser(
  userId: number,
  body: PatchUserBody,
  csrfToken: string,
): Promise<{ ok: true; data: UserProfileData } | { ok: false; status: number; message: string }> {
  const url = apiUrl(`/api/users/${userId}`)
  const res = await fetch(url, {
    method: 'PATCH',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: csrfToken,
    },
    body: JSON.stringify(body),
  })
  const json = await readJson(res)
  logger.debug('api:users', `PATCH ${url} → ${res.status}`, {
    patchBody: body,
    'X-CSRF-Token (session)': csrfToken,
  })
  if (!res.ok) {
    logger.debug('api:users', 'PATCH user error envelope (dev)', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<UserProfileData>
  if (!env.data) {
    logger.debug('api:users', 'PATCH user missing data (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:users', 'PATCH user success envelope (dev)', json)
  return { ok: true, data: env.data }
}

/**
 * VULN-01: `refresh_token` cookie; VULN-05: session CSRF on POST.
 * Uses raw `fetch` so `session-fetch` never wraps this (avoids 401→refresh loops).
 */
import { apiUrl, CSRF_HEADER, readJson } from './client'
import { loggedFetch } from './http-dev-log'
import type { ApiErrorEnvelope, ApiSuccessEnvelope } from './types'
import { logger } from '@/utils/logger'

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

export interface RefreshResponseData {
  user_id: number
  user_uuid: string
  access_expires_at_utc: string
  csrf_token: string
  csrf_expires_at_utc: string
}

export async function postAuthRefresh(
  sessionCsrf: string,
): Promise<
  | { ok: true; data: RefreshResponseData }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl('/api/auth/refresh')
  const res = await loggedFetch('api:auth', url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
  })
  const json = await readJson(res)
  if (!res.ok)
    return { ok: false, status: res.status, message: errorMessage(json) }
  const env = json as ApiSuccessEnvelope<RefreshResponseData>
  if (!env.data?.csrf_token || typeof env.data.user_id !== 'number') {
    logger.debug('api:auth', 'refresh invalid shape', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

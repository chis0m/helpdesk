// VULN-04: `POST /api/invites/accept` uses public CSRF (same weak middleware as signup/login).
import { apiUrl, CSRF_HEADER, readJson } from './client'
import { fetchWithSessionRefresh } from './session-fetch'
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

export type InviteVerifyInvalid = { valid: false }
export type InviteVerifyValid = {
  valid: true
  email: string
  first_name: string
  last_name: string
}

export async function fetchInviteVerify(token: string): Promise<
  | { ok: true; data: InviteVerifyInvalid | InviteVerifyValid }
  | { ok: false; status: number; message: string }
> {
  const params = new URLSearchParams({ token: token.trim() })
  const url = apiUrl(`/api/invites/verify?${params}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:invites', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:invites', 'verify invite error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<InviteVerifyInvalid | InviteVerifyValid>
  const data = env.data
  if (!data || typeof data !== 'object' || typeof (data as { valid?: unknown }).valid !== 'boolean') {
    logger.debug('api:invites', 'verify invite invalid shape', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data }
}

export interface AcceptInviteResponseData {
  user_id: number
  user_uuid: string
  email: string
  role: string
  redirect_to: string
}

export async function acceptInviteRequest(
  token: string,
  password: string,
  publicCsrf: string,
): Promise<
  | { ok: true; data: AcceptInviteResponseData }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl('/api/invites/accept')
  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: publicCsrf,
    },
    body: JSON.stringify({
      token: token.trim(),
      password,
    }),
  })
  const json = await readJson(res)
  logger.debug('api:invites', `POST ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:invites', 'accept invite error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<AcceptInviteResponseData>
  if (!env.data?.user_uuid || typeof env.data.redirect_to !== 'string') {
    logger.debug('api:invites', 'accept invite invalid shape', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

// VULN-01: Login/refresh/logout/change-password rely on cookie session (weak flags set server-side).
// VULN-04: Public vs session CSRF tokens sent here; weak verification is backend CSRF middleware.
import { postAuthRefresh, type RefreshResponseData } from './auth-refresh-internal'
import { apiUrl, CSRF_HEADER, readJson } from './client'
import { fetchWithSessionRefresh } from './session-fetch'
import type { ApiErrorEnvelope, ApiSuccessEnvelope } from './types'
import { setSessionCsrfPair } from '@/stores/auth-session'
import { logger } from '@/utils/logger'

export interface PublicCsrfData {
  csrf_token: string
  csrf_expires_at_utc: string
}

export interface LoginResponseData {
  user_id: number
  user_uuid: string
  email: string
  role: string
  must_change_password: boolean
  access_expires_at_utc: string
  csrf_token: string
  csrf_expires_at_utc: string
}

export interface SignupResponseData {
  user_uuid: string
  email: string
  redirect_to: string
}

export type SignupRequestBody = {
  email: string
  password: string
  first_name: string
  last_name: string
  middle_name?: string
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

export async function fetchPublicCsrfToken(): Promise<{ ok: true; token: string } | { ok: false; message: string }> {
  const url = apiUrl('/api/auth/public-csrf-token')
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:auth', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'public CSRF error body', json)
    return { ok: false, message: errorMessage(json) }
  }

  const env = json as ApiSuccessEnvelope<PublicCsrfData>
  logger.debug('api:auth', 'public CSRF response (dev — includes token)', env.data ?? json)
  const token = env.data?.csrf_token
  if (!token)
    return { ok: false, message: 'Missing CSRF token' }
  return { ok: true, token }
}

export async function loginRequest(
  email: string,
  password: string,
  publicCsrf: string,
): Promise<{ ok: true; data: LoginResponseData } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/auth/login')
  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: publicCsrf,
    },
    body: JSON.stringify({ email, password }),
  })
  const json = await readJson(res)
  logger.debug('api:auth', `POST ${url} → ${res.status}`, {
    email,
    'X-CSRF-Token (public)': publicCsrf,
  })
  if (!res.ok) {
    logger.debug('api:auth', 'login error envelope (dev)', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }

  const env = json as ApiSuccessEnvelope<LoginResponseData>
  if (!env.data) {
    logger.debug('api:auth', 'login missing data envelope (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug(
    'api:auth',
    'login success envelope (dev — includes session CSRF in data)',
    json,
  )
  return { ok: true, data: env.data }
}

export async function signupRequest(
  body: SignupRequestBody,
  publicCsrf: string,
): Promise<{ ok: true; data: SignupResponseData } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/auth/signup')
  const payload: Record<string, string> = {
    email: body.email.trim(),
    password: body.password,
    first_name: body.first_name.trim(),
    last_name: body.last_name.trim(),
  }
  const mid = body.middle_name?.trim()
  if (mid)
    payload.middle_name = mid

  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: publicCsrf,
    },
    body: JSON.stringify(payload),
  })
  const json = await readJson(res)
  logger.debug('api:auth', `POST ${url} → ${res.status}`, {
    email: payload.email,
    first_name: payload.first_name,
    last_name: payload.last_name,
    middle_name: payload.middle_name,
    'X-CSRF-Token (public)': publicCsrf,
  })
  if (!res.ok) {
    logger.debug('api:auth', 'signup error envelope (dev)', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }

  const env = json as ApiSuccessEnvelope<SignupResponseData>
  if (!env.data) {
    logger.debug('api:auth', 'signup missing data envelope (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:auth', 'signup success envelope (dev)', json)
  return { ok: true, data: env.data }
}

export interface LogoutResponseData {
  redirect_to: string
}

export async function logoutRequest(
  sessionCsrf: string,
): Promise<{ ok: true; data: LogoutResponseData } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/auth/logout')
  const res = await fetchWithSessionRefresh(url, {
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

  const env = json as ApiSuccessEnvelope<LogoutResponseData>
  if (!env.data)
    return { ok: false, status: res.status, message: 'Invalid response' }
  return { ok: true, data: env.data }
}

/** VULN-01: Session cookies via `credentials: 'include'`. VULN-04: `X-CSRF-Token` on POST. */
export type ChangePasswordBody = {
  current_password: string
  new_password: string
}

export async function changePasswordRequest(
  body: ChangePasswordBody,
  sessionCsrf: string,
): Promise<{ ok: true } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/auth/change-password')
  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
    body: JSON.stringify({
      current_password: body.current_password,
      new_password: body.new_password,
    }),
  })
  const json = await readJson(res)
  logger.debug('api:auth', `POST ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'change-password error envelope (dev)', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  return { ok: true }
}

/** VULN-04: Public CSRF + POST; weak verification is backend middleware. */
export async function forgotPasswordRequest(
  email: string,
  publicCsrf: string,
): Promise<{ ok: true } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/auth/forgot-password')
  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: publicCsrf,
    },
    body: JSON.stringify({ email: email.trim() }),
  })
  const json = await readJson(res)
  logger.debug('api:auth', `POST ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'forgot-password error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  return { ok: true }
}

/** VULN-04: Public CSRF + POST. */
export async function resetPasswordRequest(
  token: string,
  newPassword: string,
  publicCsrf: string,
): Promise<
  | { ok: true; data: { redirect_to: string } }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl('/api/auth/reset-password')
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
      new_password: newPassword,
    }),
  })
  const json = await readJson(res)
  logger.debug('api:auth', `POST ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'reset-password error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<{ redirect_to: string }>
  if (!env.data?.redirect_to)
    return { ok: false, status: res.status, message: 'Invalid response' }
  return { ok: true, data: env.data }
}

/** VULN-01: `refresh_token` cookie; VULN-04: session CSRF on POST. Uses raw `fetch` via `postAuthRefresh`. */
export type { RefreshResponseData }

export async function refreshRequest(
  sessionCsrf: string,
): Promise<
  | { ok: true; data: RefreshResponseData }
  | { ok: false; status: number; message: string }
> {
  return postAuthRefresh(sessionCsrf)
}

/** VULN-01: Session cookie required. */
export async function fetchSessionCsrfToken(): Promise<
  | { ok: true; data: PublicCsrfData }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl('/api/auth/csrf-token')
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:auth', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'session CSRF error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<PublicCsrfData>
  const token = env.data?.csrf_token
  const exp = env.data?.csrf_expires_at_utc
  if (!token || !env.data || typeof exp !== 'string')
    return { ok: false, status: res.status, message: 'Invalid response' }
  setSessionCsrfPair(token, exp)
  return { ok: true, data: env.data }
}

/** VULN-01: Session cookie; response includes numeric `user_id` (VULN-02 baseline). */
export interface AuthMeData {
  user_id: number
  user_uuid: string
  email: string
  first_name: string
  last_name: string
  middle_name: string | null
  role: string
  is_active: boolean
  must_change_password: boolean
}

export async function fetchMe(): Promise<
  | { ok: true; data: AuthMeData }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl('/api/auth/me')
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:auth', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'me error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<AuthMeData>
  if (!env.data || typeof env.data.user_id !== 'number') {
    logger.debug('api:auth', 'me invalid shape', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

/** VULN-01: Requires authenticated session cookie. */
export interface AuthSessionRow {
  session_id: string
  created_at: string
  user_agent: string | null
  ip: string | null
  is_current: boolean
}

export async function fetchAuthSessions(): Promise<
  | { ok: true; items: AuthSessionRow[] }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl('/api/auth/sessions')
  const res = await fetchWithSessionRefresh(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:auth', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'sessions list error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<{ items: AuthSessionRow[] }>
  const items = env.data?.items
  if (!Array.isArray(items)) {
    logger.debug('api:auth', 'sessions invalid shape', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, items }
}

/** VULN-04: Session CSRF on POST. */
export async function revokeMyOtherSessionsRequest(
  sessionCsrf: string,
): Promise<{ ok: true } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/auth/sessions/revoke-my-other-sessions')
  const res = await fetchWithSessionRefresh(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
  })
  const json = await readJson(res)
  logger.debug('api:auth', `POST ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'revoke other sessions error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  return { ok: true }
}

/** VULN-04: Session CSRF on DELETE. */
export async function revokeAuthSession(
  sessionId: string,
  sessionCsrf: string,
): Promise<
  | { ok: true; data: { revoked_session_id: string; logged_out: boolean } }
  | { ok: false; status: number; message: string }
> {
  const id = sessionId.trim()
  const url = apiUrl(`/api/auth/sessions/${encodeURIComponent(id)}`)
  const res = await fetchWithSessionRefresh(url, {
    method: 'DELETE',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
  })
  const json = await readJson(res)
  logger.debug('api:auth', `DELETE ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'revoke session error', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }
  const env = json as ApiSuccessEnvelope<{ revoked_session_id: string; logged_out: boolean }>
  if (!env.data || typeof env.data.logged_out !== 'boolean') {
    logger.debug('api:auth', 'revoke session invalid shape', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  return { ok: true, data: env.data }
}

/**
 * Session access refresh after 401 (TASK.md §8.1) + session CSRF rotation before unsafe requests.
 * Server CSRF TTL is independent of access token TTL; without rotation, mutating calls can 403 while still logged in.
 */
import { apiUrl, CSRF_HEADER, readJson } from './client'
import { loggedFetch } from './http-dev-log'
import type { ApiSuccessEnvelope } from './types'
import {
  invalidateClientSessionAndRedirect,
  refreshSessionOnce,
  shouldAttempt401Refresh,
} from './session-refresh'
import {
  getCsrfExpiresAtUtc,
  getSessionCsrfToken,
  setSessionCsrfPair,
} from '@/stores/auth-session'
import { logger } from '@/utils/logger'

/** Refresh session CSRF this long before server expiry (matches access refresh buffer style). */
const CSRF_REFRESH_SKEW_MS = 120_000

function isSafeMethod(method: string): boolean {
  const m = method.toUpperCase()
  return m === 'GET' || m === 'HEAD' || m === 'OPTIONS'
}

/**
 * POST routes validated with `PublicAuthCSRFRequired` — header must be the **public** token from
 * `GET /api/auth/public-csrf-token`, not session CSRF. `fetchWithSessionRefresh` must not overwrite it.
 */
const PUBLIC_CSRF_POST_PATHS = new Set([
  '/api/auth/login',
  '/api/auth/signup',
  '/api/auth/forgot-password',
  '/api/auth/reset-password',
  '/api/invites/accept',
])

function isPublicCsrfPost(urlStr: string, method: string): boolean {
  if (isSafeMethod(method))
    return false
  try {
    return PUBLIC_CSRF_POST_PATHS.has(new URL(urlStr).pathname)
  }
  catch {
    return false
  }
}

function isCsrfExpiredError(json: unknown): boolean {
  if (!json || typeof json !== 'object')
    return false
  const o = json as Record<string, unknown>
  const msg = typeof o.message === 'string' ? o.message.toLowerCase() : ''
  const err = typeof o.error === 'string' ? o.error.toLowerCase() : ''
  return msg.includes('csrf token expired') || err.includes('csrf token expired')
}

/** GET /api/auth/csrf-token — no CSRF header required (safe method). Updates stored token + expiry. */
export async function issueSessionCsrfFromServer(): Promise<boolean> {
  const url = apiUrl('/api/auth/csrf-token')
  const getInit = (): RequestInit => ({
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  let res = await loggedFetch('api:fetch', url, getInit())
  if (res.status === 401 && shouldAttempt401Refresh(url)) {
    const ok = await refreshSessionOnce()
    if (!ok) {
      logger.debug('session-fetch', 'issue CSRF: refresh failed after 401')
      return false
    }
    res = await loggedFetch('api:fetch', url, getInit())
  }
  const json = await readJson(res)
  if (!res.ok)
    return false
  const env = json as ApiSuccessEnvelope<{ csrf_token: string; csrf_expires_at_utc: string }>
  const token = env.data?.csrf_token
  const exp = env.data?.csrf_expires_at_utc
  if (!token || typeof exp !== 'string')
    return false
  setSessionCsrfPair(token, exp)
  return true
}

async function ensureSessionCsrfFresh(): Promise<boolean> {
  const expRaw = getCsrfExpiresAtUtc()
  if (expRaw) {
    const expMs = Date.parse(expRaw)
    if (!Number.isNaN(expMs) && expMs - Date.now() > CSRF_REFRESH_SKEW_MS)
      return true
  }
  return issueSessionCsrfFromServer()
}

export async function fetchWithSessionRefresh(
  url: string,
  init?: RequestInit,
): Promise<Response> {
  const method = (init?.method ?? 'GET').toUpperCase()
  let preparedInit: RequestInit = init ?? {}
  const publicCsrfPost = isPublicCsrfPost(url, method)

  if (!isSafeMethod(method) && !publicCsrfPost) {
    await ensureSessionCsrfFresh()
    const headers = new Headers(preparedInit.headers ?? undefined)
    const csrf = getSessionCsrfToken()
    if (csrf)
      headers.set(CSRF_HEADER, csrf)
    preparedInit = { ...preparedInit, headers }
  }

  let res = await loggedFetch('api:fetch', url, preparedInit)

  if (res.status === 403 && !isSafeMethod(method) && !publicCsrfPost) {
    const json = await readJson(res.clone())
    if (isCsrfExpiredError(json)) {
      const ok = await issueSessionCsrfFromServer()
      if (ok) {
        const headers = new Headers(preparedInit.headers ?? undefined)
        const csrf = getSessionCsrfToken()
        if (csrf) {
          headers.set(CSRF_HEADER, csrf)
          res = await loggedFetch('api:fetch', url, { ...preparedInit, headers })
        }
      }
    }
  }

  if (res.status !== 401 || !shouldAttempt401Refresh(url))
    return res
  if (publicCsrfPost)
    return res
  const ok = await refreshSessionOnce()
  if (!ok)
    return res
  const headers = new Headers(preparedInit.headers ?? undefined)
  const csrf = getSessionCsrfToken()
  if (csrf)
    headers.set(CSRF_HEADER, csrf)
  const res2 = await loggedFetch('api:fetch', url, { ...preparedInit, headers })
  if (res2.status === 401)
    invalidateClientSessionAndRedirect()
  return res2
}

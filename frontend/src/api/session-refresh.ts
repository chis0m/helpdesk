/**
 * Proactive `POST /api/auth/refresh` before access expiry + mutex for concurrent refresh.
 */
import { postAuthRefresh } from './auth-refresh-internal'
import {
  clearAuthSession,
  getAccessExpiresAtUtc,
  getSessionCsrfToken,
  setAuthSessionFromRefresh,
} from '@/stores/auth-session'
import { logger } from '@/utils/logger'

let refreshTimer: ReturnType<typeof setTimeout> | null = null
let refreshInFlight: Promise<boolean> | null = null

let onRefreshFailure: (() => void) | null = null

export function registerSessionRefreshFailure(handler: () => void): void {
  onRefreshFailure = handler
}

/** Clear client session, stop timers, and navigate to login (e.g. after 401 or failed refresh). */
export function invalidateClientSessionAndRedirect(): void {
  clearSessionRefreshSchedule()
  clearAuthSession()
  onRefreshFailure?.()
}

/** Paths where 401 is not “access expired” (wrong password, invalid token, etc.). */
const NO_REFRESH_401_PATHS = new Set([
  '/api/auth/login',
  '/api/auth/signup',
  '/api/auth/public-csrf-token',
  '/api/auth/refresh',
  '/api/auth/forgot-password',
  '/api/auth/reset-password',
  '/api/invites/verify',
  '/api/invites/accept',
])

export function shouldAttempt401Refresh(url: string): boolean {
  let pathname: string
  try {
    pathname = new URL(url).pathname
  }
  catch {
    return false
  }
  return !NO_REFRESH_401_PATHS.has(pathname)
}

export function clearSessionRefreshSchedule(): void {
  if (refreshTimer !== null) {
    clearTimeout(refreshTimer)
    refreshTimer = null
  }
}

function computeDelayMs(expiresAtUtc: string): number {
  const expMs = Date.parse(expiresAtUtc)
  if (Number.isNaN(expMs))
    return -1
  const now = Date.now()
  const ttl = expMs - now
  if (ttl <= 0)
    return 0
  const buffer = Math.min(60_000, Math.max(10_000, Math.floor(ttl * 0.05)))
  return Math.max(0, ttl - buffer)
}

/**
 * Schedule one-shot refresh before `access_expires_at_utc`. Call after login, after successful refresh,
 * and on app init when session storage still holds expiry.
 */
export function scheduleAccessRefresh(): void {
  clearSessionRefreshSchedule()
  const expiresRaw = getAccessExpiresAtUtc()
  const csrf = getSessionCsrfToken()
  if (!expiresRaw || !csrf)
    return
  const delay = computeDelayMs(expiresRaw)
  if (delay === -1) {
    logger.debug('session-refresh', 'invalid access_expires_at_utc, skip schedule')
    return
  }
  if (delay === 0) {
    void refreshSessionOnce()
    return
  }
  refreshTimer = setTimeout(() => {
    refreshTimer = null
    void refreshSessionOnce()
  }, delay)
  logger.debug('session-refresh', `scheduled refresh in ${Math.round(delay / 1000)}s`)
}

export function initSessionRefresh(): void {
  if (!getSessionCsrfToken() || !getAccessExpiresAtUtc())
    return
  scheduleAccessRefresh()
}

export async function refreshSessionOnce(): Promise<boolean> {
  if (refreshInFlight)
    return refreshInFlight
  const csrf = getSessionCsrfToken()
  if (!csrf) {
    invalidateClientSessionAndRedirect()
    return false
  }
  const p = (async (): Promise<boolean> => {
    try {
      const result = await postAuthRefresh(csrf)
      if (!result.ok) {
        invalidateClientSessionAndRedirect()
        return false
      }
      setAuthSessionFromRefresh(result.data)
      scheduleAccessRefresh()
      return true
    }
    finally {
      refreshInFlight = null
    }
  })()
  refreshInFlight = p
  return p
}

// SEC-02: User profile api call uses session-scoped user_uuid instead of user_id e.g GET /api/users/me instead of GET /api/users/:id.
// SEC-04: Session CSRF token for mutating API calls; compared to `auth_sessions.csrf_token` on the server (SEC-04 remediated).
import type { AuthMeData, LoginResponseData } from '@/api/auth'

import type { RefreshResponseData } from '@/api/auth-refresh-internal'

const CSRF_KEY = 'secweb-helpdesk-session-csrf'
const CSRF_EXPIRES_AT_KEY = 'secweb-helpdesk-session-csrf-expires-at-utc'
const USER_KEY = 'secweb-helpdesk-session-user'
const ACCESS_EXPIRES_AT_KEY = 'secweb-helpdesk-access-expires-at-utc'

export type AuthUserSnapshot = {
  user_id: number
  user_uuid: string
  email: string
  role: string
  must_change_password: boolean
}

export function setAuthSessionFromLogin(data: LoginResponseData): void {
  if (typeof sessionStorage === 'undefined')
    return
  try {
    sessionStorage.setItem(CSRF_KEY, data.csrf_token)
    sessionStorage.setItem(CSRF_EXPIRES_AT_KEY, data.csrf_expires_at_utc)
    sessionStorage.setItem(ACCESS_EXPIRES_AT_KEY, data.access_expires_at_utc)
    const user: AuthUserSnapshot = {
      user_id: data.user_id,
      user_uuid: data.user_uuid,
      email: data.email,
      role: data.role,
      must_change_password: data.must_change_password,
    }
    sessionStorage.setItem(USER_KEY, JSON.stringify(user))
  }
  catch {
    /* ignore quota / private mode */
  }
}

/** After `POST /api/auth/refresh` — new CSRF, access window, and user ids from server. */
export function setAuthSessionFromRefresh(data: RefreshResponseData): void {
  if (typeof sessionStorage === 'undefined')
    return
  try {
    sessionStorage.setItem(CSRF_KEY, data.csrf_token)
    sessionStorage.setItem(CSRF_EXPIRES_AT_KEY, data.csrf_expires_at_utc)
    sessionStorage.setItem(ACCESS_EXPIRES_AT_KEY, data.access_expires_at_utc)
    const prev = getAuthUserSnapshot()
    if (prev) {
      const user: AuthUserSnapshot = {
        ...prev,
        user_id: data.user_id,
        user_uuid: data.user_uuid,
      }
      sessionStorage.setItem(USER_KEY, JSON.stringify(user))
    }
  }
  catch {
    /* ignore quota / private mode */
  }
}

export function getAccessExpiresAtUtc(): string | null {
  if (typeof sessionStorage === 'undefined')
    return null
  return sessionStorage.getItem(ACCESS_EXPIRES_AT_KEY)
}

export function getSessionCsrfToken(): string | null {
  if (typeof sessionStorage === 'undefined')
    return null
  return sessionStorage.getItem(CSRF_KEY)
}

/** Server CSRF window — used to rotate before mutating requests (independent of access token TTL). */
export function getCsrfExpiresAtUtc(): string | null {
  if (typeof sessionStorage === 'undefined')
    return null
  return sessionStorage.getItem(CSRF_EXPIRES_AT_KEY)
}

/** After `GET /api/auth/csrf-token` — new session CSRF and expiry only. */
export function setSessionCsrfPair(token: string, expiresAtUtc: string): void {
  if (typeof sessionStorage === 'undefined')
    return
  try {
    sessionStorage.setItem(CSRF_KEY, token)
    sessionStorage.setItem(CSRF_EXPIRES_AT_KEY, expiresAtUtc)
  }
  catch {
    /* ignore */
  }
}

/** After `GET /api/auth/me` — update `must_change_password` and identity fields (TASK.md §8.3). */
export function mergeAuthUserFromMe(data: AuthMeData): void {
  if (typeof sessionStorage === 'undefined')
    return
  try {
    const prev = getAuthUserSnapshot()
    if (!prev)
      return
    const user: AuthUserSnapshot = {
      user_id: data.user_id,
      user_uuid: data.user_uuid,
      email: data.email,
      role: data.role,
      must_change_password: data.must_change_password,
    }
    sessionStorage.setItem(USER_KEY, JSON.stringify(user))
  }
  catch {
    /* ignore quota / private mode */
  }
}

export function getAuthUserSnapshot(): AuthUserSnapshot | null {
  if (typeof sessionStorage === 'undefined')
    return null
  try {
    const raw = sessionStorage.getItem(USER_KEY)
    if (!raw)
      return null
    const parsed = JSON.parse(raw) as unknown
    if (!parsed || typeof parsed !== 'object')
      return null
    const o = parsed as Record<string, unknown>
    if (
      typeof o.user_id !== 'number'
      || typeof o.user_uuid !== 'string'
      || typeof o.email !== 'string'
      || typeof o.role !== 'string'
      || typeof o.must_change_password !== 'boolean'
    ) {
      return null
    }
    return {
      user_id: o.user_id,
      user_uuid: o.user_uuid,
      email: o.email,
      role: o.role,
      must_change_password: o.must_change_password,
    }
  }
  catch {
    return null
  }
}

export function clearAuthSession(): void {
  if (typeof sessionStorage === 'undefined')
    return
  sessionStorage.removeItem(CSRF_KEY)
  sessionStorage.removeItem(CSRF_EXPIRES_AT_KEY)
  sessionStorage.removeItem(USER_KEY)
  sessionStorage.removeItem(ACCESS_EXPIRES_AT_KEY)
}

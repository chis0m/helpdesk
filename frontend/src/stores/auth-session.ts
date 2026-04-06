// VULN-02: Snapshot stores numeric `user_id` for routes/state (pairs with IDOR on GET/PATCH /users/:id).
// VULN-05: Session CSRF token held for mutating API calls; server-side check is flawed (see backend VULN-05).
import type { LoginResponseData } from '@/api/auth'

const CSRF_KEY = 'secweb-helpdesk-session-csrf'
const USER_KEY = 'secweb-helpdesk-session-user'

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

export function getSessionCsrfToken(): string | null {
  if (typeof sessionStorage === 'undefined')
    return null
  return sessionStorage.getItem(CSRF_KEY)
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
  sessionStorage.removeItem(USER_KEY)
}

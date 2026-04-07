/**
 * One retry after 401 when session refresh can recover access (TASK.md §8.1).
 */
import { CSRF_HEADER } from './client'
import { refreshSessionOnce, shouldAttempt401Refresh } from './session-refresh'
import { getSessionCsrfToken } from '@/stores/auth-session'

export async function fetchWithSessionRefresh(
  url: string,
  init?: RequestInit,
): Promise<Response> {
  const res = await fetch(url, init)
  if (res.status !== 401 || !shouldAttempt401Refresh(url))
    return res
  const ok = await refreshSessionOnce()
  if (!ok)
    return res
  const headers = new Headers(init?.headers ?? undefined)
  const csrf = getSessionCsrfToken()
  if (csrf)
    headers.set(CSRF_HEADER, csrf)
  return fetch(url, { ...init, headers })
}

/**
 * API client helpers.
 *
 * SECURE-01: Mutating/session requests use `credentials: 'include'` (see callers). Session cookies are
 * HttpOnly (set by the backend; Secure when HTTPS/production). Not readable from JS.
 * VULN-05: `CSRF_HEADER` — client sends `X-CSRF-Token`; broken comparison / validation is server middleware.
 */
import { getApiBaseUrl } from './base-url'

export const CSRF_HEADER = 'X-CSRF-Token'

export function apiUrl(path: string): string {
  const base = getApiBaseUrl()
  const p = path.startsWith('/') ? path : `/${path}`
  return `${base}${p}`
}

export async function readJson(res: Response): Promise<unknown> {
  const text = await res.text()
  if (!text)
    return null
  try {
    return JSON.parse(text) as unknown
  }
  catch {
    return null
  }
}

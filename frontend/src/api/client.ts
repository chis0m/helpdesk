/**
 * API client helpers.
 *
 * VULN-01: Mutating/session requests use `credentials: 'include'` (see callers) — session cookies are
 * issued by the backend with weak flags; the SPA stores/sends them as the browser allows.
 * VULN-04: `CSRF_HEADER` — client sends `X-CSRF-Token`; broken comparison / validation is server middleware.
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

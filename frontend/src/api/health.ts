/** `GET /api/health` — raw JSON (no envelope); use for connectivity checks. */
import { apiUrl, readJson } from './client'
import { fetchWithSessionRefresh } from './session-fetch'
import { logger } from '@/utils/logger'

export async function fetchHealth(): Promise<
  | { ok: true; status: string; message: string }
  | { ok: false; status: number; message: string }
> {
  const url = apiUrl('/api/health')
  try {
    const res = await fetchWithSessionRefresh(url, {
      method: 'GET',
      headers: { Accept: 'application/json' },
    })
    const json = await readJson(res)
    logger.debug('api:health', `GET ${url} → ${res.status}`)
    if (!res.ok) {
      const msg = typeof json === 'object' && json !== null && 'message' in json && typeof (json as { message: unknown }).message === 'string'
        ? (json as { message: string }).message
        : 'Request failed'
      return { ok: false, status: res.status, message: msg }
    }
    if (!json || typeof json !== 'object')
      return { ok: false, status: res.status, message: 'Invalid response' }
    const o = json as { status?: unknown; message?: unknown }
    const status = typeof o.status === 'string' ? o.status : 'ok'
    const message = typeof o.message === 'string' ? o.message : 'backend is running'
    return { ok: true, status, message }
  }
  catch (e) {
    logger.warn('[api:health] request failed (network or unexpected error)', e)
    let message = e instanceof Error ? e.message : 'Could not reach the server'
    if (message === 'Failed to fetch')
      message = 'Could not reach the server'
    return { ok: false, status: 0, message }
  }
}

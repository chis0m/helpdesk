/**
 * Dev-only request/response logging for `fetch` (not OpenTelemetry tracing).
 * Redacts passwords and long tokens; strips noisy duplication in production builds.
 */
import { CSRF_HEADER } from './client'
import { isDebugVerbose, logger } from '@/utils/logger'

const SENSITIVE_BODY_KEYS = new Set([
  'password',
  'new_password',
  'current_password',
  'token',
  'access_token',
  'refresh_token',
  'csrf_token',
])

function redactForDevLog(input: unknown): unknown {
  if (input === null || input === undefined)
    return input
  if (Array.isArray(input))
    return input.map(redactForDevLog)
  if (typeof input !== 'object')
    return input
  const o = input as Record<string, unknown>
  const out: Record<string, unknown> = {}
  for (const [k, v] of Object.entries(o)) {
    const lk = k.toLowerCase()
    if (SENSITIVE_BODY_KEYS.has(lk)) {
      out[k] = v !== undefined && v !== null && String(v).length > 0 ? '[redacted]' : v
    }
    else if (typeof v === 'object' && v !== null) {
      out[k] = redactForDevLog(v) as unknown
    }
    else {
      out[k] = v
    }
  }
  return out
}

function headersObject(headers: HeadersInit | undefined): Record<string, string> {
  if (!headers)
    return {}
  if (headers instanceof Headers) {
    const o: Record<string, string> = {}
    headers.forEach((value, key) => {
      o[key] = value
    })
    return o
  }
  if (Array.isArray(headers)) {
    const o: Record<string, string> = {}
    for (const [k, v] of headers)
      o[k] = v
    return o
  }
  return { ...headers }
}

function redactHeadersForLog(h: Record<string, string>): Record<string, string> {
  const out: Record<string, string> = {}
  for (const [name, value] of Object.entries(h)) {
    const ln = name.toLowerCase()
    if (ln === CSRF_HEADER.toLowerCase() && value.length > 12) {
      out[name] = `${value.slice(0, 8)}…`
    }
    else {
      out[name] = value
    }
  }
  return out
}

function parseRequestBody(body: BodyInit | null | undefined): unknown {
  if (body === null || body === undefined)
    return undefined
  if (typeof body === 'string') {
    try {
      return redactForDevLog(JSON.parse(body) as unknown)
    }
    catch {
      return '(non-JSON body)'
    }
  }
  return '(non-string body; not logged)'
}

export function devLogRequest(
  context: string,
  method: string,
  url: string,
  init?: RequestInit,
): void {
  if (!isDebugVerbose())
    return
  const headers = redactHeadersForLog(headersObject(init?.headers))
  const body = parseRequestBody(init?.body ?? undefined)
  logger.debug(context, '→ request', { method, url, headers, body })
}

export async function devLogResponse(context: string, url: string, res: Response): Promise<void> {
  if (!isDebugVerbose())
    return
  const clone = res.clone()
  const text = await clone.text()
  let body: unknown = text.length === 0 ? null : text
  if (text.length > 0) {
    try {
      body = redactForDevLog(JSON.parse(text) as unknown)
    }
    catch {
      body = text.length > 500 ? `${text.slice(0, 500)}…` : text
    }
  }
  logger.debug(context, '← response', { url, status: res.status, body })
}

/**
 * `fetch` with dev-only logs of request + response body (redacted). Production: plain `fetch`.
 */
export async function loggedFetch(
  context: string,
  url: string,
  init?: RequestInit,
): Promise<Response> {
  const method = (init?.method ?? 'GET').toUpperCase()
  devLogRequest(context, method, url, init)
  const res = await fetch(url, init)
  await devLogResponse(context, url, res)
  return res
}

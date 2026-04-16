import type { InjectionKey, Ref } from 'vue'
import { apiUrl, readJson } from '@/api/client'
import { logger } from '@/utils/logger'

const FETCH_TIMEOUT_MS = 4000

export type AppDetail = {
  app_name: string
  environment: string
  version: string
}

const FALLBACK: AppDetail = {
  app_name: 'SecWeb HelpDesk',
  environment: 'development',
  version: 'v1.0.0',
}

let cache: AppDetail | null = null
let inflight: Promise<AppDetail> | null = null

/** First word of `app_name` for short labels (e.g. hero, “X products”). */
export function brandShortFromAppName(name: string): string {
  const t = name.trim()
  if (!t)
    return 'SecWeb'
  return t.split(/\s+/)[0] ?? 'SecWeb'
}

export async function loadAppDetail(): Promise<AppDetail> {
  if (cache) {
    logger.debug('app-detail', 'using cached app detail', { app_name: cache.app_name })
    return cache
  }
  if (inflight) {
    logger.debug('app-detail', 'awaiting in-flight GET /api/public/app-detail')
    return inflight
  }

  inflight = (async () => {
    const url = apiUrl('/api/public/app-detail')
    try {
      console.info('[app-detail] fetch started', { url, timeoutMs: FETCH_TIMEOUT_MS })
      const t0 = typeof performance !== 'undefined' ? performance.now() : 0
      const ctrl = new AbortController()
      const tid = typeof window !== 'undefined'
        ? window.setTimeout(() => ctrl.abort(), FETCH_TIMEOUT_MS)
        : undefined
      const res = await fetch(url, {
        method: 'GET',
        headers: { Accept: 'application/json' },
        signal: ctrl.signal,
      })
      if (typeof tid === 'number')
        window.clearTimeout(tid)
      const elapsedMs = typeof performance !== 'undefined' ? Math.round(performance.now() - t0) : undefined
      const json = await readJson(res)
      if (!res.ok || !json || typeof json !== 'object') {
        logger.warn(
          `[app-detail] fetch finished with error status or invalid body (HTTP ${res.status}) — using fallback`,
        )
        cache = { ...FALLBACK }
        return cache
      }
      const root = json as Record<string, unknown>
      // Backend uses response.Success → fields under `data`; support flat shape for compatibility.
      const payload
        = root.data !== null
          && typeof root.data === 'object'
          && !Array.isArray(root.data)
          ? (root.data as Record<string, unknown>)
          : root
      const app_name = typeof payload.app_name === 'string' && payload.app_name.trim()
        ? payload.app_name.trim()
        : FALLBACK.app_name
      const environment = typeof payload.environment === 'string' && payload.environment.trim()
        ? payload.environment.trim()
        : FALLBACK.environment
      const version = typeof payload.version === 'string' && payload.version.trim()
        ? payload.version.trim()
        : FALLBACK.version
      cache = { app_name, environment, version }
      console.info('[app-detail] fetch successful', {
        ms: elapsedMs,
        app_name,
        environment,
        version,
      })
      return cache
    }
    catch (e) {
      const aborted
        = e instanceof DOMException && e.name === 'AbortError'
          || (e instanceof Error && e.name === 'AbortError')
      if (aborted) {
        logger.warn(`[app-detail] fetch timed out after ${FETCH_TIMEOUT_MS}ms — using fallback`)
      }
      else {
        logger.warn('[app-detail] fetch failed — using fallback', e)
      }
      cache = { ...FALLBACK }
      return cache
    }
    finally {
      inflight = null
    }
  })()

  return inflight
}

export type AppBrandRefs = { appName: Ref<string>; brandShort: Ref<string> }

export const appBrandKey: InjectionKey<AppBrandRefs> = Symbol('appBrand')

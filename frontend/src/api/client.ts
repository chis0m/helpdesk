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

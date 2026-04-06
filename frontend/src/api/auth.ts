import { apiUrl, CSRF_HEADER, readJson } from './client'
import type { ApiErrorEnvelope, ApiSuccessEnvelope } from './types'
import { logger } from '@/utils/logger'

export interface PublicCsrfData {
  csrf_token: string
  csrf_expires_at_utc: string
}

export interface LoginResponseData {
  user_id: number
  user_uuid: string
  email: string
  role: string
  must_change_password: boolean
  access_expires_at_utc: string
  csrf_token: string
  csrf_expires_at_utc: string
}

export interface SignupResponseData {
  user_uuid: string
  email: string
  redirect_to: string
}

export type SignupRequestBody = {
  email: string
  password: string
  first_name: string
  last_name: string
  middle_name?: string
}

function errorMessage(body: unknown): string {
  if (!body || typeof body !== 'object')
    return 'Request failed'
  const env = body as ApiErrorEnvelope
  if (typeof env.message === 'string' && env.message.length > 0)
    return env.message
  if (typeof env.error === 'string' && env.error.length > 0)
    return env.error
  return 'Request failed'
}

export async function fetchPublicCsrfToken(): Promise<{ ok: true; token: string } | { ok: false; message: string }> {
  const url = apiUrl('/api/auth/public-csrf-token')
  const res = await fetch(url, {
    method: 'GET',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  const json = await readJson(res)
  logger.debug('api:auth', `GET ${url} → ${res.status}`)
  if (!res.ok) {
    logger.debug('api:auth', 'public CSRF error body', json)
    return { ok: false, message: errorMessage(json) }
  }

  const env = json as ApiSuccessEnvelope<PublicCsrfData>
  logger.debug('api:auth', 'public CSRF response (dev — includes token)', env.data ?? json)
  const token = env.data?.csrf_token
  if (!token)
    return { ok: false, message: 'Missing CSRF token' }
  return { ok: true, token }
}

export async function loginRequest(
  email: string,
  password: string,
  publicCsrf: string,
): Promise<{ ok: true; data: LoginResponseData } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/auth/login')
  const res = await fetch(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: publicCsrf,
    },
    body: JSON.stringify({ email, password }),
  })
  const json = await readJson(res)
  logger.debug('api:auth', `POST ${url} → ${res.status}`, {
    email,
    'X-CSRF-Token (public)': publicCsrf,
  })
  if (!res.ok) {
    logger.debug('api:auth', 'login error envelope (dev)', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }

  const env = json as ApiSuccessEnvelope<LoginResponseData>
  if (!env.data) {
    logger.debug('api:auth', 'login missing data envelope (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug(
    'api:auth',
    'login success envelope (dev — includes session CSRF in data)',
    json,
  )
  return { ok: true, data: env.data }
}

export async function signupRequest(
  body: SignupRequestBody,
  publicCsrf: string,
): Promise<{ ok: true; data: SignupResponseData } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/auth/signup')
  const payload: Record<string, string> = {
    email: body.email.trim(),
    password: body.password,
    first_name: body.first_name.trim(),
    last_name: body.last_name.trim(),
  }
  const mid = body.middle_name?.trim()
  if (mid)
    payload.middle_name = mid

  const res = await fetch(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
      [CSRF_HEADER]: publicCsrf,
    },
    body: JSON.stringify(payload),
  })
  const json = await readJson(res)
  logger.debug('api:auth', `POST ${url} → ${res.status}`, {
    email: payload.email,
    first_name: payload.first_name,
    last_name: payload.last_name,
    middle_name: payload.middle_name,
    'X-CSRF-Token (public)': publicCsrf,
  })
  if (!res.ok) {
    logger.debug('api:auth', 'signup error envelope (dev)', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }

  const env = json as ApiSuccessEnvelope<SignupResponseData>
  if (!env.data) {
    logger.debug('api:auth', 'signup missing data envelope (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:auth', 'signup success envelope (dev)', json)
  return { ok: true, data: env.data }
}

export interface LogoutResponseData {
  redirect_to: string
}

export async function logoutRequest(
  sessionCsrf: string,
): Promise<{ ok: true; data: LogoutResponseData } | { ok: false; status: number; message: string }> {
  const url = apiUrl('/api/auth/logout')
  const res = await fetch(url, {
    method: 'POST',
    credentials: 'include',
    headers: {
      Accept: 'application/json',
      [CSRF_HEADER]: sessionCsrf,
    },
  })
  const json = await readJson(res)
  logger.debug('api:auth', `POST ${url} → ${res.status}`, {
    'X-CSRF-Token (session)': sessionCsrf,
  })
  if (!res.ok) {
    logger.debug('api:auth', 'logout error envelope (dev)', json)
    return { ok: false, status: res.status, message: errorMessage(json) }
  }

  const env = json as ApiSuccessEnvelope<LogoutResponseData>
  if (!env.data) {
    logger.debug('api:auth', 'logout missing data envelope (dev)', json)
    return { ok: false, status: res.status, message: 'Invalid response' }
  }
  logger.debug('api:auth', 'logout success envelope (dev)', json)
  return { ok: true, data: env.data }
}

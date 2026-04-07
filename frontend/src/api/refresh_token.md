# Session cookies, CSRF, and refresh

This app uses **HttpOnly-style session cookies** from the backend (`access_token`, `refresh_token`). JavaScript does **not** read those tokens; the browser sends them when `credentials: 'include'` is set on `fetch`.

## What is stored in `sessionStorage`?

After login (or `POST /api/auth/refresh`), the SPA keeps:

- **Session CSRF** — for `X-CSRF-Token` on mutating requests.
- **`access_expires_at_utc`** — when the access window ends (for scheduling refresh).
- **User snapshot** — `user_id`, `user_uuid`, email, role, etc., for UI and routes.

See `stores/auth-session.ts` (`setAuthSessionFromLogin`, `setAuthSessionFromRefresh`, `getSessionCsrfToken`, `getAccessExpiresAtUtc`).

## Two ways refresh runs

### 1. Proactive (timer)

`session-refresh.ts` schedules **`POST /api/auth/refresh`** shortly **before** `access_expires_at_utc` using `setTimeout` (`scheduleAccessRefresh`). The timer is:

- Started after **login** (`LoginView.vue`).
- Rescheduled after each **successful** refresh.
- Restored on **full page load** if the tab still has CSRF + expiry (`initSessionRefresh` in `main.ts`).

### 2. Reactive (401)

`session-fetch.ts` wraps most API calls. If a response is **401** and the URL is **not** a “public auth” path (login, signup, public CSRF, refresh itself, password reset flows, invite verify/accept), it calls **`refreshSessionOnce()`** once, updates CSRF from storage, and **retries the original request once** with the new header.

### Raw refresh request

`auth-refresh-internal.ts` calls **`fetch` directly** for `POST /api/auth/refresh` (no `fetchWithSessionRefresh`). That avoids a 401 on refresh triggering “refresh again” in a loop.

## When refresh fails

If `POST /api/auth/refresh` returns **401/403**, the session is cleared, the refresh timer is cleared, and the app navigates to **login** with a `redirect` query (`registerSessionRefreshFailure` in `main.ts`).

## Logout

`perform-logout.ts` clears the refresh timer, calls `POST /api/auth/logout` when possible, then clears `sessionStorage` and routes away.

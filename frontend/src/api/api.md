# API layer (`src/api`)

The API folder holds **HTTP helpers** and **one module per backend area** (auth, tickets, users, admin, invites, health). Nothing here renders UI.

## How requests work

1. **Base URL** — `base-url.ts` reads `VITE_API_BASE_URL` (default `http://localhost:8080`). All paths are built with `apiUrl('/api/...')` from `client.ts`.

2. **Cookies** — Session calls use `credentials: 'include'` so the browser sends `access_token` / `refresh_token` cookies to the API origin (needed when the Vite dev server and API are on different ports).

3. **CSRF** — Mutating calls send `X-CSRF-Token` (`CSRF_HEADER` in `client.ts`). Public flows use a token from `GET /api/auth/public-csrf-token`; logged-in flows use the session CSRF from login/refresh JSON, stored in `sessionStorage` (see `stores/auth-session.ts` and `refresh_token.md`).

4. **JSON** — Responses are parsed with `readJson()`; most endpoints return success/error **envelopes** (`types.ts`).

5. **Fetch wrapper** — Most modules use `fetchWithSessionRefresh` from `session-fetch.ts` instead of raw `fetch`. On **401**, it can run a **session refresh** and retry once (see `refresh_token.md`). `POST /api/auth/refresh` itself uses **raw** `fetch` in `auth-refresh-internal.ts` so refresh never loops.

## File map

| File | Role |
|------|------|
| `client.ts` | `apiUrl`, `CSRF_HEADER`, `readJson` |
| `base-url.ts` | `getApiBaseUrl()` |
| `types.ts` | Shared API envelope types |
| `session-fetch.ts` | `fetchWithSessionRefresh` (401 → refresh → retry) |
| `session-refresh.ts` | Timer-based refresh + `refreshSessionOnce` mutex |
| `auth-refresh-internal.ts` | Raw `POST /api/auth/refresh` |
| `auth.ts` | Login, logout, me, sessions, CSRF, password, etc. |
| `tickets.ts` | Tickets CRUD, search, comments |
| `users.ts` | profile `GET/PATCH` |
| `admin.ts` | Admin directory endpoints |
| `invites.ts` | Invite verify / accept |
| `health.ts` | `GET /api/health` |

## Related docs

- **`refresh_token.md`** — Cookies, CSRF, refresh timer, and 401 handling in detail.

# Stores (`src/stores`)

Lightweight **client-side state** — not Pinia or Vuex. One module:

## `auth-session.ts`

**Session snapshot** for the logged-in user (backed by `sessionStorage`):

- Session **CSRF** for API mutating calls.
- **`access_expires_at_utc`** for refresh scheduling.
- **User fields** (`user_id`, `user_uuid`, email, role, `must_change_password`) for UI and router helpers.

Functions: `setAuthSessionFromLogin`, `setAuthSessionFromRefresh`, `mergeAuthUserFromMe` (after `GET /api/auth/me` when clearing `must_change_password`), `getSessionCsrfToken`, `getAccessExpiresAtUtc`, `getAuthUserSnapshot`, `clearAuthSession`.

**Cookies** (access/refresh) are **not** stored here — only what the API returns in JSON plus CSRF/expiry for the client.

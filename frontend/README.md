# SecWeb Helpdesk — frontend

Vue 3 + TypeScript SPA served by **Vite**. It talks to the Go API over HTTP with **cookie-based sessions** and **CSRF** on mutating requests.

## Run locally

```bash
npm install
npm run dev
```

Default API base: `http://localhost:8080` unless `VITE_API_BASE_URL` is set (see `src/api/base-url.ts`).

```bash
npm run build   # typecheck + production bundle
npm run lint
```

## High-level architecture

```
main.ts
  ├─ vue-router (router/index.ts)
  ├─ session refresh init + failure → login (session-refresh)
  └─ App.vue → <RouterView />

Views (pages)  ──►  API modules (fetch + envelopes)  ──►  Backend
      │                        │
      └─ stores (auth snapshot)
      └─ components (layout, tickets, admin)
```

- **Router** — Public routes (`/`, `/login`, …) and dashboard (`/dashboard/...`) using `AppLayout` or `AuthLayout`.
- **API** — `src/api/` builds URLs, sends cookies (`credentials: 'include'`), CSRF headers, and uses `fetchWithSessionRefresh` for 401 recovery. See `src/api/api.md` and `src/api/refresh_token.md`.
- **Stores** — `auth-session` mirrors what the SPA needs from login/refresh JSON.
- **Must change password** — If login returns `must_change_password: true` (see `backend/TASK.md` §8.3), the router sends the user to `/change-password` until they submit `POST /api/auth/change-password` and `GET /api/auth/me` confirms the flag is cleared.
- **Types** — `src/types/` for UI/domain types; `src/api/types.ts` for API envelopes.

## Folder guide

| Folder | Role |
|--------|------|
| `src/api/` | HTTP client, auth, tickets, admin, invites, health |
| `src/components/` | Reusable Vue components (layout, tickets, admin, icons) |
| `src/constants/` | `routes.ts` — path helpers aligned with the router |
| `src/layouts/` | `AppLayout`, `AuthLayout` shells |
| `src/router/` | Route table + lazy-loaded views |
| `src/stores/` | Session + small client state |
| `src/types/` | Shared TS types for UI |
| `src/utils/` | Logout, logger, date/ticket formatting |
| `src/views/` | Page-level views per route |

## Documentation

- `src/api/api.md` — API module overview
- `src/api/refresh_token.md` — Cookies, CSRF, refresh timer, 401 path
- `src/components/components.md`, `src/stores/store.md`, `src/types/types.md`, `src/utils/utils.md`, `src/views/views.md`

## Backend contract

Authoritative server behavior and JSON shapes: **`../backend/API.md`** (and `TASK.md` for integration notes).

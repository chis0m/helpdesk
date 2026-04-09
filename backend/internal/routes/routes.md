# Routes (`internal/routes`)

Single function **`Register(r *gin.Engine, c *container.Container)`** — defines the entire HTTP surface and **middleware order**.

## Top level

- **`GET /`** — welcome JSON (outside `/api`).
- **`/api`** group — all API routes below.

## Unauthenticated `/api` routes

| Path | Middleware / notes |
|------|-------------------|
| `GET /api/health` | None. |
| `GET /api/auth/public-csrf-token` | Issue token for subsequent public POSTs. |
| `POST /api/auth/login` | IP rate limit → public CSRF → login. |
| `POST /api/auth/signup` | IP rate limit → public CSRF. |
| `POST /api/auth/forgot-password` | Rate limit → public CSRF. |
| `POST /api/auth/reset-password` | Rate limit → public CSRF. |
| `GET /api/invites/verify` | Invite rate limit. |
| `POST /api/invites/accept` | Invite rate limit → public CSRF. |
| `POST /api/auth/refresh` | **Refresh cookie** + PASETO middleware + **session CSRF** (not public CSRF). |

## Protected `/api` group

`protected.Use(AuthRequired, ActiveSessionRequired)` — every route below requires valid **access cookie** + **active session**.

Mutating routes add **`CSRFRequired`** (session CSRF header). Read-only routes often omit CSRF (GET-only).

### Auth & session management

- `GET /api/auth/csrf-token`, `/auth/me`, `/auth/sessions`
- `POST` revoke other sessions, `POST` logout, `POST` change-password, `DELETE` session by id — all CSRF-protected.

### Users & admin

- **Registration:** `POST /api/auth/signup` (public; see unauthenticated routes). There is no `POST /api/users`.
- `GET /api/admin/users`, `POST /api/admin/staff`, `POST /api/admin/invites/staff` — admin/staff flows (CSRF where mutating).
- `GET/PATCH /users/:id` — profile read/update (PATCH CSRF).

### Tickets & comments

- **`GET /api/tickets/search`** is registered **before** **`GET /api/tickets/:id`** so `search` is not captured as an id (important for routing).
- Full CRUD on tickets + nested comments; mutating routes use CSRF.

- `PATCH /api/admin/users/:user_id/role` — role change (controller enforces privileged roles).

## Rate limiters

Constructed per **`Register` call** (new limiter instances each boot). Limits differ by endpoint sensitivity (login vs invite verify).

## VULN tracking

A **`VULN-06`** tag may appear above the `protected` group in source — see `ca2/Vulnerability.md` for narrative; routes file only carries the short marker.

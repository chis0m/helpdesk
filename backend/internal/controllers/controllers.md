# Controllers (`internal/controllers`)

**HTTP adapters** for Gin: parse path/query/body, call **services**, return JSON via **`internal/response`** (`Success`, `Failure`, `FailureWithAbort`). Keep policy and branching in services; these files mostly map HTTP ↔ calls.

## `HealthController`

`Ping` — simple health JSON for load balancers and smoke tests.

## `AuthController`

- **Public:** `PublicAuthCSRFToken`, `Login`, `Signup`, `ForgotPassword`, `ResetPassword`.
- **Refresh:** uses refresh cookie context set by middleware before handler runs.
- **Protected:** `Me`, `CSRFToken`, session list/revoke/logout, `ChangePassword`.

**Cookies:** `auth_cookie.go` helpers set/clear access + refresh cookies (baseline flags documented in CA `Vulnerability.md`). Login/refresh responses may also return CSRF metadata for the SPA.

Handlers read **`middleware` context keys** for user UUID, role, session id.

## `UserController`

Admin user list, **CreateStaff** (`POST /api/admin/staff`, admin/super_admin; optional `role` staff|admin, **admin** only for super_admin actor), **GetByID** / **UpdateByID** by path id, **UpdateRoleByUserID**. End-user accounts are created via **`AuthController.Signup`** (`POST /api/auth/signup`), not this controller.

## `InviteController`

Verify invite (often public GET with token), accept invite (public POST with body + public CSRF), create staff invite (protected + CSRF).

## `TicketController`

Tickets CRUD, status, assign, delete, **search**, and **nested comments** (list/add/update/delete). Uses small helpers for parsing ids and formatting ticket JSON for responses.

Parses authenticated user from Gin context for actions that need actor UUID/role (e.g. list scoping, comment permissions).

## Error handling pattern

Typical flow: service returns `error` → controller switches on `errors.Is` for known sentinel (`gorm.ErrRecordNotFound`, `ErrForbidden`, etc.) → appropriate HTTP code and `FailureWithAbort` to stop middleware chain.

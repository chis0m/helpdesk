# Controllers (`internal/controllers`)

**HTTP adapters** for Gin: parse path/query/body, call **services**, return JSON via **`internal/response`** (`Success`, `Failure`, `FailureWithAbort`). They should stay **thin** — business rules belong in services.

## `HealthController`

`Ping` — simple health JSON for load balancers and smoke tests.

## `AuthController`

- **Public:** `PublicAuthCSRFToken`, `Login`, `Signup`, `ForgotPassword`, `ResetPassword`.
- **Refresh:** uses refresh cookie context set by middleware before handler runs.
- **Protected:** `Me`, `CSRFToken`, session list/revoke/logout, `ChangePassword`.

**Cookies:** `auth_cookie.go` helpers set/clear access + refresh cookies (baseline flags documented in CA `Vulnerability.md`). Login/refresh responses may also return CSRF metadata for the SPA.

Handlers read **`middleware` context keys** for user UUID, role, session id.

## `UserController`

User creation, admin user list, staff creation, **GetByID** / **UpdateByID** by path id, **UpdateRoleByUserID** for admin role management. Authorization for admin-only actions is enforced in controller/service layers depending on the endpoint.

## `InviteController`

Verify invite (often public GET with token), accept invite (public POST with body + public CSRF), create staff invite (protected + CSRF).

## `TicketController`

Tickets CRUD, status, assign, delete, **search**, and **nested comments** (list/add/update/delete). Uses small helpers for parsing ids and formatting ticket JSON for responses.

Parses authenticated user from Gin context for actions that need actor UUID/role (e.g. list scoping, comment permissions).

## Error handling pattern

Typical flow: service returns `error` → controller switches on `errors.Is` for known sentinel (`gorm.ErrRecordNotFound`, `ErrForbidden`, etc.) → appropriate HTTP code and `FailureWithAbort` to stop middleware chain.

## Adding an endpoint

1. Add route in `internal/routes/routes.go`.  
2. Add method on controller; inject any new service via **container** if needed.  
3. Define request struct in `internal/requests` with `binding` tags.

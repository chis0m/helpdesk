# Services (`internal/services`)

Coordinates repositories, auth helpers, mail notifiers, and config. Controllers call in here with parsed requests; responses go back through `response` helpers — no Gin handlers in this package.

## `AuthService` (`auth_service.go`)

Central authentication lifecycle:

- **Login** — verify Argon2id password, create/update **`auth_sessions`** row (refresh JTI, CSRF fields, user-agent/IP metadata), mint access + refresh PASETOs bound to session id.
- **Signup** — create `user` row with default role; no session until login.
- **Refresh** — validate refresh token + session, **rotate** refresh JTI on the session row (replay detection), re-issue access (+ optionally CSRF metadata in response).
- **Logout / revoke** — mark session revoked or delete row as implemented.
- **Session list** — returns safe subset for “my sessions” UI (mark current session).
- **Change password** — re-hash, update user, may invalidate other sessions depending on implementation.
- **Forgot / reset password** — generates opaque token, stores **hash** in `password_resets`, builds URL with `FRONTEND_URL`, calls **`PasswordResetNotifier`** (logs link in CA). Reset validates token, expiry, single-use.

Password-reset and invite tokens use **SHA-256 hashes** of raw secrets in DB; raw token only appears in email/log.

Exported **errors** (e.g. `ErrInvalidCredentials`, `ErrInvalidRefreshToken`) let controllers map to HTTP status codes.

## `UserService` (`user_service.go`)

User CRUD-style operations: get by id/UUID, list (admin), create from requests, **`CreateStaffFromRequest`** (`POST /api/admin/staff`) with optional **`role`**: `staff` (default) or `admin` — only **super_admin** may set **`admin`** (`ErrCreateStaffAdminRequiresSuperAdmin`), **role updates** with role-gating (who may promote whom). Some endpoints intentionally reflect **baseline IDOR** behavior documented in `ca2/Vulnerability.md` — secure branch should tighten checks here + in controllers.

## `InviteService` (`invite_service.go`)

Staff invitation flow:

- **CreateStaffInvite** — admin/super_admin only; optional **`role`**: `staff` (default) or `admin`. Only **super_admin** may set `admin`. Ensures email not registered and no conflicting pending invite; generates raw token, stores **hash**, sends notifier with accept URL.
- **Verify** — validate token hash, expiry, used flag.
- **Accept** — create staff user, mark invite used (transactional expectations in repo layer).

Uses **`FRONTEND_URL`** to build links consistent with the SPA.

## `TicketService` (`ticket_service.go`)

Ticket and comment operations backed by `TicketRepository`, `TicketCommentRepository`, `UserRepository`.

- **`ListForActor`** — **authorization for listing**: admins see all (subject to filters); non-admins are scoped to tickets they **report** or are **assigned** to. This is the main non-trivial policy method.
- **By-id methods** (`GetByID`, updates, comments, etc.) — repository lookups by id; baseline intentionally **does not** re-check reporter/assignee on every by-id call.
- **Search** — delegates to repository search implementation.

Domain errors: invalid status transition, forbidden comment action, forbidden list filters.

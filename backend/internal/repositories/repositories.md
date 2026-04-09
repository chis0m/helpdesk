# Repositories (`internal/repositories`)

**Data access layer** on top of **GORM**: one repository per aggregate / table group. They accept `*gorm.DB` from the container and expose methods used by services. **No HTTP or PASETO** here.

## `UserRepository`

Users table: create, get by id/email/UUID, list with pagination and optional role filter, update profile fields, update role, password hash updates.

## `AuthSessionRepository`

`auth_sessions` table lifecycle:

- Insert session on login (session id UUID, user id, refresh JTI, CSRF token + expiry, optional UA/IP).
- Fetch by session id for middleware and refresh.
- Update refresh JTI / CSRF fields on rotation.
- Revoke or delete sessions; list sessions for a user.

This table is the **source of truth** for whether an access token’s session is still valid (`ActiveSessionRequired`).

## `PasswordResetRepository`

Stores **hashed** password-reset tokens, expiry, used flag, user association. Auth service generates raw token and URL; only hash hits the DB.

## `InviteRepository`

Invite rows: create pending invite, find by token hash, mark used, check pending by email.

## `TicketRepository`

Tickets: CRUD by id, list with filters and optional **scope** (`ScopeToUserID` for reporter/assignee narrowing), assign/unassign, soft-delete if modeled. May expose a **keyword search** method used by the service (baseline implementation details matter for SQL safety — see app vulnerability notes).

## `TicketCommentRepository`

Comments: add, list by ticket, get by id, update/delete with ticket+comment id.

## Conventions

- Return **`gorm.ErrRecordNotFound`** where appropriate so services map to 404.
- Prefer **explicit column lists** or structs over unscoped global updates for clarity.
- Transactions: if a service operation needs atomicity across tables, either start `tx` in the service and pass `tx` into repo methods, or add a repository method that runs `db.Transaction(...)`.

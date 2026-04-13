# Middleware (`internal/middleware`)

HTTP **cross-cutting** concerns: authentication, session validity, CSRF, and rate limiting. Handlers assume context values set by these layers (see **context keys** below).

## Context keys — `context_keys.go`

Gin stores typed values for downstream handlers and nested middleware:

| Key | Set by | Typical use |
|-----|--------|-------------|
| `CtxAuthPayload` | `AuthRequired`, `RefreshTokenRequired` | `*auth.Payload` (claims). |
| `CtxUserUUID` | `AuthRequired` | User UUID string from token `sub`. |
| `CtxUserRole` | `AuthRequired` | Role string for authorization branches. |
| `CtxSessionID` | `AuthRequired`, `RefreshTokenRequired` | Session id from token `sess_id`. |
| `CtxTokenJTI` | `AuthRequired`, `RefreshTokenRequired` | Current access/refresh JTI for rotation checks. |

Controllers use **`middleware.CtxUserUUID`** etc. when calling services.

## `AuthRequired(tokenMaker, accessCookieName)`

1. Reads access token from the **named cookie**.
2. Verifies PASETO and expiry, checks **access** token type.
3. Loads **`CtxAuthPayload`**, **`CtxUserUUID`**, **`CtxUserRole`**, **`CtxSessionID`**, **`CtxTokenJTI`**.

Used on the **`protected`** route group together with `ActiveSessionRequired`.

## `ActiveSessionRequired(sessionRepo)`

Ensures the session id in context still refers to a **non-revoked** row in `auth_sessions`. Blocks stale tokens after logout/revoke. Runs **after** `AuthRequired` so session id is known.

## `RefreshTokenRequired(tokenMaker, refreshCookieName)`

For **`POST /auth/refresh` only**: reads refresh cookie, verifies **refresh** token type, sets payload/session/JTI in context. Does **not** require an active session check in the same way as access (refresh flow validates + rotates in controller/service).

## `CSRFRequired(sessionRepo, headerName)`

Runs on **unsafe HTTP methods** (not GET/HEAD/OPTIONS). For authenticated session CSRF:

1. Requires non-empty CSRF header.
2. Loads session by id from context; checks CSRF token pointer + expiry on the **session row**.

Baseline CA behavior: see **`Vulnerability.md`** VULN-04 — header is **not** compared to stored secret in code (intentional weak branch). On a secure branch, add constant-time comparison to `session.CSRFToken`.

## `PublicAuthCSRFRequired(store, headerName)`

For **pre-login** routes: validates token issued by **`PublicAuthCSRFStore`**, consumes it (one-time), and checks TTL. Stronger than session CSRF on baseline — login/signup depend on this.

## `IPRateLimiter`

In-memory **per-IP** sliding window (mutex + map). Used on login, signup, forgot-password, invite verify/accept. Per-process; multiple API replicas require a shared limiter instead of this map.

## Typical request chains

- **Public POST** — `RateLimit` → `PublicAuthCSRFRequired` → controller.  
- **Protected GET** — `AuthRequired` → `ActiveSessionRequired` → controller.  
- **Protected mutating** — `AuthRequired` → `ActiveSessionRequired` → `CSRFRequired` → controller.

Order is important: CSRF middleware expects session id in context from `AuthRequired`.

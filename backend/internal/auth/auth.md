# Auth (`internal/auth`)

Cryptography and token **primitives** shared by middleware, services, and seeding. Not the HTTP layer (that is `controllers` + `middleware`).

## Password hashing — `argon2id.go`

- **`HashPassword(password)`** / **`VerifyPassword(password, hash)`** — Argon2id via `golang.org/x/crypto/argon2`. Used at signup, password change, reset, and admin seed.

## PASETO — `paseto.go`

- **`PasetoMaker`** implements **`MakerInterface`**: create/verify symmetric **v2.local** tokens.
- **Key size** — `NewPasetoMaker` requires exactly **32 bytes** (the `PASETO_SYMMETRIC_KEY` string in config must match this length).
- **`CreateAccessToken` / `CreateRefreshToken`** — set `TokenType` on claims then encrypt payload.
- **`VerifyToken` / `VerifyTokenType`** — decrypt, check expiry, ensure type matches (access vs refresh).

## Payload & claims — `payload.go`

- **`Claims`** — input to mint a token (issuer, subject=user UUID string, role, audience, duration, token type, **session id**).
- **`Payload`** — what gets encrypted: `iss`, `sub`, `role`, `aud`, `iat`, `exp`, `jti`, `type`, `sess_id`.
- Each token gets a random **JTI**; refresh rotation ties to session rows in DB (see `AuthService` + `AuthSessionRepository`).

Helpers: **`IsExpired`**, **`IsValidJTI`**, **`IsValidSessionID`**, **`IsValid`** (combined check for refresh replay logic).

## Constants — `types.go`

- Cookie names: **`access_token`**, **`refresh_token`**.
- CSRF header: **`X-CSRF-Token`** (must stay aligned with CORS `AllowHeaders` in `boot`).
- **`TokenPair`** — access/refresh strings + expiry metadata returned to controllers for cookie + JSON responses.

## Public auth CSRF — `public_auth_csrf_store.go`

Short-lived, **one-time-style** tokens for **unauthenticated** POST endpoints (login, signup, etc.). Stored **in memory** in the container; **not** the same as session CSRF stored on `auth_sessions` (see `middleware/public_auth_csrf_middleware.go`).

Production note: replace with Redis or similar if you need multi-instance consistency.

## Mental model for new engineers

1. **Password** → Argon2id hash in DB.  
2. **Session** → row in `auth_sessions` (refresh JTI, CSRF fields, revocation).  
3. **Access token** → short PASETO in cookie, validated by `AuthRequired`.  
4. **Refresh token** → longer PASETO in cookie, validated by `RefreshTokenRequired` + session row.  
5. **CSRF** → two tracks: **public** (pre-login) vs **session** (post-login); both use the same header name but different storage/validation.

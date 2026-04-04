# Middleware

### Access + Refresh Middleware

`AuthRequired` validates access token from cookie.  
`RefreshTokenRequired` validates refresh token for refresh endpoint.

### Active Session Middleware

`ActiveSessionRequired` checks `session_id` against `auth_sessions` table.  
It blocks revoked or missing sessions.

### CSRF Middleware

`CSRFRequired` runs on state-changing methods only.  
Baseline behavior is intentionally weak for CA vulnerability demo.

`PublicAuthCSRFRequired` validates public auth CSRF tokens (login/signup) using in-memory synchronizer storage.

### Rate Limit Middleware

`IPRateLimiter` applies simple in-memory per-IP request throttling.  
For production, I use Redis/shared limiter.

This implementation was guided by:
- https://www.alexedwards.net/blog/how-to-rate-limit-http-requests
- https://pkg.go.dev/sync#Mutex
- https://go.dev/doc/faq#atomic_maps

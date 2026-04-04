// Package middleware contains request guards for auth/session/CSRF.
//
// |--------------------------------------------------------------------------
// | Access + Refresh Middleware
// |--------------------------------------------------------------------------
//
// AuthRequired validates access token from cookie.
// RefreshTokenRequired validates refresh token for refresh endpoint.
//
// |--------------------------------------------------------------------------
// | Active Session Middleware
// |--------------------------------------------------------------------------
//
// ActiveSessionRequired checks session_id against auth_sessions table.
// It blocks revoked or missing sessions.
//
// |--------------------------------------------------------------------------
// | CSRF Middleware
// |--------------------------------------------------------------------------
//
// CSRFRequired runs on state-changing methods only.
// Baseline behavior is intentionally weak for CA vulnerability demo.
//
// LoginCSRFRequired validates login CSRF using in-memory synchronizer tokens.
package middleware

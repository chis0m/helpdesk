// Package auth contains security helpers used by login/session flow.
//
// |--------------------------------------------------------------------------
// | Password Hash
// |--------------------------------------------------------------------------
//
// HashPassword and VerifyPassword use Argon2id for password protection.
//
// |--------------------------------------------------------------------------
// | PASETO + Payload
// |--------------------------------------------------------------------------
//
// Token helpers create and verify access/refresh PASETO tokens.
// Payload carries auth claims like sub, role, exp, jti, type, and sess_id.
//
// |--------------------------------------------------------------------------
// | Login CSRF Store (In-Memory)
// |--------------------------------------------------------------------------
//
// LoginCSRFStore issues short-lived one-time tokens for login CSRF checks.
// For production, replace in-memory storage with Redis/shared storage.
package auth

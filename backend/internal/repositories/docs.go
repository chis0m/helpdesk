// Package repositories contains database access logic.
//
// |--------------------------------------------------------------------------
// | User Repository
// |--------------------------------------------------------------------------
//
// UserRepository reads/writes users table records.
//
// |--------------------------------------------------------------------------
// | Auth Session Repository
// |--------------------------------------------------------------------------
//
// AuthSessionRepository manages session row lifecycle:
// create session, read active session, rotate refresh JTI,
// revoke session, and update CSRF token metadata.
package repositories

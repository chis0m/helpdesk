// Package auth contains authentication and password-security helpers.
//
// |--------------------------------------------------------------------------
// | Password Hash
// |--------------------------------------------------------------------------
//
// Password hashing in this package follows Argon2id-based guidance and was
// inspired by:
//   - OWASP Password Storage Cheat Sheet:
//     https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
//   - Go Argon2 package documentation (Argon2id):
//     https://pkg.go.dev/golang.org/x/crypto/argon2#hdr-Argon2id
//   - Alex Edwards blog post on hashing/verifying with Argon2 in Go:
//     https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
//
// |--------------------------------------------------------------------------
// | PASETO + Payload
// |--------------------------------------------------------------------------
//
// PASETO helpers in this package are used to:
//   - generate access and refresh tokens
//   - verify token validity and token type
//
// Payload contains the core auth claims used in this project:
// iss, sub, role, aud, iat, exp, jti, type, and sess_id.
package auth

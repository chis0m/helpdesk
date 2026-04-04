# Auth

### Password Hash

`HashPassword` and `VerifyPassword` use Argon2id for password protection.

### PASETO + Payload

Token helpers create and verify access and refresh PASETO tokens.  
Payload carries auth claims like `sub`, `role`, `exp`, `jti`, `type`, and `sess_id`.

### Public Auth CSRF Store (In-Memory)

`PublicAuthCSRFStore` issues short-lived one-time tokens for auth actions before login (for example login/signup).  
For production, I use Redis/shared storage.

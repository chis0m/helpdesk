# Repositories

### User Repository

`UserRepository` reads and writes users table records.

### Auth Session Repository

`AuthSessionRepository` manages session row lifecycle:
- create session
- read active session
- rotate refresh JTI
- revoke session
- update CSRF token metadata

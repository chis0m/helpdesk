# Container (`internal/container`)

The **dependency injection** root: one struct holds shared singletons (DB handle, token maker, repositories’ consumers). **`routes.Register`** receives only `*Container` so route files never `new` services by hand.

## `Container` fields (exported)

| Field | Role |
|-------|------|
| `DB` | Shared `*gorm.DB` (same connection pool for all repos). |
| `HealthController` | Liveness/health JSON. |
| `AuthController` | Login, signup, refresh, logout, CSRF, sessions, password flows. |
| `UserController` | Admin user list, staff creation, profile by id, role updates (signup is `AuthController`). |
| `InviteController` | Staff invite verify/accept/create. |
| `TicketController` | Tickets + comments HTTP API. |
| `UserService` | User domain operations (also used by invite acceptance). |
| `TicketService` | Ticket + comment domain logic. |
| `TokenMaker` | `auth.MakerInterface` — PASETO create/verify (middleware + auth service). |
| `PublicAuthCSRFStore` | In-memory tokens for **pre-login** routes (login, signup, forgot/reset password, invite accept). |
| `SessionRepo` | Passed into middleware factories in `routes` for session + CSRF validation. |

Internal-only dependencies (not on the struct but built inside `New`): `InviteService`, `AuthService`, invite/password-reset repositories, ticket comment repository, mail notifiers.

## `New(db, cfg, tokenMaker)` build order

1. **Repositories** — `User`, `Invite`, `Ticket`, `TicketComment`, `AuthSession`, `PasswordReset`.
2. **`PublicAuthCSRFStore`** — TTL from `cfg.CSRFTTL()`.
3. **Notifiers** — `LogStaffInviteNotifier`, `LogPasswordResetNotifier` (`internal/mail/log_mailer.go`; CA logs URLs instead of sending email).
4. **Services** — `UserService` → `InviteService` → `TicketService` → `AuthService` (auth pulls in password reset + PASETO + sessions).
5. **Controllers** — each gets the services/stores it needs.

Order matters where constructors take dependencies: e.g. `AuthService` needs `PasswordResetRepository` and `resetNotifier`; `InviteService` needs `inviteRepo`, `userRepo`, notifier.

# Boot (`boot`)

This package wires the HTTP server **from process start to listening on a port**. It is the composition root for the backend: everything else (DB, DI container, routes) is assembled here.

## `NewApp()` sequence

1. **`config.Load()`** — reads environment (via `godotenv`) into `config.Config` (`internal/config/env.go`).
2. **`logger.Init(cfg)`** — configures global zerolog (level, console vs JSON). Must run early so later steps can log.
3. **`database.RunMigrations(dsn)`** — runs Goose **up** against the `migrations/` folder (embedded via `goose.SetBaseFS` in `database` package `init`). Applies before GORM opens a long-lived connection so schema matches code expectations.
4. **`database.Connect(cfg)`** — opens a **GORM** MySQL connection used for the rest of the process lifetime.
5. **`seed.SeedAll(db, cfg)`** — idempotent seed: default super-admin; **CA fixtures** run only when **`SEED_CA`** is true (`true`/`1`/`yes`/`on`). Set `SEED_CA=false` to skip (`seed/seeds/ca/`). Safe to run on every boot.

### CA seed behaviour

- **One customer** (`riley.mustchange@…`) uses a **hardcoded** password (`CaMustChange1!` in `seed/seeds/ca/user.go`); **`must_change_password`** is set. On first creation, the server logs that password at **info** (intentional for CA).
- **Jordan**, both **staff** accounts: **random** password on first creation only, logged at **warn** with a message that this is **intentional for test data** (copy from logs; not printed again if the row already exists).
- Customer emails: `firstname.lastname@<company domain>`. Staff: `firstname.lastname@secweb.ie`.

| Email | Role | Password |
| --- | --- | --- |
| `riley.mustchange@acmelogistics.ie` | user | Hardcoded `CaMustChange1!` (see `ca/user.go`) |
| `jordan.lee@northwind.ie` | user | Random (see startup log when created) |
| `casey.admin@secweb.ie` | admin | Random (see startup log when created) |
| `sam.support@secweb.ie` | staff | Random (see startup log when created) |

**Tickets** (matched by title; skipped if already present):

- `[CA Seed] Unassigned — no comments` — reporter: Riley; unassigned; no comments.
- `[CA Seed] Assigned — 3-comment thread` — reporter: Jordan; assigned to Sam; three comments (user → staff → user).
- `[CA Seed] Unassigned — user comment only` — reporter: Riley; unassigned; one comment from Riley.

Default **super-admin** seed email domain is **`secweb.ie`** (`admin@secweb.ie` unless overridden by `SEED_ADMIN_EMAIL`).
6. **`auth.NewPasetoMaker(key)`** — builds the symmetric PASETO signer/verifier (**32-byte** key from config). Shared by middleware and `AuthService`.
7. **`container.New(db, cfg, tokenMaker)`** — constructs repositories, services, controllers, and the public-auth CSRF store. See `internal/container/container.md`.
8. **`gin.Default()`** — engine with logger + recovery middleware.
9. **`applyCORS(engine, cfg)`** — registers `gin-contrib/cors` using `FRONTEND_URL` (default `http://localhost:5173` if empty), with credentials enabled for cookie-based auth.
10. **`routes.Register(engine, c)`** — mounts all API routes and middleware chains. See `internal/routes/routes.md`.

## `Run()`

Binds to `:{PORT}` from config, logs the address, and blocks on `engine.Run`. On shutdown path, `logger.Sync()` is called (currently a no-op for zerolog but keeps the hook for future sinks).

## Who should change this?

- Adding a **global middleware** (e.g. request ID): usually here right after `gin.Default()` or inside `applyCORS` ordering considerations.
- Changing **startup order** (e.g. migrate after connect): only if you understand Goose + GORM lifecycle implications.

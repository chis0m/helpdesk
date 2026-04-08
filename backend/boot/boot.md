# Boot (`boot`)

This package wires the HTTP server **from process start to listening on a port**. It is the composition root for the backend: everything else (DB, DI container, routes) is assembled here.

## `NewApp()` sequence

1. **`config.Load()`** — reads environment (via `godotenv`) into `config.Config` (`internal/config/env.go`).
2. **`logger.Init(cfg)`** — configures global zerolog (level, console vs JSON). Must run early so later steps can log.
3. **`database.RunMigrations(dsn)`** — runs Goose **up** against the `migrations/` folder (embedded via `goose.SetBaseFS` in `database` package `init`). Applies before GORM opens a long-lived connection so schema matches code expectations.
4. **`database.Connect(cfg)`** — opens a **GORM** MySQL connection used for the rest of the process lifetime.
5. **`seed.SeedAll(db, cfg)`** — idempotent seed: default super-admin; **CA fixtures** run only when **`SEED_CA`** is true (`true`/`1`/`yes`/`on`). Set `SEED_CA=false` to skip (`seed/seeds/ca/`). Safe to run on every boot.

### CA seed behaviour

- All CA fixture users share **`password`** (`caTestPassword` in `seed/seeds/ca/helpers.go`) for easy local testing.
- **Must Change** (`must.change@…`) still has **`must_change_password`** set (forced change flow); others behave as already-set passwords (`PasswordChangedAt` set where applicable).
- Customer emails: `firstname.lastname@<company domain>`. Staff: `firstname.lastname@secweb.ie`.

| Email | Role | Password |
| --- | --- | --- |
| `must.change@example.com` | user | `password` (must change on first login) |
| `john.doe@sample.com` | user | `password` |
| `jane.doe@sample.com` | user | `password` |
| `sam.support@secweb.ie` | staff | `password` |
| `cassey.support@secweb.ie` | staff | `password` |

**Tickets** (matched by title; skipped if already present). The **Must Change** user has **no tickets**. **John Doe** has three; **Jane Doe** has two. Assignees are balanced: **Sam** holds three tickets (John open, John resolved, Jane closed); **Cassey** holds two (John in progress, Jane open).

- `Sidebar keeps collapsing when I switch between projects` — reporter: John Doe; assigned to Sam Support; status open.
- `CSV export times out after about 30 seconds` — reporter: John Doe; assigned to Cassey Support; status in progress.
- `Cannot log in after password reset — now fixed on my side` — reporter: John Doe; assigned to Sam Support; status resolved.
- `Question: when are monthly billing reminder emails sent?` — reporter: Jane Doe; assigned to Cassey Support; status open.
- `Invoice PDF shows wrong VAT rate for last month's bill` — reporter: Jane Doe; assigned to Sam Support; status closed.

Default **super-admin** seed is **`super.admin@secweb.ie`** (firstname **super**, lastname **admin**) unless overridden by `SEED_ADMIN_*` env vars.
6. **`auth.NewPasetoMaker(key)`** — builds the symmetric PASETO signer/verifier (**32-byte** key from config). Shared by middleware and `AuthService`.
7. **`container.New(db, cfg, tokenMaker)`** — constructs repositories, services, controllers, and the public-auth CSRF store. See `internal/container/container.md`.
8. **`gin.Default()`** — engine with logger + recovery middleware.
9. **`applyCORS(engine, cfg)`** — registers `gin-contrib/cors` using `FRONTEND_URL` (default `http://localhost:3000` if empty), with credentials enabled for cookie-based auth.
10. **`routes.Register(engine, c)`** — mounts all API routes and middleware chains. See `internal/routes/routes.md`.

## `Run()`

Binds to `:{PORT}` from config, logs the address, and blocks on `engine.Run`. On shutdown path, `logger.Sync()` is called (currently a no-op for zerolog but keeps the hook for future sinks).

## Who should change this?

- Adding a **global middleware** (e.g. request ID): usually here right after `gin.Default()` or inside `applyCORS` ordering considerations.
- Changing **startup order** (e.g. migrate after connect): only if you understand Goose + GORM lifecycle implications.

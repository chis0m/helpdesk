# Boot (`boot`)

This package wires the HTTP server **from process start to listening on a port**. It is the composition root for the backend: everything else (DB, DI container, routes) is assembled here.

## `NewApp()` sequence

1. **`config.Load()`** — reads environment (via `godotenv`) into `config.Config` (`internal/config/env.go`).
2. **`logger.Init(cfg)`** — configures global zerolog (level, console vs JSON). Must run early so later steps can log.
3. **`database.RunMigrations(dsn)`** — runs Goose **up** against the `migrations/` folder (embedded via `goose.SetBaseFS` in `database` package `init`). Applies before GORM opens a long-lived connection so schema matches code expectations.
4. **`database.Connect(cfg)`** — opens a **GORM** MySQL connection used for the rest of the process lifetime.
5. **`seed.SeedAll(db, cfg)`** — idempotent seed (e.g. default super-admin). Safe to run on every boot (`seed/seeder.go`).
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

# Database (`internal/database`)

Two responsibilities: **schema migrations** (Goose) and the **runtime GORM connection** used by repositories.

## Migrations — `migrate.go`

- **`RunMigrations(dsn)`** — opens a dedicated `database/sql` connection with the MySQL driver, runs `goose.Up` on the **`migrations`** directory, then closes the connection. Called from `boot.NewApp()` before `Connect`.
- **`ResetDb(dsn)`** — `goose.DownTo(..., 0)` then `goose.Up` — useful for local wipe/rebuild (often invoked via Makefile/scripts, not production).

### Goose filesystem

`init()` calls `goose.SetBaseFS(os.DirFS("."))` so migration paths resolve relative to the **process working directory** (typically the module root when you run the server). SQL files live under repo `migrations/`.

### DSN

The same MySQL DSN string as GORM (`config.Config.MySQLDSN()`) is passed in: user, password, host, port, database, `parseTime=true`.

## Connection — `connect.go`

- **`Connect(cfg)`** — `gorm.Open` with MySQL driver, **warn-level** GORM logger only (reduces noise; use application zerolog for request/business logs).

## Design notes

- **Migrations are SQL-first** (Goose), not GORM AutoMigrate — table definitions are explicit and reviewable.
- Repositories depend on `*gorm.DB` injected from the container; they do not open their own connections.

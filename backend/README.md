# Help Desk Backend (Initial Setup)

## Quick start

```bash
cd /Users/chisom/NCI/golang/helpdesk/backend
cp .env.example .env
go install github.com/pressly/goose/v3/cmd/goose@latest
make serve
```

Health check: `GET http://localhost:8080/api/health`

## Project structure

- `cmd/server`: application entrypoint.
- `boot`: app bootstrap (config, db, container, router).
- `internal/container`: dependency wiring (controllers/services/repositories).
- `internal/controllers`: HTTP layer.
- `internal/services`: business logic layer.
- `internal/repositories`: data access layer.
- `internal/database`: DB connection and migration helpers.
- `internal/routes`: route registration.
- `internal/models`: domain models.
- `migrations`: Goose SQL migration files.
- `seed`: database seed logic.

## Tools used and why

- `Gin`: fast, clean HTTP routing and middleware.
- `GORM + MySQL driver`: quick DB integration and ORM mapping.
- `Goose`: versioned SQL migrations and easy DB reset flow.
- `godotenv`: local `.env` loading for configuration.
- `go-playground/validator`: request payload validation.
- `zap`: structured logs for debugging and audit-friendly logs.
- `gin-contrib/cors`: CORS setup between Vue frontend and Go API.
- `paseto`: secure token format for auth implementation.

## Useful commands (Makefile)

```bash
make check
make tidy
make gs
make gu
make gd
make gr
make gm name=create_users_table
```

Set valid DB values (`DB_HOST`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`, `DB_PORT`) before running goose targets.

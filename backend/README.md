# Help Desk Backend (Initial Setup)

## Quick start

```bash
cd /Users/chisom/NCI/golang/helpdesk/backend
cp .env.example .env
go install github.com/pressly/goose/v3/cmd/goose@latest
make serve

Health check: `GET http://localhost:8080/api/health`

## Project structure

Docs for heavier packages: `boot`, `internal/database`, `internal/container`, `internal/auth`, `internal/middleware`, `internal/routes`, `internal/controllers`, `internal/services`, `internal/repositories` (each has a `*.md` next to the code where linked below).

- `cmd/server`: application entrypoint (`main` → `boot.NewApp`).
- `boot`: app bootstrap — [`boot/boot.md`](boot/boot.md).
- `internal/config`: environment loading (`internal/config/env.go`).
- `internal/logger`: zerolog init (`internal/logger/zerolog.go`).
- `internal/database`: DB connection and Goose migrations — [`internal/database/database.md`](internal/database/database.md).
- `internal/container`: dependency wiring — [`internal/container/container.md`](internal/container/container.md).
- `internal/auth`: PASETO, passwords, token types — [`internal/auth/auth.md`](internal/auth/auth.md).
- `internal/middleware`: auth, session, CSRF, rate limits — [`internal/middleware/middleware.md`](internal/middleware/middleware.md).
- `internal/routes`: route registration — [`internal/routes/routes.md`](internal/routes/routes.md).
- `internal/controllers`: HTTP layer — [`internal/controllers/controllers.md`](internal/controllers/controllers.md).
- `internal/services`: business logic — [`internal/services/services.md`](internal/services/services.md).
- `internal/repositories`: data access — [`internal/repositories/repositories.md`](internal/repositories/repositories.md).
- `internal/mail`: notifier interfaces; CA build logs invite/reset URLs instead of SMTP.
- `internal/models`: domain models (structs + GORM tags).
- `internal/requests` / `internal/response`: request DTOs and JSON helpers.
- `migrations`: Goose SQL migration files.
- `seed`: boot-time admin seed (`seed/seeder.go`, `seed/seeds/`).

## Tools used and why

- `Gin`: fast, clean HTTP routing and middleware.
- `GORM + MySQL driver`: quick DB integration and ORM mapping.
- `Goose`: versioned SQL migrations and easy DB reset flow.
- `godotenv`: local `.env` loading for configuration.
- `go-playground/validator`: request payload validation.
- `Argon2id` (`golang.org/x/crypto/argon2`): modern password hashing.
- `zerolog`: structured logging with clean pretty logs in development.
- `gin-contrib/cors`: CORS setup between Vue frontend and Go API.
- `paseto`: secure token format for auth implementation.

## Useful commands (Makefile)

```bash
make check
make tidy
make gs
make gu
make gd
make gr OR make dbr
make gm name=create_users_table
```

Set valid DB values (`DB_HOST`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`, `DB_PORT`) before running goose targets.

## CI (GitHub Actions)

Branches **`vulnerable-baseline`** and **`secure-fix`** run:

- **Backend workflow** (`.github/workflows/backend-ci.yml`): two jobs — **`test`** runs on every PR and push `make gocheck`, `make govuln`. **`push-ecr`** runs only on **`push`** or **`merge`**  to these branches, **after `test` succeeds**. It fetches Secrets Manager JSON, writes the env file, builds, and pushes to **`public.ecr.aws/f8m4k2h6/helpdesk-backend:vuln`** or **`:secure`** as the case may be.
- **Frontend:** `npm ci`, `npm run lint`, `npm run audit`.

Repository secrets: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`. Store each env bundle in Secrets Manager as a **JSON object** (keys → string values) so the console stays readable.
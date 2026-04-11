# Secure Web Helpdesk

A full-stack **ticketing / helpdesk** web application for managing support requests, users, and staff workflows. It is built for the **National College of Ireland — Secure Web Development** module: the **primary security focus** is applying **defence in depth** across authentication, session handling, authorization, input/output handling, CSRF protection, safe database access, and operational audit logging.

**Stack:** Vue 3 + TypeScript (SPA) + Tailwind CSS, Golang (Gin) REST API, MySQL (GORM, Goose), PASETO-based sessions tokens in **HttpOnly** cookies, **Argon2id** password hashing, **zerolog** structured logging.

---

## Features and security objectives

### Major functionality

| Area | What users can do |
|------|-------------------|
| **Accounts** | Sign up, log in, log out, change password, list/revoke sessions, password reset (mail/log driver) |
| **Roles** | Multiple roles (e.g. user(customer), staff, admin, super_admin) with **server-side** checks on sensitive routes |
| **Tickets** | Create, list, search, view, update, assign, change status, delete (within policy) |
| **Comments** | Add, list, update, delete ticket comments |
| **Admin** | Staff invites, user listing, role changes |

### Security objectives (mapped to implementation themes)

| Objective | How the project addresses it |
|-----------|------------------------------|
| **Strong authentication** | Argon2id password hashes; rate limiting on login/signup/forgot-password |
| **Session integrity** | Short-lived access token + long refresh tokens; sessions and CSRF tokens (for update/post/delete requests) stored together in database |
| **Authorization** | Role and ownership checks on tickets and admin actions |
| **Injection & XSS** | Parameterized DB access; validated request DTOs; UI must not treat user text as HTML but as string |
| **CSRF** | `X-CSRF-Token` on protected and public-auth routes as the case may be |
| **Observability** | Structured logs; audit/append-only logging targeted in remediation docs |
| **Configuration** | Environment-driven secrets and DB settings (never commit real `.env` secrets) |

---

## Github Branch Stragety
I intentionally added some features on seperate branch. For example, I added audit logging to the app in vulnerable state, before implementing the secure fix. This is so that I can show how proper Audit Logging can help with easy detection.

Vulnerability Branch: [vulnerable-baseline](https://github.com/chis0m/helpdesk/tree/vulnerable-baseline)
Audit Logging: [add-audit-log-to-track-vuln](https://github.com/chis0m/helpdesk/tree/add-audit-log-to-track-vuln)
Secure Fix Branch: [secure-fix](https://github.com/chis0m/helpdesk/tree/secure-fix)
SAST Test Branch: [sast-test](https://github.com/chis0m/helpdesk/tree/sast-test) 

## Project structure

```
helpdesk/
├── backend/                 # Go API (Gin, GORM, Goose migrations)
│   ├── cmd/server/          # Application entrypoint (`func main()`)
│   ├── boot/                # App bootstrap (HTTP server, Dependency Injection (DI) wiring)
│   ├── internal/
│   │   ├── auth/            # Tokens, CSRF header names, cookie names
│   │   ├── config/          # Environment configuration
│   │   ├── container/       # Dependency injection container
│   │   ├── controllers/     # HTTP handlers (e.g. auth, users, tickets, comments)
│   │   ├── mail/            # Manages mail dispatch
│   │   ├── middleware/      # Auth, CSRF, session checks
│   │   ├── models/          # GORM models
│   │   ├── repositories/    # Database access layer
│   │   ├── services/        # Business logic
│   │   ├── routes/          # Route registration and management
│   │   └── requests/        # Request and validation structs
|   |   ├── response/        # Response struct
│   ├── migrations/          # SQL migrations (using Goose)
│   ├── seed/                # Optional seed data for immediate testing
│   ├── makefile             # Makefile for easy command execution
│   └── go.mod
├── frontend/                # Vue 3 + Vite + TypeScript + Tailwind
│   ├── src/
│   │   ├── api/             # HTTP client modules (auth, tickets, users)
│   │   ├── views/           # View Pages (login, tickets, profile, admin, …)
│   │   ├── components/      # Layout and UI components
│   │   ├── router/          # Vue Router definitions
│   │   ├── stores/          # Auth/session client state
│   │   └── utils/           # Shared helpers
│   ├── vite.config.ts
│   └── package.json
└── README.md                  # This file
```
---

## Setup and installation

### Prerequisites

- **Go** (see `backend/go.mod` for version)
- **Node.js** (LTS recommended) and **npm**
- **MySQL** 8.x (or compatible)
- **Goose** for migrations: `go install github.com/pressly/goose/v3/cmd/goose@latest`

### 1. Database

Create a database and user (example):

```sql
CREATE DATABASE helpdesk CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'helpdesk'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL ON helpdesk.* TO 'helpdesk'@'localhost';
FLUSH PRIVILEGES;
```

### 2. Backend environment

From `backend/`:

```bash
cp .env.example .env
```

Setup variables on you can rely on the defaults at `internal/config/env.go`

Start the API:

```bash
cd backend/ && make serve
```
This also runs migrations, and it runs the default seed if `SEED_CA` is set to true

Default API port is **`8080`** unless `PORT` is set.

### 3. Frontend

```bash
cd frontend
npm install
```

Make sure the `VITE_PORT` and `VITE_API_BASE_URL` matches the backend records.
VITE_PORT is frontend port, which should match `FRONTEND_URL` in the `.env.go`
VITE_API_BASE_URL is the backend url

```bash
npm run dev      # development server
npm run build    # production build
npm run lint     # ESLint
```

Open the URL printed by Vite (typically `http://localhost:3000`).

---

## Usage guidelines

1. **Register / log in** — login through /login in the UI e.g http://127.0.0.1:3000/login. If you set `SEED_CA` to true, then use `cassey.admin@secweb.ie` to login and `password` as password. Or you can just signup.
2. **Tickets** — Create a ticket from the dashboard; open a ticket to view details, status, assignment, and comments.
3. **Search** — Use ticket search; the server must bind search parameters safely (no raw string concatenation into SQL).
4. **Profile & sessions** — Update profile where allowed; review and revoke other sessions from the sessions view.
5. **Admin / staff** — Elevated actions (invites, role changes, staff creation) require an appropriate role; all checks are enforced **on the server**, not only in the UI.

---

## Security improvements (vulnerabilities addressed)

The project follows a structured remediation plan aligned with the course **VULN-01 … VULN-06** themes. Summary:

| ID | Topic | Risk | Remediation summary |
|----|-------|------|---------------------|
| **VULN-01** | Session cookie flags | Tokens readable by JS or sent over HTTP | Set **`HttpOnly=true`**, **`Secure=true`** in production (HTTPS) while local `Secure=false`, consistent **`SameSite=strict`** |
| **VULN-02** | IDOR | Access other users’ profiles or tickets by ID | Prefer **`GET/PATCH /api/users/me`** only; enforce **ticket access policy** such that only owner or admin can access tickets. Replace ID with **UUID** |
| **VULN-03** | Stored XSS | User content executed as HTML | At backend **`TrimSpace`** on text fields; render ticket/comment bodies as **plain text** (no `v-html` (in frontend) on user-controlled fields) |
| **VULN-04** | CSRF | Cross-site form posts mutate state | **`CSRFRequired`**: compare **`X-CSRF-Token`** to session token with **constant-time** equality before `c.Next()` |
| **VULN-05** | Weak audit trail | No forensic trail for incidents | **`audit_logs`** (append-only), actor from **session only**, optional metadata; no secrets in log payloads |
| **VULN-06** | SQL injection (search) | Search query concatenated into raw SQL | **Parameterization** of search query |

Code locations are tagged with `// VULN-…` e.g **//VULN-01** in the vulnerability branch(vulnerable-baseline); then replaced with tags with `// SEC-…` e.g **//SEC-01** in secure-fix brach.

---

## Testing process (WIP)

Testing is required to cover **functional** behaviour and **static analysis (SAST)**.

| Type | Tool / method | What to record |
|------|----------------|----------------|
| **SAST (Go)** | `gosec ./...` from `backend/` | Summary of findings; fix or document accepted risk |
| **SAST / lint (JS)** | `npm run lint` in `frontend/`; optional **Semgrep** or **ESLint** security plugins | Output snippet in report appendix |
| **Dependencies** | `npm audit` (frontend); **`govulncheck ./...`** (backend) | High/critical issues and upgrades planned |
| **Functional security** | Manual or scripted tests | Cookie flags, 403 on IDOR, 403 on bad CSRF, safe search with metacharacters |

Example commands:

```bash
# Backend — install gosec: go install github.com/securego/gosec/v2/cmd/gosec@latest
cd backend && gosec ./...

# Go vulnerability check
cd backend && govulncheck ./...

# Frontend
cd frontend && npm audit
cd frontend && npm run lint
```

Summarize **key findings** (pass/fail, severity, file paths) in your report and README as they change over time.

---

## Contributions and references

- **Author / course:** Developed as part of **NCI MSCCYB — Secure Web Development** (Continuous Assessment). **Contributors:** see Git commit history; this repository is the canonical source for the helpdesk project.
- **No upstream “forked app”** — application code is authored for the module; **third-party libraries** are used under their respective **open-source licences** (see `go.mod` and `package.json`).

**Frameworks and libraries (non-exhaustive)**

| Component | Reference |
|-----------|-----------|
| HTTP router | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| Migrations | [Goose](https://github.com/pressly/goose) |
| Frontend | [Vue.js](https://vuejs.org/), [Vite](https://vitejs.dev/), [Vue Router](https://router.vuejs.org/) |
| Password hashing | [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto/argon2) (Argon2) |
|Argon2d Implementation|[How to Hash and Verify Passwords with Argon2 in Go](https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go)
| Tokens | [PASETO](https://github.com/o1egl/paseto)
|Rate Limiting| [Rate Limiting 101: Implementing in Go](https://medium.com/@sheikhahnafshifat/rate-limiting-101-implementing-in-go-c434675f1fbe), [time rate](https://pkg.go.dev/golang.org/x/time/rate)

**Security guidance (for design and report citations)**

- [Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [OWASP Top 10](https://owasp.org/Top10/)
- [OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/)
- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)

---

## Repository hygiene

- Prefer **small, meaningful commits** over a single dump before deadlines.
- Keep **`.env` out of git**; use secrets only via environment or a secrets manager in production.
- Run **`go fmt`**, **`go test ./...`**, **`npm run lint`**, and **`gosec`** before tagging a release or submitting assessment work.

---

*Last updated to align with NCI CA README expectations: overview, structure, setup, usage, security improvements, testing, and references.*

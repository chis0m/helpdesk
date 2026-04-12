# Secure Web Helpdesk

A full-stack **ticketing / helpdesk** web application for managing support requests, users, and staff workflows. It is built for the **National College of Ireland — Secure Web Development** module: the **primary security focus** is applying **defence in depth** across authentication, session handling, authorization, input/output handling, CSRF protection, safe database access, and operational logging.

**Stack:** Vue 3 + TypeScript (SPA) · Go (Gin) REST API · MySQL · PASETO-based sessions in **HttpOnly** cookies (with explicit cookie-flag hardening in remediation) · **Argon2id** password hashing · **zerolog** structured logging.

---

## Features and security objectives

### Major functionality

| Area | What users can do |
|------|-------------------|
| **Accounts** | Sign up, log in, log out, change password, list/revoke sessions, password reset (mail/log driver) |
| **Roles** | Multiple roles (e.g. customer, staff, admin / super_admin) with **server-side** checks on sensitive routes |
| **Tickets** | Create, list, search, view, update, assign, change status, delete (within policy) |
| **Comments** | Add, list, update, delete ticket comments |
| **Admin** | Staff invites, user listing, role changes (where implemented) |

### Security objectives (mapped to implementation themes)

| Objective | How the project addresses it |
|-----------|------------------------------|
| **Strong authentication** | Argon2id password hashes; rate limiting on login/signup/forgot-password |
| **Session integrity** | Short-lived access token + refresh; sessions stored server-side; CSRF tokens for mutating requests |
| **Authorization** | Role and ownership checks on tickets and admin actions (see **Security improvements**) |
| **Injection & XSS** | Parameterized DB access; validated request DTOs; UI must not treat user text as HTML |
| **CSRF** | `X-CSRF-Token` on protected and public-auth routes where middleware applies |
| **Observability** | Structured logs; audit/append-only logging targeted in remediation docs |
| **Configuration** | Environment-driven secrets and DB settings (never commit real `.env` secrets) |

---

## Project structure

```
helpdesk/
├── backend/                 # Go API (Gin, GORM, Goose migrations)
│   ├── cmd/server/          # Application entrypoint (`main`)
│   ├── boot/                # App bootstrap (HTTP server, DI wiring)
│   ├── internal/
│   │   ├── auth/            # Tokens, CSRF header names, cookie names
│   │   ├── config/          # Environment configuration
│   │   ├── container/       # Dependency injection container
│   │   ├── controllers/     # HTTP handlers (e.g. auth, users, tickets)
│   │   ├── middleware/      # Auth, CSRF, rate limits, session checks
│   │   ├── models/          # GORM models
│   │   ├── repositories/     # Database access layer
│   │   ├── services/        # Business logic
│   │   ├── routes/          # Route registration
│   │   └── requests/        # Request binding / validation structs
│   ├── migrations/          # SQL migrations (Goose)
│   ├── seed/                # Optional seed data (incl. CA fixtures when `SEED_CA` is set)
│   ├── makefile             # `goose` DB tasks, `serve`, `check`
│   └── go.mod
├── frontend/                # Vue 3 + Vite + TypeScript + Tailwind
│   ├── src/
│   │   ├── api/             # HTTP client modules (auth, tickets, users)
│   │   ├── views/           # Pages (login, tickets, profile, admin, …)
│   │   ├── components/      # Layout and UI components
│   │   ├── router/          # Vue Router definitions
│   │   ├── stores/          # Auth/session client state
│   │   └── utils/           # Shared helpers
│   ├── vite.config.ts
│   └── package.json
├── sast/                    # Saved SAST / dependency scan JSON (CA appendix)
└── README.md                # This file
```

**Important files**

| Path | Purpose |
|------|---------|
| `backend/internal/routes/routes.go` | All API routes and middleware chains |
| `backend/internal/middleware/csrf_middleware.go` | CSRF verification for session-backed mutations |
| `backend/internal/controllers/auth_cookie.go` | Sets/clears auth cookies (flags are a key hardening point) |
| `backend/internal/repositories/ticket_repository.go` | Ticket persistence and search (must stay injection-safe) |
| `frontend/src/api/client.ts` | Fetch-based helpers, `credentials: 'include'`, CSRF header name |

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
cp .env.example .env   # if your repo provides .env.example; otherwise create .env
```

Set at minimum (defaults exist in `internal/config/env.go` but **override for production**):

| Variable | Purpose |
|----------|---------|
| `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD` | MySQL connection |
| `PASETO_SYMMETRIC_KEY` | **32-byte** secret for PASETO (generate a random value; never use the dev default in production) |
| `FRONTEND_URL` | SPA origin for CORS (e.g. `http://localhost:3000`) |
| `GO_ENV` | `development` vs `production` (affects logging and future cookie security toggles) |

Run migrations:

```bash
cd backend
make gu    # goose up — requires goose on PATH and DB_* env / makefile defaults
```

Start the API:

```bash
make serve
# or: go run ./cmd/server
```

Default API port is **`8080`** unless `PORT` is set.

### 3. Frontend

```bash
cd frontend
npm install
```

Optional: `frontend/.env` — `VITE_PORT` (default **3000** per `vite.config.ts`), `VITE_API_URL` if your client reads it.

```bash
npm run dev      # development server
npm run build    # production build
npm run lint     # ESLint
```

Open the URL printed by Vite (typically `http://localhost:3000`). Ensure `FRONTEND_URL` in the backend matches this origin so cookies and CORS behave correctly.

---

## Usage guidelines

1. **Register / log in** — Use the auth flows from the landing/login views. After login, the browser stores **HttpOnly** session cookies once remediation flags are enabled (see security section).
2. **Tickets** — Create a ticket from the dashboard; open a ticket to view details, status, assignment, and comments.
3. **Search** — Use ticket search; the server must bind search parameters safely (no raw string concatenation into SQL).
4. **Profile & sessions** — Update profile where allowed; review and revoke other sessions from the sessions view.
5. **Admin / staff** — Elevated actions (invites, role changes, staff creation) require an appropriate role; all checks are enforced **on the server**, not only in the UI.

**Note:** If you run the API on port `8080` and the SPA on port `3000`, that is a **cross-origin** setup: the frontend must call the API with **credentials included**, and the backend CORS configuration must allow the frontend origin. Adjust `FRONTEND_URL` and CORS settings if you change ports.

---

## Security improvements (vulnerabilities addressed)

The project follows a structured remediation plan aligned with the course **VULN-01 … VULN-06** themes. Summary:

| ID | Topic | Risk | Remediation summary |
|----|-------|------|---------------------|
| **VULN-01** | Session cookie flags | Tokens readable by JS or sent over HTTP | Set **`HttpOnly=true`**, **`Secure=true`** in production (HTTPS), consistent **`SameSite`**, and mirror flags in **`clearAuthCookies`** |
| **VULN-02** | IDOR | Access other users’ profiles or tickets by ID | Prefer **`GET/PATCH /api/users/me`** only; enforce **ticket access policy** (reporter / assignee / admin) on every ticket route |
| **VULN-03** | Stored XSS | User content executed as HTML | **`TrimSpace`** on text fields; render ticket/comment bodies as **plain text** (no `v-html` on user-controlled fields) |
| **VULN-04** | CSRF | Cross-site form posts mutate state | **`CSRFRequired`**: compare **`X-CSRF-Token`** to session token with **constant-time** equality before `c.Next()` |
| **VULN-05** | Weak audit trail | No forensic trail for incidents | **`audit_logs`** (append-only), actor from **session only**, optional metadata; no secrets in log payloads |
| **VULN-06** | SQL injection (search) | Search query concatenated into raw SQL | **Parameterized** `LIKE` (or GORM `Where` with `?` placeholders); escape `%` / `_` for `LIKE` if needed |

Code locations are tagged with `// VULN-…` comments until fully remediated; replace tags with `// SEC-…` or remove after review. Internal design notes may live in your assignment `ca2/vuln-mitigation.md` (if copied into the repo or kept alongside the project).

---

## Testing process

Security validation combines **Go SAST** (**gosec**), **Go dependency** analysis (**govulncheck**), **npm** dependency audit, **lint**, and **manual functional** checks (sessions, CSRF, authorization, safe search). Reports are saved under **`sast/`**.

### Install CLI tools (once)

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
```

Ensure `$(go env GOPATH)/bin` is on your **`PATH`**.

### Backend (`cd backend`)

```bash
cd backend
gosec -fmt json -out ../sast/gosec-report-vuln.json ./...
govulncheck -json ./... > ../sast/backend-govulncheck.json
```

After fixes, gosec was rerun to a clean report, for example:

```bash
gosec -fmt json -out ../sast/gosec-report-vuln-fixed.json ./...
```

**`govulncheck`** reported no known reachable vulnerabilities against the scanned **`go.mod`** in the last run (**`sast/backend-govulncheck.json`**).

### Frontend (`cd frontend`)

```bash
cd frontend
npm audit --json > ../sast/frontend-npm-audit-vuln.json
npm audit fix
npm audit --json > ../sast/frontend-npm-audit-vuln-fix.json
npm run lint
# To fix lint issues
npm run lint -- --fix
```

### Key findings (summary)

| Tool | Findings & fix |
|------|----------------|
| **gosec** | **`sast/gosec-report-vuln.json`**: **G115** (`int`→`uint32` in Argon2 helpers), **G203** (`template.HTML` in email HTML). To **Fix:** create help to convert int to `uint32` conversion; email layout split so plain text is template-escaped. **`sast/gosec-report-vuln-fixed.json`**: **0** issues. |
| **govulncheck** | No vulnerable **Go** dependencies reported, check **`sast/backend-govulncheck.json`**. |
| **npm audit** | Vulnerabilities listed in **`sast/frontend-npm-audit-vuln.json`** (e.g. **TypeScript-ESLint** chain); addressed with **`npm audit fix`**, captured in **`sast/frontend-npm-audit-vuln-fix.json`**. |
| **Functional** | Manual checks: cookies, 403 on forbidden access, CSRF rejection, safe search — detailed in the CA report. |

**Supplementary checks (optionally):** **Semgrep** and **Snyk** were also run for a general pattern and dependency view across frontend/backend, the output JSON is in **`sast/`**.

---

## Contributions and references

- **Author / course:** Developed as part of **NCI MSCCYB — Secure Web Development** (Continuous Assessment). **Contributors:** see Git commit history; this repository is the canonical source for the helpdesk project.

**Frameworks and libraries (non-exhaustive)**

| Component | Reference |
|-----------|-----------|
| HTTP router | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| Migrations | [Goose](https://github.com/pressly/goose) |
| Frontend | [Vue.js](https://vuejs.org/), [Vite](https://vitejs.dev/), [Vue Router](https://router.vuejs.org/) |
| Password hashing | [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) (Argon2) |
| Tokens | [PASETO](https://github.com/o1egl/paseto) |
|MySQL Databae|

**Security guidance (for design and report citations)**

- [OWASP Top 10](https://owasp.org/Top10/)
- [OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/)
- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)

---

## Repository hygiene

- Prefer **small, meaningful commits** over larg commits
- Keep **`.env` out of git**
- Run **`go fmt`**, **`go test ./...`**, **`gosec`**, **`govulncheck`**, **`npm audit`**, and **`npm run lint`** before submitting assessment work.

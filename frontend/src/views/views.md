# Views (`src/views`)

**Route-level** Vue pages: they match **routes** in `src/router/index.ts`, load data, handle forms, and compose layout + components.

## Public / auth

- **`LandingView.vue`** — Marketing / entry.
- **`LoginView.vue`**, **`SignupView.vue`** — Auth; login sets session + starts refresh scheduler.
- **`views/auth/`** — `AcceptInviteView`, `ForgotPasswordView`, `ResetPasswordView`, **`ChangePasswordRequiredView`** (TASK.md §8.3 — shown when `must_change_password` is true until the user updates their password).

## Dashboard (inside `AppLayout`)

- **`DashboardView.vue`** — Home after login.
- **`views/tickets/`** — List, create, detail (comments, status, assign).
- **`views/account/`** — Profile, sessions.
- **`views/admin/`** — Admin users, staff, create staff.

## Patterns

- **API calls** — `@/api/...` (auth, tickets, users, admin).
- **CSRF** — `getSessionCsrfToken()` from `@/stores/auth-session` before mutating calls.
- **Navigation** — `vue-router` (`useRouter`, `useRoute`); shared path strings in `src/constants/routes.ts`.

## Layouts

Views are often **children** of `layouts/AuthLayout.vue` (auth) or `layouts/AppLayout.vue` (dashboard shell with sidebar/header).

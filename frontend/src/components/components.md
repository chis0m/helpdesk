# Components (`src/components`)

Vue **SFCs** (`.vue`) grouped by feature. They are **presentational** or **small interactive pieces**; data usually comes from parent views or API calls made in views.

## Layout (`components/layout/`)

- **`AppHeader.vue`** — Top bar; can trigger logout via `performLogout`.
- **`AppSidebar.vue`** — Dashboard nav + logout.
- **`SidebarLink.vue`** — Nav link with active state.

## Feature areas

- **`components/tickets/`** — Ticket status badge, status editor (reuse between list/detail).
- **`components/admin/`** — Admin tables (`AdminUserTable`, `AdminStaffTable`), `AdminSubnav` for admin sub-routes.
- **`components/icons/`** — Small SVG icons (`IconHome`, `IconTicket`, etc.).

## How they fit

- **Layouts** (`src/layouts/`) wrap whole pages (`AuthLayout` for auth screens, `AppLayout` for dashboard).
- **Views** (`src/views/`) compose layouts + components and own **routing** and **data loading** for a screen.

## `HelloWorld.vue`

Default Vite starter; not part of the helpdesk flow. Safe to remove if unused.

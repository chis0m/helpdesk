# Utils (`src/utils`)

Small **plain** helpers (no Vue components).

| File | Purpose |
|------|--------|
| `logger.ts` | `logger.debug` (dev / `VITE_DEBUG`), `warn`, `error`. API modules use `debug` for HTTP traces. |
| `date-format.ts` | Format ISO timestamps for display. |
| `perform-logout.ts` | Clears refresh timer, calls `POST /api/auth/logout` when CSRF exists, clears session + navigates to login. |
| `ticket-ui.ts` | Ticket UI helpers (labels, colors, etc.). |
| `user-display.ts` | Format user names for display. |

Import from `@/utils/...` in views or components.

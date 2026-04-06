/**
 * Central route paths — keep in sync with `src/router/index.ts`.
 * Use these in RouterLink, router.push, and docs instead of string literals.
 *
 * VULN-02: `profile(userId)` builds paths with numeric user id (vulnerable baseline; UUID-only on remediation branch).
 */
export const paths = {
  landing: '/',
  login: '/login',
  signup: '/signup',
  dashboard: {
    home: '/dashboard',
    tickets: '/dashboard/tickets',
    ticketNew: '/dashboard/tickets/new',
    ticketDetail: (id: string) => `/dashboard/tickets/${encodeURIComponent(id)}`,
    /** Numeric user id (vulnerable baseline); remediation may switch to UUID-only paths. */
    profile: (userId: string | number) =>
      `/dashboard/profile/${encodeURIComponent(String(userId))}`,
    /** Use with router redirect so `/dashboard/profile` resolves to the signed-in user. */
    profileRedirect: '/dashboard/profile',
    adminUsers: '/dashboard/admin/users',
    /** Internal staff directory (includes administrator flag per row). */
    adminStaff: '/dashboard/admin/staff',
    adminStaffNew: '/dashboard/admin/staff/new',
  },
} as const

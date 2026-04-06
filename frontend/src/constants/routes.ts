/**
 * Central route paths — keep in sync with `src/router/index.ts`.
 * Use these in RouterLink, router.push, and docs instead of string literals.
 */
export const paths = {
  landing: '/',
  login: '/login',
  signup: '/signup',
  dashboard: {
    home: '/dashboard',
    tickets: '/dashboard/tickets',
    ticketDetail: (id: string) => `/dashboard/tickets/${encodeURIComponent(id)}`,
    profile: '/dashboard/profile',
    adminUsers: '/dashboard/admin/users',
    /** Internal staff directory (includes administrator flag per row). */
    adminStaff: '/dashboard/admin/staff',
    adminStaffNew: '/dashboard/admin/staff/new',
  },
} as const

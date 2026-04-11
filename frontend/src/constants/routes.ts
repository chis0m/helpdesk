/**
 * Central route paths — keep in sync with `src/router/index.ts`.
 * Use these in RouterLink, router.push, and docs instead of string literals.
 */
export const paths = {
  landing: '/',
  login: '/login',
  /** Required when `must_change_password` is true (TASK.md §8.3). */
  changePasswordRequired: '/change-password',
  signup: '/signup',
  acceptInvite: '/accept-invite',
  forgotPassword: '/forgot-password',
  resetPassword: '/reset-password',
  dashboard: {
    home: '/dashboard',
    tickets: '/dashboard/tickets',
    ticketNew: '/dashboard/tickets/new',
    ticketDetail: (ticketUuid: string) =>
      `/dashboard/tickets/${encodeURIComponent(ticketUuid)}`,
    profile: '/dashboard/profile',
    adminUsers: '/dashboard/admin/users',
    /** Internal staff directory (includes administrator flag per row). */
    adminStaff: '/dashboard/admin/staff',
    adminStaffNew: '/dashboard/admin/staff/new',
    sessions: '/dashboard/sessions',
  },
} as const

/** Portal customers (admin Users list). */
export type PortalUser = {
  id: string
  email: string
  displayName: string
  createdAt: string
  organization?: string
}

/** Internal directory roles (aligned with `GET /api/admin/users`). */
export type StaffDirectoryRole = 'user' | 'staff' | 'admin' | 'super_admin'

/** Internal staff; `isAdmin` is true for admin + super_admin; `role` is the API role. */
export type StaffMember = {
  id: string
  email: string
  displayName: string
  createdAt: string
  isAdmin: boolean
  role: StaffDirectoryRole
}

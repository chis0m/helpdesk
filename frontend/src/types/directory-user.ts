/** Portal customers (admin Users list). */
export type PortalUser = {
  id: string
  email: string
  displayName: string
  createdAt: string
  organization?: string
}

/** Internal staff; `isAdmin` is the administrator flag (elevated role). */
export type StaffMember = {
  id: string
  email: string
  displayName: string
  createdAt: string
  isAdmin: boolean
}

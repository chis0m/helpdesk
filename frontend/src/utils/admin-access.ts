/** Matches backend admin directory / staff management (`admin` or `super_admin` only — not `staff`). */
export function isAdminPortalRole(role: string | undefined | null): boolean {
  return role === 'admin' || role === 'super_admin'
}

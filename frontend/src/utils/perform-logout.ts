// VULN-01: Logout calls API to clear session cookies issued under weak cookie flags (server-side VULN-01).
import type { Router } from 'vue-router'
import { logoutRequest } from '@/api/auth'
import { paths } from '@/constants/routes'
import { clearAuthSession, getSessionCsrfToken } from '@/stores/auth-session'

/**
 * Calls the API to revoke the session and clear cookies when possible, clears SPA session
 * storage, then navigates to login. Still clears local state if the API returns an error.
 */
export async function performLogout(router: Router): Promise<void> {
  const csrf = getSessionCsrfToken()
  let redirectPath: string = paths.login
  if (csrf) {
    const result = await logoutRequest(csrf)
    if (result.ok && result.data.redirect_to.startsWith('/'))
      redirectPath = result.data.redirect_to
  }

  clearAuthSession()
  await router.replace(redirectPath)
}

/**
 * Human-readable copy when the API rejects auth — avoids raw "authentication required" after session drift.
 */
export function friendlyHttpError(status: number, message: string): string {
  if (status === 401)
    return 'Your session has ended. Please sign in again.'
  if (status === 403) {
    const m = message.toLowerCase()
    if (m.includes('csrf') || m.includes('session is not active') || m.includes('invalid session'))
      return 'Your session could not be verified. Please sign in again.'
  }
  return message
}

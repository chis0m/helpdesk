/**
 * Small console logger: verbose traces in dev (and optional VITE_DEBUG), warn/error always.
 *
 * `logger.debug` is stripped from production builds unless VITE_DEBUG=true. API modules log
 * full envelopes and CSRF values in dev for vulnerable-baseline / coursework visibility;
 * redact those debug calls when hardening for production.
 */

function verbose(): boolean {
  if (import.meta.env.DEV)
    return true
  return import.meta.env.VITE_DEBUG === 'true'
}

function line(context: string, message: string): string {
  return `[${context}] ${message}`
}

export const logger = {
  /** Dev / VITE_DEBUG only — use for HTTP traces and UI debugging. */
  debug(context: string, message: string, ...args: unknown[]): void {
    if (!verbose())
      return
    console.log(line(context, message), ...args)
  },

  warn(message: string, ...args: unknown[]): void {
    console.warn(message, ...args)
  },

  error(message: string, ...args: unknown[]): void {
    console.error(message, ...args)
  },
}

import type { ApiTicketRow } from '@/api/tickets'
import type { TicketStatus } from '@/types/ticket'

/** Short label derived from `ticket_uuid` for chips only (no backend serial). */
export function ticketUuidDisplayRef(ticketUuid: string): string {
  return ticketUuid.replace(/-/g, '').slice(0, 8).toUpperCase()
}

/** Workflow steps for ticket status (UI copy + value). */
export const TICKET_STATUS_DEFINITIONS: {
  value: TicketStatus
  title: string
  description: string
}[] = [
  {
    value: 'open',
    title: 'Open',
    description: 'Waiting in the queue—no one has started work yet.',
  },
  {
    value: 'in_progress',
    title: 'In progress',
    description: 'Actively being worked by support.',
  },
  {
    value: 'resolved',
    title: 'Resolved',
    description: 'Fix or answer is in place; reopen if something is still wrong.',
  },
  {
    value: 'closed',
    title: 'Closed',
    description: 'Fully done—no further work planned on this ticket.',
  },
]

/** Tailwind classes for status pills (list + detail). */
export function ticketStatusBadgeClass(status: TicketStatus): string {
  if (status === 'in_progress') return 'bg-amber-100 text-amber-900'
  if (status === 'resolved' || status === 'closed') return 'bg-emerald-100 text-emerald-900'
  return 'bg-neutral-200 text-neutral-800'
}

export function formatTicketStatusLabel(status: TicketStatus): string {
  return status.replace(/_/g, ' ')
}

/**
 * Stored `category` is a short slug (often `general`, the DB default). Maps known
 * values to readable labels; otherwise title-cases the raw value for display.
 */
export function formatTicketCategoryLabel(category: string): string {
  const key = category.trim().toLowerCase()
  const known: Record<string, string> = {
    general: 'General support',
    technical: 'Technical',
    billing: 'Billing',
    account: 'Account',
    security: 'Security',
  }
  if (known[key])
    return known[key]
  const raw = category.trim()
  if (!raw)
    return 'Uncategorized'
  return raw.charAt(0).toUpperCase() + raw.slice(1)
}

/** Reporter line for lists: "You", display name, or `User #id` fallback. */
export function reporterDisplayLabel(t: ApiTicketRow, currentUserId: number | null): string {
  if (currentUserId !== null && t.reporter_user_id === currentUserId)
    return 'You'
  const name = t.reporter_display_name?.trim()
  if (name)
    return name
  return `User #${t.reporter_user_id}`
}

/** Assignee for detail UI, or null if unassigned. */
export function assigneeDisplayLabel(t: ApiTicketRow): string | null {
  if (t.assigned_user_id == null)
    return null
  const name = t.assigned_display_name?.trim()
  if (name)
    return name
  return `User #${t.assigned_user_id}`
}

/** Tailwind classes for category pills (list + detail). */
export function categoryBadgeClass(category: string): string {
  const k = category.trim().toLowerCase()
  if (k === 'general')
    return 'bg-sky-50 text-sky-900 ring-sky-200/90'
  if (k === 'technical')
    return 'bg-violet-50 text-violet-900 ring-violet-200/90'
  if (k === 'billing')
    return 'bg-amber-50 text-amber-950 ring-amber-200/90'
  if (k === 'account')
    return 'bg-emerald-50 text-emerald-900 ring-emerald-200/90'
  if (k === 'security')
    return 'bg-rose-50 text-rose-900 ring-rose-200/90'
  return 'bg-neutral-100 text-neutral-800 ring-neutral-200/90'
}

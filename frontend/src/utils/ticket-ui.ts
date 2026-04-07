import type { ApiTicketRow } from '@/api/tickets'
import type { TicketStatus } from '@/types/ticket'

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

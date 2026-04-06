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

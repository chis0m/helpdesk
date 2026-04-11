/** Ticket domain types for the helpdesk UI. */

export type TicketStatus = 'open' | 'in_progress' | 'resolved' | 'closed'

export type TicketComment = {
  id: string
  authorName: string
  body: string
  createdAt: string
  isStaff: boolean
}

export type Ticket = {
  /** Opaque ticket UUID (routing). */
  id: string
  title: string
  description: string
  category: string
  status: TicketStatus
  createdAt: string
  updatedAt: string
  assigneeName: string | null
  reporterName: string
  /** Portal user email (for staff detail). */
  reporterEmail: string | null
  /** Assigned support email (for portal user detail). */
  assigneeEmail: string | null
  comments: TicketComment[]
}

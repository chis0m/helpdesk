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
  id: string
  title: string
  description: string
  category: string
  status: TicketStatus
  createdAt: string
  updatedAt: string
  assigneeName: string | null
  reporterName: string
  comments: TicketComment[]
}

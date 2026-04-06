/** Sample ticket data for the helpdesk UI. */
import type { Ticket } from '@/types/ticket'
import { mergeTicketWithStatusOverrides } from '@/stores/ticket-status-overrides'
import { formatTicketStatusLabel } from '@/utils/ticket-ui'

export const MOCK_TICKETS: Ticket[] = [
  {
    id: '1042',
    title: 'VPN drops on Wi‑Fi handoff',
    description:
      'When moving between office access points, the SecWeb desktop client disconnects for 30–60 seconds. Reproducible on macOS 14 and Windows 11 with client 3.2.1.',
    category: 'Technical',
    status: 'in_progress',
    createdAt: '2026-02-10T09:15:00Z',
    updatedAt: '2026-02-12T14:22:00Z',
    assigneeName: 'Sam Rivera',
    reporterName: 'Jordan Lee',
    comments: [
      {
        id: 'c1',
        authorName: 'Sam Rivera',
        body: 'Thanks for the detail. Can you attach client logs from Settings → Diagnostics after the next drop?',
        createdAt: '2026-02-11T10:00:00Z',
        isStaff: true,
      },
      {
        id: 'c2',
        authorName: 'Jordan Lee',
        body: 'Uploaded logs — file name vpn-handoff-logs.zip on the case.',
        createdAt: '2026-02-11T16:40:00Z',
        isStaff: false,
      },
    ],
  },
  {
    id: '1038',
    title: 'Invoice PDF not attaching in billing export',
    description:
      'Export to CSV works but “Include PDF summary” leaves the PDF column empty. Org ID 8821.',
    category: 'Billing',
    status: 'open',
    createdAt: '2026-02-08T11:00:00Z',
    updatedAt: '2026-02-08T11:00:00Z',
    assigneeName: null,
    reporterName: 'Jordan Lee',
    comments: [
      {
        id: 'c3',
        authorName: 'Taylor Kim',
        body: 'We’re able to reproduce on 2.4.0. Tracking under REL-9912.',
        createdAt: '2026-02-09T09:12:00Z',
        isStaff: true,
      },
    ],
  },
  {
    id: '1031',
    title: 'SSO redirect loop on Safari',
    description:
      'After IdP login, Safari bounces between /callback and /dashboard until cache cleared. Chrome OK.',
    category: 'Technical',
    status: 'open',
    createdAt: '2026-02-05T08:30:00Z',
    updatedAt: '2026-02-06T13:00:00Z',
    assigneeName: 'Sam Rivera',
    reporterName: 'Alex Ng',
    comments: [],
  },
  {
    id: '1041',
    title: 'Expense tool CSV export blank rows',
    description: 'First row is headers only; data rows are empty for March 2026 close.',
    category: 'Billing',
    status: 'resolved',
    createdAt: '2026-01-28T15:20:00Z',
    updatedAt: '2026-02-03T17:00:00Z',
    assigneeName: 'Taylor Kim',
    reporterName: 'Jordan Lee',
    comments: [
      {
        id: 'c4',
        authorName: 'Taylor Kim',
        body: 'Fixed in 2.4.1 — please retry export and confirm.',
        createdAt: '2026-02-03T17:00:00Z',
        isStaff: true,
      },
    ],
  },
]

export function getMockTicketById(id: string): Ticket | undefined {
  const base = MOCK_TICKETS.find((t) => t.id === id)
  return base ? mergeTicketWithStatusOverrides(base) : undefined
}

/** All sample tickets with session status overrides applied. */
export function allMockTicketsWithOverrides(): Ticket[] {
  return MOCK_TICKETS.map((t) => mergeTicketWithStatusOverrides(t))
}

/** Client-side search (merged status). */
export function searchMockTickets(query: string): Ticket[] {
  const q = query.trim().toLowerCase()
  const merged = allMockTicketsWithOverrides()
  if (!q)
    return merged
  return merged.filter(
    (t) =>
      t.title.toLowerCase().includes(q)
      || t.description.toLowerCase().includes(q)
      || t.category.toLowerCase().includes(q)
      || t.id.includes(q)
      || formatTicketStatusLabel(t.status).includes(q),
  )
}

export function countOpenMockTickets(): number {
  return allMockTicketsWithOverrides().filter(
    (t) => t.status === 'open' || t.status === 'in_progress',
  ).length
}

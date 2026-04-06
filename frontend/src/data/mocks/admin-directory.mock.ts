/** Sample users and staff for admin directory screens. */
import type { PortalUser, StaffMember } from '@/types/directory-user'

export const MOCK_PORTAL_USERS: PortalUser[] = [
  {
    id: 'usr-701',
    email: 'jordan.lee@customer.example',
    displayName: 'Jordan Lee',
    createdAt: '2025-08-12T10:00:00Z',
    organization: 'Acme Logistics',
  },
  {
    id: 'usr-702',
    email: 'alex.ng@partner.co',
    displayName: 'Alex Ng',
    createdAt: '2025-11-03T14:20:00Z',
    organization: 'Partner Co',
  },
  {
    id: 'usr-703',
    email: 'morgan.smith@retail.example',
    displayName: 'Morgan Smith',
    createdAt: '2026-01-15T09:45:00Z',
    organization: 'Retail One',
  },
]

export const MOCK_STAFF_MEMBERS: StaffMember[] = [
  {
    id: 'stf-201',
    email: 'sam.rivera@secweb.internal',
    displayName: 'Sam Rivera',
    createdAt: '2024-03-01T08:00:00Z',
    isAdmin: false,
  },
  {
    id: 'stf-202',
    email: 'taylor.kim@secweb.internal',
    displayName: 'Taylor Kim',
    createdAt: '2024-06-18T11:30:00Z',
    isAdmin: false,
  },
  {
    id: 'stf-203',
    email: 'riley.chen@secweb.internal',
    displayName: 'Riley Chen',
    createdAt: '2025-02-10T16:00:00Z',
    isAdmin: false,
  },
  {
    id: 'adm-101',
    email: 'ops.lead@secweb.internal',
    displayName: 'Casey Ortiz',
    createdAt: '2023-01-09T12:00:00Z',
    isAdmin: true,
  },
  {
    id: 'adm-102',
    email: 'security.owner@secweb.internal',
    displayName: 'Jamie Patel',
    createdAt: '2023-04-22T09:15:00Z',
    isAdmin: true,
  },
]

import { ref, type Ref } from 'vue'
import type { TicketStatus } from '@/types/ticket'

const STORAGE_KEY = 'secweb-helpdesk-ticket-status'

type OverrideEntry = { status: TicketStatus; updatedAt: string }
type OverrideMap = Record<string, OverrideEntry>

function readStorage(): OverrideMap {
  if (typeof sessionStorage === 'undefined')
    return {}
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY)
    if (!raw)
      return {}
    const parsed = JSON.parse(raw) as unknown
    return parsed !== null && typeof parsed === 'object' ? (parsed as OverrideMap) : {}
  }
  catch {
    return {}
  }
}

function writeStorage(map: OverrideMap) {
  if (typeof sessionStorage === 'undefined')
    return
  try {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(map))
  }
  catch {
    /* ignore quota / private mode */
  }
}

/**
 * Session-persisted status overrides per ticket id.
 * Import the ref (or call setters) so list + detail stay aligned.
 */
export const ticketStatusOverrides: Ref<OverrideMap> = ref(readStorage())

export function setTicketStatusOverride(ticketId: string, status: TicketStatus): void {
  const next: OverrideMap = {
    ...ticketStatusOverrides.value,
    [ticketId]: { status, updatedAt: new Date().toISOString() },
  }
  ticketStatusOverrides.value = next
  writeStorage(next)
}

export function mergeTicketWithStatusOverrides<T extends { id: string; status: TicketStatus; updatedAt: string }>(
  ticket: T,
): T {
  const o = ticketStatusOverrides.value[ticket.id]
  if (!o)
    return ticket
  return { ...ticket, status: o.status, updatedAt: o.updatedAt }
}

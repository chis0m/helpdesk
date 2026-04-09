# Types (`src/types`)

**TypeScript types** used by the Vue app (not API envelope types — those live in `src/api/types.ts`).

## `ticket.ts`

Domain shapes for the **UI**: `TicketStatus`, `TicketComment`, `Ticket` (fields like `createdAt` match how views map API data).

## `directory-user.ts`

Types for **admin directory** / staff lists (`PortalUser`, `StaffMember`, etc.).

## Usage

- Views and components import from `@/types/...` (e.g. `import type { Ticket } from '@/types/ticket'`).
- API envelopes and response shapes live with API code (`src/api/types.ts` and modules); this folder holds UI-side types that differ from raw JSON.

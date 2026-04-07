# Types (`src/types`)

**TypeScript types** used by the Vue app (not API envelope types — those live in `src/api/types.ts`).

## `ticket.ts`

Domain shapes for the **UI**: `TicketStatus`, `TicketComment`, `Ticket` (fields like `createdAt` match how views map API data).

## `directory-user.ts`

Types for **admin directory** / staff lists (`PortalUser`, `StaffMember`, etc.).

## Usage

- `import type { ... } from '@/types/ticket'` in views or components.
- Keep API response shapes next to the API modules (`api/*.ts`); keep **view models** here when they differ from raw JSON.

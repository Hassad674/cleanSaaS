# Admin Panel — Vite + React + TypeScript + Tailwind

## Purpose

Separate admin dashboard for managing the SaaS application. Not part of the Next.js frontend — this is an independent Vite app.

## Structure

```
admin/src/
├── components/     → Shared UI (Sidebar, Header, DataTable, Layout)
├── pages/          → Dashboard, Users, Blog, Settings
├── hooks/          → useAuth, useApi
├── lib/            → API client (talks to Go backend)
└── types/          → TypeScript types
```

## Conventions

### Design Tokens
- Same CSS custom properties as the main frontend (copied from globals.css)
- Use semantic Tailwind classes: `bg-card`, `text-foreground`, `border-border`, etc.
- NEVER use hardcoded Tailwind colors (no `zinc-`, `gray-`, `slate-`, `white`, `black`)

### API Communication
- All API calls go through `lib/api.ts`
- Base URL from `VITE_API_URL` env var (default: `http://localhost:8081`)
- Auth token stored in localStorage, sent as `Authorization: Bearer <token>` header
- All admin endpoints require admin role — backend enforces this

### Routing
- `react-router-dom` v7 with BrowserRouter
- Protected routes redirect to `/login` if not authenticated
- Routes: `/login`, `/` (dashboard), `/users`, `/blog`, `/settings`

### State Management
- React Context for auth state
- Local component state for page-specific data
- No external state library needed

## Running

```bash
cd admin
npm run dev    # Starts on port 5174
npm run build  # Production build (tsc + vite)
```

## Testing

> Current state: the admin app has **no automated tests** yet — this is the least-mature surface of the boilerplate. Treat new admin work as test-debt-paying-down, not test-debt-adding.

- **Unit/component**: add **Vitest** + `@testing-library/react` (same toolchain as the Next.js frontend). Co-locate `*.test.tsx` next to the component/page. Cover the `lib/api.ts` client (mock fetch), auth context, and protected-route redirects.
- **Type check** (must pass before commit): `npx tsc --noEmit`.
- **End-to-end**: admin flows are exercised via the frontend's Playwright setup where they overlap; admin-specific e2e can be added under a dedicated spec once flows stabilize.
- Run the validation pipeline before committing admin changes: `npm run build` + `npx tsc --noEmit` (+ Vitest once tests exist).

## Environment Variables

- `VITE_API_URL` — Backend API URL (default: `http://localhost:8081`)

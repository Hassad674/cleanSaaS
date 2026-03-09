# Frontend — Next.js + Tailwind + Feature-based Architecture

## Structure

```
frontend/src/
├── app/                  → ROUTING ONLY (thin layer)
│   ├── (marketing)/      → Public pages: landing, pricing, blog, docs
│   ├── (auth)/           → Auth pages: login, register, forgot-password
│   ├── (dashboard)/      → Authenticated app: dashboard, settings, ai, files
│   └── (admin)/          → Admin-only pages
│
├── features/             → BUSINESS LOGIC (self-contained modules)
│   ├── auth/             → Login, register, OAuth, session
│   ├── billing/          → Pricing, subscription, invoices
│   ├── user/             → Profile, settings
│   ├── ai/               → Chat, conversation history
│   ├── notification/     → Bell, notification list
│   ├── storage/          → File upload, file list
│   ├── admin/            → Stats, user management
│   ├── marketing/        → Hero, features section, CTA
│   └── blog/             → Post list, post content
│
├── shared/               → CROSS-FEATURE reusable code
│   ├── components/       → UI components (shadcn), layouts
│   ├── hooks/            → Generic hooks (useApi, useDebounce)
│   ├── lib/              → API client, utils, constants
│   └── types/            → Shared TypeScript types
│
├── config/               → Site config (name, nav links)
└── styles/               → globals.css
```

## Rules

### app/ — Routing only
- Each `page.tsx` is 5-20 lines max
- Imports from `features/`, composes components, nothing else
- Layouts handle visual structure and route guards
- Server Components by default
- Metadata (title, description) defined here

```tsx
// CORRECT: thin page
import { LoginForm } from "@/features/auth/components/login-form";
export default function LoginPage() {
  return <LoginForm />;
}

// WRONG: logic in page
export default function LoginPage() {
  const [email, setEmail] = useState(""); // NO
}
```

### features/ — Self-contained modules
- Each feature has: `components/`, `actions/`, `hooks/`, `types.ts`
- A feature NEVER imports from another feature
- If two features need the same thing → it goes in `shared/`
- `actions/` contains Server Actions that call the Go backend
- `components/` contains UI specific to this feature
- `hooks/` contains client-side state for this feature

```
features/auth/
├── components/
│   ├── login-form.tsx
│   └── register-form.tsx
├── actions/
│   └── auth.ts            ← Server Actions
├── hooks/
│   └── use-auth.ts        ← Client state
├── lib/
│   └── session.ts         ← Feature-specific utils
└── types.ts               ← Feature types
```

### shared/ — Cross-feature only
- Must be used by 2+ features to belong here
- `components/ui/` → shadcn primitives (button, input, dialog)
- `components/layouts/` → marketing, auth, dashboard, admin layouts
- `lib/api.ts` → single API client for the Go backend
- `lib/utils.ts` → `cn()` helper and generic utilities
- `hooks/` → generic hooks not tied to any feature

### Inter-feature communication
Features never import each other. Composition happens in `app/`:

```tsx
// CORRECT: app/ composes features
import { ProfileHeader } from "@/features/user/components/profile-header";
import { SubscriptionStatus } from "@/features/billing/components/subscription-status";

export default function SettingsPage() {
  return (
    <>
      <ProfileHeader />
      <SubscriptionStatus />
    </>
  );
}

// WRONG: feature imports another feature
// In features/billing/components/status.tsx:
import { useUser } from "@/features/user/hooks/use-user"; // FORBIDDEN
```

## Server vs Client Components

- **Server Component** (default): pages, layouts, data display, SEO content
- **Client Component** (`"use client"`): forms, modals, dropdowns, real-time, onClick/onChange

Rule: if it doesn't need `useState`, `useEffect`, or event handlers, it's a Server Component.

## Data flow

### Server-side (preferred for initial data)
```
page.tsx (Server Component)
  → Server Action (features/X/actions/)
    → shared/lib/api.ts (fetch to Go backend)
```

### Client-side (for interactivity)
```
Client Component
  → Hook (features/X/hooks/)
    → shared/lib/api.ts or Server Action
```

## Code conventions

- **TypeScript strict**: no `any`, no `as` casts unless absolutely necessary
- **File naming**: kebab-case (`login-form.tsx`, `use-auth.ts`)
- **Component naming**: PascalCase (`LoginForm`, `DashboardLayout`)
- **Imports**: always use `@/` alias (maps to `src/`)
- **Styling**: Tailwind only, use `cn()` for conditional classes
- **No barrel exports**: import from specific files, not index.ts

## How to add a new feature

Example: "Add a team management feature"

1. Create `features/team/types.ts` — define Team, Member types
2. Create `features/team/actions/team.ts` — Server Actions (createTeam, inviteMember)
3. Create `features/team/hooks/use-team.ts` — client state if needed
4. Create `features/team/components/` — TeamList, InviteForm, MemberRow
5. Create `app/(dashboard)/settings/team/page.tsx` — imports from features/team/
6. Nothing else touched.

## Testing

- **Components**: `@testing-library/react` — test render + interactions
- **Server Actions**: unit tests mocking the API client
- **Hooks**: `renderHook` from Testing Library
- **E2E**: Playwright for critical flows

Test file next to source: `login-form.tsx` → `login-form.test.tsx`

## Commands

```bash
npm run dev       # Start dev server
npm run build     # Production build
npm run lint      # ESLint
npm run test      # Run tests
```

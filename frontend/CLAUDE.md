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

## Design system

### Identity
Clean, soft, warm. Inspired by Airbnb's aesthetic — rose/pink primary, warm neutrals, generous spacing, smooth transitions. Mobile-first, responsive across all breakpoints.

### Color palette (CSS variables in `globals.css`)

| Token | Light | Dark | Usage |
|-------|-------|------|-------|
| `primary` | Rose pink | Rose pink | Buttons, links, accents, focus rings |
| `primary-foreground` | White | White | Text on primary backgrounds |
| `background` | White | Near-black | Page background |
| `foreground` | Near-black | Off-white | Body text |
| `card` | White | Dark gray | Card backgrounds |
| `muted` | Light warm gray | Dark gray | Subtle backgrounds, disabled states |
| `muted-foreground` | Medium gray | Medium gray | Secondary text, placeholders |
| `accent` | Light rose tint | Dark rose tint | Hover states, active sidebar items |
| `border` | Light gray | Dark gray | Borders, separators |
| `destructive` | Red | Red | Errors, delete actions |
| `success` | Green | Green | Success states |
| `warning` | Amber | Amber | Warning states |

### How to use colors in components

```tsx
// CORRECT: use semantic tokens
<button className="bg-primary text-primary-foreground rounded-lg">
<p className="text-muted-foreground">
<div className="border border-border bg-card rounded-xl">
<span className="text-destructive">

// WRONG: hardcoded Tailwind colors
<button className="bg-zinc-900 text-white">
<p className="text-zinc-500">
<div className="border border-zinc-200 bg-white">
```

### Typography
- Font: **Geist Sans** (display + body) — loaded via `next/font/google`
- Mono: **Geist Mono** — for code blocks, technical content
- Sizes: Tailwind defaults (`text-sm`, `text-base`, `text-lg`, `text-2xl`, etc.)
- Font weight: `font-medium` for labels/nav, `font-bold` for headings, normal for body

### Spacing & layout
- Padding: `p-4` (mobile) → `p-6` (tablet) → `p-8` (desktop) via responsive prefixes
- Container: `container mx-auto px-4 sm:px-6 lg:px-8`
- Cards: `bg-card border border-border rounded-xl p-6 shadow-sm`
- Gaps: `gap-3` (tight), `gap-4` (default), `gap-6` (sections), `gap-8` (page sections)

### Border radius
- Small elements (badges, chips): `rounded-md` (--radius-sm)
- Default (inputs, buttons): `rounded-lg` (--radius-lg)
- Cards, panels: `rounded-xl` (--radius-xl)

### Interactive states
- Buttons: `hover:opacity-90 transition-opacity disabled:opacity-50`
- Links: `text-primary hover:underline` or `text-muted-foreground hover:text-foreground transition-colors`
- Focus: `focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring`
- Active sidebar: `bg-sidebar-accent text-sidebar-accent-foreground font-medium`

### Responsive breakpoints (mobile-first)
- `sm:` → 640px (large phones)
- `md:` → 768px (tablets)
- `lg:` → 1024px (desktop — sidebar appears)
- `xl:` → 1280px (wide desktop)

Always write mobile styles first, then add responsive overrides:
```tsx
// CORRECT: mobile-first
<main className="p-4 sm:p-6 lg:p-8">
<aside className="hidden lg:block w-64">

// WRONG: desktop-first
<aside className="w-64 md:hidden">
```

### Dark mode
- Automatic via `prefers-color-scheme: dark` in CSS
- All tokens have dark variants defined in `globals.css`
- Never use `dark:` Tailwind prefix — the CSS variables handle it
- Test both modes when building components

### shadcn/ui compatibility
The CSS variables follow the shadcn/ui convention. When installing shadcn components:
- They will pick up the design tokens automatically
- No manual theming needed
- Place components in `shared/components/ui/`

## Code conventions

- **TypeScript strict**: no `any`, no `as` casts unless absolutely necessary
- **File naming**: kebab-case (`login-form.tsx`, `use-auth.ts`)
- **Component naming**: PascalCase (`LoginForm`, `DashboardLayout`)
- **Imports**: always use `@/` alias (maps to `src/`)
- **Styling**: Tailwind with design tokens only, use `cn()` for conditional classes
- **No barrel exports**: import from specific files, not index.ts
- **No hardcoded colors**: always use semantic tokens (`text-primary`, `bg-muted`, etc.)

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

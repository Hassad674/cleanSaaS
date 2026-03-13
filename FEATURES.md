# CleanSaaS — Feature Roadmap

> This file tracks all planned features, their status, and priority.
> Updated: 2026-03-13

---

## Comparison with OpenSaaS (13K+ stars)

| Feature | OpenSaaS | CleanSaaS | Winner |
|---------|----------|-----------|--------|
| **Stack** | React + Node.js + Prisma (Wasp framework) | Next.js + Go + pure SQL (no framework lock-in) | **CleanSaaS** — Go is 10-50x faster, no ORM overhead, no framework dependency |
| **Architecture** | Monolithic Wasp (framework-locked) | Hexagonal Go + Feature-based Next.js | **CleanSaaS** — swappable adapters, independent modules, professional-grade |
| **Auth providers** | Email, Google, GitHub, Slack, Microsoft | Email (+ Google, GitHub in Phase 2) | OpenSaaS for now |
| **Payments** | Stripe, Lemon Squeezy, Polar.sh | Stripe (+ others trivial to add via adapters) | OpenSaaS (more providers) |
| **AI integration** | OpenAI (function calling) | Gemini (multi-provider interface) | Tie |
| **Admin dashboard** | Embedded in main app | Separate Vite app + PostHog analytics | **CleanSaaS** — more scalable, real analytics |
| **Analytics** | Plausible + Google Analytics | PostHog (product analytics, funnels, replays) | **CleanSaaS** — PostHog is far more powerful |
| **Email** | SendGrid, Mailgun, SMTP | Resend (modern, better DX) | Tie (different choices) |
| **File storage** | AWS S3 | Cloudflare R2 | **CleanSaaS** — R2 = zero egress fees, cheaper |
| **Blog** | Astro Starlight (static) | DB-backed CMS with admin editor, SEO fields, tags | **CleanSaaS** — dynamic, manageable from admin |
| **Documentation** | Astro Starlight | Planned (Phase 2) | OpenSaaS for now |
| **Background jobs** | Yes (Wasp built-in) | Yes (Go native goroutines) | Tie |
| **E2E tests** | Playwright | Playwright | Tie |
| **Deployment** | Railway, Fly.io, one-cmd | Vercel + Railway + Neon | Tie |
| **Notifications in-app** | No | Yes (bell + dropdown) | **CleanSaaS** |
| **Referral system** | No | Yes (built-in) | **CleanSaaS** |
| **Mobile app** | No | React Native (shared backend) | **CleanSaaS** |
| **i18n / Translation** | No | Yes (next-intl) | **CleanSaaS** |
| **Real-time WebSocket** | No | Yes (Go native) | **CleanSaaS** |
| **Rate limiting** | No | Yes (Go middleware) | **CleanSaaS** |
| **AI-agent optimized** | AGENTS.md + Cursor rules | CLAUDE.md + skills + memory system | **CleanSaaS** — deeper integration |
| **Dark mode** | Via theme toggle | CSS auto + toggle | Tie |
| **Performance** | Node.js + Prisma ORM | Go + raw SQL + connection pooling | **CleanSaaS** — orders of magnitude faster |
| **Modularity** | Monolithic, coupled | Every feature independently removable | **CleanSaaS** — future CLI `create-cleansaas` ready |

**Bottom line:** OpenSaaS wins on quantity of integrations (5 OAuth providers, 3 payment providers). CleanSaaS wins on architecture quality, performance, modularity, and unique features (notifications, referral, mobile, i18n, real-time, analytics depth).

---

## Phase 1 — Today's Session (Core Features)

These are the features we are building today. Each one gets a commit when complete.

### 1.1 — Settings Page
- **What**: User profile editing (name, avatar), password change, account deletion
- **Backend**: Already has user update endpoint, add password change + delete
- **Frontend**: `features/user/components/settings-*.tsx`, compose in `app/(dashboard)/settings/page.tsx`
- **API keys needed**: None
- **Status**: 🔲 To do

### 1.2 — Resend Email Adapter
- **What**: Transactional email sending (verification, password reset, notifications)
- **Backend**: `adapter/resend/` implements `port/service/EmailService`
- **Templates**: Welcome, email verification, password reset, notification digest
- **API keys needed**: `RESEND_API_KEY` ✅ Present
- **Status**: 🔲 To do

### 1.3 — Forgot Password / Reset Password
- **What**: Full flow — request reset via email → receive link with token → set new password
- **Backend**: Token generation, email sending (Resend), token validation, password update
- **Frontend**: `features/auth/components/forgot-password-form.tsx`, `reset-password-form.tsx`
- **Depends on**: 1.2 (Resend adapter)
- **API keys needed**: `RESEND_API_KEY` ✅
- **Status**: 🔲 To do

### 1.4 — Email Verification
- **What**: After registration, send verification email. Unverified users have limited access.
- **Backend**: Verification token generation + validation endpoint
- **Frontend**: Verification page, resend verification button
- **Depends on**: 1.2 (Resend adapter)
- **API keys needed**: `RESEND_API_KEY` ✅
- **Status**: 🔲 To do

### 1.5 — Stripe Billing
- **What**: Plans (free/pro/enterprise), checkout session, subscription management, webhook handling, invoices
- **Backend**: Migration (plans, subscriptions, invoices tables) + `adapter/stripe/` + billing service + handlers + webhook endpoint
- **Frontend**: Pricing page (public), billing settings (manage subscription, invoices), upgrade/downgrade
- **API keys needed**: `STRIPE_SECRET_KEY` ✅, `STRIPE_WEBHOOK_SECRET` ✅
- **Status**: 🔲 To do (domain entities + ports already scaffolded)

### 1.6 — Cloudflare R2 File Storage
- **What**: Upload files (images, videos, documents), list files, delete files, signed URLs
- **Backend**: Migration (files table) + `adapter/r2/` using S3-compatible API + storage service + handlers
- **Frontend**: Upload component (drag & drop), file list with preview, file management
- **API keys needed**: `R2_ACCOUNT_ID` ✅, `R2_ACCESS_KEY` ✅, `R2_SECRET_KEY` ✅, `R2_BUCKET_NAME` ✅, `R2_PUBLIC_URL` ✅
- **Status**: 🔲 To do (domain entity + port already scaffolded)

### 1.7 — Gemini AI Chat
- **What**: Chat interface with AI, conversation history, streaming responses
- **Backend**: Migration (conversations, messages tables) + `adapter/gemini/` + AI service + SSE streaming handler
- **Frontend**: Chat UI with message bubbles, conversation sidebar, streaming text display
- **API keys needed**: `GEMINI_API_KEY` ✅ Present
- **Status**: 🔲 To do (domain entities + port already scaffolded)

### 1.8 — In-App Notifications
- **What**: Notification bell with unread count, dropdown list, mark as read, notification types (system, billing, etc.)
- **Backend**: Migration (notifications table) + postgres adapter + notification service + handlers
- **Frontend**: Bell icon in header, notification dropdown, notification page
- **API keys needed**: None (in-app only, email notifications via Resend)
- **Status**: 🔲 To do (domain entities + port already scaffolded)

### 1.9 — Blog (DB-backed CMS via Admin Panel)
- **What**: Full blog system managed from the admin panel — not static markdown
- **Backend**: Migration (posts table: title, slug, content, excerpt, cover_image, tags, meta_title, meta_description, status, published_at) + blog service + CRUD endpoints
- **Frontend (public)**: Blog listing page with cards, individual post page (SEO-optimized), tag filtering
- **Frontend (admin)**: Blog editor — rich text, cover image upload (R2), SEO fields (meta title, meta description, slug), tags, status (draft/published/scheduled)
- **SEO per post**: Auto-generated og:image, JSON-LD Article schema, canonical URL
- **API keys needed**: None (uses R2 for images, already configured)
- **Status**: 🔲 To do

### 1.10 — Admin Panel (Separate Vite App)
- **What**: Standalone React app at `/admin/` — user management, analytics dashboard, blog CMS
- **Backend**: Admin endpoints (list users with pagination/search, update user role/status, dashboard stats, blog CRUD) + PostHog API proxy
- **Frontend (admin/)**: New Vite + React + Tailwind app with its own `CLAUDE.md`
- **Sections**: Dashboard (stats + PostHog analytics), Users (list/search/suspend/ban), Blog (create/edit/publish posts), Settings
- **API keys needed**: `POSTHOG_API_KEY` ✅, `POSTHOG_PROJECT_ID` ✅, `POSTHOG_HOST` ✅
- **Status**: 🔲 To do

### 1.11 — SEO
- **What**: Dynamic meta tags, `sitemap.xml`, `robots.txt`, Open Graph images, JSON-LD structured data
- **Backend**: Sitemap generation endpoint (optional)
- **Frontend**: Next.js metadata API in all pages, shared SEO component
- **API keys needed**: None
- **Status**: 🔲 To do

### 1.12 — Background Jobs & Cron
- **What**: Go-native job runner using goroutines + tickers. No external dependency.
- **Jobs**: Clean expired tokens (hourly), send notification digests (daily), sync Stripe subscription status (6h), cleanup orphan files (daily)
- **Backend**: `pkg/jobs/` — scheduler + job registry, wired in `main.go`
- **API keys needed**: None
- **Status**: 🔲 To do

### 1.13 — Backend Middleware (Rate Limiting, Health, Logging, Graceful Shutdown)
- **What**: Production-grade middleware stack
- **Rate limiting**: Token bucket per IP, configurable per endpoint group
- **Health check**: `GET /health` returns DB status, uptime, version
- **Structured logging**: `slog` with JSON output, request ID, duration, status
- **Graceful shutdown**: Handle SIGTERM/SIGINT, drain connections, close DB pool
- **API keys needed**: None
- **Status**: 🔲 To do

### 1.14 — Unit Tests
- **What**: Tests for all domain entities + all app services (mocked ports)
- **Backend**: `domain/*/entity_test.go`, `app/*/service_test.go`
- **Target**: 80%+ coverage on domain + app layers
- **API keys needed**: None
- **Status**: 🔲 To do

### 1.15 — Playwright E2E Tests (Setup + Core Flows)
- **What**: Browser-based end-to-end tests that validate real user flows
- **Setup**: Playwright config, test helpers (login, seed), CI-ready
- **Tests**: Register → login → dashboard, settings update, file upload, AI chat send, billing checkout flow
- **Why in Phase 1**: Acts as a safety net during autonomous development — validates features work end-to-end
- **API keys needed**: None
- **Status**: 🔲 To do

### 1.16 — Landing Page Update
- **What**: Update landing page to showcase all implemented features
- **Structure**:
  1. **Hero** — "Ship your SaaS in weeks, not months" + 2 CTAs (already done)
  2. **Feature Grid** — 8-12 cards with icon + title + one-liner (Auth, Stripe, AI Chat, Storage, Email, Notifications, Blog CMS, Admin Panel, Background Jobs, etc.)
  3. **Spotlight: AI Chat** — Screenshot/mockup + "Built-in AI chat with Gemini, streaming responses, conversation history. Swap provider in 1 line."
  4. **Spotlight: Admin Panel** — Screenshot/mockup + "Real analytics with PostHog, user management, blog CMS. Separate app, scales independently."
  5. **Spotlight: Architecture** — Hexagonal diagram + "Every feature independently removable. Swap any provider. Future CLI-ready."
  6. **Stack** — Technologies used (already done, update if needed)
  7. **Comparison** — Mini table vs OpenSaaS (3-4 key rows)
  8. **DX** — Developer experience (already done)
  9. **CTA** — "Get started in 5 minutes" (already done)
- **API keys needed**: None
- **Status**: 🔲 To do (done last, after all features are built)

---

## Phase 2 — Recommended (Post-Session)

These features add significant value and should be done soon after Phase 1.

### 2.1 — Google OAuth
- **What**: "Sign in with Google" button, account linking with existing email accounts
- **Backend**: `adapter/google/` implements `port/service/AuthProvider`
- **API keys needed**: `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` (to obtain from Google Cloud Console)
- **Priority**: High

### 2.2 — GitHub OAuth
- **What**: "Sign in with GitHub" — important for a dev-oriented boilerplate
- **Backend**: `adapter/github/` implements `port/service/AuthProvider`
- **API keys needed**: `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET` (to obtain from GitHub Developer Settings)
- **Priority**: High

### 2.3 — Referral / Parrainage System
- **What**: Built-in referral system — unique referral codes, tracking, rewards (account credit or extended trial)
- **Backend**: Migration (referrals table) + referral service + endpoints
- **Frontend**: Referral dashboard, share link, referral stats
- **API keys needed**: None (homemade system)
- **Priority**: High (growth/marketing lever)

### 2.4 — i18n / Internationalization
- **What**: Multi-language support using `next-intl`
- **Frontend**: Translation files (en, fr minimum), locale detection, language switcher
- **Backend**: Error messages remain in English (codes), frontend translates
- **Priority**: Medium

### 2.5 — Real-Time WebSocket
- **What**: Go-native WebSocket server (gorilla/websocket or nhooyr/websocket), no external service
- **Use cases**: Live notifications, real-time chat, presence indicators
- **Backend**: `pkg/ws/` — hub + client + broadcast pattern
- **Frontend**: WebSocket hook, auto-reconnect
- **Priority**: Medium

### 2.6 — Documentation Site
- **What**: User-facing docs explaining how to use the boilerplate, install, configure, add features
- **Tech**: Likely Nextra (Next.js based) or Astro Starlight
- **Priority**: Medium

### 2.7 — Dark Mode Toggle
- **What**: Manual toggle (light/dark/system) in addition to current auto `prefers-color-scheme`
- **Frontend**: Theme provider, toggle button in header, persist preference in localStorage
- **Priority**: Low

---

## Phase 3 — Ideas & Extras

Additional features that would differentiate CleanSaaS further. Not committed to yet.

### 3.1 — Mobile App (React Native)
- **What**: Cross-platform mobile app sharing the same Go backend
- **Tech**: React Native (Expo) — closest to the Next.js frontend skillset
- **Shared**: Same API, same auth (JWT), same features
- **Scope**: Auth, dashboard, AI chat, notifications, file upload

### 3.2 — CLI Tool (`create-cleansaas`)
- **What**: Interactive CLI that lets users pick only the modules they want
- **Tech**: Node.js CLI (like create-next-app), reads feature manifests, scaffolds project
- **Example**: `npx create-cleansaas my-app --features auth,billing,ai,storage`

### 3.3 — Sentry Error Tracking
- **What**: Error monitoring for both frontend (Next.js) and backend (Go)
- **Backend**: `adapter/sentry/` middleware
- **Frontend**: Sentry Next.js plugin

### 3.4 — Additional Payment Providers
- **What**: Lemon Squeezy, Polar.sh, Paddle adapters — demonstrate hexagonal swappability
- **Backend**: Each is a new adapter implementing the same `PaymentService` interface
- **Marketing**: "Switch payment provider in 1 line of code"

### 3.5 — Additional AI Providers
- **What**: OpenAI, Claude adapters alongside Gemini
- **Backend**: Each implements `AIService` interface
- **Marketing**: "Swap AI provider in 1 line"

### 3.6 — Teams / Organizations
- **What**: Multi-tenant support — create teams, invite members, role-based access within teams
- **Scope**: Large feature, needs careful design

### 3.7 — Audit Log
- **What**: Track all significant actions (login, settings change, payment, file upload)
- **Backend**: Audit log table + middleware that auto-logs actions
- **Use case**: Compliance, debugging, admin visibility

### 3.8 — Two-Factor Authentication (2FA)
- **What**: TOTP-based 2FA (Google Authenticator, Authy)
- **Backend**: TOTP secret generation, QR code, verification middleware

### 3.9 — Onboarding Flow
- **What**: Guided first-time user experience after registration
- **Frontend**: Multi-step wizard (choose plan, set up profile, connect integrations)

### 3.10 — API Rate Limiting Dashboard
- **What**: Visual dashboard showing API usage, rate limit status, quotas per user/plan
- **Ties into**: Billing (different plans = different rate limits)

### 3.11 — Webhooks System (Outgoing)
- **What**: Let users of the SaaS register webhook URLs and receive events
- **Backend**: Webhook registry, event dispatching, retry logic, signature verification

### 3.12 — Export / Import Data
- **What**: GDPR-compliant data export (JSON/CSV), account data portability
- **Backend**: Export endpoint that bundles all user data

### 3.13 — Changelog / What's New
- **What**: In-app changelog that shows new features to users
- **Frontend**: Modal or sidebar with version history, "new" badges

---

## API Keys Status

| Service | Env Variable | Status | Needed For |
|---------|-------------|--------|------------|
| PostgreSQL | `DATABASE_URL` | ✅ Present | Core |
| JWT | `JWT_SECRET` | ✅ Present | Core (auth) |
| Stripe | `STRIPE_SECRET_KEY` | ✅ Present (test key) | Billing (1.5) |
| Stripe | `STRIPE_WEBHOOK_SECRET` | ✅ Present | Billing webhooks (1.5) |
| Resend | `RESEND_API_KEY` | ✅ Present | Email (1.2, 1.3, 1.4) |
| Cloudflare R2 | `R2_ACCOUNT_ID` | ✅ Present | Storage (1.6) |
| Cloudflare R2 | `R2_ACCESS_KEY` | ✅ Present | Storage (1.6) |
| Cloudflare R2 | `R2_SECRET_KEY` | ✅ Present | Storage (1.6) |
| Cloudflare R2 | `R2_BUCKET_NAME` | ✅ Present | Storage (1.6) |
| Cloudflare R2 | `R2_PUBLIC_URL` | ✅ Present | Storage (1.6) |
| Gemini | `GEMINI_API_KEY` | ✅ Present | AI Chat (1.7) |
| PostHog | `POSTHOG_API_KEY` | ✅ Present | Admin analytics (1.10) |
| PostHog | `POSTHOG_PROJECT_ID` | ✅ Present | Admin analytics (1.10) |
| PostHog | `POSTHOG_HOST` | ✅ Present | Admin analytics (1.10) |
| Google OAuth | `GOOGLE_CLIENT_ID` | ❌ Not yet | Google OAuth (2.1) |
| Google OAuth | `GOOGLE_CLIENT_SECRET` | ❌ Not yet | Google OAuth (2.1) |
| GitHub OAuth | `GITHUB_CLIENT_ID` | ❌ Not yet | GitHub OAuth (2.2) |
| GitHub OAuth | `GITHUB_CLIENT_SECRET` | ❌ Not yet | GitHub OAuth (2.2) |

**Verdict: All API keys for Phase 1 are present. No blockers.**

> Note: `POSTHOG_API_KEY` and `POSTHOG_PROJECT_ID` need to be added to `backend/internal/config/config.go` during admin panel implementation.

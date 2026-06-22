# Architecture Details

Reflects the actual tree (verify with `find backend/internal -type d`). Dependency rule:
`handler → app → domain ← port ← adapter`.

## Backend hexagonal layers

### domain/ — pure business, stdlib-only
- `user/` (entity, role) · `billing/` (subscription, plan, invoice) · `team/` · `notification/`
- `ai/` (conversation, model) · `blog/` · `storage/` (file) · `referral/`
- `errors.go` — typed domain errors (`ErrNotFound`, `ErrValidation`, `ErrUnauthorized`, …)

### port/ — interfaces (small, specific)
- `repository/` — one per aggregate: User, Subscription, Plan, Invoice, Team, Notification, Conversation, Blog, File, Referral
- `service/` — external seams: `PaymentService`, `EmailService`, `AIService`, `StorageService`, `Broadcaster`

### app/ — use cases (orchestration, constructor-injected ports)
- `auth` · `user` · `billing` · `team` · `notification` · `ai` · `blog` · `storage` · `referral` · `admin`

### adapter/ — concrete implementations (the ONLY layer that imports external SDKs)
- `postgres/` — DB connection + all repositories (pure SQL, `database/sql` + `lib/pq`)
- `stripe/` — implements `PaymentService`
- `resend/` — implements `EmailService`
- `gemini/` — implements `AIService` (swap point: any LLM provider implements the same interface)
- `r2/` — implements `StorageService` (Cloudflare R2 via aws-sdk-go-v2 S3 API)

> Note: AI uses **Gemini** today; the boilerplate intends LLM providers to be hot-swappable behind `AIService`
> (a future `adapter/claude` or `adapter/openai` would be one new file + one line in `main.go`).

### handler/ — HTTP transport (Chi)
- `router.go` — all routes (mounted at root: `/auth/...`, `/billing/...`, `/health`, `/ws`, …)
- `auth.go`, `user.go`, `billing.go`, … — thin handlers, delegate to app services
- `middleware/` — auth (JWT), CORS, rate limit, request logging, recoverer, security headers
- `dto/request/` + `dto/response/` — request/response structs with json tags + mappers

### cmd/ — entry points
- `api/main.go` — the composition root; ALL dependency injection + feature wiring here
- `migrate/main.go` — apply/rollback SQL migrations · `seed/main.go` — seed admin user + plans + blog posts

### pkg/ — public, self-contained utilities (no internal/ imports)
- `jwt` (Generate/Validate, HS256) · `hash` (bcrypt) · `validate` (Email/Slug) · `pagination`
- `jobs` (in-process scheduler) · `ws` (WebSocket hub)

## Frontend (feature-based, Next.js App Router)
- `app/` — routes only (thin pages; server components by default)
- `src/features/<f>/` — self-contained: `components/`, `actions/` (server actions), `api/`, `hooks/`
- `src/shared/` — UI components, hooks, lib, types · `src/config/` — app config
- Rule: features never import each other; pages in `app/` compose across features.

## Request flow
```
HTTP → handler (decode/authz) → app (use case) → domain (rules) → port ← adapter (Postgres / Stripe / Gemini / R2 / Resend)
```

## Cross-cutting infra (today single-instance; see roadmap to make multi-instance)
- Rate limiter (in-memory), job scheduler (in-process), WebSocket hub (in-memory) — all in `pkg/`.
  For horizontal scaling these become optional Redis-backed adapters.

# CleanSaaS

## What is this project

Open-source boilerplate for medium-to-large SaaS applications. Not a micro-SaaS starter — this targets real products with auth, billing, AI, admin, notifications, storage, and more.

**Stack**: Next.js (frontend) + Go with Chi (backend) + PostgreSQL (SQL pure, no ORM) + Neon (prod DB)
**Infra**: Vercel (frontend) + Railway (backend) + Neon (database) + Cloudflare R2 (storage)

This project is meant to showcase professional-grade engineering. Every file, every pattern, every decision should reflect that. We are building something worthy of mass adoption.

## Architecture

### Backend (Go) — Hexagonal Architecture
```
backend/
├── cmd/api/          → Entry point, dependency injection
├── internal/
│   ├── domain/       → Pure business logic, zero external dependencies
│   ├── port/         → Interfaces (repository + service contracts)
│   ├── app/          → Use cases, orchestration
│   ├── adapter/      → Concrete implementations (postgres, stripe, resend, claude, r2)
│   └── handler/      → HTTP transport (Chi router, middleware, DTOs)
├── pkg/              → Reusable public packages (jwt, hash, validate, pagination)
├── migrations/       → SQL migration files (up/down)
└── mock/             → Generated mocks from port interfaces
```

**Dependency rule**: handler → app → domain ← port ← adapter. Domain depends on NOTHING.

### Frontend (Next.js) — Feature-based Architecture
```
frontend/src/
├── app/              → Routes only (thin layer, App Router)
│   ├── (marketing)/  → Landing, blog, docs (public, SEO)
│   ├── (auth)/       → Login, register
│   └── (dashboard)/  → Authenticated app
├── features/         → Self-contained feature modules
│   ├── auth/
│   ├── billing/
│   └── ...
├── shared/           → Shared components, hooks, lib, types
└── config/           → App configuration
```

### Database
```
db/
└── init.sql          → Base schema (shared contract between all consumers)
```

SQL migrations live in `backend/migrations/`. The `db/` folder at root is the shared schema reference.

## Development principles

### Test-Driven Development
- Write tests FIRST or alongside code, never after
- Every use case in `app/` must have unit tests with mocked ports
- Every domain entity must have validation tests
- Every adapter should have integration tests
- Tests are how AI agents self-correct — they run tests, see failures, and fix autonomously
- Target: 80%+ coverage on domain and app layers

### SOLID principles
- **S**: One file, one responsibility. One function, one job.
- **O**: Extend via new adapters/features, don't modify existing ones.
- **L**: All implementations must satisfy their port interface fully.
- **I**: Keep interfaces small and specific. No god interfaces.
- **D**: Always depend on port interfaces, never on concrete implementations.

### STUPID avoidance
- **No Singletons**: Use dependency injection in main.go, explicit and visible.
- **No Tight coupling**: Layers communicate through interfaces only.
- **No Untestable code**: If it can't be tested with mocks, refactor it.
- **No Premature optimization**: Optimize when benchmarks prove a bottleneck.
- **No Indescriptive naming**: Every name must be self-explanatory. No abbreviations.
- **No Duplication**: Extract shared logic, but only when used 3+ times.

### Performance
- Go backend: connection pooling, prepared statements, context with timeouts
- Frontend: Server Components by default, client components only when needed
- Lazy loading, code splitting, image optimization
- Database: proper indexes, no N+1 queries, pagination everywhere

### Security
- Input validation at handler level (DTOs) AND domain level (entities)
- Parameterized SQL queries only — never string concatenation
- JWT with proper expiration and rotation
- Password hashing with bcrypt
- Rate limiting on auth endpoints
- CORS configured per environment
- No secrets in code — everything via environment variables
- Stripe webhook signature verification

### Scalability & Maintainability
- Stateless backend — ready for horizontal scaling
- Database migrations with up/down for safe rollbacks
- Feature flags ready (env-based)
- Each feature is isolated — adding/removing features has zero side effects
- Adapters are swappable: change payment provider = one new file + one line in main.go

### AI-Agent Friendliness
- Clear, consistent patterns across all features — learn one, know all
- Small files with single responsibility — fits in agent context window
- Explicit interfaces — agents see exactly what to implement
- Tests as guardrails — agents run tests to validate their work
- Descriptive naming — no ambiguity, no guessing
- CLAUDE.md at each level describes conventions for that area

## Code conventions

### Go
- Use `context.Context` as first parameter everywhere
- Return `(result, error)` — never panic in application code
- Errors wrap with `fmt.Errorf("doing x: %w", err)` for stack traces
- Use domain errors (`domain.ErrNotFound`) not HTTP status codes in app layer
- Handler converts domain errors to HTTP responses
- No global state except in `main.go` for wiring
- File naming: `entity.go`, `service.go`, `service_test.go` — lowercase, underscore

### TypeScript / Next.js
- Server Components by default
- `"use client"` only when interactivity is required
- Feature folder contains: components/, hooks/, actions/, types.ts
- Shared UI components in `shared/components/` (shadcn/ui)
- API calls via server actions or dedicated API client in `shared/lib/`
- Strict TypeScript — no `any`, no `as` casts unless absolutely necessary

### SQL
- Pure SQL, no ORM, no query builder
- Migrations numbered sequentially: `001_create_users.up.sql`, `001_create_users.down.sql`
- All tables have `id` (UUID), `created_at`, `updated_at`
- Use `TEXT` not `VARCHAR` (PostgreSQL best practice)
- Index foreign keys and frequently queried columns

### Git
- Conventional commits: `feat:`, `fix:`, `refactor:`, `chore:`, `test:`, `docs:`
- One feature per commit, atomic changes
- Never commit secrets, .env files, or node_modules

## Running locally

```bash
# Database
docker compose up -d

# Backend
cd backend && make run

# Frontend
cd frontend && npm run dev
```

## Environment variables

Backend reads from `DATABASE_URL`, `PORT`, `JWT_SECRET`, and service-specific keys.
Frontend reads from `NEXT_PUBLIC_API_URL`.
See `backend/internal/config/config.go` for the full list.

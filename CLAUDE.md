# CleanSaaS

## What is this project

Open-source boilerplate for medium-to-large SaaS applications. Not a micro-SaaS starter — this targets real products with auth, billing, AI, admin, notifications, storage, and more.

**Stack**: Next.js (frontend) + Go with Chi (backend) + PostgreSQL (SQL pure, no ORM) + Neon (prod DB)
**Infra**: Vercel (frontend) + Railway (backend) + Neon (database) + Cloudflare R2 (storage)

This project is meant to showcase professional-grade engineering. Every file, every pattern, every decision should reflect that.

## Core philosophy — Modularity above all

This is a boilerplate. We do NOT know what the end user will build with it. They might build a project management tool, an e-commerce platform, a social network, a CRM, or anything else. They will NOT use every module we provide.

**This means every feature must be fully independent and removable.**

Today, the boilerplate ships with all features included. Tomorrow, we will build a CLI tool (`create-cleansaas`) that lets users pick only the modules they need. The architecture must be ready for this evolution at all times.

### What this means in practice

**A feature is removable when:**
- Deleting its folder (frontend + backend) causes ZERO compilation errors elsewhere
- No other feature imports from it directly
- Its database tables can be dropped without breaking other tables
- The app still runs perfectly without it

**The only core modules that everything depends on are:**
- Auth (users must exist for anything to work)
- Database connection and config

Everything else — billing, AI, notifications, storage, admin, blog — is optional. A user might want billing but not AI. Another might want AI but not billing. Both must work.

### Rules to enforce this

1. **Features never import each other.** If billing needs user data, it receives a `UserRepository` interface via dependency injection — not by importing the user service directly. This way, billing depends on an interface, not on the user module.

2. **Database tables reference users via foreign key, but no cross-feature foreign keys.** The `subscriptions` table can reference `users(id)`, but the `conversations` table must NOT reference `subscriptions`. Each feature's tables are self-contained.

3. **Frontend features never import from other features.** Composition happens in `app/` pages only. If a settings page needs both profile and billing components, the page imports from both features — the features don't know about each other.

4. **Backend wiring in main.go is explicit.** Every feature is wired with its dependencies in `cmd/api/main.go`. Removing a feature means deleting its lines there. No auto-discovery, no magic registration.

5. **Server Actions are feature-scoped.** Each feature's `actions/` folder only calls API endpoints related to that feature. No cross-feature API calls from within a feature.

6. **Migrations are feature-prefixed.** Each feature's tables are created in their own migration files. Dropping a feature = skipping or reverting its migrations.

### How to verify independence

Before merging any feature, mentally test: "If I delete this feature's entire folder and its lines in main.go, does everything else still compile and run?" If the answer is no, there's a hidden dependency that must be fixed.

## Project structure

```
cleanSaaS/
├── frontend/          → Next.js 15, Tailwind, feature-based (see frontend/CLAUDE.md)
├── backend/           → Go + Chi, hexagonal architecture (see backend/CLAUDE.md)
│   └── migrations/    → SQL migration files (up/down)
├── .claude/skills/    → Custom Claude Code skills (/add-feature, /check, /add-migration, /test)
├── docker-compose.yml → PostgreSQL + DbGate (local DB viewer)
└── CLAUDE.md          → This file
```

Each major directory has its own CLAUDE.md with specific conventions.

## Development principles

### Test-Driven Development
- Write tests FIRST or alongside code, never after
- Tests are how AI agents self-correct — they run tests, see failures, and fix autonomously
- Target: 80%+ coverage on business logic layers

### SOLID principles
- **S**: One file, one responsibility. One function, one job.
- **O**: Extend via new adapters/features, don't modify existing ones.
- **L**: All implementations must satisfy their interface fully.
- **I**: Keep interfaces small and specific. No god interfaces.
- **D**: Always depend on interfaces, never on concrete implementations.

### STUPID avoidance
- **No Singletons**: Dependency injection, explicit and visible.
- **No Tight coupling**: Layers communicate through interfaces only.
- **No Untestable code**: If it can't be tested with mocks, refactor it.
- **No Premature optimization**: Optimize when benchmarks prove a bottleneck.
- **No Indescriptive naming**: Every name must be self-explanatory.
- **No Duplication**: Extract shared logic, but only when used 3+ times.

### Performance
- Go: connection pooling, context with timeouts, no N+1 queries
- Next.js: Server Components by default, lazy loading, code splitting
- Database: proper indexes, pagination everywhere

### Security
- Input validation at every boundary (HTTP + domain)
- Parameterized SQL queries only — never string concatenation
- JWT with proper expiration, bcrypt for passwords
- Rate limiting on sensitive endpoints
- No secrets in code — everything via environment variables

### Scalability & Maintainability
- Stateless backend — horizontal scaling ready
- Features are isolated — adding/removing has zero side effects
- Adapters are swappable: change provider = one new file + one line change
- Database migrations with up/down for safe rollbacks

### AI-Agent Friendliness
- Consistent patterns across all features — learn one, know all
- Small files with single responsibility — fits in agent context window
- Explicit interfaces — agents see exactly what to implement
- Tests as guardrails — agents validate their own work
- CLAUDE.md at each level describes conventions for that area

## Code conventions

### SQL & Migrations
- Pure SQL, no ORM, no query builder. Powered by `golang-migrate`.
- Migrations live in `backend/migrations/`: `001_name.up.sql` / `001_name.down.sql`
- All tables: UUID `id`, `created_at`, `updated_at`
- Use `TEXT` not `VARCHAR`, index foreign keys
- No cross-feature foreign keys (only reference `users` table)
- Migrations are immutable once applied in prod — never edit, only create new ones
- Workflow: create migration → test locally (`make migrate-up`) → commit → apply to prod (`DATABASE_URL=<prod> make migrate-up`)

### Git
- Conventional commits: `feat:`, `fix:`, `refactor:`, `chore:`, `test:`, `docs:`
- One feature per commit, atomic changes
- Never commit secrets, .env files, or node_modules

## Running locally

```bash
# Database
docker compose up -d

# Apply migrations
cd backend && make migrate-up

# Seed data (admin user)
cd backend && make seed

# Backend
cd backend && make run

# Frontend
cd frontend && npm run dev
```

## Environment variables

- **Backend**: `DATABASE_URL`, `PORT`, `JWT_SECRET`, service-specific keys. See `backend/internal/config/config.go`.
- **Frontend**: `NEXT_PUBLIC_API_URL`

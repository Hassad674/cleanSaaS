# CleanSaaS

## What is this project

Open-source boilerplate for medium-to-large SaaS applications. Not a micro-SaaS starter — this targets real products with auth, billing, AI, admin, notifications, storage, and more.

**Stack**: Next.js (frontend) + Go with Chi (backend) + PostgreSQL (SQL pure, no ORM) + Neon (prod DB)
**Infra**: Vercel (frontend) + Railway (backend) + Neon (database) + Cloudflare R2 (storage)

This project is meant to showcase professional-grade engineering. Every file, every pattern, every decision should reflect that.

## Project structure

```
cleanSaaS/
├── frontend/          → Next.js 15, Tailwind, feature-based architecture
├── backend/           → Go + Chi, hexagonal architecture (see backend/CLAUDE.md)
├── db/                → Shared SQL schema reference
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

### SQL
- Pure SQL, no ORM, no query builder
- Migrations: `001_name.up.sql` / `001_name.down.sql`
- All tables: UUID `id`, `created_at`, `updated_at`
- Use `TEXT` not `VARCHAR`, index foreign keys

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

- **Backend**: `DATABASE_URL`, `PORT`, `JWT_SECRET`, service-specific keys. See `backend/internal/config/config.go`.
- **Frontend**: `NEXT_PUBLIC_API_URL`

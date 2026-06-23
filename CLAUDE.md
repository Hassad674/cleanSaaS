# CleanSaaS

## What is this project

Open-source boilerplate for medium-to-large SaaS applications. Not a micro-SaaS starter — this targets real products with auth, billing, AI, admin, notifications, storage, and more.

**Stack**: Next.js (frontend) + Go with Chi (backend) + PostgreSQL (SQL pure, no ORM) + Neon (prod DB)
**Infra**: Vercel (frontend) + Railway (backend) + Neon (database) + Cloudflare R2 (storage)

This project is meant to showcase professional-grade engineering. Every file, every pattern, every decision should reflect that.

## Default operating procedure for EVERY request (read this first)

This boilerplate is built so that **any** user — a seasoned engineer or a non-developer "vibe coding" in plain English — gets **top-tier, production-grade results by default**, without having to know the internals. When a user asks for something, do NOT free-style. Classify the request and follow the matching rail; each rail already encodes the architecture, tests, and checks. This is what makes quality the default, not an afterthought.

| User says (plain English) | Your default rail |
|---|---|
| "build/add <feature>" (e.g. "add a projects feature with tasks") | Use the **`/add-feature`** skill — scaffolds all layers domain-first, with tests, respecting modularity. |
| "add an endpoint / action to <feature>" | Use **`/add-endpoint`**. |
| "use <provider> instead" / "add Stripe/Resend/etc." | Use **`/add-adapter`** (implement the port, change one wiring line). |
| "add a table / column / store <data>" | Use **`/add-migration`** (numbered up/down, no cross-feature FK). |
| "remove <feature>" | Use **`/remove-feature`**; verify with **`/verify-independence`**. |
| "it's broken / this doesn't work / <screenshot>" | Use **`/debug`** (guided reproduce → bug report → failing test → fix → verify). |
| "run it / start the app / does it work?" | Use **`/run`**, then smoke-test. |
| "is this correct / review this" | Use **`/check`** (architecture) and **`/review`** (quality/security). |
| anything that changed code | Finish with the **Validation pipeline** below, then **`/check`**. |

**Non-negotiables on every change, regardless of who asked or how casually:**
1. Follow the layer rules (backend hexagonal `handler→app→domain←port←adapter`; frontend feature-based, no cross-feature imports). Never inline business logic in a handler/page.
2. Write tests alongside the code (TDD loop). A change isn't done until tests + the validation pipeline are green.
3. Respect the hard code-quality limits (≤600 lines/file, ≤50/func, ≤4 params, depth 3, cyclo <10) and the design tokens (no hardcoded colors).
4. Parameterized SQL only; validate input at boundaries; never commit secrets.
5. If you're unsure what the user wants, ask ONE clarifying question — but default to the most maintainable option and proceed.

The CI gates, git hooks, and these skills exist so that the floor for *any* prompt is professional-grade. Hold that floor. For a non-developer, translate jargon, explain what you did in one plain sentence, and lean on `/run` + `/debug` so they can see it work.

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
├── frontend/          → Next.js 16, Tailwind v4, feature-based (see frontend/CLAUDE.md)
├── backend/           → Go 1.25 + Chi, hexagonal architecture (see backend/CLAUDE.md)
│   └── migrations/    → SQL migration files (up/down)
├── admin/             → Vite + React + Tailwind, admin dashboard (see admin/CLAUDE.md)
├── .claude/skills/    → Custom Claude Code skills (/run, /debug, /add-feature, …)
├── .claude/memory/    → Project memory (MEMORY.md, architecture.md)
├── docker-compose.yml → PostgreSQL + DbGate (local DB viewer)
└── CLAUDE.md          → This file
```

Each major directory has its own CLAUDE.md with specific conventions.
Blockers, when they happen, are logged in `BLOCKED-<topic>.md` at the repo root (see Blocker policy below) — none exist in a healthy clone.

## Development principles

### Test-Driven Development
- Write tests FIRST or alongside code, never after
- Tests are how AI agents self-correct — they run tests, see failures, and fix autonomously
- Target: 80%+ coverage on business logic layers

### Code quality limits (hard — these are non-negotiable, not guidelines)
These caps keep files in an agent's context window and force single-responsibility. If you're about to exceed one, that's the signal to split — extract a function, a type, or a file.
- **≤ 600 lines per file.** Over that, the file is doing too much — split by responsibility.
- **≤ 50 lines per function.** Over that, extract helpers.
- **≤ 4 parameters per function/constructor.** More → pass a struct (e.g. a `Deps`/`Params` struct). The composition root may use a typed `Deps` struct rather than 16 positional args.
- **≤ 3 levels of nesting.** Deeper → use early returns / guard clauses.
- **Cyclomatic complexity < 10 per function.** Higher → break the branching apart.
- **Forbidden indescriptive names:** `data`, `info`, `tmp`, `temp`, `manager`, `util`/`utils`, `helper`/`helpers`, `handler2`, `doStuff`. Every name states what it is.
These are enforced in CI (golangci-lint + the gate scripts in `scripts/ci/`), not just trusted.

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

## Skills reference

Use these slash commands to accelerate development:

| Situation | Skill | Example |
|-----------|-------|---------|
| Start the whole stack locally | `/run` | `/run` |
| Reproduce & fix a bug (guided) | `/debug` | `/debug login button does nothing` |
| New full-stack module | `/add-feature` | `/add-feature billing with plans and subscriptions` |
| New endpoint on existing feature | `/add-endpoint` | `/add-endpoint reset-password on auth` |
| New/swap external provider | `/add-adapter` | `/add-adapter lemonsqueezy for payment` |
| New database table or column | `/add-migration` | `/add-migration create subscriptions table` |
| Remove a module cleanly | `/remove-feature` | `/remove-feature billing` |
| Prove a module is removable (no mutation) | `/verify-independence` | `/verify-independence referral` |
| Verify architecture rules | `/check` | `/check` or `/check billing` |
| Code review before commit | `/review` | `/review` or `/review auth` |
| Run tests intelligently | `/test` | `/test auth` or `/test changed` |
| Run the Playwright e2e suite | `/e2e` | `/e2e` or `/e2e auth` |

Skills are auto-detected: describing what you want is enough. You don't have to type the slash command explicitly.

## Autonomous work process

### Test tools
- **Backend**: Go `testing` package + `github.com/stretchr/testify` (assertions, mocks)
- **Frontend unit**: `vitest` + `@testing-library/react` + `@testing-library/jest-dom`
- **Frontend E2E**: `playwright` (chromium) — tests in `frontend/e2e/`
- All are already declared in the respective `package.json` / `go.mod`. Install with `go mod download`, `npm install`, and `npx playwright install --with-deps chromium` for e2e.

### Running & observing the app (don't only compile + unit-test — run it)
Compiling and unit-testing is not enough; many bugs only show up at runtime. Know how to start the stack, read its logs, and reproduce behavior.
- **Start everything:** use the `/run` skill (DB → migrate → seed → backend :8081 → frontend :3010 → admin :5174, with smoke tests). Or manually: `docker compose up -d` → `cd backend && make migrate-up && make seed && make run` → `cd frontend && npm run dev`.
- **Ports:** backend **:8081** · frontend **:3010** · admin **:5174** · Postgres **:5433** · DbGate **:8082**. Seeded admin: `admin@cleansaas.dev` / `admin123`.
- **Health:** `curl http://localhost:8081/health` → `{"status":"ok",...}` (200). Routes are mounted at **root** (`/auth/login`, not `/api/...`).
- **Logs:** when started via `/run`, tailed to `/tmp/cleansaas-{backend,frontend,admin}.log`. Read them when something misbehaves — backend logs are structured JSON (`slog`).
- **Inspect data:** DbGate at `http://localhost:8082` is a web UI on the local Postgres.
- **Reproduce a bug / get unstuck:** use the `/debug` skill — it walks through screenshot intake, browser/Chrome-extension or `curl` reproduction, localization, a bug report, a failing test, the fix, and verification.

### Test → Fix → Retest loop (MANDATORY)

Every piece of code you write must be tested. Follow this loop:

```
1. Write implementation code
2. Write unit tests for that code
3. Run tests
   ├── ALL PASS → continue to next sub-task ✅
   └── FAIL →
       4. Read error output carefully
       5. Fix the bug (in code OR test, whichever is actually wrong)
       6. Rerun tests
       (max 3 fix attempts per failing test)
       Still failing after 3 attempts → blocker policy below
```

**NEVER commit with failing tests. NEVER delete or skip a test to make the suite pass.**

### Commit strategy
- 1 commit per completed task (not per sub-step)
- Before EVERY commit: run the full validation pipeline below
- Never commit broken code on main
- Conventional messages: `feat:`, `fix:`, `test:`, `refactor:`, `chore:`, `docs:`

### Validation pipeline (run before EVERY commit)

```bash
# 1. Compilation
cd backend && go build ./...
cd frontend && npx tsc --noEmit

# 2. Backend tests
cd backend && go test ./... -count=1

# 3. Frontend tests (if tests exist)
cd frontend && npx vitest run

# 4. E2E tests (only after Playwright is set up)
cd frontend && npx playwright test

# 5. Architecture checks
# - No hardcoded colors (zinc-, gray-, slate-, white, black in className)
# - No cross-feature imports in frontend/src/features/
# - Migration has both .up.sql and .down.sql
```

ALL steps must pass. If any fails → enter fix loop above → only commit when ALL green.

### Blocker policy

**A long task is NOT a blocker.** A blocker = same error, 3+ approaches tried, no progress. Stripe taking 2h is normal. Stuck on the same error for 20 min is a blocker.

**Type A — Test failure**: max 3 fix attempts per test → comment `// TODO: fix` → log `BLOCKED-taskX.md` → continue other sub-steps

**Type B — Same error, no progress**: 3+ different approaches tried, nothing works → log `BLOCKED-taskX.md` with error + all approaches → skip only the blocked sub-step (not the whole task) → commit working code → move on

**Type C — Compilation failure**: TOP PRIORITY, 10 min to fix → if unfixable, revert latest changes (`git checkout -- <files>`) → NEVER leave build broken

When you hit a Type A/B blocker, log it in `BLOCKED-<topic>.md` at the repo root (error + every approach tried + suspected location) and note it in your active progress file.

## Compact instructions

When compacting context, prioritize preserving:
1. The current task being worked on (which phase / sub-step)
2. Any errors or blockers encountered
3. Files recently created or modified
4. The test→fix→retest loop state (what's failing, what was tried)

After compaction:
1. Re-read your active progress tracker if one exists (e.g. `AUTONOMY-LOG.md` at the repo root during a multi-phase work session) to recover the task list and checkboxes
2. Run `git branch --show-current` and `git log --oneline -15` to see the active branch and what was already committed
3. Run `cd backend && go build ./... && go test ./...` and `cd frontend && npx tsc --noEmit && npx vitest run` to confirm the baseline is green
4. Check for `BLOCKED-*.md` files at the repo root
5. Resume from the first unchecked item in the active tracker

## Environment variables

- **Backend**: `DATABASE_URL`, `PORT`, `JWT_SECRET`, service-specific keys. See `backend/internal/config/config.go`.
- **Frontend**: `NEXT_PUBLIC_API_URL`
- **Admin**: `VITE_API_URL`
- Never commit `.env` files — they are gitignored

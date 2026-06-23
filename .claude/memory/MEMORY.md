# CleanSaaS — Project Memory

> This memory ships with the boilerplate. It must contain only facts useful to **anyone** who clones the
> repo — architecture invariants, the module map, ports, conventions. It must **never** contain a
> maintainer's machine config, secrets, OS/shell settings, or personal preferences. Put machine-local notes
> in a gitignored file (e.g. `.claude/memory/local.md`), not here.

## Project identity
- Open-source SaaS boilerplate for **medium-to-large** SaaS (not micro-SaaS).
- Stack: **Next.js 16** (frontend) + **Go 1.25 / Chi** (backend) + **PostgreSQL 16** (pure SQL, no ORM) + **Tailwind v4**.
- Prod infra targets: Vercel (frontend) · Railway (backend) · Neon (DB) · Cloudflare R2 (storage).
- Core thesis: the product is the **enforcement system** (CLAUDE.md + memory + skills + CI + hooks) that keeps
  AI-generated code maintainable/performant/secure — for beginners *and* pros.

## Architecture invariants (non-negotiable)
- **Backend = hexagonal**, dependency rule: `handler → app → domain ← port ← adapter`.
  - `domain/` imports stdlib only. `port/` imports domain only. `app/` imports domain + port (interfaces).
    `adapter/` imports domain + port + external libs. An adapter never imports another adapter.
- **Frontend = feature-based.** Features never import each other; composition happens only in `app/` pages.
- **Every feature is fully removable**: delete its folder + its `cmd/api/main.go` wiring lines → everything
  else still compiles. Only **auth + DB** are hard core. Verify with the `/check` and `/verify-independence` skills.
- **No cross-feature foreign keys** — feature tables may only `REFERENCES users(id)`.
- **Swap a provider in one line**: implement the `port/service` interface in a new `adapter/`, change one line in `main.go`.

## Module map (backend feature → optional? → env key that enables it)
| Feature | Core? | Enabled by | External adapter |
|---|---|---|---|
| auth, user | **core** | always on | — |
| billing | optional | `STRIPE_SECRET_KEY` | adapter/stripe |
| storage | optional | `R2_ACCOUNT_ID` / `R2_*` | adapter/r2 |
| ai | optional | `GEMINI_API_KEY` | adapter/gemini |
| email (used by auth/notification) | optional | `RESEND_API_KEY` | adapter/resend |
| notification, blog, team | optional | always on (DB only) | — |
| referral | optional | always on (DB only) | — |

If an optional key is empty, `cmd/api/main.go` skips wiring that feature (nil service) — the app still runs.

## Ports & local facts
- Backend API: **:8081** · Frontend: **:3010** · Admin (Vite): **:5174** · Postgres: **:5433** · DbGate UI: **:8082**
- API routes are mounted at **root** (e.g. `POST /auth/login`), not under `/api`.
- Seeded admin (via `make seed`): **admin@cleansaas.dev** / **admin123**.
- Neon prod connection: drop `channel_binding=require` for `lib/pq` compatibility.
- One-command local start: the `/run` skill (or `scripts/bootstrap.sh` then start the dev servers).

## Conventions
- Code/comments/docs in **English**. Conventional commits (`feat:`/`fix:`/`refactor:`/`chore:`/`test:`/`docs:`).
- TDD: write tests alongside code; agents self-correct via the Test→Fix→Retest loop (see root `CLAUDE.md`).
- Validation pipeline before every commit: `go build ./...` · `go test ./...` · `npx tsc --noEmit` · `npx vitest run` (+ Playwright when a flow changed).
- Migrations: numbered `NNN_name.up.sql`/`.down.sql`, immutable once applied, idempotent (`IF [NOT] EXISTS`).

## Skills (use these — `/<name>`)
- `/run` — start the whole stack locally & smoke-test it.
- `/debug` — guided, beginner-friendly reproduce→report→fix→verify (screenshot + Chrome extension + tests).
- `/test`, `/e2e` — run Go/Vitest tests, run Playwright e2e.
- `/add-feature`, `/add-endpoint`, `/add-adapter`, `/add-migration` — scaffolders (layer-correct, domain-first).
- `/remove-feature`, `/verify-independence` — prove/perform clean module removal (the modularity promise).
- `/check`, `/review` — architecture-independence audit & pre-commit code review.
- `/autopilot` — heavy/long autonomous tasks: auto-checkpoint to `PROGRESS.md` + delegate to subagents + commit each verified step, so work survives context compaction.

## Heavy autonomous tasks (auto behavior)
For any heavy/multi-phase request ("build everything", big migration/audit, "while I'm away", est. >700K tokens), the agent AUTOMATICALLY (per root CLAUDE.md): keeps a committed `PROGRESS.md` checkpoint, delegates substantial sub-work to subagents (keeps main context lean), verifies each subagent's output, and commits per verified increment — resumable after compaction. Small tasks skip this overhead. See `AUTONOMY-LOG.md` for a real worked example.

## Pointers
- Detailed layer structure → [architecture.md](architecture.md)
- Per-area conventions → `CLAUDE.md` (root), `backend/CLAUDE.md`, `frontend/CLAUDE.md`, `admin/CLAUDE.md`

# Contributing to CleanSaaS

Thanks for contributing. CleanSaaS is a boilerplate whose real product is its
**enforcement system** — the architecture rules, skills, CI gates, and hooks that keep
AI-generated and human-written code maintainable, performant, and secure. Contributions
are held to that bar. This document is the contract: read it before opening a PR.

If anything here is unclear, see **[Where to ask questions](#where-to-ask-questions)**.

---

## Table of contents

- [Local setup](#local-setup)
- [The Modularity Invariant](#the-modularity-invariant)
- [The validation pipeline](#the-validation-pipeline)
- [Code-quality limits (hard)](#code-quality-limits-hard)
- [Adding a feature](#adding-a-feature)
- [Migration rules](#migration-rules)
- [Commit format](#commit-format)
- [Git hooks](#git-hooks)
- [Branch & PR etiquette](#branch--pr-etiquette)
- [Where to ask questions](#where-to-ask-questions)

---

## Local setup

### Prerequisites

| Tool   | Version  | Notes                                              |
|--------|----------|----------------------------------------------------|
| Docker | latest   | Provides Postgres + DbGate locally                 |
| Go     | **1.25+**| Required by `backend/go.mod`                        |
| Node   | **20+**  | A `.nvmrc` is present at the repo root              |
| Git    | latest   | On Windows, install Git for Windows for Git Bash    |

### Bootstrap

Fastest path is the **`/run`** skill (it brings the whole stack up and smoke-tests it).
Otherwise, run the bootstrap script directly:

```bash
./scripts/bootstrap.sh        # Linux / macOS / Git Bash
.\scripts\bootstrap.ps1       # Windows PowerShell
```

This checks prerequisites, copies `.env.example` templates, starts Docker, waits for
Postgres, runs migrations, and seeds the database. See the
[README](./README.md#quickstart) for the manual, step-by-step alternative.

### Ports

| Service          | URL / Address           |
|------------------|-------------------------|
| Backend (Go API) | `http://localhost:8081` |
| Frontend (Next)  | `http://localhost:3010` |
| Admin (Vite)     | `http://localhost:5174` |
| Postgres         | `localhost:5433`        |
| DbGate (DB UI)   | `http://localhost:8082` |

API routes are mounted at **root** (`POST /auth/login`, not `/api/...`). Seeded local
admin: `admin@cleansaas.dev` / `admin123` (development only — never use in production).

---

## The Modularity Invariant

This is the **single most important rule in the project**, and it is enforced, not
merely encouraged.

> **A feature is removable if and only if deleting its folder (frontend + backend)
> plus its wiring lines in `backend/cmd/api/main.go` leaves everything else compiling
> and running.**

Concretely, this means:

- **Features never import each other.** If `billing` needs user data, it depends on a
  `UserRepository` interface injected at the composition root — it does **not** import
  the `user` package. Frontend features never import from sibling features; composition
  happens only in `app/` pages.
- **No cross-feature foreign keys.** A feature's tables may `REFERENCES users(id)` and
  nothing else. `subscriptions` may reference `users`; `conversations` must not reference
  `subscriptions`.
- **Wiring is explicit in `main.go`.** No auto-discovery, no magic registration. Removing
  a feature = deleting its folder and its lines in `cmd/api/main.go`.

The only hard core is **auth + user + the database**. Everything else — `billing`,
`storage`, `ai`, `notification`, `blog`, `team` — must be independently removable.

**Verify before every PR that touches feature boundaries:**

```
/verify-independence <feature>   # proves a module is cleanly removable (no mutation)
/check                           # audits all architecture-independence rules
```

A PR that breaks the invariant will not be merged until the hidden dependency is removed.

---

## The validation pipeline

Run the full pipeline locally before pushing. CI runs the same checks; the
**`scripts/ci/run-all.sh`** gate must be green.

**Backend (`cd backend`)**

```bash
go build ./...
go test ./... -count=1
```

**Frontend (`cd frontend`)**

```bash
npx tsc --noEmit
npx vitest run
npx playwright test          # E2E (requires: npx playwright install --with-deps chromium)
```

**Admin (`cd admin`)**

```bash
npx tsc --noEmit
```

**Invariant + quality gates (from repo root)**

```bash
bash scripts/ci/run-all.sh   # cross-feature imports, file length, forbidden names,
                             # hardcoded colors, migration up/down pairing
```

All steps must pass. Never push or commit broken builds or failing tests, and never
delete or skip a test just to make the suite green.

---

## Code-quality limits (hard)

These are **non-negotiable** and enforced in CI (golangci-lint via
`backend/.golangci.yml` plus the gate scripts in `scripts/ci/`). They keep files inside
an agent's context window and force single-responsibility.

- **≤ 600 lines per file.**
- **≤ 50 lines per function.**
- **≤ 4 parameters per function/constructor** — more → pass a `Deps`/`Params` struct.
- **≤ 3 levels of nesting** — deeper → use early returns / guard clauses.
- **Cyclomatic complexity < 10 per function.**
- **No indescriptive names**: `data`, `info`, `tmp`, `temp`, `manager`, `util(s)`,
  `helper(s)`, `handler2`, `doStuff`, etc. Every name states what it is.

If you are about to exceed one of these, that is the signal to split — extract a
function, a type, or a file.

---

## Adding a feature

Don't hand-roll a feature from scratch. Use the scaffolder and follow the layer
conventions:

- Run the **`/add-feature`** skill (e.g. `/add-feature billing with plans and
  subscriptions`). It generates a layer-correct, domain-first module.
- Follow the architecture rules in **`backend/CLAUDE.md`** (hexagonal:
  `handler → app → domain ← port ← adapter`) and **`frontend/CLAUDE.md`**
  (feature-based, composition only in `app/` pages).
- Related scaffolders: `/add-endpoint`, `/add-adapter`, `/add-migration`,
  `/remove-feature`.

Before opening the PR, run `/verify-independence <feature>` and `/check` to prove the
new module honors the Modularity Invariant.

---

## Migration rules

Pure SQL, no ORM, powered by `golang-migrate`. Migrations live in `backend/migrations/`.

- **Numbered up/down pairs**: `NNN_name.up.sql` and `NNN_name.down.sql`. Both files are
  required (CI checks the pairing).
- **Immutable once applied.** Never edit a migration that has run in production — only
  add a new one.
- **Never renumber.** Sequence numbers are append-only history.
- **No cross-feature foreign keys.** Feature tables may only `REFERENCES users(id)`.
- All tables carry UUID `id`, `created_at`, `updated_at`; use `TEXT` not `VARCHAR`;
  index foreign keys; write idempotent DDL (`IF [NOT] EXISTS`).
- Workflow: create → `make migrate-up` locally → commit → apply to prod
  (`DATABASE_URL=<prod> make migrate-up`).

---

## Commit format

Use [Conventional Commits](https://www.conventionalcommits.org/). One logical change
per commit, atomic.

```
<type>(<optional scope>): <short imperative summary>
```

Allowed types: `feat`, `fix`, `refactor`, `chore`, `test`, `docs`.

Examples:

```
feat(auth): add password-reset endpoint
fix(billing): handle nil Stripe customer on cancellation
refactor(storage): extract R2 client into adapter
test(team): cover invite-expiry edge cases
docs(readme): clarify Windows bootstrap steps
chore(ci): bump golangci-lint to latest
```

Never commit secrets, `.env` files, or `node_modules`.

---

## Git hooks

Git does **not** auto-install hooks on clone. Run this **once per clone** to point Git
at the repo's tracked `.githooks/` directory:

```bash
bash scripts/install-git-hooks.sh
```

The pre-commit hook runs the fast quality and invariant gates before a commit lands, so
you catch violations locally instead of in CI. It is idempotent — re-running just
re-confirms the configuration.

---

## Branch & PR etiquette

- **Branch off the default branch.** Name branches by type and scope, e.g.
  `feat/auth-password-reset`, `fix/billing-nil-customer`.
- **Keep PRs focused.** One feature or fix per PR. Split unrelated changes.
- **Green before review.** The full [validation pipeline](#the-validation-pipeline) and
  the `ci-gate` job must pass. The heavy E2E workflow runs on PRs carrying the
  `run-e2e` label.
- **Describe the change**: what, why, and how you verified it. Link any related issue.
  If it touches feature boundaries, note that you ran `/verify-independence` and `/check`.
- **Update docs and `CHANGELOG.md`** (the `[Unreleased]` section) when behavior,
  conventions, or the public surface change.
- **No unrelated reformatting.** Keep diffs reviewable.

---

## Where to ask questions

- **Bugs and feature requests** → open a GitHub Issue.
- **Design / "is this the right approach" discussions** → open a GitHub Discussion (or a
  draft PR) before investing heavily, so we can align on architecture early.
- **Security vulnerabilities** → do **not** open a public issue. Follow
  [`SECURITY.md`](./SECURITY.md).
- **Conduct concerns** → see [`CODE_OF_CONDUCT.md`](./CODE_OF_CONDUCT.md).

By contributing, you agree your contributions are licensed under the same license as
the project.

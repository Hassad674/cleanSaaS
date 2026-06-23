# 0002. Pure SQL with `database/sql`, no ORM

## Status

Accepted

## Context

Every feature in CleanSaaS persists data to PostgreSQL 16. We had to choose how
the `adapter/postgres` layer talks to the database. The options span a spectrum
from full ORM (GORM, ent) through code generation (sqlc) to hand-written SQL on
the standard library.

The forces:

1. **Explicitness for AI agents.** An agent reading a repository file should see
   the exact SQL that runs — no hidden N+1 queries, no implicit eager/lazy
   loading, no surprise migrations generated from struct tags. What you read is
   what executes.
2. **Performance and control.** The product targets real, growing SaaS apps. We
   want hand-written, index-aware queries, explicit pagination, and
   `context`-driven timeouts (see the perf rules in the root `CLAUDE.md`).
3. **A thin, swappable adapter boundary.** ADR 0001 already isolates persistence
   behind `port/repository` interfaces. The adapter behind them should be a
   straightforward translation of interface → SQL, with no second abstraction
   competing with the port.
4. **Security.** Parameterized queries only, enforced as a hard rule. No string
   concatenation into SQL, ever.

An ORM would re-introduce a heavy framework right at the boundary we worked to
keep thin in ADR 0001, and would hide the very SQL we want agents and reviewers
to see.

## Decision

We will use **pure SQL on `database/sql` with the `lib/pq` driver**. No ORM, no
query builder.

- Queries are hand-written in the `adapter/postgres` package, one file per
  repository, implementing the corresponding `port/repository` interface.
- **Parameterized queries only** (`$1, $2, …`). String concatenation into SQL is
  forbidden and rejected in review.
- Every query takes a `context.Context` for timeout and cancellation.
- Schema changes are managed by **`golang-migrate`** with numbered, paired
  migrations: `NNN_name.up.sql` and `NNN_name.down.sql` in `backend/migrations/`.
  Every `up` has a reversible `down`. Migrations use `IF [NOT] EXISTS` to stay
  idempotent.
- Conventions: UUID `id`, `created_at`/`updated_at` on every table, `TEXT` over
  `VARCHAR`, foreign keys indexed, and — per ADR 0003 — **no cross-feature
  foreign keys** (a table may only `REFERENCES users(id)`).

## Consequences

**Positive (+)**

- Total control over the executed SQL: tuned queries, explicit JOINs, no
  accidental N+1s, predictable performance.
- Zero ORM magic — the adapter is exactly as legible as the SQL it contains,
  which suits both beginners and AI agents.
- Migrations are plain, reviewable SQL with explicit up/down, decoupled from
  application structs; the schema cannot silently drift from a model definition.
- Minimal dependency surface (`database/sql` + `lib/pq` + `golang-migrate`).

**Negative (−)**

- More boilerplate: every CRUD method is hand-written `Query`/`Exec` plus row
  scanning. This is the single largest source of repetition in the backend.
- The team owns mapping rows to structs and handling `sql.ErrNoRows` →
  `domain.ErrNotFound` consistently in every repository.
- No compile-time guarantee that a query string matches the Go types scanning it
  (sqlc would give this); correctness is enforced by integration tests against a
  real Postgres instead.
- Refactoring a column touches raw SQL strings rather than a single model
  definition.

## Alternatives considered

- **GORM.** The most popular Go ORM. Rejected: hides the SQL, makes N+1s easy to
  introduce and hard to see, encourages model-driven auto-migration that fights
  our explicit numbered-migration workflow, and adds a heavy abstraction exactly
  where ADR 0001 wants a thin one — directly against the explicitness force.
- **sqlc (compile SQL → typed Go).** The strongest contender: you still write raw
  SQL and get type-safe generated accessors. Rejected *for now* to avoid a code
  generation step in the build and the extra concept for beginners; the hand-
  written adapter is simple enough at the current scale. Worth revisiting if
  repository boilerplate becomes a real maintenance burden.
- **ent (Facebook's entity framework).** Powerful graph/codegen ORM. Rejected as
  the heaviest option — a large framework and schema DSL that contradict the
  "explicit, minimal, swappable" goals and would dominate the persistence layer.

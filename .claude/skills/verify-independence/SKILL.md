---
name: verify-independence
description: Prove a feature is fully removable by simulating its removal in a throwaway git worktree and running build+tests, without touching the working repo.
user-invocable: true
allowed-tools: Read, Bash, Grep, Glob
---

# Verify Independence

Target feature: **$ARGUMENTS**

You PROVE CleanSaaS's core promise — *"every feature is fully removable"* — on demand, **without mutating the real repo**. You simulate removing a feature inside an isolated, throwaway git worktree, run the full build + tests there, report whether everything still compiles, then discard the worktree.

This is the **safe, non-destructive sibling** of `/remove-feature`. Where `/remove-feature` actually deletes a module from the working tree, this skill only *simulates* it in a disposable copy and reports the verdict. It reuses the same understanding of what a feature spans. This skill is also the backbone of the future `create-cleansaas` CLI (pick-your-modules).

**Golden rule: the real working tree is NEVER modified. Every deletion and edit happens inside `/tmp/ci-independence-<feature>` only.** Use relative paths and `$CLAUDE_PROJECT_DIR`, never absolute user paths.

---

## What a feature spans (verified repo facts)

A **backend feature** `<f>` spans these layers:

```
backend/internal/domain/<f>/                          → Domain entities
backend/internal/app/<f>/                             → Application service
backend/internal/port/repository/<f>*.go              → Repository interface
backend/internal/port/service/<f>*.go                 → Service interface (if any)
backend/internal/adapter/postgres/<f>.go              → PostgreSQL adapter
backend/internal/adapter/<provider>/                  → External adapter (stripe, resend, claude...) if any
backend/internal/handler/<f>.go                       → HTTP handler
backend/internal/handler/dto/request/<f>.go           → Request DTOs
backend/internal/handler/dto/response/<f>.go          → Response DTOs (if separate)
backend/migrations/NNN_*<f>*.up.sql / .down.sql       → Migrations
```

Plus its **wiring**, which lives in two files (these are MODIFIED, not deleted):

```
backend/cmd/api/main.go              → DI lines: repo → service → handler/router params
backend/internal/handler/router.go   → route registrations + the service param in NewRouter()
```

A **frontend feature** `<f>` spans:

```
frontend/src/features/<f>/           → Entire feature directory
frontend/src/app/**/<f>/             → Route pages for this feature (+ any demo page)
```

**Core features are NOT removable** — `auth` and `user` (and the DB connection/config) are depended upon by everything. If `$ARGUMENTS` is a core feature, STOP and explain it cannot be removed.

---

## Flow

### STEP 1 — Resolve the feature name

Take the feature from `$ARGUMENTS`.

- If `$ARGUMENTS` is **empty**, list the candidate removable features and ask which one. Discover them, don't guess:
  ```bash
  ls -1 "$CLAUDE_PROJECT_DIR/backend/internal/app" 2>/dev/null
  ls -1 "$CLAUDE_PROJECT_DIR/frontend/src/features" 2>/dev/null
  ```
  Present the list **minus the core features** (`auth`, `user`) and ask the user to pick one. Then STOP and wait.
- If `$ARGUMENTS` matches `auth` or `user`, STOP: explain these are core and cannot be removed.

### STEP 2 — Create the throwaway worktree

Create an isolated git worktree so the real working tree is untouched. Tell the user explicitly: *"This is a throwaway copy at `/tmp/ci-independence-<feature>`. Your real repo will not be modified."*

```bash
cd "$CLAUDE_PROJECT_DIR"
WT="/tmp/ci-independence-$FEATURE"
# Clean any stale worktree from a previous interrupted run
git worktree remove --force "$WT" 2>/dev/null || true
git worktree add "$WT" HEAD
```

From here on, **operate ONLY inside `$WT`**. Never run delete/edit commands against `$CLAUDE_PROJECT_DIR`.

### STEP 3 — Locate ALL files belonging to the feature (inside the worktree)

Inside `$WT`, find every file across the layers listed above. Be exhaustive — use Glob and Grep, do not guess. List the full set for the user before deleting anything.

```bash
# Backend folders + files (note the adapter provider may be named differently than the feature)
ls -d  "$WT"/backend/internal/domain/$FEATURE                 2>/dev/null
ls -d  "$WT"/backend/internal/app/$FEATURE                    2>/dev/null
ls     "$WT"/backend/internal/port/repository/$FEATURE*.go    2>/dev/null
ls     "$WT"/backend/internal/port/service/$FEATURE*.go       2>/dev/null
ls     "$WT"/backend/internal/adapter/postgres/$FEATURE.go    2>/dev/null
ls     "$WT"/backend/internal/handler/$FEATURE.go             2>/dev/null
ls     "$WT"/backend/internal/handler/dto/request/$FEATURE.go 2>/dev/null
ls     "$WT"/backend/internal/handler/dto/response/$FEATURE.go 2>/dev/null
ls     "$WT"/backend/migrations/*$FEATURE*.sql                2>/dev/null
# Frontend
ls -d  "$WT"/frontend/src/features/$FEATURE                   2>/dev/null
find   "$WT"/frontend/src/app -type d -name "$FEATURE"        2>/dev/null
```

Identify any **external adapter** owned solely by this feature (e.g. billing → `adapter/stripe`, notification → `adapter/resend`, ai → `adapter/claude`). Use Grep to confirm the adapter is referenced only by this feature's service before counting it as removable.

Present the complete file list grouped by layer.

### STEP 4 — Simulate the removal (worktree only)

Delete the feature's folders/files AND strip its wiring. This mirrors `/remove-feature`, but in the throwaway copy only.

Delete files (dependents first):

```bash
rm -f  "$WT"/backend/internal/handler/dto/request/$FEATURE.go
rm -f  "$WT"/backend/internal/handler/dto/response/$FEATURE.go
rm -f  "$WT"/backend/internal/handler/$FEATURE.go
rm -rf "$WT"/backend/internal/app/$FEATURE
rm -f  "$WT"/backend/internal/adapter/postgres/$FEATURE.go
# rm -rf "$WT"/backend/internal/adapter/<provider>   # only if owned solely by this feature
rm -f  "$WT"/backend/internal/port/repository/$FEATURE*.go
rm -f  "$WT"/backend/internal/port/service/$FEATURE*.go
rm -rf "$WT"/backend/internal/domain/$FEATURE
rm -rf "$WT"/frontend/src/features/$FEATURE
# frontend route/demo pages found in STEP 3:
# rm -rf "$WT"/frontend/src/app/<path>/$FEATURE
```

Then remove the wiring lines (use Read + the worktree files; edit in place inside `$WT`):

- **`$WT/backend/cmd/api/main.go`** — remove the repo instantiation (`<f>Repo := postgres.New...`), the service instantiation (`<f>Svc := app<f>.NewService(...)`), the feature's service argument passed to `handler.NewRouter(...)`, and any now-unused import lines.
- **`$WT/backend/internal/handler/router.go`** — remove the feature's route group / registrations, its handler instantiation, the service parameter from the `NewRouter(...)` signature, and related imports.

Migrations: do NOT try to apply/revert anything. Just note which migration files belong to the feature (the down migration is what a real removal would run to drop the tables).

### STEP 5 — Run the full build + tests (inside the worktree)

Capture every result. Run in `$WT`, not the real repo.

```bash
# Backend
cd "$WT/backend" && go build ./...
cd "$WT/backend" && go test ./... -count=1
# Frontend
cd "$WT/frontend" && npx tsc --noEmit
```

Do NOT attempt to "fix" the worktree to make it pass — unlike `/remove-feature`, the point here is to *observe* whether removal is clean. Any error is a finding, not something to patch.

### STEP 6 — Verdict

- **All three pass (backend build + backend tests + frontend tsc)** → the feature is **TRULY INDEPENDENT**. The promise holds: deleting it caused zero breakage elsewhere.
- **Any compile/type/test error in code OUTSIDE the feature** → **COUPLED**. This is a hidden cross-feature coupling that violates the architecture. Report exactly which files / imports broke as actionable findings (these are the "offending edges").

### STEP 7 — ALWAYS clean up

Remove the worktree no matter what the verdict is — never leave it behind. Do this even if a previous step errored.

```bash
cd "$CLAUDE_PROJECT_DIR"
git worktree remove --force "/tmp/ci-independence-$FEATURE"
git worktree prune
```

Confirm to the user that the real repo was never modified.

### STEP 8 — Report

Output a structured report:

1. **Feature** — name, and the external adapter (if any) it owns.
2. **Files that would be removed** — grouped by layer (backend / frontend / migrations) + the 2 wiring files that would be modified.
3. **Build / test result** — backend build, backend tests, frontend `tsc`, each PASS/FAIL with captured output on failure.
4. **Verdict** — `INDEPENDENT` or `COUPLED`, with the offending edges (file → broken import) if coupled.
5. **Fixes** (only if coupled) — concrete remediation, e.g.:
   - Another feature imports this feature's package directly → invert the dependency: define a small interface in `port/` and inject it via constructor in `cmd/api/main.go` instead of importing the concrete package.
   - A frontend feature imports `@/features/<f>/...` → move composition up into the `app/` page; features must not know about each other.
   - A migration in another feature has `REFERENCES <f>_table` → drop the cross-feature FK; only `REFERENCES users(id)` is allowed.
6. **Cleanup confirmation** — worktree removed, real repo untouched.

Example:

```
Feature: billing (owns adapter/stripe)

Would remove (14 files):
  backend/internal/domain/billing/            (3 files)
  backend/internal/port/repository/subscription.go
  backend/internal/port/service/payment.go
  backend/internal/app/billing/               (2 files)
  backend/internal/adapter/postgres/subscription.go
  backend/internal/adapter/stripe/            (2 files)
  backend/internal/handler/billing.go
  backend/internal/handler/dto/request/billing.go
  frontend/src/features/billing/              (5 files)
  frontend/src/app/(dashboard)/billing/page.tsx
Would modify (2 files):
  backend/cmd/api/main.go, backend/internal/handler/router.go
Migrations (note only): 002_create_subscriptions.{up,down}.sql

Simulation results (throwaway worktree):
  backend go build ./...   PASS
  backend go test ./...    PASS (24 passed, 0 failed)
  frontend tsc --noEmit    PASS

VERDICT: INDEPENDENT — billing is fully removable. Promise holds.
Cleanup: worktree /tmp/ci-independence-billing removed. Real repo untouched.
```

Coupled example verdict line:

```
VERDICT: COUPLED — 2 offending edges:
  - backend/internal/app/admin/service.go imports internal/app/billing  → invert via port interface
  - frontend/src/features/admin/components/stats.tsx imports @/features/billing/...  → compose in app/ page
```

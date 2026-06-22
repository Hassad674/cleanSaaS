---
name: test
description: Intelligently run relevant tests based on changed files or a specified feature. Use to run tests for a feature, verify a fix, or run the full test suite (Go unit/integration, Vitest unit, Playwright e2e).
user-invocable: true
allowed-tools: Read, Bash, Grep, Glob, Agent
---

# Test

Target: **$ARGUMENTS**

You are the test runner for CleanSaaS. Determine what to test, run it, analyze results, and fix failures.

**Stack reality (use exactly):**
- Backend: Go `testing` + `testify`. Run with `go test`. Integration tests are gated by `//go:build integration` and the `-tags=integration` flag (need Docker/Postgres).
- Frontend unit: **Vitest** (`npx vitest run`) + Testing Library. Tests are co-located: `*.test.ts` / `*.test.tsx`. **There is no Jest** — never call `jest`.
- Frontend e2e: **Playwright** (`npx playwright test`), specs in `frontend/e2e/`. Needs backend :8081 + frontend :3010 running and the DB seeded.
- Always use relative paths from repo root, or `$CLAUDE_PROJECT_DIR`. Never absolute user paths.

---

## STEP 1 — Determine test scope

Based on `$ARGUMENTS`, determine what to test:

- **empty / "all"** → full suite, backend + frontend.
- **a feature name** (e.g. `auth`, `billing`, `ai`) → tests for that feature across backend + frontend.
- **"backend" / "frontend"** → that side only.
- **"e2e"** → the Playwright suite only (see STEP 5).
- **"changed" / "diff"** → detect changed files and run only relevant tests:
  ```bash
  git diff --name-only HEAD
  git diff --name-only --staged
  ```
  Map each changed file to its test(s): a `*.go` file → its `_test.go` package; a `features/<f>/...` file → that feature's Vitest tests and any e2e flow that touches it.
- **a file path** → the test file closest to that path.

---

## STEP 2 — Discover test files

### Backend (for a feature `{feature}`)
```
backend/internal/domain/{feature}/*_test.go        → Domain tests (unit)
backend/internal/app/{feature}/*_test.go           → Service tests (unit, mocked ports)
backend/internal/adapter/postgres/{feature}_test.go → Integration tests (build tag: integration)
backend/internal/handler/{feature}_test.go         → Handler tests
backend/pkg/**/*_test.go                            → Utility tests
```

### Frontend (for a feature `{feature}`)
```
frontend/src/features/{feature}/**/*.test.ts        → Unit tests (Vitest)
frontend/src/features/{feature}/**/*.test.tsx       → Component tests (Vitest + Testing Library)
frontend/e2e/*.spec.ts                              → End-to-end flows (Playwright)
```

List which test files exist and which are **missing** (see "Missing tests" in the report).

---

## STEP 3 — Run backend tests

### Unit tests for a feature:
```bash
cd backend
go test ./internal/domain/{feature}/... ./internal/app/{feature}/... -count=1 -v
```

### All backend unit tests:
```bash
cd backend
go test ./internal/domain/... ./internal/app/... ./pkg/... -count=1
```

### Integration tests (need Docker + Postgres on :5433):
```bash
cd backend
go test ./internal/adapter/... -tags=integration -count=1
```

### Full backend (everything):
```bash
cd backend
go test ./... -count=1
```

**Useful flags:** `-count=1` (disable cache), `-race` (race detector — always in CI, optional locally), `-timeout 60s`, `-run '<TestName>'` to run a single test.

---

## STEP 4 — Run frontend unit tests (Vitest)

### Specific feature:
```bash
cd frontend
npx vitest run src/features/{feature}
```

### A single test by name:
```bash
cd frontend
npx vitest run -t "<test name>"
```

### All frontend unit tests:
```bash
cd frontend
npx vitest run
```

### Type checking (always run alongside tests — a passing suite that doesn't typecheck is still broken):
```bash
cd frontend
npx tsc --noEmit
```

---

## STEP 5 — Run end-to-end tests (Playwright)

Only when `$ARGUMENTS` is "e2e"/"all", or the change affects a user flow. E2E needs the stack **running** (backend :8081, frontend :3010) and the DB **seeded** (so login works). If it isn't up, run the `/run` skill first, or delegate to the `/e2e` skill which handles all the preflight.

```bash
cd frontend
npx playwright install --with-deps chromium   # first run only
npx playwright test                            # whole suite
npx playwright test e2e/<file>.spec.ts         # one spec
npx playwright show-report                     # open the last HTML report
```

> Known config gotchas to check/fix if e2e fails immediately: the Playwright `baseURL`/port must match the dev server (**:3010**), and specs must log in with the seeded account **admin@cleansaas.dev / admin123**. A `:3006` baseURL or a `@cleansaas.com` login is a config bug, not a product bug.

---

## STEP 6 — Analyze results

For each run, categorize:

### PASSED — count + summary.
### FAILED — for each: **test name**, **file:line**, **actual vs expected / panic message**, and **root-cause analysis** (read BOTH the test and the source to explain why).
### SKIPPED — list and why (build tags, missing deps).
### MISSING — files that SHOULD have tests but don't:
- every `entity.go` in domain/ → `entity_test.go`
- every `service.go` in app/ → `service_test.go`
- key frontend components/hooks → a `.test.tsx`/`.test.ts`

---

## STEP 7 — Fix failures (if "fix" in $ARGUMENTS or the cause is obvious)

### Go failures:
1. Read the failing test and the source.
2. Decide whether the bug is in the test or the implementation (prefer fixing the implementation, unless the expectation is genuinely wrong).
3. Fix → re-run just that test (`-run '<TestName>'`) → then re-run the feature suite for regressions.

### TypeScript / React (Vitest) failures:
1. Read the test and the component/hook.
2. Check type mismatches, missing props, wrong assertions, unmocked fetch/server-action.
3. Fix → re-run (`-t "<name>"`).

Follow `CLAUDE.md`'s **Test → Fix → Retest loop** (max 3 attempts per failing test) and **Blocker policy** (stuck on the same error after 3 genuinely different approaches → write `BLOCKED-<topic>.md`).

### Never:
- Skip or delete a failing test to make the suite pass.
- Change assertions to match buggy behavior.
- Add `t.Skip()` / `.skip()` without a clear, stated reason.

---

## STEP 8 — Compilation check (always, even if tests pass)

```bash
cd backend  && go build ./...
cd frontend && npx tsc --noEmit
```

A passing test suite with compilation errors means something is broken.

---

## Report format

```
# Test Report — {target}

## Summary
- Backend:  X passed, Y failed, Z skipped   (unit / integration)
- Frontend: X passed, Y failed, Z skipped   (Vitest)
- E2E:      X passed, Y failed              (Playwright, if run)
- Compilation: backend OK / frontend OK

## Failures
### [FAIL] TestAuthService_Register_DuplicateEmail
- File: backend/internal/app/auth/service_test.go:45
- Error: expected ErrAlreadyExists, got ErrInternal
- Cause: FindByEmail returns the wrong error when the user exists
- Fix: (applied / suggested)

## Missing tests
- backend/internal/domain/billing/invoice_test.go — NO TESTS
- frontend/src/features/team/... — feature layer has no unit tests

## Recommendations
- <highest-value test to add next>
```

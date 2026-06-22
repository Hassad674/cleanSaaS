---
name: check
description: Verify feature independence, architecture compliance, and code quality. Run to detect cross-feature imports, forbidden dependencies, missing tests, and modularity violations.
user-invocable: true
allowed-tools: Read, Bash, Grep, Glob, Agent
---

# Check — Architecture & Independence Verification

Target: **$ARGUMENTS**

If `$ARGUMENTS` is empty, check the ENTIRE project. Otherwise, check only the specified feature(s).

You are the architecture guardian for CleanSaaS. Run every check below and produce a clear pass/fail report.

---

## CHECK 1 — Backend: Dependency direction

The hexagonal rule is absolute:
```
handler → app → domain ← port ← adapter
```

### 1a. Domain purity
Verify that files in `backend/internal/domain/` import ONLY:
- Go standard library packages
- Other domain sub-packages (`domain/`, `domain/user`, etc.)

**How to check:**
```bash
# Find any non-stdlib, non-domain imports in domain/
cd .
```
Use Grep to search for import statements in `backend/internal/domain/` and verify none reference `internal/port`, `internal/app`, `internal/adapter`, `internal/handler`, `internal/config`, or any external module.

**FAIL if:** Domain imports anything from port/, app/, adapter/, handler/, or external packages.

### 1b. Port purity
Verify that `backend/internal/port/` imports ONLY domain/ packages and Go stdlib.

**FAIL if:** Port imports app/, adapter/, handler/, or external packages.

### 1c. App layer
Verify that `backend/internal/app/` imports ONLY:
- domain/ packages
- port/ interfaces
- pkg/ utilities
- Go stdlib

**FAIL if:** App imports adapter/ or handler/ directly.

### 1d. No adapter cross-imports
Verify that no adapter imports another adapter.

**FAIL if:** `adapter/postgres` imports `adapter/stripe`, etc.

### 1e. Handler imports
Verify handler/ imports only from app/, dto/, middleware/, and pkg/. Never from adapter/ or domain/ directly (except for response mapping which may reference domain types for DTO conversion).

---

## CHECK 2 — Frontend: Feature isolation

### 2a. No cross-feature imports
Verify that no file in `frontend/src/features/{X}/` imports from `frontend/src/features/{Y}/`.

**How to check:**
Use Grep to search all files in `frontend/src/features/` for import patterns like `@/features/`. For each match, verify the import stays within the SAME feature folder.

**FAIL if:** `features/billing/components/foo.tsx` imports from `@/features/user/...`

### 2b. Feature imports only from allowed sources
Each feature file can ONLY import from:
- Its own feature directory (`../types`, `./sub-component`, `../hooks/...`)
- `@/shared/` (components, hooks, lib, types)
- External npm packages (react, next, etc.)

**FAIL if:** A feature imports from `@/app/`, `@/config/`, or another feature.

### 2c. Page thinness
Verify that files in `frontend/src/app/` are thin routing layers:
- No `useState`, `useEffect`, or event handlers in page files
- Pages should be under 20 lines (warn if over 30)
- Pages import from `@/features/` and compose, nothing else

---

## CHECK 3 — Database: No cross-feature foreign keys

### 3a. Migration FK audit
Read all SQL migration files in `backend/migrations/` (and `db/init.sql` if present).
List every `REFERENCES` clause. Verify:
- FKs to `users(id)` are OK (users is core)
- FKs to any other feature's table are FORBIDDEN

**FAIL if:** `conversations` table references `subscriptions`, or `notifications` references `files`, etc.

### 3b. Migration pairs
Every `.up.sql` must have a matching `.down.sql`. Check for unpaired migrations.

---

## CHECK 4 — Backend: Wiring in main.go

### 4a. Explicit DI
Read `backend/cmd/api/main.go`. Verify:
- All dependencies are wired explicitly (no auto-discovery, no reflection)
- Each feature's repo → service → handler chain is visible
- No global variables or init() functions for wiring

### 4b. Feature removability
For each feature wired in main.go, verify that commenting out its lines (repo + service + router registration) would NOT break compilation of other features.

---

## CHECK 5 — Code quality

### 5a. No `any` in TypeScript
Search frontend for `any` type usage (excluding node_modules, .next).

**FAIL if:** `any` is used outside of legitimate escape hatches (with a comment explaining why).

### 5b. No `as` type casts in TypeScript
Search for unsafe type assertions. Warn on each occurrence.

### 5c. No raw SQL string concatenation in Go
Search backend for string concatenation in SQL queries (indicators: `fmt.Sprintf` with SQL keywords, `+` concatenation near query strings).

**FAIL if:** Any query uses string interpolation instead of `$1, $2` parameters.

### 5d. Context usage
Verify all repository and service methods accept `context.Context` as first parameter.

### 5e. Error handling
Search for swallowed errors: `_ = someFunc()` patterns where the function returns an error.

---

## CHECK 6 — Test coverage

### 6a. Domain tests exist
For each entity in `domain/`, verify a `_test.go` file exists.

### 6b. Service tests exist
For each service in `app/`, verify a `_test.go` file exists.

### 6c. Frontend component tests (warning only)
Check if key components have test files. Warn if missing but don't fail.

---

## CHECK 7 — Independence simulation

For each non-core feature (everything except auth/user), simulate removal:

1. List all files that belong to this feature (backend + frontend)
2. Grep the entire codebase for imports of those paths
3. Any import OUTSIDE the feature itself = **FAIL**

Core features (auth, user) are exempt — they can be depended upon.

---

## Report format

Output a structured report:

```
# CleanSaaS Architecture Check Report

## Summary
- Total checks: X
- Passed: X
- Failed: X
- Warnings: X

## Results

### CHECK 1 — Backend dependency direction
- [PASS] 1a. Domain purity
- [PASS] 1b. Port purity
- [FAIL] 1c. App layer — app/billing imports adapter/stripe directly
  → Fix: inject via port/service/PaymentService interface
...

### CHECK 7 — Independence simulation
- [PASS] billing — removable, 0 external references
- [FAIL] notification — referenced by features/admin/components/stats.tsx
  → Fix: remove cross-feature import, pass data via page composition
```

For each failure, provide:
1. Exact file and line
2. What rule it violates
3. How to fix it

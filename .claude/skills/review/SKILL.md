---
name: review
description: Review code for security, performance, architecture compliance, and quality. Use before committing or to audit existing code. Complementary to /check which verifies structure — /review verifies content.
user-invocable: true
allowed-tools: Read, Bash, Grep, Glob, Agent
---

# Review

Review target: **$ARGUMENTS**

If `$ARGUMENTS` is empty, review all uncommitted changes (staged + unstaged). Otherwise, review the specified feature, file, or commit range.

You are the code reviewer for CleanSaaS. Check every aspect below and produce an actionable report.

---

## STEP 0 — Determine review scope

### If `$ARGUMENTS` is empty — review uncommitted changes:
```bash
cd /home/hassad/Documents/boilerplateSaaS
git diff --name-only HEAD
git diff --name-only --staged
```
Read every changed file.

### If `$ARGUMENTS` is a feature name (e.g., "auth", "billing"):
Read all files in:
- `backend/internal/domain/{feature}/`
- `backend/internal/app/{feature}/`
- `backend/internal/adapter/postgres/{feature}.go`
- `backend/internal/handler/{feature}.go`
- `backend/internal/handler/dto/request/{feature}.go`
- `frontend/src/features/{feature}/`

### If `$ARGUMENTS` is a file path:
Read that specific file and its related files (test file, interface, etc.).

### If `$ARGUMENTS` is a commit range (e.g., "HEAD~3..HEAD"):
```bash
git diff $ARGUMENTS --name-only
```
Read all files in the diff.

---

## REVIEW 1 — Security

### 1a. SQL injection
Search for string concatenation in SQL queries:
- `fmt.Sprintf` near SQL keywords (`SELECT`, `INSERT`, `UPDATE`, `DELETE`, `WHERE`)
- String `+` concatenation building queries
- Template literals in queries

**CRITICAL if found.** All queries must use `$1, $2, $3` parameterized placeholders.

### 1b. Input validation
For every handler that accepts user input:
- Is the request body decoded and validated before use?
- Are required fields checked?
- Are string lengths bounded? (prevent memory abuse)
- Are IDs validated as UUID format? (prevent path traversal)

### 1c. Authentication/Authorization
- Are protected endpoints behind `middleware.Auth()`?
- Do endpoints that modify user data verify ownership? (user can only edit their own data)
- Are admin-only endpoints checking `user.IsAdmin()`?

### 1d. Secrets exposure
- No hardcoded API keys, passwords, tokens in source code
- No secrets in error messages returned to clients
- Sensitive fields (password_hash, tokens) not included in API responses

### 1e. XSS / injection in frontend
- No `dangerouslySetInnerHTML` without sanitization
- User-generated content properly escaped
- No `eval()` or `new Function()` with user input

---

## REVIEW 2 — Performance

### 2a. N+1 queries
Look for patterns where a list is fetched, then each item triggers another query:
```go
// BAD: N+1
users, _ := repo.List(ctx, 0, 100)
for _, u := range users {
    subs, _ := subRepo.FindByUserID(ctx, u.ID)  // 100 extra queries!
}
```
Suggest JOINs or batch queries instead.

### 2b. Missing database indexes
For every `WHERE` clause and `JOIN` condition in adapter SQL, verify a corresponding index exists in migrations.

### 2c. Unbounded queries
- Every `SELECT` that returns multiple rows must have `LIMIT`
- Pagination must be enforced (no "fetch all" endpoints)
- Large text fields should not be included in list queries if not needed

### 2d. Missing context timeouts
- All database calls should use `ctx` (already enforced by interface, but verify)
- Long-running operations should have timeout context

### 2e. Frontend performance
- Large components should be lazy-loaded
- No unnecessary `"use client"` — check if Server Component would work
- No `useEffect` fetching data that could be a Server Action

---

## REVIEW 3 — Architecture compliance

### 3a. Layer violations
- Domain imports nothing external
- App layer uses interfaces, not concrete types
- Handler contains no business logic
- Adapters don't import each other

### 3b. Feature isolation
- No cross-feature imports (backend or frontend)
- Composition happens in app/ pages only (frontend)
- DI in main.go only (backend)

### 3c. Error handling chain
- Domain defines errors → app returns them → handler maps to HTTP
- No HTTP status codes in app layer
- No domain errors swallowed silently
- All errors from external calls are wrapped with context: `fmt.Errorf("doing X: %w", err)`

### 3d. DTO discipline
- Handlers never pass domain entities directly to JSON encoder
- Response DTOs strip sensitive fields (password_hash, internal IDs where inappropriate)
- Request DTOs only contain fields needed for that operation

---

## REVIEW 4 — Code quality

### 4a. TypeScript strict mode
- No `any` type (suggest the correct type)
- No unsafe `as` casts (suggest type guards or proper typing)
- No `// @ts-ignore` or `// @ts-expect-error` without justification

### 4b. Go idioms
- Error handling: `if err != nil` immediately after the call, not deferred
- No naked returns in functions with named return values
- No `panic()` in library code (only in main() for fatal startup errors)
- Exported functions have a comment (Go convention)

### 4c. Naming
- Go: PascalCase exports, camelCase locals, ALL_CAPS only for constants
- TypeScript: PascalCase components, camelCase functions/variables, kebab-case files
- Database: snake_case tables and columns
- Clear, descriptive names — no `x`, `tmp`, `data` (unless obvious from context)

### 4d. Dead code
- Unused variables, imports, functions
- Commented-out code blocks
- Unreachable code after return/panic

### 4e. Duplication
- Same logic repeated 3+ times → should be extracted
- Similar handler patterns that could share a helper
- Identical SQL queries in different methods

---

## REVIEW 5 — Testing

### 5a. Test quality
- Tests actually assert meaningful behavior (not just "no panic")
- Edge cases covered: empty input, nil, not found, duplicate, unauthorized
- Tests are independent — no shared mutable state between tests
- Table-driven tests for multiple scenarios

### 5b. Test coverage gaps
- New service methods without tests
- New handlers without at least happy-path tests
- Domain validation rules without tests
- Error paths not tested

### 5c. Test naming
- Go: `TestServiceName_MethodName_Scenario`
- TypeScript: `describe('ComponentName')` → `it('should do X when Y')`

---

## Report format

```
# Code Review — {target}

## Summary
- Critical: X issues
- Warning: X issues
- Info: X suggestions
- Clean: X checks passed

## Critical

### [CRITICAL] SQL injection in adapter/postgres/billing.go:42
```go
query := fmt.Sprintf("SELECT * FROM invoices WHERE status = '%s'", status)
```
**Fix:** Use parameterized query: `WHERE status = $1`, status

## Warnings

### [WARN] Missing index for subscriptions.plan_id
The query at adapter/postgres/subscription.go:28 filters by `plan_id` but no index exists.
**Fix:** Add migration: `CREATE INDEX idx_subscriptions_plan_id ON subscriptions(plan_id);`

### [WARN] N+1 query in handler/admin.go:15
Fetches users list then loops to get subscription for each.
**Fix:** Use a JOIN query or batch fetch.

## Info

### [INFO] Unused import in app/billing/service.go:4
`"fmt"` is imported but not used.

## Clean

- [OK] No hardcoded secrets
- [OK] All handlers validate input
- [OK] Domain layer has zero external imports
- [OK] Response DTOs strip sensitive fields
```

### Severity levels:
- **CRITICAL** — Security vulnerability or data loss risk. Must fix before merge.
- **WARNING** — Performance issue, missing test, or architecture smell. Should fix.
- **INFO** — Style, naming, or minor improvement. Nice to fix.
- **CLEAN** — Explicitly confirm checks that passed (builds confidence).

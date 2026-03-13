---
name: test
description: Intelligently run relevant tests based on changed files or a specified feature. Use to run tests for a feature, verify a fix, or run the full test suite.
user-invocable: true
allowed-tools: Read, Bash, Grep, Glob, Agent
---

# Test

Target: **$ARGUMENTS**

You are the test runner for CleanSaaS. Determine what to test, run it, analyze results, and fix failures.

---

## STEP 1 — Determine test scope

Based on `$ARGUMENTS`, determine what to test:

### If `$ARGUMENTS` is empty or "all":
Run the full test suite for both backend and frontend.

### If `$ARGUMENTS` is a feature name (e.g., "auth", "billing", "ai"):
Run tests for that specific feature across backend and frontend.

### If `$ARGUMENTS` is "backend" or "frontend":
Run the full test suite for that side only.

### If `$ARGUMENTS` is "changed" or "diff":
Use git to detect changed files and run only relevant tests:
```bash
git diff --name-only HEAD
git diff --name-only --staged
```
Map changed files to their test files and related test suites.

### If `$ARGUMENTS` is a file path:
Run the test file closest to that path.

---

## STEP 2 — Discover test files

### Backend test discovery

For a feature named `{feature}`:
```
backend/internal/domain/{feature}/*_test.go    → Domain tests
backend/internal/app/{feature}/*_test.go       → Service tests (unit)
backend/internal/adapter/postgres/{feature}_test.go → Integration tests
backend/internal/handler/{feature}_test.go     → Handler tests
backend/pkg/**/*_test.go                       → Utility tests
```

List which test files exist and which are missing.

### Frontend test discovery

For a feature named `{feature}`:
```
frontend/src/features/{feature}/**/*.test.ts    → Unit tests
frontend/src/features/{feature}/**/*.test.tsx   → Component tests
```

List which test files exist and which are missing.

---

## STEP 3 — Run backend tests

### Unit tests (for a specific feature):
```bash
cd /home/hassad/Documents/boilerplateSaaS/backend
go test ./internal/domain/{feature}/... -v -count=1
go test ./internal/app/{feature}/... -v -count=1
```

### All backend unit tests:
```bash
cd /home/hassad/Documents/boilerplateSaaS/backend
go test ./internal/domain/... ./internal/app/... ./pkg/... -v -count=1
```

### Integration tests (if Docker is running):
```bash
cd /home/hassad/Documents/boilerplateSaaS/backend
go test ./internal/adapter/... -v -count=1 -tags=integration
```

### Full backend:
```bash
cd /home/hassad/Documents/boilerplateSaaS/backend
go test ./... -v -count=1
```

**Important flags:**
- `-v` for verbose output
- `-count=1` to disable test caching
- `-race` to detect race conditions (add for CI, skip for quick local runs)
- `-timeout 30s` to prevent hanging tests

---

## STEP 4 — Run frontend tests

### Specific feature:
```bash
cd /home/hassad/Documents/boilerplateSaaS/frontend
npx jest --testPathPattern="features/{feature}" --verbose 2>/dev/null || npm test -- --testPathPattern="features/{feature}" --verbose
```

### All frontend:
```bash
cd /home/hassad/Documents/boilerplateSaaS/frontend
npm test -- --verbose
```

### Type checking (always run alongside tests):
```bash
cd /home/hassad/Documents/boilerplateSaaS/frontend
npx tsc --noEmit
```

---

## STEP 5 — Analyze results

For each test run, parse the output and categorize:

### PASSED tests
List count and summary.

### FAILED tests
For each failure:
1. **Test name** — Full test function/describe name
2. **File location** — Exact file and line
3. **Error message** — The actual vs expected or panic message
4. **Root cause analysis** — Read the test file AND the source file to understand why it failed

### SKIPPED tests
List and note why (missing dependencies, build tags, etc.)

### MISSING tests
List files that SHOULD have tests but don't:
- Every `entity.go` in domain/ needs `entity_test.go`
- Every `service.go` in app/ needs `service_test.go`
- Key frontend components should have `.test.tsx`

---

## STEP 6 — Fix failures (if requested or obvious)

If the user included "fix" in `$ARGUMENTS`, or if failures are due to obvious issues:

### For Go test failures:
1. Read the failing test and the source code
2. Determine if the bug is in the test or the implementation
3. Fix the actual bug (prefer fixing implementation over adjusting tests, unless the test expectation is wrong)
4. Re-run the specific test to verify the fix
5. Run the full feature test suite to check for regressions

### For TypeScript/React test failures:
1. Read the test file and component
2. Check for type mismatches, missing props, incorrect assertions
3. Fix and re-run

### Do NOT:
- Skip or delete failing tests to make the suite pass
- Change test assertions to match buggy behavior
- Add `t.Skip()` or `.skip()` without a clear reason

---

## STEP 7 — Compilation check

Always verify compilation even if tests pass:

```bash
cd /home/hassad/Documents/boilerplateSaaS/backend && go build ./...
cd /home/hassad/Documents/boilerplateSaaS/frontend && npx tsc --noEmit
```

A passing test suite with compilation errors means something is broken.

---

## Report format

```
# Test Report — {target}

## Summary
- Backend: X passed, Y failed, Z skipped
- Frontend: X passed, Y failed, Z skipped
- Compilation: backend OK / frontend OK

## Failures

### [FAIL] TestAuthService_Register_DuplicateEmail
- File: backend/internal/app/auth/service_test.go:45
- Error: expected ErrAlreadyExists, got ErrInternal
- Cause: FindByEmail returns wrong error when user exists
- Fix: (applied / suggested)

## Missing tests

- backend/internal/domain/billing/entity_test.go — NO TESTS
- backend/internal/app/ai/service_test.go — NO TESTS

## Recommendations
- Add domain validation tests for billing entity
- Add service-level tests for AI feature with mocked repository
```

---
name: remove-feature
description: Safely remove a feature from the project (backend + frontend + migrations). Use to prove modularity works or when a user doesn't need a module.
user-invocable: true
allowed-tools: Read, Edit, Bash, Grep, Glob, Agent
---

# Remove Feature

Remove the feature: **$ARGUMENTS**

You are safely removing a feature from CleanSaaS. This operation proves the modularity promise — every non-core feature must be deletable without breaking anything.

---

## GUARD — Core features cannot be removed

The following features are **core** and cannot be removed:
- **auth** — Everything depends on authentication
- **user** — Users are the foundation of every feature

If `$ARGUMENTS` matches a core feature, STOP and explain why it cannot be removed.

---

## STEP 1 — Inventory the feature

Map every file and reference belonging to this feature.

### Backend files
```
backend/internal/domain/{feature}/          → Domain entities
backend/internal/port/repository/{feature}.go → Repository interface
backend/internal/port/service/{feature}.go  → Service interface (if any)
backend/internal/app/{feature}/             → Application service
backend/internal/adapter/postgres/{feature}.go → PostgreSQL adapter
backend/internal/adapter/{provider}/        → External adapters (stripe, resend, etc.)
backend/internal/handler/{feature}.go       → HTTP handler
backend/internal/handler/dto/request/{feature}.go → Request DTOs
backend/internal/handler/dto/response/{feature}.go → Response DTOs (if separate file)
```

### Frontend files
```
frontend/src/features/{feature}/            → Entire feature directory
frontend/src/app/**/{feature}/              → Route pages for this feature
```

### Migrations
```
backend/migrations/*_{feature}*.up.sql
backend/migrations/*_{feature}*.down.sql
```

### Wiring references
```
backend/cmd/api/main.go                    → DI lines for this feature
backend/internal/handler/router.go          → Route registrations
```

Use Glob and Grep to find ALL files. Do not guess — be exhaustive.

---

## STEP 2 — Check for external dependencies

Before deleting anything, verify no other feature depends on this one.

### Backend: search for imports of this feature
```bash
# Search for any import of this feature's packages outside its own directory
```
Use Grep to find imports matching `domain/{feature}`, `app/{feature}`, `port/repository/{feature}`, `adapter/*/{feature}` in files OUTSIDE the feature's own directories.

**Allowed references:**
- `cmd/api/main.go` — wiring (will be cleaned up)
- `handler/router.go` — routes (will be cleaned up)
- `handler/{feature}.go` — the feature's own handler

**Forbidden references:**
- Another feature's app service importing this feature
- Another feature's handler importing this feature's types
- Another adapter importing this feature's adapter

If forbidden references are found, **STOP** and report them. The dependency must be fixed before removal is possible.

### Frontend: search for cross-feature imports
Use Grep to search all files in `frontend/src/features/` for imports from `@/features/{feature}/`. Any match in a DIFFERENT feature = **STOP**.

### Database: check for cross-feature foreign keys
Read all migration files. If any OTHER feature's table has a `REFERENCES {feature_table}`, **STOP** and report.

---

## STEP 3 — Clean wiring (before deletion)

### 3a. Update `backend/cmd/api/main.go`
Remove:
- Repository instantiation: `{feature}Repo := postgres.New{Feature}Repository(db)`
- Service instantiation: `{feature}Svc := app{feature}.NewService(...)`
- Any related import lines
- Update `handler.NewRouter(...)` call to remove the feature's service parameter

### 3b. Update `backend/internal/handler/router.go`
Remove:
- The feature's route group / route registrations
- The feature's handler instantiation inside the router
- The service parameter from `NewRouter()` function signature
- Related import lines

### 3c. Update response DTOs
If `handler/dto/response/response.go` contains DTOs specific to this feature, remove them.
If the feature had its own `handler/dto/response/{feature}.go`, it will be deleted with the rest.

---

## STEP 4 — Delete feature files

Delete in this order (dependencies first, dependents last):

### Backend
1. `backend/internal/handler/dto/request/{feature}.go`
2. `backend/internal/handler/dto/response/{feature}.go` (if separate)
3. `backend/internal/handler/{feature}.go`
4. `backend/internal/app/{feature}/` (entire directory)
5. `backend/internal/adapter/postgres/{feature}.go`
6. `backend/internal/adapter/{provider}/` (feature-specific adapters)
7. `backend/internal/port/repository/{feature}.go`
8. `backend/internal/port/service/{feature}.go` (if exists)
9. `backend/internal/domain/{feature}/` (entire directory)

### Frontend
10. `frontend/src/features/{feature}/` (entire directory)
11. `frontend/src/app/**/{feature}/` (route pages)

### Migrations
12. Do NOT delete migration files. Instead, note which migrations belong to this feature. The user can run the down migration to drop the tables:
    ```
    To drop this feature's tables: apply the down migration for migrations/NNN_{feature}.down.sql
    ```

---

## STEP 5 — Compile and verify

### Backend
```bash
cd ./backend && go build ./...
```

Fix any compilation errors. Common issues:
- Unused imports in main.go or router.go — remove them
- Missing parameter in NewRouter() — update the signature

### Frontend
```bash
cd ./frontend && npx tsc --noEmit
```

Fix any TypeScript errors. Common issues:
- Page importing from deleted feature — delete the page too
- Shared component referencing feature type — this shouldn't happen if architecture is clean

---

## STEP 6 — Run tests

```bash
cd ./backend && go test ./... -count=1
cd ./frontend && npm test 2>/dev/null || true
```

All remaining tests must pass. If a test outside this feature fails, there was a hidden dependency — fix it.

---

## STEP 7 — Final verification

Run a mental `/check` on the result:
- [ ] Backend compiles with zero errors
- [ ] Frontend compiles with zero errors
- [ ] No dangling imports referencing the removed feature
- [ ] No orphan route pages pointing to deleted components
- [ ] main.go and router.go are clean
- [ ] All remaining tests pass

---

## Output

Report:
1. **Files deleted** (grouped by layer)
2. **Files modified** (main.go, router.go, etc.)
3. **Migration note** (which down migration to run for table cleanup)
4. **Dependency issues found** (if any, and how they were resolved)
5. **Compilation result** — backend + frontend
6. **Test result** — all pass / failures

Example:
```
Removed feature: billing

Deleted (14 files):
  backend/internal/domain/billing/         (3 files)
  backend/internal/port/repository/subscription.go
  backend/internal/port/service/payment.go
  backend/internal/app/billing/            (2 files)
  backend/internal/adapter/postgres/subscription.go
  backend/internal/adapter/stripe/         (2 files)
  backend/internal/handler/billing.go
  backend/internal/handler/dto/request/billing.go
  frontend/src/features/billing/           (6 files)
  frontend/src/app/(dashboard)/billing/page.tsx

Modified (2 files):
  backend/cmd/api/main.go — removed billing wiring
  backend/internal/handler/router.go — removed billing routes

Migrations to revert:
  make migrate-down  (reverts 002_create_subscriptions)

Compilation: backend OK, frontend OK
Tests: 24 passed, 0 failed
```

---
name: add-endpoint
description: Add a new API endpoint (use case) to an existing feature. Lighter than /add-feature — scaffolds handler, DTO, service method, and test for a single operation.
user-invocable: true
allowed-tools: Read, Write, Edit, Bash, Grep, Glob
---

# Add Endpoint

Add endpoint: **$ARGUMENTS**

You are adding a new endpoint to an existing feature in CleanSaaS. This is lighter than `/add-feature` — the domain, ports, and adapter already exist. You're adding a new use case.

---

## STEP 0 — Parse the request

Determine from `$ARGUMENTS`:
- **Feature name** — which existing feature (auth, user, billing, ai, storage, notification, admin)
- **Operation** — what the endpoint does (reset-password, export-data, archive-conversation, etc.)
- **HTTP method + path** — derive from the operation (POST /auth/reset-password, GET /users/export, etc.)
- **Public or protected** — does it require authentication?

---

## STEP 1 — Read existing feature code

Before writing anything, read the current state of the feature:

1. **App service** — `backend/internal/app/{feature}/service.go` — understand existing methods and dependencies
2. **Handler** — `backend/internal/handler/{feature}.go` — understand existing handler structure
3. **Router** — `backend/internal/handler/router.go` — see where routes are registered
4. **Repository interface** — `backend/internal/port/repository/{feature}.go` — check if new data methods are needed
5. **Domain entity** — `backend/internal/domain/{feature}/entity.go` — check if new business methods are needed

---

## STEP 2 — Domain changes (if needed)

If the new use case requires new business logic on the entity, add methods to the domain:

```go
// In domain/{feature}/entity.go
func (e *{Entity}) Archive() {
    e.Archived = true
    e.UpdatedAt = time.Now()
}
```

Add tests for any new domain method in `entity_test.go`.

If the new use case requires a new field on the entity, add it and create a migration with `/add-migration`.

**Skip this step if the existing domain is sufficient.**

---

## STEP 3 — Repository changes (if needed)

If the use case needs a new data access method:

### 3a. Add to port interface
```go
// In port/repository/{feature}.go
type {Feature}Repository interface {
    // ... existing methods
    FindByToken(ctx context.Context, token string) (*{feature}.{Entity}, error) // NEW
}
```

### 3b. Implement in adapter
```go
// In adapter/postgres/{feature}.go
func (r *{Feature}Repository) FindByToken(ctx context.Context, token string) (*{feature}.{Entity}, error) {
    e := &{feature}.{Entity}{}
    err := r.db.QueryRowContext(ctx,
        `SELECT /* columns */ FROM {table} WHERE token = $1`, token,
    ).Scan(/* fields */)
    if err == sql.ErrNoRows {
        return nil, domain.ErrNotFound
    }
    return e, err
}
```

**Skip this step if existing repository methods are sufficient.**

---

## STEP 4 — Service method

Add the new use case to `backend/internal/app/{feature}/service.go`:

**Pattern** (reference existing methods in the file):
```go
func (s *Service) ResetPassword(ctx context.Context, email string) error {
    user, err := s.users.FindByEmail(ctx, email)
    if err != nil {
        return err
    }
    // Business logic here
    // Call port services if needed (email, payment, etc.)
    return nil
}
```

Rules:
- Accept primitive types or domain types as parameters
- Return domain types and domain errors
- No HTTP concepts (no status codes, no request/response structs)
- Use the existing injected dependencies (repository + services)

---

## STEP 5 — Service test

Add test for the new method in `backend/internal/app/{feature}/service_test.go`:

```go
func TestService_ResetPassword_Success(t *testing.T) {
    // Setup mocks
    // Call service method
    // Assert result
}

func TestService_ResetPassword_UserNotFound(t *testing.T) {
    // Setup mocks to return ErrNotFound
    // Call service method
    // Assert ErrNotFound returned
}
```

Use table-driven tests when there are multiple scenarios. Name: `TestServiceName_MethodName_Scenario`.

---

## STEP 6 — Request DTO

Add or create the request struct in `backend/internal/handler/dto/request/{feature}.go`:

```go
type ResetPasswordRequest struct {
    Email string `json:"email"`
}
```

Only include fields that come from the HTTP request body or query params.

---

## STEP 7 — Response DTO (if needed)

If the endpoint returns feature-specific data, add a response DTO or helper.

For simple success/error responses, use the existing `response.JSON()` and `response.Error()`.

---

## STEP 8 — Handler method

Add the handler method to `backend/internal/handler/{feature}.go`:

**Pattern** (reference existing handlers in the file):
```go
func (h *{Feature}Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
    var req request.ResetPasswordRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "invalid request body")
        return
    }

    if req.Email == "" {
        response.Error(w, http.StatusBadRequest, "email is required")
        return
    }

    if err := h.svc.ResetPassword(r.Context(), req.Email); err != nil {
        response.HandleDomainError(w, err)
        return
    }

    response.JSON(w, http.StatusOK, map[string]string{"message": "reset email sent"})
}
```

Rules:
- Decode request → validate required fields → call service → encode response
- NO business logic in the handler
- Use `r.Context()` to pass context
- Use `chi.URLParam(r, "id")` for path parameters

---

## STEP 9 — Register route

Add the route in `backend/internal/handler/router.go`:

```go
// In the appropriate group (public or protected):
r.Post("/auth/reset-password", authHandler.ResetPassword)
```

- Public endpoints go in the public route group (no middleware)
- Protected endpoints go inside `r.Group` with `middleware.Auth(jwtSecret)`
- Use RESTful conventions: `GET /resources`, `POST /resources`, `GET /resources/{id}`, etc.

---

## STEP 10 — Frontend (if needed)

If the endpoint needs a frontend action:

### Server Action
Add to `frontend/src/features/{feature}/actions/{feature}.ts`:
```ts
export async function resetPassword(email: string) {
  return api<void>("/auth/reset-password", {
    method: "POST",
    body: { email },
  });
}
```

### Component updates
If a UI component needs to call this action, update it or create a new one in the feature's `components/` directory.

---

## STEP 11 — Verify

```bash
cd ./backend && go build ./...
cd ./backend && go test ./internal/app/{feature}/... -v -count=1
```

### Checklist:
- [ ] Service method has no HTTP concepts
- [ ] Handler is thin — decode, call, encode
- [ ] New repository method (if any) uses parameterized SQL
- [ ] Test covers success + error cases
- [ ] Route is in the correct group (public/protected)
- [ ] No cross-feature imports introduced

---

## Output

Report:
1. **Endpoint** — `METHOD /path` (public/protected)
2. **Files created** — new files
3. **Files modified** — existing files with changes
4. **Test coverage** — tests written
5. **Frontend** — action added (if applicable)

Example:
```
Added endpoint: POST /auth/reset-password (public)

Modified:
  backend/internal/app/auth/service.go — added ResetPassword()
  backend/internal/app/auth/service_test.go — added 2 tests
  backend/internal/handler/auth.go — added ResetPassword handler
  backend/internal/handler/dto/request/auth.go — added ResetPasswordRequest
  backend/internal/handler/router.go — registered POST /auth/reset-password

Frontend:
  frontend/src/features/auth/actions/auth.ts — added resetPassword()

Tests: 2 new (Success, UserNotFound)
```

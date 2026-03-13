# Backend — Go + Chi + Hexagonal Architecture

## Structure

```
backend/
├── cmd/
│   ├── api/main.go        → Entry point. ALL dependency injection happens here.
│   ├── migrate/main.go    → Run SQL migrations
│   └── seed/main.go       → Seed initial data (plans, admin user)
│
├── internal/              → Private application code
│   ├── domain/            → LAYER 1: Pure business logic
│   ├── port/              → LAYER 2: Interface contracts
│   ├── app/               → LAYER 3: Use cases / orchestration
│   ├── adapter/           → LAYER 4: Concrete implementations
│   ├── handler/           → LAYER 5: HTTP transport
│   └── config/            → Configuration from env vars
│
├── pkg/                   → Public reusable packages
├── migrations/            → SQL migration files (up/down)
├── mock/                  → Generated mocks from port interfaces
├── test/                  → Integration / E2E tests
├── Makefile               → make run, make test, make migrate
└── go.mod
```

## Dependency rule (absolute, never break this)

```
handler → app → domain ← port ← adapter
```

- **domain/** imports NOTHING except Go stdlib
- **port/** imports only domain/
- **app/** imports domain/ and port/ (interfaces only)
- **adapter/** imports domain/, port/, and external libraries
- **handler/** imports app/ and dto/

An adapter NEVER imports another adapter. An app service NEVER imports an adapter directly.

## Layer rules

### domain/ — Pure business entities and rules

- Zero external imports. Only Go standard library.
- Contains: entities (structs + methods), value objects, domain errors
- Entities validate themselves: `user.New(email, name, hash)` returns error if invalid
- Business rules live HERE, not in app/ or handler/

```go
// CORRECT: validation in domain
func New(email, name, hash string) (*User, error) {
    if email == "" { return nil, domain.ErrValidation }
    ...
}

// WRONG: validation in handler or app
```

**Files per domain module**: `entity.go`, `entity_test.go`, and optionally value objects

### port/ — Interface contracts

- Defines WHAT the system needs, not HOW
- Two sub-packages:
  - `repository/` — data persistence interfaces
  - `service/` — external service interfaces (payment, email, AI, storage, OAuth)
- Interfaces are small and specific. No god interfaces.

```go
// CORRECT: focused interface
type UserRepository interface {
    Create(ctx context.Context, u *user.User) error
    FindByID(ctx context.Context, id string) (*user.User, error)
    FindByEmail(ctx context.Context, email string) (*user.User, error)
    Update(ctx context.Context, u *user.User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, offset, limit int) ([]*user.User, int, error)
}

// WRONG: catch-all interface with 30 methods
```

### app/ — Use cases (application services)

- Orchestrates domain entities and port interfaces
- Each sub-package = one functional domain (auth, billing, user, ai, etc.)
- Receives dependencies via constructor injection
- Returns domain types and domain errors — NEVER HTTP concepts

```go
// CORRECT
type Service struct {
    users repository.UserRepository  // interface, not concrete type
    email service.EmailService
}

func NewService(users repository.UserRepository, email service.EmailService) *Service {
    return &Service{users: users, email: email}
}

// WRONG: importing postgres package directly
```

**Files**: `service.go` + `service_test.go` per module. Tests use mocks.

### adapter/ — Concrete implementations

- Each sub-package implements one or more port interfaces
- Sub-packages: `postgres/`, `stripe/`, `resend/`, `claude/`, `openai/`, `google/`, `r2/`
- Can import external libraries (stripe SDK, AWS SDK, etc.)
- Each adapter has: `client.go` (setup/config) + implementation files

```go
// postgres/user.go implements repository.UserRepository
type UserRepository struct { db *sql.DB }

// stripe/payment.go implements service.PaymentService
type PaymentService struct { client *stripe.Client }
```

**To swap a provider**: create new adapter, change ONE line in cmd/api/main.go. Nothing else changes.

### handler/ — HTTP transport

- Converts HTTP requests to app service calls and back
- Contains: route definitions, handlers, middleware, DTOs
- Sub-structure:
  - `router.go` — all route definitions
  - `auth.go`, `user.go`, `billing.go` — handler groups
  - `middleware/` — auth (JWT), CORS, rate limit, logging
  - `dto/request/` — incoming request structs with json tags
  - `dto/response/` — outgoing response structs + helpers

```go
// CORRECT: handler is thin, delegates to app service
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req request.LoginRequest
    json.NewDecoder(r.Body).Decode(&req)
    user, token, err := h.svc.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        response.HandleDomainError(w, err)
        return
    }
    response.JSON(w, http.StatusOK, response.AuthResponse{Token: token, User: response.UserFromDomain(user)})
}

// WRONG: business logic in handler
```

### pkg/ — Public reusable packages

- Can be imported by external projects
- Contains pure utilities: `jwt/`, `hash/`, `validate/`, `pagination/`
- Each package is self-contained with its own tests
- No imports from internal/

### config/ — Configuration

- Single `config.go` with typed Config struct
- All env vars loaded and validated at startup
- Default values for local development
- No config scattered across files

## How to add a new feature

Example: "Add a teams/organization feature"

1. **domain/team/entity.go** — Team struct, validation, business methods
2. **domain/team/entity_test.go** — Test validation rules
3. **port/repository/team.go** — TeamRepository interface
4. **app/team/service.go** — Use cases (CreateTeam, InviteMember, etc.)
5. **app/team/service_test.go** — Unit tests with mocked repository
6. **adapter/postgres/team.go** — SQL implementation of TeamRepository
7. **handler/team.go** — HTTP endpoints
8. **handler/dto/request/team.go** — Request DTOs
9. **handler/dto/response/team.go** — Response DTOs (add to response.go or new file)
10. **cmd/api/main.go** — Wire: `teamRepo → teamSvc → router`
11. **migrations/00X_create_teams.up.sql** — Database table

Always follow this order. Domain first, HTTP last.

## How to swap a provider

Example: "Replace Stripe with Lemon Squeezy"

1. Create `adapter/lemonsqueezy/client.go` + `payment.go`
2. Implement `service.PaymentService` interface
3. In `cmd/api/main.go`, change: `stripe.NewPaymentService(...)` → `lemonsqueezy.NewPaymentService(...)`
4. Done. Zero changes elsewhere.

## Testing strategy

### Unit tests (fast, no dependencies)
- **domain/*_test.go** — Entity validation, business rules
- **app/**/service_test.go** — Use cases with mocked ports
- **pkg/*_test.go** — Utility functions
- Run: `make test-unit`

### Integration tests (need Docker)
- **adapter/postgres/*_test.go** — Against real PostgreSQL
- **test/** — Full request flow tests
- Run: `make test-integration`

### Rules
- Test file lives next to the file it tests: `service.go` → `service_test.go`
- App layer tests mock ALL ports — test logic, not infrastructure
- Name tests: `TestServiceName_MethodName_Scenario`
- Table-driven tests for multiple scenarios

```go
func TestAuthService_Register_Success(t *testing.T) { ... }
func TestAuthService_Register_DuplicateEmail(t *testing.T) { ... }
func TestAuthService_Login_WrongPassword(t *testing.T) { ... }
```

## SQL conventions

- Pure SQL with `database/sql` + `lib/pq`. No ORM.
- Parameterized queries ONLY: `$1, $2, $3` — never string concatenation
- All queries use `context.Context` for timeout/cancellation
- Tables: UUID primary key, `created_at TIMESTAMP NOT NULL DEFAULT NOW()`, `updated_at`
- Use `TEXT` not `VARCHAR`. Index foreign keys.
- No cross-feature foreign keys (only reference `users` table)

## Migrations

Powered by `golang-migrate`. Migration files live in `backend/migrations/`.

### File naming

```
migrations/
├── 001_create_users.up.sql
├── 001_create_users.down.sql
├── 002_create_subscriptions.up.sql
├── 002_create_subscriptions.down.sql
└── ...
```

- Numbered sequentially: `001`, `002`, ..., `010`, ...
- Each migration has an `.up.sql` (apply) and `.down.sql` (rollback)
- snake_case descriptive name: `create_X`, `add_Y_to_X`, `drop_Z`
- Feature-scoped: each feature's tables in their own migration files

### Rules

- **Migrations are immutable.** Once applied in prod, NEVER edit — create a new migration instead.
- **Always write the down migration.** Every `up` must be reversible.
- **Use `IF NOT EXISTS` / `IF EXISTS`** for idempotent migrations.
- **Test locally before prod.** `make migrate-up` locally → verify → push → apply to prod.
- **No cross-feature foreign keys.** Only `REFERENCES users(id)` is allowed.

### Workflow: local → prod

```
1. Create migration files   →  /add-migration or manually
2. Test locally              →  make migrate-up (on Docker PostgreSQL)
3. Verify schema             →  DbGate at localhost:8082
4. Rollback test             →  make migrate-down (verify down works)
5. Re-apply                  →  make migrate-up
6. Commit & push             →  git commit
7. Apply to prod             →  DATABASE_URL=<neon_url> make migrate-up
```

### Fixing a broken migration

If a migration fails halfway (dirty state):
```bash
make migrate-status            # shows version + dirty flag
make migrate-force VERSION=N   # force-set to version N (the last clean version)
```
Then fix the SQL and re-run `make migrate-up`.

## Error handling

- Domain defines typed errors: `domain.ErrNotFound`, `domain.ErrUnauthorized`, etc.
- App layer returns domain errors, never HTTP codes
- Handler maps domain errors to HTTP responses via `response.HandleDomainError()`
- Wrap errors for context: `fmt.Errorf("creating user: %w", err)`
- Never swallow errors silently

## Commands

```bash
make run              # Start API server
make build            # Build binary
make test             # Run all tests
make test-unit        # Run unit tests only
make test-integration # Run integration tests
make migrate-up       # Apply all pending migrations
make migrate-down     # Rollback last migration
make migrate-down-all # Rollback ALL migrations
make migrate-status   # Show current migration version
make migrate-force    # Force version: make migrate-force VERSION=1
make seed             # Seed initial data (admin user, plans)
make mock             # Generate mocks
make lint             # Run linter
make tidy             # go mod tidy
```

---
name: add-feature
description: Scaffold a complete new feature across backend (Go hexagonal) and frontend (Next.js feature-based). Use when adding new functionality like teams, analytics, comments, etc.
user-invocable: true
allowed-tools: Read, Write, Edit, Bash, Grep, Glob, Agent
---

# Add Feature

Create the feature: **$ARGUMENTS**

You are scaffolding a new feature for CleanSaaS. Follow EVERY step below in order. Domain first, HTTP last. Read existing code for patterns before writing anything.

---

## STEP 0 — Understand the request

Parse `$ARGUMENTS` to determine:
- **Feature name** (kebab-case for files, PascalCase for types, lowercase for packages)
- **Which layers** are needed: backend only, frontend only, or full-stack
- **Core entities** and their fields
- **CRUD operations** or custom use cases needed
- **Whether it needs new database tables**

If the request is ambiguous, ask the user to clarify before proceeding.

---

## STEP 1 — Backend: Domain entity

Create `backend/internal/domain/{feature}/entity.go`:

**Pattern** (reference `backend/internal/domain/user/entity.go`):
- Struct with business fields + `CreatedAt`, `UpdatedAt time.Time`
- Constructor `New(...)` that validates required fields, returns `(*Entity, error)`
- Business methods on the struct (e.g., `IsActive()`, `Cancel()`)
- Import ONLY Go stdlib + `domain` errors package. ZERO external imports.

```go
package {feature}

import (
    "time"
    "github.com/hassad/boilerplateSaaS/backend/internal/domain"
)

type {Entity} struct {
    ID        string
    // ... fields
    CreatedAt time.Time
    UpdatedAt time.Time
}

func New(/* required params */) (*{Entity}, error) {
    // validate required fields
    if someField == "" {
        return nil, domain.ErrValidation
    }
    return &{Entity}{
        // ... set fields
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}
```

Also create `entity_test.go` with table-driven tests for validation and business methods.

---

## STEP 2 — Backend: Port interfaces

Create `backend/internal/port/repository/{feature}.go`:

**Pattern** (reference `backend/internal/port/repository/user.go`):
- Small, focused interface
- All methods take `context.Context` as first param
- Return domain types, never SQL types

```go
package repository

import (
    "context"
    "github.com/hassad/boilerplateSaaS/backend/internal/domain/{feature}"
)

type {Feature}Repository interface {
    Create(ctx context.Context, entity *{feature}.{Entity}) error
    FindByID(ctx context.Context, id string) (*{feature}.{Entity}, error)
    // ... only methods actually needed
    List(ctx context.Context, offset, limit int) ([]*{feature}.{Entity}, int, error)
}
```

If the feature needs an external service (AI, email, payment), add an interface in `port/service/` too.

---

## STEP 3 — Backend: Application service

Create `backend/internal/app/{feature}/service.go`:

**Pattern** (reference `backend/internal/app/auth/service.go`):
- Struct holds port interfaces (repository + services), injected via constructor
- Methods are use cases, not CRUD wrappers
- Returns domain types and domain errors — NEVER HTTP concepts
- Use `context.Context` on all methods

```go
package {feature}

import (
    "context"
    "github.com/hassad/boilerplateSaaS/backend/internal/domain/{feature}"
    "github.com/hassad/boilerplateSaaS/backend/internal/port/repository"
)

type Service struct {
    {feature}s repository.{Feature}Repository
}

func NewService({feature}s repository.{Feature}Repository) *Service {
    return &Service{{feature}s: {feature}s}
}

func (s *Service) Create(ctx context.Context, /* params */) (*{feature}.{Entity}, error) {
    entity, err := {feature}.New(/* params */)
    if err != nil {
        return nil, err
    }
    if err := s.{feature}s.Create(ctx, entity); err != nil {
        return nil, err
    }
    return entity, nil
}
```

Also create `service_test.go` — unit tests with mocked repository. Use table-driven tests. Name: `TestServiceName_MethodName_Scenario`.

---

## STEP 4 — Backend: PostgreSQL adapter

Create `backend/internal/adapter/postgres/{feature}.go`:

**Pattern** (reference `backend/internal/adapter/postgres/user.go`):
- Struct holds `*sql.DB`
- Implements the repository interface
- Pure SQL with `$1, $2, $3` parameters — NEVER string concatenation
- Map `sql.ErrNoRows` to `domain.ErrNotFound`
- Use `QueryRowContext` / `ExecContext` / `QueryContext` — always with context

```go
package postgres

import (
    "context"
    "database/sql"
    "github.com/hassad/boilerplateSaaS/backend/internal/domain"
    "github.com/hassad/boilerplateSaaS/backend/internal/domain/{feature}"
)

type {Feature}Repository struct {
    db *sql.DB
}

func New{Feature}Repository(db *sql.DB) *{Feature}Repository {
    return &{Feature}Repository{db: db}
}

func (r *{Feature}Repository) Create(ctx context.Context, e *{feature}.{Entity}) error {
    err := r.db.QueryRowContext(ctx,
        `INSERT INTO {table} (/* columns */) VALUES ($1, $2, ...) RETURNING id`,
        /* fields */,
    ).Scan(&e.ID)
    return err
}

func (r *{Feature}Repository) FindByID(ctx context.Context, id string) (*{feature}.{Entity}, error) {
    e := &{feature}.{Entity}{}
    err := r.db.QueryRowContext(ctx,
        `SELECT /* columns */ FROM {table} WHERE id = $1`, id,
    ).Scan(/* fields */)
    if err == sql.ErrNoRows {
        return nil, domain.ErrNotFound
    }
    return e, err
}
```

---

## STEP 5 — Backend: SQL migration

Create migration files using the next available number (check `backend/migrations/` for the latest):
- `backend/migrations/{NNN}_create_{feature}.up.sql`
- `backend/migrations/{NNN}_create_{feature}.down.sql`

**SQL conventions:**
- UUID primary key: `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- Timestamps: `created_at TIMESTAMP NOT NULL DEFAULT NOW()`, `updated_at TIMESTAMP NOT NULL DEFAULT NOW()`
- Use `TEXT` not `VARCHAR`
- Foreign key to users: `user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE`
- **NO cross-feature foreign keys** (only reference `users` table)
- Index all foreign keys and frequently queried columns

```sql
-- up
CREATE TABLE {feature_table} (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- ... feature columns
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_{feature_table}_user_id ON {feature_table}(user_id);
```

```sql
-- down
DROP TABLE IF EXISTS {feature_table};
```

---

## STEP 6 — Backend: HTTP handler + DTOs

Create `backend/internal/handler/{feature}.go`:

**Pattern** (reference `backend/internal/handler/auth.go`):
- Handler struct holds `*app.Service`
- Each method: decode request → call service → encode response
- Thin layer — NO business logic here

```go
package handler

type {Feature}Handler struct {
    svc *{feature}.Service
}

func New{Feature}Handler(svc *{feature}.Service) *{Feature}Handler {
    return &{Feature}Handler{svc: svc}
}

func (h *{Feature}Handler) Create(w http.ResponseWriter, r *http.Request) {
    var req request.Create{Feature}Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.Error(w, http.StatusBadRequest, "invalid request body")
        return
    }
    entity, err := h.svc.Create(r.Context(), /* params from req */)
    if err != nil {
        response.HandleDomainError(w, err)
        return
    }
    response.JSON(w, http.StatusCreated, /* response DTO */)
}
```

Create request DTOs in `handler/dto/request/{feature}.go`.
Add response DTO or helper in `handler/dto/response/` (extend `response.go` or create `{feature}.go`).

---

## STEP 7 — Backend: Wire routes

Update `backend/internal/handler/router.go`:
- Add the new service parameter to `NewRouter()`
- Register routes in the appropriate group (public or protected)

**Pattern:**
```go
// In protected routes group:
{feature}Handler := New{Feature}Handler({feature}Svc)
r.Post("/{feature}s", {feature}Handler.Create)
r.Get("/{feature}s", {feature}Handler.List)
r.Get("/{feature}s/{id}", {feature}Handler.GetByID)
```

Update `backend/cmd/api/main.go`:
- Instantiate the repository: `{feature}Repo := postgres.New{Feature}Repository(db)`
- Instantiate the service: `{feature}Svc := app{feature}.NewService({feature}Repo)`
- Pass to router: `handler.NewRouter(..., {feature}Svc, ...)`

---

## STEP 8 — Frontend: Types

Create `frontend/src/features/{feature}/types.ts`:

**Pattern** (reference `frontend/src/features/auth/types.ts`):
```ts
export type {Entity} = {
  id: string;
  // ... fields matching API response
  created_at: string;
  updated_at: string;
};

export type Create{Entity}Data = {
  // ... fields for create form
};
```

---

## STEP 9 — Frontend: Server Actions

Create `frontend/src/features/{feature}/actions/{feature}.ts`:

**Pattern** (reference `frontend/src/features/auth/actions/auth.ts`):
```ts
"use server";

import { api } from "@/shared/lib/api";
import type { {Entity} } from "../types";

export async function create{Entity}(data: Create{Entity}Data) {
  return api<{Entity}>("/{feature}s", {
    method: "POST",
    body: data,
  });
}

export async function get{Entity}s() {
  return api<{Entity}[]>("/{feature}s");
}
```

---

## STEP 10 — Frontend: Client hook (if needed)

Create `frontend/src/features/{feature}/hooks/use-{feature}.ts` only if the feature needs client-side state management (forms, real-time updates, optimistic UI).

**Pattern** (reference `frontend/src/features/auth/hooks/use-auth.ts`):
```ts
"use client";

import { useState, useEffect } from "react";
import { api } from "@/shared/lib/api";
import type { {Entity} } from "../types";

export function use{Feature}() {
  const [items, setItems] = useState<{Entity}[]>([]);
  const [loading, setLoading] = useState(true);
  // ... state logic
  return { items, loading };
}
```

---

## STEP 11 — Frontend: Components

Create components in `frontend/src/features/{feature}/components/`:

- Use **Server Components** by default (data display, lists)
- Use `"use client"` only for forms, modals, interactive elements
- Import types from `../types`, hooks from `../hooks/`, actions from `../actions/`
- **NEVER import from another feature** — only from `@/shared/`
- Styling: Tailwind only, use `cn()` from `@/shared/lib/utils` for conditional classes

---

## STEP 12 — Frontend: Route page

Create `frontend/src/app/(dashboard)/{feature}/page.tsx` (or appropriate route group):

**Pattern** — thin page, 5-20 lines max:
```tsx
import type { Metadata } from "next";
import { {FeatureComponent} } from "@/features/{feature}/components/{component}";

export const metadata: Metadata = { title: "{Feature}" };

export default function {Feature}Page() {
  return <{FeatureComponent} />;
}
```

---

## STEP 13 — Verify independence

Run this mental checklist:
- [ ] Backend domain imports ONLY Go stdlib + domain errors
- [ ] App service depends on port interfaces, not concrete adapters
- [ ] No cross-feature imports (backend or frontend)
- [ ] Database tables only FK to `users`, not to other feature tables
- [ ] Frontend feature uses only `@/shared/` for cross-cutting concerns
- [ ] Removing this feature's folders + main.go lines = everything still compiles

If any check fails, fix it before finishing.

---

## STEP 14 — Compile and verify

```bash
cd backend && go build ./...
cd frontend && npx tsc --noEmit
```

Fix any errors before reporting done.

---

## Output

When finished, report:
1. Files created (grouped by layer)
2. Files modified (router.go, main.go)
3. Migration file names
4. Any decisions made and why
5. Independence verification result

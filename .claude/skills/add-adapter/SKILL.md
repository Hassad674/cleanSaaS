---
name: add-adapter
description: Add or swap an adapter (external service implementation) following hexagonal architecture. Use when integrating a new provider (payment, email, AI, storage, OAuth) or replacing an existing one.
user-invocable: true
allowed-tools: Read, Write, Edit, Bash, Grep, Glob
---

# Add Adapter

Create or swap adapter: **$ARGUMENTS**

You are adding a new adapter to CleanSaaS following hexagonal architecture. An adapter implements a port interface — the rest of the codebase is completely unaware of which concrete provider is used.

---

## STEP 0 — Parse the request

Determine from `$ARGUMENTS`:
- **Adapter type**: payment, email, AI, storage, OAuth, or new category
- **Provider name**: stripe, lemonsqueezy, resend, sendgrid, openai, claude, r2, s3, google, github, etc.
- **Action**: new adapter, or replacing an existing one

---

## STEP 1 — Identify the port interface

Find the interface this adapter must implement.

### Existing port interfaces
Read the relevant file in `backend/internal/port/service/`:

| Category | File | Interface |
|----------|------|-----------|
| Payment | `port/service/payment.go` | `PaymentService` |
| Email | `port/service/email.go` | `EmailService` |
| AI | `port/service/ai.go` | `AIService` |
| Storage | `port/service/storage.go` | `StorageService` |
| OAuth | `port/service/auth_provider.go` | `OAuthProvider` |

Read the interface file to get the exact method signatures.

If no interface exists for this adapter category, create one in `port/service/` first:

```go
package service

import "context"

type {Category}Service interface {
    // Define methods based on the use case
}
```

Keep the interface small and focused. Define only what the app layer actually needs.

---

## STEP 2 — Create the adapter package

Create `backend/internal/adapter/{provider}/` with two files:

### 2a. Client setup — `client.go`

```go
package {provider}

// Client holds the provider's SDK client and configuration.
type Client struct {
    // SDK client, API key, base URL, etc.
}

// NewClient creates a configured {provider} client.
func NewClient(apiKey string) *Client {
    return &Client{
        // initialize SDK or HTTP client
    }
}
```

### 2b. Interface implementation — `{category}.go`

**Pattern** (reference existing adapters in `backend/internal/adapter/`):

```go
package {provider}

import (
    "context"
    // provider SDK
)

// {Category}Service implements service.{Category}Service using {Provider}.
type {Category}Service struct {
    client *Client
}

func New{Category}Service(client *Client) *{Category}Service {
    return &{Category}Service{client: client}
}

// Implement every method from the port interface.
func (s *{Category}Service) MethodName(ctx context.Context, /* params */) (/* returns */) {
    // Map domain types → provider SDK types
    // Call provider API
    // Map provider response → domain types
    // Return domain types and domain errors
}
```

### Rules:
- The adapter imports domain types and the provider's SDK — nothing else from internal/
- **Never import another adapter** (no `adapter/postgres` from `adapter/stripe`)
- Return domain errors (`domain.ErrNotFound`, etc.), not provider-specific errors
- Use `context.Context` on all methods for timeout/cancellation
- Handle provider errors gracefully — map to domain errors

---

## STEP 3 — Add configuration

If the adapter needs API keys or config, add them to `backend/internal/config/config.go`:

```go
// In Config struct:
{Provider}Key    string
{Provider}Secret string

// In Load():
{Provider}Key:    env("{PROVIDER}_API_KEY", ""),
{Provider}Secret: env("{PROVIDER}_SECRET", ""),
```

Only add config fields that don't already exist.

---

## STEP 4 — Install SDK dependency

If the provider has a Go SDK:
```bash
cd ./backend
go get {sdk-package}
go mod tidy
```

If no official SDK, use `net/http` directly. Keep it simple.

---

## STEP 5 — Wire in main.go

Update `backend/cmd/api/main.go`:

### For a NEW adapter (no existing provider):
```go
// Add import
import {provider} "github.com/hassad/boilerplateSaaS/backend/internal/adapter/{provider}"

// In main():
{provider}Client := {provider}.NewClient(cfg.{Provider}Key)
{category}Svc := {provider}.New{Category}Service({provider}Client)

// Pass to the app service that needs it:
featureSvc := appfeature.NewService(repo, {category}Svc)
```

### For SWAPPING an existing adapter:
Change only the instantiation lines. Example — replace Stripe with LemonSqueezy:

```go
// BEFORE:
stripeClient := stripe.NewClient(cfg.StripeKey)
paymentSvc := stripe.NewPaymentService(stripeClient)

// AFTER:
lsClient := lemonsqueezy.NewClient(cfg.LemonSqueezyKey)
paymentSvc := lemonsqueezy.NewPaymentService(lsClient)
```

**Nothing else changes.** The app service receives the same interface.

---

## STEP 6 — Write tests

Create `backend/internal/adapter/{provider}/{category}_test.go`:

- Test that the adapter correctly implements the interface (compile-time check):
  ```go
  var _ service.{Category}Service = (*{Category}Service)(nil)
  ```
- Unit tests with mocked HTTP responses (if using HTTP client)
- No tests that hit the real provider API (those go in integration tests)

---

## STEP 7 — Verify

```bash
cd ./backend && go build ./...
```

### Verify the swap promise:
- [ ] Only `cmd/api/main.go` changed for wiring
- [ ] No app/ or handler/ files were modified
- [ ] No domain/ files were modified
- [ ] The adapter implements the full port interface
- [ ] Config is centralized in config.go

---

## Output

Report:
1. **Files created** — adapter package files
2. **Files modified** — config.go, main.go
3. **Interface implemented** — which port interface + all methods
4. **Dependencies added** — new Go modules
5. **Swap instructions** — if replacing, what to change back

Example:
```
Created adapter: lemonsqueezy

Files created:
  backend/internal/adapter/lemonsqueezy/client.go
  backend/internal/adapter/lemonsqueezy/payment.go
  backend/internal/adapter/lemonsqueezy/payment_test.go

Modified:
  backend/internal/config/config.go — added LemonSqueezyKey, LemonSqueezyWebhookSecret
  backend/cmd/api/main.go — swapped stripe → lemonsqueezy

Implements: service.PaymentService (8 methods)
Dependencies: github.com/lemonsqueezy/go-sdk v1.2.0

To revert to Stripe:
  In main.go, change lemonsqueezy.NewPaymentService → stripe.NewPaymentService
```

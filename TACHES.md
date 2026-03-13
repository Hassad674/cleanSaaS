# CleanSaaS — Autonomous Task List

> **This file is the single source of truth for the autonomous Claude session.**
> Read it fully before starting. Re-read it after every context compression.

---

## 0. BEFORE YOU START — Mandatory Setup

### 0.1 — Understand the project

You are building **CleanSaaS**, an open-source SaaS boilerplate. Read these files IN THIS ORDER before doing anything:

```
1. /home/hassad/Documents/boilerplateSaaS/CLAUDE.md              → Project rules, architecture, conventions
2. /home/hassad/Documents/boilerplateSaaS/backend/CLAUDE.md      → Backend architecture (hexagonal), layer rules, SQL conventions
3. /home/hassad/Documents/boilerplateSaaS/frontend/CLAUDE.md     → Frontend architecture (feature-based), design system, landing page structure
4. /home/hassad/Documents/boilerplateSaaS/FEATURES.md            → Feature descriptions and API keys status
5. /home/hassad/Documents/boilerplateSaaS/backend/.env            → Available API keys (DO NOT commit this file)
6. /home/hassad/Documents/boilerplateSaaS/backend/cmd/api/main.go → Current wiring, see what's already connected
```

### 0.2 — Understand what already exists

Before each task, check what's already scaffolded. Many domain entities, ports, and service shells exist but are empty. DO NOT recreate files that exist — extend them.

Run `find backend/internal -name "*.go" | sort` and `find frontend/src -name "*.ts" -o -name "*.tsx" | sort` to see the current file tree.

### 0.3 — Verify local environment works

```bash
cd /home/hassad/Documents/boilerplateSaaS && docker compose up -d    # DB must be running
cd backend && go build ./...                                          # Must compile
cd ../frontend && npx tsc --noEmit                                    # Must compile
```

If any of these fail, fix them BEFORE starting tasks. Do NOT proceed with a broken baseline.

---

## 1. RULES — Follow these at ALL times

### 1.1 — Architecture rules (NEVER break these)

**Backend hexagonal layers — dependency direction is absolute:**
```
handler → app → domain ← port ← adapter
```
- `domain/` imports NOTHING except Go stdlib
- `port/` imports only `domain/`
- `app/` imports `domain/` and `port/` (interfaces only)
- `adapter/` imports `domain/`, `port/`, and external libraries
- `handler/` imports `app/`, `dto/`, `middleware/`, `pkg/`
- An adapter NEVER imports another adapter
- An app service NEVER imports an adapter directly

**Frontend feature isolation:**
- `features/X/` NEVER imports from `features/Y/`
- Cross-feature composition happens ONLY in `app/` pages
- Features import from: their own folder, `@/shared/`, npm packages
- `app/` pages are thin (5-20 lines): import + compose, no logic

**Database:**
- Pure SQL, parameterized queries only (`$1, $2` — never `fmt.Sprintf` with SQL)
- No cross-feature foreign keys (only `REFERENCES users(id)` allowed)
- Every migration has `.up.sql` AND `.down.sql`
- Use `IF NOT EXISTS` / `IF EXISTS` for idempotent migrations

### 1.2 — Commit rules

After completing each task (not each sub-step, each TASK):

1. Stage only relevant files: `git add <specific files>` — NEVER `git add .` or `git add -A`
2. Verify BEFORE committing:
   ```bash
   cd /home/hassad/Documents/boilerplateSaaS/backend && go build ./...
   cd /home/hassad/Documents/boilerplateSaaS/frontend && npx tsc --noEmit
   cd /home/hassad/Documents/boilerplateSaaS/backend && go test ./... -count=1 -short
   ```
3. ALL THREE must pass. If they don't, fix before committing.
4. Commit with conventional message:
   ```bash
   git commit -m "$(cat <<'EOF'
   feat: implement <feature name>

   <2-3 lines describing what was added>

   Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
   EOF
   )"
   ```
5. NEVER commit `.env` files, `node_modules/`, `.next/`, or files with secrets.

### 1.3 — Testing rules

- Write unit tests for EVERY domain entity and EVERY app service method you create
- Domain tests: test validation, business rules, edge cases
- App service tests: mock ALL ports (repositories, external services), test orchestration logic
- Test file lives next to source: `service.go` → `service_test.go`
- Name pattern: `TestServiceName_MethodName_Scenario`
- Run tests after writing them to confirm they pass

### 1.4 — Blocker policy

If stuck on ANY issue for more than 30 minutes:

1. Create `BLOCKED-taskX.md` at project root with:
   ```
   # Blocked: Task X — <task name>
   ## What failed
   <exact error message or description>
   ## What I tried
   <list of approaches attempted>
   ## Possible solutions
   <ideas for fixing later>
   ```
2. Commit whatever working code exists for that task
3. Move to the next task immediately
4. Do NOT loop retrying the same approach

### 1.5 — Design system rules (frontend)

- NEVER use hardcoded Tailwind colors (`text-zinc-500`, `bg-white`, `border-gray-200`)
- ALWAYS use semantic tokens: `text-foreground`, `text-muted-foreground`, `bg-card`, `border-border`, `bg-primary`, `text-primary-foreground`, etc.
- NEVER use `dark:` prefix — CSS variables handle dark mode automatically
- Mobile-first: write base styles, then `sm:`, `md:`, `lg:` overrides
- Cards: `bg-card border border-border rounded-xl p-6 shadow-sm`
- Buttons primary: `bg-primary text-primary-foreground rounded-lg hover:opacity-90 transition-opacity`

### 1.6 — Code quality rules

- TypeScript: no `any`, no `as` casts unless absolutely necessary (with comment)
- Go: always accept `context.Context` as first parameter in service/repo methods
- Go: always wrap errors with context: `fmt.Errorf("creating user: %w", err)`
- Go: never swallow errors (`_ = fn()` where fn returns error is forbidden)
- Keep files small: one file = one responsibility, under 200 lines ideally
- Imports: always `@/` alias in frontend, never relative `../../`

### 1.7 — Context recovery

After a context compression, you will lose the conversation history. When that happens:

1. Re-read this file (`TACHES.md`) to know what to do
2. Check which tasks are already marked `[x]` below
3. Run `git log --oneline -20` to see what was already committed
4. Check for any `BLOCKED-*.md` files
5. Resume from the first unchecked task

---

## 2. TASKS — Execute in order

Update the checkbox (`[ ]` → `[x]`) in this file after completing each task.

---

### TASK 1: Backend Middleware & Infrastructure
> Priority: FIRST — this benefits all subsequent features

- [ ] **1a. Structured logging middleware**
  - Create `backend/internal/handler/middleware/logging.go`
  - Use Go `log/slog` with JSON output
  - Log: method, path, status, duration, request_id (UUID per request)
  - Add request_id to context so all layers can access it
  - Wire in `handler/router.go` as global middleware

- [ ] **1b. Rate limiting middleware**
  - Create `backend/internal/handler/middleware/ratelimit.go`
  - Token bucket algorithm per IP address
  - Default: 100 requests/minute for API, 10 requests/minute for auth endpoints
  - Return `429 Too Many Requests` with `Retry-After` header
  - Wire in `handler/router.go` (different limits for different route groups)

- [ ] **1c. Health check endpoint**
  - Add `GET /health` in `handler/router.go`
  - Returns: `{"status":"ok","db":"connected","uptime":"...","version":"1.0.0"}`
  - Ping database to check connection
  - No auth required

- [ ] **1d. Graceful shutdown**
  - Modify `backend/cmd/api/main.go`
  - Listen for `SIGTERM` and `SIGINT`
  - Drain active connections with timeout (30s)
  - Close database pool
  - Log shutdown sequence

- [ ] **1e. Commit**: `feat: add production middleware (logging, rate-limit, health, graceful shutdown)`

---

### TASK 2: Resend Email Adapter
> Priority: HIGH — unlocks password reset and email verification

- [ ] **2a. Email service port** — read `backend/internal/port/service/email.go`, verify interface is sufficient:
  ```go
  type EmailService interface {
      Send(ctx context.Context, to, subject, htmlBody string) error
  }
  ```
  If it needs more methods (SendTemplate, etc.), update it.

- [ ] **2b. Resend adapter**
  - Create `backend/internal/adapter/resend/client.go` — Resend API client setup
  - Create `backend/internal/adapter/resend/email.go` — implements `service.EmailService`
  - Use `RESEND_API_KEY` from config
  - From address: `noreply@cleansaas.dev` (or configurable)
  - Go dependency: `go get github.com/resend/resend-go/v2` (run from backend/)

- [ ] **2c. Email templates**
  - Create `backend/internal/adapter/resend/templates.go`
  - Templates as Go functions returning HTML strings (keep it simple, no template engine):
    - `WelcomeEmail(name string) (subject, body string)`
    - `VerificationEmail(name, link string) (subject, body string)`
    - `PasswordResetEmail(name, link string) (subject, body string)`
  - Clean, responsive HTML (inline CSS), matches design system colors (rose primary)

- [ ] **2d. Wire in main.go**
  - Create Resend email service: `emailSvc := resend.NewEmailService(cfg.ResendKey)`
  - Pass to auth service: `appauth.NewService(userRepo, emailSvc, jwtMaker)`

- [ ] **2e. Test**: Unit test for template generation (correct subject, body contains link)

- [ ] **2f. Commit**: `feat: implement Resend email adapter with templates`

---

### TASK 3: Settings Page
> Extends existing user feature

- [ ] **3a. Backend — password change endpoint**
  - Add method to `app/user/service.go`: `ChangePassword(ctx, userID, oldPassword, newPassword string) error`
  - Verify old password with bcrypt, hash new password, update in DB
  - Add handler in `handler/user.go`: `PUT /api/users/me/password`
  - Add DTO: `dto/request/user.go` — `ChangePasswordRequest{OldPassword, NewPassword string}`

- [ ] **3b. Backend — delete account endpoint**
  - Add method to `app/user/service.go`: `DeleteAccount(ctx, userID string) error`
  - Add `Delete(ctx, id string) error` to `adapter/postgres/user.go` if not exists
  - Add handler: `DELETE /api/users/me`

- [ ] **3c. Frontend — settings components**
  - Create `frontend/src/features/user/components/settings-profile.tsx` — edit name form
  - Create `frontend/src/features/user/components/settings-password.tsx` — change password form
  - Create `frontend/src/features/user/components/settings-danger.tsx` — delete account button with confirmation
  - All forms use design tokens, show loading/success/error states

- [ ] **3d. Frontend — settings page**
  - Update `frontend/src/app/(dashboard)/settings/page.tsx` — compose the 3 settings components
  - Sections separated by cards with clear headings

- [ ] **3e. Test**: Unit test for `ChangePassword` service method (correct old pw, wrong old pw, same password)

- [ ] **3f. Commit**: `feat: implement settings page (profile, password, account deletion)`

---

### TASK 4: Forgot Password / Reset Password
> Depends on: Task 2 (Resend)

- [ ] **4a. Migration**
  - Create `backend/migrations/002_create_password_resets.up.sql`:
    ```sql
    CREATE TABLE IF NOT EXISTS password_resets (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        token TEXT NOT NULL UNIQUE,
        expires_at TIMESTAMP NOT NULL,
        used BOOLEAN NOT NULL DEFAULT false,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
    CREATE INDEX IF NOT EXISTS idx_password_resets_token ON password_resets(token);
    CREATE INDEX IF NOT EXISTS idx_password_resets_user_id ON password_resets(user_id);
    ```
  - Create matching `.down.sql`
  - Run `cd backend && make migrate-up` to apply

- [ ] **4b. Backend — domain + port**
  - Add password reset token logic in auth service (or create new if cleaner)
  - Port: if needed, add `PasswordResetRepository` interface in `port/repository/`
  - Adapter: implement in `adapter/postgres/password_reset.go`

- [ ] **4c. Backend — endpoints**
  - `POST /api/auth/forgot-password` — takes email, generates token, sends email via Resend
  - `POST /api/auth/reset-password` — takes token + new password, validates, updates password
  - Both are PUBLIC (no auth required)
  - Rate limit the forgot-password endpoint (max 3/hour per email)

- [ ] **4d. Frontend**
  - Update `frontend/src/features/auth/components/forgot-password-form.tsx` (file exists but is placeholder)
  - Create `frontend/src/features/auth/components/reset-password-form.tsx`
  - Create page `frontend/src/app/(auth)/reset-password/page.tsx`
  - Add "Forgot password?" link on login form

- [ ] **4e. Test**: Service test — request reset, use valid token, use expired token, use already-used token

- [ ] **4f. Commit**: `feat: implement forgot/reset password flow with email`

---

### TASK 5: Email Verification
> Depends on: Task 2 (Resend)

- [ ] **5a. Migration**
  - Create `backend/migrations/003_create_email_verifications.up.sql`:
    - Table: `email_verifications` (id, user_id FK, token UNIQUE, expires_at, created_at)
  - Create matching `.down.sql`
  - Apply: `make migrate-up`

- [ ] **5b. Backend**
  - After registration, generate verification token and send email
  - `POST /api/auth/verify-email` — takes token, marks user `email_verified = true`
  - `POST /api/auth/resend-verification` — authenticated, resends verification email
  - Rate limit resend (max 3/hour)

- [ ] **5c. Frontend**
  - Create `frontend/src/app/(auth)/verify-email/page.tsx`
  - Show "Check your email" message after registration
  - Handle token from URL query param, call verify endpoint
  - Show success/error state
  - Add resend button

- [ ] **5d. Test**: Service test — verify valid token, expired token, already verified user

- [ ] **5e. Commit**: `feat: implement email verification flow`

---

### TASK 6: Stripe Billing
> Complex feature — take your time, follow hexagonal strictly

- [ ] **6a. Migration**
  - Create `backend/migrations/004_create_billing.up.sql`:
    ```sql
    -- Plans table (seeded, not user-created)
    CREATE TABLE IF NOT EXISTS plans (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        name TEXT NOT NULL,
        stripe_price_id TEXT NOT NULL UNIQUE,
        price_cents INTEGER NOT NULL,
        interval TEXT NOT NULL DEFAULT 'month',
        features JSONB NOT NULL DEFAULT '[]',
        is_active BOOLEAN NOT NULL DEFAULT true,
        sort_order INTEGER NOT NULL DEFAULT 0,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    -- Subscriptions
    CREATE TABLE IF NOT EXISTS subscriptions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        plan_id UUID NOT NULL REFERENCES plans(id),
        stripe_subscription_id TEXT NOT NULL UNIQUE,
        status TEXT NOT NULL DEFAULT 'active',
        current_period_start TIMESTAMP NOT NULL,
        current_period_end TIMESTAMP NOT NULL,
        cancel_at_period_end BOOLEAN NOT NULL DEFAULT false,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
    CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
    CREATE INDEX IF NOT EXISTS idx_subscriptions_stripe_id ON subscriptions(stripe_subscription_id);

    -- Invoices
    CREATE TABLE IF NOT EXISTS invoices (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        stripe_invoice_id TEXT NOT NULL UNIQUE,
        amount_cents INTEGER NOT NULL,
        currency TEXT NOT NULL DEFAULT 'usd',
        status TEXT NOT NULL,
        invoice_url TEXT NOT NULL DEFAULT '',
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
    CREATE INDEX IF NOT EXISTS idx_invoices_user_id ON invoices(user_id);
    ```
  - Create matching `.down.sql` (DROP in reverse order: invoices, subscriptions, plans)
  - Apply migration

- [ ] **6b. Seed plans**
  - Update `backend/cmd/seed/main.go` to insert 3 plans:
    - Free ($0, basic features)
    - Pro ($19/mo, all features)
    - Enterprise ($49/mo, all features + priority support)
  - Use real Stripe price IDs from test mode OR placeholder strings that will be configured

- [ ] **6c. Backend — port + adapter**
  - Read existing `port/service/payment.go` — extend if needed with:
    - `CreateCheckoutSession(ctx, userID, priceID string) (sessionURL string, err error)`
    - `CreatePortalSession(ctx, customerID string) (sessionURL string, err error)`
    - `HandleWebhook(ctx, payload []byte, signature string) error`
  - Read existing `port/repository/subscription.go` — extend if needed
  - Create `adapter/stripe/client.go` — Stripe client setup
  - Create `adapter/stripe/payment.go` — implements `PaymentService`
  - Create `adapter/postgres/subscription.go` — implements `SubscriptionRepository`
  - Create `adapter/postgres/plan.go` — implements plan queries
  - Create `adapter/postgres/invoice.go` — implements invoice queries
  - Go dependency: `go get github.com/stripe/stripe-go/v82` (run from backend/)

- [ ] **6d. Backend — billing service**
  - Flesh out `app/billing/service.go`:
    - `GetPlans(ctx) ([]*Plan, error)`
    - `CreateCheckout(ctx, userID, planID string) (string, error)` — creates Stripe checkout session
    - `GetSubscription(ctx, userID string) (*Subscription, error)`
    - `CancelSubscription(ctx, userID string) error`
    - `HandleWebhook(ctx, payload []byte, sig string) error` — processes checkout.session.completed, invoice.paid, customer.subscription.updated/deleted

- [ ] **6e. Backend — handlers**
  - Create `handler/billing.go`:
    - `GET /api/billing/plans` — public, list active plans
    - `POST /api/billing/checkout` — authenticated, create checkout session
    - `GET /api/billing/subscription` — authenticated, get current subscription
    - `POST /api/billing/cancel` — authenticated, cancel at period end
    - `POST /api/billing/portal` — authenticated, Stripe customer portal
    - `POST /api/webhooks/stripe` — public (verified by signature), webhook handler
  - DTOs for request/response

- [ ] **6f. Frontend — pricing page**
  - Create `frontend/src/features/billing/components/pricing-cards.tsx` — 3 plan cards (Free/Pro/Enterprise)
  - Create `frontend/src/features/billing/actions/billing.ts` — API calls
  - Create `frontend/src/features/billing/hooks/use-billing.ts` — subscription state
  - Update `frontend/src/app/(marketing)/pricing/page.tsx` — compose pricing cards

- [ ] **6g. Frontend — billing settings**
  - Create `frontend/src/features/billing/components/subscription-status.tsx` — current plan, next billing date, cancel button
  - Create `frontend/src/features/billing/components/invoice-list.tsx` — list of invoices
  - Create `frontend/src/app/(dashboard)/settings/billing/page.tsx` or add section to settings page

- [ ] **6h. Wire in main.go** — stripe adapter → billing service → handler → router

- [ ] **6i. Test**: Billing service unit tests (mocked Stripe adapter): create checkout, handle webhook events, cancel

- [ ] **6j. Commit**: `feat: implement Stripe billing (plans, checkout, subscriptions, webhooks, invoices)`

---

### TASK 7: Cloudflare R2 File Storage
> Uses S3-compatible API

- [ ] **7a. Migration**
  - Create `backend/migrations/005_create_files.up.sql`:
    - Table: `files` (id UUID, user_id FK, name TEXT, key TEXT UNIQUE, size_bytes BIGINT, content_type TEXT, url TEXT, created_at, updated_at)
  - Create matching `.down.sql`
  - Apply migration

- [ ] **7b. Backend — adapter**
  - Create `backend/internal/adapter/r2/client.go` — S3-compatible client using `github.com/aws/aws-sdk-go-v2`
  - Create `backend/internal/adapter/r2/storage.go` — implements `service.StorageService`
  - Methods: `Upload(ctx, key string, data io.Reader, contentType string) (url string, err error)`, `Delete(ctx, key string) error`, `GetSignedURL(ctx, key string, duration time.Duration) (string, error)`

- [ ] **7c. Backend — service + handler**
  - Flesh out `app/storage/service.go`: Upload (validate type/size, store metadata in DB, upload to R2), Delete, List, GetByID
  - Create `adapter/postgres/file.go` — file metadata repository
  - Create `handler/storage.go`:
    - `POST /api/files/upload` — multipart upload, authenticated
    - `GET /api/files` — list user's files, paginated
    - `DELETE /api/files/:id` — delete file + R2 object
  - Max file size: 50MB
  - Allowed types: images (jpg, png, gif, webp), videos (mp4, webm), documents (pdf, doc, txt)

- [ ] **7d. Frontend**
  - Create `frontend/src/features/storage/components/file-upload.tsx` — drag & drop zone + click to browse, progress bar, preview
  - Create `frontend/src/features/storage/components/file-list.tsx` — grid/list view of files with thumbnails
  - Create `frontend/src/features/storage/components/file-card.tsx` — individual file card with actions (download, delete)
  - Create `frontend/src/features/storage/actions/storage.ts` — API calls
  - Create `frontend/src/features/storage/hooks/use-storage.ts`
  - Update `frontend/src/app/(dashboard)/files/page.tsx` — removed placeholder, compose real components
  - The files page in dashboard nav already exists in site config (`/files`)

- [ ] **7e. Wire in main.go**

- [ ] **7f. Test**: Storage service tests (mocked R2 + file repo): upload valid, upload too large, upload forbidden type, delete own file, delete other's file (should fail)

- [ ] **7g. Commit**: `feat: implement R2 file storage (upload, list, delete with drag-and-drop UI)`

---

### TASK 8: Gemini AI Chat
> Streaming responses via SSE (Server-Sent Events)

- [ ] **8a. Migration**
  - Create `backend/migrations/006_create_conversations.up.sql`:
    ```sql
    CREATE TABLE IF NOT EXISTS conversations (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        title TEXT NOT NULL DEFAULT 'New conversation',
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
    CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);

    CREATE TABLE IF NOT EXISTS messages (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
        role TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
    CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
    ```
  - Create matching `.down.sql`
  - Apply migration

- [ ] **8b. Backend — Gemini adapter**
  - Create `backend/internal/adapter/gemini/client.go` — Gemini API client setup using `GEMINI_API_KEY`
  - Create `backend/internal/adapter/gemini/ai.go` — implements `service.AIService`
  - Read existing `port/service/ai.go` and extend if needed:
    ```go
    type AIService interface {
        Chat(ctx context.Context, messages []ai.Message) (string, error)
        ChatStream(ctx context.Context, messages []ai.Message) (<-chan string, error)
    }
    ```
  - Go dependency: `go get github.com/google/generative-ai-go` (run from backend/)

- [ ] **8c. Backend — service + handler**
  - Flesh out `app/ai/service.go`: CreateConversation, SendMessage (saves to DB + calls Gemini), GetHistory, ListConversations, DeleteConversation
  - Create `adapter/postgres/conversation.go` — implements `ConversationRepository`
  - Create `handler/ai.go`:
    - `GET /api/ai/conversations` — list conversations
    - `POST /api/ai/conversations` — create new conversation
    - `GET /api/ai/conversations/:id/messages` — get messages
    - `POST /api/ai/conversations/:id/messages` — send message (returns full response)
    - `POST /api/ai/conversations/:id/stream` — send message (SSE streaming)
    - `DELETE /api/ai/conversations/:id` — delete conversation
  - For SSE: set `Content-Type: text/event-stream`, flush after each chunk

- [ ] **8d. Frontend**
  - Create `frontend/src/features/ai/components/chat-layout.tsx` — sidebar (conversation list) + main (chat area)
  - Create `frontend/src/features/ai/components/conversation-list.tsx` — list of past conversations, new conversation button
  - Create `frontend/src/features/ai/components/chat-messages.tsx` — message bubbles (user = right/primary, AI = left/muted)
  - Create `frontend/src/features/ai/components/chat-input.tsx` — text input + send button, handles Enter to send
  - Create `frontend/src/features/ai/hooks/use-chat.ts` — manages conversation state, handles SSE streaming
  - Create `frontend/src/features/ai/actions/ai.ts` — API calls
  - Update `frontend/src/app/(dashboard)/ai/page.tsx` — compose chat layout
  - Streaming: read SSE with `EventSource` or `fetch` + `ReadableStream`, display tokens as they arrive

- [ ] **8e. Wire in main.go**

- [ ] **8f. Test**: AI service tests (mocked Gemini): send message saves to DB, conversation ownership check

- [ ] **8g. Commit**: `feat: implement Gemini AI chat with streaming and conversation history`

---

### TASK 9: In-App Notifications

- [ ] **9a. Migration**
  - Create `backend/migrations/007_create_notifications.up.sql`:
    - Table: `notifications` (id UUID, user_id FK, type TEXT, title TEXT, message TEXT, read BOOLEAN DEFAULT false, data JSONB DEFAULT '{}', created_at)
  - Create matching `.down.sql`
  - Apply migration

- [ ] **9b. Backend**
  - Create `adapter/postgres/notification.go` — implements `NotificationRepository`
  - Flesh out `app/notification/service.go`: Send, MarkAsRead, MarkAllAsRead, GetUnread, List (paginated), GetUnreadCount
  - Create `handler/notification.go`:
    - `GET /api/notifications` — list (paginated), query param `?unread=true`
    - `GET /api/notifications/count` — unread count
    - `PUT /api/notifications/:id/read` — mark one as read
    - `PUT /api/notifications/read-all` — mark all as read

- [ ] **9c. Frontend**
  - Create `frontend/src/features/notification/components/notification-bell.tsx` — bell icon with unread badge count
  - Create `frontend/src/features/notification/components/notification-dropdown.tsx` — dropdown list of recent notifications
  - Create `frontend/src/features/notification/components/notification-list.tsx` — full page notification list
  - Create `frontend/src/features/notification/hooks/use-notifications.ts` — poll unread count every 30s
  - Create `frontend/src/features/notification/actions/notification.ts`
  - Add notification bell to `shared/components/layouts/dashboard-layout.tsx` header (import from feature in layout is OK since layout is in app/ layer)
  - Update `frontend/src/app/(dashboard)/notifications/page.tsx`

- [ ] **9d. Wire in main.go**

- [ ] **9e. Test**: Notification service tests: send, mark read, get unread count, list paginated

- [ ] **9f. Commit**: `feat: implement in-app notifications with bell, dropdown, and mark-as-read`

---

### TASK 10: Blog System (Backend + Public Frontend)
> Blog content is managed from admin panel (Task 11), but the backend API and public pages are built here

- [ ] **10a. Migration**
  - Create `backend/migrations/008_create_blog.up.sql`:
    ```sql
    CREATE TABLE IF NOT EXISTS blog_posts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        author_id UUID NOT NULL REFERENCES users(id),
        title TEXT NOT NULL,
        slug TEXT NOT NULL UNIQUE,
        excerpt TEXT NOT NULL DEFAULT '',
        content TEXT NOT NULL DEFAULT '',
        cover_image_url TEXT NOT NULL DEFAULT '',
        meta_title TEXT NOT NULL DEFAULT '',
        meta_description TEXT NOT NULL DEFAULT '',
        tags TEXT[] NOT NULL DEFAULT '{}',
        status TEXT NOT NULL DEFAULT 'draft',
        published_at TIMESTAMP,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
    CREATE UNIQUE INDEX IF NOT EXISTS idx_blog_posts_slug ON blog_posts(slug);
    CREATE INDEX IF NOT EXISTS idx_blog_posts_status ON blog_posts(status);
    CREATE INDEX IF NOT EXISTS idx_blog_posts_published_at ON blog_posts(published_at DESC);
    ```
  - Create matching `.down.sql`
  - Apply migration

- [ ] **10b. Backend — domain + service**
  - Create `backend/internal/domain/blog/post.go` — BlogPost entity with validation
  - Create `backend/internal/port/repository/blog.go` — BlogRepository interface
  - Create `backend/internal/app/blog/service.go` — CRUD, List (paginated, filterable by tag/status), GetBySlug
  - Create `backend/internal/adapter/postgres/blog.go` — implements BlogRepository

- [ ] **10c. Backend — handlers**
  - Create `backend/internal/handler/blog.go`:
    - `GET /api/blog/posts` — public, list published posts (paginated, filter by tag)
    - `GET /api/blog/posts/:slug` — public, get single post by slug
    - `GET /api/blog/tags` — public, list all tags with counts
    - `POST /api/admin/blog/posts` — admin only, create post
    - `PUT /api/admin/blog/posts/:id` — admin only, update post
    - `DELETE /api/admin/blog/posts/:id` — admin only, delete post
    - `GET /api/admin/blog/posts` — admin only, list ALL posts including drafts

- [ ] **10d. Frontend — public blog pages**
  - Create `frontend/src/features/blog/components/post-card.tsx` — card with cover image, title, excerpt, tags, date
  - Create `frontend/src/features/blog/components/post-content.tsx` — full post rendering (HTML content from API)
  - Create `frontend/src/features/blog/components/tag-filter.tsx` — clickable tags to filter posts
  - Create `frontend/src/features/blog/actions/blog.ts` — API calls
  - Update `frontend/src/app/(marketing)/blog/page.tsx` — blog listing
  - Create `frontend/src/app/(marketing)/blog/[slug]/page.tsx` — individual post page with SEO metadata

- [ ] **10e. Seed example posts** — Add 2-3 example blog posts in `cmd/seed/main.go` to demonstrate the blog

- [ ] **10f. Wire in main.go**

- [ ] **10g. Test**: Blog service tests: create, list published, get by slug, filter by tag

- [ ] **10h. Commit**: `feat: implement blog system (CRUD API, public pages, tags, SEO)`

---

### TASK 11: Admin Panel (Separate Vite App)
> This is a NEW app in `/admin/` directory — NOT part of the Next.js frontend

- [ ] **11a. Scaffold Vite app**
  ```bash
  cd /home/hassad/Documents/boilerplateSaaS
  npm create vite@latest admin -- --template react-ts
  cd admin && npm install
  npm install tailwindcss @tailwindcss/vite -D
  npm install react-router-dom
  ```
  Configure Tailwind (vite.config.ts + CSS), set up same design tokens as main frontend (copy CSS variables from frontend globals.css)

- [ ] **11b. Create admin/CLAUDE.md** with conventions for the admin app

- [ ] **11c. Admin app structure**
  ```
  admin/src/
  ├── components/     → Shared UI (Sidebar, Header, DataTable, Card)
  ├── pages/          → Dashboard, Users, Blog, Settings
  ├── hooks/          → useAuth, useApi
  ├── lib/            → API client (talks to same Go backend)
  └── types/          → TypeScript types
  ```

- [ ] **11d. Admin auth** — Login form, verify JWT, check user role is 'admin', redirect if not

- [ ] **11e. Dashboard page** — Stats cards (total users, active subscriptions, revenue, posts), PostHog analytics iframe/charts
  - Add PostHog config to `backend/internal/config/config.go` (`PostHogAPIKey`, `PostHogProjectID`, `PostHogHost`)
  - Backend proxy endpoint: `GET /api/admin/analytics` — fetches from PostHog API

- [ ] **11f. Users page** — Table with search, pagination, columns (name, email, role, status, created_at), actions (suspend, ban, make admin)

- [ ] **11g. Blog CMS page** — List of all posts (drafts + published), create/edit form with:
  - Title, slug (auto-generated from title, editable)
  - Rich text editor (use a simple one: `react-quill` or basic `contentEditable` with toolbar)
  - Cover image upload (to R2 via backend)
  - Tags (comma-separated or chips input)
  - SEO fields: meta title, meta description
  - Status toggle: Draft / Published
  - Publish date picker for scheduling

- [ ] **11h. Wire admin routes in backend** — all `/api/admin/*` routes check user role is admin

- [ ] **11i. Commit**: `feat: implement admin panel (Vite app with users, blog CMS, analytics dashboard)`

---

### TASK 12: Background Jobs & Cron

- [ ] **12a. Job runner infrastructure**
  - Create `backend/pkg/jobs/scheduler.go`:
    ```go
    type Job struct {
        Name     string
        Interval time.Duration
        Fn       func(ctx context.Context) error
    }
    type Scheduler struct { jobs []Job; stop chan struct{} }
    func NewScheduler() *Scheduler
    func (s *Scheduler) Register(job Job)
    func (s *Scheduler) Start(ctx context.Context)  // runs each job in its own goroutine with ticker
    func (s *Scheduler) Stop()
    ```

- [ ] **12b. Implement jobs**
  - Clean expired password reset tokens (every 1 hour)
  - Clean expired email verification tokens (every 1 hour)
  - Log system stats to slog (every 5 minutes) — active connections, goroutine count

- [ ] **12c. Wire in main.go** — create scheduler, register jobs, start in goroutine, stop on shutdown

- [ ] **12d. Test**: Test scheduler start/stop, test job execution

- [ ] **12e. Commit**: `feat: implement background job scheduler with token cleanup jobs`

---

### TASK 13: SEO

- [ ] **13a. Metadata for all pages**
  - Update every `page.tsx` in `app/` to export proper `metadata` or `generateMetadata`:
    - title, description, openGraph (title, description, image), twitter card
  - Marketing pages: static metadata
  - Blog posts: dynamic metadata from post data (title, excerpt, cover_image)
  - Dashboard pages: `robots: { index: false }` (don't index authenticated pages)

- [ ] **13b. sitemap.xml**
  - Create `frontend/src/app/sitemap.ts` — Next.js sitemap generation
  - Include: landing page, pricing, blog posts (fetch slugs from API), static pages
  - Exclude: dashboard, auth, admin pages

- [ ] **13c. robots.txt**
  - Create `frontend/src/app/robots.ts` — Next.js robots generation
  - Allow: /, /blog, /pricing
  - Disallow: /dashboard, /settings, /admin, /api

- [ ] **13d. JSON-LD structured data**
  - Landing page: `Organization` schema
  - Blog posts: `Article` schema (author, datePublished, dateModified)
  - Pricing: `Product` schema with `Offer`

- [ ] **13e. Commit**: `feat: implement SEO (metadata, sitemap, robots.txt, JSON-LD structured data)`

---

### TASK 14: Unit Tests (Comprehensive)
> Write tests for ALL domain entities and app services created in previous tasks

- [ ] **14a. Domain tests**
  - `domain/user/entity_test.go` — validation (empty email, empty name, valid creation)
  - `domain/billing/plan_test.go`, `subscription_test.go`, `invoice_test.go`
  - `domain/ai/conversation_test.go`, `model_test.go`
  - `domain/notification/notification_test.go`
  - `domain/storage/file_test.go`
  - `domain/blog/post_test.go` — slug generation, status validation

- [ ] **14b. App service tests** (all with mocked ports)
  - `app/auth/service_test.go` — register, login, forgot password, reset password, verify email
  - `app/user/service_test.go` — get profile, update profile, change password, delete account
  - `app/billing/service_test.go` — get plans, create checkout, handle webhook events
  - `app/storage/service_test.go` — upload, delete, list
  - `app/ai/service_test.go` — create conversation, send message, list conversations
  - `app/notification/service_test.go` — send, mark read, get unread count
  - `app/blog/service_test.go` — create post, list published, get by slug

- [ ] **14c. Generate mocks** — Use `mockgen` or manual mocks in `backend/mock/` for all port interfaces

- [ ] **14d. Verify coverage**: `go test ./... -cover` — target 80%+ on domain/ and app/ packages

- [ ] **14e. Commit**: `test: add comprehensive unit tests for all domain entities and app services`

---

### TASK 15: Playwright E2E Tests

- [ ] **15a. Setup Playwright**
  ```bash
  cd /home/hassad/Documents/boilerplateSaaS/frontend
  npm init playwright@latest
  ```
  Configure: `playwright.config.ts` with baseURL `http://localhost:3006`, test dir `e2e/`

- [ ] **15b. Test helpers**
  - Create `frontend/e2e/helpers.ts` — login helper, API seed helper, common selectors

- [ ] **15c. Write core E2E tests**
  - `e2e/auth.spec.ts` — register new user, login, see dashboard, logout
  - `e2e/settings.spec.ts` — update profile name, change password
  - `e2e/blog.spec.ts` — visit blog listing, open a post, verify SEO tags

- [ ] **15d. Commit**: `test: add Playwright E2E tests for auth, settings, and blog flows`

---

### TASK 16: Landing Page Update
> LAST TASK — update landing page to showcase all implemented features

- [ ] **16a. Update feature grid**
  - Update `frontend/src/features/marketing/components/features-section.tsx`
  - 12 cards with accurate descriptions of what's actually implemented:
    - Auth (email + password, JWT, route guards)
    - Billing (Stripe, 3 plans, checkout, webhooks)
    - AI Chat (Gemini, streaming, conversation history)
    - File Storage (R2, drag & drop, preview)
    - Email (Resend, templates, verification)
    - Notifications (in-app, bell, mark as read)
    - Blog CMS (DB-backed, admin editor, SEO)
    - Admin Panel (analytics, user management, blog CMS)
    - Background Jobs (Go native, scheduled tasks)
    - Security (rate limiting, input validation, bcrypt)
    - Architecture (hexagonal, feature-based, independently removable)
    - Developer Experience (CLAUDE.md, skills, tests, TypeScript strict)

- [ ] **16b. Create spotlight sections**
  - Create `frontend/src/features/marketing/components/spotlight-ai.tsx` — AI chat showcase
  - Create `frontend/src/features/marketing/components/spotlight-admin.tsx` — Admin panel showcase
  - Create `frontend/src/features/marketing/components/spotlight-architecture.tsx` — Architecture showcase

- [ ] **16c. Create comparison section**
  - Create `frontend/src/features/marketing/components/comparison-section.tsx`
  - Mini table: 4-5 rows comparing CleanSaaS vs OpenSaaS on key differentiators

- [ ] **16d. Update page composition**
  - Update `frontend/src/app/(marketing)/page.tsx`:
    ```
    Hero → FeatureGrid → SpotlightAI → SpotlightAdmin → SpotlightArchitecture → Stack → Comparison → DX → CTA
    ```

- [ ] **16e. Commit**: `feat: update landing page with feature grid, spotlights, and comparison`

---

## 3. AFTER ALL TASKS

- [ ] Run final check: `cd backend && go build ./... && go test ./... -count=1`
- [ ] Run final check: `cd frontend && npx tsc --noEmit`
- [ ] Run `/check` to verify architecture compliance
- [ ] Verify no hardcoded colors: search frontend for `zinc-`, `gray-`, `slate-`, `white`, `black` in className strings
- [ ] Verify no cross-feature imports in frontend
- [ ] Verify no `.env` files are committed: `git status`
- [ ] Update this file — all checkboxes should be `[x]`

---

## 4. QUICK REFERENCE

### File paths
- Backend: `/home/hassad/Documents/boilerplateSaaS/backend/`
- Frontend: `/home/hassad/Documents/boilerplateSaaS/frontend/`
- Admin: `/home/hassad/Documents/boilerplateSaaS/admin/`
- Migrations: `/home/hassad/Documents/boilerplateSaaS/backend/migrations/`

### Commands
```bash
# Backend
cd backend && make run                    # Start server (port 8081)
cd backend && go build ./...              # Compile check
cd backend && go test ./... -count=1      # Run tests
cd backend && make migrate-up             # Apply migrations
cd backend && make migrate-down           # Rollback last migration
cd backend && make seed                   # Seed data

# Frontend
cd frontend && npm run dev                # Start dev (port 3006)
cd frontend && npx tsc --noEmit           # Type check
cd frontend && npx playwright test        # E2E tests

# Docker
docker compose up -d                      # Start PostgreSQL + DbGate
```

### API base URL
- Backend: `http://localhost:8081`
- Frontend: `http://localhost:3006`
- DbGate: `http://localhost:8082`

### Current migration count
Check `ls backend/migrations/` — currently at `001_create_users`. Tasks add migrations 002-008.

### Design tokens quick ref
```
bg-primary text-primary-foreground     → main buttons, CTAs
bg-card border-border rounded-xl       → card containers
text-foreground                        → main text
text-muted-foreground                  → secondary text
bg-muted                               → subtle backgrounds
bg-destructive text-destructive-foreground → delete/error
bg-accent                              → hover states
```

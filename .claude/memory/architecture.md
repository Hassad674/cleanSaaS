# Architecture Details

## Backend Hexagonal Layers

### domain/ — Pure business (zero deps)
- user/: entity, role
- billing/: subscription, plan, invoice
- notification/: notification, template
- ai/: conversation, model
- storage/: file
- errors.go: typed domain errors

### port/ — Interfaces
- repository/: UserRepository, SubscriptionRepository, PlanRepository, InvoiceRepository, NotificationRepository, ConversationRepository
- service/: PaymentService, EmailService, AIService, StorageService, OAuthProvider

### app/ — Use cases
- auth: Register, Login, OAuthCallback
- billing: Subscribe, Cancel, ChangePlan, HandleWebhook
- user: GetProfile, UpdateProfile, DeleteAccount
- notification: Send, MarkRead, GetUnread
- ai: Chat, Stream, GetHistory
- storage: Upload, Delete, GetURL
- admin: Dashboard, Analytics, ManageUsers

### adapter/ — Implementations
- postgres/: DB connection, repositories (user, subscription, etc.)
- stripe/: PaymentService
- resend/: EmailService
- claude/: AIService
- openai/: AIService (same interface)
- google/: OAuthProvider
- r2/: StorageService

### handler/ — HTTP
- router.go: all routes
- auth.go, user.go, billing.go, etc.
- middleware/: auth (JWT), cors, ratelimit, logging
- dto/: request/ and response/ objects

### pkg/ — Public utilities
- jwt: Generate, Validate
- hash: Password, Check (bcrypt)
- validate: Email, MinLength, Slug
- pagination: FromRequest

## Frontend Feature-based
- app/: routes only (thin)
- features/: self-contained modules
- shared/: components (shadcn), hooks, lib, types
- config/: app config

## Request Flow
```
HTTP → handler → app (use case) → domain → port ← adapter (DB/external)
```

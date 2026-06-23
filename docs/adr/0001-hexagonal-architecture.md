# 0001. Hexagonal (ports & adapters) backend architecture

## Status

Accepted

## Context

CleanSaaS is a boilerplate for medium-to-large SaaS products, not a micro-SaaS
starter. The people who clone it fall into two groups with opposite needs:

- **Beginners** who need a single, predictable structure they can learn once and
  reuse everywhere, and who will lean on AI agents to extend it.
- **Professionals** who need a core that is independently testable and infra that
  can be swapped (Stripe → Lemon Squeezy, Postgres → another store, Gemini →
  Claude) without a rewrite.

Three forces dominate:

1. **Swappable infrastructure.** The boilerplate ships opinionated providers
   (Stripe, Resend, Gemini, Cloudflare R2, Neon Postgres), but no clone will keep
   all of them. Provider choice must be a leaf detail, not a structural one.
2. **A testable core.** Business rules must be unit-testable in milliseconds,
   with no database and no network. If testing a use case requires Docker, agents
   cannot run the Test→Fix→Retest loop fast enough to self-correct.
3. **Consistency for AI agents.** Agents extend this code. A uniform shape —
   "domain first, HTTP last, dependencies always injected" — means an agent that
   has seen one feature can correctly scaffold the next. Inconsistency is where
   agent-written code rots.

A conventional framework-centric layout (handlers calling an ORM directly)
fails forces 1 and 2: the database and the web framework leak into business
logic, making the core untestable in isolation and the infra un-swappable.

## Decision

We will structure the backend as a **hexagonal (ports and adapters)**
architecture with a strict, one-directional dependency rule:

```
handler → app → domain ← port ← adapter
```

- **domain/** — pure business entities, value objects, and rules. Imports the Go
  standard library only. Entities validate themselves.
- **port/** — interface contracts (`port/repository`, `port/service`). Imports
  `domain` only. Defines *what* the core needs, never *how*.
- **app/** — use cases / application services, one sub-package per feature
  (`auth`, `billing`, `user`, `ai`, …). Imports `domain` + `port` interfaces.
  Receives all dependencies via constructor injection. Returns domain types and
  domain errors — never HTTP concepts.
- **adapter/** — concrete implementations of ports (`postgres`, `stripe`,
  `resend`, `gemini`, `r2`, …). Imports `domain`, `port`, and external libs. An
  adapter never imports another adapter.
- **handler/** — HTTP transport (Chi router, DTOs, middleware). Thin: decode
  request → call app service → map domain error → encode response.

All dependency injection happens in exactly one place: `cmd/api/main.go` (the
composition root). Swapping a provider is a new adapter file plus one changed
line there.

## Consequences

**Positive (+)**

- The core (`domain` + `app`) is unit-testable with mocked ports — no DB, no
  network — which makes the agent self-correction loop fast.
- Infrastructure is genuinely swappable: implement the port interface in a new
  adapter, change one line in `main.go`, done (see ADR 0002 and the `/add-adapter`
  skill).
- The dependency direction is mechanically checkable, so it can be enforced in CI
  rather than trusted (see ADR 0005 and `scripts/ci/check-cross-feature-imports.sh`,
  which verifies domain purity).
- Uniform shape across features makes the codebase legible to AI agents and
  newcomers alike.

**Negative (−)**

- More indirection than a direct handler→DB call: every feature carries a port
  interface and an adapter even when there is only one implementation.
- More files per feature (entity, port, service, adapter, handler, DTOs). This is
  intentional friction, but it is friction.
- Beginners must learn the dependency rule before they are productive. The
  per-directory `CLAUDE.md` files and the `/add-feature` scaffolder exist to
  flatten that curve.

## Alternatives considered

- **Layered / MVC (controller → service → repository, no ports).** Familiar and
  fewer files, but services typically import concrete repositories, so the core
  is not swappable and tests need real (or heavily faked) infrastructure.
  Rejected: fails the testable-core and swappable-infra forces.
- **Clean Architecture (Uncle Bob, with explicit use-case interactors and
  request/response models per use case).** Achieves the same dependency
  inversion but adds ceremony (boundary models, interactor objects) that is
  overkill for a boilerplate and harder for beginners. Hexagonal gives us the
  same isolation with a smaller vocabulary. Rejected as over-engineered for the
  audience; we keep its dependency-inversion spirit.
- **Framework-coupled (handlers calling an ORM / query builder directly).**
  Fastest to write initially, but the framework and database leak everywhere,
  making the core untestable in isolation and providers impossible to swap
  cleanly. Rejected as the exact failure mode this ADR exists to prevent.

This decision is complementary to ADR 0004: hexagonal governs the *direction* of
dependencies; DDD governs *where behavior lives* inside the `domain` layer.

# 0003. Every feature is fully removable

## Status

Accepted

## Context

CleanSaaS is a boilerplate, and we do not know what its users will build. One
clone becomes a CRM, another an e-commerce platform, another a social app. None
of them uses every module we ship. Today the repo includes all features; the
roadmap includes a `create-cleansaas` CLI that lets users pick only the modules
they want.

For that future to be possible — and for the codebase to stay clean today —
features cannot quietly entangle themselves. The moment `billing` imports
`user`, or the `conversations` table gains a foreign key to `subscriptions`,
removing AI without removing billing becomes impossible, and the boilerplate's
core promise breaks.

Only two things are truly core: **auth** (there are no users without it) and the
**database connection / config**. Everything else — billing, storage, ai,
notification, blog, team — must be optional and independently extractable.

The challenge is that this kind of discipline erodes silently. A single
"convenient" cross-import, merged once, is enough. So the rule cannot live only
in prose; it must be mechanically verifiable (see ADR 0005).

## Decision

We will treat **full removability** as a hard architectural invariant. A feature
is removable when deleting its folders (frontend + backend) **and** its wiring
lines in `cmd/api/main.go` produces **zero compilation errors** elsewhere and the
app still runs.

Concretely:

1. **No cross-feature imports.** A package under `internal/app/<A>` must not
   import `internal/app/<B>`; an adapter must not import another adapter. If
   `billing` needs user data, it depends on a `UserRepository`-style *interface*
   injected at the composition root — never on the `user` package.
2. **No cross-feature foreign keys.** A feature's tables may only
   `REFERENCES users(id)`. `subscriptions` may reference `users`, but
   `conversations` must never reference `subscriptions`. Each feature's tables
   are self-contained and droppable.
3. **Explicit wiring only.** Every feature is wired by hand in
   `cmd/api/main.go`. No auto-discovery, no magic registration. Removing a
   feature = deleting its lines there.
4. **Frontend composition lives in `app/` pages.** Features never import each
   other; a page that needs both profile and billing imports from both — the
   features stay mutually ignorant.
5. **Feature-scoped migrations and server actions.** Each feature's tables get
   their own numbered migration files; each feature's `actions/` only calls its
   own endpoints.

This is **enforced, not trusted**:

- `scripts/ci/check-cross-feature-imports.sh` fails the build on any
  `app/<A> → app/<B>`, `adapter → adapter`, or domain-purity violation
  (frontend `@/features/<B>` cross-imports too).
- The `/check` skill audits a feature (or the whole repo) against the rules.
- The `/verify-independence` skill proves a feature is removable **without
  mutating the repo** — a dry run of the deletion.

## Consequences

**Positive (+)**

- The `create-cleansaas` "pick your modules" CLI is achievable, because removal
  is already a clean operation.
- Features can be developed, reasoned about, and reviewed in isolation; a bug in
  `blog` cannot reach `billing`.
- New contributors and AI agents get a bright-line rule with an automated check,
  so "did I just couple two features?" is answered by CI, not by hoping a
  reviewer notices.
- Database tables map cleanly to features, so dropping a feature's schema is safe.

**Negative (−)**

- Sharing logic between features requires defining an interface and injecting it
  at the composition root, rather than a direct import — more ceremony for genuine
  cross-cutting needs (broadly reusable value objects are the escape hatch; see
  ADR 0004).
- `cmd/api/main.go` carries all wiring and grows with the feature count
  (mitigated by a typed `Deps` struct rather than long positional argument lists).
- Some duplication is accepted across features rather than coupling them; the
  "extract only when used 3+ times" rule keeps this from getting out of hand.
- Contributors must internalize the no-cross-import / no-cross-FK rules up front.

## Alternatives considered

- **Shared "common" / "core" service package that features import freely.** The
  usual outcome is a god-module everything depends on, which is the opposite of
  removable. Rejected. We allow only genuinely generic, dependency-free utilities
  in `pkg/` and broad value objects in `domain`.
- **A plugin / auto-registration system** (features register themselves at
  startup via `init()`). Hides the dependency graph and makes removal a
  guessing game. Rejected in favor of explicit wiring in `main.go`.
- **Documentation-only convention (no CI gate).** This is how modularity rules
  usually die — one convenient import at a time. Rejected; mechanical enforcement
  is non-negotiable here (ADR 0005).

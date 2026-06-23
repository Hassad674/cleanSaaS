# 0004. Targeted DDD — rich domain only where complexity warrants it

## Status

Proposed / partially Accepted (roadmap — not yet fully implemented across all
features)

## Context

ADR 0001 gives us a `domain` layer that imports nothing but the standard
library. That tells us *where* business logic lives; it does not tell us *how
rich* the domain model should be. Two failure modes are equally bad for a
boilerplate:

- **Over-engineering.** Wrapping a simple blog post in aggregates, value objects,
  and domain events teaches beginners ceremony with no payoff, and bloats files
  that should be trivial. A boilerplate used by beginners must not look like a
  DDD textbook for CRUD.
- **Under-modeling.** Letting genuinely complex features — billing (plans,
  subscriptions, proration, webhook-driven state) and team (membership, roles,
  invitations) — stay as anemic structs pushes their rules up into `app`
  services, which then grow long, branchy, and hard to test exactly where
  correctness matters most.

So the question is not "DDD or not" but "*where* does richer modeling earn its
keep?" We also need to be explicit that DDD and hexagonal are **complementary,
not alternatives**: hexagonal governs dependency *direction* (ADR 0001); DDD
governs *where behavior lives* inside the `domain` layer that hexagonal already
isolates.

## Decision

We will apply **DDD selectively**, matched to each feature's complexity:

- **Rich domain modeling for complex features.** `billing` and `team` model
  their cores as **aggregates** with invariants enforced inside the entity,
  **domain events** for state changes worth reacting to, and **value objects**
  for the concepts they own.
- **Broadly reused value objects** that cut across features live in `domain` and
  may be imported anywhere: `Email`, `Money`, `Slug`. These are
  dependency-free, self-validating, and are the sanctioned exception to "features
  don't share code" (ADR 0003) — they are domain primitives, not feature logic.
- **Intentionally anemic models for simple CRUD features.** `blog`,
  `notification`, `storage`, and the persistence side of `user` stay as plain
  validated entities with their logic in `app` services. We do **not** force
  aggregates or domain events onto them.

The test for "rich vs anemic": *does this feature have invariants or state
transitions that are wrong to express as procedural steps in a service?* If yes,
model it richly. If it is fundamentally "validate fields, store, fetch," keep it
anemic.

DDD is layered **on top of** hexagonal, not instead of it. A billing aggregate
still sits in `domain`, is still orchestrated by `app/billing`, is still
persisted through a `port/repository` interface, and is still reachable only via
a thin `handler`. Nothing about the dependency rule changes.

## Consequences

**Positive (+)**

- Complexity lives where the complexity is. Billing's and team's hardest rules
  are unit-testable inside the domain, not buried in branchy service code.
- Beginners are not taxed with aggregates and events for a blog post; simple
  features stay simple and readable.
- Shared value objects (`Email`, `Money`, `Slug`) remove a whole class of
  validation bugs and duplication without coupling features to each other.
- DDD and hexagonal reinforce each other: clear dependency direction *and* a
  clear home for behavior.

**Negative (−)**

- **Two idioms in one codebase.** Some features are rich, others anemic — a
  reader must recognize which is which. We accept this as honest, not accidental:
  the README index and each feature's `CLAUDE.md` should state which idiom it
  uses.
- Judgment is required to classify a new feature, and judgment can be wrong; a
  feature can start anemic and need promotion later (or vice versa).
- **Partially implemented.** As of this ADR, not every complex feature fully
  expresses aggregates/events yet — this is a roadmap commitment, which is why
  the status is Proposed/partially Accepted rather than fully Accepted.

## Alternatives considered

- **Rich DDD everywhere.** Maximal consistency, but it over-engineers a
  boilerplate and intimidates the beginner audience for no benefit on CRUD.
  Rejected.
- **Anemic everywhere (transaction-script style, all logic in services).**
  Simplest to learn, but billing and team service files would become long,
  deeply branched, and hard to test — exactly the smell our code-quality limits
  forbid. Rejected.
- **DDD as a replacement for hexagonal** (treating them as competing
  architectures). A category error: they answer different questions. Rejected in
  favor of combining them.

Supersedes nothing. Refines ADR 0001 by specifying how the `domain` layer is
populated.

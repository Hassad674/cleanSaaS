# Architecture Decision Records (ADRs)

This directory records the significant architectural decisions made for CleanSaaS:
why the backend is hexagonal, why there is no ORM, why every feature must be
removable, and so on. The goal is that a newcomer (human or AI agent) can read
these records and understand *why* the codebase looks the way it does — not just
*what* the rules are (those live in the various `CLAUDE.md` files).

## What is an ADR?

An **Architecture Decision Record** is a short document that captures a single
architecturally significant decision, the context that forced it, and the
consequences of choosing it. ADRs are immutable history: once a decision is
**Accepted**, you do not rewrite it — if you change your mind later, you write a
*new* ADR that supersedes the old one. This gives the project a paper trail of
how its architecture evolved and why.

The format used here is the one popularized by Michael Nygard (2011).

## The Nygard format

Each ADR is a single Markdown file with these sections, in this order:

| Section | Purpose |
|---|---|
| **Title** | `# NNNN. Short imperative phrase` — names the decision, not the problem. |
| **Status** | One of the lifecycle states below. |
| **Context** | The forces at play: requirements, constraints, the problem we must solve. Written so the decision feels inevitable by the end. No solution here. |
| **Decision** | The choice we made, stated in active voice: "We will…". |
| **Consequences** | What becomes easier (**+**) and harder (**−**) as a result. Be honest — every decision has a cost. |
| **Alternatives considered** | The other options we weighed and why we rejected them. This is what stops a future reader (or AI agent) from re-litigating a settled decision. |

Keep each ADR to roughly one page. If it is sprawling, the decision is probably
two decisions — split it.

## Status lifecycle

```
Proposed ──► Accepted ──► Deprecated
                 │
                 └──────► Superseded by NNNN
```

- **Proposed** — drafted, under discussion, not yet binding.
- **Accepted** — the decision is in force; the codebase should reflect it.
- **Deprecated** — no longer the recommended approach, but not replaced by a
  specific newer decision.
- **Superseded by NNNN** — replaced by a later ADR. Link to it. The superseded
  ADR stays in the repo (history is never deleted), but its Status line points
  forward to the replacement.

An ADR may also note **partially implemented (roadmap)** alongside its status
when the decision is agreed but the code has not fully caught up yet (see ADR
0004).

## Numbering and file names

- Files are named `NNNN-kebab-title.md`, e.g. `0001-hexagonal-architecture.md`.
- `NNNN` is a zero-padded, monotonically increasing integer. Never reuse a
  number, even if an ADR is superseded.
- Numbers are allocated in the order ADRs are *written*, not by topic.

## How to add an ADR

1. Pick the next free number (`ls docs/adr` and take the highest + 1).
2. Copy the structure of an existing ADR into `NNNN-your-title.md`.
3. Fill in every section. Start the Status at **Proposed**.
4. Open a PR. Discussion happens on the PR.
5. When the team agrees, change the Status to **Accepted** and merge.
6. If a later decision overturns this one, write a new ADR and set this one's
   Status to **Superseded by NNNN** with a link.

> Rule of thumb: write an ADR whenever a decision is **costly to reverse** or
> **non-obvious to a newcomer**. Routine, easily-reversed choices do not need one.

## Index

| ADR | Title | Status |
|---|---|---|
| [0001](0001-hexagonal-architecture.md) | Hexagonal (ports & adapters) backend architecture | Accepted |
| [0002](0002-pure-sql-no-orm.md) | Pure SQL with `database/sql`, no ORM | Accepted |
| [0003](0003-feature-modularity-removability.md) | Every feature is fully removable | Accepted |
| [0004](0004-targeted-ddd.md) | Targeted DDD — rich only where complexity warrants it | Proposed / partially Accepted |
| [0005](0005-enforcement-as-product.md) | Mechanical enforcement is the product | Accepted |

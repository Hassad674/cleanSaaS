# 0005. Mechanical enforcement is the product

## Status

Accepted

## Context

Generating code is no longer the hard part. Anyone with an AI agent can produce a
working SaaS feature in an afternoon. What almost no one gets for free is code
that *stays* maintainable, performant, and secure after dozens of agent-driven
changes — code where the architecture still holds six months and two hundred
commits later.

That gap is CleanSaaS's reason to exist. Our differentiator is not the stack
(Next.js + Go + Postgres is unremarkable) and not the features (auth, billing,
ai, storage… everyone has these). It is that the boilerplate **mechanically
prevents the decay** that normally follows AI-assisted development.

The problem with conventions documented in prose is that AI agents and humans
both drift from them under deadline pressure, and a convention that is only
*described* is a convention that is *eventually ignored*. The architecture
decisions in ADRs 0001–0004 — the dependency rule, no-ORM SQL hygiene,
no-cross-feature coupling, where behavior lives — are only worth anything if a
violation *fails the build*, not merely earns a frown in review.

## Decision

We will treat **enforcement as the product**, and make architectural rules
mechanically checkable wherever possible. The enforcement system has four layers:

1. **`CLAUDE.md` files as a machine-and-human contract.** Root, `backend/`,
   `frontend/`, and `admin/` each carry a `CLAUDE.md` plus `.claude/memory/`
   stating the invariants. These are read by AI agents on every task and by
   humans on day one — the same source of truth for both.
2. **Skills that scaffold correctly by construction.** `/add-feature`,
   `/add-endpoint`, `/add-adapter`, `/add-migration`, `/remove-feature`,
   `/verify-independence` generate code that already obeys the rules (domain
   first, ports injected, migrations paired), so the easy path is the correct
   path.
3. **CI gates that fail on violations.** `scripts/ci/` enforces what a linter
   cannot express: `check-cross-feature-imports.sh` (ADR 0003 isolation +
   domain purity), `check-file-length.sh`, `check-forbidden-names.sh`,
   `check-hardcoded-colors.sh`, `check-migration-pairs.sh`, run together by
   `run-all.sh`, alongside `golangci-lint` for the code-quality caps.
4. **Tests and git hooks as guardrails.** TDD with the Test→Fix→Retest loop
   means agents self-correct against failing tests; the validation pipeline
   (`go build` · `go test` · `tsc --noEmit` · `vitest` · Playwright) gates every
   commit; git hooks run the cheap checks before code ever lands.

The principle: **every rule that matters is enforced by a machine, not by trust.**
A rule that cannot be checked mechanically is documented as a reviewer
responsibility and is a candidate for a future automated check.

## Consequences

**Positive (+)**

- AI-generated code stays within the architecture, because stepping outside it
  fails the build — the system catches drift the moment it happens.
- The same guardrails serve beginners (who get correctness they could not yet
  enforce themselves) and pros (who get a fast, trustworthy feedback loop).
- The rules in ADRs 0001–0004 are load-bearing rather than aspirational; this
  ADR is what makes the others real.
- Onboarding is faster: the build tells you when you are wrong, immediately and
  specifically (GitHub annotations point at the offending file and line).

**Negative (−)**

- The enforcement system is itself code that must be built, tested, and
  maintained; a buggy gate is worse than no gate (a false failure blocks work, a
  false pass erodes trust). The checks need their own tests
  (`scripts/ci/__tests__`).
- Stricter gates mean more red builds, which can frustrate newcomers until they
  learn the rules — the cost we accept for keeping the codebase clean.
- Some valuable rules resist mechanical checking (e.g. "is this the right
  abstraction?"); those still depend on human review, and the line between the
  two must be kept honest.
- Maintenance burden as the stack evolves: tool upgrades can break gates that
  parse imports or file structure.

## Alternatives considered

- **Documentation-only conventions (style guide, wiki, `CONTRIBUTING.md`).** Zero
  enforcement cost, but conventions without teeth decay — especially under
  AI-driven change velocity. This is precisely the failure mode the project
  exists to solve. Rejected.
- **Code review as the sole gatekeeper.** Humans miss things, review fatigue is
  real, and a solo user with an AI agent has no reviewer at all. Review is kept as
  a complement (the `/review` skill) but not relied on as the enforcement layer.
  Rejected as the *only* mechanism.
- **Linting alone (`golangci-lint`, ESLint).** Necessary but insufficient:
  linters enforce style and some complexity caps but cannot express
  "no cross-feature import" or "every migration has a down file." We keep linting
  and add purpose-built gate scripts on top. Rejected as a complete solution.

This ADR is the foundation that makes ADRs 0001–0004 enforceable rather than
merely recommended.

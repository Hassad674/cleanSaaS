# AUTONOMY-LOG — Autonomous upgrade of CleanSaaS

> **This file is the single source of truth for the autonomous upgrade session.**
> It survives context compaction. If you (the agent) just woke up / were compacted:
> 1. Read this whole file.
> 2. Run `git branch --show-current` and `git log --oneline -15`.
> 3. Run `cd backend && go build ./... && go test ./...` and `cd frontend && npx tsc --noEmit && npx vitest run` to confirm green.
> 4. Resume from the first `[ ]` or `[~]` item in the current phase below.

## Mission

Elevate CleanSaaS from "very good boilerplate" to genuinely world-class, matching the bar set by the
owner's best project (DesignedTrust-Services). Backend is the priority; frontend only needs good demo UX.
The real product = the **enforcement system** (CLAUDE.md + memory + skills + CI + hooks) that forces
AI-generated code to stay maintainable/performant/secure, usable by beginners AND pros.

## Standing decisions (do not re-litigate)

- **Autonomy scope:** MAXIMUM — including behavior-changing backend P0 + targeted DDD. Owner reviews diff on return.
- **Remove referral + mobile:** YES, both. referral via `/remove-feature` logic; mobile = delete folder. Reversible via git.
- **Delivery:** stacked git branches, one per phase, atomic conventional commits. **DO NOT push** (no auth; owner pushes after review).
- **Verification gate:** nothing is `[x]` until it passes `go build`, `go test`, `tsc`, `vitest` (and Playwright when relevant). Keep the app runnable.
- **DDD verdict:** PARTIAL — rich aggregates for billing + team only; value objects broadly (Email/Money/Slug); skip aggregate ceremony for simple CRUD (blog/notification/storage/user).
- **Language:** respond to owner in French; code/comments/docs in English.

## Branch plan (stacked — each branches off the previous phase tip)

| Phase | Branch | Base |
|---|---|---|
| 0 | `phase-0-calibration` | `main` |
| 1 | `phase-1-enforcement` | `phase-0-calibration` |
| 2 | `phase-2-hygiene` | `phase-1-enforcement` |
| 3 | `phase-3-backend-p0` | `phase-2-hygiene` |
| 4 | `phase-4-ddd` | `phase-3-backend-p0` |

## Current state

- **Active phase:** 0 — Calibration (COMPLETE, committing)
- **Active branch:** `phase-0-calibration`
- **Baseline:** GREEN (backend build+tests ok after fixing a pre-existing mock bug; tsc ok; vitest 25/25).
- **Last action:** finished calibration (memory, skills, CLAUDE.md). Committing phase 0, then opening phase 1.
- **Extra bugs found (to fix in phase 2 docs/hygiene):** Playwright config baseURL :3006 vs dev :3010; e2e specs log in as `admin@cleansaas.com` but seed creates `admin@cleansaas.dev`.

---

## Phase 0 — Calibration (CLAUDE.md + memory + skills) — PRIORITY, do first

Foundation. A well-calibrated agent produces correct code; everything downstream depends on this.

- [x] Establish green baseline (fixed `mockPaymentSvc.RetrieveCheckoutSession` missing method)
- [x] Create AUTONOMY-LOG.md + branch strategy
- [x] **Memory sanitation** — MEMORY.md rewritten (removed sudo/personal/language lines; added module map + env-gating, ports, seeded creds, architecture invariants, skills index); architecture.md corrected to reality (adapters gemini/postgres/r2/resend/stripe; all features listed).
- [x] **Fix existing skills** (8): stale paths removed (0 left); `test` skill rewritten (Vitest + Playwright, no Jest).
- [x] **Upgrade CLAUDE.md** (root + admin): broken refs fixed (BLOCKED.md, TACHES.md §0.3/§1.4); hard numeric limits added; runtime/observability section added; Next 16 not 15; Compact-instructions point at AUTONOMY-LOG.md; skills table updated with /run /debug /e2e /verify-independence; admin/CLAUDE.md testing section added.
- [x] **New skills:** `/run`, `/debug`, `/e2e`, `/verify-independence` created (frontmatter valid).
- [~] Commit + create phase-1 branch (in progress).

## Phase 1 — Mechanical enforcement (CI + hooks + gates)

The leap toward the bar. Mostly additive (no app behavior change).

- [ ] golangci-lint config (`.golangci.yml`)
- [ ] GitHub Actions: `ci.yml` (fast: go vet/lint/build/test -race + per-package coverage; node tsc/vitest/build), with `ci-gate` aggregate job + `concurrency` auto-cancel
- [ ] GitHub Actions: `e2e.yml` (label-gated `run-e2e`, real Postgres container, Playwright)
- [ ] Custom gate scripts in `scripts/ci/`: cross-feature import lint (backend+frontend), migration up/down pair check, no-hardcoded-colors check
- [ ] Meta-tests for the gate scripts (prove they fail on violations and don't false-positive)
- [ ] `.githooks/pre-commit` (zero-dep bash: gofmt + tsc + staged-only) + `scripts/install-git-hooks.sh` + meta-test
- [ ] Dependabot config (grouped, ecosystem-aware)
- [ ] Commit + create phase-2 branch.

## Phase 2 — Hygiene & modularity demo

- [ ] Remove `referral` (backend slice + frontend feature + demo page + migration 009) — verify build green = modularity proof
- [ ] Remove `mobile/` folder entirely
- [ ] Reconcile docs: FEATURES.md (Flutter not RN; referral/mobile removed), README, delete or fix TACHES.md
- [ ] OSS files: SECURITY.md (SLA), CONTRIBUTING.md, CHANGELOG.md (Keep-a-Changelog), CODE_OF_CONDUCT.md, `.claudeignore`, `.gitignore` tuning (ignore internal agent/audit docs)
- [ ] `docs/ARCHITECTURE.md` (+ mermaid), `docs/adr/` (Nygard ADRs for key decisions), `docs/ops.md` runbook
- [ ] Commit + create phase-3 branch.

## Phase 3 — Backend P0 correctness (behavior-changing)

- [ ] Real transactions: `DBTX` interface; repos accept it; wrap multi-write use cases (team+owner, subscription+invoice, conversation+message) in `WithTx`
- [ ] JWT: short access token + refresh token (DB-stored, rotation) + `session_version`/revocation; `iss`/`aud`; `config.Validate()` fails boot on weak/default secret in non-dev
- [ ] Stripe webhook idempotency (processed_events table, ON CONFLICT) + return 5xx on transient errors
- [ ] Context timeouts on DB queries + external calls; stop discarding ctx in Stripe adapter
- [ ] Fix IDOR: team.Get/GetTeamBySlug verify membership
- [ ] Panic recovery in job scheduler goroutines
- [ ] Optimistic locking (version column) on mutable aggregates + consistent RowsAffected checks
- [ ] Commit + create phase-4 branch.

## Phase 4 — Targeted DDD + integration tests

- [ ] Value objects: Email, Money/AmountCents, Slug, PlanInterval (with tests)
- [ ] Rich billing aggregate (Subscription.Activate/Renew/ChangePlan/Cancel — invariants on the aggregate)
- [ ] Rich team aggregate (Invite/RemoveMember/ChangeRole authorization matrix as domain methods)
- [ ] Domain events (collected on aggregate, dispatched post-commit) — also fixes fire-and-forget emails
- [ ] Thin out app/billing + app/team services accordingly
- [ ] Integration tests: testcontainers + `//go:build integration` for postgres adapters
- [ ] Commit. Done — summarize for owner.

---

## Blockers log

(none yet — use BLOCKED-<phase>.md for anything stuck after 3 attempts, and note it here)

## Decisions / discoveries during work

- 2026-06-22: Pre-existing baseline bug — `mockPaymentSvc` missing `RetrieveCheckoutSession` made `go test ./...` red. Fixed. (Evidence CI is needed.)

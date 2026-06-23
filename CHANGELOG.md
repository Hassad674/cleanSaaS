# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

This cycle focuses on turning the boilerplate's architecture rules from *documented*
into *mechanically enforced*, plus OSS maturity and module cleanup.

### Added
- **CI pipelines** (`.github/workflows/`): a fast `ci.yml` (backend build/vet/golangci/test
  `-race` + coverage floor; frontend tsc/vitest/build; admin tsc/build; invariant gates;
  aggregate `ci-gate` check) and a label-gated `e2e.yml` with a real Postgres service.
- **Invariant gate scripts** (`scripts/ci/`) with meta-tests: cross-feature/layer import
  isolation, migration up/down pairing, hardcoded-color detection, file-length cap,
  forbidden-identifier check. `run-all.sh` aggregates them.
- **golangci-lint config** (`backend/.golangci.yml`) mapping the hard code-quality limits
  to linters, plus SQL-correctness linters (bodyclose, rowserrcheck, sqlclosecheck).
- **Git pre-commit hook** (`.githooks/pre-commit`) — zero-dependency bash: staged-only
  gofmt + secret guard + best-effort vet/tsc, with an installer and a meta-test.
- **New Claude Code skills**: `/run` (start the whole stack), `/debug` (guided
  beginner reproduce→report→fix→verify with screenshots + Chrome extension), `/e2e`
  (Playwright), `/verify-independence` (prove a module is removable in a throwaway worktree).
- **Documentation**: `docs/ARCHITECTURE.md` (with Mermaid diagrams), `docs/adr/`
  (Nygard-format ADRs 0001–0005), `docs/ops.md` runbook.
- **OSS hygiene**: `SECURITY.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, this
  `CHANGELOG.md`, `.claudeignore`, Dependabot config.
- **Dependabot** (`.github/dependabot.yml`): grouped gomod / npm (frontend + admin) /
  github-actions updates, with React/Next/Tailwind majors split out.

### Changed
- Calibrated the AI-agent layer: `CLAUDE.md` now declares hard numeric code-quality
  limits, a runtime/observability section, and corrected references; project memory was
  sanitized (machine-personal/secret lines removed) and given a module map with env-gating.
- Rewrote the `/test` skill for the real stack (Vitest + Playwright; removed dead Jest commands).
- Fixed stale absolute paths across 6 skills (now relative / `$CLAUDE_PROJECT_DIR`).

### Removed
- **referral** module — removed across all layers (backend slice + migration 009 + frontend
  feature/pages/demo) as a demonstration of the "every feature is fully removable" promise.
- **mobile** (Flutter) app — removed entirely; the backend is a generic REST/WS API with no
  mobile-specific code, so removal has zero backend impact.

### Fixed
- Pre-existing test build failure: the billing `mockPaymentSvc` was missing
  `RetrieveCheckoutSession`, which made `go test ./...` fail on a fresh clone.
- Playwright config used port `:3006` while the dev server runs on `:3010`; e2e specs
  logged in as `admin@cleansaas.com` while the seed creates `admin@cleansaas.dev`.

[Unreleased]: https://github.com/Hassad674/cleanSaaS/commits/main

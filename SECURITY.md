# Security Policy

CleanSaaS is an open-source **boilerplate**. This document covers the security of
the boilerplate code itself — the patterns, defaults, and dependencies we ship.

> **Scope, stated plainly.** CleanSaaS gives you a hardened starting point, not a
> finished product. Once you deploy a fork, **you** own its security: your secrets,
> your infrastructure, your access controls, your data, your patches. A vulnerability
> in *your* deployment that does not stem from a flaw in the upstream boilerplate is
> not something we can fix for you. Vulnerabilities in the **upstream boilerplate
> code** are exactly what we want to hear about.

---

## Reporting a vulnerability

**Do not open a public GitHub issue for security problems.** Public disclosure before
a fix is available puts every downstream user at risk.

Report privately through one of these channels:

1. **Email** — `security@cleansaas.dev`
   *(Maintainers of a fork: replace this with your own monitored address. The
   placeholder will not reach you.)*
2. **GitHub Security Advisories** — use the repository's
   **Security → Report a vulnerability** tab to open a private advisory.

Please include, as much as you can:

- A description of the vulnerability and its impact.
- The affected component (e.g. `backend/internal/app/auth`, a migration, a dependency).
- Step-by-step reproduction, ideally with a minimal proof of concept.
- Affected version / commit SHA.
- Any suggested remediation.

You will receive an acknowledgement, and we will keep you informed through triage,
fix, and disclosure. If you wish, we will credit you in the advisory and the
`CHANGELOG.md`.

---

## Response SLA

These are the targets we hold ourselves to once a report is received. "Business days"
are Monday–Friday.

| Stage                  | Target                                            |
|------------------------|---------------------------------------------------|
| Acknowledge receipt    | within **72 hours**                               |
| Initial triage         | within **5 business days**                        |
| Fix — **CRITICAL**     | within **14 days** of confirmed triage            |
| Fix — **HIGH**         | within **30 days** of confirmed triage            |
| Fix — **MEDIUM**       | within **90 days** of confirmed triage            |
| Fix — **LOW**          | best effort, typically bundled into a later release |

Severity follows [CVSS v3.1](https://www.first.org/cvss/). If a report's severity is
disputed, we will explain our reasoning in the advisory thread.

---

## Supported versions

Security fixes land on the latest minor release line. Older lines receive fixes only
for **CRITICAL** issues, and only until the next minor supersedes them.

| Version       | Supported          |
|---------------|--------------------|
| `0.x` (latest)| :white_check_mark: |
| older `0.x`   | :x:                |

> The project is pre-1.0. Until a `1.0.0` release exists, the most recent commit on
> the default branch is the only fully supported reference point. Pin to a tag or SHA
> in production and watch releases for security advisories.

---

## Out of scope

The following are **not** treated as vulnerabilities in the boilerplate. Reports
limited to these will be closed with an explanation:

- Issues requiring **physical access** to a developer's or server's machine.
- **Social engineering** of maintainers, contributors, or downstream users.
- **Denial of service** achievable only through unrealistic load, or by an attacker
  who already controls the infrastructure.
- Findings that depend on a **misconfigured fork** rather than a flaw in upstream code
  (e.g. committing a real `.env`, disabling the rate limiter, weakening `JWT_SECRET`,
  running with debug/seed credentials in production, exposing Postgres publicly).
- Missing security headers or hardening on **third-party services** you choose to wire
  in (Stripe, Neon, Cloudflare R2, Resend, Gemini) — report those to the respective vendor.
- Vulnerabilities in **dependencies** that have no released patch and no exploitable
  path in CleanSaaS. (We still want to know; we just can't always fix them upstream.)
- Best-practice / informational findings with **no demonstrable security impact**
  (e.g. scanner output without a working proof of concept).
- Self-inflicted issues such as using the seeded `admin@cleansaas.dev` / `admin123`
  credentials, which exist solely for local development.

---

## Built-in security posture

What the boilerplate ships with **today**:

- **Parameterized SQL only.** Pure SQL via `database/sql`, no ORM. Queries use bound
  parameters — never string concatenation — so the codebase is structurally resistant
  to SQL injection.
- **Password hashing with bcrypt.** Credentials are never stored or logged in plaintext.
- **JWT authentication** with signed tokens and expiration. The signing secret is read
  from `JWT_SECRET` (environment only — never hardcoded).
- **Rate limiting** on sensitive endpoints (e.g. authentication) to blunt brute-force
  and credential-stuffing attempts.
- **CORS** configured with an explicit allow-list of origins rather than a wildcard.
- **Security headers** applied via middleware on HTTP responses.
- **Secrets via environment variables only.** No credentials in source; `.env` files
  are gitignored and `.env.example` templates ship instead.
- **Input validation at every boundary** (HTTP and domain layers).
- **Stateless backend** — no server-side session store to compromise.

### Honest roadmap (not yet shipped)

We are explicit about what is **not** in place so you can plan accordingly:

- **Refresh-token rotation** — current JWTs are short-lived access tokens; a rotating
  refresh-token flow is on the roadmap, not implemented yet.
- **Row-Level Security (RLS)** — tenant/row isolation is currently enforced in the
  application layer, not at the Postgres RLS level. Database-enforced RLS is planned.

Until these land, treat token revocation and tenant isolation as application-layer
guarantees and review them against your own threat model before going to production.

---

## Coordinated disclosure

We follow **coordinated disclosure**. We ask that you give us a reasonable window
(generally aligned with the SLA above) to investigate and ship a fix before any public
disclosure. In return, we commit to working the issue promptly, keeping you updated,
crediting you if you wish, and publishing a security advisory once a fix is released.
We will not pursue legal action against researchers who act in good faith, avoid
privacy violations and data destruction, and follow this policy.

Thank you for helping keep CleanSaaS and everyone who builds on it safe.

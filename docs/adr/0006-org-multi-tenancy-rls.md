# 0006. Organization multi-tenancy with PostgreSQL Row-Level Security

## Status

Accepted

## Context

CleanSaaS targets real, multi-user SaaS products. The moment more than one
customer shares a database, the single most expensive bug class becomes
**cross-tenant data leakage**: user A reading or writing user B's rows. A
boilerplate that ships without a credible isolation story teaches the wrong
default and exposes every product built on it.

We need a tenancy model and an enforcement mechanism that is:

- **Correct in depth.** A single missed `WHERE` clause in one query must not be
  able to leak another tenant's data. One layer of defense is not enough — the
  history of SaaS breaches is the history of a forgotten filter.
- **Infrastructure-agnostic at the application layer.** ADR 0001 (hexagonal) and
  ADR 0002 (pure SQL) must hold: the `app` layer must not learn about
  `database/sql`, roles, or GUCs to get isolation.
- **Locally verifiable.** A contributor with one Postgres connection string must
  be able to *prove* isolation works, not just trust it.
- **Compatible with system work.** Migrations, the seeder, and scheduled cleanup
  jobs legitimately operate across all tenants and must not be hobbled by the
  same fence that protects request traffic.

The unit of tenancy is the **organization**. A user belongs to one or more
organizations; on signup they get a personal organization they own. Tenant-owned
resources — **subscriptions, files, conversations, notifications** — belong to an
organization. The blog stays **global/public** (not tenant-scoped); teams are
left as-is for a future iteration.

## Decision

We will enforce organization isolation in **three layers of defense**, with
PostgreSQL Row-Level Security (RLS) as the non-bypassable last line.

**Layer 1 — Handler/app ownership.** Authenticated handlers resolve the caller's
active organization in middleware and the app services check resource ownership
(e.g. a file's `user_id`), as before.

**Layer 2 — Repository `WHERE org_id`.** Every tenant-table query filters by the
active `org_id` (read from the request context, `pkg/orgctx`), and every insert
stamps it. A missing org folds to `org_id = NULL` (never true), so this layer is
deny-by-default too.

**Layer 3 — PostgreSQL RLS (the safety net).** Each tenant table has
`ENABLE` + `FORCE ROW LEVEL SECURITY` and a deny-by-default policy keyed on a
per-transaction setting:

```sql
CREATE POLICY org_isolation ON <table>
  USING      (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid)
  WITH CHECK (org_id = NULLIF(current_setting('app.current_org_id', true), '')::uuid);
```

`current_setting(..., true)` returns NULL when the GUC was never set, and
`NULLIF(..., '')` folds an empty/reset value to NULL as well; `org_id = NULL` is
never true, so **a missing org context sees zero rows and can write nothing — it
fails closed without erroring**. The `messages` table (which has no `org_id`) is
isolated **transitively** via an `EXISTS` against its parent `conversations` row.

**The SET LOCAL ROLE + GUC mechanism.** RLS only bites a role that does not
bypass it. We create a dedicated `app_user` role: `NOLOGIN NOBYPASSRLS`, granted
only the privileges the request path needs, and **not** the owner of the tables.
On the authenticated request path, tenant database work runs inside a short
transaction that first does:

```sql
SET LOCAL ROLE app_user;
SELECT set_config('app.current_org_id', $1, true);  -- $1 = validated UUID, bound param
```

`SET LOCAL` confines both the role switch and the GUC to that transaction, so a
pooled connection can never carry a leftover tenant into the next request. The
org id is **validated as a UUID and passed as a bound parameter** — never
string-concatenated into SQL.

This mechanism is implemented as a clean extension of the existing transaction
seam (ADR: TxManager / DBTX): the adapter exposes per-feature **org-scoped
unit-of-work** ports — `FileScope`, `ConversationScope`, `NotificationScope`,
`SubscriptionScope` — each a `WithOrg…(ctx, fn(repo))` that opens the scoped
transaction and hands the callback a transaction-bound repository. The `app`
services depend only on these small interfaces; they never see `database/sql`,
roles, or GUCs. Unit tests implement a scope by calling the callback directly
with a plain mock repo, so they stay infrastructure-free.

**System paths stay privileged.** Migrations, the seeder, and scheduler cleanups
run as the default connection role, which owns the tables and therefore bypasses
RLS — exactly what these trusted, cross-tenant jobs need. The Stripe webhook is a
system path too: it looks subscriptions up by Stripe ID across all tenants and
resolves `org_id` from the customer's user before persisting.

**Org context end-to-end.** On login the active org is stamped into the access
token as an `org` claim; middleware resolves the active org (verifying membership
of any explicit claim, else the user's default/personal org) into the request
context. On register, the user, their personal organization, and the owner
membership are created in **one transaction**. The public demos run under a fixed
seeded **demo organization**, so they keep working end-to-end as a real tenant.

**Production role setup.** Apply migrations and run the seeder as a privileged
(table-owning, RLS-bypassing) database user. The running API must connect as that
same owner — the request path itself drops to `app_user` via `SET LOCAL ROLE` per
transaction. The `app_user` role is created by migration `015` with exactly the
grants it needs (USAGE on schema; SELECT/INSERT/UPDATE/DELETE on the tenant
tables; SELECT on the `plans` and `users` lookup tables). Hardened deployments may
additionally connect the API as a non-owner role that has been `GRANT app_user`,
which makes RLS apply even to the API's base connection (`FORCE ROW LEVEL
SECURITY` already covers the owner case for the scoped transactions).

## Consequences

**Positive (+)**

- A forgotten `WHERE org_id` cannot leak data: the database itself rejects
  cross-tenant reads and writes. Isolation is a property of the schema, not the
  diligence of every query author.
- The `app` layer stays hexagonal and infra-agnostic — it knows only "run my work
  scoped to the active org." No `database/sql`, role, or GUC leaks upward.
- It is locally verifiable with one connection string. The integration test
  (`go test -tags=integration ./internal/adapter/postgres/ -run RLS`) creates two
  orgs and proves, against the live DB, that org A sees only its rows, cannot read
  org B's, is rejected by `WITH CHECK` when inserting for org B, affects zero rows
  when updating org B, and sees nothing at all with no GUC set.
- System jobs are unaffected — they run privileged by design.

**Negative (−)**

- Every tenant request opens a short transaction to set the role + GUC; reads
  that were single round-trips now run inside a transaction. The cost is small
  and bounded, and external calls (AI streaming, Stripe) are deliberately made
  **outside** any open transaction.
- Two database roles and a per-request `SET LOCAL` are moving parts a contributor
  must understand; this ADR plus code comments are the mitigation.
- Adding a new tenant-scoped table is a checklist, not a one-liner: add `org_id`,
  index it, `ENABLE`+`FORCE` RLS, add the `org_isolation` policy, grant
  `app_user`, and route the repo through a scope.
- `messages` isolation is transitive (via `conversations`); a future tenant table
  with no natural parent needs its own `org_id` rather than a transitive policy.

## Alternatives considered

- **Application-layer filtering only (`WHERE org_id` everywhere).** Simple and
  ORM-free, but one missed clause leaks a tenant, and nothing in the database
  stops it. Rejected as the *sole* defense — it is kept as layer 2.
- **Schema-per-tenant or database-per-tenant.** Strong isolation, but it does not
  fit a boilerplate: migrations fan out across N schemas, connection management
  explodes, and cross-tenant admin/analytics become painful. Wrong default for
  the audience. Rejected.
- **A `tenant_id` set in application code with no DB enforcement.** Same failure
  mode as filtering-only, plus a false sense of safety. Rejected.
- **Wrapping the entire HTTP request in one RLS transaction.** Simplest wiring,
  but it would hold a database transaction open across slow external calls (AI
  streaming, Stripe), tying up a connection for the request's lifetime. Rejected
  in favor of short, per-operation org-scoped transactions.

Supersedes nothing. Builds on ADR 0001 (hexagonal — the scope ports keep `app`
infra-agnostic), ADR 0002 (pure SQL — parameterized throughout, GUC set via a
bound parameter), and ADR 0003 (modularity — blog stays global; each tenant
feature owns its `org_id` with no cross-feature foreign keys).

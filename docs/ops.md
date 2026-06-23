# CleanSaaS — Operations Runbook

Practical, honest operations guide for CleanSaaS (Next.js 16 + Go 1.25/Chi +
PostgreSQL 16). It covers running the stack locally, the env vars that gate
optional modules, database operations, deploying, backup/restore, slow-query
triage, and three incident playbooks.

> **Honesty note.** This is a boilerplate. Much of "production ops" here is
> manual or platform-managed and **not yet automated** — those gaps are called
> out inline as **[not automated]** rather than papered over. Treat this as the
> current reality plus a to-do list, not a description of a mature platform.

---

## 1. Local stack management

### One command: the `/run` skill

The fastest way to bring everything up is the `/run` skill. It runs the full
sequence and smoke-tests it:

```
docker compose up -d   →  migrate up  →  seed  →  backend :8081  →  frontend :3010  →  admin :5174
```

Logs are tailed to `/tmp/cleansaas-backend.log`, `/tmp/cleansaas-frontend.log`,
and `/tmp/cleansaas-admin.log`. Backend logs are structured JSON (`slog`) — read
them first when something misbehaves.

### Manual equivalent

```bash
docker compose up -d                          # Postgres + DbGate
cd backend && make migrate-up && make seed    # schema + admin user
cd backend && make run                        # API on :8081
cd frontend && npm run dev                     # web on :3010
cd admin && npm run dev                        # admin on :5174
```

### Ports

| Service | Port | Notes |
|---|---|---|
| Backend API | **8081** | Routes mounted at **root** (`POST /auth/login`, not `/api/...`) |
| Frontend (Next.js) | **3010** | |
| Admin (Vite + React) | **5174** | |
| PostgreSQL | **5433** | Mapped from container 5432 to avoid clashing with a local Postgres |
| DbGate (web DB UI) | **8082** | Browse local Postgres at `http://localhost:8082` |

### Health check & seeded admin

```bash
curl http://localhost:8081/health      # → {"status":"ok",...} (HTTP 200)
```

Seeded admin (created by `make seed`): **admin@cleansaas.dev** / **admin123**.
Change or remove these before any real deployment.

---

## 2. Environment variables and the modules they gate

All backend config is loaded in `backend/internal/config/config.go`. Optional
features are wired in `cmd/api/main.go` **only if their key is set** — if the key
is empty, the feature is skipped (nil service) and the app still boots. This is
the runtime side of the "every feature is removable" invariant (see
`docs/adr/0003-feature-modularity-removability.md`).

### Always required

| Var | Purpose | Local default |
|---|---|---|
| `DATABASE_URL` | Postgres DSN | `postgres://postgres:postgres@localhost:5433/cleansaas?sslmode=disable` |
| `JWT_SECRET` | Token signing secret | `dev-secret-change-me` (**change in prod**) |
| `PORT` | API listen port | `8081` |
| `FRONTEND_URL` | CORS / link generation | `http://localhost:3006` |

### Optional — each gates a module

| Var(s) | Module enabled | Adapter | If empty |
|---|---|---|---|
| `STRIPE_SECRET_KEY` (+ `STRIPE_WEBHOOK_SECRET`) | **billing** | `adapter/stripe` | billing not wired |
| `R2_ACCOUNT_ID`, `R2_ACCESS_KEY`, `R2_SECRET_KEY`, `R2_BUCKET_NAME`, `R2_PUBLIC_URL` | **storage** | `adapter/r2` | storage not wired (gated on `R2_ACCESS_KEY`) |
| `GEMINI_API_KEY` | **ai** (chat) | `adapter/gemini` | ai not wired |
| `RESEND_API_KEY` | **email** (used by auth + notification) | `adapter/resend` | email sends are skipped |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` / `GOOGLE_REDIRECT_URL` | OAuth login | `adapter/google` | Google OAuth disabled |

> `CLAUDE_API_KEY` and `OPENAI_API_KEY` also exist in config as alternative AI
> providers; the wired default is Gemini.

Modules that need **no key** (DB-only, always on): `notification`, `blog`,
`team`. Core (always on): `auth`, `user`.

Frontend uses `NEXT_PUBLIC_API_URL`; admin uses `VITE_API_URL`. Never commit
`.env` files — they are gitignored.

---

## 3. Database operations

All commands run from `backend/` and use `golang-migrate` under the hood.

| Command | Action |
|---|---|
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back the **last** migration |
| `make migrate-down-all` | Roll back **all** migrations |
| `make migrate-status` | Show current version + dirty flag |
| `make migrate-force VERSION=N` | Force-set version to N (recover from a dirty state) |
| `make seed` | Seed admin user + plans |

To run any of these against **prod**, prefix with the prod DSN:

```bash
DATABASE_URL=<neon_url> make migrate-up
```

### The immutability rule (do not break this)

**Migrations are immutable once applied in production.** Never edit an applied
`.up.sql` / `.down.sql` — write a *new* numbered migration instead. The schema's
history must stay reproducible from the files. Every `up` must have a working
`down`, and migrations use `IF [NOT] EXISTS` to stay idempotent.

Standard workflow: create files → `make migrate-up` locally → verify in DbGate →
test `make migrate-down` → re-apply → commit → `DATABASE_URL=<prod> make migrate-up`.

### The 009 numbering gap is intentional

Migrations currently are: `001`–`008`, then **`010`**, `011`. There is **no
`009`** and that is deliberate — `009` was the **referral** feature's migration,
which was removed during cleanup (the referral and mobile modules no longer
exist). The gap is left in place because applied migration numbers are immutable
(see above) and renumbering would rewrite history. **Do not "fill in" 009.** The
next new migration should be `012`.

Current modules after cleanup: `auth`, `user`, `billing`, `storage`, `ai`,
`notification`, `blog`, `team`. No `referral`, no `mobile`.

### Recovering a dirty migration

```bash
make migrate-status            # shows version + dirty=true
make migrate-force VERSION=N   # N = last known-clean version
# fix the offending SQL, then:
make migrate-up
```

---

## 4. Deploy targets

| Layer | Target |
|---|---|
| Frontend (Next.js) | **Vercel** |
| Backend (Go API) | **Railway** |
| Database | **Neon** (managed Postgres) |
| Object storage | **Cloudflare R2** |

> **[not automated]** There is no CI/CD deploy pipeline shipped. Deploys are
> currently manual via each platform's dashboard/CLI. Wiring GitHub Actions →
> Vercel/Railway is a roadmap item.

### First-deploy checklist

1. **Neon DB**
   - Create the project/branch; copy the connection string.
   - **Drop `channel_binding=require` from the DSN** — `lib/pq` does not support
     it and the connection will fail otherwise. (The backend's `cleanDSN()`
     strips `&channel_binding=require` defensively, but set a clean DSN anyway.)
   - Ensure `sslmode=require` for prod.
2. **Run migrations against prod:** `DATABASE_URL=<neon_url> make migrate-up`.
   Then `DATABASE_URL=<neon_url> make seed` (and immediately change/disable the
   default admin credentials).
3. **Backend on Railway**
   - Set env vars: `DATABASE_URL` (Neon), a strong `JWT_SECRET`, `FRONTEND_URL`
     (the Vercel URL), `PORT`, and whichever optional keys you want enabled
     (`STRIPE_*`, `R2_*`, `GEMINI_API_KEY`, `RESEND_API_KEY`, `GOOGLE_*`).
   - Deploy; confirm `GET /health` returns 200.
4. **Frontend on Vercel**
   - Set `NEXT_PUBLIC_API_URL` to the Railway backend URL.
   - Deploy.
5. **Cloudflare R2** (only if using storage)
   - Create the bucket; generate access keys; set `R2_*` on Railway. Configure
     `R2_PUBLIC_URL` for public asset access.
6. **Stripe** (only if using billing)
   - Set `STRIPE_SECRET_KEY` (live) and create a webhook endpoint pointing at the
     backend; set its signing secret as `STRIPE_WEBHOOK_SECRET`.
7. **CORS / URLs:** confirm `FRONTEND_URL` (backend) and `NEXT_PUBLIC_API_URL`
   (frontend) point at each other. Re-test login end to end.

### The Neon `channel_binding` gotcha (call-out)

Neon's copy-paste connection string includes `channel_binding=require`. The Go
driver (`lib/pq`) does not understand this parameter and the connection fails.
**Remove it.** The backend strips `&channel_binding=require` in `cleanDSN()` as a
safety net, but you should provide a clean DSN regardless.

---

## 5. Backup & restore

> **[not automated]** No scheduled backup job ships with the boilerplate. In
> production, **Neon's point-in-time restore / branching is the primary safety
> net** — use it. The commands below are for manual/local dumps and for moving
> data between environments.

```bash
# Dump (custom format, compressed)
pg_dump "$DATABASE_URL" -Fc -f cleansaas-$(date +%Y%m%d).dump

# Restore into an empty database
pg_restore --clean --if-exists -d "$TARGET_DATABASE_URL" cleansaas-YYYYMMDD.dump

# Plain SQL dump (human-readable / portable)
pg_dump "$DATABASE_URL" --no-owner --no-privileges -f backup.sql
```

Guidance:
- Take a fresh dump (or a Neon branch) **before running any prod migration**.
- Test restores into a throwaway DB periodically — an untested backup is not a
  backup.
- Never restore a prod dump over local without scrubbing PII first.

---

## 6. Slow-query triage

When latency climbs or the DB is hot, work from the data, not from a hunch.

```sql
-- Currently running / blocked queries (run live during an incident)
SELECT pid, state, wait_event_type, now() - query_start AS age, query
FROM pg_stat_activity
WHERE state <> 'idle'
ORDER BY age DESC;

-- Read the plan for a suspect query (does it seq-scan a big table?)
EXPLAIN (ANALYZE, BUFFERS) <your query>;
```

If `pg_stat_statements` is enabled (recommended on Neon), rank by total time:

```sql
SELECT query, calls, total_exec_time, mean_exec_time, rows
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 20;
```

Common causes & fixes in this codebase:
- **Missing index on a foreign key.** Our convention indexes FKs and `010_add_
  performance_indexes` adds more — confirm the column you filter/join on is
  indexed; add a migration (`012_...`) if not.
- **N+1 from a loop in an `app` service.** Because we write raw SQL (see
  `docs/adr/0002-pure-sql-no-orm.md`), N+1s are visible in the adapter — fold the
  loop into a single query with a JOIN or `IN (...)`.
- **Unbounded result set.** Pagination is required everywhere; a query without
  `LIMIT` on a growing table is the bug.

> **[not automated]** No APM/query-monitoring is wired by default. On Neon, lean
> on its dashboard metrics; consider enabling `pg_stat_statements`.

---

## 7. Incident playbooks

Each playbook is **symptom → diagnosis → action**. Backend logs are JSON
(`slog`) at `/tmp/cleansaas-backend.log` locally, or your platform's log stream
in prod.

### (a) Backend 5xx spike

**Symptom**
- Error rate / 5xx climbing; `GET /health` slow or failing; users see
  "something went wrong."

**Diagnosis**
1. `curl https://<backend>/health` — is the process up at all?
2. Tail logs for the error class: panics, `context deadline exceeded` (DB
   timeouts), or `502` from the platform (backend crashed / OOM).
3. Check recent deploys — did a deploy precede the spike? Most 5xx spikes are
   "the last deploy."
4. Distinguish **app crash** (process restarts, 502 from Railway) vs **DB
   pressure** (timeouts, `pq` errors — jump to playbook (b)).

**Action**
- If a recent deploy is the cause: **roll back** to the previous Railway
  deployment first, investigate after.
- If panicking on a code path: capture the stack from logs, hotfix or roll back,
  add a failing test reproducing it (Test→Fix→Retest).
- If DB timeouts dominate: go to playbook (b).
- If OOM/restart loops: bump the instance or fix the leak; confirm `/health`
  green before declaring resolved.

### (b) Database connection exhaustion

**Symptom**
- Errors like `sorry, too many clients already`, `pq: remaining connection slots
  are reserved`, or requests hanging then timing out (`context deadline
  exceeded`).

**Diagnosis**
1. Count connections:
   ```sql
   SELECT count(*), state FROM pg_stat_activity GROUP BY state;
   SHOW max_connections;
   ```
2. Look for **idle in transaction** sessions (connections held open by a leaked
   transaction) — these are usually the culprit.
3. On Neon, check whether you are near the plan's connection limit; consider
   whether the **pooled** (PgBouncer) endpoint should be used.
4. Check for a connection spike correlated with a traffic spike or a deploy that
   changed pool settings.

**Action**
- **Immediate relief:** terminate leaked/idle-in-transaction sessions:
  ```sql
  SELECT pg_terminate_backend(pid)
  FROM pg_stat_activity
  WHERE state = 'idle in transaction'
    AND now() - state_change > interval '5 minutes';
  ```
- **Fix the leak:** ensure every query path closes rows/uses
  `context`-bound timeouts (our convention) and that transactions always
  `Commit`/`Rollback` — an unclosed `*sql.Rows` or un-rolled-back `*sql.Tx` is the
  usual cause.
- **Tune the pool:** set sane `SetMaxOpenConns` / `SetMaxIdleConns` /
  `SetConnMaxLifetime` on the `*sql.DB` so the app cannot exceed the DB's limit;
  on Neon prefer the pooled connection string for the API.

### (c) Stripe webhook failures / duplicates

**Symptom**
- Stripe Dashboard shows failed webhook deliveries (non-2xx), or billing state
  (subscriptions) is out of sync with Stripe; possible duplicate processing
  (e.g. a subscription applied twice).

**Diagnosis**
1. **Signature failures (400):** the endpoint rejects events. Almost always a
   `STRIPE_WEBHOOK_SECRET` mismatch (wrong env, or the raw request body was
   modified before signature verification). Check backend logs for the
   verification error.
2. **5xx on the webhook route:** the handler errored after a valid signature —
   Stripe will retry, which is how **duplicates** arise if the handler is not
   idempotent.
3. **Wrong URL / not reachable:** Dashboard shows timeouts or connection errors —
   verify the endpoint URL points at the live backend and is publicly reachable.
4. Replay a failed event from the Stripe Dashboard (or `stripe trigger` /
   `stripe listen` locally) to reproduce.

**Action**
- Fix the signing secret first if signatures fail; verify against the **raw**
  body, before any JSON decoding.
- **Idempotency is the real fix for duplicates:** Stripe guarantees *at-least-
  once* delivery, so the handler must tolerate replays. Key on Stripe's event ID
  / object ID and make state transitions idempotent (applying the same
  `subscription.updated` twice must be a no-op). **[not automated — verify the
  billing handler dedupes; add a processed-events guard if missing.]**
- Always **return 2xx quickly** once the event is durably accepted; do slow work
  asynchronously so Stripe does not time out and retry.
- After a fix, **replay** the failed events from the Dashboard to resync state,
  then confirm subscription rows match Stripe.

---

## Related docs

- Architecture decisions: `docs/adr/` (hexagonal, no-ORM SQL, removability,
  targeted DDD, enforcement-as-product).
- Conventions: root `CLAUDE.md`, `backend/CLAUDE.md`, `.claude/memory/MEMORY.md`.
- Local skills: `/run`, `/debug`, `/check`, `/test`, `/e2e`.

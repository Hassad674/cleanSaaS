---
name: run
description: Start the whole CleanSaaS stack locally — DB, backend, frontend, admin — from a fresh clone with one flow, then smoke-test it. Use to get the app running or verify a fresh checkout boots.
user-invocable: true
allowed-tools: Read, Bash, Grep, Glob
---

# Run — Boot the Whole Local Stack

Goal: take CleanSaaS from a fresh clone to a fully running, smoke-tested local stack. Run every step in order. All paths are relative to the repo root (`$CLAUDE_PROJECT_DIR`) — `cd` into the repo root first and stay there unless a step says otherwise.

A bootstrap script (`scripts/bootstrap.sh`, `scripts/bootstrap.ps1`) already does steps 1–4 (prereqs → env → DB → migrate → seed). It does **NOT** start the dev servers. This skill can either call bootstrap then do steps 5–7, or run every step itself. Prefer running the steps yourself so you control logging and can report precisely.

**Ports (memorize):** Backend Go API `8081` · Frontend Next.js `3010` · Admin Vite `5174` · Postgres `5433` (→ container `5432`) · DbGate UI `8082`.

---

## STEP 1 — Check prerequisites & report versions

Required: **Docker**, **Go 1.25+**, **Node 20+** (`.nvmrc` pins 20).

```bash
docker --version
go version
node --version
npm --version
```

If any command is missing, STOP and tell the user which tool to install (point them at `README.md`). If Go is older than 1.25 or Node older than 20, warn — the build may fail (see Troubleshooting).

---

## STEP 2 — Copy env files if missing

Three env files are needed. Copy each from its `.example` **only if the target is absent** — never overwrite an existing file.

```bash
[ -f backend/.env ]        || cp backend/.env.example backend/.env
[ -f frontend/.env.local ] || cp frontend/.env.example frontend/.env.local
[ -f admin/.env ]          || cp admin/.env.example admin/.env
```

The defaults are enough to boot core auth. Optional modules stay **disabled** when their key is empty — this is fine for a smoke test. Tell the user explicitly which are off:

| Module  | Key (in `backend/.env`)              | Disabled when empty |
|---------|--------------------------------------|---------------------|
| billing | `STRIPE_SECRET_KEY`                  | yes                 |
| storage | `R2_ACCOUNT_ID` / `R2_*`             | yes                 |
| ai      | `GEMINI_API_KEY`                     | yes                 |

Report which modules will run with reduced functionality.

---

## STEP 3 — Start DB + DbGate, wait for Postgres healthy

```bash
docker compose up -d
```

Then wait until Postgres actually accepts connections (the container is `cleansaas-db`):

```bash
for i in $(seq 1 30); do
  if docker exec cleansaas-db pg_isready -U postgres -d cleansaas >/dev/null 2>&1; then
    echo "postgres ready"; break
  fi
  sleep 1
done
docker exec cleansaas-db pg_isready -U postgres -d cleansaas
```

If it never becomes ready after ~30s, STOP and check `docker compose logs db` (see Troubleshooting: port 5433 busy).

---

## STEP 4 — Migrate + seed

```bash
cd backend && go run cmd/migrate/main.go up    # or: make migrate-up
go run cmd/seed/main.go                          # or: make seed
cd ..
```

The seed creates admin user `admin@cleansaas.dev` / `admin123`, the plans, and blog posts. Seeding is idempotent-ish: re-running **wipes nothing** but may log duplicate-key warnings on already-seeded rows — those warnings are safe to ignore.

---

## STEP 5 — Start the services as background processes

Start each long-running server detached, capturing logs to `/tmp/cleansaas-*.log`. Admin is optional — start it unless the user said otherwise.

```bash
# Backend (:8081). `make run` sources .env then runs the API.
( cd backend && make run ) > /tmp/cleansaas-backend.log 2>&1 &

# Frontend (:3010)
( cd frontend && npm install && npm run dev ) > /tmp/cleansaas-frontend.log 2>&1 &

# Admin (:5174) — optional
( cd admin && npm install && npm run dev ) > /tmp/cleansaas-admin.log 2>&1 &
```

> If launching via the Bash tool, use `run_in_background` for each instead of `&`, and capture the same log paths. `npm install` is only needed on the first run; skip it on reruns if `node_modules` already exists.

Give services time to come up before smoke-testing (backend a few seconds; Next/Vite ~10–30s on first `npm install`). Poll the logs (`tail /tmp/cleansaas-frontend.log`) to confirm "ready"/"compiled" before declaring failure.

---

## STEP 6 — Smoke test each service

Backend health must return HTTP 200 with `{"status":"ok","db":"connected",...}`:

```bash
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8081/health
curl -s http://localhost:8081/health
```

Frontend, admin, and DbGate — confirm each answers (200/30x):

```bash
curl -s -o /dev/null -w "frontend %{http_code}\n" http://localhost:3010
curl -s -o /dev/null -w "admin    %{http_code}\n" http://localhost:5174
curl -s -o /dev/null -w "dbgate   %{http_code}\n" http://localhost:8082
```

If any is not yet up, wait and retry a few times before marking it DOWN — and quote the tail of its log in the report.

---

## STEP 7 — Final report

Output a table of every service, its URL, and live status, plus creds, log locations, and stop instructions:

```
# CleanSaaS — Local Stack

| Service   | URL                          | Status      |
|-----------|------------------------------|-------------|
| Backend   | http://localhost:8081/health | UP (200)    |
| Frontend  | http://localhost:3010        | UP          |
| Admin     | http://localhost:5174        | UP          |
| DbGate    | http://localhost:8082        | UP          |
| Postgres  | localhost:5433 (cleansaas)   | healthy     |

Admin login: admin@cleansaas.dev / admin123
Disabled modules (empty keys): billing, storage, ai   # list only the ones actually empty

Logs:
  backend  → /tmp/cleansaas-backend.log
  frontend → /tmp/cleansaas-frontend.log
  admin    → /tmp/cleansaas-admin.log

Stop everything:
  docker compose down                      # stops DB + DbGate
  # kill dev servers (background processes):
  pkill -f 'cmd/api/main.go'               # backend
  pkill -f 'next dev'                       # frontend
  pkill -f 'vite'                           # admin
```

Mark any service that failed as `DOWN` and include the last ~15 log lines and the likely cause.

---

## Troubleshooting

- **Port 5433 already in use** — another Postgres is bound. Find it: `lsof -i :5433` (or `ss -ltnp | grep 5433`). Stop that service, or remap the host port in `docker-compose.yml`. Same idea for `8081`, `3010`, `5174`, `8082`.
- **`backend/.env` missing / `DATABASE_URL` unset** — `make run` sources `backend/.env`; if it's absent the API can't reach the DB. Re-do Step 2. Confirm `DATABASE_URL` points at `...@localhost:5433/cleansaas`.
- **Backend `/health` returns `db: disconnected` or refuses connection** — Postgres isn't healthy or migrations didn't run. Recheck Step 3 (`pg_isready`) and Step 4.
- **Go version too old** — repo needs Go 1.25+. `go version` to check; upgrade if below.
- **Node version too old** — `.nvmrc` pins 20. Run `nvm use` in `frontend/` and `admin/`, then reinstall deps.
- **`docker compose up` fails** — ensure the Docker daemon is running; inspect `docker compose logs`.
- **Seed prints duplicate-key warnings** — expected on a re-seed; data is not wiped. Safe to ignore.

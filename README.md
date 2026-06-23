# CleanSaaS

Open-source boilerplate for medium-to-large SaaS applications.
**Stack**: Next.js 16 + Go 1.25 (Chi) + PostgreSQL 16 (pure SQL, no ORM) + Tailwind v4.
**Infra**: Vercel (frontend) · Railway (backend) · Neon (database) · Cloudflare R2 (storage).

> Philosophy: every feature is independently removable. The only hard core is auth + database.
> See [`CLAUDE.md`](./CLAUDE.md) for architecture rules.

---

## Vibe coding with CleanSaaS (you don't need to be a developer)

CleanSaaS is calibrated so that, opened in [Claude Code](https://claude.com/claude-code), **describing what you want in plain English produces production-grade code by default** — correctly architected, tested, and checked — without you knowing the internals. Anyone can generate an app with AI; this boilerplate is the guardrails that keep it maintainable, fast, and secure.

Just say what you want. The agent routes your request to the right built-in workflow automatically:

| You type… | What happens |
|---|---|
| *"start the app"* | `/run` boots the database, backend, and frontend, then smoke-tests it. |
| *"add a feature where users can create projects with tasks"* | `/add-feature` scaffolds every layer (database → API → UI), writes tests, and keeps it modular. |
| *"the login button doesn't do anything — here's a screenshot"* | `/debug` walks you through reproducing it (with the Claude Chrome extension), writes a failing test, fixes it, and proves the fix. |
| *"switch payments from Stripe to LemonSqueezy"* | `/add-adapter` swaps the provider by changing one line. |
| *"remove the blog"* | `/remove-feature` deletes it cleanly and proves nothing else broke. |

Every change is automatically held to a professional floor: layered architecture, tests, security (parameterized SQL, auth, rate limits), and CI gates + a git hook that **fail the build if a rule is broken**. You get a senior engineer's standards by default — type a sentence, get a tested feature.

**New here?** Run the [Quickstart](#quickstart) once, then just talk to the agent. Type `/run` to see it live.

---

## Prerequisites

| Tool             | Version         | Notes |
|------------------|-----------------|-------|
| Docker Desktop   | latest          | Provides Postgres + DbGate locally |
| Go               | **1.25+**       | Required by `backend/go.mod` |
| Node.js          | **20+**         | A `.nvmrc` is present at the root |
| Git              | latest          | **On Windows**: install [Git for Windows](https://gitforwindows.org) — provides Git Bash, required to run `make` |

> **Windows users**: open all shell commands below either in **Git Bash** (recommended — `make` works) or in **PowerShell** (use the `.ps1` bootstrap script).

---

## Ports

| Service          | URL                              |
|------------------|----------------------------------|
| Backend (Go API) | http://localhost:8081            |
| Frontend (Next)  | http://localhost:3010            |
| Admin (Vite)     | http://localhost:5174            |
| Postgres         | `localhost:5433`                 |
| DbGate (DB UI)   | http://localhost:8082            |

---

## Quickstart

### Option A — One command (recommended)

**Linux / macOS / Git Bash on Windows**:
```bash
git clone https://github.com/Hassad674/cleanSaaS.git
cd cleanSaaS
./scripts/bootstrap.sh
```

**Windows PowerShell**:
```powershell
git clone https://github.com/Hassad674/cleanSaaS.git
cd cleanSaaS
.\scripts\bootstrap.ps1
```

The script will: check prerequisites → copy `.env.example` files → start Docker → wait for Postgres → run migrations → seed the database.

> If PowerShell blocks the script, run once:
> `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned`

### Option B — Manual

```bash
# 1. Copy env templates
cp backend/.env.example  backend/.env
cp frontend/.env.example frontend/.env.local
cp admin/.env.example    admin/.env

# 2. Start the database
docker compose up -d

# 3. Apply migrations and seed
cd backend
make migrate-up
make seed
cd ..
```

### Then start the services (each in its own terminal)

```bash
# Backend (Git Bash on Windows — make needs bash)
cd backend && make run

# Frontend
cd frontend && npm install && npm run dev

# Admin (optional)
cd admin && npm install && npm run dev
```

---

## Verify it works

```bash
curl http://localhost:8081/health        # → 200 OK
```

Open http://localhost:3010 — you should see the marketing landing page.

**Default admin credentials** (created by `make seed`):
- Email:    `admin@cleansaas.dev`
- Password: `admin123`

Log in at http://localhost:5174 (admin dashboard) or http://localhost:3010/login.

---

## Common commands

### Backend (`cd backend`)
```bash
make run                # Start API on :8081 (loads .env)
make build              # Compile binary to bin/api
make test               # All tests
make test-unit          # Unit tests only (fast)
make migrate-up         # Apply pending migrations
make migrate-down       # Rollback last migration
make migrate-status     # Current migration version
make seed               # Seed admin user + plans + blog posts
make tidy               # go mod tidy
```

### Frontend (`cd frontend`)
```bash
npm run dev             # Dev server on :3010
npm run build           # Production build
npm run lint            # ESLint
npx tsc --noEmit        # Type check
npx vitest run          # Unit tests
npx playwright test     # E2E tests
```

### Admin (`cd admin`)
```bash
npm run dev             # Dev server on :5174
npm run build           # Production build (tsc + vite)
```

### Docker
```bash
docker compose up -d            # Start db + dbgate
docker compose down             # Stop (keeps data)
docker compose down -v          # Stop and WIPE database (irréversible)
docker compose logs -f db       # Follow Postgres logs
```

---

## Project structure

```
cleanSaaS/
├── frontend/         Next.js 16 — feature-based, see frontend/CLAUDE.md
├── backend/          Go + Chi, hexagonal — see backend/CLAUDE.md
│   └── migrations/   Numbered up/down SQL files
├── admin/            Vite + React — see admin/CLAUDE.md
├── docs/             ARCHITECTURE.md, ADRs, ops runbook
├── scripts/          bootstrap.sh / bootstrap.ps1 · ci/ gate scripts · install-git-hooks.sh
├── .github/          CI workflows + Dependabot
├── docker-compose.yml
└── CLAUDE.md         Architecture rules and modularity philosophy
```

Each sub-project has its own `CLAUDE.md` with detailed conventions.

---

## For Claude Code agents

**Order of operations on a fresh clone**:
1. `./scripts/bootstrap.sh` (or `.ps1` on Windows PowerShell)
2. Verify build before doing anything:
   ```bash
   cd backend  && go build ./...
   cd frontend && npx tsc --noEmit
   cd admin    && npx tsc --noEmit
   ```
3. Start services in separate terminals: backend, frontend, (admin if needed).
4. Smoke test: `curl http://localhost:8081/health`.

**If something fails**:
- Postgres won't start → check port 5433 is free (`lsof -i :5433` / `netstat -ano | findstr :5433`).
- `DATABASE_URL is required` → `backend/.env` was not created. Re-run bootstrap.
- Go version mismatch → ensure `go version` ≥ 1.25.
- `make` not found on Windows → you're not in Git Bash. Either switch shells or call `go run cmd/api/main.go` directly (after sourcing `.env`).
- CRLF warnings on Windows → already handled by `.gitattributes`. If you see them, run `git add --renormalize .` once.

**When wiping local DB**: `docker compose down -v` then re-run bootstrap. Migrations are immutable; never edit applied ones, only add new ones.

---

## Troubleshooting

| Symptom | Cause / Fix |
|---|---|
| `port is already allocated` on `:5433` | Another Postgres on 5433. Edit `docker-compose.yml` ports mapping. |
| `permission denied: ./scripts/bootstrap.sh` | `chmod +x scripts/bootstrap.sh` (Linux/Mac), or run via `bash scripts/bootstrap.sh`. |
| `the term 'make' is not recognized` (Windows) | Use Git Bash, not cmd.exe / PowerShell, for Make-based commands. |
| Frontend can't reach backend | Check `NEXT_PUBLIC_API_URL` in `frontend/.env.local` matches backend `PORT`. |
| Stripe / R2 / OAuth errors | These modules are optional. Leave the env vars empty to disable, or fill them with test keys. |

---

## License

MIT (or whatever you set in your repo).

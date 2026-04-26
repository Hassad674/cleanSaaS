#!/usr/bin/env bash
# CleanSaaS bootstrap — Linux / macOS / Git Bash on Windows
# One-command setup: checks prerequisites, copies env files, starts DB, runs migrations and seed.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log()  { printf "${GREEN}[bootstrap]${NC} %s\n" "$*"; }
warn() { printf "${YELLOW}[bootstrap]${NC} %s\n" "$*"; }
err()  { printf "${RED}[bootstrap]${NC} %s\n" "$*" >&2; }

require() {
  if ! command -v "$1" >/dev/null 2>&1; then
    err "Missing required tool: $1"
    err "Install it and retry. See README.md for prerequisites."
    exit 1
  fi
}

log "Checking prerequisites..."
require docker
require go
require node
require npm
log "  docker: $(docker --version | head -1)"
log "  go:     $(go version)"
log "  node:   $(node --version)"

# Copy env files if absent (never overwrite)
copy_env() {
  local src="$1"
  local dst="$2"
  if [ -f "$dst" ]; then
    warn "  $dst already exists, skipping"
  elif [ -f "$src" ]; then
    cp "$src" "$dst"
    log "  $src → $dst"
  else
    err "  $src not found"
    exit 1
  fi
}

log "Copying env templates..."
copy_env backend/.env.example  backend/.env
copy_env frontend/.env.example frontend/.env.local
copy_env admin/.env.example    admin/.env

log "Starting Postgres + DbGate (docker compose)..."
docker compose up -d

log "Waiting for Postgres to be ready..."
for i in $(seq 1 30); do
  if docker exec cleansaas-db pg_isready -U postgres -d cleansaas >/dev/null 2>&1; then
    log "  Postgres ready"
    break
  fi
  if [ "$i" -eq 30 ]; then
    err "Postgres did not become ready within 30s"
    docker compose logs db
    exit 1
  fi
  sleep 1
done

log "Applying migrations..."
(cd backend && go run cmd/migrate/main.go up)

log "Seeding database (admin user, plans, blog posts)..."
(cd backend && go run cmd/seed/main.go)

echo
printf "${GREEN}========================================${NC}\n"
printf "${GREEN}  CleanSaaS is ready to go!${NC}\n"
printf "${GREEN}========================================${NC}\n"
echo
echo "  DbGate (DB UI):     http://localhost:8082"
echo "  Backend (Go):       http://localhost:8081  ->  cd backend && make run"
echo "  Frontend (Next.js): http://localhost:3010  ->  cd frontend && npm run dev"
echo "  Admin (Vite):       http://localhost:5174  ->  cd admin && npm run dev"
echo
echo "  Default admin login: admin@cleansaas.dev / admin123"
echo
echo "  Health check:       curl http://localhost:8081/health"
echo

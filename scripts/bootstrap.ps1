# CleanSaaS bootstrap — Windows PowerShell
# One-command setup: checks prerequisites, copies env files, starts DB, runs migrations and seed.
#
# Usage:  .\scripts\bootstrap.ps1
# If you get an execution policy error, run once:
#   Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned

$ErrorActionPreference = 'Stop'

$RootDir = Resolve-Path (Join-Path $PSScriptRoot '..')
Set-Location $RootDir

function Log     ($msg) { Write-Host "[bootstrap] $msg" -ForegroundColor Green }
function LogWarn ($msg) { Write-Host "[bootstrap] $msg" -ForegroundColor Yellow }
function LogErr  ($msg) { Write-Host "[bootstrap] $msg" -ForegroundColor Red }

function Require-Tool ($name) {
    if (-not (Get-Command $name -ErrorAction SilentlyContinue)) {
        LogErr "Missing required tool: $name"
        LogErr "Install it and retry. See README.md for prerequisites."
        exit 1
    }
}

Log "Checking prerequisites..."
Require-Tool docker
Require-Tool go
Require-Tool node
Require-Tool npm
Log ("  docker: " + (docker --version))
Log ("  go:     " + (go version))
Log ("  node:   " + (node --version))

function Copy-EnvFile ($src, $dst) {
    if (Test-Path $dst) {
        LogWarn "  $dst already exists, skipping"
    }
    elseif (Test-Path $src) {
        Copy-Item $src $dst
        Log "  $src -> $dst"
    }
    else {
        LogErr "  $src not found"
        exit 1
    }
}

Log "Copying env templates..."
Copy-EnvFile 'backend/.env.example'  'backend/.env'
Copy-EnvFile 'frontend/.env.example' 'frontend/.env.local'
Copy-EnvFile 'admin/.env.example'    'admin/.env'

Log "Starting Postgres + DbGate (docker compose)..."
docker compose up -d
if ($LASTEXITCODE -ne 0) { LogErr "docker compose failed"; exit 1 }

Log "Waiting for Postgres to be ready..."
$ready = $false
for ($i = 1; $i -le 30; $i++) {
    docker exec cleansaas-db pg_isready -U postgres -d cleansaas 2>$null | Out-Null
    if ($LASTEXITCODE -eq 0) {
        Log "  Postgres ready"
        $ready = $true
        break
    }
    Start-Sleep -Seconds 1
}
if (-not $ready) {
    LogErr "Postgres did not become ready within 30s"
    docker compose logs db
    exit 1
}

Log "Applying migrations..."
Push-Location backend
go run cmd/migrate/main.go up
$rc = $LASTEXITCODE
Pop-Location
if ($rc -ne 0) { LogErr "migrate failed"; exit 1 }

Log "Seeding database (admin user, plans, blog posts)..."
Push-Location backend
go run cmd/seed/main.go
$rc = $LASTEXITCODE
Pop-Location
if ($rc -ne 0) { LogErr "seed failed"; exit 1 }

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  CleanSaaS is ready to go!"             -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "  DbGate (DB UI):     http://localhost:8082"
Write-Host "  Backend (Go):       http://localhost:8081  ->  cd backend ; make run     (Git Bash)"
Write-Host "  Frontend (Next.js): http://localhost:3010  ->  cd frontend ; npm run dev"
Write-Host "  Admin (Vite):       http://localhost:5174  ->  cd admin ; npm run dev"
Write-Host ""
Write-Host "  Default admin login: admin@cleansaas.dev / admin123"
Write-Host ""
Write-Host "  Health check:       curl http://localhost:8081/health"
Write-Host ""

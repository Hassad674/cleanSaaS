# CleanSaaS — Project Memory

## Project Identity
- Open-source SaaS boilerplate: Next.js + Go (Chi) + PostgreSQL (pure SQL)
- Target: medium-to-large SaaS, NOT micro-SaaS
- Goal: beat OpenSaaS (13K stars), aim 1-3K stars in 3 months
- Philosophy: Claude Code as primary dev tool, natural language driven

## Architecture Decisions
- **Backend**: Hexagonal architecture (domain → port → app → adapter → handler)
- **Frontend**: Feature-based (not hexagonal — doesn't fit React/Next.js)
- **DB**: Pure SQL, no ORM. Neon for prod, Docker PostgreSQL for local
- **Mobile**: Decided against — backend API is decoupled, users add their own client
- **Infra**: Vercel (frontend) + Railway (backend) + Neon (DB) + Cloudflare R2 (storage)
- See [architecture.md](architecture.md) for detailed structure

## Key Principles (non-negotiable)
- TDD: tests first, AI agents self-correct via test results
- SOLID + STUPID avoidance
- Performance, security, scalability, maintainability, evolvability
- No rigidity: adapters are swappable, features are isolated
- AI-agent friendly: small files, explicit interfaces, consistent patterns

## Repo Structure
```
cleanSaaS/
├── frontend/          → Next.js 15 + Tailwind + shadcn
├── backend/           → Go + Chi (hexagonal)
│   ├── cmd/api/       → Entry point
│   ├── internal/      → domain/ port/ app/ adapter/ handler/
│   └── pkg/           → jwt, hash, validate, pagination
├── db/                → Shared SQL schema reference
└── docker-compose.yml → PostgreSQL + DbGate (DB viewer)
```

## Infra Details
- Local DB port: 5433 (not 5432, was already taken)
- Backend default port: 8081 (not 8080, was already taken)
- Neon connection requires removing `channel_binding=require` for lib/pq
- GitHub repo: Hassad674/cleanSaaS

## User Preferences
- Language: French for conversation, English for code/docs
- Ambitious scope — don't simplify or take easy routes
- Wants professional, showcase-quality code
- Sudoers NOPASSWD configured (temporary — remind to remove)

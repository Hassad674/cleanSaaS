---
name: autopilot
description: Run a HEAVY, long, multi-step autonomous task safely across context limits — auto-checkpoint to PROGRESS.md, delegate sub-work to subagents to preserve the main context, and commit each verified increment so the work survives context compaction and is resumable. Use for "build the whole X", "do everything", "while I'm away", big migrations/audits/refactors. NOT for small one-shot tasks.
user-invocable: true
allowed-tools: Read, Bash, Grep, Glob, Edit, Write, Agent
---

# Autopilot — heavy autonomous task runner

Task: **$ARGUMENTS**

You are running a large task that may outlast a single context window. Context compaction WILL drop detail mid-task; this skill makes the work resumable and keeps your main context lean so you can do far more before that happens. (This is the exact method used to upgrade this very repo — see `AUTONOMY-LOG.md`.)

> **When NOT to use this:** small, one-shot tasks. The checkpoint + delegation overhead is only worth it for heavy/multi-phase work (many files, "do everything", a cross-codebase migration/audit, or anything you estimate will exceed ~700K tokens / many tool-steps). For small tasks, just do them directly.

For a non-developer: you don't need to do anything. Describe the big thing you want; the agent will manage itself so it doesn't "forget" mid-way, and will check in / commit progress as it goes.

---

## STEP 1 — Scope & plan

1. Restate the goal in one sentence and list the **standing decisions** (anything the user chose, so they're never re-litigated). If a genuinely blocking product decision is unclear, ask ONE question; otherwise pick the most maintainable default and proceed.
2. Decompose the work into **phases** (each independently shippable & verifiable). Order by dependency and impact.
3. Decide delivery: a branch per phase, or one working branch with atomic commits. Default: work on a feature branch off the current branch; commit per verified step; don't push unless asked.

## STEP 2 — Create the checkpoint file (PROGRESS.md)

Write `PROGRESS.md` at the repo root — the single source of truth, committed. Include:

```
# PROGRESS — <task>

> Resume after compaction: 1) read this file  2) `git branch --show-current` + `git log --oneline -15`
> 3) re-run the validation pipeline to confirm green  4) continue from the first [ ] / [~] below.

## Goal
## Standing decisions (do not re-litigate)
## Current state   (active phase + branch + last action)
## Phases
- [ ] Phase 1 — …   (verification: build/test status)
- [ ] Phase 2 — …
## Blockers log
## Decisions/discoveries during work
```

Use `[x]` done · `[~]` in-progress · `[ ]` todo. **Update it after every step.** (The root `CLAUDE.md` Compact-instructions already point agents at a progress tracker on resume.)

## STEP 3 — Establish a green baseline

Run the validation pipeline and record the result in PROGRESS.md. If it's already red, report that first and decide whether fixing it is in scope — never build on a broken base.

```bash
cd backend && go build ./... && go test ./... -count=1
cd frontend && npx tsc --noEmit && npx vitest run
bash scripts/ci/run-all.sh   # from repo root
```

## STEP 4 — Execute phases by delegating to subagents

For each substantial phase, **delegate to a subagent** (the Agent/Task tool — your "agent team"). This is the key to context safety: the subagent does the heavy file-reading and implementation in ITS OWN context and returns a concise result, so your main window stays lean and you can orchestrate many more phases before compaction.

- Give each subagent a **tight spec**, the relevant file pointers, the repo conventions, and a **hard rule**: *end with `go build ./...` + `go test ./...` green and `gofmt -w` touched files, or REVERT and report — never leave a broken build.*
- Run agents **sequentially** when they touch overlapping files (most code changes do — parallel edits to the same repo conflict). Only parallelize genuinely independent work (e.g. separate docs).
- **Verify every subagent's output yourself** in the main loop before accepting it: run the validation pipeline + gates. Trust, but verify.
- Keep the main loop for orchestration, verification, commits, and high-judgment design — push the bulk of reading/writing into subagents.

## STEP 5 — Commit each verified increment

After a phase is green (build + tests + `bash scripts/ci/run-all.sh` + the pre-commit hook), commit it with a conventional message describing the change. Small atomic commits make the work durable and resumable from any point. Update PROGRESS.md (check the box, note the commit). **Never commit a broken build.**

## STEP 6 — Survive compaction & report

- If context is compacted mid-run, recover deterministically: read `PROGRESS.md` → `git log` → re-run the validation pipeline → resume from the first unchecked phase. Don't re-derive what the checkpoint records.
- When the task (or the session) ends, give the user an honest status: what's done & verified, what remains (with the exact PROGRESS.md/commit pointers to resume), and any blockers. Be honest — never report green you didn't verify.

---

## Guardrails
- Heavy tasks only — don't add this ceremony to small jobs.
- Verified-or-revert on every subagent; green before every commit.
- One source of truth (`PROGRESS.md`), kept current.
- Honesty over optimism in the final report.

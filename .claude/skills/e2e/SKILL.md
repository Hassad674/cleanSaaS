---
name: e2e
description: Run the Playwright end-to-end tests against a running stack, capturing failures with screenshots and traces. Use to run the E2E suite, verify a user flow, or debug a failing spec.
user-invocable: true
allowed-tools: Read, Bash, Grep, Glob
---

# E2E — Playwright End-to-End Suite

Target: **$ARGUMENTS**

You are the E2E runner for CleanSaaS. `$ARGUMENTS` may be empty (run the whole suite), a spec name/filter (e.g. `auth`, `auth.spec.ts`, `settings`), or contain flags like `--headed` or `--debug` to pass through to Playwright.

All commands run from `frontend/`. Use relative paths or `$CLAUDE_PROJECT_DIR` — never hardcode absolute user paths.

Run the steps below in order. Do not skip a step because it "looks fine" — verify each one.

---

## STEP 1 — Ensure the stack is up

The E2E specs (`auth`, `blog`, `settings`) hit the real frontend and backend. Both must be running and reachable before any test will pass.

- **Frontend (Next.js dev)** expected on **:3010**
- **Backend (Go API)** expected on **:8081** (the `helpers.ts` `API_URL` default)
- **Postgres** on **:5433** (via `docker compose up -d`)

Check them:

```bash
# Backend health — must return 200
curl -fsS -o /dev/null -w "%{http_code}" http://localhost:8081/health

# Frontend reachable
curl -fsS -o /dev/null -w "%{http_code}" http://localhost:3010
```

**If either is down:** STOP and tell the user to bring the stack up first (run the `/run` skill, or manually: `docker compose up -d`, then `cd backend && make run`, then `cd frontend && npm run dev`). Do not try to run the suite against a dead stack.

---

## STEP 2 — Ensure the database is seeded

The `login` helper signs in as a seeded user, so auth/settings specs fail with a blank DB. Confirm the seed exists:

```bash
cd $CLAUDE_PROJECT_DIR/backend && make seed
```

The seed creates an admin user (`admin@cleansaas.dev` / `admin123`).

> NOTE: the committed specs in `frontend/e2e/` log in as `admin@cleansaas.com` (`helpers.ts` + `auth.spec.ts` + `settings.spec.ts`). If login specs fail with "invalid credentials", surface this email mismatch (`.com` vs `.dev`) — either re-seed with the email the specs expect or update the specs/helper to match the seeded address. Do not silently change either; report it.

---

## STEP 3 — Ensure Chromium is installed

Playwright is configured for the `chromium` project only. If browsers were never installed, every test errors immediately.

```bash
cd $CLAUDE_PROJECT_DIR/frontend
# Install only if missing (safe to run; downloads chromium + system deps)
npx playwright install --with-deps chromium
```

If Chromium is already present this is a no-op. Run it when in doubt.

---

## STEP 4 — Verify / fix the baseURL & port mismatch (KNOWN BUG)

`frontend/playwright.config.ts` currently pins the suite to **:3006**:

```ts
use: { baseURL: "http://localhost:3006", ... },
webServer: { command: "npm run dev -- --port 3006", port: 3006, reuseExistingServer: !process.env.CI, ... },
```

But the dev server in this repo runs on **:3010**. This mismatch means Playwright either spins up a second server on :3006 (not the one you seeded/expect) or talks to the wrong port. **Surface this explicitly — never ignore it.**

Read the config and check the port:

```bash
grep -nE "baseURL|port|--port" $CLAUDE_PROJECT_DIR/frontend/playwright.config.ts
```

If it still says `3006`, do ONE of the following and tell the user which you chose:

- **Preferred — point this run at :3010 without editing the file:** export the override before running:
  ```bash
  export PLAYWRIGHT_BASE_URL=http://localhost:3010
  ```
  (Only effective if the config reads `process.env.PLAYWRIGHT_BASE_URL`. If it does NOT — and the current config does not — you must edit the config, see below.)
- **Fix the config (recommended permanent fix):** tell the user to change `baseURL` to `http://localhost:3010` and the `webServer` `command`/`port` to `3010` so they match the dev server. With `allowed-tools: Read, Bash, Grep, Glob` this skill cannot edit the file directly — report the exact change needed and ask the user (or use `/run`-style guidance) to apply it.

Do not proceed to run the suite until the port the tests target matches the port the frontend is actually serving (:3010). Reusing the already-running :3010 server is best (avoids a duplicate `npm run dev`).

---

## STEP 5 — Run the suite

From `frontend/`. Scope and flags come from `$ARGUMENTS`:

- **No arguments** → full suite:
  ```bash
  cd $CLAUDE_PROJECT_DIR/frontend && npx playwright test
  ```
- **A spec name/filter** (`auth`, `settings`, `blog`) → single spec:
  ```bash
  cd $CLAUDE_PROJECT_DIR/frontend && npx playwright test e2e/auth.spec.ts
  ```
  (Match the argument to a real file in `frontend/e2e/`: `auth.spec.ts`, `blog.spec.ts`, `settings.spec.ts`. A bare word like `auth` can be passed as a filter: `npx playwright test auth`.)
- **Flags in `$ARGUMENTS`** (`--headed`, `--debug`, `--ui`) → pass through:
  ```bash
  cd $CLAUDE_PROJECT_DIR/frontend && npx playwright test e2e/auth.spec.ts --headed
  ```

Available specs: `e2e/auth.spec.ts`, `e2e/blog.spec.ts`, `e2e/settings.spec.ts` (`e2e/helpers.ts` is a shared helper, not a spec).

---

## STEP 6 — On failure: capture and surface artifacts

Playwright is configured with `screenshot: "only-on-failure"` and `trace: "on-first-retry"`. Artifacts land in:

- `frontend/test-results/` → per-failure screenshots (`*.png`) and traces (`trace.zip`)
- `frontend/playwright-report/` → the HTML report

For each failed test, gather and report:
1. The **spec file + test name** that failed.
2. The **failing assertion / error** (the `expect(...)` line or timeout message from the output).
3. The **artifact path** — locate it:
   ```bash
   ls -R $CLAUDE_PROJECT_DIR/frontend/test-results 2>/dev/null
   ```
   Point to the matching `.png` (screenshot) and `trace.zip` (if a retry produced one).
4. Offer the HTML report for deeper inspection:
   ```bash
   cd $CLAUDE_PROJECT_DIR/frontend && npx playwright show-report
   ```
   To open a specific trace: `npx playwright show-trace test-results/<path>/trace.zip`.

Common root causes to check before blaming the test:
- Stack not up / not seeded (STEP 1–2).
- Port mismatch (STEP 4) — tests timing out on `goto` usually means wrong baseURL.
- Seed email mismatch `.com` vs `.dev` (STEP 2) — login specs fail with invalid credentials.

---

## STEP 7 — Final report

Output a structured report:

```
# E2E Report — {target}

## Stack
- Backend :8081 /health: 200 OK
- Frontend :3010: reachable
- DB seeded: yes (admin@cleansaas.dev / admin123)
- Chromium: installed
- baseURL/port: targeting :3010 (config said :3006 — FIXED / NEEDS FIX)

## Results
- Passed: X
- Failed: Y
- Skipped: Z

## Failures
### [FAIL] e2e/auth.spec.ts › "login with existing user"
- Assertion: expect(page).toHaveURL(/.*dashboard/) timed out after 30000ms
- Likely cause: seed email mismatch (.com vs .dev) / login did not redirect
- Screenshot: frontend/test-results/auth-…/test-failed-1.png
- Trace: frontend/test-results/auth-…/trace.zip (run: npx playwright show-trace <path>)

## Next steps
- <e.g. "Fix playwright.config.ts baseURL to :3010 and re-run", or
       "Re-seed with admin@cleansaas.com to match specs", or
       "Open the report: npx playwright show-report">
```

If all tests pass, report the counts and confirm the stack + config were correct — and still flag the :3006 config if it was only worked around for this run (so it gets a permanent fix).

---
name: debug
description: Guide the user step-by-step to reproduce and fix a bug — collect a screenshot/error, reproduce in the browser or via curl, localize, write a bug report, add a failing test, fix, and verify with the full pipeline.
user-invocable: true
allowed-tools: Read, Bash, Grep, Glob, Edit, Agent
---

# Debug

Bug to investigate: **$ARGUMENTS**

You are the debugging guide for CleanSaaS. Your job is to walk a non-expert all the way from "something is broken" to "it's fixed and proven fixed" — calmly, explicitly, one step at a time. Assume the user may not know what a JWT, a network tab, or a stack trace is, and explain those briefly when they come up. But never slow down a pro: respect the fast-path below.

> **Pro fast-path.** If the user already knows the symptom and where it lives, skip the hand-holding: confirm the baseline (STEP 2), jump to **localize (STEP 4)**, write a **failing test (STEP 6)**, then **fix + verify (STEP 7-8)**. Steps 0-3 are for when the bug isn't yet pinned down.

This skill does NOT reinvent the project's quality rules. It explicitly invokes them:
- The **Test → Fix → Retest loop** (max 3 fix attempts per failing test) — see `CLAUDE.md` › "Test → Fix → Retest loop".
- The **validation pipeline** (`go build ./...`, `npx tsc --noEmit`, `go test ./...`, `npx vitest run`, `npx playwright test`) — see `CLAUDE.md` › "Validation pipeline".
- The **Blocker policy** (same error, 3+ approaches, no progress → write `BLOCKED-<topic>.md`) — see `CLAUDE.md` › "Blocker policy".

**Repo facts you will need** (use exactly):
- Backend: Go + Chi on **:8081**, hexagonal (`handler → app → domain ← port ← adapter`). Routes are in `backend/internal/handler/router.go`, mounted at **ROOT** — e.g. `/auth/login`, NOT `/api/...`.
- Frontend: Next.js on **:3010**, feature-based (`app/` pages → `src/features/<f>/{components,actions,api,hooks}`).
- Admin: Vite on **:5174**.
- Postgres on **:5433**. DbGate (web DB viewer) on **:8082**.
- Logs (if the stack was started via the `/run` skill): `/tmp/cleansaas-backend.log`, `/tmp/cleansaas-frontend.log`, `/tmp/cleansaas-admin.log`.
- Seeded admin account: **admin@cleansaas.dev** / **admin123**.
- Always use relative paths or `$CLAUDE_PROJECT_DIR`. Never absolute user paths.

---

## STEP 0 — Intake (talk to the user)

Before touching any code, understand the bug from the human's point of view. Ask for these four things in plain language:

1. **What did you expect to happen?**
2. **What actually happened instead?**
3. **Evidence** — paste the *exact* error text, OR drop a **screenshot**. If they reference an image file (e.g. `./bug.png`, or a path they paste), **Read it** — you can see screenshots. Look at error toasts, red text, broken layout, the URL bar, and any visible status codes.
4. **Where** — which page/URL (e.g. `http://localhost:3010/dashboard`) or which API call (e.g. "Sign in" button)?

If the user can't articulate it, give them this short checklist to fill in — it's friendlier than a blank page:

```
[ ] Page or screen:        (URL or what you were looking at)
[ ] What I clicked/typed:   (the exact action)
[ ] What I expected:
[ ] What I saw instead:     (error text, or "nothing happened", or "wrong number")
[ ] Screenshot attached?    (drag the file in, or paste its path)
[ ] First time, or did it used to work?
```

Don't move on until you have at least: a symptom, and a place (a URL or an API action). If `$ARGUMENTS` already contains a clear symptom, use it and only ask for what's missing.

---

## STEP 1 — Classify the bug surface

From the intake, sort the bug into ONE primary surface. This decides which reproduction playbook you follow in STEP 3. Tell the user which one you picked and why (one sentence).

| Surface | Tell-tale signs | Repro path |
|---|---|---|
| **frontend-render** | Blank page, wrong layout, missing data on screen, React error overlay, console errors, hydration warnings | STEP 3A (browser) |
| **API/network** | A request returns 4xx/5xx, wrong JSON, slow/hanging request, CORS error | STEP 3A (browser) to capture, then STEP 3B (curl) to isolate |
| **auth** | "401 Unauthorized", login fails, redirected to login, "token expired", can't access a protected page | STEP 3B (curl with login) |
| **database/migration** | 500 error mentioning a column/table, "relation does not exist", data not saved, migration failed | STEP 3C (db) |
| **build/compile** | App won't start, red TypeScript errors, `go build` fails, white screen + build error in terminal | STEP 2 surfaces this directly |

If it spans two surfaces (common: a UI error that's really an API 500), capture in the browser first (3A), then confirm the root cause server-side (3B/3C).

---

## STEP 2 — Clean baseline (never debug on broken ground)

Before reproducing anything, confirm the project itself is healthy. If the baseline is already red, **that** is the first thing to report — you can't trust a bug investigation on top of a broken build or a failing suite.

Run the validation pipeline (this is `CLAUDE.md`'s pipeline — same commands):

```bash
cd "$CLAUDE_PROJECT_DIR/backend" && go build ./... 2>&1 | tail -30
cd "$CLAUDE_PROJECT_DIR/frontend" && npx tsc --noEmit 2>&1 | tail -30
cd "$CLAUDE_PROJECT_DIR/backend" && go test ./... -count=1 2>&1 | tail -40
cd "$CLAUDE_PROJECT_DIR/frontend" && npx vitest run 2>&1 | tail -40
```

Interpret the result:
- **All green** → good, the bug is a real regression in otherwise-healthy code. Continue to STEP 3.
- **Baseline already red** → STOP the bug hunt. Report to the user, in plain words: *"Before we look at your bug, the project itself isn't building/passing right now — here's what's broken: …"*. Decide with them whether the red baseline IS their bug (it often is for **build/compile** bugs) or a separate pre-existing problem. If it's a `build/compile` bug, you've already found it — skip to STEP 4 to localize the failing file.

> Beginner note: "the build" just means turning the source code into a runnable app. "The test suite" is a set of automatic checks the project runs to catch mistakes. Green = passing, red = something failed.

---

## STEP 3 — Reproduce the bug

You must see the bug happen with your own (or the user's) eyes before fixing it. Pick the playbook from STEP 1.

### STEP 3A — UI bug (frontend-render / capturing an API error)

**First, make sure the app is actually running.** If the stack isn't up, suggest the `/run` skill (it launches backend :8081, frontend :3010, admin :5174 and tees logs to `/tmp/cleansaas-*.log`). Quick check:

```bash
curl -s -o /dev/null -w "frontend %{http_code}\n" http://localhost:3010 || echo "frontend not up — run the /run skill"
curl -s -o /dev/null -w "backend  %{http_code}\n" http://localhost:8081/health 2>/dev/null || echo "backend not up — run the /run skill"
```

**Then reproduce in the browser using the Claude Chrome extension** (the easiest path for a beginner — Claude can look at the page directly). Walk the user through it, literally:

1. Open the page where the bug happens, e.g. `http://localhost:3010/...`.
2. Open the **Claude Chrome extension** (the Claude icon in your browser's toolbar, usually top-right). Make sure it's connected to this tab.
3. Now **redo the exact action that breaks** — click the button, submit the form, whatever it was.
4. Ask the extension, in its chat box, to capture:
   - **The failing network request**: its **URL**, **method** (GET/POST/…), **status code** (e.g. 401, 500), and the **response body** (the text the server sent back).
   - **The console errors**: any red error messages and the **stack trace** (the list of code locations showing where it broke).
5. **Copy all of that back here** — paste the request details and the console errors into our chat.

> Beginner glossary:
> - **Network request** = a message your browser sends to the backend (e.g. "log me in"). The **status code** is a 3-digit reply: 2xx = success, 4xx = "you did something wrong" (e.g. 401 = not authorized), 5xx = "the server broke".
> - **Console** = a hidden log inside the browser where the page prints errors.
> - **Stack trace** = a breadcrumb trail of which functions were running when the error happened — gold for finding the bug.

**Fallback if the user doesn't have the Claude Chrome extension** — use the browser's built-in DevTools:

1. On the broken page, press **F12** (or right-click → "Inspect") to open DevTools.
2. Click the **Network** tab. Tick "Preserve log". Now redo the broken action.
3. Find the **red** row (the failed request). Click it. Copy: the **Request URL**, **Request Method**, **Status Code** (top of the "Headers" panel), and the **Response** tab's body.
4. Click the **Console** tab. Copy any **red** error text, including the lines underneath it (the stack trace).
5. Paste all of that here.

Once you have the request + status + body + console error, you usually also know whether the true root cause is frontend or backend. If the status is 4xx/5xx, treat it as an **API/network** or **auth** bug and confirm server-side with STEP 3B.

### STEP 3B — API / auth bug (reproduce it yourself with curl)

You don't need the user for this — reproduce the failing call directly against the backend on **:8081**. Remember routes are mounted at ROOT (`/auth/login`, not `/api/auth/login`).

**If the endpoint is public** (e.g. login itself), just call it and print status + body:

```bash
curl -s -i -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@cleansaas.dev","password":"admin123"}'
```

**If the endpoint is protected** (most are), you first need a **JWT** — a signed token the backend gives you at login that proves who you are. You attach it to later requests in an `Authorization: Bearer <token>` header. Log in, capture the token, then call the protected endpoint:

```bash
# 1) Log in and extract the token (jq if available, else a grep fallback)
TOKEN=$(curl -s -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@cleansaas.dev","password":"admin123"}' \
  | (jq -r '.token // .access_token // .data.token' 2>/dev/null || grep -o '"token":"[^"]*"' | cut -d'"' -f4))
echo "token: ${TOKEN:0:20}..."   # show only the start, tokens are secret

# 2) Call the protected endpoint that's misbehaving (replace path/method/body)
curl -s -i -X GET http://localhost:8081/<the-failing-endpoint> \
  -H "Authorization: Bearer $TOKEN"
```

Read the printed **status line** and **body**. A reproducible non-2xx (or a 2xx with wrong data) is your confirmed bug. If the exact token JSON shape differs, peek at the login handler/response DTO (`backend/internal/handler/auth.go`, `backend/internal/handler/dto/response/`) to learn the field name.

> Beginner note: `curl` is a command-line tool that sends the same kind of request your browser does, but it shows you the raw reply. `-i` prints the status and headers; `-X` sets the method; `-d` sends the body.

### STEP 3C — Database / migration bug

Symptoms like "relation does not exist", "column … does not exist", or saved data vanishing usually mean the database schema and the code disagree. Check migration state first:

```bash
cd "$CLAUDE_PROJECT_DIR/backend"
go run cmd/migrate/main.go status 2>/dev/null || make migrate-status 2>/dev/null || echo "check backend/Makefile for the migrate target"
# If migrations are behind, apply them:
# go run cmd/migrate/main.go up   (or: make migrate-up)
```

Then inspect the actual data/schema in **DbGate** at `http://localhost:8082` (a web UI for the Postgres on :5433) — tell the user to open it, expand the table in question, and confirm whether the column/row they expect is really there. Cross-check the table definition against the relevant file in `backend/migrations/` (every `*.up.sql` must have a matching `*.down.sql`). A mismatch between a `SELECT`/`INSERT` in `backend/internal/adapter/postgres/<feature>.go` and the migrated schema is a classic root cause.

---

## STEP 4 — Localize (map the symptom to a line of code)

Now turn the evidence into a `file:line`. Trace the architecture in the direction of the data.

**Backend (handler → app → domain):**
1. Find the route in `backend/internal/handler/router.go` to get the handler function name.
   ```bash
   grep -n "<path-fragment>" "$CLAUDE_PROJECT_DIR/backend/internal/handler/router.go"
   ```
2. Open that handler in `backend/internal/handler/<feature>.go` — it decodes the request and calls the app service.
3. Follow into `backend/internal/app/<feature>/service.go` (the use-case / business logic) and, if needed, `backend/internal/domain/<feature>/` (validation rules) and `backend/internal/adapter/postgres/<feature>.go` (the SQL).
4. The stack trace from STEP 3 points you straight at the failing frame — start there.

**Frontend (page → feature → action → api):**
1. The URL maps to a file under `frontend/src/app/...`.
2. That page composes components from `frontend/src/features/<f>/components/`.
3. Data flows through `frontend/src/features/<f>/actions/` (server actions) and `frontend/src/features/<f>/api/` (the fetch layer).
   ```bash
   grep -rn "<endpoint-or-symbol>" "$CLAUDE_PROJECT_DIR/frontend/src/features/"
   ```

**Read the suspect files** (don't guess) and pin the exact `file:line` where the behavior diverges from what's expected. If you're unsure across many files, dispatch an `Agent` (Explore) to fan out the search and report the candidate locations — but you keep the conclusion.

---

## STEP 5 — Write the bug report (required artifact)

Produce this structured report and **show it to the user** before fixing. It forces clear thinking and teaches the beginner how a real bug is documented. Keep it tight.

```
# Bug report — <short title>

## Summary
<one or two sentences: what's broken, who it affects>

## Repro steps
1. <exact step>
2. <exact step>
3. <what breaks>

## Expected vs Actual
- Expected: <…>
- Actual:   <… include the status code / error text>

## Evidence
- Screenshot: <path/ref the user gave, if any>
- Request/response: <method + URL + status + key part of body>
- Log excerpt: <relevant lines from /tmp/cleansaas-*.log if available>

## Suspected location
- <relative/path/file.go>:<line>  (and the call chain that reaches it)

## Hypothesis (root cause)
<your best explanation of WHY it happens — the actual defect, not just the symptom>
```

If logs were captured by `/run`, grab the relevant tail:

```bash
tail -40 /tmp/cleansaas-backend.log 2>/dev/null
tail -40 /tmp/cleansaas-frontend.log 2>/dev/null
```

---

## STEP 6 — Write a FAILING test (red first)

Per the project's TDD rule (`CLAUDE.md` › Test-Driven Development), capture the bug in an automatic test that **fails right now**. This proves you understood the bug and gives you a green light to aim for. Pick the smallest test that reproduces it:

- **Backend unit** — add a case to the relevant `_test.go` next to the suspect code (`backend/internal/app/<feature>/service_test.go`, `domain/<feature>/entity_test.go`, etc.). Name it `TestServiceName_MethodName_Scenario`. Run only it:
  ```bash
  cd "$CLAUDE_PROJECT_DIR/backend" && go test ./internal/app/<feature>/... -run '<TestName>' -count=1 -v
  ```
- **Frontend unit** — add/extend a `vitest` test near the component or hook. Run only it:
  ```bash
  cd "$CLAUDE_PROJECT_DIR/frontend" && npx vitest run -t "<test name>"
  ```
- **End-to-end** (best when the bug is a user-flow across UI + API) — add a Playwright spec in `frontend/e2e/`. Run only it:
  ```bash
  cd "$CLAUDE_PROJECT_DIR/frontend" && npx playwright test -g "<spec name>"
  ```

**Confirm it FAILS for the right reason** (the bug's reason, not a typo in the test). A test that fails because you mis-imported something proves nothing. Show the user the red output and say, in one line, what it's asserting.

---

## STEP 7 — Fix, then turn the test green

Make the **smallest change** at the `file:line` from STEP 4 that makes your failing test pass. Don't refactor unrelated code while you're here.

Re-run the one test until it's green, following `CLAUDE.md`'s **Test → Fix → Retest loop** — **max 3 fix attempts per failing test**:

```
1. Apply the fix
2. Re-run just this test
   ├── PASS → go to STEP 8 (full pipeline)
   └── FAIL → read the error, adjust the FIX (or the test, if the test itself was wrong), retry
       (stop after 3 attempts → Blocker policy below)
```

**Blocker policy** (`CLAUDE.md` › Blocker policy): if the *same error* survives 3+ genuinely different approaches with no progress, you're blocked. Write `$CLAUDE_PROJECT_DIR/BLOCKED-<topic>.md` containing the error, every approach you tried, and the suspected location, leave a `// TODO: fix` at the spot, and **tell the user clearly**: *"I'm blocked on X — here's what I tried and why none worked. I've logged it in BLOCKED-<topic>.md."* Do not fake a green test, and never delete or skip a test to make the suite pass. (A long task is not a blocker — only "stuck on the same error" is.)

---

## STEP 8 — Verify with the full validation pipeline

A green unit test isn't enough — make sure the fix didn't break anything else. Run the **full** pipeline (`CLAUDE.md` › Validation pipeline):

```bash
cd "$CLAUDE_PROJECT_DIR/backend"  && go build ./...            # compiles
cd "$CLAUDE_PROJECT_DIR/frontend" && npx tsc --noEmit          # typechecks
cd "$CLAUDE_PROJECT_DIR/backend"  && go test ./... -count=1     # all backend tests
cd "$CLAUDE_PROJECT_DIR/frontend" && npx vitest run            # all frontend unit tests
# If the bug was a user flow and Playwright is set up:
cd "$CLAUDE_PROJECT_DIR/frontend" && npx playwright test
```

Everything must be green. If the wider run surfaces a *new* failure your fix caused, that's now part of this bug — loop back to STEP 7 for it. Compilation failures are top priority (`CLAUDE.md` Blocker policy Type C): fix within ~10 min or revert the change so the build is never left broken.

---

## STEP 9 — Report and teach

Close the loop with a short, friendly summary the user can actually learn from:

```
## Fixed ✅
- What changed: <relative/path/file>:<line> — <one line on the change>
- Now green:    <the test you added> + full pipeline passing
- Why it happened (plain language): <one sentence a non-expert understands —
  the actual cause, e.g. "the login form sent the password under the wrong
  field name, so the server never received it">
```

Then remind the user to commit, using a conventional message (`CLAUDE.md` › Git):

```bash
# from the repo root, e.g.:
git add -A && git commit -m "fix(<feature>): <short description of the bug fixed>"
```

If you hit a blocker in STEP 7, say so here instead: point to `BLOCKED-<topic>.md`, summarize what's left, and suggest the next thing to try. Honesty over a fake green.

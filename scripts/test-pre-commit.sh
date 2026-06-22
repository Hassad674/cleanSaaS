#!/usr/bin/env bash
#
# test-pre-commit.sh — meta-test proving the pre-commit hook FAILS CLOSED.
#
# Strategy: build a disposable, isolated git repo in a temp dir, drop the real
# hook into it, and drive it through git's own staging machinery. This way the
# test NEVER touches the developer's real index, working tree, commits, or git
# config — and it cleans up everything on exit (even on failure).
#
# Assertions:
#   1. A staged, deliberately MISFORMATTED .go file  -> hook exits NON-ZERO and
#      its output mentions gofmt.
#   2. A staged, correctly-formatted .go file        -> hook exits ZERO.
#   3. (bonus) A staged .env file                    -> hook exits NON-ZERO.
#
set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
HOOK_SRC="$REPO_ROOT/.githooks/pre-commit"

if [ -t 1 ]; then
  RED='\033[0;31m'; GREEN='\033[0;32m'; BOLD='\033[1m'; NC='\033[0m'
else
  RED=''; GREEN=''; BOLD=''; NC=''
fi

ok()   { printf "${GREEN}  ok${NC}   - %s\n" "$*"; }
bad()  { printf "${RED}  FAIL${NC} - %s\n" "$*"; }

if [ ! -x "$HOOK_SRC" ]; then
  printf "${RED}meta-test: hook not found or not executable: %s${NC}\n" "$HOOK_SRC" >&2
  exit 1
fi

# --- isolated sandbox --------------------------------------------------------
SANDBOX="$(mktemp -d "${TMPDIR:-/tmp}/cleansaas-precommit-test.XXXXXX")"
cleanup() { rm -rf "$SANDBOX"; }
trap cleanup EXIT INT TERM

PASS_COUNT=0
FAIL_COUNT=0
assert() { # assert <description> <expected-exit:pass|fail> <actual-exit> <output> <must-contain-or-empty>
  desc="$1"; expect="$2"; code="$3"; out="$4"; needle="${5:-}"
  exp_nonzero=0; [ "$expect" = "fail" ] && exp_nonzero=1
  got_nonzero=0; [ "$code" -ne 0 ] && got_nonzero=1
  if [ "$exp_nonzero" != "$got_nonzero" ]; then
    bad "$desc (expected $expect, hook exited $code)"
    printf '%s\n' "$out" | sed 's/^/        | /'
    FAIL_COUNT=$((FAIL_COUNT+1)); return
  fi
  if [ -n "$needle" ] && ! printf '%s' "$out" | grep -qi -- "$needle"; then
    bad "$desc (exit ok but output did not mention '$needle')"
    printf '%s\n' "$out" | sed 's/^/        | /'
    FAIL_COUNT=$((FAIL_COUNT+1)); return
  fi
  ok "$desc"
  PASS_COUNT=$((PASS_COUNT+1))
}

# Run the hook inside the sandbox repo against whatever is currently staged.
# Colors are disabled because output is piped (not a tty), keeping greps simple.
run_hook() { ( cd "$SANDBOX" && "$SANDBOX/.githooks/pre-commit" ) 2>&1; }

printf "${BOLD}== CleanSaaS pre-commit meta-test ==${NC}\n"
printf "sandbox: %s\n\n" "$SANDBOX"

# --- set up the disposable repo ---------------------------------------------
git -C "$SANDBOX" init -q
git -C "$SANDBOX" config user.email "test@example.com"
git -C "$SANDBOX" config user.name  "Meta Test"
mkdir -p "$SANDBOX/.githooks"
cp "$HOOK_SRC" "$SANDBOX/.githooks/pre-commit"
chmod +x "$SANDBOX/.githooks/pre-commit"

# ============================================================================
# CASE 1 — misformatted Go file must be REJECTED, citing gofmt.
# ============================================================================
# Bad formatting: tabs replaced by spaces + extra spacing that gofmt rewrites.
mkdir -p "$SANDBOX/backend/internal/sample"
cat > "$SANDBOX/backend/internal/sample/bad.go" <<'EOF'
package sample

func   Add( a int ,b int )  int {
        return a+b
}
EOF
git -C "$SANDBOX" add backend/internal/sample/bad.go
out1="$(run_hook)"; code1=$?
assert "misformatted .go file is rejected" fail "$code1" "$out1"
assert "rejection output cites gofmt"      fail "$code1" "$out1" "gofmt"
git -C "$SANDBOX" reset -q   # unstage

# ============================================================================
# CASE 2 — gofmt-clean Go file must PASS.
# ============================================================================
rm -f "$SANDBOX/backend/internal/sample/bad.go"
# Produce a guaranteed-clean file by running gofmt itself (if available).
cat > "$SANDBOX/backend/internal/sample/good_src.go" <<'EOF'
package sample

func Add(a int, b int) int {
	return a + b
}
EOF
if command -v gofmt >/dev/null 2>&1; then
  gofmt "$SANDBOX/backend/internal/sample/good_src.go" > "$SANDBOX/backend/internal/sample/good.go"
else
  cp "$SANDBOX/backend/internal/sample/good_src.go" "$SANDBOX/backend/internal/sample/good.go"
fi
rm -f "$SANDBOX/backend/internal/sample/good_src.go"
git -C "$SANDBOX" add backend/internal/sample/good.go
out2="$(run_hook)"; code2=$?
assert "gofmt-clean .go file passes" pass "$code2" "$out2"
git -C "$SANDBOX" reset -q

# ============================================================================
# CASE 3 — staged .env file must be REJECTED (secret guard).
# ============================================================================
printf 'SECRET_KEY=hunter2\n' > "$SANDBOX/.env"
git -C "$SANDBOX" add -f .env   # -f because a generated .gitignore might exclude it
out3="$(run_hook)"; code3=$?
assert "staged .env file is rejected" fail "$code3" "$out3" "env"
git -C "$SANDBOX" reset -q

# ============================================================================
# CASE 4 — graceful degradation: hook must NOT crash / must SKIP cleanly when
# a clean Go file is staged but tooling output should show SKIP, not FAIL,
# for absent toolchains. We can't uninstall go here, so we assert the hook
# never blocks purely due to a *missing* toolchain by staging a non-Go,
# non-TS, non-env file and confirming it passes.
# ============================================================================
printf '# readme\n' > "$SANDBOX/NOTES.md"
git -C "$SANDBOX" add NOTES.md
out4="$(run_hook)"; code4=$?
assert "unrelated file (no toolchain needed) passes" pass "$code4" "$out4"
git -C "$SANDBOX" reset -q

# --- final tally; confirm sandbox left no staged changes --------------------
STAGED_LEFT="$(git -C "$SANDBOX" diff --cached --name-only)"
if [ -n "$STAGED_LEFT" ]; then
  bad "sandbox left staged changes (test cleanup bug): $STAGED_LEFT"
  FAIL_COUNT=$((FAIL_COUNT+1))
else
  ok "no staged changes left behind in sandbox"
  PASS_COUNT=$((PASS_COUNT+1))
fi
COMMITS="$(git -C "$SANDBOX" rev-list --all --count 2>/dev/null || echo 0)"
if [ "$COMMITS" -ne 0 ]; then
  bad "sandbox created commits ($COMMITS); meta-test must not commit"
  FAIL_COUNT=$((FAIL_COUNT+1))
else
  ok "no commits created"
  PASS_COUNT=$((PASS_COUNT+1))
fi

printf "\n${BOLD}Result:${NC} %d passed, %d failed.\n" "$PASS_COUNT" "$FAIL_COUNT"
if [ "$FAIL_COUNT" -ne 0 ]; then
  printf "${RED}META-TEST FAILED${NC}\n"
  exit 1
fi
printf "${GREEN}META-TEST PASSED — hook fails closed.${NC}\n"
exit 0

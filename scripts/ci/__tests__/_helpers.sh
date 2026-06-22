#!/usr/bin/env bash
#
# _helpers.sh — shared helpers for the invariant-gate meta-tests.
#
# Sourced by each test-*.sh. Provides assertion helpers, a temp-dir factory,
# and a way to copy a check script into a fixture repo so it can be run with the
# fixture as ROOT (the checks resolve ROOT as <script>/../.. — so we place the
# script at <fixture>/scripts/ci/<check> and run it from there).
set -euo pipefail

CI_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"   # .../scripts/ci

TESTS_RUN=0
TESTS_FAILED=0

pass() { printf '  ok   %s\n' "$1"; TESTS_RUN=$((TESTS_RUN + 1)); }
fail() { printf '  FAIL %s\n' "$1"; TESTS_RUN=$((TESTS_RUN + 1)); TESTS_FAILED=$((TESTS_FAILED + 1)); }

# mk_fixture <check-script-name>  -> echoes path to a fresh fixture root that
# already contains scripts/ci/<check-script-name> (executable).
mk_fixture() {
  check="$1"
  dir="$(mktemp -d)"
  mkdir -p "$dir/scripts/ci"
  cp "$CI_DIR/$check" "$dir/scripts/ci/$check"
  chmod +x "$dir/scripts/ci/$check"
  printf '%s' "$dir"
}

# run_check <fixture-root> <check-name> -> sets RC and OUT
run_check() {
  fixture="$1"; check="$2"
  set +e
  OUT="$(bash "$fixture/scripts/ci/$check" 2>&1)"
  RC=$?
  set -e
}

# assert_exit <expected> <label>
assert_exit() {
  if [ "$RC" -eq "$1" ]; then pass "$2 (exit $RC)"; else fail "$2 (expected exit $1, got $RC)"; fi
}

# assert_contains <needle> <label>
assert_contains() {
  if printf '%s' "$OUT" | grep -qF -- "$1"; then pass "$2"; else
    fail "$2 — output did not contain: $1"
    printf '%s\n' "$OUT" | sed 's/^/      | /'
  fi
}

# assert_annotation <path-substr> <label> : checks a ::error annotation names the file
assert_annotation() {
  if printf '%s' "$OUT" | grep -qE "::error file=[^,]*$1"; then pass "$2"; else
    fail "$2 — no ::error annotation naming: $1"
    printf '%s\n' "$OUT" | sed 's/^/      | /'
  fi
}

finish() {
  echo "  ---- $TESTS_RUN assertion(s), $TESTS_FAILED failed ----"
  [ "$TESTS_FAILED" -eq 0 ]
}

#!/usr/bin/env bash
#
# run-all.sh — run every CleanSaaS invariant gate, aggregate results,
# print a summary, and exit non-zero if any gate failed.
#
# Each gate is run in its own subshell so one failure does not abort the rest
# (we must report ALL failures, not just the first).
set -uo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CI_DIR="$ROOT/scripts/ci"

CHECKS="
check-cross-feature-imports.sh
check-migration-pairs.sh
check-hardcoded-colors.sh
check-file-length.sh
check-forbidden-names.sh
"

overall=0
passed=0
failed=0
results=""

for check in $CHECKS; do
  script="$CI_DIR/$check"
  echo "========================================"
  echo ">>> $check"
  echo "========================================"
  if [ ! -x "$script" ]; then
    if [ -f "$script" ]; then
      bash "$script"
    else
      echo "::error file=scripts/ci/$check,line=1::gate script missing"
      echo "MISSING: $check"
      rc=1
    fi
    rc=${rc:-$?}
  else
    "$script"
    rc=$?
  fi

  if [ "$rc" -eq 0 ]; then
    passed=$((passed + 1))
    results="$results
  PASS  $check"
  else
    failed=$((failed + 1))
    overall=1
    results="$results
  FAIL  $check (exit $rc)"
  fi
  echo
done

echo "========================================"
echo "INVARIANT GATE SUMMARY"
echo "========================================"
printf '%s\n' "$results"
echo "----------------------------------------"
echo "Passed: $passed   Failed: $failed"
if [ "$overall" -ne 0 ]; then
  echo "RESULT: FAIL — one or more invariant gates failed"
  exit 1
fi
echo "RESULT: PASS — all invariant gates green"
exit 0

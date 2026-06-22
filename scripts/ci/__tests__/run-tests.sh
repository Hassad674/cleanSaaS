#!/usr/bin/env bash
#
# run-tests.sh — run every invariant-gate meta-test.
#
# Each test-*.sh copies a minimal fixture to a temp dir, injects a known
# violation, runs the corresponding gate, and asserts it fails (exit !=0) and
# names the offending file, plus that the gate passes on a clean fixture.
#
# Exits non-zero if any meta-test fails.
set -uo pipefail

TESTS_DIR="$(cd "$(dirname "$0")" && pwd)"

overall=0
passed=0
failed=0
summary=""

for t in \
  test-cross-feature-imports.sh \
  test-migration-pairs.sh \
  test-hardcoded-colors.sh \
  test-file-length.sh \
  test-forbidden-names.sh
do
  echo "########################################"
  bash "$TESTS_DIR/$t"
  rc=$?
  if [ "$rc" -eq 0 ]; then
    passed=$((passed + 1)); summary="$summary
  PASS  $t"
  else
    failed=$((failed + 1)); overall=1; summary="$summary
  FAIL  $t (exit $rc)"
  fi
  echo
done

echo "########################################"
echo "META-TEST SUMMARY"
echo "########################################"
printf '%s\n' "$summary"
echo "----------------------------------------"
echo "Passed: $passed   Failed: $failed"
if [ "$overall" -ne 0 ]; then
  echo "RESULT: FAIL — one or more meta-tests failed"
  exit 1
fi
echo "RESULT: PASS — all meta-tests green"
exit 0

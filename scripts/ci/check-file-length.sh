#!/usr/bin/env bash
#
# check-file-length.sh
#
# Hard cap: a source file over 600 lines is doing too much — split it.
#
# Scope:
#   backend/**/*.go
#   frontend/src/**/*.{ts,tsx}
#   admin/src/**/*.{ts,tsx}
#
# Excluded entirely (never counted):
#   *_test.go, *.test.ts, *.test.tsx
#   frontend/e2e/**
#   generated / lock files (*_gen.go, *.gen.go, *_templ.go, *.lock, *-lock.*)
#
# Demo pages under frontend/src/app/(marketing)/demo are large BY DESIGN:
#   they are reported as a WARNING (does not fail the gate).
#
# Exit 0 = clean (only warnings allowed), exit 1 = a non-exempt file > 600 lines.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
MAX=600
DEMO_PREFIX="frontend/src/app/(marketing)/demo/"

violations=0
warnings=0

is_excluded() {
  rel="$1"
  case "$rel" in
    *_test.go|*.test.ts|*.test.tsx) return 0 ;;
    frontend/e2e/*)                 return 0 ;;
    *_gen.go|*.gen.go|*_templ.go)   return 0 ;;
    *.lock|*-lock.json|*-lock.yaml|*-lock.yml) return 0 ;;
  esac
  return 1
}

check_file() {
  f="$1"
  rel="${f#"$ROOT"/}"
  is_excluded "$rel" && return 0
  n="$(wc -l < "$f" 2>/dev/null | tr -d ' ')"
  [ -n "$n" ] || return 0
  [ "$n" -gt "$MAX" ] || return 0

  case "$rel" in
    "$DEMO_PREFIX"*)
      printf '::warning file=%s,line=%s::demo page is %s lines (>%s) — large by design, exempt from the hard cap\n' "$rel" "$n" "$n" "$MAX"
      printf 'WARNING %s — %s lines (>%s) [demo page, exempt]\n' "$rel" "$n" "$MAX"
      warnings=$((warnings + 1))
      ;;
    *)
      printf '::error file=%s,line=%s::file is %s lines (>%s) — split it by responsibility\n' "$rel" "$n" "$n" "$MAX"
      printf 'VIOLATION %s — %s lines (>%s)\n' "$rel" "$n" "$MAX"
      violations=$((violations + 1))
      ;;
  esac
}

scan() {
  base="$1"; shift
  [ -d "$base" ] || return 0
  while IFS= read -r f; do
    [ -n "$f" ] || continue
    check_file "$f"
  done <<EOF
$(find "$base" -type f \( "$@" \) 2>/dev/null)
EOF
}

scan "$ROOT/backend"       -name '*.go'
scan "$ROOT/frontend/src"  -name '*.ts' -o -name '*.tsx'
scan "$ROOT/admin/src"     -name '*.ts' -o -name '*.tsx'

echo "----------------------------------------"
[ "$warnings" -ne 0 ] && echo "($warnings demo file(s) over the cap — warning only)"
if [ "$violations" -ne 0 ]; then
  echo "FAIL: check-file-length — $violations file(s) over $MAX lines"
  exit 1
fi
echo "PASS: check-file-length — no non-exempt file over $MAX lines"
exit 0

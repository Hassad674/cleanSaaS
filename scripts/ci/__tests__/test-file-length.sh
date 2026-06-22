#!/usr/bin/env bash
# Meta-test for check-file-length.sh
set -euo pipefail
. "$(cd "$(dirname "$0")" && pwd)/_helpers.sh"

CHECK="check-file-length.sh"
echo "== $CHECK =="

# helper: write a file with N identical lines (avoid 'yes | head' SIGPIPE under pipefail)
write_lines() {
  f="$1"; n="$2"; mkdir -p "$(dirname "$f")"
  awk -v n="$n" 'BEGIN { for (i = 0; i < n; i++) print "x" }' > "$f"
}

seed_clean() {
  d="$1"
  write_lines "$d/backend/internal/app/auth/service.go" 100
  write_lines "$d/frontend/src/components/card.tsx" 200
}

# ---- clean fixture passes ----
CLEAN="$(mk_fixture "$CHECK")"; seed_clean "$CLEAN"
run_check "$CLEAN" "$CHECK"
assert_exit 0 "clean fixture (small files) passes"

# ---- oversized go file fails ----
L1="$(mk_fixture "$CHECK")"; seed_clean "$L1"
write_lines "$L1/backend/internal/app/billing/service.go" 700
run_check "$L1" "$CHECK"
assert_exit 1 "700-line .go file fails"
assert_annotation "backend/internal/app/billing/service.go" "names offending go file"
assert_contains "700 lines" "human-readable line count reported"

# ---- oversized tsx file fails ----
L2="$(mk_fixture "$CHECK")"; seed_clean "$L2"
write_lines "$L2/frontend/src/features/ai/components/chat.tsx" 650
run_check "$L2" "$CHECK"
assert_exit 1 "650-line .tsx file fails"
assert_annotation "frontend/src/features/ai/components/chat.tsx" "names offending tsx file"

# ---- exclusion: oversized *_test.go must NOT fail ----
L3="$(mk_fixture "$CHECK")"; seed_clean "$L3"
write_lines "$L3/backend/internal/app/auth/service_test.go" 900
run_check "$L3" "$CHECK"
assert_exit 0 "oversized *_test.go is excluded"

# ---- demo page is a WARNING, not a failure ----
L4="$(mk_fixture "$CHECK")"; seed_clean "$L4"
write_lines "$L4/frontend/src/app/(marketing)/demo/storage/storage-demo.tsx" 1000
run_check "$L4" "$CHECK"
assert_exit 0 "oversized demo page does NOT fail the gate"
assert_contains "::warning" "demo page emitted as warning annotation"
assert_contains "storage-demo.tsx" "demo page named in warning"

finish

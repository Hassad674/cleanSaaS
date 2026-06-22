#!/usr/bin/env bash
# Meta-test for check-migration-pairs.sh
set -euo pipefail
. "$(cd "$(dirname "$0")" && pwd)/_helpers.sh"

CHECK="check-migration-pairs.sh"
echo "== $CHECK =="

seed_clean() {
  d="$1"
  mkdir -p "$d/backend/migrations"
  for n in 001_create_users 002_create_billing; do
    printf 'CREATE TABLE x();\n' > "$d/backend/migrations/$n.up.sql"
    printf 'DROP TABLE x;\n'     > "$d/backend/migrations/$n.down.sql"
  done
}

# ---- clean fixture passes ----
CLEAN="$(mk_fixture "$CHECK")"; seed_clean "$CLEAN"
run_check "$CLEAN" "$CHECK"
assert_exit 0 "clean fixture passes"

# ---- up without down ----
M1="$(mk_fixture "$CHECK")"; seed_clean "$M1"
printf 'CREATE TABLE y();\n' > "$M1/backend/migrations/003_create_orphan.up.sql"
run_check "$M1" "$CHECK"
assert_exit 1 "missing .down.sql fails"
assert_annotation "003_create_orphan.up.sql" "names offending up file"
assert_contains "no matching .down.sql" "human-readable missing-down msg"

# ---- down without up ----
M2="$(mk_fixture "$CHECK")"; seed_clean "$M2"
printf 'DROP TABLE z;\n' > "$M2/backend/migrations/004_create_orphan.down.sql"
run_check "$M2" "$CHECK"
assert_exit 1 "missing .up.sql fails"
assert_annotation "004_create_orphan.down.sql" "names offending down file"

finish

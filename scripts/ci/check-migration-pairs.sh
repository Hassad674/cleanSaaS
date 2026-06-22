#!/usr/bin/env bash
#
# check-migration-pairs.sh
#
# Every backend/migrations/NNN_*.up.sql must have a matching *.down.sql,
# and every *.down.sql must have a matching *.up.sql.
#
# Exit 0 = clean, exit 1 = a migration is missing its pair.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
MIG_DIR="$ROOT/backend/migrations"

violations=0

report() {
  file="$1"; msg="$2"
  case "$file" in
    "$ROOT"/*) rel="${file#"$ROOT"/}" ;;
    *)         rel="$file" ;;
  esac
  printf '::error file=%s,line=1::%s\n' "$rel" "$msg"
  printf 'VIOLATION %s — %s\n' "$rel" "$msg"
  violations=$((violations + 1))
}

if [ ! -d "$MIG_DIR" ]; then
  echo "PASS: check-migration-pairs — no migrations directory ($MIG_DIR)"
  exit 0
fi

# Every .up.sql needs a .down.sql
while IFS= read -r up; do
  [ -n "$up" ] || continue
  down="${up%.up.sql}.down.sql"
  if [ ! -f "$down" ]; then
    report "$up" "migration '$(basename "$up")' has no matching .down.sql (every up migration must be reversible)"
  fi
done <<EOF
$(find "$MIG_DIR" -maxdepth 1 -type f -name '*.up.sql' 2>/dev/null | sort)
EOF

# Every .down.sql needs an .up.sql
while IFS= read -r down; do
  [ -n "$down" ] || continue
  up="${down%.down.sql}.up.sql"
  if [ ! -f "$up" ]; then
    report "$down" "migration '$(basename "$down")' has no matching .up.sql (orphaned down migration)"
  fi
done <<EOF
$(find "$MIG_DIR" -maxdepth 1 -type f -name '*.down.sql' 2>/dev/null | sort)
EOF

echo "----------------------------------------"
if [ "$violations" -ne 0 ]; then
  echo "FAIL: check-migration-pairs — $violations unpaired migration file(s)"
  exit 1
fi
echo "PASS: check-migration-pairs — all migrations have up/down pairs"
exit 0

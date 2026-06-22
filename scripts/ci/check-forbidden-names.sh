#!/usr/bin/env bash
#
# check-forbidden-names.sh
#
# CLAUDE.md forbids indescriptive identifiers. We enforce a CONSERVATIVE subset
# that is safe to grep without false positives. A name is flagged ONLY when it
# is, as a whole word:
#
#   - a Go function name:        func doStuff(...        |  func (x *T) Manager(...
#   - a Go type name:            type Manager struct     |  type Util interface
#   - a package-level var/const: var data ...            |  const tmp ...
#
# It is NOT flagged when it appears as:
#   - a struct field            (those are indented, never matched by ^func/type/var/const)
#   - a json tag                (`json:"data"` — never at line start as a decl)
#   - a local variable / map key
#   - a substring               ('metadata', 'information', 'template' must NOT match,
#                                 guaranteed by the \b word boundary)
#   - a function PARAMETER       ('func renderTemplate(data ...)' — name is renderTemplate)
#
# Forbidden list (case-insensitive on the whole word, so both `manager` and
# `Manager` match, but `metadata` does not):
#   data info tmp temp manager util utils helper helpers handler2 doStuff
#
# Scope: backend/**/*.go, excluding *_test.go.
# Exit 0 = clean, exit 1 = a forbidden identifier is declared.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
BACKEND="$ROOT/backend"

# Alternation of forbidden words. Matched case-insensitively as a whole word.
FORBIDDEN='data|info|tmp|temp|manager|util|utils|helper|helpers|handler2|doStuff'

violations=0

report() {
  file="$1"; line="$2"; kind="$3"; name="$4"
  case "$file" in
    "$ROOT"/*) rel="${file#"$ROOT"/}" ;;
    *)         rel="$file" ;;
  esac
  printf '::error file=%s,line=%s::forbidden indescriptive %s name "%s" — every name must state what it is\n' "$rel" "$line" "$kind" "$name"
  printf 'VIOLATION %s:%s — forbidden %s name "%s"\n' "$rel" "$line" "$kind" "$name"
  violations=$((violations + 1))
}

if [ ! -d "$BACKEND" ]; then
  echo "PASS: check-forbidden-names — no backend directory"
  exit 0
fi

# Patterns (case-insensitive whole-word, anchored at line start):
#   func  name(            and   func (recv) name(
#   type  name <kind>
#   var/const  name
FUNC_RE="^func[[:space:]]+(\([^)]*\)[[:space:]]+)?(${FORBIDDEN})[[:space:]]*\("
TYPE_RE="^type[[:space:]]+(${FORBIDDEN})[[:space:]]"
VARCONST_RE="^(var|const)[[:space:]]+(${FORBIDDEN})[[:space:]]"

while IFS= read -r f; do
  [ -n "$f" ] || continue

  # functions
  while IFS=: read -r lineno content; do
    [ -n "$lineno" ] || continue
    name="$(printf '%s' "$content" | sed -E 's/^func[[:space:]]+(\([^)]*\)[[:space:]]+)?([A-Za-z0-9_]+).*/\2/')"
    report "$f" "$lineno" "function" "$name"
  done <<EOF
$(grep -niE "$FUNC_RE" "$f" 2>/dev/null || true)
EOF

  # types
  while IFS=: read -r lineno content; do
    [ -n "$lineno" ] || continue
    name="$(printf '%s' "$content" | sed -E 's/^type[[:space:]]+([A-Za-z0-9_]+).*/\1/')"
    report "$f" "$lineno" "type" "$name"
  done <<EOF
$(grep -niE "$TYPE_RE" "$f" 2>/dev/null || true)
EOF

  # package-level var / const
  while IFS=: read -r lineno content; do
    [ -n "$lineno" ] || continue
    name="$(printf '%s' "$content" | sed -E 's/^(var|const)[[:space:]]+([A-Za-z0-9_]+).*/\2/')"
    report "$f" "$lineno" "var/const" "$name"
  done <<EOF
$(grep -niE "$VARCONST_RE" "$f" 2>/dev/null || true)
EOF

done <<EOF
$(find "$BACKEND" -type f -name '*.go' ! -name '*_test.go' 2>/dev/null)
EOF

echo "----------------------------------------"
if [ "$violations" -ne 0 ]; then
  echo "FAIL: check-forbidden-names — $violations forbidden identifier(s) declared"
  exit 1
fi
echo "PASS: check-forbidden-names — no forbidden indescriptive identifiers"
exit 0

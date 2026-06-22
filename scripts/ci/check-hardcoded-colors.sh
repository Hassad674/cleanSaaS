#!/usr/bin/env bash
#
# check-hardcoded-colors.sh
#
# The project uses semantic Tailwind tokens (bg-card, text-foreground, …).
# Hardcoded palette colors are forbidden in className strings:
#   - zinc-*, gray-*, slate-*, neutral-*       (e.g. bg-zinc-900, text-gray-500)
#   - bare white / black color utilities       (e.g. bg-white, text-black, border-white)
#
# Scope: frontend/src/** and admin/src/**  (*.ts and *.tsx).
#
# Care is taken to avoid false positives:
#   - 'whitespace-*' is NOT a color (the white/black regex requires a real color
#     utility prefix immediately followed by white/black as a whole word).
#   - palette colors are only flagged with a real Tailwind utility prefix
#     (bg-/text-/border-/ring-/...), so identifiers like 'graystone' won't match.
#
# Exit 0 = clean, exit 1 = hardcoded color found.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"

violations=0

report() {
  file="$1"; line="$2"; token="$3"
  case "$file" in
    "$ROOT"/*) rel="${file#"$ROOT"/}" ;;
    *)         rel="$file" ;;
  esac
  printf '::error file=%s,line=%s::hardcoded Tailwind color "%s" — use a semantic token (bg-card, text-foreground, border-border, …)\n' "$rel" "$line" "$token"
  printf 'VIOLATION %s:%s — hardcoded color "%s"\n' "$rel" "$line" "$token"
  violations=$((violations + 1))
}

# Tailwind utility prefixes that take a color value.
PREFIX='(bg|text|border|ring|ring-offset|divide|outline|fill|stroke|from|to|via|decoration|placeholder|accent|caret|shadow)'
# A color utility, optionally with a variant chain (hover:, dark:, md:, etc.) and
# an optional leading '!' important marker. Anchored on a non-identifier boundary.
PALETTE="(^|[^A-Za-z0-9_-])([a-z]+:)*!?${PREFIX}-(zinc|gray|slate|neutral)-[0-9]+"
# bare white/black as a color value (whole word -> 'whitespace-' is safe because
# it has no color-utility prefix and 'white' here must be a standalone token).
WHITEBLACK="(^|[^A-Za-z0-9_-])([a-z]+:)*!?${PREFIX}-(white|black)([^A-Za-z0-9_-]|$)"

scan_dir() {
  dir="$1"
  [ -d "$dir" ] || return 0
  while IFS= read -r f; do
    [ -n "$f" ] || continue
    while IFS=: read -r lineno content; do
      [ -n "$lineno" ] || continue
      # palette colors
      tok="$(printf '%s' "$content" | grep -oE "${PREFIX}-(zinc|gray|slate|neutral)-[0-9]+" | head -1 || true)"
      [ -n "$tok" ] && report "$f" "$lineno" "$tok"
    done <<EOF
$(grep -nE "$PALETTE" "$f" 2>/dev/null || true)
EOF
    while IFS=: read -r lineno content; do
      [ -n "$lineno" ] || continue
      tok="$(printf '%s' "$content" | grep -oE "${PREFIX}-(white|black)" | head -1 || true)"
      [ -n "$tok" ] && report "$f" "$lineno" "$tok"
    done <<EOF
$(grep -nE "$WHITEBLACK" "$f" 2>/dev/null || true)
EOF
  done <<EOF
$(find "$dir" -type f \( -name '*.tsx' -o -name '*.ts' \) ! -name '*.test.ts' ! -name '*.test.tsx' 2>/dev/null)
EOF
}

scan_dir "$ROOT/frontend/src"
scan_dir "$ROOT/admin/src"

echo "----------------------------------------"
if [ "$violations" -ne 0 ]; then
  echo "FAIL: check-hardcoded-colors — $violations hardcoded color(s) found"
  exit 1
fi
echo "PASS: check-hardcoded-colors — no hardcoded Tailwind colors"
exit 0

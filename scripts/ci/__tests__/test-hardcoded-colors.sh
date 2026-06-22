#!/usr/bin/env bash
# Meta-test for check-hardcoded-colors.sh
set -euo pipefail
. "$(cd "$(dirname "$0")" && pwd)/_helpers.sh"

CHECK="check-hardcoded-colors.sh"
echo "== $CHECK =="

# Clean fixture uses ONLY semantic tokens + a whitespace-* utility (false-positive trap).
seed_clean() {
  d="$1"
  mkdir -p "$d/frontend/src/components" "$d/admin/src/components"
  cat > "$d/frontend/src/components/card.tsx" <<'EOF'
export const Card = () => (
  <div className="bg-card text-foreground border-border whitespace-pre-wrap">
    hello
  </div>
);
EOF
  cat > "$d/admin/src/components/panel.tsx" <<'EOF'
export const Panel = () => <div className="bg-background text-muted-foreground">x</div>;
EOF
}

# ---- clean fixture passes ----
CLEAN="$(mk_fixture "$CHECK")"; seed_clean "$CLEAN"
run_check "$CLEAN" "$CHECK"
assert_exit 0 "clean fixture (semantic tokens + whitespace-*) passes"

# ---- palette color (zinc) ----
C1="$(mk_fixture "$CHECK")"; seed_clean "$C1"
cat > "$C1/frontend/src/components/bad.tsx" <<'EOF'
export const Bad = () => <div className="bg-zinc-900 text-gray-500">x</div>;
EOF
run_check "$C1" "$CHECK"
assert_exit 1 "hardcoded palette color (zinc) fails"
assert_annotation "frontend/src/components/bad.tsx" "names offending file"
assert_contains "bg-zinc-900" "human-readable token reported"

# ---- bare white/black utility ----
C2="$(mk_fixture "$CHECK")"; seed_clean "$C2"
cat > "$C2/admin/src/components/bad.tsx" <<'EOF'
export const Bad = () => <div className="bg-white text-black">x</div>;
EOF
run_check "$C2" "$CHECK"
assert_exit 1 "bare white/black utility fails"
assert_annotation "admin/src/components/bad.tsx" "names offending admin file"

# ---- false-positive guard: whitespace-nowrap alone must NOT fail ----
C3="$(mk_fixture "$CHECK")"; seed_clean "$C3"
cat > "$C3/frontend/src/components/ws.tsx" <<'EOF'
export const Ws = () => <span className="whitespace-nowrap text-foreground">x</span>;
EOF
run_check "$C3" "$CHECK"
assert_exit 0 "whitespace-nowrap is NOT flagged as a color"

finish

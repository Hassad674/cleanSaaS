#!/usr/bin/env bash
# Meta-test for check-forbidden-names.sh
set -euo pipefail
. "$(cd "$(dirname "$0")" && pwd)/_helpers.sh"

CHECK="check-forbidden-names.sh"
echo "== $CHECK =="

# Clean fixture is FULL of near-miss false-positive traps that must NOT match:
#   - struct field named 'data'
#   - json tag "info"
#   - function PARAMETER named 'data' / 'template'
#   - function NAME containing forbidden words as substrings (renderTemplate, GetMetadata)
#   - local variable 'tmp'
seed_clean() {
  d="$1"
  mkdir -p "$d/backend/internal/adapter/resend"
  cat > "$d/backend/internal/adapter/resend/email.go" <<'EOF'
package resend

type Payload struct {
	Data     string `json:"data"`
	Info     string `json:"info"`
	Metadata string `json:"metadata"`
}

func renderTemplate(template string, data map[string]string) (string, string) {
	tmp := template + data["k"]
	return tmp, tmp
}

func (s *Service) GetMetadata(info string) string {
	return info
}
EOF
}

# ---- clean fixture passes (no false positives on fields/tags/params/substrings) ----
CLEAN="$(mk_fixture "$CHECK")"; seed_clean "$CLEAN"
run_check "$CLEAN" "$CHECK"
assert_exit 0 "clean fixture (fields/tags/params/substrings) passes — no false positives"

# ---- forbidden function name: func doStuff ----
N1="$(mk_fixture "$CHECK")"; seed_clean "$N1"
cat >> "$N1/backend/internal/adapter/resend/email.go" <<'EOF'

func doStuff() {}
EOF
run_check "$N1" "$CHECK"
assert_exit 1 "func doStuff fails"
assert_annotation "backend/internal/adapter/resend/email.go" "names offending file (func)"
assert_contains "doStuff" "human-readable name reported"

# ---- forbidden method name: func (x *T) Manager( ----
N2="$(mk_fixture "$CHECK")"; seed_clean "$N2"
cat >> "$N2/backend/internal/adapter/resend/email.go" <<'EOF'

func (s *Service) Manager() {}
EOF
run_check "$N2" "$CHECK"
assert_exit 1 "func (x) Manager() fails"
assert_contains "Manager" "method name Manager reported"

# ---- forbidden type name: type Manager struct ----
N3="$(mk_fixture "$CHECK")"; seed_clean "$N3"
cat >> "$N3/backend/internal/adapter/resend/email.go" <<'EOF'

type Manager struct{ x int }
EOF
run_check "$N3" "$CHECK"
assert_exit 1 "type Manager struct fails"
assert_contains "type" "type kind reported"

# ---- forbidden package-level var: var data ----
N4="$(mk_fixture "$CHECK")"; seed_clean "$N4"
cat >> "$N4/backend/internal/adapter/resend/email.go" <<'EOF'

var data = 1
EOF
run_check "$N4" "$CHECK"
assert_exit 1 "package-level var data fails"
assert_contains "var/const" "var/const kind reported"

# ---- exclusion: forbidden name in *_test.go must NOT fail ----
N5="$(mk_fixture "$CHECK")"; seed_clean "$N5"
cat > "$N5/backend/internal/adapter/resend/email_test.go" <<'EOF'
package resend

func doStuff() {}
EOF
run_check "$N5" "$CHECK"
assert_exit 0 "forbidden name inside *_test.go is excluded"

finish

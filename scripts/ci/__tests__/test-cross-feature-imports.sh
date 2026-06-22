#!/usr/bin/env bash
# Meta-test for check-cross-feature-imports.sh
set -euo pipefail
. "$(cd "$(dirname "$0")" && pwd)/_helpers.sh"

CHECK="check-cross-feature-imports.sh"
echo "== $CHECK =="

GO_HDR='package billing'

# Build a minimal clean fixture: two backend app features (auth, billing),
# two adapters (postgres, stripe), a pure domain, and two frontend features.
seed_clean() {
  d="$1"
  mkdir -p "$d/backend/internal/app/auth" "$d/backend/internal/app/billing"
  mkdir -p "$d/backend/internal/adapter/postgres" "$d/backend/internal/adapter/stripe"
  mkdir -p "$d/backend/internal/domain/user" "$d/backend/internal/port"
  mkdir -p "$d/frontend/src/features/auth/components" "$d/frontend/src/features/billing/components"

  cat > "$d/backend/internal/app/billing/service.go" <<'EOF'
package billing

import (
	"github.com/hassad/boilerplateSaaS/backend/internal/domain/user"
	"github.com/hassad/boilerplateSaaS/backend/internal/port"
)

type Service struct{ _ port.X; _ user.User }
EOF
  cat > "$d/backend/internal/app/auth/service.go" <<'EOF'
package auth

import "github.com/hassad/boilerplateSaaS/backend/internal/domain/user"

type Service struct{ _ user.User }
EOF
  cat > "$d/backend/internal/adapter/postgres/user.go" <<'EOF'
package postgres

import "github.com/hassad/boilerplateSaaS/backend/internal/port"

type Repo struct{ _ port.X }
EOF
  cat > "$d/backend/internal/domain/user/entity.go" <<'EOF'
package user

import "github.com/hassad/boilerplateSaaS/backend/internal/domain/errors"

type User struct{ _ errors.E }
EOF
  cat > "$d/frontend/src/features/billing/components/card.tsx" <<'EOF'
import type { Plan } from "@/features/billing/types";
export const X = (p: Plan) => p;
EOF
  cat > "$d/frontend/src/features/auth/components/form.tsx" <<'EOF'
import { login } from "@/features/auth/actions/auth";
export const Y = login;
EOF
}

# ---- clean fixture passes ----
CLEAN="$(mk_fixture "$CHECK")"; seed_clean "$CLEAN"
run_check "$CLEAN" "$CHECK"
assert_exit 0 "clean fixture passes"

# ---- backend: app imports another app feature ----
B1="$(mk_fixture "$CHECK")"; seed_clean "$B1"
cat > "$B1/backend/internal/app/billing/service.go" <<'EOF'
package billing

import "github.com/hassad/boilerplateSaaS/backend/internal/app/auth"

type Service struct{ _ auth.Service }
EOF
run_check "$B1" "$CHECK"
assert_exit 1 "app->app cross import fails"
assert_annotation "internal/app/billing/service.go" "names offending app file"
assert_contains "internal/app/billing imports internal/app/auth" "human-readable app->app msg"

# ---- backend: adapter imports another adapter ----
B2="$(mk_fixture "$CHECK")"; seed_clean "$B2"
cat > "$B2/backend/internal/adapter/postgres/user.go" <<'EOF'
package postgres

import "github.com/hassad/boilerplateSaaS/backend/internal/adapter/stripe"

type Repo struct{ _ stripe.Client }
EOF
run_check "$B2" "$CHECK"
assert_exit 1 "adapter->adapter cross import fails"
assert_annotation "internal/adapter/postgres/user.go" "names offending adapter file"

# ---- backend: domain imports app (impurity) ----
B3="$(mk_fixture "$CHECK")"; seed_clean "$B3"
cat > "$B3/backend/internal/domain/user/entity.go" <<'EOF'
package user

import "github.com/hassad/boilerplateSaaS/backend/internal/app/auth"

type User struct{ _ auth.Service }
EOF
run_check "$B3" "$CHECK"
assert_exit 1 "domain importing app fails"
assert_annotation "internal/domain/user/entity.go" "names offending domain file"

# ---- frontend: feature imports another feature via @/features alias ----
F1="$(mk_fixture "$CHECK")"; seed_clean "$F1"
cat > "$F1/frontend/src/features/billing/components/card.tsx" <<'EOF'
import { login } from "@/features/auth/actions/auth";
export const X = login;
EOF
run_check "$F1" "$CHECK"
assert_exit 1 "frontend alias cross import fails"
assert_annotation "frontend/src/features/billing/components/card.tsx" "names offending frontend file"
assert_contains "imports @/features/auth" "human-readable frontend msg"

# ---- frontend: feature imports another feature via relative path ----
F2="$(mk_fixture "$CHECK")"; seed_clean "$F2"
cat > "$F2/frontend/src/features/billing/components/card.tsx" <<'EOF'
import { login } from "../../auth/actions/auth";
export const X = login;
EOF
run_check "$F2" "$CHECK"
assert_exit 1 "frontend relative cross import fails"
assert_annotation "frontend/src/features/billing/components/card.tsx" "names offending file (relative)"

finish

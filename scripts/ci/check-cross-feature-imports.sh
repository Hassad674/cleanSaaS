#!/usr/bin/env bash
#
# check-cross-feature-imports.sh
#
# Enforces feature/layer isolation that golangci-lint cannot express:
#
#   Backend (module github.com/hassad/boilerplateSaaS/backend):
#     - A file under internal/app/<A>/ must not import internal/app/<B>/   (B != A)
#     - A file under internal/adapter/<A>/ must not import internal/adapter/<B>/ (B != A)
#     - Domain purity: no file under internal/domain/ imports internal/app,
#       internal/adapter or internal/handler. (Other internal/domain/... and
#       internal/port imports are allowed.)
#
#   Frontend:
#     - A file under frontend/src/features/<A>/ must not import @/features/<B>
#       nor a relative path that resolves into another feature <B>   (B != A)
#
# Exit 0 = clean, exit 1 = violation(s) found.
# Each violation prints a GitHub annotation (::error file=...,line=...::msg)
# plus a human-readable line.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
MODULE="github.com/hassad/boilerplateSaaS/backend"

violations=0

# Emit one violation: $1=file (relative or absolute), $2=line, $3=message
report() {
  file="$1"; line="$2"; msg="$3"
  # Make path repo-relative for readability if it is under ROOT.
  case "$file" in
    "$ROOT"/*) rel="${file#"$ROOT"/}" ;;
    *)         rel="$file" ;;
  esac
  printf '::error file=%s,line=%s::%s\n' "$rel" "$line" "$msg"
  printf 'VIOLATION %s:%s — %s\n' "$rel" "$line" "$msg"
  violations=$((violations + 1))
}

############################################################
# Backend
############################################################
check_backend() {
  app_dir="$ROOT/backend/internal/app"
  adapter_dir="$ROOT/backend/internal/adapter"
  domain_dir="$ROOT/backend/internal/domain"

  # --- app/<A> must not import app/<B> (B != A) ---
  if [ -d "$app_dir" ]; then
    while IFS= read -r f; do
      # owning feature = first path segment under app/
      rest="${f#"$app_dir"/}"
      owner="${rest%%/*}"
      # find imported app features in this file
      while IFS=: read -r lineno imp; do
        [ -n "$imp" ] || continue
        # imp looks like: internal/app/<B>/...  -> extract <B>
        b="${imp#*internal/app/}"
        b="${b%%/*}"
        b="${b%%\"*}"   # strip trailing quote if no slash
        [ -n "$b" ] || continue
        if [ "$b" != "$owner" ]; then
          report "$f" "$lineno" "internal/app/$owner imports internal/app/$b — features must not import each other (inject an interface instead)"
        fi
      done <<EOF
$(grep -nE "\"$MODULE/internal/app/" "$f" 2>/dev/null || true)
EOF
    done <<EOF
$(find "$app_dir" -type f -name '*.go' ! -name '*_test.go' 2>/dev/null)
EOF
  fi

  # --- adapter/<A> must not import adapter/<B> (B != A) ---
  if [ -d "$adapter_dir" ]; then
    while IFS= read -r f; do
      rest="${f#"$adapter_dir"/}"
      owner="${rest%%/*}"
      while IFS=: read -r lineno imp; do
        [ -n "$imp" ] || continue
        b="${imp#*internal/adapter/}"
        b="${b%%/*}"
        b="${b%%\"*}"
        [ -n "$b" ] || continue
        if [ "$b" != "$owner" ]; then
          report "$f" "$lineno" "internal/adapter/$owner imports internal/adapter/$b — an adapter must never import another adapter"
        fi
      done <<EOF
$(grep -nE "\"$MODULE/internal/adapter/" "$f" 2>/dev/null || true)
EOF
    done <<EOF
$(find "$adapter_dir" -type f -name '*.go' ! -name '*_test.go' 2>/dev/null)
EOF
  fi

  # --- domain purity: no domain file imports app/adapter/handler ---
  if [ -d "$domain_dir" ]; then
    while IFS= read -r f; do
      while IFS=: read -r lineno imp; do
        [ -n "$imp" ] || continue
        layer="${imp#*internal/}"
        layer="${layer%%/*}"
        layer="${layer%%\"*}"
        report "$f" "$lineno" "internal/domain imports internal/$layer — domain must be pure (only stdlib, internal/domain, internal/port allowed)"
      done <<EOF
$(grep -nE "\"$MODULE/internal/(app|adapter|handler)(/|\")" "$f" 2>/dev/null || true)
EOF
    done <<EOF
$(find "$domain_dir" -type f -name '*.go' ! -name '*_test.go' 2>/dev/null)
EOF
  fi
}

############################################################
# Frontend
############################################################
check_frontend() {
  feat_dir="$ROOT/frontend/src/features"
  [ -d "$feat_dir" ] || return 0

  # list of real feature names (for resolving relative imports)
  features=""
  while IFS= read -r d; do
    [ -n "$d" ] || continue
    features="$features ${d##*/}"
  done <<EOF
$(find "$feat_dir" -mindepth 1 -maxdepth 1 -type d 2>/dev/null)
EOF

  while IFS= read -r f; do
    [ -n "$f" ] || continue
    rest="${f#"$feat_dir"/}"
    owner="${rest%%/*}"

    # --- alias imports: @/features/<B> ---
    while IFS=: read -r lineno imp; do
      [ -n "$imp" ] || continue
      b="${imp#*@/features/}"
      b="${b%%/*}"
      b="${b%%\'*}"; b="${b%%\"*}"; b="${b%%\`*}"
      [ -n "$b" ] || continue
      if [ "$b" != "$owner" ]; then
        report "$f" "$lineno" "frontend feature '$owner' imports @/features/$b — features must not import each other (compose in app/ pages)"
      fi
    done <<EOF
$(grep -nE "from[[:space:]]*['\"\`]@/features/" "$f" 2>/dev/null || true)
EOF

    # --- relative imports that climb out of the feature: ../../<B> ---
    # only flag when the path leaves features/<owner> and lands in another feature.
    while IFS=: read -r lineno imp; do
      [ -n "$imp" ] || continue
      # extract the quoted module path
      path="$(printf '%s' "$imp" | sed -E "s/.*from[[:space:]]*['\"\`]([^'\"\`]+)['\"\`].*/\1/")"
      case "$path" in
        ../*)
          # Resolve the import LEXICALLY (no filesystem dependency) against the
          # file's directory, then see which feature (if any) it lands in.
          filedir="$(dirname "$f")"
          combined="$filedir/$path"
          # Normalize away '.' and '..' segments lexically.
          resolved=""
          oldIFS="$IFS"; IFS='/'
          for seg in $combined; do
            case "$seg" in
              ''|'.') ;;
              '..')   resolved="${resolved%/*}" ;;
              *)      resolved="$resolved/$seg" ;;
            esac
          done
          IFS="$oldIFS"
          case "$resolved" in
            "$feat_dir"/*)
              tgtrest="${resolved#"$feat_dir"/}"
              tgt="${tgtrest%%/*}"
              if [ -n "$tgt" ] && [ "$tgt" != "$owner" ]; then
                report "$f" "$lineno" "frontend feature '$owner' imports feature '$tgt' via relative path '$path' — features must not import each other"
              fi
              ;;
          esac
          ;;
      esac
    done <<EOF
$(grep -nE "from[[:space:]]*['\"\`]\.\./" "$f" 2>/dev/null || true)
EOF
  done <<EOF
$(find "$feat_dir" -type f \( -name '*.ts' -o -name '*.tsx' \) ! -name '*.test.ts' ! -name '*.test.tsx' 2>/dev/null)
EOF
}

check_backend
check_frontend

echo "----------------------------------------"
if [ "$violations" -ne 0 ]; then
  echo "FAIL: check-cross-feature-imports — $violations cross-feature/layer import violation(s) found"
  exit 1
fi
echo "PASS: check-cross-feature-imports — no cross-feature/layer imports"
exit 0

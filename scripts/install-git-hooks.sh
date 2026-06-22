#!/usr/bin/env bash
#
# install-git-hooks.sh — point git at the repo's tracked .githooks/ directory.
#
# Git does NOT auto-install hooks on clone (they live outside the tree by
# default), so every contributor must run this once per clone. It is fully
# idempotent: re-running just re-confirms the configuration.
#
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null)" || {
  echo "install-git-hooks: must be run inside the git repository." >&2
  exit 1
}
cd "$REPO_ROOT"

HOOKS_DIR=".githooks"

if [ -t 1 ]; then
  GREEN='\033[0;32m'; YELLOW='\033[1;33m'; BOLD='\033[1m'; NC='\033[0m'
else
  GREEN=''; YELLOW=''; BOLD=''; NC=''
fi

if [ ! -d "$HOOKS_DIR" ]; then
  echo "install-git-hooks: $REPO_ROOT/$HOOKS_DIR does not exist." >&2
  exit 1
fi

# Ensure every hook in the dir is executable (idempotent).
find "$HOOKS_DIR" -maxdepth 1 -type f -exec chmod +x {} +

CURRENT="$(git config --local --get core.hooksPath || true)"
if [ "$CURRENT" = "$HOOKS_DIR" ]; then
  printf "${GREEN}[install-git-hooks]${NC} core.hooksPath already set to %s — nothing to do.\n" "$HOOKS_DIR"
else
  git config core.hooksPath "$HOOKS_DIR"
  printf "${GREEN}[install-git-hooks]${NC} configured core.hooksPath = %s\n" "$HOOKS_DIR"
fi

printf "\n"
printf "${BOLD}Git hooks installed for this clone.${NC}\n"
printf "  Active hooks: %s\n" "$(find "$HOOKS_DIR" -maxdepth 1 -type f -printf '%f ' 2>/dev/null || echo '(none)')"
printf "\n"
printf "${YELLOW}Reminder:${NC} hooks are opt-in PER CLONE. Git never installs them\n"
printf "  automatically, so each contributor must run this script once after cloning:\n"
printf "      bash scripts/install-git-hooks.sh\n"
printf "\n"
printf "${BOLD}To bypass the pre-commit hook for a single commit:${NC}\n"
printf "      git commit --no-verify\n"
printf "\n"

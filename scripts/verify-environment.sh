#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/activate-env.sh"

fail=0

assert_external() {
  local label="$1"
  local value="$2"
  case "$value" in
    "$STACKGENOME_VOLUME"/*) ;;
    *)
      printf 'FAIL: %s is outside external volume: %s\n' "$label" "$value" >&2
      fail=1
      ;;
  esac
}

assert_external "STACKGENOME_REPO" "$STACKGENOME_REPO"
assert_external "STACKGENOME_PROGRAMS" "$STACKGENOME_PROGRAMS"
assert_external "STACKGENOME_CACHES" "$STACKGENOME_CACHES"
assert_external "MISE_DATA_DIR" "$MISE_DATA_DIR"
assert_external "MISE_INSTALLS_DIR" "$MISE_INSTALLS_DIR"
assert_external "MISE_CACHE_DIR" "$MISE_CACHE_DIR"
assert_external "GOPATH" "$GOPATH"
assert_external "GOMODCACHE" "$GOMODCACHE"
assert_external "GOCACHE" "$GOCACHE"
assert_external "GOTMPDIR" "$GOTMPDIR"

printf '\nTool checks:\n'
command -v mise || { echo "mise not found"; fail=1; }
command -v go || { echo "go not found"; fail=1; }

if command -v mise >/dev/null 2>&1; then mise --version; fi
if command -v go >/dev/null 2>&1; then
  go version
  go env GOPATH GOMODCACHE GOCACHE GOTMPDIR
fi

printf '\nDisk locations:\n'
printf 'Volume:    %s\n' "$STACKGENOME_VOLUME"
printf 'Repo:      %s\n' "$STACKGENOME_REPO"
printf 'Programs:  %s\n' "$STACKGENOME_PROGRAMS"
printf 'Caches:    %s\n' "$STACKGENOME_CACHES"

if [[ "$fail" -ne 0 ]]; then
  echo "Environment verification failed." >&2
  exit 1
fi

echo "Environment verification passed."

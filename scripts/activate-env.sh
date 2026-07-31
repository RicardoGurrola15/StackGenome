#!/usr/bin/env bash
# Source this file: source scripts/activate-env.sh

_stackgenome_fail() {
  printf 'StackGenome environment error: %s\n' "$1" >&2
  return 1 2>/dev/null || exit 1
}

export STACKGENOME_VOLUME="${STACKGENOME_VOLUME:-/Volumes/intento1}"
export STACKGENOME_REPO="${STACKGENOME_REPO:-$STACKGENOME_VOLUME/Repos/StackGenome}"
export STACKGENOME_PROGRAMS="${STACKGENOME_PROGRAMS:-$STACKGENOME_VOLUME/programas}"
export STACKGENOME_CACHES="${STACKGENOME_CACHES:-$STACKGENOME_VOLUME/caches}"

[[ -d "$STACKGENOME_VOLUME" ]] || _stackgenome_fail \
  "volume $STACKGENOME_VOLUME is not mounted; refusing to fall back to the main disk."

mkdir -p \
  "$STACKGENOME_PROGRAMS/mise/bin" \
  "$STACKGENOME_PROGRAMS/mise/data" \
  "$STACKGENOME_PROGRAMS/mise/installs" \
  "$STACKGENOME_PROGRAMS/go-workspace/bin" \
  "$STACKGENOME_CACHES/mise" \
  "$STACKGENOME_CACHES/go/mod" \
  "$STACKGENOME_CACHES/go/build" \
  "$STACKGENOME_CACHES/go/tmp" 2>/dev/null || true

export MISE_INSTALL_PATH="$STACKGENOME_PROGRAMS/mise/bin/mise"
export MISE_DATA_DIR="$STACKGENOME_PROGRAMS/mise/data"
export MISE_INSTALLS_DIR="$STACKGENOME_PROGRAMS/mise/installs"
export MISE_CACHE_DIR="$STACKGENOME_CACHES/mise"

export GOPATH="$STACKGENOME_PROGRAMS/go-workspace"
export GOMODCACHE="$STACKGENOME_CACHES/go/mod"
export GOCACHE="$STACKGENOME_CACHES/go/build"
export GOTMPDIR="$STACKGENOME_CACHES/go/tmp"
export XDG_CACHE_HOME="$STACKGENOME_CACHES"

export PATH="$STACKGENOME_PROGRAMS/mise/bin:$MISE_DATA_DIR/shims:$GOPATH/bin:$PATH"

unset -f _stackgenome_fail

printf 'StackGenome environment activated.\n'
printf '  Repo:       %s\n' "$STACKGENOME_REPO"
printf '  Programs:   %s\n' "$STACKGENOME_PROGRAMS"
printf '  Caches:     %s\n' "$STACKGENOME_CACHES"

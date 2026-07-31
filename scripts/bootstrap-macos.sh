#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd -- "$SCRIPT_DIR/.." && pwd)"

# shellcheck source=activate-env.sh
source "$SCRIPT_DIR/activate-env.sh"

if [[ "$(uname -s)" != "Darwin" ]]; then
  echo "This bootstrap is for macOS. The project itself remains cross-platform." >&2
  exit 2
fi

if [[ "$REPO_DIR" != "$STACKGENOME_REPO" ]]; then
  echo "Warning: repository is at $REPO_DIR, expected $STACKGENOME_REPO." >&2
  echo "No product code may depend on the expected absolute path." >&2
fi

MISE_BIN="$STACKGENOME_PROGRAMS/mise/bin/mise"

if [[ ! -x "$MISE_BIN" ]]; then
  echo "mise is not installed at: $MISE_BIN"
  echo "Planned source: https://mise.run"
  echo "Planned destination: $MISE_BIN"
  read -r -p "Install mise there now? [y/N] " answer
  case "$answer" in
    y|Y|yes|YES)
      curl --fail --location --proto '=https' --tlsv1.2 https://mise.run |
        MISE_INSTALL_PATH="$MISE_BIN" sh
      ;;
    *)
      echo "Installation cancelled."
      exit 9
      ;;
  esac
fi

"$MISE_BIN" --version
"$MISE_BIN" trust "$REPO_DIR/mise.toml"
"$MISE_BIN" install -C "$REPO_DIR"

echo
echo "Installed tool locations:"
"$MISE_BIN" where go || true
echo
echo "Run:"
echo "  source \"$REPO_DIR/scripts/activate-env.sh\""
echo "  \"$REPO_DIR/scripts/verify-environment.sh\""

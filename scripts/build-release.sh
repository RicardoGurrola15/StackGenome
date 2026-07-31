#!/usr/bin/env bash
set -e

# StackGenome build script for Alpha Release

PROJECT_NAME="stackgenome"
DIST_DIR="dist"

echo "Building $PROJECT_NAME for distribution..."

# Ensure we're in the right directory
cd "$(dirname "$0")/.."

# Clean and recreate dist folder
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

# Version to embed
VERSION=$(grep -E '^\s+Version\s+=' internal/cli/version.go | awk -F '"' '{print $2}')
if [ -z "$VERSION" ]; then
    VERSION="unknown"
fi
echo "Version detected: $VERSION"

LDFLAGS="-s -w"

# Build matrix
declare -a os_archs=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for os_arch in "${os_archs[@]}"; do
    IFS="/" read -r GOOS GOARCH <<< "$os_arch"
    
    OUTPUT_BIN="$DIST_DIR/${PROJECT_NAME}_${GOOS}_${GOARCH}"
    
    if [ "$GOOS" == "windows" ]; then
        OUTPUT_BIN="${OUTPUT_BIN}.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    env GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="$LDFLAGS" -o "$OUTPUT_BIN" ./cmd/stackgenome
done

echo "Generating checksums..."
cd "$DIST_DIR"
# Support both sha256sum and shasum depending on OS
if command -v sha256sum &> /dev/null; then
    sha256sum * > checksums.txt
elif command -v shasum &> /dev/null; then
    shasum -a 256 * > checksums.txt
else
    echo "Could not find sha256sum or shasum to generate checksums."
fi

echo "Build complete. Binaries available in $DIST_DIR/"
ls -lh

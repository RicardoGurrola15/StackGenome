#!/usr/bin/env bash

# This script builds the StackGenome CLI for macOS, Linux, and Windows.
# Output binaries are placed in the dist/ directory.

set -e

# Change to project root relative to script
cd "$(dirname "$0")/.."

# Activate environment if not CI
if [[ -z "$CI" ]]; then
  source scripts/activate-env.sh || true
fi

# Use git describe to get version or fallback to dev
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS="-X 'stackgenome/internal/cli.Version=$VERSION'"

echo "Building StackGenome version $VERSION..."

rm -rf dist/
mkdir -p dist/

# macOS (arm64 & amd64)
echo "Compiling for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="$LDFLAGS" -o dist/stackgenome-darwin-arm64 ./cmd/stackgenome
echo "Compiling for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -o dist/stackgenome-darwin-amd64 ./cmd/stackgenome

# Linux (arm64 & amd64)
echo "Compiling for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags="$LDFLAGS" -o dist/stackgenome-linux-arm64 ./cmd/stackgenome
echo "Compiling for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o dist/stackgenome-linux-amd64 ./cmd/stackgenome

# Windows (amd64 & arm64)
echo "Compiling for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o dist/stackgenome-windows-amd64.exe ./cmd/stackgenome
echo "Compiling for Windows (arm64)..."
GOOS=windows GOARCH=arm64 go build -ldflags="$LDFLAGS" -o dist/stackgenome-windows-arm64.exe ./cmd/stackgenome

echo "Build complete. Artifacts are in the dist/ directory."
ls -lh dist/

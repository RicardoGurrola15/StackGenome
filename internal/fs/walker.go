package fs

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"
)

// defaultExclusions contains exact directory names that are always skipped.
// These are simple name-based exclusions for common cross-ecosystem patterns.
var defaultExclusions = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	".venv":        true,
	"venv":         true,
	".next":        true,
	"dist":         true,
	".dart_tool":   true, // Dart/Flutter generated tooling
	"__pycache__":  true, // Python bytecode
	".mypy_cache":  true, // Python type checker cache
	".ruff_cache":  true, // Python linter cache
	".tox":         true, // Python test environments
	".cargo":       true, // Rust registry cache
	"_build":       true, // Elixir/Mix build output
	".turbo":       true, // Turborepo cache
	".nuxt":        true, // Nuxt.js build output
}

// vendoredPathPrefixes contains path prefixes (using forward slashes) that
// should be classified as vendored/generated and excluded from language analysis.
// These are matched against the relative path from the project root.
var vendoredPathPrefixes = []string{
	"ios/Pods/",
	"ios/Flutter/",   // Flutter-generated iOS artifacts
	"ios/.symlinks/", // Flutter plugin symlinks
	"macos/Pods/",
	"macos/Flutter/",
	"android/.gradle/",
	"android/build/",
	"android/app/build/",
	"linux/build/",
	"macos/build/",
	"windows/build/",
	"web/build/",
	"build/",       // Generic build output (Dart, Go, etc.)
	".pub-cache/",  // Pub global cache
	"third_party/", // Vendored third-party code
	".gradle/",
	"Pods/",    // CocoaPods at root level
	".kotlin/", // Kotlin incremental compilation cache
}

// SafeWalker traverses the file system starting from Root.
type SafeWalker struct {
	Root       string
	Exclusions map[string]bool
}

// NewSafeWalker initializes a new walker with default exclusions.
func NewSafeWalker(root string) *SafeWalker {
	excl := make(map[string]bool)
	for k, v := range defaultExclusions {
		excl[k] = v
	}

	return &SafeWalker{
		Root:       root,
		Exclusions: excl,
	}
}

// WalkFunc is the callback signature for files found.
type WalkFunc func(path string, d fs.DirEntry) error

// isVendoredPath returns true if the given relative path (using forward slashes)
// matches a known vendored/generated path prefix. This prevents iOS Pods,
// Gradle caches, and other generated directories from polluting language stats.
func isVendoredPath(relPath string) bool {
	normalized := filepath.ToSlash(relPath)
	for _, prefix := range vendoredPathPrefixes {
		if strings.HasPrefix(normalized, prefix) || normalized == strings.TrimSuffix(prefix, "/") {
			return true
		}
	}
	return false
}

// IsVendoredPath is the exported version for use by detectors.
func IsVendoredPath(relPath string) bool {
	return isVendoredPath(relPath)
}

// Walk executes the traversal. It catches permission errors and skips excluded directories.
// It also checks the provided context for cancellation.
func (sw *SafeWalker) Walk(ctx context.Context, fn WalkFunc) error {
	return filepath.WalkDir(sw.Root, func(path string, d fs.DirEntry, err error) error {
		// Check for context cancellation
		if ctx.Err() != nil {
			return filepath.SkipAll
		}

		// 1. Handle errors (like permission denied) gracefully
		if err != nil {
			// If it's a directory with no permissions, skip it
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			// Ignore file read errors, just continue
			return nil
		}

		// 2. Validate safety against path traversal
		if !IsSafePath(sw.Root, path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 3. Compute relative path for prefix matching
		relPath, relErr := filepath.Rel(sw.Root, path)
		if relErr != nil {
			relPath = path
		}

		// 4. Name-based exclusions (simple directory names)
		if d.IsDir() {
			if sw.Exclusions[d.Name()] {
				return filepath.SkipDir
			}
			// 5. Path-prefix exclusions for vendored/generated ecosystems
			if isVendoredPath(relPath + "/") {
				return filepath.SkipDir
			}
		}

		// 6. Pass to callback with relative path for detectors
		return fn(relPath, d)
	})
}

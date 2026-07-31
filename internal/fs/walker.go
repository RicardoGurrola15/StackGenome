package fs

import (
	"context"
	"io/fs"
	"path/filepath"
)

var defaultExclusions = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	".venv":        true,
	".next":        true,
	"dist":         true,
	"build":        true,
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

		// 3. Static Exclusions
		if d.IsDir() {
			if sw.Exclusions[d.Name()] {
				return filepath.SkipDir
			}
		}

		// 4. Pass to callback
		return fn(path, d)
	})
}

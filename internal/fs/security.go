package fs

import (
	"path/filepath"
	"strings"
)

// IsSafePath checks if the target path is inside the project root path.
// This prevents directory traversal attacks via symlinks (e.g. ../../../etc/passwd).
func IsSafePath(root, target string) bool {
	// Clean both paths
	cleanRoot := filepath.Clean(root)
	cleanTarget := filepath.Clean(target)

	// Ensure the root path isn't empty
	if cleanRoot == "" || cleanRoot == "." {
		// If root is current directory, target just shouldn't start with "../" or be absolute
		if filepath.IsAbs(cleanTarget) {
			return false
		}
		return !strings.HasPrefix(cleanTarget, ".."+string(filepath.Separator)) && cleanTarget != ".."
	}

	// Make both absolute to safely compare if root is an absolute path
	if filepath.IsAbs(cleanRoot) {
		absRoot, err := filepath.Abs(cleanRoot)
		if err != nil {
			return false
		}
		absTarget, err := filepath.Abs(cleanTarget)
		if err != nil {
			return false
		}

		// Target must start with root + separator
		return strings.HasPrefix(absTarget, absRoot+string(filepath.Separator)) || absTarget == absRoot
	}

	// For relative paths
	return strings.HasPrefix(cleanTarget, cleanRoot+string(filepath.Separator)) || cleanTarget == cleanRoot
}

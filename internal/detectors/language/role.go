package language

import (
	"path/filepath"
	"strings"

	"stackgenome/internal/projectgraph"
)

// DetermineRole infers the NodeRole based on the file path.
func DetermineRole(relPath string) projectgraph.NodeRole {
	// If it's in the root directory (or just the filename itself), it's primary
	dir := filepath.Dir(relPath)
	if dir == "." || dir == "" {
		return projectgraph.RolePrimary
	}

	// If it's inside vendor/, it's vendored
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	for _, part := range parts {
		if part == "vendor" {
			return projectgraph.RoleVendored
		}
	}

	// Otherwise, it's a satellite (e.g. in a subfolder like docs/, or a subproject)
	return projectgraph.RoleSatellite
}

package detectors

import (
// Using a sub-package structure might create import cycles if we're not careful.
// We'll let the language package register itself or we just import it here.
// Actually, wait, if `detectors` imports `language` and `language` imports `detectors`, that's an import cycle.
// Yes, `language` imports `stackgenome/internal/detectors` for the interface.
// So `detectors` cannot import `language`.
// We need a subpackage like `internal/detectors/registry` or we register them dynamically.
)

var defaultDetectors []FileDetector

// Register adds a detector to the default registry.
func Register(d FileDetector) {
	defaultDetectors = append(defaultDetectors, d)
}

// DefaultRegistry returns the list of registered detectors.
func DefaultRegistry() []FileDetector {
	// Return a copy to avoid mutation
	out := make([]FileDetector, len(defaultDetectors))
	copy(out, defaultDetectors)
	return out
}

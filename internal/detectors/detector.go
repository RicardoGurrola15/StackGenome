package detectors

import "stackgenome/internal/projectgraph"

// FileDetector defines an interface for components that can extract knowledge from a file.
type FileDetector interface {
	// Handles returns true if the detector should process the given filename (e.g. "package.json").
	Handles(filename string) bool

	// Detect processes the file content and returns the detected nodes and relationships.
	// The path should be a relative path from the project root.
	Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error)
}

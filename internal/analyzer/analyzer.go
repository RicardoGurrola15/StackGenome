package analyzer

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"stackgenome/internal/detectors"
	localfs "stackgenome/internal/fs"
	"stackgenome/internal/projectgraph"
)

// Analyzer orchestrates the file system walking and applying detectors.
type Analyzer struct {
	RootPath      string
	Walker        *localfs.SafeWalker
	FileDetectors []detectors.FileDetector
}

// NewAnalyzer creates a new configured analyzer for the given root path.
func NewAnalyzer(root string, detectors []detectors.FileDetector) *Analyzer {
	return &Analyzer{
		RootPath:      root,
		Walker:        localfs.NewSafeWalker(root),
		FileDetectors: detectors,
	}
}

// Analyze runs the complete analysis pipeline and returns the ProjectGraph.
func (a *Analyzer) Analyze() (*projectgraph.Graph, error) {
	graph := projectgraph.NewGraph()

	ctx := context.Background()
	// The walker now passes relative paths in the callback.
	err := a.Walker.Walk(ctx, func(relPath string, d fs.DirEntry) error {
		// Only process regular files
		if d == nil || d.IsDir() {
			return nil
		}

		filename := d.Name()

		// Filter detectors that handle this file
		var matchingDetectors []detectors.FileDetector
		for _, det := range a.FileDetectors {
			if det.Handles(filename) {
				matchingDetectors = append(matchingDetectors, det)
			}
		}

		// Skip reading file if no detector handles it
		if len(matchingDetectors) == 0 {
			return nil
		}

		// Reconstruct absolute path for file reading
		absPath := filepath.Join(a.RootPath, relPath)

		// Read file content once for all interested detectors
		content, err := os.ReadFile(absPath)
		if err != nil {
			// Skip files we can't read
			return nil
		}

		for _, det := range matchingDetectors {
			nodes, edges, err := det.Detect(relPath, content)
			if err != nil {
				// We don't abort analysis for one bad parse
				continue
			}

			for _, n := range nodes {
				// AddNode merges evidences if the node already exists (deduplication)
				_ = graph.AddNode(n)
			}
			for _, e := range edges {
				_ = graph.AddEdge(e.SourceID, e.TargetID, e.Relation)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("analysis walk failed: %w", err)
	}

	// Post-processing: deduplicate platform nodes by canonical key
	graph.DeduplicatePlatforms()

	return graph, nil
}

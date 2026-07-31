package platform

import (
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

type PlatformDetector struct{}

var _ detectors.FileDetector = (*PlatformDetector)(nil)

func (d *PlatformDetector) Handles(filename string) bool {
	return filename == "AndroidManifest.xml" || filename == "Info.plist"
}

func (d *PlatformDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = strings.ReplaceAll(filepath.Dir(relPath), "/", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, "\\", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, ".", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, "-", "_")
	}

	nodeID := "platform_" + sanitizedPath + "_" + strings.ReplaceAll(filename, ".", "_")

	name := "Android Native"
	if filename == "Info.plist" {
		name = "iOS/macOS Native"
	}

	node := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypePlatform,
		Role:       projectgraph.RolePrimary,
		Name:       name,
		Confidence: 1.0,
		Evidences: []evidence.Evidence{
			{
				Kind:        "manifest",
				Path:        relPath,
				Sensitivity: evidence.LocalPath,
			},
		},
		Properties: make(map[string]string),
	}

	return []*projectgraph.Node{node}, nil, nil
}

package editor

import (
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

type EditorDetector struct{}

var _ detectors.FileDetector = (*EditorDetector)(nil)

func (d *EditorDetector) Handles(filename string) bool {
	return filename == "settings.json" || filename == "launch.json" || filename == "workspace.xml" || filename == "devcontainer.json" || strings.HasSuffix(filename, ".iml")
}

func (d *EditorDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	slashPath := filepath.ToSlash(relPath)
	isVSCode := strings.Contains(slashPath, ".vscode/")
	isIntelliJ := strings.Contains(slashPath, ".idea/") || strings.HasSuffix(filename, ".iml")
	isDevContainer := strings.Contains(slashPath, ".devcontainer/")

	if !isVSCode && !isIntelliJ && !isDevContainer {
		return nil, nil, nil
	}

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = strings.ReplaceAll(filepath.Dir(relPath), "/", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, "\\", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, ".", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, "-", "_")
	}

	nodeID := "editor_" + sanitizedPath + "_" + strings.ReplaceAll(filename, ".", "_")

	name := "Unknown Editor"
	if isVSCode {
		name = "VS Code"
	} else if isIntelliJ {
		name = "IntelliJ/JetBrains"
	} else if isDevContainer {
		name = "DevContainer"
	}

	node := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeEditor,
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

package cicd

import (
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

type CICDDetector struct{}

var _ detectors.FileDetector = (*CICDDetector)(nil)

func (d *CICDDetector) Handles(filename string) bool {
	return filename == "Jenkinsfile" || filename == ".gitlab-ci.yml" || strings.HasSuffix(filename, ".yml") || strings.HasSuffix(filename, ".yaml")
}

func (d *CICDDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	isGitHubAction := strings.Contains(filepath.ToSlash(relPath), ".github/workflows/")
	isGitLab := filename == ".gitlab-ci.yml"
	isJenkins := filename == "Jenkinsfile"

	if !isGitHubAction && !isGitLab && !isJenkins {
		return nil, nil, nil
	}

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = strings.ReplaceAll(filepath.Dir(relPath), "/", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, "\\", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, ".", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, "-", "_")
	}

	nodeID := "cicd_" + sanitizedPath + "_" + strings.ReplaceAll(filename, ".", "_")

	name := "Unknown CI/CD"
	if isGitHubAction {
		name = "GitHub Actions"
	} else if isGitLab {
		name = "GitLab CI"
	} else if isJenkins {
		name = "Jenkins"
	}

	node := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeCICD,
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

package infra

import (
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

type InfraDetector struct{}

var _ detectors.FileDetector = (*InfraDetector)(nil)

func (d *InfraDetector) Handles(filename string) bool {
	return filename == "Dockerfile" || filename == "docker-compose.yml" || filename == "docker-compose.yaml" || strings.HasSuffix(filename, ".tf")
}

func (d *InfraDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = strings.ReplaceAll(filepath.Dir(relPath), "/", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, "\\", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, ".", "_")
		sanitizedPath = strings.ReplaceAll(sanitizedPath, "-", "_")
	}

	nodeID := "infra_" + sanitizedPath + "_" + strings.ReplaceAll(filename, ".", "_")

	name := "Docker"
	if strings.HasSuffix(filename, ".tf") {
		name = "Terraform"
	}

	infraNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeInfrastructure,
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

	if filename == "docker-compose.yml" || filename == "docker-compose.yaml" {
		infraNode.Properties["type"] = "compose"
	} else if filename == "Dockerfile" {
		infraNode.Properties["type"] = "image"
	} else if strings.HasSuffix(filename, ".tf") {
		infraNode.Properties["type"] = "iac"
	}

	return []*projectgraph.Node{infraNode}, nil, nil
}

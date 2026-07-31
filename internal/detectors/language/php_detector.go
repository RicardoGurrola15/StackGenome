package language

import (
	"encoding/json"
	"path/filepath"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

type PHPDetector struct{}

var _ detectors.FileDetector = (*PHPDetector)(nil)

func (d *PHPDetector) Handles(filename string) bool {
	return filename == "composer.json" || filename == "composer.lock"
}

func (d *PHPDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	if filename == "composer.lock" {
		return d.detectLockfile(relPath)
	}

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}
	nodeID := "lang_php_" + sanitizedPath

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "PHP",
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

	var composer struct {
		Require    map[string]string `json:"require"`
		RequireDev map[string]string `json:"require-dev"`
	}
	if err := json.Unmarshal(content, &composer); err != nil {
		return nil, nil, err
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	addDeps := func(deps map[string]string) {
		for name, version := range deps {
			if name == "php" || name == "php-64bit" || name == "ext-json" {
				continue
			}
			depID := "dep_php_" + sanitizeID(name)
			depNode := &projectgraph.Node{
				ID:         depID,
				Type:       projectgraph.TypeDependency,
				Name:       name,
				Version:    version,
				Confidence: 0.8,
				Properties: map[string]string{
					"purl": "pkg:composer/" + name + "@" + version,
				},
			}
			nodes = append(nodes, depNode)
			edges = append(edges, projectgraph.Edge{
				SourceID: langNode.ID,
				TargetID: depID,
				Relation: projectgraph.RelationDependsOn,
			})
		}
	}

	addDeps(composer.Require)
	addDeps(composer.RequireDev)

	return nodes, edges, nil
}

func (d *PHPDetector) detectLockfile(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         "lang_php_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "PHP",
		Confidence: 1.0,
		Evidences: []evidence.Evidence{
			{
				Kind:        "lockfile",
				Path:        relPath,
				Sensitivity: evidence.LocalPath,
			},
		},
	}

	return []*projectgraph.Node{langNode}, nil, nil
}

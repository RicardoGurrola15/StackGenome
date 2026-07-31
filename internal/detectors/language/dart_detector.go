package language

import (
	"bufio"
	"bytes"
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

type DartDetector struct{}

var _ detectors.FileDetector = (*DartDetector)(nil)

func (d *DartDetector) Handles(filename string) bool {
	return filename == "pubspec.yaml" || filename == "pubspec.lock"
}

func (d *DartDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	if filename == "pubspec.lock" {
		return d.detectLockfile(relPath)
	}

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}
	nodeID := "lang_dart_" + sanitizedPath

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Dart/Flutter",
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

	var deps []struct{ name, version string }

	scanner := bufio.NewScanner(bytes.NewReader(content))
	inDependencies := false
	inDevDependencies := false
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "dependencies:" {
			inDependencies = true
			inDevDependencies = false
			continue
		} else if trimmed == "dev_dependencies:" {
			inDependencies = false
			inDevDependencies = true
			continue
		} else if !strings.HasPrefix(line, " ") && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			inDependencies = false
			inDevDependencies = false
		}

		if (inDependencies || inDevDependencies) && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) > 0 {
				name := strings.TrimSpace(parts[0])
				version := ""
				if len(parts) > 1 {
					version = strings.TrimSpace(parts[1])
					version = strings.Trim(version, "^'\"")
				}
				if name != "" && name != "flutter" && name != "sdk" {
					deps = append(deps, struct{ name, version string }{name: name, version: version})
				}
			}
		}
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	for _, dep := range deps {
		depID := "dep_dart_" + sanitizeID(dep.name)
		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       dep.name,
			Version:    dep.version,
			Confidence: 0.8,
			Properties: map[string]string{
				"purl": "pkg:pub/" + dep.name + "@" + dep.version,
			},
		}
		nodes = append(nodes, depNode)
		edges = append(edges, projectgraph.Edge{
			SourceID: langNode.ID,
			TargetID: depID,
			Relation: projectgraph.RelationDependsOn,
		})
	}

	return nodes, edges, nil
}

func (d *DartDetector) detectLockfile(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         "lang_dart_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Dart/Flutter",
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

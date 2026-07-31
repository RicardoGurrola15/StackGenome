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

type RubyDetector struct{}

var _ detectors.FileDetector = (*RubyDetector)(nil)

func (d *RubyDetector) Handles(filename string) bool {
	return filename == "Gemfile" || filename == "Gemfile.lock"
}

func (d *RubyDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	if filename == "Gemfile.lock" {
		return d.detectLockfile(relPath)
	}

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}
	nodeID := "lang_ruby_" + sanitizedPath

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Ruby",
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
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "gem ") {
			parts := strings.SplitN(line, ",", 2)
			namePart := strings.TrimPrefix(parts[0], "gem ")
			namePart = strings.TrimSpace(namePart)
			namePart = strings.Trim(namePart, "'\"")

			version := ""
			if len(parts) > 1 {
				version = strings.TrimSpace(parts[1])
				version = strings.Trim(version, "'\"")
			}

			deps = append(deps, struct{ name, version string }{name: namePart, version: version})
		}
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	for _, dep := range deps {
		if dep.name == "" {
			continue
		}
		depID := "dep_ruby_" + sanitizeID(dep.name)
		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       dep.name,
			Version:    dep.version,
			Confidence: 0.8,
			Properties: map[string]string{
				"purl": "pkg:gem/" + dep.name + "@" + dep.version,
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

func (d *RubyDetector) detectLockfile(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         "lang_ruby_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Ruby",
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

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

type SwiftDetector struct{}

var _ detectors.FileDetector = (*SwiftDetector)(nil)

func (d *SwiftDetector) Handles(filename string) bool {
	return filename == "Package.swift" || filename == "Podfile" || filename == "Package.resolved" || filename == "Podfile.lock"
}

func (d *SwiftDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	if filename == "Package.resolved" || filename == "Podfile.lock" {
		return d.detectLockfile(relPath)
	}

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}
	nodeID := "lang_swift_" + sanitizedPath

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Swift/Objective-C",
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

	if filename == "Package.swift" {
		langNode.Properties["manager"] = "spm"
		deps = d.parseSPM(content)
	} else if filename == "Podfile" {
		langNode.Properties["manager"] = "cocoapods"
		deps = d.parsePodfile(content)
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	for _, dep := range deps {
		if dep.name == "" {
			continue
		}
		depID := "dep_swift_" + sanitizeID(dep.name)
		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       dep.name,
			Version:    dep.version,
			Confidence: 0.8,
			Properties: map[string]string{
				"purl": "pkg:swift/" + dep.name,
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

func (d *SwiftDetector) detectLockfile(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         "lang_swift_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Swift/Objective-C",
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

func (d *SwiftDetector) parseSPM(content []byte) []struct{ name, version string } {
	var deps []struct{ name, version string }
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Basic best effort parsing for .package(url: "https://github.com/Alamofire/Alamofire.git", from: "5.0.0")
		if strings.HasPrefix(line, ".package(") {
			urlPart := extractBetween(line, "url:", ",")
			if urlPart == "" {
				urlPart = extractBetween(line, "url:", ")")
			}
			urlPart = strings.TrimSpace(urlPart)
			urlPart = strings.Trim(urlPart, "\"")

			if urlPart != "" {
				name := filepath.Base(urlPart)
				name = strings.TrimSuffix(name, ".git")

				version := extractBetween(line, "from:", ")")
				version = strings.Trim(strings.TrimSpace(version), "\"")

				deps = append(deps, struct{ name, version string }{name: name, version: version})
			}
		}
	}
	return deps
}

func (d *SwiftDetector) parsePodfile(content []byte) []struct{ name, version string } {
	var deps []struct{ name, version string }
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "pod ") {
			parts := strings.SplitN(line, ",", 2)
			namePart := strings.TrimPrefix(parts[0], "pod ")
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
	return deps
}

func extractBetween(s, start, end string) string {
	startIdx := strings.Index(s, start)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(start)

	endIdx := strings.Index(s[startIdx:], end)
	if endIdx == -1 {
		return ""
	}

	return s[startIdx : startIdx+endIdx]
}

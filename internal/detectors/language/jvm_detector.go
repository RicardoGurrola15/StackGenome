package language

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

type JVMDetector struct{}

var _ detectors.FileDetector = (*JVMDetector)(nil)

func (d *JVMDetector) Handles(filename string) bool {
	return filename == "pom.xml" || filename == "build.gradle" || filename == "build.gradle.kts"
}

func (d *JVMDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}
	nodeID := "lang_jvm_" + sanitizedPath

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "JVM",
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

	var deps []struct{ group, name, version string }

	if filename == "pom.xml" {
		langNode.Properties["manager"] = "maven"
		deps = d.parsePom(content)
	} else if strings.HasPrefix(filename, "build.gradle") {
		langNode.Properties["manager"] = "gradle"
		deps = d.parseGradle(content)
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	for _, dep := range deps {
		if dep.group == "" || dep.name == "" {
			continue
		}
		depID := "dep_jvm_" + sanitizeID(dep.group+"_"+dep.name)
		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       dep.group + ":" + dep.name,
			Version:    dep.version,
			Confidence: 0.8,
			Properties: map[string]string{
				"purl": "pkg:maven/" + dep.group + "/" + dep.name + "@" + dep.version,
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

func (d *JVMDetector) parsePom(content []byte) []struct{ group, name, version string } {
	type Dependency struct {
		GroupId    string `xml:"groupId"`
		ArtifactId string `xml:"artifactId"`
		Version    string `xml:"version"`
	}
	type Project struct {
		Dependencies []Dependency `xml:"dependencies>dependency"`
	}
	var proj Project
	_ = xml.Unmarshal(content, &proj)

	var deps []struct{ group, name, version string }
	for _, dep := range proj.Dependencies {
		deps = append(deps, struct{ group, name, version string }{
			group:   dep.GroupId,
			name:    dep.ArtifactId,
			version: dep.Version,
		})
	}
	return deps
}

func (d *JVMDetector) parseGradle(content []byte) []struct{ group, name, version string } {
	var deps []struct{ group, name, version string }
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "implementation ") || strings.HasPrefix(line, "api ") || strings.HasPrefix(line, "compile ") || strings.HasPrefix(line, "implementation(") || strings.HasPrefix(line, "api(") {
			// format: implementation 'group:name:version'
			parts := strings.Split(line, "'")
			if len(parts) < 3 {
				parts = strings.Split(line, "\"")
			}
			if len(parts) >= 3 {
				coords := strings.Split(parts[1], ":")
				if len(coords) >= 2 {
					version := ""
					if len(coords) >= 3 {
						version = coords[2]
					}
					deps = append(deps, struct{ group, name, version string }{group: coords[0], name: coords[1], version: version})
				}
			}
		}
	}
	return deps
}

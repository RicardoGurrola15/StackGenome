package language

import (
	"bufio"
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

type CppDetector struct{}

var _ detectors.FileDetector = (*CppDetector)(nil)

func (d *CppDetector) Handles(filename string) bool {
	return filename == "CMakeLists.txt" || filename == "conanfile.txt" || filename == "conanfile.py" || filename == "vcpkg.json"
}

func (d *CppDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}
	nodeID := "lang_cpp_" + sanitizedPath

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "C/C++",
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

	if filename == "CMakeLists.txt" {
		langNode.Properties["manager"] = "cmake"
		deps = d.parseCMake(content)
	} else if filename == "conanfile.txt" {
		langNode.Properties["manager"] = "conan"
		deps = d.parseConanTxt(content)
	} else if filename == "conanfile.py" {
		langNode.Properties["manager"] = "conan"
	} else if filename == "vcpkg.json" {
		langNode.Properties["manager"] = "vcpkg"
		deps = d.parseVcpkg(content)
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	for _, dep := range deps {
		if dep.name == "" {
			continue
		}
		depID := "dep_cpp_" + sanitizeID(dep.name)
		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       dep.name,
			Version:    dep.version,
			Confidence: 0.8,
			Properties: map[string]string{
				"purl": "pkg:cpp/" + dep.name,
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

func (d *CppDetector) parseCMake(content []byte) []struct{ name, version string } {
	var deps []struct{ name, version string }
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "find_package(") {
			part := strings.TrimPrefix(line, "find_package(")
			part = strings.TrimSuffix(part, ")")
			parts := strings.Fields(part)
			if len(parts) > 0 {
				deps = append(deps, struct{ name, version string }{name: parts[0], version: ""})
			}
		}
	}
	return deps
}

func (d *CppDetector) parseConanTxt(content []byte) []struct{ name, version string } {
	var deps []struct{ name, version string }
	scanner := bufio.NewScanner(bytes.NewReader(content))
	inRequires := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			inRequires = (line == "[requires]")
			continue
		}
		if inRequires && line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, "/")
			if len(parts) >= 2 {
				deps = append(deps, struct{ name, version string }{name: parts[0], version: parts[1]})
			}
		}
	}
	return deps
}

func (d *CppDetector) parseVcpkg(content []byte) []struct{ name, version string } {
	var deps []struct{ name, version string }
	var vcpkg struct {
		Dependencies []interface{} `json:"dependencies"`
	}
	_ = json.Unmarshal(content, &vcpkg)

	for _, dep := range vcpkg.Dependencies {
		if strDep, ok := dep.(string); ok {
			deps = append(deps, struct{ name, version string }{name: strDep, version: ""})
		} else if mapDep, ok := dep.(map[string]interface{}); ok {
			if name, ok := mapDep["name"].(string); ok {
				deps = append(deps, struct{ name, version string }{name: name, version: ""})
			}
		}
	}
	return deps
}

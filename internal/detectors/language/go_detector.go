package language

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

// GoDetector parses go.mod and go.work files to detect Go language usage
// and its direct dependencies.
type GoDetector struct{}

// Ensure it implements the interface
var _ detectors.FileDetector = (*GoDetector)(nil)

func (d *GoDetector) Handles(filename string) bool {
	return filename == "go.mod" || filename == "go.work" || filename == "go.sum"
}

func (d *GoDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	switch filename {
	case "go.mod":
		return d.detectGoMod(relPath, content)
	case "go.work":
		return d.detectGoWork(relPath, content)
	default:
		return nil, nil, nil
	}
}

func (d *GoDetector) detectGoMod(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))

	var goVersion string
	var modulePath string
	var deps []struct{ path, version string }
	inRequireBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Strip comments
		if ci := strings.Index(line, "//"); ci >= 0 {
			line = strings.TrimSpace(line[:ci])
		}

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				modulePath = parts[1]
			}
			continue
		}

		if strings.HasPrefix(line, "go ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				goVersion = parts[1]
			}
			continue
		}

		// Single require
		if strings.HasPrefix(line, "require ") && !strings.HasSuffix(line, "(") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				deps = append(deps, struct{ path, version string }{parts[1], parts[2]})
			}
			continue
		}

		// require block
		if line == "require (" || strings.HasSuffix(line, "require (") {
			inRequireBlock = true
			continue
		}
		if inRequireBlock {
			if line == ")" {
				inRequireBlock = false
				continue
			}
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps = append(deps, struct{ path, version string }{parts[0], parts[1]})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	sanitizedPath := "root"
	if relPath != "go.mod" {
		sanitizedPath = sanitizeID(relPath)
	}

	langNode := &projectgraph.Node{
		ID:         "lang_go_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Go",
		Version:    goVersion,
		Confidence: 1.0,
		Properties: map[string]string{
			"module": modulePath,
		},
		Evidences: []evidence.Evidence{
			{
				Kind:        "manifest",
				Path:        relPath,
				Selector:    "go directive",
				Value:       goVersion,
				Sensitivity: evidence.PublicMetadata,
			},
		},
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	// Generate dependency nodes + edges
	for _, dep := range deps {
		depSuffix := sanitizeID(relPath)
		if relPath == "go.mod" {
			depSuffix = "root"
		}
		depID := fmt.Sprintf("dep_go_%s_%s", sanitizeID(dep.path), depSuffix)
		purl := fmt.Sprintf("pkg:golang/%s@%s", dep.path, dep.version)
		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       dep.path,
			Version:    dep.version,
			Confidence: 1.0,
			Properties: map[string]string{
				"purl": purl,
			},
			Evidences: []evidence.Evidence{
				{
					Kind:        "manifest",
					Path:        relPath,
					Selector:    "require",
					Value:       dep.path + " " + dep.version,
					Sensitivity: evidence.PublicMetadata,
				},
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

func (d *GoDetector) detectGoWork(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	var members []string
	inUseBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if ci := strings.Index(line, "//"); ci >= 0 {
			line = strings.TrimSpace(line[:ci])
		}
		if line == "" {
			continue
		}

		if line == "use (" || strings.HasSuffix(line, "use (") {
			inUseBlock = true
			continue
		}
		if strings.HasPrefix(line, "use ") && !strings.HasSuffix(line, "(") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				members = append(members, parts[1])
			}
			continue
		}
		if inUseBlock {
			if line == ")" {
				inUseBlock = false
				continue
			}
			members = append(members, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	workNodeID := "workspace_go_root"
	workNode := &projectgraph.Node{
		ID:         workNodeID,
		Type:       projectgraph.TypeWorkspace,
		Name:       "Go Workspace",
		Confidence: 1.0,
		Properties: map[string]string{
			"members": strings.Join(members, ","),
		},
		Evidences: []evidence.Evidence{
			{
				Kind:        "manifest",
				Path:        relPath,
				Selector:    "use",
				Value:       strings.Join(members, ", "),
				Sensitivity: evidence.PublicMetadata,
			},
		},
	}

	return []*projectgraph.Node{workNode}, nil, nil
}

func (d *GoDetector) detectGoSum(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if relPath != "go.sum" {
		sanitizedPath = sanitizeID(relPath)
	}

	langNode := &projectgraph.Node{
		ID:         "lang_go_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Go",
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

// sanitizeID replaces path separators and special chars with underscores.
func sanitizeID(s string) string {
	r := strings.NewReplacer("/", "_", ".", "_", "@", "_", "-", "_")
	return r.Replace(s)
}

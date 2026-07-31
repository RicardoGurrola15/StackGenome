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

// RustDetector parses Cargo.toml files to detect Rust language usage,
// its direct dependencies, and workspace members.
type RustDetector struct{}

// Ensure it implements the interface
var _ detectors.FileDetector = (*RustDetector)(nil)

func (d *RustDetector) Handles(filename string) bool {
	return filename == "Cargo.toml" || filename == "Cargo.lock"
}

func (d *RustDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	if filename == "Cargo.lock" {
		return d.detectLockfile(relPath)
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))

	type section int
	const (
		sectionNone section = iota
		sectionDeps
		sectionWorkspace
		sectionPackage
	)

	current := sectionNone
	var crateName, crateVersion string
	var deps []struct{ name, version string }
	var workspaceMembers []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Strip comments
		if ci := strings.Index(line, "#"); ci >= 0 {
			line = strings.TrimSpace(line[:ci])
		}
		if line == "" {
			continue
		}

		// Section headers
		lower := strings.ToLower(line)
		switch {
		case lower == "[dependencies]" || lower == "[dev-dependencies]" || lower == "[build-dependencies]":
			current = sectionDeps
			continue
		case lower == "[workspace]":
			current = sectionWorkspace
			continue
		case lower == "[package]":
			current = sectionPackage
			continue
		case strings.HasPrefix(lower, "["):
			current = sectionNone
			continue
		}

		eqIdx := strings.Index(line, "=")
		if eqIdx < 0 {
			continue
		}

		key := strings.TrimSpace(line[:eqIdx])
		val := strings.TrimSpace(line[eqIdx+1:])

		switch current {
		case sectionPackage:
			switch key {
			case "name":
				crateName = strings.Trim(val, `"`)
			case "version":
				crateVersion = strings.Trim(val, `"`)
			}

		case sectionDeps:
			depName := key
			depVersion := ""
			// value can be "1.0.0" (string) or { version = "1.0.0", ... }
			if strings.HasPrefix(val, "{") {
				if vi := strings.Index(val, "version"); vi >= 0 {
					sub := val[vi+len("version"):]
					if qi := strings.Index(sub, "\""); qi >= 0 {
						sub = sub[qi+1:]
						if qi2 := strings.Index(sub, "\""); qi2 >= 0 {
							depVersion = sub[:qi2]
						}
					}
				}
			} else {
				depVersion = strings.Trim(val, `"`)
			}
			if depName != "" {
				deps = append(deps, struct{ name, version string }{depName, depVersion})
			}

		case sectionWorkspace:
			// members = ["crate_a", "crate_b"]
			if key == "members" {
				memberLine := strings.Trim(val, "[]")
				for _, m := range strings.Split(memberLine, ",") {
					m = strings.Trim(strings.TrimSpace(m), `"`)
					if m != "" {
						workspaceMembers = append(workspaceMembers, m)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	nodeID := "lang_rust_root"
	if relPath != "Cargo.toml" {
		nodeID = fmt.Sprintf("lang_rust_%s", sanitizeID(relPath))
	}

	props := map[string]string{}
	if crateName != "" {
		props["crate"] = crateName
	}
	if len(workspaceMembers) > 0 {
		props["workspace"] = "true"
		props["members"] = strings.Join(workspaceMembers, ",")
	}

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Rust",
		Version:    crateVersion,
		Confidence: 1.0,
		Properties: props,
		Evidences: []evidence.Evidence{
			{
				Kind:        "manifest",
				Path:        relPath,
				Selector:    "[package]",
				Value:       crateName + "@" + crateVersion,
				Sensitivity: evidence.PublicMetadata,
			},
		},
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	for _, dep := range deps {
		depID := fmt.Sprintf("dep_cargo_%s_%s", sanitizeID(dep.name), sanitizeID(relPath))
		if relPath == "Cargo.toml" {
			depID = fmt.Sprintf("dep_cargo_%s_root", sanitizeID(dep.name))
		}

		purl := fmt.Sprintf("pkg:cargo/%s", dep.name)
		if dep.version != "" {
			purl += "@" + dep.version
		}

		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       dep.name,
			Version:    dep.version,
			Confidence: 1.0,
			Properties: map[string]string{
				"purl":    purl,
				"package": dep.name,
			},
			Evidences: []evidence.Evidence{
				{
					Kind:        "manifest",
					Path:        relPath,
					Selector:    "[dependencies]." + dep.name,
					Value:       dep.name + "@" + dep.version,
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

func (d *RustDetector) detectLockfile(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         "lang_rust_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Rust",
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

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

// DartDetector parses pubspec.yaml and pubspec.lock files to detect
// Flutter/Dart projects, their dependencies (with runtime/dev scope),
// and version constraints.
type DartDetector struct{}

var _ detectors.FileDetector = (*DartDetector)(nil)

func (d *DartDetector) Handles(filename string) bool {
	return filename == "pubspec.yaml" ||
		filename == "pubspec.lock" ||
		filename == "shorebird.yaml" ||
		filename == "firebase.json"
}

func (d *DartDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	switch filename {
	case "pubspec.lock":
		return d.detectLockfile(relPath, content)
	case "shorebird.yaml":
		return d.detectShorebird(relPath)
	case "firebase.json":
		return d.detectFirebase(relPath, content)
	default:
		return d.detectPubspec(relPath, content)
	}
}

// detectPubspec parses pubspec.yaml for runtime and dev dependencies.
func (d *DartDetector) detectPubspec(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
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

	type depEntry struct {
		name    string
		version string
		scope   projectgraph.NodeScope
	}
	var deps []depEntry

	scanner := bufio.NewScanner(bytes.NewReader(content))

	// State machine for section tracking
	const (
		sectionNone = iota
		sectionDeps
		sectionDevDeps
		sectionEnvironment
	)
	section := sectionNone
	currentDepIndent := 0

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Detect section headers (top-level keys)
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			key := strings.TrimSuffix(trimmed, ":")
			switch key {
			case "dependencies":
				section = sectionDeps
				currentDepIndent = 2
			case "dev_dependencies":
				section = sectionDevDeps
				currentDepIndent = 2
			case "environment":
				section = sectionEnvironment
			default:
				section = sectionNone
			}
			continue
		}

		// Parse environment SDK constraint
		if section == sectionEnvironment {
			if strings.Contains(trimmed, "sdk:") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					langNode.Properties["sdk_constraint"] = strings.TrimSpace(strings.Trim(parts[1], "\"'"))
				}
			}
			continue
		}

		if section == sectionNone {
			continue
		}

		// Parse dependency entry (must be indented)
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent < currentDepIndent || trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Only process direct dep keys (indent == currentDepIndent)
		if indent != currentDepIndent {
			continue
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) == 0 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		if name == "" || name == "flutter" || name == "sdk" {
			continue
		}

		// Extract inline version constraint if present
		version := ""
		if len(parts) > 1 {
			raw := strings.TrimSpace(parts[1])
			if raw != "" && raw != "{}" && !strings.HasPrefix(raw, "{") {
				version = strings.Trim(raw, "^\"'~>=<")
			}
		}

		scope := projectgraph.ScopeRuntime
		if section == sectionDevDeps {
			scope = projectgraph.ScopeDevelopment
		}

		deps = append(deps, depEntry{name: name, version: version, scope: scope})
	}

	langNode.Properties["manager"] = "pub"
	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	for _, dep := range deps {
		depID := "dep_dart_" + sanitizeID(dep.name)
		purl := "pkg:pub/" + dep.name
		if dep.version != "" {
			purl += "@" + dep.version
		}
		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Role:       DetermineRole(relPath),
			Scope:      dep.scope,
			Name:       dep.name,
			Version:    dep.version, // declared constraint
			Confidence: 0.85,
			Properties: map[string]string{
				"purl":   purl,
				"source": "pubspec.yaml",
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

// detectLockfile parses pubspec.lock to extract resolved versions for known deps.
// It enriches existing dependency nodes with their resolved version.
func (d *DartDetector) detectLockfile(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}

	// Also emit a language node to ensure we detect Dart even without pubspec.yaml
	langNode := &projectgraph.Node{
		ID:         "lang_dart_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Dart/Flutter",
		Confidence: 0.95, // Slightly lower than pubspec.yaml
		Evidences: []evidence.Evidence{
			{
				Kind:        "lockfile",
				Path:        relPath,
				Sensitivity: evidence.LocalPath,
			},
		},
		Properties: map[string]string{"manager": "pub"},
	}

	// Parse package sections to extract resolved versions
	// pubspec.lock format:
	// packages:
	//   package_name:
	//     version: "1.2.3"
	//     source: hosted
	var nodes []*projectgraph.Node
	nodes = append(nodes, langNode)

	scanner := bufio.NewScanner(bytes.NewReader(content))
	inPackages := false
	currentPkg := ""
	var resolvedVersion, sourceType string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "packages:" {
			inPackages = true
			continue
		}
		if !inPackages {
			continue
		}

		// Top-level package name (2-space indent)
		if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "   ") && strings.HasSuffix(trimmed, ":") {
			// Save previous package if any
			if currentPkg != "" && resolvedVersion != "" {
				depID := "dep_dart_" + sanitizeID(currentPkg)
				// Emit a thin node with resolved version that AddNode will merge
				purl := "pkg:pub/" + currentPkg + "@" + resolvedVersion
				n := &projectgraph.Node{
					ID:         depID,
					Type:       projectgraph.TypeDependency,
					Name:       currentPkg,
					Resolved:   resolvedVersion,
					Confidence: 0.95,
					Properties: map[string]string{
						"purl":        purl,
						"source_type": sourceType,
						"source":      "pubspec.lock",
					},
				}
				nodes = append(nodes, n)
			}
			currentPkg = strings.TrimSuffix(trimmed, ":")
			resolvedVersion = ""
			sourceType = ""
			continue
		}

		// Parse attributes of current package
		if currentPkg != "" {
			if strings.Contains(trimmed, "version:") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					resolvedVersion = strings.Trim(strings.TrimSpace(parts[1]), "\"'")
				}
			} else if trimmed == "source: hosted" {
				sourceType = "hosted"
			} else if trimmed == "source: path" {
				sourceType = "path"
			} else if trimmed == "source: git" {
				sourceType = "git"
			} else if trimmed == "source: sdk" {
				sourceType = "sdk"
			}
		}
	}
	// Save last package
	if currentPkg != "" && resolvedVersion != "" {
		depID := "dep_dart_" + sanitizeID(currentPkg)
		purl := "pkg:pub/" + currentPkg + "@" + resolvedVersion
		n := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       currentPkg,
			Resolved:   resolvedVersion,
			Confidence: 0.95,
			Properties: map[string]string{
				"purl":        purl,
				"source_type": sourceType,
				"source":      "pubspec.lock",
			},
		}
		nodes = append(nodes, n)
	}

	return nodes, nil, nil
}

// detectShorebird detects Shorebird OTA distribution from shorebird.yaml.
// It does NOT export the app_id (privacy).
func (d *DartDetector) detectShorebird(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	node := &projectgraph.Node{
		ID:         "tool_shorebird",
		Type:       projectgraph.TypeTool,
		Role:       projectgraph.RolePrimary,
		Name:       "Shorebird",
		Confidence: 1.0,
		Evidences: []evidence.Evidence{
			{
				Kind:        "config",
				Path:        relPath,
				Sensitivity: evidence.PublicMetadata, // app_id is NOT exported
			},
		},
		Properties: map[string]string{
			"category": "distribution",
			"type":     "ota-updates",
		},
	}
	return []*projectgraph.Node{node}, nil, nil
}

// detectFirebase detects Firebase services from firebase.json.
// It does NOT export project IDs, numeric IDs, or endpoints.
func (d *DartDetector) detectFirebase(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	// Detect which Firebase services are declared without reading project IDs
	services := []string{}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, `"firestore"`):
			services = append(services, "firestore")
		case strings.HasPrefix(trimmed, `"hosting"`):
			services = append(services, "hosting")
		case strings.HasPrefix(trimmed, `"storage"`):
			services = append(services, "storage")
		case strings.HasPrefix(trimmed, `"functions"`):
			services = append(services, "functions")
		case strings.HasPrefix(trimmed, `"emulators"`):
			services = append(services, "emulators")
		case strings.HasPrefix(trimmed, `"remoteconfig"`):
			services = append(services, "remoteconfig")
		}
	}

	node := &projectgraph.Node{
		ID:         "infra_firebase",
		Type:       projectgraph.TypeInfrastructure,
		Role:       projectgraph.RolePrimary,
		Name:       "Firebase",
		Confidence: 1.0,
		Evidences: []evidence.Evidence{
			{
				Kind:        "config",
				Path:        relPath,
				Sensitivity: evidence.PublicMetadata, // no project IDs exported
			},
		},
		Properties: map[string]string{
			"services": strings.Join(services, ","),
			"category": "backend-as-a-service",
		},
	}
	return []*projectgraph.Node{node}, nil, nil
}

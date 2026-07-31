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

// PythonDetector parses requirements.txt and pyproject.toml files to detect
// Python language usage and its direct dependencies.
type PythonDetector struct{}

// Ensure it implements the interface
var _ detectors.FileDetector = (*PythonDetector)(nil)

func (d *PythonDetector) Handles(filename string) bool {
	return filename == "requirements.txt" || filename == "pyproject.toml" || filename == "poetry.lock" || filename == "Pipfile.lock"
}

func (d *PythonDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	if filename == "poetry.lock" || filename == "Pipfile.lock" {
		return d.detectLockfile(relPath)
	}

	switch filename {
	case "requirements.txt":
		return d.detectRequirements(relPath, content)
	case "pyproject.toml":
		return d.detectPyproject(relPath, content)
	default:
		return nil, nil, nil
	}
}

func (d *PythonDetector) detectRequirements(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	nodeID := "lang_python_root"
	if relPath != "requirements.txt" {
		nodeID = fmt.Sprintf("lang_python_%s", sanitizeID(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Python",
		Confidence: 1.0,
		Evidences: []evidence.Evidence{
			{
				Kind:        "manifest",
				Path:        relPath,
				Selector:    "requirements.txt",
				Sensitivity: evidence.PublicMetadata,
			},
		},
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Skip flags like -r, -c, --index-url etc.
		if strings.HasPrefix(line, "-") {
			continue
		}

		pkgName, pkgVersion := parsePythonRequirement(line)
		if pkgName == "" {
			continue
		}

		depID := fmt.Sprintf("dep_pypi_%s_%s", sanitizeID(pkgName), sanitizeID(relPath))
		if relPath == "requirements.txt" {
			depID = fmt.Sprintf("dep_pypi_%s_root", sanitizeID(pkgName))
		}

		purlVersion := pkgVersion
		purl := fmt.Sprintf("pkg:pypi/%s", strings.ToLower(pkgName))
		if purlVersion != "" {
			purl += "@" + purlVersion
		}

		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       pkgName,
			Version:    pkgVersion,
			Confidence: 1.0,
			Properties: map[string]string{
				"purl":    purl,
				"package": pkgName,
			},
			Evidences: []evidence.Evidence{
				{
					Kind:        "manifest",
					Path:        relPath,
					Selector:    pkgName,
					Value:       line,
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

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return nodes, edges, nil
}

// parsePythonRequirement parses a single requirements.txt line into (name, version).
// Handles: pkg==1.0, pkg>=1.0, pkg~=1.0, pkg (no version), pkg[extra]==1.0
func parsePythonRequirement(line string) (name, version string) {
	// Strip extras: requests[security]==2.31.0 -> requests==2.31.0
	if idx := strings.Index(line, "["); idx >= 0 {
		if end := strings.Index(line, "]"); end > idx {
			line = line[:idx] + line[end+1:]
		}
	}

	// Find operator
	for _, op := range []string{"===", "~=", "==", "!=", ">=", "<=", ">", "<"} {
		if idx := strings.Index(line, op); idx >= 0 {
			name = strings.TrimSpace(line[:idx])
			rest := strings.TrimSpace(line[idx+len(op):])
			// Strip any trailing constraints (e.g. "1.0,<2.0")
			if comma := strings.Index(rest, ","); comma >= 0 {
				rest = rest[:comma]
			}
			version = strings.TrimSpace(rest)
			return
		}
	}

	return strings.TrimSpace(line), ""
}

func (d *PythonDetector) detectPyproject(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	nodeID := "lang_python_root"
	if relPath != "pyproject.toml" {
		nodeID = fmt.Sprintf("lang_python_%s", sanitizeID(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Python",
		Confidence: 1.0,
		Evidences: []evidence.Evidence{
			{
				Kind:        "manifest",
				Path:        relPath,
				Selector:    "pyproject.toml",
				Sensitivity: evidence.PublicMetadata,
			},
		},
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	// Simple line-by-line TOML scanner — no external TOML library.
	// Supports [project.dependencies], [tool.poetry.dependencies]
	scanner := bufio.NewScanner(bytes.NewReader(content))
	inDepsSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect section headers
		if strings.HasPrefix(line, "[") {
			lower := strings.ToLower(line)
			inDepsSection = lower == "[project.dependencies]" ||
				lower == "[tool.poetry.dependencies]" ||
				lower == "[tool.pdm.dev-dependencies]"
			continue
		}

		if !inDepsSection || line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key = "value" or key = {version = "value", ...}
		eqIdx := strings.Index(line, "=")
		if eqIdx < 0 {
			continue
		}
		pkgName := strings.TrimSpace(line[:eqIdx])
		rest := strings.TrimSpace(line[eqIdx+1:])

		// Skip python version constraint entry
		if pkgName == "python" {
			continue
		}

		// Extract version: can be "^1.0.0" or {version = "^1.0.0", ...}
		pkgVersion := ""
		if strings.HasPrefix(rest, "{") {
			// Inline table: find version = "..."
			if vi := strings.Index(rest, "version"); vi >= 0 {
				sub := rest[vi+len("version"):]
				if qi := strings.Index(sub, "\""); qi >= 0 {
					sub = sub[qi+1:]
					if qi2 := strings.Index(sub, "\""); qi2 >= 0 {
						pkgVersion = sub[:qi2]
					}
				}
			}
		} else {
			// Simple: "^1.0.0" or "*"
			pkgVersion = strings.Trim(rest, `"' `)
		}

		if pkgName == "" {
			continue
		}

		depID := fmt.Sprintf("dep_pypi_%s_%s", sanitizeID(pkgName), sanitizeID(relPath))
		if relPath == "pyproject.toml" {
			depID = fmt.Sprintf("dep_pypi_%s_root", sanitizeID(pkgName))
		}

		purl := fmt.Sprintf("pkg:pypi/%s", strings.ToLower(pkgName))
		if pkgVersion != "" && pkgVersion != "*" {
			purl += "@" + strings.TrimLeft(pkgVersion, "^~>=<")
		}

		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       projectgraph.TypeDependency,
			Name:       pkgName,
			Version:    pkgVersion,
			Confidence: 1.0,
			Properties: map[string]string{
				"purl":    purl,
				"package": pkgName,
			},
			Evidences: []evidence.Evidence{
				{
					Kind:        "manifest",
					Path:        relPath,
					Selector:    pkgName,
					Value:       line,
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

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return nodes, edges, nil
}

func (d *PythonDetector) detectLockfile(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         "lang_python_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Python",
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

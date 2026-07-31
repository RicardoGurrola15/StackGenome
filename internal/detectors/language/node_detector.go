package language

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

// knownFrameworks maps package names to friendly names for JS/TS frameworks.
var knownFrameworks = map[string]string{
	"react":         "React",
	"react-dom":     "React",
	"next":          "Next.js",
	"vue":           "Vue",
	"@vue/core":     "Vue",
	"nuxt":          "Nuxt",
	"svelte":        "Svelte",
	"angular":       "Angular",
	"@angular/core": "Angular",
	"tailwindcss":   "TailwindCSS",
	"express":       "Express",
	"fastify":       "Fastify",
	"nestjs":        "NestJS",
	"@nestjs/core":  "NestJS",
}

// NodeDetector parses package.json files to detect JavaScript/TypeScript language usage,
// its direct dependencies, workspaces, and frameworks.
type NodeDetector struct{}

// Ensure it implements the interface
var _ detectors.FileDetector = (*NodeDetector)(nil)

func (d *NodeDetector) Handles(filename string) bool {
	return filename == "package.json" || filename == "package-lock.json" || filename == "yarn.lock" || filename == "pnpm-lock.yaml"
}

func (d *NodeDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	if filename != "package.json" {
		return d.detectLockfile(relPath)
	}

	var pkg struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Engines struct {
			Node string `json:"node"`
		} `json:"engines"`
		Workspaces      []string          `json:"workspaces"`
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(content, &pkg); err != nil {
		return nil, nil, fmt.Errorf("failed to parse package.json at %s: %w", relPath, err)
	}

	// Determine language (JS or TS)
	isTS := false
	if _, ok := pkg.DevDependencies["typescript"]; ok {
		isTS = true
	}
	if _, ok := pkg.Dependencies["typescript"]; ok {
		isTS = true
	}

	langName := "JavaScript"
	if isTS {
		langName = "TypeScript"
	}

	nodeID := "lang_node_root"
	if relPath != "package.json" {
		nodeID = fmt.Sprintf("lang_node_%s", sanitizeID(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       langName,
		Version:    pkg.Engines.Node,
		Confidence: 1.0,
		Evidences: []evidence.Evidence{
			{
				Kind:        "manifest",
				Path:        relPath,
				Selector:    "engines.node",
				Value:       pkg.Engines.Node,
				Sensitivity: evidence.PublicMetadata,
			},
		},
	}

	// Mark as workspace root if workspaces field exists
	if len(pkg.Workspaces) > 0 {
		langNode.Properties = map[string]string{
			"workspace": "true",
			"members":   strings.Join(pkg.Workspaces, ","),
		}
	}

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	// Merge all deps for iteration
	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}
	// devDeps don't get "depends_on" edges – skip for now to keep graph clean
	// They do contribute to type detection (e.g., typescript)

	for pkgName, pkgVersion := range allDeps {
		depID := fmt.Sprintf("dep_npm_%s_%s", sanitizeID(pkgName), sanitizeID(relPath))
		if relPath == "package.json" {
			depID = fmt.Sprintf("dep_npm_%s_root", sanitizeID(pkgName))
		}
		purl := fmt.Sprintf("pkg:npm/%s@%s", pkgName, strings.TrimPrefix(pkgVersion, "^"))

		nodeType := projectgraph.TypeDependency
		nodeName := pkgName
		// Upgrade to framework type if known
		if friendly, ok := knownFrameworks[pkgName]; ok {
			nodeType = projectgraph.TypeFramework
			nodeName = friendly
		}

		depNode := &projectgraph.Node{
			ID:         depID,
			Type:       nodeType,
			Name:       nodeName,
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
					Selector:    "dependencies." + pkgName,
					Value:       pkgVersion,
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

func (d *NodeDetector) detectLockfile(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         "lang_node_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       "Node.js",
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

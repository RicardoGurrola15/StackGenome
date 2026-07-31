package language

import (
	"encoding/xml"
	"path/filepath"
	"strings"

	"stackgenome/internal/detectors"
	"stackgenome/internal/evidence"
	"stackgenome/internal/projectgraph"
)

type DotNetDetector struct{}

var _ detectors.FileDetector = (*DotNetDetector)(nil)

func (d *DotNetDetector) Handles(filename string) bool {
	return strings.HasSuffix(filename, ".csproj") || strings.HasSuffix(filename, ".fsproj") || strings.HasSuffix(filename, ".sln") || filename == "packages.lock.json"
}

func (d *DotNetDetector) Detect(relPath string, content []byte) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	filename := filepath.Base(relPath)

	if filename == "packages.lock.json" {
		return d.detectLockfile(relPath)
	}

	sanitizedPath := "root"
	if relPath != filename {
		sanitizedPath = sanitizeID(relPath)
	}
	nodeID := "lang_dotnet_" + sanitizedPath

	langNode := &projectgraph.Node{
		ID:         nodeID,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       ".NET",
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

	if strings.HasSuffix(filename, ".sln") {
		langNode.Properties["workspace"] = "true"
		return []*projectgraph.Node{langNode}, nil, nil
	}

	type PackageReference struct {
		Include string `xml:"Include,attr"`
		Version string `xml:"Version,attr"`
	}
	type Project struct {
		ItemGroups []struct {
			PackageReferences []PackageReference `xml:"PackageReference"`
		} `xml:"ItemGroup"`
	}

	var proj Project
	_ = xml.Unmarshal(content, &proj)

	nodes := []*projectgraph.Node{langNode}
	var edges []projectgraph.Edge

	for _, ig := range proj.ItemGroups {
		for _, dep := range ig.PackageReferences {
			if dep.Include == "" {
				continue
			}
			depID := "dep_dotnet_" + sanitizeID(dep.Include)
			depNode := &projectgraph.Node{
				ID:         depID,
				Type:       projectgraph.TypeDependency,
				Name:       dep.Include,
				Version:    dep.Version,
				Confidence: 0.8,
				Properties: map[string]string{
					"purl": "pkg:nuget/" + dep.Include + "@" + dep.Version,
				},
			}
			nodes = append(nodes, depNode)
			edges = append(edges, projectgraph.Edge{
				SourceID: langNode.ID,
				TargetID: depID,
				Relation: projectgraph.RelationDependsOn,
			})
		}
	}

	return nodes, edges, nil
}

func (d *DotNetDetector) detectLockfile(relPath string) ([]*projectgraph.Node, []projectgraph.Edge, error) {
	sanitizedPath := "root"
	if filepath.Dir(relPath) != "." {
		sanitizedPath = sanitizeID(filepath.Dir(relPath))
	}

	langNode := &projectgraph.Node{
		ID:         "lang_dotnet_" + sanitizedPath,
		Type:       projectgraph.TypeLanguage,
		Role:       DetermineRole(relPath),
		Name:       ".NET",
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

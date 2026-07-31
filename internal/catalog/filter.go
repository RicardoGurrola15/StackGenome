package catalog

import (
	"strings"

	schemav1 "stackgenome/pkg/schema/v1"
)

// ProjectContext holds extracted sets from a ProjectGraphDTO for fast lookup.
type ProjectContext struct {
	Languages  map[string]bool
	Infra      map[string]bool
	Frameworks map[string]bool
	NodeIDs    map[string]bool
}

// BuildContext extracts a lookup-friendly context from an analyzed DTO.
func BuildContext(dto *schemav1.ProjectGraphDTO) *ProjectContext {
	ctx := &ProjectContext{
		Languages:  make(map[string]bool),
		Infra:      make(map[string]bool),
		Frameworks: make(map[string]bool),
		NodeIDs:    make(map[string]bool),
	}

	for _, n := range dto.Nodes {
		ctx.NodeIDs[n.ID] = true
		switch n.Type {
		case "language":
			ctx.Languages[strings.ToLower(n.Name)] = true
		case "infrastructure":
			ctx.Infra[strings.ToLower(n.Name)] = true
		case "framework":
			ctx.Frameworks[strings.ToLower(n.Name)] = true
		}
	}

	return ctx
}

// Filter returns only catalog entries that have at least one target language
// or infra dimension matching the project context.
// Entries with empty targets on all dimensions pass through (generic tools).
func Filter(entries []Entry, ctx *ProjectContext) []Entry {
	var result []Entry
	for _, e := range entries {
		if entryIsRelevant(e, ctx) {
			result = append(result, e)
		}
	}
	return result
}

func entryIsRelevant(e Entry, ctx *ProjectContext) bool {
	// If the entry targets nothing specific, it's universally relevant
	if len(e.Targets.Languages) == 0 && len(e.Targets.Infra) == 0 && len(e.Targets.Frameworks) == 0 {
		return true
	}

	for _, lang := range e.Targets.Languages {
		if ctx.Languages[strings.ToLower(lang)] {
			return true
		}
	}
	for _, infra := range e.Targets.Infra {
		if ctx.Infra[strings.ToLower(infra)] {
			return true
		}
	}
	for _, fw := range e.Targets.Frameworks {
		if ctx.Frameworks[strings.ToLower(fw)] {
			return true
		}
	}
	return false
}

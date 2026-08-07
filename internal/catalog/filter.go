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
	// PrimaryLanguage is the name of the highest-confidence language node.
	// Used by the scorer to apply the primaryLanguageBoost.
	PrimaryLanguage string
}

// BuildContext extracts a lookup-friendly context from an analyzed DTO.
func BuildContext(dto *schemav1.ProjectGraphDTO) *ProjectContext {
	ctx := &ProjectContext{
		Languages:  make(map[string]bool),
		Infra:      make(map[string]bool),
		Frameworks: make(map[string]bool),
		NodeIDs:    make(map[string]bool),
	}

	var primaryConf float64
	var primaryEvidences int
	for _, n := range dto.Nodes {
		ctx.NodeIDs[n.ID] = true
		switch n.Type {
		case "language":
			nameLower := strings.ToLower(n.Name)
			ctx.Languages[nameLower] = true
			// Split compound names like "dart/flutter" → add both "dart" and "flutter"
			if strings.Contains(nameLower, "/") {
				parts := strings.Split(nameLower, "/")
				for _, p := range parts {
					ctx.Languages[strings.TrimSpace(p)] = true
				}
			}
			// Track the primary language (highest confidence + most evidences on tie)
			evCount := len(n.Evidences)
			if n.Confidence > primaryConf || (n.Confidence == primaryConf && evCount > primaryEvidences) {
				primaryConf = n.Confidence
				primaryEvidences = evCount
				// Use the first part before "/" as the canonical primary name
				if idx := strings.Index(n.Name, "/"); idx > 0 {
					ctx.PrimaryLanguage = strings.ToLower(strings.TrimSpace(n.Name[:idx]))
				} else {
					ctx.PrimaryLanguage = strings.ToLower(n.Name)
				}
			}
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
	if len(e.Ecosystem) == 0 && len(e.Targets.Infra) == 0 && len(e.Targets.Frameworks) == 0 {
		return true
	}
	// Also treat ["*"] as universally relevant for ecosystem
	if len(e.Ecosystem) == 1 && e.Ecosystem[0] == "*" {
		return true
	}

	for _, lang := range e.Ecosystem {
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

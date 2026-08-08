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
	// Signals are derived boolean signals used for context-gating of entries.
	// e.g. "has_backend", "has_mobile", "has_firebase", "has_ci", "has_tests".
	Signals map[string]bool
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
		Signals:    make(map[string]bool),
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

	// Derive project-level signals from detected nodes.
	// has_mobile: project has a mobile platform (android/ios) or dart/flutter
	if ctx.Languages["dart"] || ctx.Languages["flutter"] ||
		ctx.Infra["android"] || ctx.Infra["ios"] {
		ctx.Signals["has_mobile"] = true
	}
	// has_backend: project has a server-side language / infra node
	if ctx.Languages["go"] || ctx.Languages["python"] || ctx.Languages["ruby"] ||
		ctx.Languages["php"] || ctx.Languages["java"] || ctx.Languages["node.js"] ||
		ctx.Languages["node"] || ctx.Languages[".net"] || ctx.Languages["dotnet"] ||
		ctx.Infra["docker"] || ctx.Infra["kubernetes"] {
		ctx.Signals["has_backend"] = true
	}
	// has_firebase: firebase infra node is present
	for id := range ctx.NodeIDs {
		if strings.Contains(id, "firebase") {
			ctx.Signals["has_firebase"] = true
			break
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
	// RequiresContext gate: ALL listed signals must be present.
	// If any required signal is absent, skip this entry entirely.
	for _, sig := range e.RequiresContext {
		if !ctx.Signals[sig] {
			return false
		}
	}

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

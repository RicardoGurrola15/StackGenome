// Package catalog provides the local recommendation engine for StackGenome.
// It loads a versioned seed catalog embedded in the binary and scores entries
// against the analyzed project graph to produce deterministic recommendations.
package catalog

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed catalog.json
var catalogData []byte

// Entry represents a single tool or resource in the catalog.
type Entry struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Category    string      `json:"category,omitempty"` // e.g. "backend", "mobile", "devops"
	Tags        []string    `json:"tags"`
	Targets     EntryTarget `json:"targets"`
	Ecosystem   []string    `json:"ecosystem"` // e.g. ["dart", "flutter", "go"]
	// RequiresContext lists signals that MUST be present in the project for
	// this entry to be considered. Valid values: "has_backend", "has_mobile",
	// "has_ci", "has_firebase", "has_tests". Empty means no extra requirement.
	RequiresContext []string `json:"requires_context,omitempty"`
	URL             string   `json:"url"`
}

// EntryTarget specifies which project characteristics make this entry relevant.
type EntryTarget struct {
	Languages  []string `json:"languages"`
	Frameworks []string `json:"frameworks"`
	Infra      []string `json:"infra"`
}

// Catalog holds the loaded seed data.
type Catalog struct {
	Version string  `json:"version"`
	Entries []Entry `json:"entries"`
}

// Load parses and returns the embedded catalog. Returns an error if the
// embedded JSON is malformed (should never happen in production builds).
func Load() (*Catalog, error) {
	var c Catalog
	if err := json.Unmarshal(catalogData, &c); err != nil {
		return nil, fmt.Errorf("failed to parse embedded catalog: %w", err)
	}
	return &c, nil
}

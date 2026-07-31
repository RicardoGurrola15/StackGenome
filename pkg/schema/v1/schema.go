package schema

import (
	"stackgenome/internal/evidence"
)

// ProjectGraphDTO is the versioned JSON output for the final report.
// This represents the stable contract (v1) that external tools will consume.
type ProjectGraphDTO struct {
	Version         string              `json:"schema_version"` // Schema version, e.g., "1.0.0"
	Nodes           []NodeDTO           `json:"nodes"`
	Edges           []EdgeDTO           `json:"edges"`
	Environment     *EnvironmentDTO     `json:"environment,omitempty"`
	Recommendations []RecommendationDTO `json:"recommendations,omitempty"`
}

// RecommendationDTO represents a suggested tool or resource from the local catalog.
type RecommendationDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Score       float64  `json:"score"`
	Reasons     []string `json:"reasons"`
	URL         string   `json:"url,omitempty"`
}

// EnvironmentDTO contains opt-in telemetry about the host system.
type EnvironmentDTO struct {
	OS    string            `json:"os"`
	Arch  string            `json:"arch"`
	Tools map[string]string `json:"tools,omitempty"`
}

// NodeDTO represents a node in the project graph for serialization.
type NodeDTO struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"` // e.g., "module", "language", "framework"
	Role       string              `json:"role,omitempty"`
	Name       string              `json:"name"`
	Version    string              `json:"version,omitempty"`
	Evidences  []evidence.Evidence `json:"evidences,omitempty"`
	Confidence float64             `json:"confidence"`
	Properties map[string]string   `json:"properties,omitempty"`
}

// EdgeDTO represents a serialized edge in the graph.
type EdgeDTO struct {
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
	Relation string `json:"relation"` // e.g., "contains", "depends_on"
}

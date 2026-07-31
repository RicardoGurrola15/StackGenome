package sanitizer

import (
	"testing"

	"stackgenome/internal/evidence"
	schemav1 "stackgenome/pkg/schema/v1"
)

// buildContaminatedDTO returns a DTO with sensitive data injected in every field.
func buildContaminatedDTO() *schemav1.ProjectGraphDTO {
	return &schemav1.ProjectGraphDTO{
		Version: "1.0.0",
		Nodes: []schemav1.NodeDTO{
			{
				ID:   "lang_go_root",
				Type: "language",
				Name: "Go",
				Evidences: []evidence.Evidence{
					{Kind: "file", Path: "/Users/ricardo/company-secret/internal/auth.go", Sensitivity: evidence.LocalPath},
					{Kind: "secret", Path: "internal/tokens/api_key=sk-abc123", Sensitivity: evidence.Secret},
				},
				Properties: map[string]string{
					"type":          "module",
					"internal_team": "platform-security", // should be stripped
					"deploy_key":    "aws-prod-key-7890", // should be stripped
					"workspace":     "monorepo",          // should be kept
				},
			},
		},
		Edges: []schemav1.EdgeDTO{
			{SourceID: "lang_go_root", TargetID: "dep_go_example", Relation: "depends_on"},
		},
		Environment: &schemav1.EnvironmentDTO{
			OS:   "darwin",
			Arch: "arm64",
			Tools: map[string]string{
				"go": "go version go1.22 darwin/arm64",
			},
		},
	}
}

// ─── FilterSecrets ────────────────────────────────────────────────────────────

func TestFilterSecrets_RemovesSecretEvidences(t *testing.T) {
	dto := buildContaminatedDTO()
	FilterSecrets(dto)

	for _, n := range dto.Nodes {
		for _, ev := range n.Evidences {
			if evidence.Sensitivity(ev.Sensitivity) == evidence.Secret {
				t.Errorf("Expected no secret evidences after FilterSecrets, found: %q", ev.Path)
			}
		}
	}
}

func TestFilterSecrets_KeepsNonSecretEvidences(t *testing.T) {
	dto := buildContaminatedDTO()
	FilterSecrets(dto)

	// LocalPath evidence should survive FilterSecrets
	found := false
	for _, n := range dto.Nodes {
		for _, ev := range n.Evidences {
			if ev.Kind == "file" {
				found = true
			}
		}
	}
	if !found {
		t.Error("Expected LocalPath evidence to survive FilterSecrets")
	}
}

// ─── Anonymize ────────────────────────────────────────────────────────────────

func TestAnonymize_StripAllEvidences(t *testing.T) {
	dto := buildContaminatedDTO()
	Anonymize(dto)

	for _, n := range dto.Nodes {
		if len(n.Evidences) > 0 {
			t.Errorf("Expected all evidences to be stripped in fingerprint mode, node %q still has %d", n.ID, len(n.Evidences))
		}
	}
}

func TestAnonymize_StripNonAllowlistedProperties(t *testing.T) {
	dto := buildContaminatedDTO()
	Anonymize(dto)

	for _, n := range dto.Nodes {
		if _, ok := n.Properties["internal_team"]; ok {
			t.Errorf("node %q: 'internal_team' property should have been stripped in fingerprint mode", n.ID)
		}
		if _, ok := n.Properties["deploy_key"]; ok {
			t.Errorf("node %q: 'deploy_key' property should have been stripped in fingerprint mode", n.ID)
		}
	}
}

func TestAnonymize_KeepsAllowlistedProperties(t *testing.T) {
	dto := buildContaminatedDTO()
	Anonymize(dto)

	for _, n := range dto.Nodes {
		if n.ID == "lang_go_root" {
			if _, ok := n.Properties["workspace"]; !ok {
				t.Errorf("node %q: 'workspace' property should be retained in fingerprint mode", n.ID)
			}
			if _, ok := n.Properties["type"]; !ok {
				t.Errorf("node %q: 'type' property should be retained in fingerprint mode", n.ID)
			}
		}
	}
}

func TestAnonymize_RedactsEnvironmentTools(t *testing.T) {
	dto := buildContaminatedDTO()
	Anonymize(dto)

	if dto.Environment == nil {
		t.Fatal("Expected environment block to be present")
	}
	if len(dto.Environment.Tools) > 0 {
		t.Errorf("Expected environment tools to be redacted in fingerprint mode, got: %v", dto.Environment.Tools)
	}
}

func TestAnonymize_KeepsEnvironmentOSArch(t *testing.T) {
	dto := buildContaminatedDTO()
	Anonymize(dto)

	if dto.Environment.OS == "" {
		t.Error("Expected OS to be retained in fingerprint mode")
	}
	if dto.Environment.Arch == "" {
		t.Error("Expected Arch to be retained in fingerprint mode")
	}
}

func TestAnonymize_NoPathLeakInProperties(t *testing.T) {
	dto := buildContaminatedDTO()
	// Inject a path-like string into a property
	dto.Nodes[0].Properties["custom_path"] = "/Users/ricardo/company/internal/secret"
	Anonymize(dto)

	for _, n := range dto.Nodes {
		for k, v := range n.Properties {
			if k == "custom_path" {
				t.Errorf("node %q: 'custom_path' (value: %q) should have been stripped", n.ID, v)
			}
		}
	}
}

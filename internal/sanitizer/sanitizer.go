// Package sanitizer provides the anonymization engine for StackGenome's
// metadata-only fingerprint mode. It applies a multi-layered privacy filter
// to a ProjectGraphDTO ensuring no proprietary project details, filesystem
// paths, or secret fields are leaked.
package sanitizer

import (
	"stackgenome/internal/evidence"
	schemav1 "stackgenome/pkg/schema/v1"
)

// FilterSecrets removes evidences marked as Secret from the DTO.
// This is always applied regardless of mode.
func FilterSecrets(dto *schemav1.ProjectGraphDTO) {
	for i, node := range dto.Nodes {
		var safe []evidence.Evidence
		for _, ev := range node.Evidences {
			if ev.Sensitivity != evidence.Secret {
				safe = append(safe, ev)
			}
		}
		dto.Nodes[i].Evidences = safe
	}
}

// Anonymize applies full privacy projection (fingerprint mode).
// It performs:
//  1. Secret evidence removal (same as FilterSecrets)
//  2. Strips ALL remaining evidences (filesystem paths leak project structure)
//  3. Strips non-standard Properties (may contain team/org-level data)
//  4. Redacts Environment tools output (version strings may leak vendor info)
//  5. Clears the Environment host info (OS/arch retained as they are generic)
func Anonymize(dto *schemav1.ProjectGraphDTO) {
	for i, node := range dto.Nodes {
		// Step 2: strip ALL evidences in fingerprint mode (paths are PII for projects)
		dto.Nodes[i].Evidences = nil

		// Step 3: strip Properties — keep only a hardcoded allowlist of harmless keys.
		filtered := filterProperties(node.Properties)
		dto.Nodes[i].Properties = filtered
	}

	// Step 4: keep Environment OS/Arch but redact tool version strings
	if dto.Environment != nil {
		dto.Environment.Tools = nil
	}
}

// propertyAllowlist defines which property keys are safe to expose in
// fingerprint mode. All other keys are dropped.
var propertyAllowlist = map[string]bool{
	"type":      true,
	"framework": true,
	"scope":     true,
	"workspace": true,
	"purl_type": true,
}

func filterProperties(props map[string]string) map[string]string {
	if len(props) == 0 {
		return nil
	}
	out := make(map[string]string, len(propertyAllowlist))
	for k, v := range props {
		if propertyAllowlist[k] {
			out[k] = v
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

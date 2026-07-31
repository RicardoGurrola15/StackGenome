package projectgraph

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"stackgenome/internal/evidence"
)

func TestGraph_AddNode_Validation(t *testing.T) {
	g := NewGraph()

	// 1. Validation: Unique IDs
	n1 := &Node{ID: "node1"}
	if err := g.AddNode(n1); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	n2 := &Node{
		ID: "node1",
		Evidences: []evidence.Evidence{
			{Path: "some/path.txt"},
		},
	}
	if err := g.AddNode(n2); err != nil {
		t.Fatalf("expected no error when adding node with duplicate ID (should merge), got %v", err)
	}
	if len(g.nodes["node1"].Evidences) != 1 {
		t.Fatalf("expected evidence to be merged")
	}

	// 2. Validation: Empty ID
	n3 := &Node{ID: ""}
	if err := g.AddNode(n3); err == nil {
		t.Fatal("expected error adding node with empty ID, got nil")
	}

	// 3. Validation: Relative paths in evidence
	n4 := &Node{
		ID: "node4",
		Evidences: []evidence.Evidence{
			{Path: "/absolute/path/file.txt"},
		},
	}
	if err := g.AddNode(n4); err == nil {
		t.Fatal("expected error adding node with absolute path, got nil")
	}
}

func TestGraph_AddEdge_Validation(t *testing.T) {
	g := NewGraph()
	g.AddNode(&Node{ID: "source"})

	if err := g.AddEdge("source", "target", RelationContains); err == nil {
		t.Fatal("expected error adding edge to missing target node, got nil")
	}
}

func TestGraph_Serialization_Golden(t *testing.T) {
	g := NewGraph()

	g.AddNode(&Node{
		ID:         "mod_stackgenome",
		Type:       TypeModule,
		Name:       "stackgenome",
		Confidence: 0.9,
		Evidences: []evidence.Evidence{
			{
				Kind:        "manifest",
				Path:        "go.mod",
				Selector:    "module",
				Value:       "stackgenome",
				Sensitivity: evidence.PublicMetadata,
			},
		},
	})

	g.AddNode(&Node{
		ID:         "lang_go",
		Type:       TypeLanguage,
		Name:       "Go",
		Version:    "1.26.x",
		Confidence: 1.0,
		Properties: map[string]string{
			"primary": "true",
		},
		Evidences: []evidence.Evidence{
			{
				Kind:        "manifest",
				Path:        "go.mod",
				Value:       "1.26",
				Sensitivity: evidence.PublicMetadata,
			},
		},
	})

	g.AddEdge("mod_stackgenome", "lang_go", RelationContains)

	dto := g.ToDTO()
	b, err := json.MarshalIndent(dto, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal DTO: %v", err)
	}

	goldenPath := filepath.Join("testdata", "projectgraph_v1.golden.json")

	// Ensure a newline at EOF for comparison if desired, but we can just use json.Indent.

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	// Compare ignoring exact whitespace optionally, but since it's a golden test let's compare identically (removing trailing newline differences).
	actualStr := string(bytes.TrimSpace(b))
	goldenStr := string(bytes.TrimSpace(golden))

	if actualStr != goldenStr {
		t.Errorf("serialization output does not match golden file.\nExpected:\n%s\nGot:\n%s", goldenStr, actualStr)
	}
}

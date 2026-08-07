package catalog

import (
	"testing"

	schemav1 "stackgenome/pkg/schema/v1"
)

// ─── helpers ─────────────────────────────────────────────────────────────────

func makeGoDockerDTO() *schemav1.ProjectGraphDTO {
	return &schemav1.ProjectGraphDTO{
		Version: "1.0.0",
		Nodes: []schemav1.NodeDTO{
			{ID: "lang_go_root", Type: "language", Name: "Go"},
			{ID: "infra_root_Dockerfile", Type: "infrastructure", Name: "Docker"},
		},
	}
}

func makeRustOnlyDTO() *schemav1.ProjectGraphDTO {
	return &schemav1.ProjectGraphDTO{
		Version: "1.0.0",
		Nodes: []schemav1.NodeDTO{
			{ID: "lang_rust_root", Type: "language", Name: "Rust"},
		},
	}
}

// ─── Filter tests ─────────────────────────────────────────────────────────────

func TestFilter_KeepsRelevantEntries(t *testing.T) {
	ctx := BuildContext(makeGoDockerDTO())

	goOnlyEntry := Entry{
		ID:        "tool:golangci-lint",
		Name:      "golangci-lint",
		Ecosystem: []string{"go"},
	}
	pythonOnlyEntry := Entry{
		ID:        "tool:ruff",
		Name:      "Ruff",
		Ecosystem: []string{"python"},
	}

	filtered := Filter([]Entry{goOnlyEntry, pythonOnlyEntry}, ctx)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 entry after filter, got %d", len(filtered))
	}
	if filtered[0].ID != "tool:golangci-lint" {
		t.Errorf("Expected golangci-lint to survive filter, got %q", filtered[0].ID)
	}
}

func TestFilter_UniversalToolPassesAlways(t *testing.T) {
	ctx := BuildContext(makeRustOnlyDTO())

	universal := Entry{
		ID:        "tool:universal",
		Name:      "Universal",
		Ecosystem: []string{}, // no targets specified
	}

	filtered := Filter([]Entry{universal}, ctx)
	if len(filtered) != 1 {
		t.Error("Expected universal tool to pass filter for any project")
	}
}

func TestFilter_InfraMatchSuffices(t *testing.T) {
	ctx := BuildContext(makeGoDockerDTO())

	dockerOnlyEntry := Entry{
		ID:   "tool:trivy",
		Name: "Trivy",
		Targets: EntryTarget{
			Infra: []string{"docker"},
		},
	}

	filtered := Filter([]Entry{dockerOnlyEntry}, ctx)
	if len(filtered) != 1 {
		t.Error("Expected docker-targeting entry to survive filter when Docker is in infra")
	}
}

// ─── Scorer tests ─────────────────────────────────────────────────────────────

func TestScore_LanguageHitGivesCorrectWeight(t *testing.T) {
	ctx := BuildContext(makeGoDockerDTO())
	entry := Entry{
		ID:        "tool:golangci-lint",
		Name:      "golangci-lint",
		Ecosystem: []string{"go"},
	}

	scored := Score([]Entry{entry}, ctx)
	if len(scored) == 0 {
		t.Fatal("Expected at least one scored entry")
	}

	// Language match (0.40) + novelty (0.10) = 0.50
	expected := 0.50
	if scored[0].Score != expected {
		t.Errorf("Expected score %.2f, got %.2f", expected, scored[0].Score)
	}
}

func TestScore_LanguageAndInfraHitGivesHigherScore(t *testing.T) {
	ctx := BuildContext(makeGoDockerDTO())
	entry := Entry{
		ID:        "tool:testcontainers",
		Name:      "Testcontainers",
		Ecosystem: []string{"go"},
		Targets: EntryTarget{
			Infra: []string{"docker"},
		},
	}

	scored := Score([]Entry{entry}, ctx)
	if len(scored) == 0 {
		t.Fatal("Expected at least one scored entry")
	}

	// Language (0.40) + Infra (0.30) + novelty (0.10) = 0.80
	expected := 0.80
	if scored[0].Score != expected {
		t.Errorf("Expected score %.2f, got %.2f", expected, scored[0].Score)
	}
}

func TestScore_IsDeterministic(t *testing.T) {
	ctx := BuildContext(makeGoDockerDTO())
	entries := []Entry{
		{ID: "tool:b", Name: "B", Ecosystem: []string{"go"}},
		{ID: "tool:a", Name: "A", Ecosystem: []string{"go"}},
	}

	first := Score(entries, ctx)
	second := Score(entries, ctx)

	if len(first) != len(second) {
		t.Fatal("Non-deterministic length")
	}
	for i := range first {
		if first[i].ID != second[i].ID || first[i].Score != second[i].Score {
			t.Errorf("Non-deterministic result at index %d: %v vs %v", i, first[i].ID, second[i].ID)
		}
	}
}

func TestScore_TiebreakByID(t *testing.T) {
	ctx := BuildContext(makeGoDockerDTO())
	// Both score identically (language match only)
	entries := []Entry{
		{ID: "tool:zzz", Name: "ZZZ", Ecosystem: []string{"go"}},
		{ID: "tool:aaa", Name: "AAA", Ecosystem: []string{"go"}},
	}

	scored := Score(entries, ctx)
	if scored[0].ID != "tool:aaa" {
		t.Errorf("Expected tool:aaa first (tiebreak by ID), got %q", scored[0].ID)
	}
}

func TestTopN_LimitsOutput(t *testing.T) {
	scored := []ScoredEntry{
		{Entry: Entry{ID: "a"}, Score: 0.9},
		{Entry: Entry{ID: "b"}, Score: 0.8},
		{Entry: Entry{ID: "c"}, Score: 0.7},
		{Entry: Entry{ID: "d"}, Score: 0.6},
	}
	top := TopN(scored, 3)
	if len(top) != 3 {
		t.Errorf("Expected 3 items, got %d", len(top))
	}
}

func TestTopN_DoesNotPanicWhenFewerThanN(t *testing.T) {
	scored := []ScoredEntry{
		{Entry: Entry{ID: "a"}, Score: 0.9},
	}
	top := TopN(scored, 10) // ask for 10, only 1 available
	if len(top) != 1 {
		t.Errorf("Expected 1 item, got %d", len(top))
	}
}

// ─── Integration: Load + Filter + Score ──────────────────────────────────────

func TestCatalog_LoadAndRecommend(t *testing.T) {
	cat, err := Load()
	if err != nil {
		t.Fatalf("Failed to load catalog: %v", err)
	}
	if len(cat.Entries) == 0 {
		t.Fatal("Expected non-empty catalog")
	}

	ctx := BuildContext(makeGoDockerDTO())
	filtered := Filter(cat.Entries, ctx)
	scored := Score(filtered, ctx)
	top := TopN(scored, 3)

	if len(top) == 0 {
		t.Error("Expected at least one recommendation for a Go+Docker project")
	}
	for _, rec := range top {
		if rec.Score <= 0 {
			t.Errorf("Recommendation %q has non-positive score", rec.ID)
		}
		if len(rec.Reasons) == 0 {
			t.Errorf("Recommendation %q has no reasons", rec.ID)
		}
	}
}

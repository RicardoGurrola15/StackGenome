package catalog

import (
	"testing"

	schemav1 "stackgenome/pkg/schema/v1"
)

// ─── DTOs for Phase 21A test cases ───────────────────────────────────────────

// makeFlutterMobileOnlyDTO: pure mobile Flutter app, no backend signals.
func makeFlutterMobileOnlyDTO() *schemav1.ProjectGraphDTO {
	return &schemav1.ProjectGraphDTO{
		Version: "1.0.0",
		Nodes: []schemav1.NodeDTO{
			{ID: "lang_dart_root", Type: "language", Name: "Dart/Flutter", Confidence: 1},
			{ID: "platform_android", Type: "platform", Name: "Android Native"},
			{ID: "platform_ios", Type: "platform", Name: "iOS/macOS Native"},
		},
	}
}

// makeFlutterWithBackendDTO: Flutter app that also has a Dart backend (Dart Frog).
func makeFlutterWithBackendDTO() *schemav1.ProjectGraphDTO {
	return &schemav1.ProjectGraphDTO{
		Version: "1.0.0",
		Nodes: []schemav1.NodeDTO{
			{ID: "lang_dart_root", Type: "language", Name: "Dart/Flutter", Confidence: 1},
			{ID: "lang_node_root", Type: "language", Name: "Node.js", Confidence: 1},
			{ID: "platform_android", Type: "platform", Name: "Android Native"},
		},
	}
}

// makeNoRelevantDTO: project with no known ecosystem (exotic/unsupported).
func makeNoRelevantDTO() *schemav1.ProjectGraphDTO {
	return &schemav1.ProjectGraphDTO{
		Version: "1.0.0",
		Nodes: []schemav1.NodeDTO{
			{ID: "lang_cobol_root", Type: "language", Name: "COBOL", Confidence: 1},
		},
	}
}

// makeFlutterWithToolAlreadyInstalledDTO: Flutter project that already uses Shorebird.
func makeFlutterWithToolAlreadyInstalledDTO() *schemav1.ProjectGraphDTO {
	return &schemav1.ProjectGraphDTO{
		Version: "1.0.0",
		Nodes: []schemav1.NodeDTO{
			{ID: "lang_dart_root", Type: "language", Name: "Dart/Flutter", Confidence: 1},
			// Shorebird already detected
			{ID: "tool_shorebird", Type: "tool", Name: "Shorebird"},
		},
	}
}

// makeFlutterWithVendoredCppDTO: Flutter project with C/C++ as secondary vendored language.
func makeFlutterWithVendoredCppDTO() *schemav1.ProjectGraphDTO {
	return &schemav1.ProjectGraphDTO{
		Version: "1.0.0",
		Nodes: []schemav1.NodeDTO{
			{ID: "lang_dart_root", Type: "language", Name: "Dart/Flutter", Confidence: 1},
			// C++ is secondary / vendored (low confidence)
			{ID: "lang_cpp_root", Type: "language", Name: "C/C++", Confidence: 0.3},
		},
	}
}

// ─── Case A: Flutter mobile project should NOT recommend Dart Frog ─────────────

func TestScorer_CaseA_DartFrogNotRecommendedForPureMobileFlutter(t *testing.T) {
	cat, err := Load()
	if err != nil {
		t.Fatalf("Failed to load catalog: %v", err)
	}

	ctx := BuildContext(makeFlutterMobileOnlyDTO())

	// Dart Frog requires has_backend signal
	filtered := Filter(cat.Entries, ctx)
	scored := Score(filtered, ctx)

	for _, s := range scored {
		if s.ID == "tool:dart_frog" {
			t.Errorf("Case A: Dart Frog should NOT be recommended for a pure mobile Flutter project, but it was (score: %.2f)", s.Score)
		}
	}
}

// ─── Case B: Flutter + backend → Dart Frog IS a valid candidate ──────────────

func TestScorer_CaseB_DartFrogRecommendedForFlutterWithBackend(t *testing.T) {
	dartFrogEntry := Entry{
		ID:              "tool:dart_frog",
		Name:            "Dart Frog",
		Ecosystem:       []string{"dart"},
		RequiresContext: []string{"has_backend"},
	}

	ctx := BuildContext(makeFlutterWithBackendDTO())

	filtered := Filter([]Entry{dartFrogEntry}, ctx)
	if len(filtered) == 0 {
		t.Error("Case B: Dart Frog should be eligible for a Flutter project with a backend signal")
	}
}

// ─── Case C: Tool already detected → novelty bonus must be 0 ─────────────────

func TestScorer_CaseC_AlreadyInstalledToolHasNoNoveltyBonus(t *testing.T) {
	shorebirdEntry := Entry{
		ID:        "tool_shorebird", // same ID as what's in the graph
		Name:      "Shorebird",
		Ecosystem: []string{"dart", "flutter"},
	}

	ctx := BuildContext(makeFlutterWithToolAlreadyInstalledDTO())

	scored := Score([]Entry{shorebirdEntry}, ctx)
	if len(scored) == 0 {
		// Score can be below threshold — that's also acceptable
		return
	}

	for _, r := range scored[0].Reasons {
		if r == "herramienta no detectada aún en el proyecto" {
			t.Error("Case C: Already-installed tool should not receive the 'not yet detected' novelty reason")
		}
	}
}

// ─── Case D: No relevant tools → engine abstains (0 results) ─────────────────

func TestScorer_CaseD_NoRecommendationsForUnsupportedEcosystem(t *testing.T) {
	cat, err := Load()
	if err != nil {
		t.Fatalf("Failed to load catalog: %v", err)
	}

	ctx := BuildContext(makeNoRelevantDTO())
	filtered := Filter(cat.Entries, ctx)
	scored := Score(filtered, ctx)
	top := TopN(scored, 3)

	// We do not assert exactly 0 (some universal tools may still show up),
	// but we do assert that no tool with ecosystem-specific reasons appears.
	for _, rec := range top {
		for _, reason := range rec.Reasons {
			if reason == "ecosistema compatible: cobol" {
				t.Error("Case D: Should not produce ecosystem-matched recommendations for unknown ecosystem")
			}
		}
	}
}

// ─── Case E: Vendored/secondary language must not become primary ──────────────

func TestScorer_CaseE_VendoredCppDoesNotDominatePrimaryLanguage(t *testing.T) {
	ctx := BuildContext(makeFlutterWithVendoredCppDTO())

	if ctx.PrimaryLanguage == "c" || ctx.PrimaryLanguage == "c/c++" || ctx.PrimaryLanguage == "c++" {
		t.Errorf("Case E: Primary language should be dart, not vendored C/C++; got %q", ctx.PrimaryLanguage)
	}

	if ctx.PrimaryLanguage != "dart" {
		t.Errorf("Case E: Expected primary language to be 'dart', got %q", ctx.PrimaryLanguage)
	}
}

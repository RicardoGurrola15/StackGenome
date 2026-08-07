package analyzer

import (
	"path/filepath"
	"testing"

	"stackgenome/internal/detectors"
	_ "stackgenome/internal/detectors/cicd"
	_ "stackgenome/internal/detectors/editor"
	_ "stackgenome/internal/detectors/infra"
	_ "stackgenome/internal/detectors/language"
	_ "stackgenome/internal/detectors/platform"
)

func allDetectors() []detectors.FileDetector {
	return detectors.DefaultRegistry()
}

// setupAnalyzerAndRun is a helper to run the analyzer on a given path
func setupAnalyzerAndRun(t *testing.T, dir string, offline bool) map[string]struct{} {
	t.Helper()
	fixturesPath := filepath.Join("testdata", dir)
	a := NewAnalyzer(fixturesPath, allDetectors())

	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()
	nodesByID := make(map[string]struct{})
	for _, n := range dto.Nodes {
		nodesByID[n.ID] = struct{}{}
	}
	return nodesByID
}

func TestAnalyzer_DetectGoLanguageAndDeps(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures")
	a := NewAnalyzer(fixturesPath, allDetectors())

	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	nodesByID := make(map[string]struct{})
	for _, n := range dto.Nodes {
		nodesByID[n.ID] = struct{}{}
	}

	// Go language node
	if _, ok := nodesByID["lang_go_root"]; !ok {
		t.Errorf("expected lang_go_root node")
	}

	// Go dependency nodes (from go.mod require block)
	expectedGoDeps := []string{
		"dep_go_github_com_stretchr_testify_root",
		"dep_go_github_com_spf13_cobra_root",
		"dep_go_golang_org_x_sync_root",
	}
	for _, id := range expectedGoDeps {
		if _, ok := nodesByID[id]; !ok {
			t.Errorf("expected Go dep node %q", id)
		}
	}
}

func TestAnalyzer_DetectGoWorkspace(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures")
	a := NewAnalyzer(fixturesPath, allDetectors())

	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	nodesByID := make(map[string]struct{})
	for _, n := range dto.Nodes {
		nodesByID[n.ID] = struct{}{}
	}

	// go.work should produce a workspace node
	if _, ok := nodesByID["workspace_go_root"]; !ok {
		t.Errorf("expected workspace_go_root node from go.work")
	}
}

func TestAnalyzer_DetectNodeWorkspaceAndDeps(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures")
	a := NewAnalyzer(fixturesPath, allDetectors())

	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	nodesByID := make(map[string]struct{})
	propsByID := make(map[string]map[string]string)
	for _, n := range dto.Nodes {
		nodesByID[n.ID] = struct{}{}
		propsByID[n.ID] = n.Properties
	}

	// TypeScript language node should exist with workspace property
	if _, ok := nodesByID["lang_node_root"]; !ok {
		t.Errorf("expected lang_node_root node")
	}
	if props, ok := propsByID["lang_node_root"]; ok {
		if props["workspace"] != "true" {
			t.Errorf("expected lang_node_root to have workspace=true, got %q", props["workspace"])
		}
	}

	// Dep nodes from package.json#dependencies
	expectedNpmDeps := []string{
		"dep_npm_react_root",
		"dep_npm_react_dom_root",
		"dep_npm_axios_root",
	}
	for _, id := range expectedNpmDeps {
		if _, ok := nodesByID[id]; !ok {
			t.Errorf("expected npm dep node %q", id)
		}
	}

	// React should be classified as framework type
	for _, n := range dto.Nodes {
		if n.ID == "dep_npm_react_root" {
			if n.Type != "framework" {
				t.Errorf("expected react node type to be 'framework', got %q", n.Type)
			}
			// Verify PURL is set
			if n.Properties["purl"] == "" {
				t.Errorf("expected purl property on react node")
			}
			break
		}
	}
}

func TestAnalyzer_DetectPythonDeps(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures")
	a := NewAnalyzer(fixturesPath, allDetectors())

	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	nodesByID := make(map[string]struct{})
	for _, n := range dto.Nodes {
		nodesByID[n.ID] = struct{}{}
	}

	// Python language node
	if _, ok := nodesByID["lang_python_root"]; !ok {
		t.Errorf("expected lang_python_root node")
	}

	// Python dep nodes from requirements.txt
	expectedPyDeps := []string{
		"dep_pypi_requests_root",
		"dep_pypi_flask_root",
		"dep_pypi_sqlalchemy_root",
	}
	for _, id := range expectedPyDeps {
		if _, ok := nodesByID[id]; !ok {
			t.Errorf("expected Python dep node %q", id)
		}
	}

	// Verify PURL format
	for _, n := range dto.Nodes {
		if n.ID == "dep_pypi_requests_root" {
			purl := n.Properties["purl"]
			if purl != "pkg:pypi/requests@2.31.0" {
				t.Errorf("expected PURL pkg:pypi/requests@2.31.0, got %q", purl)
			}
			break
		}
	}
}

func TestAnalyzer_DetectRustDeps(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures")
	a := NewAnalyzer(fixturesPath, allDetectors())

	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	nodesByID := make(map[string]struct{})
	for _, n := range dto.Nodes {
		nodesByID[n.ID] = struct{}{}
	}

	// Rust language node
	if _, ok := nodesByID["lang_rust_root"]; !ok {
		t.Errorf("expected lang_rust_root node")
	}

	// Rust dep nodes from Cargo.toml
	expectedRustDeps := []string{
		"dep_cargo_serde_root",
		"dep_cargo_tokio_root",
		"dep_cargo_clap_root",
	}
	for _, id := range expectedRustDeps {
		if _, ok := nodesByID[id]; !ok {
			t.Errorf("expected Rust dep node %q", id)
		}
	}

	// Verify PURL format
	for _, n := range dto.Nodes {
		if n.ID == "dep_cargo_clap_root" {
			purl := n.Properties["purl"]
			if purl != "pkg:cargo/clap@4.4.0" {
				t.Errorf("expected PURL pkg:cargo/clap@4.4.0, got %q", purl)
			}
			break
		}
	}
}

func TestAnalyzer_DetectJVMDeps(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	// pom.xml
	if _, ok := nodesByID["lang_jvm_root"]; !ok {
		t.Errorf("expected lang_jvm_root node")
	}
	if _, ok := nodesByID["dep_jvm_org_springframework_boot_spring_boot_starter_web"]; !ok {
		t.Errorf("expected spring-boot-starter-web dependency node")
	}

	// build.gradle (mocked as the same root, we just check its dependencies)
	if _, ok := nodesByID["dep_jvm_com_google_guava_guava"]; !ok {
		t.Errorf("expected guava dependency node")
	}
}

func TestAnalyzer_DetectSwiftDeps(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	if _, ok := nodesByID["lang_swift_root"]; !ok {
		t.Errorf("expected lang_swift_root node")
	}
	if _, ok := nodesByID["dep_swift_Alamofire"]; !ok {
		t.Errorf("expected Alamofire dependency node")
	}
	if _, ok := nodesByID["dep_swift_AFNetworking"]; !ok {
		t.Errorf("expected AFNetworking dependency node")
	}
}

func TestAnalyzer_DetectDotNetDeps(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	if _, ok := nodesByID["lang_dotnet_root"]; !ok {
		t.Errorf("expected lang_dotnet_root node")
	}
	if _, ok := nodesByID["dep_dotnet_Newtonsoft_Json"]; !ok {
		t.Errorf("expected Newtonsoft.Json dependency node")
	}
}

func TestAnalyzer_DetectDartDeps(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	if _, ok := nodesByID["lang_dart_root"]; !ok {
		t.Errorf("expected lang_dart_root node")
	}
	if _, ok := nodesByID["dep_dart_http"]; !ok {
		t.Errorf("expected http dependency node")
	}
}

func TestAnalyzer_DetectPHPDeps(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	if _, ok := nodesByID["lang_php_root"]; !ok {
		t.Errorf("expected lang_php_root node")
	}
	if _, ok := nodesByID["dep_php_guzzlehttp_guzzle"]; !ok {
		t.Errorf("expected guzzle dependency node")
	}
}

func TestAnalyzer_DetectRubyDeps(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	if _, ok := nodesByID["lang_ruby_root"]; !ok {
		t.Errorf("expected lang_ruby_root node")
	}
	if _, ok := nodesByID["dep_ruby_rails"]; !ok {
		t.Errorf("expected rails dependency node")
	}
}

func TestAnalyzer_DetectCppDeps(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	if _, ok := nodesByID["lang_cpp_root"]; !ok {
		t.Errorf("expected lang_cpp_root node")
	}
	if _, ok := nodesByID["dep_cpp_Boost"]; !ok {
		t.Errorf("expected Boost dependency node")
	}
	if _, ok := nodesByID["dep_cpp_zlib"]; !ok {
		t.Errorf("expected zlib dependency node")
	}
}

func TestAnalyzer_DetectInfra(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	if _, ok := nodesByID["infra_root_Dockerfile"]; !ok {
		t.Errorf("expected Dockerfile infra node")
	}
	if _, ok := nodesByID["infra_root_main_tf"]; !ok {
		t.Errorf("expected main.tf infra node")
	}
}

func TestAnalyzer_DetectCICD(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	// sanitized path for .github/workflows is _github_workflows
	if _, ok := nodesByID["cicd__github_workflows_ci_yml"]; !ok {
		t.Errorf("expected GitHub Actions node cicd__github_workflows_ci_yml")
	}
}

func TestAnalyzer_DetectEditor(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	if _, ok := nodesByID["editor__vscode_settings_json"]; !ok {
		t.Errorf("expected VSCode settings node editor__vscode_settings_json")
	}
}

func TestAnalyzer_DetectPlatform(t *testing.T) {
	nodesByID := setupAnalyzerAndRun(t, "fixtures", false)

	if _, ok := nodesByID["platform_root_AndroidManifest_xml"]; !ok {
		t.Errorf("expected AndroidManifest node platform_root_AndroidManifest_xml")
	}
}

func TestAnalyzer_AllEdgesHaveValidNodes(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures")
	a := NewAnalyzer(fixturesPath, allDetectors())

	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	nodeSet := make(map[string]struct{})
	for _, n := range dto.Nodes {
		nodeSet[n.ID] = struct{}{}
	}

	// Every edge must reference existing nodes
	for _, e := range dto.Edges {
		if _, ok := nodeSet[e.SourceID]; !ok {
			t.Errorf("edge source %q does not exist in nodes", e.SourceID)
		}
		if _, ok := nodeSet[e.TargetID]; !ok {
			t.Errorf("edge target %q does not exist in nodes", e.TargetID)
		}
	}

	if len(dto.Edges) == 0 {
		t.Errorf("expected at least one edge, got none")
	}
}

// TestFlutterFull_VendoredPodsNotContaminating verifies F-001:
// C/C++ files inside ios/Pods/ must NOT appear as language nodes.
func TestFlutterFull_VendoredPodsNotContaminating(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures", "flutter_full")
	a := NewAnalyzer(fixturesPath, allDetectors())
	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	for _, n := range dto.Nodes {
		if n.Type == "language" && (n.Name == "C/C++" || n.Name == "C/C++ Header") {
			t.Errorf("F-001 regression: vendored C/C++ from ios/Pods appeared as language node: id=%q name=%q", n.ID, n.Name)
		}
	}
}

// TestFlutterFull_DartDetected verifies Dart/Flutter is the primary language.
func TestFlutterFull_DartDetected(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures", "flutter_full")
	a := NewAnalyzer(fixturesPath, allDetectors())
	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	found := false
	for _, n := range dto.Nodes {
		if n.Type == "language" && n.Name == "Dart/Flutter" {
			found = true
			if n.Confidence < 0.9 {
				t.Errorf("Dart/Flutter confidence too low: %.2f", n.Confidence)
			}
		}
	}
	if !found {
		t.Error("expected Dart/Flutter language node, not found")
	}
}

// TestFlutterFull_ShorebirdDetected verifies F-006: Shorebird is detected as a tool.
func TestFlutterFull_ShorebirdDetected(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures", "flutter_full")
	a := NewAnalyzer(fixturesPath, allDetectors())
	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	found := false
	for _, n := range dto.Nodes {
		if n.ID == "tool_shorebird" {
			found = true
			if _, hasAppID := n.Properties["app_id"]; hasAppID {
				t.Error("F-006 privacy: tool_shorebird node must NOT expose app_id")
			}
		}
	}
	if !found {
		t.Error("F-006: expected tool_shorebird node, not found")
	}
}

// TestFlutterFull_FirebaseDetected verifies F-007: Firebase is detected as infrastructure
// without leaking project IDs.
func TestFlutterFull_FirebaseDetected(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures", "flutter_full")
	a := NewAnalyzer(fixturesPath, allDetectors())
	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	found := false
	for _, n := range dto.Nodes {
		if n.ID == "infra_firebase" {
			found = true
			services, ok := n.Properties["services"]
			if !ok || services == "" {
				t.Error("F-007: infra_firebase node has no detected services")
			}
			if _, hasProjectID := n.Properties["project_id"]; hasProjectID {
				t.Error("F-007 privacy: infra_firebase node must NOT expose project_id")
			}
		}
	}
	if !found {
		t.Error("F-007: expected infra_firebase infrastructure node, not found")
	}
}

// TestFlutterFull_DartDepsHaveVersionAndScope verifies F-004 and F-005.
func TestFlutterFull_DartDepsHaveVersionAndScope(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures", "flutter_full")
	a := NewAnalyzer(fixturesPath, allDetectors())
	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	type depInfo struct{ version, scope string }
	var goRouter, buildRunner *depInfo
	for _, n := range dto.Nodes {
		if n.Type != "dependency" {
			continue
		}
		switch n.ID {
		case "dep_dart_go_router":
			goRouter = &depInfo{n.Version, n.Scope}
		case "dep_dart_build_runner":
			buildRunner = &depInfo{n.Version, n.Scope}
		}
	}

	if goRouter == nil {
		t.Error("F-004: dep_dart_go_router not found")
	} else {
		if goRouter.version == "" {
			t.Error("F-004: dep_dart_go_router has no declared version")
		}
		if goRouter.scope != "runtime" {
			t.Errorf("F-005: dep_dart_go_router scope should be 'runtime', got %q", goRouter.scope)
		}
	}

	if buildRunner == nil {
		t.Error("F-005: dep_dart_build_runner not found")
	} else if buildRunner.scope != "development" {
		t.Errorf("F-005: dep_dart_build_runner scope should be 'development', got %q", buildRunner.scope)
	}
}

// TestFlutterFull_LockfileResolvesVersions verifies pubspec.lock resolved versions
// are merged into dependency nodes.
func TestFlutterFull_LockfileResolvesVersions(t *testing.T) {
	fixturesPath := filepath.Join("testdata", "fixtures", "flutter_full")
	a := NewAnalyzer(fixturesPath, allDetectors())
	graph, err := a.Analyze()
	if err != nil {
		t.Fatalf("Analyze() failed: %v", err)
	}
	dto := graph.ToDTO()

	for _, n := range dto.Nodes {
		if n.ID == "dep_dart_go_router" {
			if n.Resolved == "" {
				t.Error("dep_dart_go_router should have a resolved version from pubspec.lock")
			}
			return
		}
	}
	t.Error("dep_dart_go_router not found in graph")
}

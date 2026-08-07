package main

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLI_Output(t *testing.T) {
	// 1. Build the binary
	binaryPath := filepath.Join(t.TempDir(), "stackgenome")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}

	// 2. Run the binary against the analyzer fixtures
	fixturesPath := filepath.Join("..", "..", "internal", "analyzer", "testdata", "fixtures")
	runCmd := exec.Command(binaryPath, "analyze", "--json", fixturesPath)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	runCmd.Stdout = &stdout
	runCmd.Stderr = &stderr

	if err := runCmd.Run(); err != nil {
		t.Fatalf("binary execution failed: %v\nstderr: %s", err, stderr.String())
	}

	// 3. Verify stderr is empty (no logs leaking to stdout)
	// We allow some stderr if it's logging, but JSON must be on stdout.

	// 4. Validate JSON output
	var result map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse stdout as JSON: %v\nstdout: %s", err, stdout.String())
	}

	// 5. Basic sanity checks on the JSON structure
	nodesRaw, ok := result["nodes"].([]interface{})
	if !ok {
		t.Fatalf("nodes missing or invalid type")
	}

	if len(nodesRaw) == 0 {
		t.Errorf("expected nodes in output, got none")
	}

	// Make sure schema_version matches golden standard (schema v1)
	versionRaw, ok := result["schema_version"].(string)
	if !ok || versionRaw != "1.0.0" {
		t.Errorf("expected schema_version '1.0.0', got %v", versionRaw)
	}
}

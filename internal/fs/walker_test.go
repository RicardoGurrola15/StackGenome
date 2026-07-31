package fs

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestIsSafePath(t *testing.T) {
	tests := []struct {
		root     string
		target   string
		expected bool
	}{
		{"/app", "/app/src/main.go", true},
		{"/app", "/app", true},
		{"/app", "/etc/passwd", false},
		{"/app", "/app/../etc/passwd", false},
		{".", "src/main.go", true},
		{".", "../escape", false},
		{".", "/absolute/escape", false},
	}

	for _, tt := range tests {
		actual := IsSafePath(tt.root, tt.target)
		if actual != tt.expected {
			t.Errorf("IsSafePath(%q, %q): expected %v, got %v", tt.root, tt.target, tt.expected, actual)
		}
	}
}

func TestSafeWalker(t *testing.T) {
	// Create a temporary sandbox structure
	tempDir := t.TempDir()

	// 1. Create a safe file
	err := os.WriteFile(filepath.Join(tempDir, "safe.txt"), []byte("ok"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create an excluded directory
	nodeModules := filepath.Join(tempDir, "node_modules")
	err = os.Mkdir(nodeModules, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(nodeModules, "bad.txt"), []byte("bad"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Create a directory with no permissions
	noPermDir := filepath.Join(tempDir, "noperm")
	err = os.Mkdir(noPermDir, 0000)
	if err != nil {
		t.Fatal(err)
	}

	sw := NewSafeWalker(tempDir)
	var foundFiles []string

	ctx := context.Background()
	err = sw.Walk(ctx, func(path string, d fs.DirEntry) error {
		foundFiles = append(foundFiles, d.Name())
		return nil
	})

	// Restore permissions so cleanup doesn't fail
	os.Chmod(noPermDir, 0755)

	if err != nil {
		t.Fatalf("Walk failed with error: %v", err)
	}

	// Validate results
	foundSafe := false
	for _, f := range foundFiles {
		if f == "safe.txt" {
			foundSafe = true
		}
		if f == "bad.txt" {
			t.Errorf("Walk traversed into excluded directory node_modules")
		}
		if f == "noperm" {
			// It might see the directory itself before trying to enter it, which is fine,
			// but we want to make sure it didn't crash.
		}
	}

	if !foundSafe {
		t.Errorf("Walk did not find the safe.txt file")
	}
}

func BenchmarkSafeWalker(b *testing.B) {
	tempDir := b.TempDir()

	// Create a synthetic tree
	// 10 subdirectories, each with 10 files
	for i := 0; i < 10; i++ {
		dir := filepath.Join(tempDir, "dir"+string(rune('a'+i)))
		os.Mkdir(dir, 0755)
		for j := 0; j < 10; j++ {
			os.WriteFile(filepath.Join(dir, "file"+string(rune('a'+j))+".txt"), []byte("data"), 0644)
		}
	}
	// Add an excluded directory with files
	nodeModules := filepath.Join(tempDir, "node_modules")
	os.Mkdir(nodeModules, 0755)
	for j := 0; j < 10; j++ {
		os.WriteFile(filepath.Join(nodeModules, "file"+string(rune('a'+j))+".txt"), []byte("data"), 0644)
	}

	sw := NewSafeWalker(tempDir)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Walk(ctx, func(path string, d fs.DirEntry) error {
			return nil
		})
	}
}

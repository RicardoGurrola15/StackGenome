package analyzer

import (
	"path/filepath"
	"testing"

	"stackgenome/internal/detectors"
)

func BenchmarkAnalyzer(b *testing.B) {
	// Initialize registry
	dets := detectors.DefaultRegistry()

	// Locate the fixtures path
	fixturesPath := filepath.Join("testdata", "fixtures")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a := NewAnalyzer(fixturesPath, dets)
		_, _ = a.Analyze()
	}
}

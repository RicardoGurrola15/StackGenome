package language

import (
	"testing"
)

func FuzzNodeDetector(f *testing.F) {
	f.Add("package.json", []byte(`{"name":"test","dependencies":{"react":"18.0.0"}}`))
	f.Add("package.json", []byte(`{invalid-json}`))
	f.Add("package-lock.json", []byte(`{"name":"test"}`))
	f.Add("random.json", []byte(`{}`))
	f.Add("package.json", []byte(nil))

	detector := &NodeDetector{}

	f.Fuzz(func(t *testing.T, filename string, content []byte) {
		if !detector.Handles(filename) {
			return
		}
		// We expect this to not panic. Errors are okay.
		_, _, _ = detector.Detect(filename, content)
	})
}

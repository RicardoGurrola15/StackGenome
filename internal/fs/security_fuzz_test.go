package fs

import (
	"testing"
)

func FuzzIsSafePath(f *testing.F) {
	// Add seed corpus
	f.Add("/var/www/project", "/var/www/project/src/main.go")
	f.Add("/var/www/project", "/var/www/project/../passwd")
	f.Add(".", "./src")
	f.Add(".", "../etc")
	f.Add("C:\\project", "C:\\project\\file.txt")
	f.Add("C:\\project", "C:\\Windows\\System32")
	f.Add("", "test.go")
	f.Add("", "../test.go")
	f.Add("/root", "/root")
	f.Add("/root/", "/root/dir")

	f.Fuzz(func(t *testing.T, root string, target string) {
		// Just ensure that the function does not panic for any random inputs.
		// We don't assert boolean logic because Fuzz testing is mainly to prevent crashes (panics).
		// However, we can assert that if root == target (and not empty), it should be true.
		// To avoid path separator nuances, we just check no panics happen.
		_ = IsSafePath(root, target)
	})
}

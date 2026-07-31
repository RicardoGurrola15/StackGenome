package envscanner

import (
	"runtime"
	"testing"
	"time"
)

func TestScanner_OSAndArch(t *testing.T) {
	scanner := NewScanner()
	env := scanner.Scan()

	if env.OS != runtime.GOOS {
		t.Errorf("Expected OS %q, got %q", runtime.GOOS, env.OS)
	}

	if env.Arch != runtime.GOARCH {
		t.Errorf("Expected Arch %q, got %q", runtime.GOARCH, env.Arch)
	}
}

func TestScanner_TimeoutHandling(t *testing.T) {
	// We inject a fake long-running command by replacing the allowlist
	// for the duration of this test.
	originalAllowlist := allowlist
	defer func() { allowlist = originalAllowlist }()

	// "sleep 5" will simulate a command that takes too long
	allowlist = map[string][]string{
		"sleep": {"sleep", "5"},
	}

	scanner := &Scanner{Timeout: 50 * time.Millisecond}
	start := time.Now()

	env := scanner.Scan()

	duration := time.Since(start)

	if duration > 1*time.Second {
		t.Errorf("Scanner did not respect timeout, took %v", duration)
	}

	// Should not have captured sleep output (it was killed or had no output)
	if _, ok := env.Tools["sleep"]; ok {
		t.Errorf("Expected 'sleep' to not be collected due to timeout")
	}
}

func TestScanner_CommandNotFound(t *testing.T) {
	originalAllowlist := allowlist
	defer func() { allowlist = originalAllowlist }()

	allowlist = map[string][]string{
		"nonexistent_tool": {"nonexistent_tool_binary_name", "--version"},
	}

	scanner := NewScanner()
	env := scanner.Scan()

	if _, ok := env.Tools["nonexistent_tool"]; ok {
		t.Errorf("Expected 'nonexistent_tool' to fail silently and not be in tools list")
	}
}

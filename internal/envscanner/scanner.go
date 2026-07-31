package envscanner

import (
	"context"
	"os/exec"
	"runtime"
	"strings"
	"time"

	schemav1 "stackgenome/pkg/schema/v1"
)

// Scanner handles opt-in environment telemetry.
type Scanner struct {
	Timeout time.Duration
}

// NewScanner creates a Scanner with a default timeout of 2 seconds.
func NewScanner() *Scanner {
	return &Scanner{
		Timeout: 2 * time.Second,
	}
}

var allowlist = map[string][]string{
	"go":      {"go", "version"},
	"node":    {"node", "--version"},
	"python":  {"python3", "--version"},
	"rustc":   {"rustc", "--version"},
	"java":    {"java", "-version"},
	"dotnet":  {"dotnet", "--version"},
	"docker":  {"docker", "--version"},
	"flutter": {"flutter", "--version"},
}

// Scan collects the host's OS, architecture, and versions of allowlisted tools.
func (s *Scanner) Scan() *schemav1.EnvironmentDTO {
	env := &schemav1.EnvironmentDTO{
		OS:    runtime.GOOS,
		Arch:  runtime.GOARCH,
		Tools: make(map[string]string),
	}

	for tool, args := range allowlist {
		version := s.runCommand(args[0], args[1:]...)
		if version != "" {
			env.Tools[tool] = version
		}
	}

	return env
}

func (s *Scanner) runCommand(name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Could be command not found, timeout, or non-zero exit. We just ignore it.
		return ""
	}

	// Clean up the output string
	output := strings.TrimSpace(string(out))
	// Return the first line only, to avoid massive multi-line outputs (e.g. java -version)
	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

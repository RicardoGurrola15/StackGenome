package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"stackgenome/internal/analyzer"
	"stackgenome/internal/api"
	"stackgenome/internal/catalog"
	"stackgenome/internal/detectors"
	"stackgenome/internal/envscanner"
	"stackgenome/internal/sanitizer"
)

const analyzeUsage = `Usage: stackgenome analyze [flags] [directory]

Analyzes a software project and produces a normalized stack graph.
Defaults to the current directory if no path is provided.

Flags:
  -recommend   Append top-3 tool recommendations from the embedded local catalog (offline, no network).
  -remote      Fetch recommendations from the Cloudflare API (requires network; implies strict anonymization).
  -env         Also scan the local environment: OS, installed compilers, SDKs (opt-in).
  -json        Output the full ProjectGraph as JSON instead of a human-readable summary.

Examples:
  stackgenome analyze .
  stackgenome analyze /path/to/my/project
  stackgenome analyze -recommend .
  stackgenome analyze -json . > report.json
  stackgenome analyze -recommend -env /path/to/project
`

func RunAnalyze(args []string) error {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, analyzeUsage)
	}
	envFlag := fs.Bool("env", false, "Opt-in to scan the local environment (OS, compilers)")
	recommendFlag := fs.Bool("recommend", false, "Append Top 3 tool recommendations from the local catalog")
	remoteFlag := fs.Bool("remote", false, "Opt-in to fetch recommendations from the remote Cloudflare API (implies strict fingerprinting)")
	jsonFlag := fs.Bool("json", false, "Output results as raw JSON instead of human-readable text")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return fmt.Errorf("invalid flag: %w\n\nRun 'stackgenome analyze --help' for usage", err)
	}

	targetDir := "."
	if fs.NArg() > 0 {
		targetDir = fs.Arg(0)
	}

	info, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory not found: %q\n\nMake sure the path exists and is accessible.", targetDir)
	}
	if err != nil {
		return fmt.Errorf("cannot access directory %q: %w", targetDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%q is a file, not a directory\n\nProvide a project root directory, e.g. stackgenome analyze /path/to/project", targetDir)
	}

	dets := detectors.DefaultRegistry()
	a := analyzer.NewAnalyzer(targetDir, dets)
	graph, err := a.Analyze()
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	dto := graph.ToDTO()

	if *envFlag {
		scanner := envscanner.NewScanner()
		dto.Environment = scanner.Scan()
	}

	// Fetch recommendations
	if *remoteFlag {
		// Enforce full anonymization BEFORE sending data to the remote API
		sanitizer.Anonymize(dto)

		client := api.NewClient("")
		recs, err := client.GetRecommendations(context.Background(), dto, 3)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Remote recommendations failed: %v\n", err)
			fmt.Fprintf(os.Stderr, "   Falling back to local catalog.\n\n")
			// Fallback to local if remote fails
			if cat, catErr := catalog.Load(); catErr == nil {
				ctx := catalog.BuildContext(dto)
				filtered := catalog.Filter(cat.Entries, ctx)
				scored := catalog.Score(filtered, ctx)
				dto.Recommendations = catalog.TopN(scored, 3)
			}
		} else {
			dto.Recommendations = recs
		}
	} else if *recommendFlag {
		// Use local catalog (offline mode)
		cat, err := catalog.Load()
		if err != nil {
			return fmt.Errorf("failed to load local catalog: %w", err)
		}
		ctx := catalog.BuildContext(dto)
		filtered := catalog.Filter(cat.Entries, ctx)
		scored := catalog.Score(filtered, ctx)
		dto.Recommendations = catalog.TopN(scored, 3)
	}

	// Apply final privacy filtering for local output if not already fully anonymized by --remote
	if !*remoteFlag {
		sanitizer.FilterSecrets(dto)
	}

	if *jsonFlag {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(dto); err != nil {
			return fmt.Errorf("failed to encode JSON: %w", err)
		}
	} else {
		printer := NewPrinter(os.Stdout)
		printer.PrintSummary(dto)
	}

	return nil
}

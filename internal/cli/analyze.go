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

func RunAnalyze(args []string) error {
	fs := flag.NewFlagSet("analyze", flag.ExitOnError)
	envFlag := fs.Bool("env", false, "Opt-in to scan the local environment (OS, compilers)")
	recommendFlag := fs.Bool("recommend", false, "Append Top 3 tool recommendations from the local catalog")
	remoteFlag := fs.Bool("remote", false, "Opt-in to fetch recommendations from the remote Cloudflare API (implies strict fingerprinting)")
	jsonFlag := fs.Bool("json", false, "Output results as raw JSON instead of human-readable text")

	if err := fs.Parse(args); err != nil {
		return err
	}

	targetDir := "."
	if fs.NArg() > 0 {
		targetDir = fs.Arg(0)
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("target directory does not exist: %s", targetDir)
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
			fmt.Fprintf(os.Stderr, "⚠️ Remote recommendations failed: %v\n", err)
			// Fallback to local catalog if possible, or just skip
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

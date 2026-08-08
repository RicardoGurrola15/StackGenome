package cli

import (
	"fmt"
	"io"
	"strings"

	"stackgenome/pkg/schema/v1"
)

// Printer is responsible for rendering human-readable output
type Printer struct {
	Out io.Writer
}

func NewPrinter(out io.Writer) *Printer {
	return &Printer{Out: out}
}

// PrintSummary prints a clean, deduplicated summary of the analysis.
func (p *Printer) PrintSummary(dto *schema.ProjectGraphDTO) {
	fmt.Fprintln(p.Out, "🔬 StackGenome Analysis Complete")
	fmt.Fprintln(p.Out, "=================================")

	// Deduplicate languages: keep the highest-confidence node per name.
	langBest := make(map[string]schema.NodeDTO)
	var platforms []string
	var tools []string
	var infra []string
	var frameworks []string

	seenPlatforms := make(map[string]bool)
	seenTools := make(map[string]bool)
	seenInfra := make(map[string]bool)
	seenFrameworks := make(map[string]bool)

	for _, node := range dto.Nodes {
		switch node.Type {
		case "language":
			existing, ok := langBest[node.Name]
			if !ok || node.Confidence > existing.Confidence {
				langBest[node.Name] = node
			}
		case "platform":
			if !seenPlatforms[node.Name] {
				platforms = append(platforms, node.Name)
				seenPlatforms[node.Name] = true
			}
		case "tool":
			if !seenTools[node.Name] {
				tools = append(tools, node.Name)
				seenTools[node.Name] = true
			}
		case "infrastructure":
			if !seenInfra[node.Name] {
				infra = append(infra, node.Name)
				seenInfra[node.Name] = true
			}
		case "framework":
			if !seenFrameworks[node.Name] {
				frameworks = append(frameworks, node.Name)
				seenFrameworks[node.Name] = true
			}
		}
	}

	// Build deduplicated language display (primary first, then secondary).
	// Separate high-confidence (primary/satellite, conf=1.0) from low-confidence.
	var primaryLangs []string
	var secondaryLangs []string
	for _, node := range langBest {
		label := fmt.Sprintf("%s (%.0f%%)", node.Name, node.Confidence*100)
		if node.Confidence >= 0.9 {
			primaryLangs = append(primaryLangs, label)
		} else {
			secondaryLangs = append(secondaryLangs, label)
		}
	}

	fmt.Fprintln(p.Out, "")
	if len(primaryLangs) > 0 {
		fmt.Fprintf(p.Out, "🧬 Primary Languages : %s\n", strings.Join(primaryLangs, ", "))
	}
	if len(secondaryLangs) > 0 {
		fmt.Fprintf(p.Out, "🔹 Also Detected     : %s\n", strings.Join(secondaryLangs, ", "))
	}
	if len(platforms) > 0 {
		fmt.Fprintf(p.Out, "📱 Platforms         : %s\n", strings.Join(platforms, ", "))
	}
	if len(infra) > 0 {
		fmt.Fprintf(p.Out, "☁️  Infrastructure   : %s\n", strings.Join(infra, ", "))
	}
	if len(tools) > 0 {
		fmt.Fprintf(p.Out, "🛠️  Tools Detected   : %s\n", strings.Join(tools, ", "))
	}
	if len(frameworks) > 0 {
		fmt.Fprintf(p.Out, "🔧 Frameworks        : %s\n", strings.Join(frameworks, ", "))
	}

	if len(dto.Recommendations) > 0 {
		fmt.Fprintln(p.Out, "\n✨ Top Recommendations:")
		for i, rec := range dto.Recommendations {
			fmt.Fprintf(p.Out, "  %d. %s (Score: %.2f)\n", i+1, rec.Name, rec.Score)
			fmt.Fprintf(p.Out, "     %s\n", rec.Description)
			if len(rec.Reasons) > 0 {
				fmt.Fprintf(p.Out, "     💡 Why: %s\n", strings.Join(rec.Reasons, " • "))
			}
			if rec.URL != "" {
				fmt.Fprintf(p.Out, "     🔗 %s\n", rec.URL)
			}
			fmt.Fprintln(p.Out, "")
		}
	} else {
		fmt.Fprintln(p.Out, "\nℹ️  No recommendations found for this project's current configuration.")
		fmt.Fprintln(p.Out, "   Run with -recommend to enable the local catalog, or -remote for cloud recommendations.")
	}
}

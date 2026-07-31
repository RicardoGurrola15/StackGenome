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

// PrintSummary prints a summary of the analysis to the terminal.
func (p *Printer) PrintSummary(dto *schema.ProjectGraphDTO) {
	fmt.Fprintln(p.Out, "🔬 StackGenome Analysis Complete")
	fmt.Fprintln(p.Out, "=================================")

	// Count languages and frameworks
	var langs []string
	var frameworks []string
	for _, node := range dto.Nodes {
		if node.Type == "language" {
			langs = append(langs, fmt.Sprintf("%s (%.0f%%)", node.Name, node.Confidence*100))
		}
		if node.Type == "framework" {
			frameworks = append(frameworks, node.Name)
		}
	}

	fmt.Fprintf(p.Out, "\n📦 Languages: %s\n", strings.Join(langs, ", "))
	if len(frameworks) > 0 {
		fmt.Fprintf(p.Out, "🛠️  Frameworks: %s\n", strings.Join(frameworks, ", "))
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
		fmt.Fprintln(p.Out, "\nℹ️  No specific recommendations found.")
	}
}

package catalog

import (
	"math"
	"sort"
	"strings"

	schemav1 "stackgenome/pkg/schema/v1"
)

const (
	weightLanguage  = 0.40
	weightInfra     = 0.30
	weightFramework = 0.20
	weightNovelty   = 0.10
)

// ScoredEntry pairs an Entry with its computed score and structured reasons.
type ScoredEntry struct {
	Entry
	Score   float64
	Reasons []string
}

// Score computes a deterministic relevance score (0.0–1.0) for each entry
// against the project context and returns the results sorted descending.
func Score(entries []Entry, ctx *ProjectContext) []ScoredEntry {
	scored := make([]ScoredEntry, 0, len(entries))

	for _, e := range entries {
		s, reasons := computeScore(e, ctx)
		if s > 0 {
			scored = append(scored, ScoredEntry{
				Entry:   e,
				Score:   round(s, 2),
				Reasons: reasons,
			})
		}
	}

	// Deterministic sort: score DESC, then ID ASC as tiebreaker
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].Score != scored[j].Score {
			return scored[i].Score > scored[j].Score
		}
		return scored[i].ID < scored[j].ID
	})

	return scored
}

// TopN returns the top n recommendations from a scored list as DTOs.
func TopN(scored []ScoredEntry, n int) []schemav1.RecommendationDTO {
	if n > len(scored) {
		n = len(scored)
	}
	out := make([]schemav1.RecommendationDTO, n)
	for i := 0; i < n; i++ {
		out[i] = schemav1.RecommendationDTO{
			ID:          scored[i].ID,
			Name:        scored[i].Name,
			Description: scored[i].Description,
			Score:       scored[i].Score,
			Reasons:     scored[i].Reasons,
			URL:         scored[i].URL,
		}
	}
	return out
}

func computeScore(e Entry, ctx *ProjectContext) (float64, []string) {
	var score float64
	var reasons []string

	// Language match
	matchedLangs := intersect(e.Targets.Languages, ctx.Languages)
	if len(matchedLangs) > 0 {
		score += weightLanguage
		reasons = append(reasons, "lenguaje compatible: "+strings.Join(matchedLangs, ", "))
	}

	// Infra match
	matchedInfra := intersect(e.Targets.Infra, ctx.Infra)
	if len(matchedInfra) > 0 {
		score += weightInfra
		reasons = append(reasons, "infraestructura compatible: "+strings.Join(matchedInfra, ", "))
	}

	// Framework match
	matchedFW := intersect(e.Targets.Frameworks, ctx.Frameworks)
	if len(matchedFW) > 0 {
		score += weightFramework
		reasons = append(reasons, "framework compatible: "+strings.Join(matchedFW, ", "))
	}

	// Novelty bonus: tool not already detected in the graph
	if !ctx.NodeIDs[e.ID] {
		score += weightNovelty
		reasons = append(reasons, "herramienta no detectada aún en el proyecto")
	}

	// Universal tools (no targets specified) get a base relevance score
	if len(e.Targets.Languages) == 0 && len(e.Targets.Infra) == 0 && len(e.Targets.Frameworks) == 0 {
		if score < weightNovelty {
			score = weightNovelty
			reasons = append(reasons, "herramienta de uso general")
		}
	}

	return score, reasons
}

func intersect(targets []string, present map[string]bool) []string {
	var matched []string
	for _, t := range targets {
		if present[strings.ToLower(t)] {
			matched = append(matched, t)
		}
	}
	return matched
}

func round(val float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(val*p) / p
}

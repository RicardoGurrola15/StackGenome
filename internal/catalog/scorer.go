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

	// primaryLanguageBoost is added when the tool's ecosystem matches the
	// project's primary language (highest-confidence language node).
	primaryLanguageBoost = 0.15

	// minScoreThreshold is the minimum score an entry must reach to be
	// included in recommendations. Entries below this are silently dropped,
	// enabling the engine to abstain when no relevant tool is found.
	minScoreThreshold = 0.30
)

// ScoredEntry pairs an Entry with its computed score and structured reasons.
type ScoredEntry struct {
	Entry
	Score   float64
	Reasons []string
}

// Score computes a deterministic relevance score (0.0–1.0) for each entry
// against the project context and returns the results sorted descending.
// Only entries at or above minScoreThreshold are included.
func Score(entries []Entry, ctx *ProjectContext) []ScoredEntry {
	scored := make([]ScoredEntry, 0, len(entries))

	for _, e := range entries {
		s, reasons := computeScore(e, ctx)
		if s >= minScoreThreshold {
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
// If the scored list is empty, it returns an empty slice (never nil).
func TopN(scored []ScoredEntry, n int) []schemav1.RecommendationDTO {
	if len(scored) == 0 {
		return []schemav1.RecommendationDTO{}
	}
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

	// Ecosystem match (Languages)
	matchedLangs := intersect(e.Ecosystem, ctx.Languages)
	if len(matchedLangs) > 0 {
		// Filter out "*" wildcard from display
		var realMatches []string
		for _, m := range matchedLangs {
			if m != "*" {
				realMatches = append(realMatches, m)
			}
		}
		if len(realMatches) > 0 {
			score += weightLanguage
			reasons = append(reasons, "ecosistema compatible: "+strings.Join(realMatches, ", "))

			// Primary language boost: if the tool matches the project's primary
			// language (highest-confidence language node), add an extra bonus.
			for _, m := range realMatches {
				if ctx.PrimaryLanguage != "" && strings.EqualFold(m, ctx.PrimaryLanguage) {
					score += primaryLanguageBoost
					reasons = append(reasons, "lenguaje principal del proyecto: "+ctx.PrimaryLanguage)
					break
				}
			}
		}
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

	// Universal tools (no specific targets) get a base relevance score only
	// if they haven't already accumulated points from specific matches.
	if (len(e.Ecosystem) == 0 || (len(e.Ecosystem) == 1 && e.Ecosystem[0] == "*")) &&
		len(e.Targets.Infra) == 0 && len(e.Targets.Frameworks) == 0 {
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
